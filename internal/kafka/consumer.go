package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/notification_service/internal/config"
	"github.com/notification_service/internal/models"
)

// MessageHandler is a function that processes Kafka messages
type MessageHandler func(msg *models.KafkaNotificationMessage) error

// Consumer represents a Kafka consumer
type Consumer struct {
	consumer *kafka.Consumer
	topic    string
	handler  MessageHandler
	running  bool
}

// NewConsumer creates a new Kafka consumer
func NewConsumer(cfg *config.Config, handler MessageHandler) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":     cfg.Kafka.BootstrapServers,
		"group.id":              "notification-service",
		"auto.offset.reset":     "earliest",
		"enable.auto.commit":    true,
		"auto.commit.interval.ms": 5000,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	return &Consumer{
		consumer: c,
		topic:    cfg.Kafka.Topic,
		handler:  handler,
	}, nil
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	if c.running {
		return fmt.Errorf("consumer is already running")
	}

	err := c.consumer.SubscribeTopics([]string{c.topic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topics: %w", err)
	}

	c.running = true
	log.Printf("Kafka consumer started, listening to topic: %s", c.topic)

	go c.consume(ctx)
	return nil
}

// consume is the main message processing loop
func (c *Consumer) consume(ctx context.Context) {
	for c.running {
		select {
		case <-ctx.Done():
			c.Stop()
			return
		default:
			msg, err := c.consumer.ReadMessage(100 * time.Millisecond)
			if err != nil {
				// Timeout or no message available is not an error
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					log.Printf("Error reading message: %v", err)
				}
				continue
			}

			log.Printf("Received message from topic %s [%d] at offset %d: %s",
				*msg.TopicPartition.Topic, msg.TopicPartition.Partition, msg.TopicPartition.Offset, string(msg.Value))

			var notification models.KafkaNotificationMessage
			if err := json.Unmarshal(msg.Value, &notification); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				continue
			}

			if err := c.handler(&notification); err != nil {
				log.Printf("Error handling message: %v", err)
			}
		}
	}
}

// Stop stops the consumer
func (c *Consumer) Stop() {
	if !c.running {
		return
	}

	c.running = false
	if err := c.consumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	log.Println("Kafka consumer stopped")
} 