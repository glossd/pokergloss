package gomq

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/glossd/pokergloss/auth"
	conf "github.com/glossd/pokergloss/goconf"
	log "github.com/sirupsen/logrus"
	"time"
)

func Pull(subID string, topicID string, receiver func(ctx context.Context, msg *pubsub.Message) error) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := pubsub.NewClient(ctx, conf.Props.GCP.ProjectID, auth.GoogleClientOptions()...)
	if err != nil {
		return fmt.Errorf("pubsub.NewClient: %v", err)
	}

	topic := client.Topic(topicID)
	exists, err := topic.Exists(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.Topic.Exists failed: %v", err)
	}
	if !exists {
		log.Printf("Topic %v doesn't exist - creating it", topicID)
		_, err = client.CreateTopic(ctx, topicID)
		if err != nil {
			return fmt.Errorf("failed to create topic:: %v", err)
		}
	}

	sub := client.Subscription(subID)
	exists, err = sub.Exists(ctx)
	if err != nil {
		return fmt.Errorf("pubsub.Subscription.Exists failed: %v", err)
	}
	if !exists {
		log.Printf("Subscription %v doesn't exist - creating it", subID)
		config := pubsub.SubscriptionConfig{Topic: topic, AckDeadline: 20 * time.Second}
		_, err = client.CreateSubscription(ctx, subID, config)
		if err != nil {
			return fmt.Errorf("failed to create subscription: %v", err)
		}
	}

	err = sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		err := receiver(ctx, msg)
		if err != nil {
			if IsAckableError(err) {
				msg.Ack()
				return
			}
			msg.Nack()
			return
		}
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("pubsub receive error: %s", err)
	}
	return nil
}
