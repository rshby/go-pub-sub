package console

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"go-pub-sub/drivers/gpubsub"
	"go-pub-sub/internal/controller"
	"go-pub-sub/internal/middleware"
	"go-pub-sub/internal/service"
	"net/http"
	"os"
	"os/signal"
)

// InitServer is method to initialize server service
func InitServer() {
	ctx := context.Background()
	pubsubProvider := gpubsub.NewPubSubProvider(ctx)

	// register service
	chatService := service.NewChatService(pubsubProvider)

	// register controller
	chatController := controller.NewChatController(chatService)

	// init redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: ":6379",
		DB:   3,
	})
	ping, _ := redisClient.Ping(ctx).Result()
	logrus.Info(ping)

	//limiterMiddleware := middleware.RateLimiterMiddleware(redisClient)
	//limiterMiddleware := middleware.RateLimitMiddlewareTimeRate()
	//limiterMiddleware := middleware.RateLimitByUser()
	ctx, cancel := context.WithCancel(ctx)
	limiterMiddleware := middleware.RateLimiterByUserMiddlewareV2(ctx)

	// add gin
	app := gin.Default()
	app.ForwardedByClientIP = true
	app.Use(limiterMiddleware)

	// set routing
	SetRouting(app.RouterGroup, chatController)

	httpServer := http.Server{
		Addr:    ":4000",
		Handler: app,
	}

	signalChan := make(chan os.Signal, 1)
	errChan := make(chan error, 1)
	quitChan := make(chan bool, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	go func() {
		for {
			select {
			case <-signalChan:
				logrus.Info("receive signal interrupt ⚠️")
				cancel()
				_ = pubsubProvider.ShutDown()
				ShutdownServer(ctx, &httpServer)
				quitChan <- true
			case <-errChan:
				quitChan <- true
			}
		}
	}()

	<-quitChan
	logrus.Info("server exit ❌")
}

// SetRouting is function to set routing endpoint
func SetRouting(app gin.RouterGroup, controller *controller.ChatController) {
	chat := app.Group("/chat")
	{
		chat.POST("/send", controller.Send)
	}
}

func ShutdownServer(ctx context.Context, server *http.Server) {
	if server != nil {
		if err := server.Shutdown(ctx); err != nil {
			logrus.Fatalf("force close server 🔴")
		}

		logrus.Info("server sucess to shutdown ⚠️✅")
	}
}
