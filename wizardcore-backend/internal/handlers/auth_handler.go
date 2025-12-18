package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userService *services.UserService
	logger      *zap.Logger
}

func NewAuthHandler(userService *services.UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userService: userService,
		logger:      logger,
	}
}

func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Validate required fields
	if req.SupabaseUserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "supabase_user_id is required"})
		return
	}
	if req.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email is required"})
		return
	}

	// Convert to User model
	user := &models.User{
		SupabaseUserID: req.SupabaseUserID,
		Email:          req.Email,
		DisplayName:    &req.DisplayName,
		TotalXP:        0,
		PracticeScore:  0,
		CurrentStreak:  0,
		LongestStreak:  0,
	}

	// Create user
	if err := h.userService.CreateUser(user); err != nil {
		h.logger.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}