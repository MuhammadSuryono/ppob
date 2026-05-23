package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventPublisher struct {
	client *redis.Client
}

func NewEventPublisher(client *redis.Client) *EventPublisher {
	return &EventPublisher{client: client}
}

func (p *EventPublisher) Publish(ctx context.Context, stream string, eventType string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	err = p.client.XAdd(ctx, &redis.XAddArgs{
		Stream: stream,
		Values: map[string]interface{}{
			"type":      eventType,
			"payload":   string(payload),
			"timestamp": time.Now().Unix(),
		},
	}).Err()

	if err != nil {
		return fmt.Errorf("failed to publish event to Redis Streams: %w", err)
	}

	return nil
}

type EventHandler func(ctx context.Context, eventType string, payload string) error

type EventConsumer struct {
	client   *redis.Client
	stream   string
	group    string
	consumer string
}

func NewEventConsumer(client *redis.Client, stream, group, consumer string) *EventConsumer {
	return &EventConsumer{
		client:   client,
		stream:   stream,
		group:    group,
		consumer: consumer,
	}
}

func (c *EventConsumer) Start(ctx context.Context, handler EventHandler) error {
	// Ensure group exists
	err := c.client.XGroupCreateMkStream(ctx, c.stream, c.group, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return fmt.Errorf("failed to create consumer group: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.group,
				Consumer: c.consumer,
				Streams:  []string{c.stream, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err != nil {
				if err != redis.Nil {
					log.Printf("Error reading from Redis Streams: %v", err)
				}
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					eventType, _ := msg.Values["type"].(string)
					payload, _ := msg.Values["payload"].(string)

					if err := handler(ctx, eventType, payload); err != nil {
						log.Printf("Error handling event %s: %v", msg.ID, err)
						// Depending on strategy, we might not ACK here to retry later
						continue
					}

					// Acknowledge message
					c.client.XAck(ctx, c.stream, c.group, msg.ID)
				}
			}
		}
	}
}
