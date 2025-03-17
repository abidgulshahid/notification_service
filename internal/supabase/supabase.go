package supabase

import (
	"errors"
	"fmt"
	"time"

	supabase "github.com/nedpals/supabase-go"
	"github.com/notification_service/internal/config"
	"github.com/notification_service/internal/models"
)

// Client represents a Supabase client
type Client struct {
	client    *supabase.Client
	tableName string
}

// NewClient creates a new Supabase client
func NewClient(cfg *config.Config) (*Client, error) {
	if cfg.Supabase.URL == "" || cfg.Supabase.APIKey == "" {
		return nil, errors.New("Supabase URL and API key must be provided")
	}

	client := supabase.CreateClient(cfg.Supabase.URL, cfg.Supabase.APIKey)

	return &Client{
		client:    client,
		tableName: cfg.Supabase.NotificationsTable,
	}, nil
}

// InsertNotification inserts a notification into the Supabase database
func (c *Client) InsertNotification(notification *models.Notification) (string, error) {
	// Set default values
	if notification.Status == "" {
		notification.Status = models.NotificationStatusPending
	}
	
	now := time.Now()
	notification.CreatedAt = now
	notification.UpdatedAt = now

	// Insert the notification
	var result struct {
		ID string `json:"id"`
	}
	
	err := c.client.DB.From(c.tableName).Insert(notification, &result, false)
	if err != nil {
		return "", fmt.Errorf("failed to insert notification: %w", err)
	}

	return result.ID, nil
}

// UpdateNotificationStatus updates the status of a notification
func (c *Client) UpdateNotificationStatus(id string, status models.NotificationStatus) error {
	updateData := map[string]interface{}{
		"status":     status,
		"updated_at": time.Now(),
	}

	// If the notification was sent, update the sent_at field
	if status == models.NotificationStatusSent {
		sentAt := time.Now()
		updateData["sent_at"] = sentAt
	}

	err := c.client.DB.From(c.tableName).Update(updateData, supabase.FilterOptions{
		Filters: []string{fmt.Sprintf("id.eq.%s", id)},
	}, nil)

	if err != nil {
		return fmt.Errorf("failed to update notification status: %w", err)
	}

	return nil
}

// GetNotification retrieves a notification by ID
func (c *Client) GetNotification(id string) (*models.Notification, error) {
	var notifications []models.Notification
	
	err := c.client.DB.From(c.tableName).Select("*", "", false).
		Filter("id", "eq", id).
		Execute(&notifications)
	
	if err != nil {
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}

	if len(notifications) == 0 {
		return nil, fmt.Errorf("notification not found: %s", id)
	}

	return &notifications[0], nil
} 