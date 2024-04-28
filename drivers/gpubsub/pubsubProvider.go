package gpubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"go-pub-sub/drivers/gpubsub/interfaces"
	"go-pub-sub/internal/config"
	"sync"
)

type PubSubProvider struct {
	Client *pubsub.Client
}

// NewPubSubProvider is function to create new instance of pubsubProvider
func NewPubSubProvider(ctx context.Context) interfaces.IPubSubProvider {
	client, err := pubsub.NewClient(ctx, config.PubSubProjectId())
	if err != nil {
		logrus.Fatal(err)
	}

	logrus.Info("success create google pubsub client 💬✅")
	return &PubSubProvider{
		Client: client,
	}
}

// Publish is method to publish new message to google pub sub
func (p *PubSubProvider) Publish(ctx context.Context, topic string, message []byte) (string, error) {
	t, err := p.CreateTopicIfNotExist(ctx, topic)
	if err != nil {
		return "", err
	}

	msg := pubsub.Message{
		Data: message,
	}
	res := t.Publish(ctx, &msg)

	wg := &sync.WaitGroup{}
	chanErr := make(chan error, 1)
	chanId := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()

		id, err := res.Get(ctx)
		if err != nil {
			chanId <- ""
			chanErr <- err
			return
		}

		chanId <- id
		chanErr <- nil
		logrus.Info("message success published")
	}()

	wg.Wait()
	return <-chanId, <-chanErr
}

// CreateTopicIfNotExist is method to create topic if not exist
func (p *PubSubProvider) CreateTopicIfNotExist(ctx context.Context, topic string) (*pubsub.Topic, error) {
	t := p.Client.Topic(topic)
	exists, _ := t.Exists(ctx)
	if !exists {
		createTopic, err := p.Client.CreateTopic(ctx, topic)
		if err != nil {
			return nil, err
		}

		return createTopic, nil
	}

	return t, nil
}

func (p *PubSubProvider) ShutDown() error {
	if p.Client != nil {
		if err := p.Client.Close(); err != nil {
			logrus.Fatalf("force close pubsub client 🔴")
			return err
		}
	}

	logrus.Info("success close pubsub client ⚠️✅")
	return nil
}

func (p *PubSubProvider) CreateSubsriberIfNotExist(ctx context.Context, subcriberName string, topic string) (*pubsub.Subscription, error) {
	var sub *pubsub.Subscription
	sub = p.Client.Subscription(subcriberName)
	exists, _ := sub.Exists(ctx)

	if !exists {
		t, _ := p.CreateTopicIfNotExist(ctx, topic)
		sub, err := p.Client.CreateSubscription(ctx, subcriberName, pubsub.SubscriptionConfig{Topic: t})
		if err != nil {
			return nil, err
		}

		return sub, nil
	}

	return sub, nil
}

// Subscribe is method to subsribe specific topic
func (p *PubSubProvider) Subscribe(ctx context.Context, topic string, subscribeName string) {
	sub, _ := p.CreateSubsriberIfNotExist(ctx, subscribeName, topic)

	go func() {
		logrus.Infof("[%v] starting subscribe topic [%v] ⌛️", subscribeName, topic)
		err := sub.Receive(ctx, func(ctx context.Context, message *pubsub.Message) {
			msg := fmt.Sprintf("receive message 💬✅ : %v", string(message.Data))
			logrus.Info(msg)
			message.Ack()
		})
		if err != nil {
			logrus.Errorf("cant receive message 💬⚠️ : %v")
		}
	}()
}
