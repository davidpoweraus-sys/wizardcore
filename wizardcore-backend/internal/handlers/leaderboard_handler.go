package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type LeaderboardHandler struct {
	leaderboardService *services.LeaderboardService
	logger             *zap.Logger
}

func NewLeaderboardHandler(leaderboardService *services.LeaderboardService, logger *zap.Logger) *LeaderboardHandler {
	return &LeaderboardHandler{
		leaderboardService: leaderboardService,
		logger:             logger,
	}
}

func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	timeframe := c.DefaultQuery("timeframe", "all")
	pathwayID := c.Query("pathway")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))

	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 10
	}

	leaderboard, err := h.leaderboardService.GetLeaderboard(timeframe, &pathwayID, page, perPage)
	if err != nil {
		h.logger.Error("Failed to get leaderboard", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}