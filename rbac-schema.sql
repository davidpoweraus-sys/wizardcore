-- Scalable RBAC Schema for WizardCore
-- Supports: Role inheritance, permission grouping, dynamic role creation

-- ==================== CORE TABLES ====================

-- Users table (existing, with role support added)
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;

-- Roles table
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

CREATE INDEX idx_roles_name ON roles(name);
CREATE INDEX idx_permissions_name ON permissions(name);
CREATE INDEX idx_permissions_resource_action ON permissions(resource_type, action);
CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);
CREATE INDEX idx_role_permissions_role_id ON role_permissions(role_id);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);
CREATE INDEX idx_permission_audit_user_time ON permission_audit_log(user_id, created_at DESC);
CREATE INDEX idx_role_audit_target_time ON role_audit_log(target_user_id, created_at DESC);

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
INSERT INTO roles (id, name, description, is_system_role, is_default) VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'super_admin', 'Full system access', true, false),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'admin', 'Administrative access', true, false),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'user', 'Basic user access', true, true),
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'content_creator', 'Create and manage content', false, false),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'moderator', 'Moderate user content', false, false)
ON CONFLICT (name) DO NOTHING;

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

-- ==================== SECURITY POLICIES ====================

-- Row Level Security (optional, for extra security)
-- ALTER TABLE permissions ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE roles ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE user_roles ENABLE ROW LEVEL SECURITY;

-- ==================== MIGRATION NOTES ====================

/*
To migrate existing users:
1. All existing users get the 'user' role automatically
2. You can manually assign admin roles to specific users
3. Update your application code to use the new permission system
*/