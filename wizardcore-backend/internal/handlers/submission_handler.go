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

type SubmissionHandler struct {
	submissionService *services.SubmissionService
	logger            *zap.Logger
}

func NewSubmissionHandler(submissionService *services.SubmissionService, logger *zap.Logger) *SubmissionHandler {
	return &SubmissionHandler{
		submissionService: submissionService,
		logger:            logger,
	}
}

func (h *SubmissionHandler) CreateSubmission(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create submission model
	submission := &models.Submission{
		ID:         uuid.New(),
		UserID:     userID,
		ExerciseID: req.ExerciseID,
		SourceCode: req.SourceCode,
		LanguageID: req.LanguageID,
		Status:     "pending",
	}

	var err error
	if req.MatchID != nil {
		err = h.submissionService.CreateSubmissionWithMatch(submission, req.MatchID)
	} else {
		err = h.submissionService.CreateSubmission(submission)
	}

	if err != nil {
		h.logger.Error("Failed to create submission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process submission"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"submission": submission})
}

func (h *SubmissionHandler) GetSubmission(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	submissionIDStr := c.Param("id")
	submissionID, err := uuid.Parse(submissionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := h.submissionService.GetSubmissionByID(submissionID)
	if err != nil {
		h.logger.Error("Failed to fetch submission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch submission"})
		return
	}
	if submission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
		return
	}

	// Ensure the submission belongs to the authenticated user
	if submission.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submission": submission})
}

func (h *SubmissionHandler) GetLatestSubmission(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	submission, err := h.submissionService.GetLatestSubmission(exerciseID, userID)
	if err != nil {
		h.logger.Error("Failed to fetch latest submission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch latest submission"})
		return
	}
	if submission == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No submission found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"submission": submission})
}

func (h *SubmissionHandler) SaveDraft(c *gin.Context) {
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	exerciseIDStr := c.Param("exercise_id")
	exerciseID, err := uuid.Parse(exerciseIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	var req struct {
		SourceCode string `json:"source_code" binding:"required"`
		LanguageID int    `json:"language_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	submission := &models.Submission{
		ID:             uuid.New(),
		UserID:         userID,
		ExerciseID:     exerciseID,
		SourceCode:     req.SourceCode,
		LanguageID:     req.LanguageID,
		SubmissionType: "draft",
		Status:         "draft",
	}

	if err := h.submissionService.SaveDraft(submission); err != nil {
		h.logger.Error("Failed to save draft", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save draft"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Draft saved", "submission": submission})
}