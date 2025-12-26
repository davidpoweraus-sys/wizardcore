-- ==================== MIGRATION 013: FULL RBAC SYSTEM ====================
-- Implements comprehensive Role-Based Access Control with permissions, role inheritance, and audit logging
-- This migration extends the simple 'role' column added in migration 012

-- ==================== CORE TABLES ====================

-- Add is_active column to users (for account suspension)
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT true;

-- Roles table (system for managing user roles)
CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_system_role BOOLEAN DEFAULT false, -- Cannot be deleted (e.g., admin, user)
    is_default BOOLEAN DEFAULT false, -- Auto-assigned to new users
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Permission categories for organization
CREATE TABLE IF NOT EXISTS permission_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Permissions table (fine-grained access controls)
CREATE TABLE IF NOT EXISTS permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID REFERENCES permission_categories(id) ON DELETE SET NULL,
    name VARCHAR(100) UNIQUE NOT NULL, -- e.g., "exercise:create"
    description TEXT,
    resource_type VARCHAR(100), -- e.g., "exercise", "user", "pathway"
    action VARCHAR(50), -- e.g., "create", "read", "update", "delete", "manage"
    is_dangerous BOOLEAN DEFAULT false, -- Requires extra approval
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==================== RELATIONSHIP TABLES ====================

-- Role-Permission mapping (many-to-many)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    permission_id UUID REFERENCES permissions(id) ON DELETE CASCADE,
    granted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    granted_by UUID REFERENCES users(id),
    PRIMARY KEY (role_id, permission_id)
);

-- User-Role assignment (many-to-many)
CREATE TABLE IF NOT EXISTS user_roles (
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by UUID REFERENCES users(id),
    expires_at TIMESTAMP, -- Optional: temporary role assignment
    PRIMARY KEY (user_id, role_id)
);

-- Role inheritance (roles can inherit from other roles)
CREATE TABLE IF NOT EXISTS role_inheritance (
    child_role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    parent_role_id UUID REFERENCES roles(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (child_role_id, parent_role_id),
    CHECK (child_role_id != parent_role_id) -- Prevent self-reference
);

-- ==================== AUDIT TABLES ====================

-- Permission usage audit log
CREATE TABLE IF NOT EXISTS permission_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    permission_id UUID REFERENCES permissions(id),
    resource_id VARCHAR(255), -- ID of the resource being accessed
    action VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL, -- "granted", "denied", "error"
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Role change audit log
CREATE TABLE IF NOT EXISTS role_audit_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id),
    target_user_id UUID REFERENCES users(id),
    role_id UUID REFERENCES roles(id),
    action VARCHAR(20) NOT NULL, -- "assigned", "removed", "updated"
    reason TEXT,
    performed_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ==================== INDEXES ====================

CREATE INDEX IF NOT EXISTS idx_roles_name ON roles(name);
CREATE INDEX IF NOT EXISTS idx_permissions_name ON permissions(name);
CREATE INDEX IF NOT EXISTS idx_permissions_resource_action ON permissions(resource_type, action);
CREATE INDEX IF NOT EXISTS idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX IF NOT EXISTS idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX IF NOT EXISTS idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX IF NOT EXISTS idx_permission_audit_user_time ON permission_audit_log(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_role_audit_target_time ON role_audit_log(target_user_id, created_at DESC);

-- ==================== DEFAULT DATA ====================

-- Insert default permission categories
INSERT INTO permission_categories (id, name, description, sort_order) VALUES
    ('11111111-1111-1111-1111-111111111111', 'User Management', 'User account and profile operations', 10),
    ('22222222-2222-2222-2222-222222222222', 'Content Management', 'Course and exercise content', 20),
    ('33333333-3333-3333-3333-333333333333', 'System Administration', 'Platform configuration and maintenance', 30),
    ('44444444-4444-4444-4444-444444444444', 'Analytics', 'View platform statistics and reports', 40),
    ('55555555-5555-5555-5555-555555555555', 'Moderation', 'Content and user moderation', 50)
ON CONFLICT (name) DO NOTHING;

-- Insert system roles (cannot be deleted)
-- Note: These map to the simple role column values from migration 012
INSERT INTO roles (id, name, description, is_system_role, is_default) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'super_admin', 'Full system access', true, false),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'admin', 'Administrative access', true, false),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'student', 'Basic student access', true, true),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'content_creator', 'Create and manage content', true, false),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'moderator', 'Moderate user content', false, false)
ON CONFLICT (name) DO NOTHING;

-- ==================== DEFAULT PERMISSIONS ====================

-- User Management permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'user:read', 'Read user profiles', 'user', 'read', false),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'user:update', 'Update user profiles', 'user', 'update', false),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'user:delete', 'Delete users', 'user', 'delete', true),
    (gen_random_uuid(), '11111111-1111-1111-1111-111111111111', 'user:manage', 'Manage all users', 'user', 'manage', true)
ON CONFLICT (name) DO NOTHING;

-- Content Management permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'content:read', 'Read content', 'content', 'read', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'content:create', 'Create content', 'content', 'create', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'content:update', 'Update content', 'content', 'update', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'content:delete', 'Delete content', 'content', 'delete', true),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'content:manage', 'Manage all content', 'content', 'manage', true)
ON CONFLICT (name) DO NOTHING;

-- Exercise permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'exercise:read', 'Read exercises', 'exercise', 'read', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'exercise:create', 'Create exercises', 'exercise', 'create', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'exercise:update', 'Update exercises', 'exercise', 'update', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'exercise:delete', 'Delete exercises', 'exercise', 'delete', true),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'exercise:submit', 'Submit exercise solutions', 'exercise', 'submit', false)
ON CONFLICT (name) DO NOTHING;

-- Pathway permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'pathway:read', 'Read pathways', 'pathway', 'read', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'pathway:create', 'Create pathways', 'pathway', 'create', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'pathway:update', 'Update pathways', 'pathway', 'update', false),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'pathway:delete', 'Delete pathways', 'pathway', 'delete', true),
    (gen_random_uuid(), '22222222-2222-2222-2222-222222222222', 'pathway:enroll', 'Enroll in pathways', 'pathway', 'enroll', false)
ON CONFLICT (name) DO NOTHING;

-- System Administration permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'system:admin', 'Full system administration', 'system', 'admin', true),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'system:config', 'Configure system settings', 'system', 'config', true),
    (gen_random_uuid(), '33333333-3333-3333-3333-333333333333', 'rbac:manage', 'Manage roles and permissions', 'rbac', 'manage', true)
ON CONFLICT (name) DO NOTHING;

-- Analytics permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'analytics:view', 'View analytics', 'analytics', 'view', false),
    (gen_random_uuid(), '44444444-4444-4444-4444-444444444444', 'analytics:export', 'Export analytics data', 'analytics', 'export', false)
ON CONFLICT (name) DO NOTHING;

-- Moderation permissions
INSERT INTO permissions (id, category_id, name, description, resource_type, action, is_dangerous) VALUES
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 'moderation:review', 'Review user content', 'moderation', 'review', false),
    (gen_random_uuid(), '55555555-5555-5555-5555-555555555555', 'moderation:action', 'Take moderation actions', 'moderation', 'action', true)
ON CONFLICT (name) DO NOTHING;

-- ==================== ASSIGN PERMISSIONS TO ROLES ====================

-- Student role gets basic permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'student'),
    id
FROM permissions 
WHERE name IN (
    'user:read', 'user:update',
    'content:read', 'exercise:read', 'exercise:submit',
    'pathway:read', 'pathway:enroll'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Content creator role gets content creation permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'content_creator'),
    id
FROM permissions 
WHERE name IN (
    'user:read', 'user:update',
    'content:read', 'content:create', 'content:update',
    'exercise:read', 'exercise:create', 'exercise:update', 'exercise:submit',
    'pathway:read', 'pathway:create', 'pathway:update', 'pathway:enroll',
    'analytics:view'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Moderator role gets moderation permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'moderator'),
    id
FROM permissions 
WHERE name IN (
    'user:read',
    'content:read', 'content:update',
    'exercise:read',
    'pathway:read',
    'moderation:review', 'moderation:action'
)
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Admin role gets most permissions (except super_admin)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'),
    id
FROM permissions 
WHERE name NOT LIKE 'system:%' AND name != 'rbac:manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Super admin gets ALL permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'super_admin'),
    id
FROM permissions 
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- ==================== HELPER FUNCTIONS ====================

-- Function to check if user has permission (including inherited roles)
CREATE OR REPLACE FUNCTION check_user_permission(
    p_user_id UUID,
    p_permission_name VARCHAR(100)
) RETURNS BOOLEAN AS $$
DECLARE
    has_perm BOOLEAN;
BEGIN
    WITH RECURSIVE user_roles_tree AS (
        -- Get user's direct roles
        SELECT ur.role_id
        FROM user_roles ur
        WHERE ur.user_id = p_user_id
          AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP)
        
        UNION
        
        -- Get inherited roles (recursive)
        SELECT ri.parent_role_id
        FROM user_roles_tree urt
        JOIN role_inheritance ri ON urt.role_id = ri.child_role_id
    )
    SELECT EXISTS (
        SELECT 1
        FROM user_roles_tree urt
        JOIN role_permissions rp ON urt.role_id = rp.role_id
        JOIN permissions p ON rp.permission_id = p.id
        WHERE p.name = p_permission_name
    ) INTO has_perm;
    
    RETURN has_perm;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- Function to get all permissions for a user
CREATE OR REPLACE FUNCTION get_user_permissions(p_user_id UUID)
RETURNS TABLE (
    permission_name VARCHAR(100),
    permission_description TEXT,
    resource_type VARCHAR(100),
    action VARCHAR(50),
    category_name VARCHAR(100)
) AS $$
BEGIN
    RETURN QUERY
    WITH RECURSIVE user_roles_tree AS (
        SELECT ur.role_id
        FROM user_roles ur
        WHERE ur.user_id = p_user_id
          AND (ur.expires_at IS NULL OR ur.expires_at > CURRENT_TIMESTAMP)
        
        UNION
        
        SELECT ri.parent_role_id
        FROM user_roles_tree urt
        JOIN role_inheritance ri ON urt.role_id = ri.child_role_id
    )
    SELECT DISTINCT
        p.name,
        p.description,
        p.resource_type,
        p.action,
        pc.name
    FROM user_roles_tree urt
    JOIN role_permissions rp ON urt.role_id = rp.role_id
    JOIN permissions p ON rp.permission_id = p.id
    LEFT JOIN permission_categories pc ON p.category_id = pc.id
    ORDER BY pc.sort_order, p.resource_type, p.action;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- ==================== TRIGGERS ====================

-- Auto-update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON roles
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_permissions_updated_at
    BEFORE UPDATE ON permissions
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ==================== MIGRATION NOTES ====================

-- This migration adds comprehensive RBAC on top of the simple 'role' column from migration 012
-- The simple 'role' column can be used for backwards compatibility
-- New code should use the user_roles table for role assignment
-- Existing users will need to have their roles migrated to the new system via a separate data migration
