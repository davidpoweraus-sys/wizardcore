package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID  `json:"id" db:"id"`
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	Type       string     `json:"type" db:"type"`
	Title      string     `json:"title" db:"title"`
	Message    *string    `json:"message,omitempty" db:"message"`
	Icon       *string    `json:"icon,omitempty" db:"icon"`
	ActionURL  *string    `json:"action_url,omitempty" db:"action_url"`
	IsRead     bool       `json:"is_read" db:"is_read"`
	ReadAt     *time.Time `json:"read_at,omitempty" db:"read_at"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
}

type NotificationResponse struct {
	Notifications []Notification `json:"notifications"`
	UnreadCount   int            `json:"unread_count"`
}

type MarkAsReadRequest struct {
	NotificationIDs []uuid.UUID `json:"notification_ids"`
}