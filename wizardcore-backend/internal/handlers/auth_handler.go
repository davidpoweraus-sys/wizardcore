package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"github.com/yourusername/wizardcore-backend/internal/version"
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

	// Log version and timing information
	versionInfo := version.GetInfo()
	h.logger.Info("CreateUser request",
		zap.String("supabase_user_id", req.SupabaseUserID.String()),
		zap.String("email", req.Email),
		zap.String("display_name", req.DisplayName),
		zap.String("path", c.Request.URL.Path),
		zap.String("method", c.Request.Method),
		zap.String("user_agent", c.Request.UserAgent()),
		zap.String("version", versionInfo.Version),
		zap.String("build_time", versionInfo.BuildTime),
		zap.String("client_ip", c.ClientIP()),
	)

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
		h.logger.Error("Failed to create user",
			zap.String("supabase_user_id", req.SupabaseUserID.String()),
			zap.String("email", req.Email),
			zap.Error(err),
			zap.String("version", versionInfo.Version),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	h.logger.Info("User created successfully",
		zap.String("supabase_user_id", req.SupabaseUserID.String()),
		zap.String("email", req.Email),
		zap.String("user_id", user.ID.String()),
		zap.String("version", versionInfo.Version),
		zap.Time("created_at", time.Now()),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user":    user,
	})
}