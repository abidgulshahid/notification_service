package email

import (
	"errors"
	"fmt"

	"github.com/notification_service/internal/config"
	"github.com/notification_service/internal/models"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendGridClient represents a SendGrid email client
type SendGridClient struct {
	client    *sendgrid.Client
	fromEmail string
	fromName  string
}

// NewSendGridClient creates a new SendGrid client
func NewSendGridClient(cfg *config.Config) (*SendGridClient, error) {
	if cfg.SendGrid.APIKey == "" {
		return nil, errors.New("SendGrid API key must be provided")
	}

	if cfg.SendGrid.FromEmail == "" {
		return nil, errors.New("SendGrid sender email must be provided")
	}

	client := sendgrid.NewSendClient(cfg.SendGrid.APIKey)

	return &SendGridClient{
		client:    client,
		fromEmail: cfg.SendGrid.FromEmail,
		fromName:  cfg.SendGrid.FromName,
	}, nil
}

// SendEmail sends an email notification using SendGrid
func (s *SendGridClient) SendEmail(notification *models.Notification) error {
	from := mail.NewEmail(s.fromName, s.fromEmail)
	to := mail.NewEmail("", notification.Channel) // Channel contains the recipient's email address
	
	message := mail.NewSingleEmail(from, notification.Subject, to, notification.Content, notification.Content)
	
	response, err := s.client.Send(message)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	if response.StatusCode >= 400 {
		return fmt.Errorf("failed to send email, status code: %d, body: %s", response.StatusCode, response.Body)
	}
	
	return nil
} 