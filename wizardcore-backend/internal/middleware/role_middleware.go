package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/services"
)

// RoleMiddleware checks if the user has the required role
// DEPRECATED: Use RBACMiddleware instead for fine-grained permission control
func RoleMiddleware(db *sql.DB, requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context (set by AuthMiddleware)
		supabaseUserID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get RBAC service from context
		rbacServiceInterface, exists := c.Get("rbac_service")
		if !exists {
			// Fall back to old role-based system if RBAC service not available
			fallbackRoleMiddleware(db, requiredRole, c)
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Map required role to permission
		var permissionName string
		switch requiredRole {
		case "content_creator":
			permissionName = "content:manage"
		case "admin":
			permissionName = "system:admin"
		default:
			permissionName = requiredRole + ":access"
		}

		// Check permission using RBAC system
		result, err := rbacService.CheckPermission(supabaseUserID, permissionName, nil, c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !result.HasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// hasRequiredRole checks if the user's role satisfies the required role
// Admin has access to everything
// Content creators have access to creator endpoints
// Students only have access to student endpoints
func hasRequiredRole(userRole, requiredRole string) bool {
	// Admin has access to everything
	if userRole == "admin" {
		return true
	}

	// If content_creator is required, user must be content_creator or admin
	if requiredRole == "content_creator" {
		return userRole == "content_creator" || userRole == "admin"
	}

	// If student is required, anyone can access (student, content_creator, admin)
	if requiredRole == "student" {
		return true
	}

	// If admin is required, only admin can access
	if requiredRole == "admin" {
		return userRole == "admin"
	}

	return false
}

// fallbackRoleMiddleware is the old role-based system used as fallback
func fallbackRoleMiddleware(db *sql.DB, requiredRole string, c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		c.Abort()
		return
	}

	// Query user role from database
	var role string
	query := "SELECT role FROM users WHERE id = $1"
	err := db.QueryRow(query, userID.(uuid.UUID)).Scan(&role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user role"})
		}
		c.Abort()
		return
	}

	// Check if user has required role
	if !hasRequiredRole(role, requiredRole) {
		c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
		c.Abort()
		return
	}

	// Store role in context for handlers to use
	c.Set("role", role)
	c.Next()
}

// ContentCreatorMiddleware is a convenience middleware for content creator routes
func ContentCreatorMiddleware(db *sql.DB) gin.HandlerFunc {
	return RoleMiddleware(db, "content_creator")
}

// AdminMiddleware is a convenience middleware for admin routes
func AdminMiddleware(db *sql.DB) gin.HandlerFunc {
	return RoleMiddleware(db, "admin")
}
