package gomq

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/glossd/pokergloss/auth"
	"github.com/glossd/pokergloss/goconf"
	"github.com/golang/protobuf/proto"
	"time"
)

var ErrPublisherNotReady = errors.New("publisher is not ready")

type Publisher struct {
	c *pubsub.Client
}

func InitPublisher() (*Publisher, error) {
	pubsubClient, err := pubsub.NewClient(context.Background(), goconf.Props.GCP.ProjectID, auth.GoogleClientOptions()...)
	if err != nil {
		return nil, fmt.Errorf("pubsub.NewClient: %v", err)
	}
	return &Publisher{c: pubsubClient}, nil
}

func (p *Publisher) Publish(topicID string, msg proto.Message) error {
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal msg: %v", err)
	}
	return p.publish(topicID, msgBytes)
}

func (p *Publisher) PublishJSON(topicID string, msg interface{}) error {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal msg: %v", err)
	}
	return p.publish(topicID, msgBytes)
}

func (p *Publisher) publish(topicID string, msg []byte) error {
	if p.c == nil {
		return fmt.Errorf("pubsub client is not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	t := p.c.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{Data: msg})
	// Block until the result is returned and a server-generated
	// ID is returned for the published message.
	_, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("publish get: %v", err)
	}
	return nil
}
