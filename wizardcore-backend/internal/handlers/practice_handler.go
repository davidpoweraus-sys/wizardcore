package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type PracticeHandler struct {
	practiceService *services.PracticeService
	logger          *zap.Logger
}

func NewPracticeHandler(practiceService *services.PracticeService, logger *zap.Logger) *PracticeHandler {
	return &PracticeHandler{
		practiceService: practiceService,
		logger:          logger,
	}
}

func (h *PracticeHandler) GetChallenges(c *gin.Context) {
	challenges, err := h.practiceService.GetChallenges()
	if err != nil {
		h.logger.Error("Failed to get challenges", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch challenges"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"challenges": challenges})
}

func (h *PracticeHandler) GetAreas(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	areas, err := h.practiceService.GetAreas(userID)
	if err != nil {
		h.logger.Error("Failed to get practice areas", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch practice areas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"areas": areas})
}

func (h *PracticeHandler) GetStats(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	stats, err := h.practiceService.GetStats(userID)
	if err != nil {
		h.logger.Error("Failed to get practice stats", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch practice stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *PracticeHandler) GetRecentMatches(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if limit > 20 {
		limit = 20
	}
	matches, err := h.practiceService.GetRecentMatches(userID, limit)
	if err != nil {
		h.logger.Error("Failed to get recent matches", zap.Error(err), zap.String("user_id", userID.String()))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent matches"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"matches": matches})
}

func (h *PracticeHandler) StartChallenge(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	challengeType := c.Param("type")
	match, err := h.practiceService.StartChallenge(userID, challengeType)
	if err != nil {
		h.logger.Error("Failed to start challenge", zap.Error(err), zap.String("user_id", userID.String()), zap.String("type", challengeType))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start challenge"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"match": match})
}