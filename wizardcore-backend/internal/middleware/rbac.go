package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

// RBACMiddleware creates middleware that checks for specific permissions
func RBACMiddleware(permissionName string, resourceIDParam string) gin.HandlerFunc {
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Get resource ID if parameter is specified
		var resourceID *string
		if resourceIDParam != "" {
			id := c.Param(resourceIDParam)
			if id != "" {
				resourceID = &id
			}
		}

		// Create context with HTTP request for audit logging
		ctx := context.WithValue(c.Request.Context(), "http_request", c.Request)

		// Check permission
		result, err := rbacService.CheckPermission(supabaseUserID, permissionName, resourceID, ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
			c.Abort()
			return
		}

		if !result.HasPermission {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"required_permission": permissionName,
				"user_roles":          result.Roles,
			})
			c.Abort()
			return
		}

		// Store permission result in context for downstream handlers
		c.Set("permission_check", result)

		c.Next()
	}
}

// RequireAnyPermission middleware checks if user has any of the specified permissions
func RequireAnyPermission(permissionNames []string, resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		supabaseUserID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get RBAC service
		rbacServiceInterface, exists := c.Get("rbac_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Get resource ID if parameter is specified
		var resourceID *string
		if resourceIDParam != "" {
			id := c.Param(resourceIDParam)
			if id != "" {
				resourceID = &id
			}
		}

		// Create context with HTTP request
		ctx := context.WithValue(c.Request.Context(), "http_request", c.Request)

		// Check each permission
		var lastError error
		var lastResult *models.PermissionCheckResult

		for _, permName := range permissionNames {
			result, err := rbacService.CheckPermission(supabaseUserID, permName, resourceID, ctx)
			if err != nil {
				lastError = err
				continue
			}

			if result.HasPermission {
				// Store permission result in context
				c.Set("permission_check", result)
				c.Set("granted_permission", permName)
				c.Next()
				return
			}

			lastResult = result
		}

		// If we get here, no permission was granted
		if lastError != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"error":                "Insufficient permissions",
				"required_permissions": permissionNames,
				"user_roles":           lastResult.Roles,
			})
		}
		c.Abort()
	}
}

// RequireAllPermissions middleware checks if user has all specified permissions
func RequireAllPermissions(permissionNames []string, resourceIDParam string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		supabaseUserID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get RBAC service
		rbacServiceInterface, exists := c.Get("rbac_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Get resource ID if parameter is specified
		var resourceID *string
		if resourceIDParam != "" {
			id := c.Param(resourceIDParam)
			if id != "" {
				resourceID = &id
			}
		}

		// Create context with HTTP request
		ctx := context.WithValue(c.Request.Context(), "http_request", c.Request)

		// Check all permissions
		var grantedPermissions []string
		var deniedPermissions []string

		for _, permName := range permissionNames {
			result, err := rbacService.CheckPermission(supabaseUserID, permName, resourceID, ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check permissions"})
				c.Abort()
				return
			}

			if result.HasPermission {
				grantedPermissions = append(grantedPermissions, permName)
			} else {
				deniedPermissions = append(deniedPermissions, permName)
			}
		}

		if len(deniedPermissions) > 0 {
			c.JSON(http.StatusForbidden, gin.H{
				"error":               "Insufficient permissions",
				"granted_permissions": grantedPermissions,
				"denied_permissions":  deniedPermissions,
				"required_all":        permissionNames,
			})
			c.Abort()
			return
		}

		// Store granted permissions in context
		c.Set("granted_permissions", grantedPermissions)
		c.Next()
	}
}

// RequireRole middleware checks if user has a specific role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		supabaseUserID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get RBAC service
		rbacServiceInterface, exists := c.Get("rbac_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Get user roles
		roles, err := rbacService.GetRepository().GetUserRoles(supabaseUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles"})
			c.Abort()
			return
		}

		// Check if user has the required role
		hasRole := false
		for _, role := range roles {
			if strings.EqualFold(role.Name, roleName) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			var roleNames []string
			for _, role := range roles {
				roleNames = append(roleNames, role.Name)
			}

			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Insufficient role",
				"required_role": roleName,
				"user_roles":    roleNames,
			})
			c.Abort()
			return
		}

		c.Set("required_role", roleName)
		c.Next()
	}
}

// RequireAnyRole middleware checks if user has any of the specified roles
func RequireAnyRole(roleNames []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		supabaseUserID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get RBAC service
		rbacServiceInterface, exists := c.Get("rbac_service")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "RBAC service not available"})
			c.Abort()
			return
		}

		rbacService, ok := rbacServiceInterface.(*services.RBACService)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid RBAC service"})
			c.Abort()
			return
		}

		// Get user roles
		roles, err := rbacService.GetRepository().GetUserRoles(supabaseUserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user roles"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		var userRoleNames []string
		var grantedRole string

		for _, role := range roles {
			userRoleNames = append(userRoleNames, role.Name)
			for _, requiredRole := range roleNames {
				if strings.EqualFold(role.Name, requiredRole) {
					grantedRole = role.Name
					break
				}
			}
			if grantedRole != "" {
				break
			}
		}

		if grantedRole == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "Insufficient role",
				"required_roles": roleNames,
				"user_roles":     userRoleNames,
			})
			c.Abort()
			return
		}

		c.Set("granted_role", grantedRole)
		c.Next()
	}
}

// ResourceOwnerMiddleware checks if user owns the resource they're trying to access
func ResourceOwnerMiddleware(resourceType, idParam string, getOwnerIDFunc func(c *gin.Context, resourceID string) (uuid.UUID, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID from context
		userID, ok := GetSupabaseUserID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			c.Abort()
			return
		}

		// Get resource ID from URL parameter
		resourceID := c.Param(idParam)
		if resourceID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Resource ID required"})
			c.Abort()
			return
		}

		// Get resource owner ID
		ownerID, err := getOwnerIDFunc(c, resourceID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get resource owner"})
			c.Abort()
			return
		}

		// Check if user is the owner
		if userID != ownerID {
			// User is not the owner, check if they have admin permission
			rbacServiceInterface, exists := c.Get("rbac_service")
			if exists {
				if rbacService, ok := rbacServiceInterface.(*services.RBACService); ok {
					ctx := context.WithValue(c.Request.Context(), "http_request", c.Request)
					permName := models.BuildPermissionName(resourceType, models.ActionManage)

					result, err := rbacService.CheckPermission(userID, permName, &resourceID, ctx)
					if err == nil && result.HasPermission {
						// User has admin permission, allow access
						c.Set("admin_override", true)
						c.Next()
						return
					}
				}
			}

			// Neither owner nor admin
			c.JSON(http.StatusForbidden, gin.H{
				"error":         "Access denied",
				"message":       "You do not own this resource and lack administrative permissions",
				"resource_type": resourceType,
				"resource_id":   resourceID,
			})
			c.Abort()
			return
		}

		// User is the owner
		c.Set("resource_owner", true)
		c.Next()
	}
}

// AuditLoggerMiddleware logs all requests with user and permission info
func AuditLoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user ID if available
		var userID string
		if id, ok := GetSupabaseUserID(c); ok {
			userID = id.String()
		} else {
			userID = "anonymous"
		}

		// Log request start
		logger.Info("Request started",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_id", userID),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		)

		// Process request
		c.Next()

		// Log request completion
		status := c.Writer.Status()
		logLevel := zap.InfoLevel
		if status >= 400 && status < 500 {
			logLevel = zap.WarnLevel
		} else if status >= 500 {
			logLevel = zap.ErrorLevel
		}

		logger.Log(logLevel, "Request completed",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("user_id", userID),
			zap.Int("status", status),
			zap.Int("size", c.Writer.Size()),
			zap.Duration("duration", c.GetDuration("request_duration")),
		)
	}
}

// GetPermissionCheckResult retrieves the permission check result from context
func GetPermissionCheckResult(c *gin.Context) (*models.PermissionCheckResult, bool) {
	value, exists := c.Get("permission_check")
	if !exists {
		return nil, false
	}

	result, ok := value.(*models.PermissionCheckResult)
	return result, ok
}

// Helper function to check permission in handlers (for conditional logic)
func CheckPermissionInHandler(c *gin.Context, permissionName string, resourceID *string) (bool, error) {
	// Get user ID from context
	userID, ok := GetSupabaseUserID(c)
	if !ok {
		return false, nil
	}

	// Get RBAC service
	rbacServiceInterface, exists := c.Get("rbac_service")
	if !exists {
		return false, nil
	}

	rbacService, ok := rbacServiceInterface.(*services.RBACService)
	if !ok {
		return false, nil
	}

	// Create context with HTTP request
	ctx := context.WithValue(c.Request.Context(), "http_request", c.Request)

	// Check permission
	result, err := rbacService.CheckPermission(userID, permissionName, resourceID, ctx)
	if err != nil {
		return false, err
	}

	return result.HasPermission, nil
}
