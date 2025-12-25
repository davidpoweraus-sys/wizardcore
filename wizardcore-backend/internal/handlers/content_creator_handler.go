package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type ContentCreatorHandler struct {
	service *services.ContentCreatorService
	logger  *zap.Logger
}

func NewContentCreatorHandler(service *services.ContentCreatorService, logger *zap.Logger) *ContentCreatorHandler {
	return &ContentCreatorHandler{
		service: service,
		logger:  logger,
	}
}

// Profile Management

// CreateProfile godoc
// @Summary Create content creator profile
// @Description Create a new content creator profile for the authenticated user
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.CreateContentCreatorProfileRequest true "Profile data"
// @Success 201 {object} models.ContentCreatorProfile
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/profile [post]
func (h *ContentCreatorHandler) CreateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.CreateContentCreatorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.service.CreateProfile(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to create profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, profile)
}

// GetProfile godoc
// @Summary Get content creator profile
// @Description Get the content creator profile for the authenticated user
// @Tags content-creator
// @Produce json
// @Success 200 {object} models.ContentCreatorProfile
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/profile [get]
func (h *ContentCreatorHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	profile, err := h.service.GetProfile(userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get profile", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, profile)
}

// UpdateProfile godoc
// @Summary Update content creator profile
// @Description Update the content creator profile for the authenticated user
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.UpdateContentCreatorProfileRequest true "Updated profile data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/profile [put]
func (h *ContentCreatorHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.UpdateContentCreatorProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateProfile(userID.(uuid.UUID), &req); err != nil {
		h.logger.Error("Failed to update profile", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
}

// GetStats godoc
// @Summary Get creator statistics
// @Description Get statistics for the authenticated content creator
// @Tags content-creator
// @Produce json
// @Success 200 {object} models.CreatorStats
// @Failure 500 {object} map[string]string
// @Router /content-creator/stats [get]
func (h *ContentCreatorHandler) GetStats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	stats, err := h.service.GetStats(userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// Pathway Management

// CreatePathway godoc
// @Summary Create a new pathway
// @Description Create a new learning pathway
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.CreatePathwayRequest true "Pathway data"
// @Success 201 {object} models.Pathway
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways [post]
func (h *ContentCreatorHandler) CreatePathway(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.CreatePathwayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	pathway, err := h.service.CreatePathway(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to create pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pathway)
}

// GetPathways godoc
// @Summary Get creator's pathways
// @Description Get all pathways created by the authenticated user
// @Tags content-creator
// @Produce json
// @Param status query string false "Filter by status (draft, published, archived, under_review)"
// @Success 200 {array} models.Pathway
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways [get]
func (h *ContentCreatorHandler) GetPathways(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	status := c.Query("status")

	pathways, err := h.service.GetCreatorPathways(userID.(uuid.UUID), status)
	if err != nil {
		h.logger.Error("Failed to get pathways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, pathways)
}

// UpdatePathway godoc
// @Summary Update a pathway
// @Description Update an existing pathway
// @Tags content-creator
// @Accept json
// @Produce json
// @Param id path string true "Pathway ID"
// @Param request body models.UpdatePathwayRequest true "Updated pathway data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways/{id} [put]
func (h *ContentCreatorHandler) UpdatePathway(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	pathwayIDStr := c.Param("id")
	pathwayID, err := uuid.Parse(pathwayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pathway ID"})
		return
	}

	var req models.UpdatePathwayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdatePathway(pathwayID, userID.(uuid.UUID), &req); err != nil {
		h.logger.Error("Failed to update pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pathway updated successfully"})
}

// DeletePathway godoc
// @Summary Delete a pathway
// @Description Delete a pathway and all its modules and exercises
// @Tags content-creator
// @Param id path string true "Pathway ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways/{id} [delete]
func (h *ContentCreatorHandler) DeletePathway(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	pathwayIDStr := c.Param("id")
	pathwayID, err := uuid.Parse(pathwayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pathway ID"})
		return
	}

	if err := h.service.DeletePathway(pathwayID, userID.(uuid.UUID)); err != nil {
		h.logger.Error("Failed to delete pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "pathway deleted successfully"})
}

// Module Management

// CreateModule godoc
// @Summary Create a new module
// @Description Create a new module within a pathway
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.CreateModuleRequest true "Module data"
// @Success 201 {object} models.Module
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/modules [post]
func (h *ContentCreatorHandler) CreateModule(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.CreateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	module, err := h.service.CreateModule(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to create module", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, module)
}

// GetModules godoc
// @Summary Get creator's modules
// @Description Get all modules created by the authenticated user
// @Tags content-creator
// @Produce json
// @Param pathway_id query string false "Filter by pathway ID"
// @Success 200 {array} models.Module
// @Failure 500 {object} map[string]string
// @Router /content-creator/modules [get]
func (h *ContentCreatorHandler) GetModules(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var pathwayID *uuid.UUID
	pathwayIDStr := c.Query("pathway_id")
	if pathwayIDStr != "" {
		parsed, err := uuid.Parse(pathwayIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pathway ID"})
			return
		}
		pathwayID = &parsed
	}

	modules, err := h.service.GetCreatorModules(userID.(uuid.UUID), pathwayID)
	if err != nil {
		h.logger.Error("Failed to get modules", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, modules)
}

// UpdateModule godoc
// @Summary Update a module
// @Description Update an existing module
// @Tags content-creator
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param request body models.UpdateModuleRequest true "Updated module data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/modules/{id} [put]
func (h *ContentCreatorHandler) UpdateModule(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	moduleIDStr := c.Param("id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	var req models.UpdateModuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateModule(moduleID, userID.(uuid.UUID), &req); err != nil {
		h.logger.Error("Failed to update module", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "module updated successfully"})
}

// DeleteModule godoc
// @Summary Delete a module
// @Description Delete a module and all its exercises
// @Tags content-creator
// @Param id path string true "Module ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/modules/{id} [delete]
func (h *ContentCreatorHandler) DeleteModule(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	moduleIDStr := c.Param("id")
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
		return
	}

	if err := h.service.DeleteModule(moduleID, userID.(uuid.UUID)); err != nil {
		h.logger.Error("Failed to delete module", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "module deleted successfully"})
}

// Exercise Management

// CreateExercise godoc
// @Summary Create a new exercise
// @Description Create a new exercise with test cases
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.CreateExerciseRequest true "Exercise data"
// @Success 201 {object} models.Exercise
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/exercises [post]
func (h *ContentCreatorHandler) CreateExercise(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	exercise, err := h.service.CreateExercise(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to create exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, exercise)
}

// GetExercises godoc
// @Summary Get creator's exercises
// @Description Get all exercises created by the authenticated user
// @Tags content-creator
// @Produce json
// @Param module_id query string false "Filter by module ID"
// @Success 200 {array} models.Exercise
// @Failure 500 {object} map[string]string
// @Router /content-creator/exercises [get]
func (h *ContentCreatorHandler) GetExercises(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var moduleID *uuid.UUID
	moduleIDStr := c.Query("module_id")
	if moduleIDStr != "" {
		parsed, err := uuid.Parse(moduleIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid module ID"})
			return
		}
		moduleID = &parsed
	}

	exercises, err := h.service.GetCreatorExercises(userID.(uuid.UUID), moduleID)
	if err != nil {
		h.logger.Error("Failed to get exercises", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

// GetExercise godoc
// @Summary Get exercise with test cases
// @Description Get a specific exercise with all its test cases
// @Tags content-creator
// @Produce json
// @Param id path string true "Exercise ID"
// @Success 200 {object} models.ExerciseWithTests
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/exercises/{id} [get]
func (h *ContentCreatorHandler) GetExercise(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	exercise, err := h.service.GetExerciseWithTestCases(exerciseID, userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

// UpdateExercise godoc
// @Summary Update an exercise
// @Description Update an existing exercise
// @Tags content-creator
// @Accept json
// @Produce json
// @Param id path string true "Exercise ID"
// @Param request body models.UpdateExerciseRequest true "Updated exercise data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/exercises/{id} [put]
func (h *ContentCreatorHandler) UpdateExercise(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	var req models.UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.UpdateExercise(exerciseID, userID.(uuid.UUID), &req); err != nil {
		h.logger.Error("Failed to update exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exercise updated successfully"})
}

// DeleteExercise godoc
// @Summary Delete an exercise
// @Description Delete an exercise and all its test cases
// @Tags content-creator
// @Param id path string true "Exercise ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/exercises/{id} [delete]
func (h *ContentCreatorHandler) DeleteExercise(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid exercise ID"})
		return
	}

	if err := h.service.DeleteExercise(exerciseID, userID.(uuid.UUID)); err != nil {
		h.logger.Error("Failed to delete exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "exercise deleted successfully"})
}

// Review Management

// SubmitForReview godoc
// @Summary Submit content for review
// @Description Submit pathway, module, or exercise for admin review
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.SubmitContentForReviewRequest true "Review request"
// @Success 201 {object} models.ContentReview
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/reviews [post]
func (h *ContentCreatorHandler) SubmitForReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.SubmitContentForReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	review, err := h.service.SubmitForReview(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to submit for review", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, review)
}

// GetReviews godoc
// @Summary Get creator's reviews
// @Description Get all review submissions for the authenticated creator
// @Tags content-creator
// @Produce json
// @Success 200 {array} models.ContentReview
// @Failure 500 {object} map[string]string
// @Router /content-creator/reviews [get]
func (h *ContentCreatorHandler) GetReviews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	reviews, err := h.service.GetCreatorReviews(userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to get reviews", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

// ReviewContent godoc
// @Summary Review submitted content (Admin only)
// @Description Approve, reject, or request revision for submitted content
// @Tags admin
// @Accept json
// @Produce json
// @Param request body models.ReviewContentRequest true "Review decision"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/reviews [post]
func (h *ContentCreatorHandler) ReviewContent(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	// TODO: Check if user is admin
	// This should be done in middleware

	var req models.ReviewContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := h.service.ReviewContent(userID.(uuid.UUID), &req); err != nil {
		h.logger.Error("Failed to review content", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "content reviewed successfully"})
}

// ExportPathway godoc
// @Summary Export a pathway with all modules and exercises
// @Description Export a complete pathway including all modules, exercises, and test cases
// @Tags content-creator
// @Produce json
// @Param id path string true "Pathway ID"
// @Success 200 {object} models.ExportResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways/{id}/export [get]
func (h *ContentCreatorHandler) ExportPathway(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	pathwayIDStr := c.Param("id")
	pathwayID, err := uuid.Parse(pathwayIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pathway ID"})
		return
	}

	export, err := h.service.ExportPathway(pathwayID, userID.(uuid.UUID))
	if err != nil {
		h.logger.Error("Failed to export pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, export)
}

// ImportPathway godoc
// @Summary Import a pathway from JSON
// @Description Import a complete pathway including modules, exercises, and test cases
// @Tags content-creator
// @Accept json
// @Produce json
// @Param request body models.ImportPathwayRequest true "Pathway data to import"
// @Success 201 {object} models.Pathway
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /content-creator/pathways/import [post]
func (h *ContentCreatorHandler) ImportPathway(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	var req models.ImportPathwayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	pathway, err := h.service.ImportPathway(userID.(uuid.UUID), &req)
	if err != nil {
		h.logger.Error("Failed to import pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, pathway)
}
