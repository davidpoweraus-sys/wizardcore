package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) Create(notification *models.Notification) error {
	query := `
		INSERT INTO notifications (
			id, user_id, type, title, message, icon, action_url,
			is_read, read_at, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, created_at
	`
	if notification.ID == uuid.Nil {
		notification.ID = uuid.New()
	}
	now := time.Now()
	err := r.db.QueryRow(
		query,
		notification.ID,
		notification.UserID,
		notification.Type,
		notification.Title,
		notification.Message,
		notification.Icon,
		notification.ActionURL,
		notification.IsRead,
		notification.ReadAt,
		now,
	).Scan(&notification.ID, &notification.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

func (r *NotificationRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	query := `
		SELECT id, user_id, type, title, message, icon, action_url,
			is_read, read_at, created_at
		FROM notifications
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query notifications: %w", err)
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.Type,
			&n.Title,
			&n.Message,
			&n.Icon,
			&n.ActionURL,
			&n.IsRead,
			&n.ReadAt,
			&n.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, n)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return notifications, nil
}

func (r *NotificationRepository) CountUnreadByUserID(userID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false
	`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) MarkAsRead(userID uuid.UUID, notificationIDs []uuid.UUID) error {
	if len(notificationIDs) == 0 {
		return nil
	}
	// Build query with variable number of IDs
	query := `
		UPDATE notifications
		SET is_read = true, read_at = $1
		WHERE user_id = $2 AND id = ANY($3)
	`
	now := time.Now()
	_, err := r.db.Exec(query, now, userID, notificationIDs)
	if err != nil {
		return fmt.Errorf("failed to mark notifications as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) DeleteByID(userID, notificationID uuid.UUID) error {
	query := `DELETE FROM notifications WHERE user_id = $1 AND id = $2`
	result, err := r.db.Exec(query, userID, notificationID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("notification not found")
	}
	return nil
}