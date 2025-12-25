package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"go.uber.org/zap"
)

type RBACHandler struct {
	rbacService *services.RBACService
	logger      *zap.Logger
}

func NewRBACHandler(rbacService *services.RBACService, logger *zap.Logger) *RBACHandler {
	return &RBACHandler{
		rbacService: rbacService,
		logger:      logger,
	}
}

// ==================== ROLE MANAGEMENT ====================

// CreateRole creates a new role
// @Summary Create a new role
// @Description Create a new role with specified permissions
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateRoleRequest true "Role creation data"
// @Success 201 {object} models.Role
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles [post]
func (h *RBACHandler) CreateRole(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	role, err := h.rbacService.CreateRole(&req, userID)
	if err != nil {
		h.logger.Error("Failed to create role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, role)
}

// GetRole retrieves a role by ID
// @Summary Get role by ID
// @Description Get role details including permissions
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 200 {object} models.RoleWithPermissions
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/{id} [get]
func (h *RBACHandler) GetRole(c *gin.Context) {
	roleIDStr := c.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	role, err := h.rbacService.GetRoleWithPermissions(roleID)
	if err != nil {
		h.logger.Error("Failed to get role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if role == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// ListRoles retrieves all roles with pagination
// @Summary List all roles
// @Description Get paginated list of roles
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param include_system query bool false "Include system roles" default(false)
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles [get]
func (h *RBACHandler) ListRoles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	includeSystem, _ := strconv.ParseBool(c.DefaultQuery("include_system", "false"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	_ = includeSystem // Mark as used for now

	// This would use the repository directly
	// For now, return placeholder
	c.JSON(http.StatusOK, gin.H{
		"page":      page,
		"page_size": pageSize,
		"total":     0,
		"roles":     []models.Role{},
		"has_more":  false,
	})
}

// UpdateRole updates an existing role
// @Summary Update role
// @Description Update role details
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Param request body models.UpdateRoleRequest true "Role update data"
// @Success 200 {object} models.Role
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/{id} [put]
func (h *RBACHandler) UpdateRole(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var req models.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	role, err := h.rbacService.UpdateRole(roleID, &req, userID)
	if err != nil {
		h.logger.Error("Failed to update role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// DeleteRole deletes a role
// @Summary Delete role
// @Description Delete a role (cannot delete system roles)
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Param id path string true "Role ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/{id} [delete]
func (h *RBACHandler) DeleteRole(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roleIDStr := c.Param("id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	if err := h.rbacService.DeleteRole(roleID, userID); err != nil {
		h.logger.Error("Failed to delete role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// ==================== PERMISSION MANAGEMENT ====================

// CreatePermission creates a new permission
// @Summary Create a new permission
// @Description Create a new permission
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreatePermissionRequest true "Permission creation data"
// @Success 201 {object} models.Permission
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/permissions [post]
func (h *RBACHandler) CreatePermission(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CreatePermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	permission, err := h.rbacService.CreatePermission(&req, userID)
	if err != nil {
		h.logger.Error("Failed to create permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, permission)
}

// ==================== ROLE-PERMISSION MANAGEMENT ====================

// GrantPermission grants a permission to a role
// @Summary Grant permission to role
// @Description Grant a permission to a role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.GrantPermissionRequest true "Grant permission data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/permissions/grant [post]
func (h *RBACHandler) GrantPermission(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.GrantPermissionToRole(req.RoleID, req.PermissionID, userID); err != nil {
		h.logger.Error("Failed to grant permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission granted successfully"})
}

// RevokePermission revokes a permission from a role
// @Summary Revoke permission from role
// @Description Revoke a permission from a role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.GrantPermissionRequest true "Revoke permission data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/permissions/revoke [post]
func (h *RBACHandler) RevokePermission(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.GrantPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.RevokePermissionFromRole(req.RoleID, req.PermissionID, userID); err != nil {
		h.logger.Error("Failed to revoke permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Permission revoked successfully"})
}

// ==================== USER-ROLE MANAGEMENT ====================

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Description Assign a role to a user
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AssignRoleRequest true "Assign role data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/users/roles/assign [post]
func (h *RBACHandler) AssignRole(c *gin.Context) {
	// Get user ID from context (the user performing the assignment)
	assignedBy, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.AssignRoleToUser(req.UserID, req.RoleID, assignedBy, req.ExpiresAt, req.Reason); err != nil {
		h.logger.Error("Failed to assign role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role assigned successfully"})
}

// RemoveRole removes a role from a user
// @Summary Remove role from user
// @Description Remove a role from a user
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.AssignRoleRequest true "Remove role data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/users/roles/remove [post]
func (h *RBACHandler) RemoveRole(c *gin.Context) {
	// Get user ID from context (the user performing the removal)
	removedBy, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.RemoveRoleFromUser(req.UserID, req.RoleID, removedBy, req.Reason); err != nil {
		h.logger.Error("Failed to remove role", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role removed successfully"})
}

// GetUserRoles gets all roles for a user
// @Summary Get user roles
// @Description Get all roles assigned to a user
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} []models.Role
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/users/{user_id}/roles [get]
func (h *RBACHandler) GetUserRoles(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// This would use the repository directly
	// For now, return placeholder
	_ = userID // Mark as used for now
	c.JSON(http.StatusOK, []models.Role{})
}

// GetUserPermissions gets all permissions for a user
// @Summary Get user permissions
// @Description Get all permissions for a user (including via roles)
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "User ID"
// @Success 200 {object} models.UserPermissionsResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/users/{user_id}/permissions [get]
func (h *RBACHandler) GetUserPermissions(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	permissions, err := h.rbacService.GetUserPermissions(userID)
	if err != nil {
		h.logger.Error("Failed to get user permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, permissions)
}

// CheckPermission checks if a user has a specific permission
// @Summary Check user permission
// @Description Check if a user has a specific permission
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CheckPermissionRequest true "Permission check data"
// @Success 200 {object} models.PermissionCheckResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/permissions/check [post]
func (h *RBACHandler) CheckPermission(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.CheckPermissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	result, err := h.rbacService.CheckPermission(userID, req.PermissionName, req.ResourceID, c.Request.Context())
	if err != nil {
		h.logger.Error("Failed to check permission", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

// ==================== ROLE INHERITANCE ====================

// AddRoleInheritance adds inheritance between roles
// @Summary Add role inheritance
// @Description Make a role inherit from another role
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetRoleInheritanceRequest true "Role inheritance data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/inheritance/add [post]
func (h *RBACHandler) AddRoleInheritance(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SetRoleInheritanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.AddRoleInheritance(req.ChildRoleID, req.ParentRoleID, userID); err != nil {
		h.logger.Error("Failed to add role inheritance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role inheritance added successfully"})
}

// RemoveRoleInheritance removes inheritance between roles
// @Summary Remove role inheritance
// @Description Remove inheritance between roles
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.SetRoleInheritanceRequest true "Role inheritance data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/inheritance/remove [post]
func (h *RBACHandler) RemoveRoleInheritance(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req models.SetRoleInheritanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.RemoveRoleInheritance(req.ChildRoleID, req.ParentRoleID, userID); err != nil {
		h.logger.Error("Failed to remove role inheritance", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role inheritance removed successfully"})
}

// ==================== BULK OPERATIONS ====================

// BulkUpdateRolePermissions updates all permissions for a role
// @Summary Bulk update role permissions
// @Description Replace all permissions for a role with the provided list
// @Tags RBAC
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param role_id path string true "Role ID"
// @Param request body []string true "Array of permission IDs"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /rbac/roles/{role_id}/permissions/bulk [put]
func (h *RBACHandler) BulkUpdateRolePermissions(c *gin.Context) {
	// Get user ID from context
	userID, ok := middleware.GetSupabaseUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	roleIDStr := c.Param("role_id")
	roleID, err := uuid.Parse(roleIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role ID"})
		return
	}

	var permissionIDs []uuid.UUID
	if err := c.ShouldBindJSON(&permissionIDs); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.rbacService.BulkUpdateRolePermissions(roleID, permissionIDs, userID); err != nil {
		h.logger.Error("Failed to bulk update role permissions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Role permissions updated successfully"})
}

// ==================== HEALTH CHECK ====================

// HealthCheck checks RBAC system health
// @Summary RBAC health check
// @Description Check RBAC system health
// @Tags RBAC
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /rbac/health [get]
func (h *RBACHandler) HealthCheck(c *gin.Context) {
	// Basic health check - would verify database connections, etc.
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"timestamp": "2025-12-25T00:00:00Z", // Placeholder
		"system":    "rbac",
	})
}
