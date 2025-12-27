package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type PathwayHandler struct {
	pathwayService *services.PathwayService
	userService    *services.UserService
	logger         *zap.Logger
}

func NewPathwayHandler(pathwayService *services.PathwayService, userService *services.UserService, logger *zap.Logger) *PathwayHandler {
	return &PathwayHandler{
		pathwayService: pathwayService,
		userService:    userService,
		logger:         logger,
	}
}

func (h *PathwayHandler) GetAllPathways(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		// If user is not authenticated, return pathways without enrollment status
		pathways, err := h.pathwayService.GetAllPathways()
		if err != nil {
			h.logger.Error("Failed to fetch pathways", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pathways"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pathways": pathways})
		return
	}

	// Get internal user ID
	user, err := h.userService.GetUserBySupabaseUserID(supabaseUserID)
	if err != nil {
		h.logger.Error("Failed to fetch user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pathways"})
		return
	}
	if user == nil {
		// User not found in our database, return pathways without enrollment status
		pathways, err := h.pathwayService.GetAllPathways()
		if err != nil {
			h.logger.Error("Failed to fetch pathways", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pathways"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"pathways": pathways})
		return
	}

	// Get pathways with enrollment status
	pathways, err := h.pathwayService.GetAllPathwaysWithEnrollment(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch pathways with enrollment", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pathways"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pathways": pathways})
}

func (h *PathwayHandler) GetPathway(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pathway ID"})
		return
	}

	pathway, err := h.pathwayService.GetPathwayByID(id)
	if err != nil {
		h.logger.Error("Failed to fetch pathway", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch pathway"})
		return
	}
	if pathway == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pathway not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pathway": pathway})
}

func (h *PathwayHandler) EnrollPathway(c *gin.Context) {
	supabaseUserID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get internal user ID
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

	idStr := c.Param("id")
	pathwayID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid pathway ID"})
		return
	}

	err = h.pathwayService.EnrollUser(user.ID, pathwayID)
	if err != nil {
		h.logger.Error("Failed to enroll user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully enrolled in pathway"})
}

func (h *PathwayHandler) GetUserPathways(c *gin.Context) {
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

	enrollments, err := h.pathwayService.GetUserEnrollments(user.ID)
	if err != nil {
		h.logger.Error("Failed to fetch user pathways", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user pathways"})
		return
	}

	// For each enrollment, we might want to enrich with pathway details.
	// For simplicity, we'll just return enrollments.
	c.JSON(http.StatusOK, gin.H{"enrollments": enrollments})
}

func (h *PathwayHandler) GetDeadlines(c *gin.Context) {
	// TODO: implement deadlines logic
	c.JSON(http.StatusOK, gin.H{"deadlines": []string{}})
}
