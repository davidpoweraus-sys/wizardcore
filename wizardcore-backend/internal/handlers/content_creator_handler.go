package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type ContentCreatorHandler struct {
	contentCreatorService *services.ContentCreatorService
	userService           *services.UserService
	logger                *zap.Logger
}

func NewContentCreatorHandler(contentCreatorService *services.ContentCreatorService, userService *services.UserService, logger *zap.Logger) *ContentCreatorHandler {
	return &ContentCreatorHandler{
		contentCreatorService: contentCreatorService,
		userService:           userService,
		logger:                logger,
	}
}

// GetContentCreatorProfile gets the content creator profile for the authenticated user
func (h *ContentCreatorHandler) GetContentCreatorProfile(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	profile, err := h.contentCreatorService.GetContentCreatorProfile(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch content creator profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// UpdateContentCreatorProfile updates the content creator profile
func (h *ContentCreatorHandler) UpdateContentCreatorProfile(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var updateData struct {
		Bio            string   `json:"bio"`
		Specialization []string `json:"specialization"`
		Website        string   `json:"website"`
		GithubURL      string   `json:"github_url"`
		LinkedinURL    string   `json:"linkedin_url"`
		TwitterURL     string   `json:"twitter_url"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	profile, err := h.contentCreatorService.UpdateContentCreatorProfile(user.ID, updateData.Bio, updateData.Specialization, updateData.Website, updateData.GithubURL, updateData.LinkedinURL, updateData.TwitterURL)
	if err != nil {
		h.logger.Error("Failed to update content creator profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"profile": profile})
}

// GetContentCreatorStats gets statistics for a content creator
func (h *ContentCreatorHandler) GetContentCreatorStats(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	stats, err := h.contentCreatorService.GetContentCreatorStats(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch content creator stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// CreatePathway creates a new pathway (course)
func (h *ContentCreatorHandler) CreatePathway(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var pathwayData struct {
		Title         string   `json:"title" binding:"required"`
		Subtitle      string   `json:"subtitle"`
		Description   string   `json:"description" binding:"required"`
		Level         string   `json:"level" binding:"required"`
		DurationWeeks int      `json:"duration_weeks" binding:"required"`
		ColorGradient string   `json:"color_gradient"`
		Icon          string   `json:"icon"`
		Prerequisites []string `json:"prerequisites"`
		Status        string   `json:"status"`
	}

	if err := c.ShouldBindJSON(&pathwayData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	pathway, err := h.contentCreatorService.CreatePathway(
		user.ID,
		pathwayData.Title,
		pathwayData.Subtitle,
		pathwayData.Description,
		pathwayData.Level,
		pathwayData.DurationWeeks,
		pathwayData.ColorGradient,
		pathwayData.Icon,
		pathwayData.Prerequisites,
		pathwayData.Status,
	)
	if err != nil {
		h.logger.Error("Failed to create pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create pathway: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"pathway": pathway})
}

// UpdatePathway updates an existing pathway
func (h *ContentCreatorHandler) UpdatePathway(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	pathwayIDStr := c.Param("id")
	pathwayID, err := uuid.Parse(pathwayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pathway ID"})
		return
	}

	var updateData struct {
		Title         string   `json:"title"`
		Subtitle      string   `json:"subtitle"`
		Description   string   `json:"description"`
		Level         string   `json:"level"`
		DurationWeeks int      `json:"duration_weeks"`
		ColorGradient string   `json:"color_gradient"`
		Icon          string   `json:"icon"`
		Prerequisites []string `json:"prerequisites"`
		Status        string   `json:"status"`
		ReviewNotes   string   `json:"review_notes"`
		ChangeNotes   string   `json:"change_notes"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	pathway, err := h.contentCreatorService.UpdatePathway(
		user.ID,
		pathwayID,
		updateData.Title,
		updateData.Subtitle,
		updateData.Description,
		updateData.Level,
		updateData.DurationWeeks,
		updateData.ColorGradient,
		updateData.Icon,
		updateData.Prerequisites,
		updateData.Status,
		updateData.ReviewNotes,
		updateData.ChangeNotes,
	)
	if err != nil {
		h.logger.Error("Failed to update pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update pathway: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pathway": pathway})
}

// CreateModule creates a new module within a pathway
func (h *ContentCreatorHandler) CreateModule(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	pathwayIDStr := c.Param("pathwayId")
	pathwayID, err := uuid.Parse(pathwayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pathway ID"})
		return
	}

	var moduleData struct {
		Title          string `json:"title" binding:"required"`
		Description    string `json:"description"`
		SortOrder      int    `json:"sort_order" binding:"required"`
		EstimatedHours int    `json:"estimated_hours"`
		XPReward       int    `json:"xp_reward"`
		Status         string `json:"status"`
	}

	if err := c.ShouldBindJSON(&moduleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	module, err := h.contentCreatorService.CreateModule(
		user.ID,
		pathwayID,
		moduleData.Title,
		moduleData.Description,
		moduleData.SortOrder,
		moduleData.EstimatedHours,
		moduleData.XPReward,
		moduleData.Status,
	)
	if err != nil {
		h.logger.Error("Failed to create module", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create module: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"module": module})
}

// CreateExercise creates a new exercise within a module
func (h *ContentCreatorHandler) CreateExercise(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	moduleIDStr := c.Param("moduleId")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var exerciseData struct {
		Title            string   `json:"title" binding:"required"`
		Difficulty       string   `json:"difficulty" binding:"required"`
		Points           int      `json:"points"`
		TimeLimitMinutes int      `json:"time_limit_minutes"`
		SortOrder        int      `json:"sort_order" binding:"required"`
		Objectives       []string `json:"objectives"`
		Content          string   `json:"content"`
		Examples         string   `json:"examples"`
		Description      string   `json:"description" binding:"required"`
		Constraints      []string `json:"constraints"`
		Hints            []string `json:"hints"`
		StarterCode      string   `json:"starter_code"`
		SolutionCode     string   `json:"solution_code"`
		LanguageID       int      `json:"language_id" binding:"required"`
		Tags             []string `json:"tags"`
		Status           string   `json:"status"`
		RequiresApproval bool     `json:"requires_approval"`
		TestCases        []struct {
			Input          string `json:"input"`
			ExpectedOutput string `json:"expected_output"`
			IsHidden       bool   `json:"is_hidden"`
			Points         int    `json:"points"`
		} `json:"test_cases"`
	}

	if err := c.ShouldBindJSON(&exerciseData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	// Convert test cases to proper format
	testCases := make([]services.TestCaseInput, len(exerciseData.TestCases))
	for i, tc := range exerciseData.TestCases {
		testCases[i] = services.TestCaseInput{
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsHidden:       tc.IsHidden,
			Points:         tc.Points,
		}
	}

	exercise, err := h.contentCreatorService.CreateExercise(
		user.ID,
		moduleID,
		exerciseData.Title,
		exerciseData.Difficulty,
		exerciseData.Points,
		exerciseData.TimeLimitMinutes,
		exerciseData.SortOrder,
		exerciseData.Objectives,
		exerciseData.Content,
		exerciseData.Examples,
		exerciseData.Description,
		exerciseData.Constraints,
		exerciseData.Hints,
		exerciseData.StarterCode,
		exerciseData.SolutionCode,
		exerciseData.LanguageID,
		exerciseData.Tags,
		exerciseData.Status,
		exerciseData.RequiresApproval,
		testCases,
	)
	if err != nil {
		h.logger.Error("Failed to create exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"exercise": exercise})
}

// GetMyContent gets all content created by the authenticated user
func (h *ContentCreatorHandler) GetMyContent(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	contentType := c.Query("type") // "pathways", "modules", "exercises", or empty for all
	status := c.Query("status")    // "draft", "published", "archived", or empty for all

	content, err := h.contentCreatorService.GetMyContent(user.ID, contentType, status)
	if err != nil {
		h.logger.Error("Failed to fetch content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch content"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"content": content})
}

// SubmitForReview submits content for admin review
func (h *ContentCreatorHandler) SubmitForReview(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var reviewData struct {
		ContentType string `json:"content_type" binding:"required"`
		ContentID   string `json:"content_id" binding:"required"`
		Notes       string `json:"notes"`
	}

	if err := c.ShouldBindJSON(&reviewData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	contentID, err := uuid.Parse(reviewData.ContentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content ID"})
		return
	}

	review, err := h.contentCreatorService.SubmitForReview(user.ID, reviewData.ContentType, contentID, reviewData.Notes)
	if err != nil {
		h.logger.Error("Failed to submit for review", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit for review: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"review": review})
}

// GetContentReviews gets reviews for content created by the user
func (h *ContentCreatorHandler) GetContentReviews(c *gin.Context) {
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
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	status := c.Query("status") // "pending", "approved", "rejected", or empty for all

	reviews, err := h.contentCreatorService.GetContentReviews(user.ID, status)
	if err != nil {
		h.logger.Error("Failed to fetch reviews", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
