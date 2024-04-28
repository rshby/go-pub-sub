package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
	"time"
)

type UserLimiter struct {
	Limiter  *rate.Limiter
	LastSeen time.Time
}

var (
	globalLimiter = rate.NewLimiter(300, 300)
	Users         = map[string]*UserLimiter{}
	mu            = &sync.RWMutex{}
)

// RateLimiterMiddleware is function to get limiter middleware
func RateLimiterMiddleware(redis *redis.Client) gin.HandlerFunc {
	store, rate := InitRateLimit(redis)
	middleware := mgin.NewMiddleware(limiter.New(store, rate))
	return middleware
}

func InitRateLimit(redisClient *redis.Client) (limiter.Store, limiter.Rate) {
	rate, err := limiter.NewRateFromFormatted("5-S")
	if err != nil {
		logrus.Fatalf("cant retrive rate : %v", err)
	}

	store, err := sredis.NewStoreWithOptions(redisClient, limiter.StoreOptions{
		Prefix: "limiter_gin",
	})

	return store, rate
}

// RateLimitMiddlewareTimeRate is function to get gin middleware rate limit with package time rate
func RateLimitMiddlewareTimeRate() gin.HandlerFunc {
	limiter := rate.NewLimiter(5, 2)
	return func(context *gin.Context) {
		if !limiter.Allow() {
			statusCode := http.StatusTooManyRequests
			context.JSON(statusCode, gin.H{
				"message": "too many request",
			})
			context.Abort()
			return
		}

		context.Next()
		return
	}
}

// RateLimitByUser is function to get gin middleware limit by client
func RateLimitByUser() gin.HandlerFunc {
	type Client struct {
		Limiter  *rate.Limiter
		LastSeen time.Time
	}

	mu := &sync.Mutex{}
	clients := map[string]*Client{}

	limiter := rate.NewLimiter(100, 100)

	go func() {
		for {
			time.Sleep(50 * time.Second)
			mu.Lock()
			for key, value := range clients {
				if time.Since(value.LastSeen) > 1*time.Second {
					delete(clients, key)
				}
			}
			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		if !limiter.Allow() {
			statusCode := http.StatusTooManyRequests
			ctx.JSON(statusCode, gin.H{
				"message": "too many request",
			})
			ctx.Abort()
			return
		}

		// cek limit by user
		mu.Lock()
		ip := ctx.RemoteIP()
		if _, ok := clients[ip]; !ok {
			clients[ip] = &Client{
				Limiter: rate.NewLimiter(1, 1),
			}
		}

		clients[ip].LastSeen = time.Now()

		if !clients[ip].Limiter.Allow() {
			mu.Unlock()
			logrus.Info("kena sini per user")
			statusCode := http.StatusTooManyRequests
			ctx.JSON(statusCode, gin.H{
				"message": "too many request",
			})
			ctx.Abort()
			return
		}

		mu.Unlock()
		ctx.Next()
		return
	}
}

func RateLimiterByUserMiddlewareV2() gin.HandlerFunc {
	go CleanUpUsers()

	return func(ctx *gin.Context) {
		if !globalLimiter.Allow() {
			statusCode := http.StatusTooManyRequests
			ctx.JSON(statusCode, gin.H{
				"message": "too many request",
			})
			ctx.Abort()
			return
		}

		// get limiter by IP
		ip := ctx.ClientIP()
		userLimiter := GetUserByIp(ip)
		if !userLimiter.Limiter.Allow() {
			statusCode := http.StatusTooManyRequests
			ctx.JSON(statusCode, gin.H{
				"message": "too many request per user",
			})
			ctx.Abort()
			return
		}

		// pass the limiter
		ctx.Next()
		return
	}
}

// GetUserByIp is function to get limiter by user ip
// if not exist then create new one
func GetUserByIp(ip string) *UserLimiter {
	mu.RLock()
	if user, ok := Users[ip]; ok {
		mu.RUnlock()
		return user
	}

	mu.RUnlock()
	mu.Lock()
	userLimiter := UserLimiter{
		Limiter:  rate.NewLimiter(1, 1),
		LastSeen: time.Now(),
	}
	Users[ip] = &userLimiter
	mu.Unlock()
	return &userLimiter
}

// CleanUpUsers is function to cleanup user ip and limiter data
func CleanUpUsers() {
	for {
		time.Sleep(10 * time.Second)
		mu.Lock()
		for ip, x := range Users {
			if time.Since(x.LastSeen) >= 4*time.Second {
				delete(Users, ip)
			}
		}
		mu.Unlock()
	}
}
