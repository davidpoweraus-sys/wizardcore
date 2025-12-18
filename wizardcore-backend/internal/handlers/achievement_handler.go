package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type AchievementHandler struct {
	achievementService *services.AchievementService
	logger             *zap.Logger
}

func NewAchievementHandler(achievementService *services.AchievementService, logger *zap.Logger) *AchievementHandler {
	return &AchievementHandler{
		achievementService: achievementService,
		logger:             logger,
	}
}

func (h *AchievementHandler) GetUserAchievements(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	achievements, err := h.achievementService.GetUserAchievements(userID)
	if err != nil {
		h.logger.Error("Failed to get user achievements", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"achievements": achievements})
}