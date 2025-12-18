package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService    *services.UserService
	progressService *services.ProgressService
	logger         *zap.Logger
}

func NewUserHandler(userService *services.UserService, progressService *services.ProgressService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService:    userService,
		progressService: progressService,
		logger:         logger,
	}
}

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Fetch existing user
	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	var req models.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Update fields if provided
	if req.DisplayName != nil {
		user.DisplayName = req.DisplayName
	}
	if req.Bio != nil {
		user.Bio = req.Bio
	}
	if req.Location != nil {
		user.Location = req.Location
	}
	if req.Website != nil {
		user.Website = req.Website
	}
	if req.GithubUsername != nil {
		user.GithubUsername = req.GithubUsername
	}
	if req.TwitterUsername != nil {
		user.TwitterUsername = req.TwitterUsername
	}

	if err := h.userService.UpdateUser(user); err != nil {
		h.logger.Error("Failed to update user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    user,
	})
}

func (h *UserHandler) GetPreferences(c *gin.Context) {
	// TODO: implement preferences retrieval
	c.JSON(http.StatusOK, gin.H{"preferences": nil})
}

func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	// TODO: implement preferences update
	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated"})
}

func (h *UserHandler) GetStats(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	progress, err := h.progressService.GetUserProgress(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch progress", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	totals := progress.Totals
	activeCourses := len(progress.Pathways)
	completionRate := totals.OverallProgress
	studyTime := 0 // TODO: implement study time aggregation
	xpEarned := totals.TotalXP

	c.JSON(http.StatusOK, gin.H{
		"stats": gin.H{
			"active_courses":    activeCourses,
			"completion_rate":   completionRate,
			"study_time":        studyTime,
			"xp_earned":         xpEarned,
			"total_xp":          totals.TotalXP,
			"xp_this_week":      totals.XPThisWeek,
			"current_streak":    totals.CurrentStreak,
			"modules_completed": totals.ModulesCompleted,
			"modules_total":     totals.ModulesTotal,
		},
	})
}

func (h *UserHandler) GetActivities(c *gin.Context) {
	// TODO: implement activities retrieval
	c.JSON(http.StatusOK, gin.H{"activities": []string{}})
}

func (h *UserHandler) GetNavCounts(c *gin.Context) {
	// TODO: implement nav counts retrieval
	c.JSON(http.StatusOK, gin.H{"counts": nil})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	if err := h.userService.DeleteUser(user.ID); err != nil {
		h.logger.Error("Failed to delete user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Account deleted successfully"})
}

func (h *UserHandler) ExportData(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	progress, err := h.progressService.GetUserProgress(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch progress", zap.Error(err))
		// Continue without progress data
	}

	exportData := gin.H{
		"user": user,
		"progress": progress,
		"exported_at": time.Now().Format(time.RFC3339),
		"note": "This is a partial export. Full export includes achievements, submissions, certificates, etc.",
	}

	c.JSON(http.StatusOK, exportData)
}