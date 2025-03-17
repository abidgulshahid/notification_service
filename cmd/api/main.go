package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/notification_service/internal/config"
	"github.com/notification_service/internal/email"
	"github.com/notification_service/internal/kafka"
	"github.com/notification_service/internal/notifications"
	"github.com/notification_service/internal/supabase"
	"github.com/notification_service/internal/telegram"
)

func main() {
	log.Println("Starting notification service...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create Supabase client
	supabaseClient, err := supabase.NewClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create Supabase client: %v", err)
	}

	// Create SendGrid client
	var emailClient *email.SendGridClient
	if cfg.SendGrid.APIKey != "" {
		emailClient, err = email.NewSendGridClient(cfg)
		if err != nil {
			log.Printf("Warning: Failed to create SendGrid client: %v", err)
		}
	} else {
		log.Println("Warning: SendGrid API key not provided, email notifications will not be available")
	}

	// Create Telegram client
	var telegramClient *telegram.TelegramClient
	if cfg.Telegram.BotToken != "" {
		telegramClient, err = telegram.NewTelegramClient(cfg)
		if err != nil {
			log.Printf("Warning: Failed to create Telegram client: %v", err)
		} else {
			// Start the Telegram bot in the background
			go telegramClient.StartBot()
			defer telegramClient.StopBot()
		}
	} else {
		log.Println("Warning: Telegram bot token not provided, Telegram notifications will not be available")
	}

	// Create notification service
	notificationService := notifications.NewService(supabaseClient, emailClient, telegramClient)

	// Create Kafka consumer
	consumer, err := kafka.NewConsumer(cfg, notificationService.ProcessNotification)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Kafka consumer
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Wait for termination signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Println("Received termination signal, shutting down...")
	
	// Stop Kafka consumer
	consumer.Stop()
	
	log.Println("Notification service stopped")
} 