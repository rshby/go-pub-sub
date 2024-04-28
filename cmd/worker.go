package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go-pub-sub/drivers/gpubsub"
	"go-pub-sub/internal/config"
	"os"
	"os/signal"
)

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	logrus.SetFormatter(customFormatter)
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	// initialise db

	// initialise worker
	ctx := context.Background()
	pubsubProvider := gpubsub.NewPubSubProvider(ctx)

	// subscribe
	ctxSub, cancel := context.WithCancel(ctx)
	go pubsubProvider.Subscribe(ctxSub, config.PubSubTopicName(), "mobile-1")

	// signal
	signalChan := make(chan os.Signal, 1)
	quitChan := make(chan bool, 1)

	signal.Notify(signalChan, os.Interrupt)

	go func() {
		select {
		case <-signalChan:
			logrus.Info("receive interrupt signal ⚠️")
			cancel()
			pubsubProvider.ShutDown()
			quitChan <- true
		}
	}()

	<-quitChan
	logrus.Info("worker exit 🔴")
}
