-- Rollback RBAC system migration

-- Drop triggers
DROP TRIGGER IF EXISTS update_permissions_updated_at ON permissions;
DROP TRIGGER IF EXISTS update_roles_updated_at ON roles;

-- Drop functions
DROP FUNCTION IF EXISTS update_updated_at_column CASCADE;
DROP FUNCTION IF EXISTS get_user_permissions CASCADE;
DROP FUNCTION IF EXISTS check_user_permission CASCADE;

-- Drop audit tables
DROP TABLE IF EXISTS role_audit_log CASCADE;
DROP TABLE IF EXISTS permission_audit_log CASCADE;

-- Drop relationship tables
DROP TABLE IF EXISTS role_inheritance CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;

-- Drop core tables
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS permission_categories CASCADE;
DROP TABLE IF EXISTS roles CASCADE;

-- Remove is_active column from users
ALTER TABLE users DROP COLUMN IF EXISTS is_active;
