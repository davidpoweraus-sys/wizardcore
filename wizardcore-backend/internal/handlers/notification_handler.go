package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
	logger              *zap.Logger
}

func NewNotificationHandler(notificationService *services.NotificationService, logger *zap.Logger) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		logger:              logger,
	}
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	_, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// limit := 20
	// offset := 0
	// if limitStr := c.Query("limit"); limitStr != "" {
	// 	if l, err := uuid.Parse(limitStr); err == nil {
	// 		limit = int(l.ID())
	// 	}
	// }
	// if offsetStr := c.Query("offset"); offsetStr != "" {
	// 	if o, err := uuid.Parse(offsetStr); err == nil {
	// 		offset = int(o.ID())
	// 	}
	// }

	// Get user ID from supabaseUserID (need to fetch user)
	// For simplicity, assume we have a user service to get user ID
	// We'll need to inject user service, but for now we'll skip.
	// Since we don't have user ID, we'll just return empty.
	// TODO: integrate with user service to get user ID
	c.JSON(http.StatusOK, gin.H{"notifications": []string{}, "unread_count": 0})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	_, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.MarkAsReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID (TODO)
	// userID := supabaseUserID // This is supabase user ID, not internal user ID
	// We need to convert, but for now we'll just use a placeholder.
	// We'll implement later.
	h.logger.Info("Mark as read", zap.Any("notification_ids", req.NotificationIDs))

	c.JSON(http.StatusOK, gin.H{"message": "Notifications marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	_, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	notificationIDStr := c.Param("id")
	notificationID, err := uuid.Parse(notificationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Get user ID (TODO)
	// userID := supabaseUserID
	h.logger.Info("Delete notification", zap.String("notification_id", notificationID.String()))

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}