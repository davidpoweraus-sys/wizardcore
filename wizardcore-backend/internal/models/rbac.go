package models

import (
	"time"

	"github.com/google/uuid"
)

// ==================== CORE MODELS ====================

type Role struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name" validate:"required,min=3,max=100"`
	Description  *string   `json:"description,omitempty" db:"description"`
	IsSystemRole bool      `json:"is_system_role" db:"is_system_role"`
	IsDefault    bool      `json:"is_default" db:"is_default"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PermissionCategory struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name" validate:"required,min=3,max=100"`
	Description *string   `json:"description,omitempty" db:"description"`
	SortOrder   int       `json:"sort_order" db:"sort_order"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type Permission struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty" db:"category_id"`
	Name        string     `json:"name" db:"name" validate:"required,min=3,max=100"`
	Description *string    `json:"description,omitempty" db:"description"`
	Resource    string     `json:"resource" db:"resource_type" validate:"required,min=2,max=100"`
	Action      string     `json:"action" db:"action" validate:"required,min=2,max=50"`
	IsDangerous bool       `json:"is_dangerous" db:"is_dangerous"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`

	// Optional joined fields
	Category *PermissionCategory `json:"category,omitempty" db:"-"`
}

// ==================== RELATIONSHIP MODELS ====================

type RolePermission struct {
	RoleID       uuid.UUID  `json:"role_id" db:"role_id"`
	PermissionID uuid.UUID  `json:"permission_id" db:"permission_id"`
	GrantedAt    time.Time  `json:"granted_at" db:"granted_at"`
	GrantedBy    *uuid.UUID `json:"granted_by,omitempty" db:"granted_by"`
}

type UserRole struct {
	UserID     uuid.UUID  `json:"user_id" db:"user_id"`
	RoleID     uuid.UUID  `json:"role_id" db:"role_id"`
	AssignedAt time.Time  `json:"assigned_at" db:"assigned_at"`
	AssignedBy *uuid.UUID `json:"assigned_by,omitempty" db:"assigned_by"`
	ExpiresAt  *time.Time `json:"expires_at,omitempty" db:"expires_at"`

	// Optional joined fields
	Role *Role `json:"role,omitempty" db:"-"`
}

type RoleInheritance struct {
	ChildRoleID  uuid.UUID `json:"child_role_id" db:"child_role_id"`
	ParentRoleID uuid.UUID `json:"parent_role_id" db:"parent_role_id"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ==================== AUDIT MODELS ====================

type PermissionAuditLog struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	PermissionID uuid.UUID              `json:"permission_id" db:"permission_id"`
	ResourceID   *string                `json:"resource_id,omitempty" db:"resource_id"`
	Action       string                 `json:"action" db:"action" validate:"required"`
	Status       string                 `json:"status" db:"status" validate:"required,oneof=granted denied error"`
	IPAddress    *string                `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent    *string                `json:"user_agent,omitempty" db:"user_agent"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

type RoleAuditLog struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserID       uuid.UUID `json:"user_id" db:"user_id"`
	TargetUserID uuid.UUID `json:"target_user_id" db:"target_user_id"`
	RoleID       uuid.UUID `json:"role_id" db:"role_id"`
	Action       string    `json:"action" db:"action" validate:"required,oneof=assigned removed updated"`
	Reason       *string   `json:"reason,omitempty" db:"reason"`
	PerformedBy  uuid.UUID `json:"performed_by" db:"performed_by"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

// ==================== REQUEST/RESPONSE MODELS ====================

type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,min=3,max=100,alphanumunderscore"`
	Description *string `json:"description,omitempty"`
	IsDefault   bool    `json:"is_default"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=3,max=100,alphanumunderscore"`
	Description *string `json:"description,omitempty"`
	IsDefault   *bool   `json:"is_default,omitempty"`
}

type CreatePermissionRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=100,alphanumcolon"`
	Description *string    `json:"description,omitempty"`
	Resource    string     `json:"resource" validate:"required,min=2,max=100"`
	Action      string     `json:"action" validate:"required,min=2,max=50"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	IsDangerous bool       `json:"is_dangerous"`
}

type AssignRoleRequest struct {
	UserID    uuid.UUID  `json:"user_id" validate:"required"`
	RoleID    uuid.UUID  `json:"role_id" validate:"required"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Reason    *string    `json:"reason,omitempty"`
}

type GrantPermissionRequest struct {
	RoleID       uuid.UUID `json:"role_id" validate:"required"`
	PermissionID uuid.UUID `json:"permission_id" validate:"required"`
}

type SetRoleInheritanceRequest struct {
	ChildRoleID  uuid.UUID `json:"child_role_id" validate:"required"`
	ParentRoleID uuid.UUID `json:"parent_role_id" validate:"required"`
}

type CheckPermissionRequest struct {
	PermissionName string  `json:"permission_name" validate:"required"`
	ResourceID     *string `json:"resource_id,omitempty"`
}

type PermissionCheckResult struct {
	HasPermission bool        `json:"has_permission"`
	Permission    *Permission `json:"permission,omitempty"`
	Roles         []string    `json:"roles,omitempty"`
}

type UserPermissionsResponse struct {
	UserID      uuid.UUID               `json:"user_id"`
	Permissions []PermissionWithContext `json:"permissions"`
	Roles       []Role                  `json:"roles"`
}

type PermissionWithContext struct {
	Permission
	GrantedVia []string `json:"granted_via"` // Role names that grant this permission
}

// ==================== HELPER TYPES ====================

type PermissionMap map[string]bool // permission_name -> has_permission

type RoleWithPermissions struct {
	Role
	Permissions []Permission `json:"permissions"`
	UserCount   int          `json:"user_count"`
}

type UserWithRoles struct {
	User
	Roles []Role `json:"roles"`
}

// ==================== CONSTANTS ====================

const (
	// Default system role names
	RoleSuperAdmin     = "super_admin"
	RoleAdmin          = "admin"
	RoleUser           = "user"
	RoleContentCreator = "content_creator"
	RoleModerator      = "moderator"

	// Common permission actions
	ActionCreate = "create"
	ActionRead   = "read"
	ActionUpdate = "update"
	ActionDelete = "delete"
	ActionManage = "manage"
	ActionExport = "export"
	ActionImport = "import"

	// Common resource types
	ResourceUser       = "user"
	ResourceRole       = "role"
	ResourcePermission = "permission"
	ResourcePathway    = "pathway"
	ResourceModule     = "module"
	ResourceExercise   = "exercise"
	ResourceSubmission = "submission"
	ResourceSystem     = "system"

	// Audit statuses
	AuditStatusGranted = "granted"
	AuditStatusDenied  = "denied"
	AuditStatusError   = "error"

	// Role audit actions
	RoleActionAssigned = "assigned"
	RoleActionRemoved  = "removed"
	RoleActionUpdated  = "updated"
)

// ==================== HELPER FUNCTIONS ====================

// BuildPermissionName creates a standardized permission name
func BuildPermissionName(resource, action string) string {
	return resource + ":" + action
}

// ParsePermissionName parses a permission name into resource and action
func ParsePermissionName(permissionName string) (resource, action string, ok bool) {
	for i := len(permissionName) - 1; i >= 0; i-- {
		if permissionName[i] == ':' {
			return permissionName[:i], permissionName[i+1:], true
		}
	}
	return "", "", false
}

// IsSystemRole checks if a role name is a system role
func IsSystemRole(roleName string) bool {
	switch roleName {
	case RoleSuperAdmin, RoleAdmin, RoleUser:
		return true
	default:
		return false
	}
}

// GetDefaultPermissions returns default permissions for common roles
func GetDefaultPermissions(roleName string) []string {
	switch roleName {
	case RoleSuperAdmin:
		return []string{
			BuildPermissionName(ResourceSystem, ActionManage),
			BuildPermissionName(ResourceUser, ActionManage),
			BuildPermissionName(ResourceRole, ActionManage),
			BuildPermissionName(ResourcePermission, ActionManage),
		}
	case RoleAdmin:
		return []string{
			BuildPermissionName(ResourceUser, ActionManage),
			BuildPermissionName(ResourcePathway, ActionManage),
			BuildPermissionName(ResourceModule, ActionManage),
			BuildPermissionName(ResourceExercise, ActionManage),
			BuildPermissionName(ResourceSubmission, ActionRead),
		}
	case RoleContentCreator:
		return []string{
			BuildPermissionName(ResourcePathway, ActionCreate),
			BuildPermissionName(ResourcePathway, ActionUpdate),
			BuildPermissionName(ResourceModule, ActionCreate),
			BuildPermissionName(ResourceModule, ActionUpdate),
			BuildPermissionName(ResourceExercise, ActionCreate),
			BuildPermissionName(ResourceExercise, ActionUpdate),
		}
	case RoleModerator:
		return []string{
			BuildPermissionName(ResourceUser, ActionRead),
			BuildPermissionName(ResourceSubmission, ActionRead),
			BuildPermissionName(ResourceSubmission, ActionUpdate),
			BuildPermissionName(ResourcePathway, ActionUpdate),
		}
	case RoleUser:
		return []string{
			BuildPermissionName(ResourceUser, ActionRead),   // Own profile
			BuildPermissionName(ResourceUser, ActionUpdate), // Own profile
			BuildPermissionName(ResourcePathway, ActionRead),
			BuildPermissionName(ResourceModule, ActionRead),
			BuildPermissionName(ResourceExercise, ActionRead),
			BuildPermissionName(ResourceSubmission, ActionCreate),
			BuildPermissionName(ResourceSubmission, ActionRead), // Own submissions
		}
	default:
		return []string{}
	}
}
