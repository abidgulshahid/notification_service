package telegram

import (
	"errors"
	"fmt"
	"time"

	"github.com/notification_service/internal/config"
	"github.com/notification_service/internal/models"
	telebot "gopkg.in/telebot.v3"
)

// TelegramClient represents a Telegram client
type TelegramClient struct {
	bot *telebot.Bot
}

// NewTelegramClient creates a new Telegram client
func NewTelegramClient(cfg *config.Config) (*TelegramClient, error) {
	if cfg.Telegram.BotToken == "" {
		return nil, errors.New("Telegram bot token must be provided")
	}

	bot, err := telebot.NewBot(telebot.Settings{
		Token:  cfg.Telegram.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Telegram bot: %w", err)
	}

	return &TelegramClient{
		bot: bot,
	}, nil
}

// StartBot starts the Telegram bot
func (t *TelegramClient) StartBot() {
	t.bot.Start()
}

// StopBot stops the Telegram bot
func (t *TelegramClient) StopBot() {
	t.bot.Stop()
}

// SendNotification sends a notification to a Telegram chat
func (t *TelegramClient) SendNotification(notification *models.Notification) error {
	// Check if channel is provided
	if notification.Channel == "" {
		return errors.New("telegram channel ID must be provided")
	}

	// Create a recipient from the channel (chat ID)
	recipient := &telebot.Chat{ID: parseChatID(notification.Channel)}

	// Create a message
	message := notification.Content
	if notification.Subject != "" {
		message = fmt.Sprintf("*%s*\n\n%s", notification.Subject, notification.Content)
	}

	// Send the message
	_, err := t.bot.Send(recipient, message, &telebot.SendOptions{
		ParseMode: telebot.ModeMarkdown,
	})
	if err != nil {
		return fmt.Errorf("failed to send telegram message: %w", err)
	}

	return nil
}

// parseChatID converts a chat ID from string to int64
func parseChatID(chatID string) int64 {
	var id int64
	fmt.Sscanf(chatID, "%d", &id)
	return id
} 