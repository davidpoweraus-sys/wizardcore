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
	userService     *services.UserService
	progressService *services.ProgressService
	activityService *services.ActivityService
	logger          *zap.Logger
}

func NewUserHandler(userService *services.UserService, progressService *services.ProgressService, activityService *services.ActivityService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService:     userService,
		progressService: progressService,
		activityService: activityService,
		logger:          logger,
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

	preferences, err := h.userService.GetUserPreferences(c.Request.Context(), user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"preferences": preferences})
}

func (h *UserHandler) UpdatePreferences(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.UpdatePreferencesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	// Convert request to updates map
	updates := make(map[string]interface{})
	if req.Theme != nil {
		updates["theme"] = *req.Theme
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.EmailNotifications != nil {
		updates["email_notifications"] = *req.EmailNotifications
	}
	if req.PushNotifications != nil {
		updates["push_notifications"] = *req.PushNotifications
	}
	if req.PublicProfile != nil {
		updates["public_profile"] = *req.PublicProfile
	}
	if req.ShowProgress != nil {
		updates["show_progress"] = *req.ShowProgress
	}
	if req.AutoSave != nil {
		updates["auto_save"] = *req.AutoSave
	}
	if req.SoundEffects != nil {
		updates["sound_effects"] = *req.SoundEffects
	}

	err = h.userService.UpdateUserPreferences(c.Request.Context(), user.ID, updates)
	if err != nil {
		h.logger.Error("Failed to update preferences", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update preferences"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Preferences updated successfully"})
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
	studyTime := totals.TotalStudyTimeMinutes / 60 // Convert minutes to hours
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

	// Get query parameters for pagination
	limit := 10 // default limit
	offset := 0 // default offset

	// TODO: Parse limit and offset from query parameters if needed

	activities, err := h.activityService.GetUserActivities(c.Request.Context(), user.ID, limit, offset)
	if err != nil {
		h.logger.Error("Failed to fetch activities", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch activities"})
		return
	}

	// Convert activities to frontend format
	var frontendActivities []map[string]interface{}
	for _, activity := range activities {
		frontendActivity := map[string]interface{}{
			"id":          activity.ID.String(),
			"type":        mapActivityType(activity.ActivityType),
			"title":       activity.Title,
			"description": getStringValue(activity.Description),
			"time":        formatActivityTime(activity.CreatedAt),
			"icon":        getStringValue(activity.Icon),
			"color":       getStringValue(activity.Color),
		}
		frontendActivities = append(frontendActivities, frontendActivity)
	}

	c.JSON(http.StatusOK, gin.H{"activities": frontendActivities})
}

// Helper function to map activity types to frontend types
func mapActivityType(activityType string) string {
	switch activityType {
	case "practice":
		return "practice"
	case "completion":
		return "completion"
	case "achievement":
		return "achievement"
	case "streak":
		return "streak"
	default:
		return "other"
	}
}

// Helper function to format time for frontend
func formatActivityTime(t time.Time) string {
	// Return relative time (e.g., "2 hours ago", "Yesterday")
	// For now, return ISO format
	return t.Format(time.RFC3339)
}

// Helper function to get string value from pointer
func getStringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
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
		"user":        user,
		"progress":    progress,
		"exported_at": time.Now().Format(time.RFC3339),
		"note":        "This is a partial export. Full export includes achievements, submissions, certificates, etc.",
	}

	c.JSON(http.StatusOK, exportData)
}
