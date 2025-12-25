-- Create default permissions for testing
-- These are common permissions needed for the application

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

-- Assign permissions to roles
-- User role gets basic permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'user'),
    id
FROM permissions 
WHERE name IN ('user:read', 'user:update', 'content:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Content creator role gets content creation permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'content_creator'),
    id
FROM permissions 
WHERE name IN ('content:create', 'content:update', 'content:read', 'user:read', 'user:update')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Moderator role gets moderation permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'moderator'),
    id
FROM permissions 
WHERE name IN ('moderation:review', 'moderation:action', 'content:read', 'user:read')
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Admin role gets most permissions (except super_admin)
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'admin'),
    id
FROM permissions 
WHERE name NOT LIKE 'system:%' AND name != 'rbac:manage'
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Super admin gets all permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT 
    (SELECT id FROM roles WHERE name = 'super_admin'),
    id
FROM permissions 
ON CONFLICT (role_id, permission_id) DO NOTHING;

-- Display what was created
SELECT 
    r.name as role_name,
    COUNT(rp.permission_id) as permission_count
FROM roles r
LEFT JOIN role_permissions rp ON r.id = rp.role_id
GROUP BY r.name
ORDER BY permission_count DESC;