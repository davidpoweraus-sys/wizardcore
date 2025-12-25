package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"go.uber.org/zap"
)

type RBACRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRBACRepository(db *sql.DB, logger *zap.Logger) *RBACRepository {
	return &RBACRepository{
		db:     db,
		logger: logger,
	}
}

// ==================== ROLE OPERATIONS ====================

func (r *RBACRepository) CreateRole(role *models.Role) error {
	query := `
		INSERT INTO roles (id, name, description, is_system_role, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		role.ID,
		role.Name,
		role.Description,
		role.IsSystemRole,
		role.IsDefault,
		now,
		now,
	).Scan(&role.ID, &role.CreatedAt, &role.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create role: %w", err)
	}

	return nil
}

func (r *RBACRepository) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	query := `
		SELECT id, name, description, is_system_role, is_default, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.IsSystemRole,
		&role.IsDefault,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return role, nil
}

func (r *RBACRepository) GetRoleByName(name string) (*models.Role, error) {
	query := `
		SELECT id, name, description, is_system_role, is_default, created_at, updated_at
		FROM roles
		WHERE name = $1
	`

	role := &models.Role{}
	err := r.db.QueryRow(query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.IsSystemRole,
		&role.IsDefault,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get role by name: %w", err)
	}

	return role, nil
}

func (r *RBACRepository) UpdateRole(role *models.Role) error {
	query := `
		UPDATE roles
		SET name = $1, description = $2, is_default = $3, updated_at = $4
		WHERE id = $5 AND is_system_role = false
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		query,
		role.Name,
		role.Description,
		role.IsDefault,
		time.Now(),
		role.ID,
	).Scan(&role.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update role: %w", err)
	}

	return nil
}

func (r *RBACRepository) DeleteRole(id uuid.UUID) error {
	// Check if it's a system role
	var isSystemRole bool
	err := r.db.QueryRow("SELECT is_system_role FROM roles WHERE id = $1", id).Scan(&isSystemRole)
	if err != nil {
		return fmt.Errorf("failed to check role: %w", err)
	}

	if isSystemRole {
		return fmt.Errorf("cannot delete system role")
	}

	// Delete role (cascade will handle role_permissions and user_roles)
	query := "DELETE FROM roles WHERE id = $1"
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

func (r *RBACRepository) ListRoles(page, pageSize int, includeSystem bool) ([]models.Role, int, error) {
	// Build query
	baseQuery := "SELECT id, name, description, is_system_role, is_default, created_at, updated_at FROM roles"
	countQuery := "SELECT COUNT(*) FROM roles"

	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if !includeSystem {
		whereClauses = append(whereClauses, fmt.Sprintf("is_system_role = $%d", argIdx))
		args = append(args, false)
		argIdx++
	}

	if len(whereClauses) > 0 {
		where := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += where
		countQuery += where
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count roles: %w", err)
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" ORDER BY name LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	// Execute query
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.IsSystemRole,
			&role.IsDefault,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, total, nil
}

// ==================== PERMISSION OPERATIONS ====================

func (r *RBACRepository) CreatePermission(permission *models.Permission) error {
	query := `
		INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRow(
		query,
		permission.ID,
		permission.CategoryID,
		permission.Name,
		permission.Description,
		permission.Resource,
		permission.Action,
		permission.IsDangerous,
		now,
		now,
	).Scan(&permission.CreatedAt, &permission.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create permission: %w", err)
	}

	return nil
}

func (r *RBACRepository) GetPermissionByID(id uuid.UUID) (*models.Permission, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.description, p.resource_type, p.action, p.is_dangerous, p.created_at, p.updated_at,
		       pc.id, pc.name, pc.description, pc.sort_order, pc.created_at
		FROM permissions p
		LEFT JOIN permission_categories pc ON p.category_id = pc.id
		WHERE p.id = $1
	`

	permission := &models.Permission{}
	var categoryID sql.NullString
	var categoryName, categoryDesc sql.NullString
	var categorySort sql.NullInt32
	var categoryCreated sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&permission.ID,
		&permission.CategoryID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.IsDangerous,
		&permission.CreatedAt,
		&permission.UpdatedAt,
		&categoryID,
		&categoryName,
		&categoryDesc,
		&categorySort,
		&categoryCreated,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission: %w", err)
	}

	// Set category if exists
	if categoryID.Valid {
		catID, _ := uuid.Parse(categoryID.String)
		permission.Category = &models.PermissionCategory{
			ID:          catID,
			Name:        categoryName.String,
			Description: &categoryDesc.String,
			SortOrder:   int(categorySort.Int32),
			CreatedAt:   categoryCreated.Time,
		}
	}

	return permission, nil
}

func (r *RBACRepository) GetPermissionByName(name string) (*models.Permission, error) {
	query := `
		SELECT id, category_id, name, description, resource_type, action, is_dangerous, created_at, updated_at
		FROM permissions
		WHERE name = $1
	`

	permission := &models.Permission{}
	err := r.db.QueryRow(query, name).Scan(
		&permission.ID,
		&permission.CategoryID,
		&permission.Name,
		&permission.Description,
		&permission.Resource,
		&permission.Action,
		&permission.IsDangerous,
		&permission.CreatedAt,
		&permission.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get permission by name: %w", err)
	}

	return permission, nil
}

func (r *RBACRepository) ListPermissions(page, pageSize int, categoryID *uuid.UUID) ([]models.Permission, int, error) {
	// Build query
	baseQuery := `
		SELECT p.id, p.category_id, p.name, p.description, p.resource_type, p.action, p.is_dangerous, p.created_at, p.updated_at,
		       pc.id, pc.name, pc.description, pc.sort_order, pc.created_at
		FROM permissions p
		LEFT JOIN permission_categories pc ON p.category_id = pc.id
	`
	countQuery := "SELECT COUNT(*) FROM permissions p"

	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if categoryID != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("p.category_id = $%d", argIdx))
		args = append(args, categoryID)
		argIdx++
	}

	if len(whereClauses) > 0 {
		where := " WHERE " + strings.Join(whereClauses, " AND ")
		baseQuery += where
		countQuery += where
	}

	// Get total count
	var total int
	err := r.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count permissions: %w", err)
	}

	// Add pagination and ordering
	baseQuery += " ORDER BY pc.sort_order, p.resource_type, p.action"
	baseQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, pageSize, (page-1)*pageSize)

	// Execute query
	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		var categoryID sql.NullString
		var categoryName, categoryDesc sql.NullString
		var categorySort sql.NullInt32
		var categoryCreated sql.NullTime

		err := rows.Scan(
			&permission.ID,
			&permission.CategoryID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.IsDangerous,
			&permission.CreatedAt,
			&permission.UpdatedAt,
			&categoryID,
			&categoryName,
			&categoryDesc,
			&categorySort,
			&categoryCreated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan permission: %w", err)
		}

		// Set category if exists
		if categoryID.Valid {
			catID, _ := uuid.Parse(categoryID.String)
			permission.Category = &models.PermissionCategory{
				ID:          catID,
				Name:        categoryName.String,
				Description: &categoryDesc.String,
				SortOrder:   int(categorySort.Int32),
				CreatedAt:   categoryCreated.Time,
			}
		}

		permissions = append(permissions, permission)
	}

	return permissions, total, nil
}

// ==================== ROLE-PERMISSION OPERATIONS ====================

func (r *RBACRepository) GrantPermissionToRole(roleID, permissionID, grantedBy uuid.UUID) error {
	query := `
		INSERT INTO role_permissions (role_id, permission_id, granted_at, granted_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (role_id, permission_id) DO UPDATE SET granted_at = $3, granted_by = $4
	`

	_, err := r.db.Exec(query, roleID, permissionID, time.Now(), grantedBy)
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}

	return nil
}

func (r *RBACRepository) RevokePermissionFromRole(roleID, permissionID uuid.UUID) error {
	query := "DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2"

	result, err := r.db.Exec(query, roleID, permissionID)
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("permission not found for role")
	}

	return nil
}

func (r *RBACRepository) GetRolePermissions(roleID uuid.UUID) ([]models.Permission, error) {
	query := `
		SELECT p.id, p.category_id, p.name, p.description, p.resource_type, p.action, p.is_dangerous, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource_type, p.action
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role permissions: %w", err)
	}
	defer rows.Close()

	var permissions []models.Permission
	for rows.Next() {
		var permission models.Permission
		err := rows.Scan(
			&permission.ID,
			&permission.CategoryID,
			&permission.Name,
			&permission.Description,
			&permission.Resource,
			&permission.Action,
			&permission.IsDangerous,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// ==================== USER-ROLE OPERATIONS ====================

func (r *RBACRepository) AssignRoleToUser(userID, roleID, assignedBy uuid.UUID, expiresAt *time.Time) error {
	query := `
		INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (user_id, role_id) DO UPDATE SET assigned_at = $3, assigned_by = $4, expires_at = $5
	`

	_, err := r.db.Exec(query, userID, roleID, time.Now(), assignedBy, expiresAt)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

func (r *RBACRepository) RemoveRoleFromUser(userID, roleID uuid.UUID) error {
	query := "DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2"

	result, err := r.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("role not found for user")
	}

	return nil
}

func (r *RBACRepository) GetUserRoles(userID uuid.UUID) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.is_system_role, r.is_default, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1 AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP)
		ORDER BY r.name
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.IsSystemRole,
			&role.IsDefault,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (r *RBACRepository) GetUsersWithRole(roleID uuid.UUID, page, pageSize int) ([]uuid.UUID, int, error) {
	// Count total
	countQuery := `
		SELECT COUNT(*)
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		WHERE ur.role_id = $1 AND u.is_active = true
		  AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP)
	`

	var total int
	err := r.db.QueryRow(countQuery, roleID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get user IDs
	query := `
		SELECT ur.user_id
		FROM user_roles ur
		JOIN users u ON ur.user_id = u.id
		WHERE ur.role_id = $1 AND u.is_active = true
		  AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP)
		ORDER BY ur.assigned_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, roleID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var userIDs []uuid.UUID
	for rows.Next() {
		var userID uuid.UUID
		err := rows.Scan(&userID)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user ID: %w", err)
		}
		userIDs = append(userIDs, userID)
	}

	return userIDs, total, nil
}

// ==================== PERMISSION CHECKING ====================

func (r *RBACRepository) CheckUserPermission(userID uuid.UUID, permissionName string) (bool, error) {
	query := `
		SELECT check_user_permission($1, $2)
	`

	var hasPermission bool
	err := r.db.QueryRow(query, userID, permissionName).Scan(&hasPermission)
	if err != nil {
		return false, fmt.Errorf("failed to check permission: %w", err)
	}

	return hasPermission, nil
}

func (r *RBACRepository) GetUserPermissions(userID uuid.UUID) ([]string, error) {
	query := `
		SELECT permission_name
		FROM get_user_permissions($1)
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user permissions: %w", err)
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var perm string
		err := rows.Scan(&perm)
		if err != nil {
			return nil, fmt.Errorf("failed to scan permission: %w", err)
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

func (r *RBACRepository) GetUserPermissionMap(userID uuid.UUID) (models.PermissionMap, error) {
	permissions, err := r.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	permMap := make(models.PermissionMap)
	for _, perm := range permissions {
		permMap[perm] = true
	}

	return permMap, nil
}

// ==================== ROLE INHERITANCE ====================

func (r *RBACRepository) AddRoleInheritance(childRoleID, parentRoleID uuid.UUID) error {
	// Check for circular inheritance
	if childRoleID == parentRoleID {
		return fmt.Errorf("cannot inherit from self")
	}

	// Check if parent inherits from child (circular)
	query := `
		WITH RECURSIVE inheritance_chain AS (
			SELECT child_role_id, parent_role_id
			FROM role_inheritance
			WHERE child_role_id = $1
			UNION
			SELECT ri.child_role_id, ri.parent_role_id
			FROM role_inheritance ri
			JOIN inheritance_chain ic ON ri.child_role_id = ic.parent_role_id
		)
		SELECT EXISTS (SELECT 1 FROM inheritance_chain WHERE parent_role_id = $2)
	`

	var circular bool
	err := r.db.QueryRow(query, parentRoleID, childRoleID).Scan(&circular)
	if err != nil {
		return fmt.Errorf("failed to check circular inheritance: %w", err)
	}

	if circular {
		return fmt.Errorf("circular inheritance detected")
	}

	// Add inheritance
	insertQuery := `
		INSERT INTO role_inheritance (child_role_id, parent_role_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (child_role_id, parent_role_id) DO NOTHING
	`

	_, err = r.db.Exec(insertQuery, childRoleID, parentRoleID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to add role inheritance: %w", err)
	}

	return nil
}

func (r *RBACRepository) RemoveRoleInheritance(childRoleID, parentRoleID uuid.UUID) error {
	query := "DELETE FROM role_inheritance WHERE child_role_id = $1 AND parent_role_id = $2"

	result, err := r.db.Exec(query, childRoleID, parentRoleID)
	if err != nil {
		return fmt.Errorf("failed to remove role inheritance: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("inheritance not found")
	}

	return nil
}

func (r *RBACRepository) GetRoleInheritance(roleID uuid.UUID) ([]models.Role, error) {
	query := `
		WITH RECURSIVE inheritance_chain AS (
			SELECT child_role_id, parent_role_id
			FROM role_inheritance
			WHERE child_role_id = $1
			UNION
			SELECT ri.child_role_id, ri.parent_role_id
			FROM role_inheritance ri
			JOIN inheritance_chain ic ON ri.child_role_id = ic.parent_role_id
		)
		SELECT r.id, r.name, r.description, r.is_system_role, r.is_default, r.created_at, r.updated_at
		FROM inheritance_chain ic
		JOIN roles r ON ic.parent_role_id = r.id
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get role inheritance: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.IsSystemRole,
			&role.IsDefault,
			&role.CreatedAt,
			&role.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// ==================== AUDIT LOGGING ====================

func (r *RBACRepository) LogPermissionCheck(ctx context.Context, log *models.PermissionAuditLog) error {
	query := `
		INSERT INTO permission_audit_log (id, user_id, permission_id, resource_id, action, status, ip_address, user_agent, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.UserID,
		log.PermissionID,
		log.ResourceID,
		log.Action,
		log.Status,
		log.IPAddress,
		log.UserAgent,
		log.Metadata,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to log permission check: %w", err)
	}

	return nil
}

func (r *RBACRepository) LogRoleChange(ctx context.Context, log *models.RoleAuditLog) error {
	query := `
		INSERT INTO role_audit_log (id, user_id, target_user_id, role_id, action, reason, performed_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.ExecContext(ctx, query,
		log.ID,
		log.UserID,
		log.TargetUserID,
		log.RoleID,
		log.Action,
		log.Reason,
		log.PerformedBy,
		log.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to log role change: %w", err)
	}

	return nil
}

// ==================== BULK OPERATIONS ====================

func (r *RBACRepository) BulkGrantPermissions(roleID uuid.UUID, permissionIDs []uuid.UUID, grantedBy uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	stmt, err := tx.Prepare(`
		INSERT INTO role_permissions (role_id, permission_id, granted_at, granted_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (role_id, permission_id) DO UPDATE SET granted_at = $3, granted_by = $4
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, permID := range permissionIDs {
		_, err := stmt.Exec(roleID, permID, now, grantedBy)
		if err != nil {
			return fmt.Errorf("failed to grant permission %s: %w", permID, err)
		}
	}

	return tx.Commit()
}

func (r *RBACRepository) BulkAssignRoles(userID uuid.UUID, roleIDs []uuid.UUID, assignedBy uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	now := time.Now()
	stmt, err := tx.Prepare(`
		INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, role_id) DO UPDATE SET assigned_at = $3, assigned_by = $4
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, roleID := range roleIDs {
		_, err := stmt.Exec(userID, roleID, now, assignedBy)
		if err != nil {
			return fmt.Errorf("failed to assign role %s: %w", roleID, err)
		}
	}

	return tx.Commit()
}
