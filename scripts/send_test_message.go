package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/joho/godotenv"
)

// TestMessage represents a test notification message
type TestMessage struct {
	UserID   string            `json:"user_id"`
	Type     string            `json:"type"`
	Channel  string            `json:"channel"`
	Subject  string            `json:"subject"`
	Content  string            `json:"content"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

func main() {
	// Load .env file if exists
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	// Define command-line flags
	msgType := flag.String("type", "email", "Type of notification (email or telegram)")
	channel := flag.String("channel", "", "Email address or Telegram chat ID")
	subject := flag.String("subject", "Test Notification", "Notification subject")
	content := flag.String("content", "This is a test notification from Kafka.", "Notification content")
	userID := flag.String("user", "user-123", "User ID")
	
	bootstrapServers := flag.String("bootstrap-servers", getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"), "Kafka bootstrap servers")
	topic := flag.String("topic", getEnv("KAFKA_TOPIC", "notifications"), "Kafka topic")
	
	flag.Parse()
	
	if *channel == "" {
		log.Fatal("Channel (email or telegram chat ID) is required")
	}
	
	// Create message
	message := TestMessage{
		UserID:   *userID,
		Type:     *msgType,
		Channel:  *channel,
		Subject:  *subject,
		Content:  *content,
		Metadata: map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
			"source":    "test-script",
		},
	}
	
	// Convert message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Failed to marshal message to JSON: %v", err)
	}
	
	// Create Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": *bootstrapServers,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %v", err)
	}
	defer p.Close()
	
	// Handle delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v", ev.TopicPartition.Error)
				} else {
					log.Printf("Successfully produced message to topic %s (partition %d at offset %d)",
						*ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
				}
			}
		}
	}()
	
	// Produce message
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: topic, Partition: kafka.PartitionAny},
		Value:          messageJSON,
	}, nil)
	if err != nil {
		log.Fatalf("Failed to produce message: %v", err)
	}
	
	// Wait for message delivery
	p.Flush(15 * 1000)
	
	fmt.Println("Message sent successfully!")
	fmt.Printf("Message: %s\n", string(messageJSON))
}

// getEnv retrieves environment variables with fallback to default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 