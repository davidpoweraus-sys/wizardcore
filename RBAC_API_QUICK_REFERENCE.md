# RBAC API Quick Reference

## Overview
The RBAC (Role-Based Access Control) system provides fine-grained permission management for the Offensive Wizard platform.

## Default Roles and Permissions

### Roles
1. **super_admin** - Full system access (all permissions)
2. **admin** - Administrative access (most permissions except system-level)
3. **user** - Basic user access (default role for all users)
4. **content_creator** - Create and manage content
5. **moderator** - Moderate user content

### Default Permissions by Role

#### User Role (default)
- `user:read` - Read user profiles
- `user:update` - Update own profile
- `content:read` - Read content

#### Content Creator Role
- All user permissions plus:
- `content:create` - Create content
- `content:update` - Update content

#### Moderator Role
- `moderation:review` - Review user content
- `moderation:action` - Take moderation actions
- `content:read` - Read content
- `user:read` - Read user profiles

#### Admin Role
- Most permissions except super-admin level
- Includes content management, user management, analytics, moderation

#### Super Admin Role
- All permissions including:
- `system:admin` - Full system administration
- `system:config` - Configure system settings
- `rbac:manage` - Manage roles and permissions

## API Endpoints

All RBAC endpoints are under `/api/v1/admin/rbac/` and require admin privileges.

### Role Management
- `GET /api/v1/admin/rbac/roles` - List all roles
- `POST /api/v1/admin/rbac/roles` - Create a new role
- `GET /api/v1/admin/rbac/roles/:id` - Get role details
- `PUT /api/v1/admin/rbac/roles/:id` - Update a role
- `DELETE /api/v1/admin/rbac/roles/:id` - Delete a role (non-system roles only)

### Permission Management
- `POST /api/v1/admin/rbac/permissions` - Create a new permission

### User-Role Assignment
- `GET /api/v1/admin/rbac/users/:user_id/roles` - Get user's roles
- `POST /api/v1/admin/rbac/users/:user_id/roles` - Assign role to user
- `DELETE /api/v1/admin/rbac/users/:user_id/roles/:role_id` - Remove role from user

### Role-Permission Assignment
- `POST /api/v1/admin/rbac/roles/:role_id/permissions` - Grant permission to role
- `DELETE /api/v1/admin/rbac/roles/:role_id/permissions/:permission_id` - Revoke permission from role

### Permission Checking
- `GET /api/v1/admin/rbac/users/:user_id/permissions` - Get user's permissions
- `POST /api/v1/admin/rbac/check-permission` - Check if user has permission

## Usage Examples

### Create a new role
```bash
curl -X POST http://localhost:8080/api/v1/admin/rbac/roles \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "premium_user",
    "description": "Premium user with extra features",
    "is_system_role": false,
    "is_default": false
  }'
```

### Assign role to user
```bash
curl -X POST http://localhost:8080/api/v1/admin/rbac/users/<USER_ID>/roles \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "role_id": "<ROLE_ID>"
  }'
```

### Check user permission
```bash
curl -X POST http://localhost:8080/api/v1/admin/rbac/check-permission \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "<USER_ID>",
    "permission_name": "content:create"
  }'
```

## Integration with Existing Middleware

The existing middleware has been updated to use RBAC:

### Content Creator Routes
Protected by `ContentCreatorMiddleware` which now checks for `content:manage` permission.

### Admin Routes
Protected by `AdminMiddleware` which now checks for `system:admin` permission.

### Custom Permission Checks
Use `RBACMiddleware` for fine-grained control:
```go
// In router setup
protected.GET("/some/endpoint", middleware.RBACMiddleware("resource:action", "resource_id_param"), handler.Function)
```

## Database Schema

### Key Tables
1. **roles** - Role definitions
2. **permissions** - Permission definitions
3. **permission_categories** - Permission organization
4. **user_roles** - User-role assignments
5. **role_permissions** - Role-permission assignments
6. **role_inheritance** - Role inheritance hierarchy
7. **permission_audit_log** - Audit trail for permission checks
8. **role_audit_log** - Audit trail for role changes

### Important Functions
- `check_user_permission()` - Check if user has permission (including inheritance)
- `get_user_permissions()` - Get all user permissions
- `get_user_roles()` - Get all user roles (including inherited)

## Security Notes

1. **System roles** (`is_system_role = true`) cannot be deleted
2. **Dangerous permissions** (`is_dangerous = true`) require extra scrutiny
3. **Audit logging** - All permission checks and role changes are logged
4. **Role inheritance** - Roles can inherit permissions from other roles
5. **Default role** - New users automatically get the role with `is_default = true`

## Next Steps for Frontend Integration

1. Create admin panel for role/permission management
2. Add permission checks to frontend components
3. Implement role-based UI rendering
4. Add user role management interface
5. Create audit log viewer