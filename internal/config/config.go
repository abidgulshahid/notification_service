package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Kafka    KafkaConfig
	Supabase SupabaseConfig
	SendGrid SendGridConfig
	Telegram TelegramConfig
}

type KafkaConfig struct {
	BootstrapServers string
	Topic            string
}

type SupabaseConfig struct {
	URL              string
	APIKey           string
	NotificationsTable string
}

type SendGridConfig struct {
	APIKey     string
	FromEmail  string
	FromName   string
}

type TelegramConfig struct {
	BotToken string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Load .env file if exists
	err := godotenv.Load()
	if err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	viper.AutomaticEnv()

	config := &Config{
		Kafka: KafkaConfig{
			BootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "localhost:9092"),
			Topic:            getEnv("KAFKA_TOPIC", "notifications"),
		},
		Supabase: SupabaseConfig{
			URL:                getEnv("SUPABASE_URL", ""),
			APIKey:             getEnv("SUPABASE_API_KEY", ""),
			NotificationsTable: getEnv("SUPABASE_NOTIFICATIONS_TABLE", "notifications"),
		},
		SendGrid: SendGridConfig{
			APIKey:    getEnv("SENDGRID_API_KEY", ""),
			FromEmail: getEnv("SENDGRID_FROM_EMAIL", ""),
			FromName:  getEnv("SENDGRID_FROM_NAME", "Notification Service"),
		},
		Telegram: TelegramConfig{
			BotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
		},
	}

	return config, nil
}

// getEnv retrieves environment variables with fallback to default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
} 