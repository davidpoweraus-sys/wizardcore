package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type ExerciseHandler struct {
	exerciseService *services.ExerciseService
	logger          *zap.Logger
}

func NewExerciseHandler(exerciseService *services.ExerciseService, logger *zap.Logger) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
		logger:          logger,
	}
}

func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	exerciseWithTests, err := h.exerciseService.GetExerciseByID(id)
	if err != nil {
		h.logger.Error("Failed to fetch exercise", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercise"})
		return
	}
	if exerciseWithTests == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercise": exerciseWithTests})
}

func (h *ExerciseHandler) GetExercisesByModule(c *gin.Context) {
	moduleIDStr := c.Query("module_id")
	if moduleIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing module_id query parameter"})
		return
	}
	moduleID, err := uuid.Parse(moduleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	exercises, err := h.exerciseService.GetExercisesByModuleID(moduleID)
	if err != nil {
		h.logger.Error("Failed to fetch exercises by module", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
}

func (h *ExerciseHandler) GetExerciseStats(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid exercise ID"})
		return
	}

	stats, err := h.exerciseService.GetExerciseStats(id)
	if err != nil {
		h.logger.Error("Failed to fetch exercise stats", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercise stats"})
		return
	}
	if stats == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"stats": stats})
}