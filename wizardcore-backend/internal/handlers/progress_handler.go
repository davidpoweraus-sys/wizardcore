package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type ProgressHandler struct {
	progressService *services.ProgressService
	logger          *zap.Logger
}

func NewProgressHandler(progressService *services.ProgressService, logger *zap.Logger) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
		logger:          logger,
	}
}

func (h *ProgressHandler) GetUserProgress(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	progress, err := h.progressService.GetUserProgress(userID)
	if err != nil {
		h.logger.Error("Failed to get user progress", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func (h *ProgressHandler) GetMilestones(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	milestones, err := h.progressService.GetMilestones(userID)
	if err != nil {
		h.logger.Error("Failed to get milestones", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch milestones"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"milestones": milestones})
}

func (h *ProgressHandler) GetWeeklyActivity(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	weeklyActivity, err := h.progressService.GetWeeklyActivity(userID)
	if err != nil {
		h.logger.Error("Failed to get weekly activity", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly activity"})
		return
	}

	c.JSON(http.StatusOK, weeklyActivity)
}

func (h *ProgressHandler) GetWeeklyHours(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	weeklyHours, err := h.progressService.GetWeeklyHours(userID)
	if err != nil {
		h.logger.Error("Failed to get weekly hours", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch weekly hours"})
		return
	}

	c.JSON(http.StatusOK, weeklyHours)
}