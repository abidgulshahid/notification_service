package notifications

import (
	"fmt"
	"log"

	"github.com/notification_service/internal/email"
	"github.com/notification_service/internal/models"
	"github.com/notification_service/internal/supabase"
	"github.com/notification_service/internal/telegram"
)

// Service handles notification processing
type Service struct {
	supabaseClient *supabase.Client
	emailClient    *email.SendGridClient
	telegramClient *telegram.TelegramClient
}

// NewService creates a new notification service
func NewService(
	supabaseClient *supabase.Client,
	emailClient *email.SendGridClient,
	telegramClient *telegram.TelegramClient,
) *Service {
	return &Service{
		supabaseClient: supabaseClient,
		emailClient:    emailClient,
		telegramClient: telegramClient,
	}
}

// ProcessNotification processes a notification message from Kafka
func (s *Service) ProcessNotification(msg *models.KafkaNotificationMessage) error {
	log.Printf("Processing notification for user %s of type %s", msg.UserID, msg.Type)

	// Create notification record
	notification := &models.Notification{
		UserID:   msg.UserID,
		Type:     msg.Type,
		Channel:  msg.Channel,
		Subject:  msg.Subject,
		Content:  msg.Content,
		Status:   models.NotificationStatusPending,
		Metadata: msg.Metadata,
	}

	// Insert notification into Supabase
	id, err := s.supabaseClient.InsertNotification(notification)
	if err != nil {
		return fmt.Errorf("failed to insert notification: %w", err)
	}

	notification.ID = id
	log.Printf("Notification inserted with ID: %s", id)

	// Send notification based on type
	var sendErr error
	switch notification.Type {
	case models.NotificationTypeEmail:
		sendErr = s.sendEmailNotification(notification)
	case models.NotificationTypeTelegram:
		sendErr = s.sendTelegramNotification(notification)
	default:
		sendErr = fmt.Errorf("unsupported notification type: %s", notification.Type)
	}

	// Update notification status
	var status models.NotificationStatus
	if sendErr != nil {
		log.Printf("Failed to send notification: %v", sendErr)
		status = models.NotificationStatusFailed
	} else {
		log.Printf("Notification sent successfully")
		status = models.NotificationStatusSent
	}

	if err := s.supabaseClient.UpdateNotificationStatus(id, status); err != nil {
		log.Printf("Failed to update notification status: %v", err)
	}

	return sendErr
}

// sendEmailNotification sends an email notification
func (s *Service) sendEmailNotification(notification *models.Notification) error {
	if s.emailClient == nil {
		return fmt.Errorf("email client not configured")
	}

	log.Printf("Sending email notification to %s", notification.Channel)
	return s.emailClient.SendEmail(notification)
}

// sendTelegramNotification sends a Telegram notification
func (s *Service) sendTelegramNotification(notification *models.Notification) error {
	if s.telegramClient == nil {
		return fmt.Errorf("telegram client not configured")
	}

	log.Printf("Sending telegram notification to %s", notification.Channel)
	return s.telegramClient.SendNotification(notification)
} 