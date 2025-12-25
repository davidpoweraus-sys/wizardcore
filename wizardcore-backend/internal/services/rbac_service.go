package services

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
	"go.uber.org/zap"
)

type RBACService struct {
	rbacRepo *repositories.RBACRepository
	userRepo *repositories.UserRepository
	logger   *zap.Logger
}

// GetRepository exposes the RBAC repository (for middleware access)
func (s *RBACService) GetRepository() *repositories.RBACRepository {
	return s.rbacRepo
}

func NewRBACService(rbacRepo *repositories.RBACRepository, userRepo *repositories.UserRepository, logger *zap.Logger) *RBACService {
	return &RBACService{
		rbacRepo: rbacRepo,
		userRepo: userRepo,
		logger:   logger,
	}
}

// ==================== ROLE MANAGEMENT ====================

func (s *RBACService) CreateRole(req *models.CreateRoleRequest, createdBy uuid.UUID) (*models.Role, error) {
	// Validate role name format
	if !isValidRoleName(req.Name) {
		return nil, fmt.Errorf("invalid role name format. Use lowercase letters, numbers, and underscores only")
	}

	// Check if role already exists
	existing, err := s.rbacRepo.GetRoleByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing role: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("role '%s' already exists", req.Name)
	}

	// Create role
	role := &models.Role{
		ID:           uuid.New(),
		Name:         req.Name,
		Description:  req.Description,
		IsSystemRole: false,
		IsDefault:    req.IsDefault,
	}

	if err := s.rbacRepo.CreateRole(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	s.logger.Info("Role created",
		zap.String("role_name", role.Name),
		zap.String("created_by", createdBy.String()),
	)

	return role, nil
}

func (s *RBACService) UpdateRole(roleID uuid.UUID, req *models.UpdateRoleRequest, updatedBy uuid.UUID) (*models.Role, error) {
	// Get existing role
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}

	// Check if it's a system role
	if role.IsSystemRole {
		return nil, fmt.Errorf("cannot modify system role")
	}

	// Update fields if provided
	if req.Name != nil {
		if !isValidRoleName(*req.Name) {
			return nil, fmt.Errorf("invalid role name format")
		}

		// Check if new name already exists
		existing, err := s.rbacRepo.GetRoleByName(*req.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing role: %w", err)
		}
		if existing != nil && existing.ID != roleID {
			return nil, fmt.Errorf("role '%s' already exists", *req.Name)
		}
		role.Name = *req.Name
	}

	if req.Description != nil {
		role.Description = req.Description
	}

	if req.IsDefault != nil {
		role.IsDefault = *req.IsDefault
	}

	// Update in database
	if err := s.rbacRepo.UpdateRole(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	s.logger.Info("Role updated",
		zap.String("role_id", roleID.String()),
		zap.String("updated_by", updatedBy.String()),
	)

	return role, nil
}

func (s *RBACService) DeleteRole(roleID uuid.UUID, deletedBy uuid.UUID) error {
	// Get role to check if it's a system role
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if it's a default role
	if role.IsDefault {
		return fmt.Errorf("cannot delete default role")
	}

	// Delete role
	if err := s.rbacRepo.DeleteRole(roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	s.logger.Info("Role deleted",
		zap.String("role_id", roleID.String()),
		zap.String("deleted_by", deletedBy.String()),
	)

	return nil
}

func (s *RBACService) GetRoleWithPermissions(roleID uuid.UUID) (*models.RoleWithPermissions, error) {
	// Get role
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return nil, fmt.Errorf("role not found")
	}

	// Get role permissions
	permissions, err := s.rbacRepo.GetRolePermissions(roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}

	// Get user count for this role
	_, total, err := s.rbacRepo.GetUsersWithRole(roleID, 1, 1000000) // Large limit to get all
	if err != nil {
		return nil, fmt.Errorf("failed to count users: %w", err)
	}

	return &models.RoleWithPermissions{
		Role:        *role,
		Permissions: permissions,
		UserCount:   total,
	}, nil
}

// ==================== PERMISSION MANAGEMENT ====================

func (s *RBACService) CreatePermission(req *models.CreatePermissionRequest, createdBy uuid.UUID) (*models.Permission, error) {
	// Validate permission name format
	if !isValidPermissionName(req.Name) {
		return nil, fmt.Errorf("invalid permission name format. Use format 'resource:action'")
	}

	// Check if permission already exists
	existing, err := s.rbacRepo.GetPermissionByName(req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing permission: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("permission '%s' already exists", req.Name)
	}

	// Create permission
	permission := &models.Permission{
		ID:          uuid.New(),
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Description: req.Description,
		Resource:    req.Resource,
		Action:      req.Action,
		IsDangerous: req.IsDangerous,
	}

	if err := s.rbacRepo.CreatePermission(permission); err != nil {
		return nil, fmt.Errorf("failed to create permission: %w", err)
	}

	s.logger.Info("Permission created",
		zap.String("permission_name", permission.Name),
		zap.String("created_by", createdBy.String()),
	)

	return permission, nil
}

// ==================== ROLE-PERMISSION MANAGEMENT ====================

func (s *RBACService) GrantPermissionToRole(roleID, permissionID, grantedBy uuid.UUID) error {
	// Check if role exists
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if permission exists
	permission, err := s.rbacRepo.GetPermissionByID(permissionID)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found")
	}

	// Check if permission is dangerous and user has permission to grant dangerous permissions
	if permission.IsDangerous {
		// TODO: Add additional checks for dangerous permissions
		s.logger.Warn("Granting dangerous permission",
			zap.String("permission", permission.Name),
			zap.String("granted_by", grantedBy.String()),
		)
	}

	// Grant permission
	if err := s.rbacRepo.GrantPermissionToRole(roleID, permissionID, grantedBy); err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	s.logger.Info("Permission granted to role",
		zap.String("role_id", roleID.String()),
		zap.String("permission_id", permissionID.String()),
		zap.String("granted_by", grantedBy.String()),
	)

	return nil
}

func (s *RBACService) RevokePermissionFromRole(roleID, permissionID, revokedBy uuid.UUID) error {
	// Check if role exists
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if permission exists
	permission, err := s.rbacRepo.GetPermissionByID(permissionID)
	if err != nil {
		return fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return fmt.Errorf("permission not found")
	}

	// Revoke permission
	if err := s.rbacRepo.RevokePermissionFromRole(roleID, permissionID); err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	s.logger.Info("Permission revoked from role",
		zap.String("role_id", roleID.String()),
		zap.String("permission_id", permissionID.String()),
		zap.String("revoked_by", revokedBy.String()),
	)

	return nil
}

// ==================== USER-ROLE MANAGEMENT ====================

func (s *RBACService) AssignRoleToUser(userID, roleID, assignedBy uuid.UUID, expiresAt *time.Time, reason *string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if role exists
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if user already has this role
	userRoles, err := s.rbacRepo.GetUserRoles(userID)
	if err != nil {
		return fmt.Errorf("failed to get user roles: %w", err)
	}

	for _, userRole := range userRoles {
		if userRole.ID == roleID {
			return fmt.Errorf("user already has this role")
		}
	}

	// Assign role
	if err := s.rbacRepo.AssignRoleToUser(userID, roleID, assignedBy, expiresAt); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	// Log role assignment
	auditLog := &models.RoleAuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		TargetUserID: userID,
		RoleID:       roleID,
		Action:       models.RoleActionAssigned,
		Reason:       reason,
		PerformedBy:  assignedBy,
		CreatedAt:    time.Now(),
	}

	if err := s.rbacRepo.LogRoleChange(context.Background(), auditLog); err != nil {
		s.logger.Error("Failed to log role assignment", zap.Error(err))
	}

	s.logger.Info("Role assigned to user",
		zap.String("user_id", userID.String()),
		zap.String("role_id", roleID.String()),
		zap.String("assigned_by", assignedBy.String()),
	)

	return nil
}

func (s *RBACService) RemoveRoleFromUser(userID, roleID, removedBy uuid.UUID, reason *string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if role exists
	role, err := s.rbacRepo.GetRoleByID(roleID)
	if err != nil {
		return fmt.Errorf("failed to get role: %w", err)
	}
	if role == nil {
		return fmt.Errorf("role not found")
	}

	// Check if it's a system role and user is trying to remove it from themselves
	if role.IsSystemRole && userID == removedBy {
		return fmt.Errorf("cannot remove system role from yourself")
	}

	// Remove role
	if err := s.rbacRepo.RemoveRoleFromUser(userID, roleID); err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	// Log role removal
	auditLog := &models.RoleAuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		TargetUserID: userID,
		RoleID:       roleID,
		Action:       models.RoleActionRemoved,
		Reason:       reason,
		PerformedBy:  removedBy,
		CreatedAt:    time.Now(),
	}

	if err := s.rbacRepo.LogRoleChange(context.Background(), auditLog); err != nil {
		s.logger.Error("Failed to log role removal", zap.Error(err))
	}

	s.logger.Info("Role removed from user",
		zap.String("user_id", userID.String()),
		zap.String("role_id", roleID.String()),
		zap.String("removed_by", removedBy.String()),
	)

	return nil
}

// ==================== PERMISSION CHECKING ====================

func (s *RBACService) CheckPermission(userID uuid.UUID, permissionName string, resourceID *string, ctx context.Context) (*models.PermissionCheckResult, error) {
	// Get permission
	permission, err := s.rbacRepo.GetPermissionByName(permissionName)
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}
	if permission == nil {
		return nil, fmt.Errorf("permission not found: %s", permissionName)
	}

	// Check if user has permission
	hasPermission, err := s.rbacRepo.CheckUserPermission(userID, permissionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check permission: %w", err)
	}

	// Get user roles for context
	userRoles, err := s.rbacRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	var roleNames []string
	for _, role := range userRoles {
		roleNames = append(roleNames, role.Name)
	}

	// Log permission check
	auditLog := &models.PermissionAuditLog{
		ID:           uuid.New(),
		UserID:       userID,
		PermissionID: permission.ID,
		ResourceID:   resourceID,
		Action:       "check",
		Status:       map[bool]string{true: models.AuditStatusGranted, false: models.AuditStatusDenied}[hasPermission],
		Metadata: map[string]interface{}{
			"permission_name": permissionName,
			"resource_id":     resourceID,
			"user_roles":      roleNames,
		},
		CreatedAt: time.Now(),
	}

	// Try to get IP and user agent from context
	if ctx != nil {
		if req, ok := ctx.Value("http_request").(*http.Request); ok {
			auditLog.IPAddress = getIPAddress(req)
			auditLog.UserAgent = getUserAgent(req)
		}
	}

	if err := s.rbacRepo.LogPermissionCheck(context.Background(), auditLog); err != nil {
		s.logger.Error("Failed to log permission check", zap.Error(err))
	}

	return &models.PermissionCheckResult{
		HasPermission: hasPermission,
		Permission:    permission,
		Roles:         roleNames,
	}, nil
}

func (s *RBACService) GetUserPermissions(userID uuid.UUID) (*models.UserPermissionsResponse, error) {
	// Get user roles
	roles, err := s.rbacRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	// Get user permissions
	permissionNames, err := s.rbacRepo.GetUserPermissions(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}

	// Get full permission details
	var permissions []models.PermissionWithContext
	for _, permName := range permissionNames {
		perm, err := s.rbacRepo.GetPermissionByName(permName)
		if err != nil {
			s.logger.Warn("Failed to get permission details", zap.String("permission", permName), zap.Error(err))
			continue
		}
		if perm != nil {
			// TODO: Determine which roles grant this permission
			permissions = append(permissions, models.PermissionWithContext{
				Permission: *perm,
				GrantedVia: []string{"unknown"}, // Would need additional query
			})
		}
	}

	return &models.UserPermissionsResponse{
		UserID:      userID,
		Permissions: permissions,
		Roles:       roles,
	}, nil
}

// ==================== ROLE INHERITANCE ====================

func (s *RBACService) AddRoleInheritance(childRoleID, parentRoleID, addedBy uuid.UUID) error {
	// Check if child role exists
	childRole, err := s.rbacRepo.GetRoleByID(childRoleID)
	if err != nil {
		return fmt.Errorf("failed to get child role: %w", err)
	}
	if childRole == nil {
		return fmt.Errorf("child role not found")
	}

	// Check if parent role exists
	parentRole, err := s.rbacRepo.GetRoleByID(parentRoleID)
	if err != nil {
		return fmt.Errorf("failed to get parent role: %w", err)
	}
	if parentRole == nil {
		return fmt.Errorf("parent role not found")
	}

	// Add inheritance
	if err := s.rbacRepo.AddRoleInheritance(childRoleID, parentRoleID); err != nil {
		return fmt.Errorf("failed to add role inheritance: %w", err)
	}

	s.logger.Info("Role inheritance added",
		zap.String("child_role_id", childRoleID.String()),
		zap.String("parent_role_id", parentRoleID.String()),
		zap.String("added_by", addedBy.String()),
	)

	return nil
}

func (s *RBACService) RemoveRoleInheritance(childRoleID, parentRoleID, removedBy uuid.UUID) error {
	// Remove inheritance
	if err := s.rbacRepo.RemoveRoleInheritance(childRoleID, parentRoleID); err != nil {
		return fmt.Errorf("failed to remove role inheritance: %w", err)
	}

	s.logger.Info("Role inheritance removed",
		zap.String("child_role_id", childRoleID.String()),
		zap.String("parent_role_id", parentRoleID.String()),
		zap.String("removed_by", removedBy.String()),
	)

	return nil
}

// ==================== INITIALIZATION ====================

func (s *RBACService) InitializeDefaultRolesAndPermissions() error {
	s.logger.Info("Initializing default RBAC data")

	// Create default permission categories if they don't exist
	// Create default permissions for each system role
	// Assign default permissions to system roles
	// Ensure all users have the default 'user' role

	// This would be called during application startup or migration
	return nil
}

// ==================== HELPER FUNCTIONS ====================

func isValidRoleName(name string) bool {
	// Allow lowercase letters, numbers, underscores
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_') {
			return false
		}
	}
	return len(name) >= 3 && len(name) <= 100
}

func isValidPermissionName(name string) bool {
	// Format: resource:action
	parts := strings.Split(name, ":")
	if len(parts) != 2 {
		return false
	}

	// Check resource and action
	resource, action := parts[0], parts[1]
	if len(resource) < 2 || len(resource) > 100 || len(action) < 2 || len(action) > 50 {
		return false
	}

	// Allow lowercase letters, numbers, underscores, colons
	validChars := func(s string) bool {
		for _, r := range s {
			if !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' || r == ':') {
				return false
			}
		}
		return true
	}

	return validChars(resource) && validChars(action)
}

func getIPAddress(req *http.Request) *string {
	// Get IP from X-Forwarded-For header if behind proxy
	ip := req.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = req.RemoteAddr
	}

	// Extract IP from "host:port" format
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		ip = ip[:idx]
	}

	return &ip
}

func getUserAgent(req *http.Request) *string {
	ua := req.Header.Get("User-Agent")
	if ua == "" {
		return nil
	}
	return &ua
}

// ==================== BULK OPERATIONS ====================

func (s *RBACService) BulkUpdateRolePermissions(roleID uuid.UUID, permissionIDs []uuid.UUID, updatedBy uuid.UUID) error {
	// Get current permissions
	currentPerms, err := s.rbacRepo.GetRolePermissions(roleID)
	if err != nil {
		return fmt.Errorf("failed to get current permissions: %w", err)
	}

	// Create sets for comparison
	currentSet := make(map[uuid.UUID]bool)
	for _, perm := range currentPerms {
		currentSet[perm.ID] = true
	}

	newSet := make(map[uuid.UUID]bool)
	for _, permID := range permissionIDs {
		newSet[permID] = true
	}

	// Determine permissions to add and remove
	var toAdd, toRemove []uuid.UUID
	for permID := range newSet {
		if !currentSet[permID] {
			toAdd = append(toAdd, permID)
		}
	}
	for permID := range currentSet {
		if !newSet[permID] {
			toRemove = append(toRemove, permID)
		}
	}

	// Execute in transaction
	// Note: In a real implementation, you'd use a transaction
	for _, permID := range toRemove {
		if err := s.rbacRepo.RevokePermissionFromRole(roleID, permID); err != nil {
			s.logger.Warn("Failed to revoke permission", zap.String("permission_id", permID.String()), zap.Error(err))
		}
	}

	if err := s.rbacRepo.BulkGrantPermissions(roleID, toAdd, updatedBy); err != nil {
		return fmt.Errorf("failed to bulk grant permissions: %w", err)
	}

	s.logger.Info("Bulk updated role permissions",
		zap.String("role_id", roleID.String()),
		zap.Int("added", len(toAdd)),
		zap.Int("removed", len(toRemove)),
		zap.String("updated_by", updatedBy.String()),
	)

	return nil
}
