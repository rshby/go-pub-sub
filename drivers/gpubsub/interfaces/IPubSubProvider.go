package interfaces

import (
	"cloud.google.com/go/pubsub"
	"context"
)

type IPubSubProvider interface {
	// Publish is method to publish message
	Publish(ctx context.Context, topic string, message []byte) (string, error)

	// CreateTopicIfNotExist is method to create topic if not exist and return topic if exist
	CreateTopicIfNotExist(ctx context.Context, topic string) (*pubsub.Topic, error)

	// CreateSubsriberIfNotExist is method to create subsriber if not exist
	CreateSubsriberIfNotExist(ctx context.Context, subcriberName string, topic string) (*pubsub.Subscription, error)

	// Subscribe is method to subscribe specific topic
	Subscribe(ctx context.Context, topic string, subscribeName string)

	// ShutDown is method to shutdown and close pubsub client
	ShutDown() error
}
