package config

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

// GetEnv is method to get env
func GetEnv(key string) string {
	if err := godotenv.Load(".env"); err != nil {
		logrus.Fatal(err)
	}

	return os.Getenv(key)
}

// PubSubProjectId is fuction to get project id in google pub sub
func PubSubProjectId() string {
	return GetEnv("PUBSUB_PROJECT_ID")
}

// PubSubTopicName is method to get pub sub topic name
func PubSubTopicName() string {
	if topic := GetEnv("PUBSUB_TOPIC_NAME"); topic != "" {
		return topic
	}

	return DefaultTopicName
}
