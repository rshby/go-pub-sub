package service

import (
	"context"
	gpubsubInterface "go-pub-sub/drivers/gpubsub/interfaces"
	"go-pub-sub/internal/config"
	"go-pub-sub/internal/service/interfaces"
)

type ChatService struct {
	PubSubProvider gpubsubInterface.IPubSubProvider
}

// NewChatService is function to create new instance chatService
func NewChatService(pubsubProvider gpubsubInterface.IPubSubProvider) interfaces.IChatService {
	return &ChatService{PubSubProvider: pubsubProvider}
}

// SendChat is method to send chat
func (c *ChatService) SendChat(ctx context.Context, message string) (string, error) {
	return c.PubSubProvider.Publish(ctx, config.DefaultTopicName, []byte(message))
}
