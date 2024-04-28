package main

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"go-pub-sub/internal/console"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	console.InitServer()
}
