package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	if notification.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if notification.Type == "" {
		return fmt.Errorf("type is required")
	}
	if notification.Title == "" {
		return fmt.Errorf("title is required")
	}
	return s.notificationRepo.Create(notification)
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID, limit, offset int) (*models.NotificationResponse, error) {
	notifications, err := s.notificationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch notifications: %w", err)
	}
	unreadCount, err := s.notificationRepo.CountUnreadByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return &models.NotificationResponse{
		Notifications: notifications,
		UnreadCount:   unreadCount,
	}, nil
}

func (s *NotificationService) MarkAsRead(userID uuid.UUID, notificationIDs []uuid.UUID) error {
	if len(notificationIDs) == 0 {
		return nil
	}
	return s.notificationRepo.MarkAsRead(userID, notificationIDs)
}

func (s *NotificationService) DeleteNotification(userID, notificationID uuid.UUID) error {
	return s.notificationRepo.DeleteByID(userID, notificationID)
}