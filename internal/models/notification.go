package models

import (
	"time"
)

// NotificationType represents different notification channels
type NotificationType string

const (
	// NotificationTypeEmail represents email notifications
	NotificationTypeEmail NotificationType = "email"
	// NotificationTypeTelegram represents telegram notifications
	NotificationTypeTelegram NotificationType = "telegram"
)

// NotificationStatus represents the current status of a notification
type NotificationStatus string

const (
	// NotificationStatusPending means the notification is waiting to be processed
	NotificationStatusPending NotificationStatus = "pending"
	// NotificationStatusSent means the notification was successfully sent
	NotificationStatusSent NotificationStatus = "sent"
	// NotificationStatusFailed means the notification sending failed
	NotificationStatusFailed NotificationStatus = "failed"
)

// Notification represents a notification that needs to be sent
type Notification struct {
	ID        string            `json:"id"`
	UserID    string            `json:"user_id"`
	Type      NotificationType  `json:"type"`
	Channel   string            `json:"channel"` // email address or telegram chat ID
	Subject   string            `json:"subject"`
	Content   string            `json:"content"`
	Status    NotificationStatus `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	SentAt    *time.Time        `json:"sent_at,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// KafkaNotificationMessage represents a message received from Kafka
type KafkaNotificationMessage struct {
	UserID    string           `json:"user_id"`
	Type      NotificationType `json:"type"`
	Channel   string           `json:"channel"`
	Subject   string           `json:"subject"`
	Content   string           `json:"content"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
} 