# Scalable Role-Based Access Control (RBAC) Implementation Guide

## Overview

This document outlines the comprehensive, scalable RBAC system designed for the Offensive Wizard application. The system supports dynamic role creation, permission management, role inheritance, and audit logging.

## Architecture

### Core Components

1. **Roles** - Named groups of permissions (e.g., "admin", "content_creator", "moderator")
2. **Permissions** - Fine-grained access controls (e.g., "exercise:create", "user:delete")
3. **Role-Permission Mapping** - Many-to-many relationships
4. **User-Role Assignment** - Users can have multiple roles
5. **Role Inheritance** - Roles can inherit permissions from other roles
6. **Audit Logging** - Comprehensive logging of all permission checks and role changes

### Database Schema

See `rbac-schema.sql` for complete schema with:
- Tables for roles, permissions, categories
- Relationship tables for role-permission and user-role mappings
- Inheritance support via `role_inheritance` table
- Audit logging tables
- PostgreSQL functions for permission checking

## Implementation Steps

### Phase 1: Database Migration

1. **Run the RBAC schema migration**:
   ```bash
   psql -d your_database -f rbac-schema.sql
   ```

2. **Update existing users table**:
   ```sql
   ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT true;
   ```

### Phase 2: Backend Integration

1. **Add RBAC services to dependency injection** (in `router/router.go`):
   ```go
   // Initialize RBAC repository
   rbacRepo := repositories.NewRBACRepository(db, logger)
   
   // Initialize RBAC service
   rbacService := services.NewRBACService(rbacRepo, userRepo, logger)
   
   // Initialize RBAC handler
   rbacHandler := handlers.NewRBACHandler(rbacService, logger)
   ```

2. **Add RBAC middleware to routes**:
   ```go
   // Add RBAC service to Gin context
   r.Use(func(c *gin.Context) {
       c.Set("rbac_service", rbacService)
       c.Next()
   })
   ```

3. **Protect routes with RBAC middleware**:
   ```go
   // Example: Admin-only route
   adminGroup := protected.Group("/admin")
   adminGroup.Use(middleware.RBACMiddleware("system:manage", ""))
   {
       adminGroup.POST("/users", adminHandler.CreateUser)
       adminGroup.DELETE("/users/:id", adminHandler.DeleteUser)
   }
   
   // Example: Content creator route
   creatorGroup := protected.Group("/creator")
   creatorGroup.Use(middleware.RequireAnyPermission([]string{
       "pathway:create",
       "module:create",
       "exercise:create",
   }, ""))
   {
       creatorGroup.POST("/pathways", creatorHandler.CreatePathway)
       creatorGroup.POST("/exercises", creatorHandler.CreateExercise)
   }
   ```

### Phase 3: Frontend Integration

1. **Create RBAC context/provider**:
   ```typescript
   // lib/rbac/context.tsx
   import { createContext, useContext, useEffect, useState } from 'react'
   import { getCurrentUserPermissions } from '@/lib/rbac/api'
   
   interface RBACContextType {
     permissions: string[]
     hasPermission: (permission: string) => boolean
     hasAnyPermission: (permissions: string[]) => boolean
     hasAllPermissions: (permissions: string[]) => boolean
     isLoading: boolean
   }
   
   export const RBACContext = createContext<RBACContextType>({
     permissions: [],
     hasPermission: () => false,
     hasAnyPermission: () => false,
     hasAllPermissions: () => false,
     isLoading: true,
   })
   
   export function RBACProvider({ children }: { children: React.ReactNode }) {
     const [permissions, setPermissions] = useState<string[]>([])
     const [isLoading, setIsLoading] = useState(true)
   
     useEffect(() => {
       loadPermissions()
     }, [])
   
     const loadPermissions = async () => {
       try {
         const userPermissions = await getCurrentUserPermissions()
         setPermissions(userPermissions)
       } catch (error) {
         console.error('Failed to load permissions:', error)
       } finally {
         setIsLoading(false)
       }
     }
   
     const hasPermission = (permission: string) => {
       return permissions.includes(permission)
     }
   
     const hasAnyPermission = (requiredPermissions: string[]) => {
       return requiredPermissions.some(perm => permissions.includes(perm))
     }
   
     const hasAllPermissions = (requiredPermissions: string[]) => {
       return requiredPermissions.every(perm => permissions.includes(perm))
     }
   
     return (
       <RBACContext.Provider value={{
         permissions,
         hasPermission,
         hasAnyPermission,
         hasAllPermissions,
         isLoading,
       }}>
         {children}
       </RBACContext.Provider>
     )
   }
   
   export const useRBAC = () => useContext(RBACContext)
   ```

2. **Create permission-based UI components**:
   ```typescript
   // components/common/ProtectedComponent.tsx
   import { useRBAC } from '@/lib/rbac/context'
   
   interface ProtectedComponentProps {
     children: React.ReactNode
     requiredPermission?: string
     requiredPermissions?: string[]
     requireAll?: boolean
     fallback?: React.ReactNode
   }
   
   export function ProtectedComponent({
     children,
     requiredPermission,
     requiredPermissions,
     requireAll = false,
     fallback = null,
   }: ProtectedComponentProps) {
     const { hasPermission, hasAnyPermission, hasAllPermissions, isLoading } = useRBAC()
   
     if (isLoading) {
       return <div>Loading permissions...</div>
     }
   
     let hasAccess = false
   
     if (requiredPermission) {
       hasAccess = hasPermission(requiredPermission)
     } else if (requiredPermissions) {
       if (requireAll) {
         hasAccess = hasAllPermissions(requiredPermissions)
       } else {
         hasAccess = hasAnyPermission(requiredPermissions)
       }
     } else {
       hasAccess = true // No permission requirement
     }
   
     return hasAccess ? <>{children}</> : <>{fallback}</>
   }
   ```

3. **Use in components**:
   ```tsx
   // Example usage
   function AdminPanel() {
     return (
       <ProtectedComponent requiredPermission="system:manage">
         <AdminDashboard />
       </ProtectedComponent>
     )
   }
   
   function ContentCreatorTools() {
     return (
       <ProtectedComponent requiredPermissions={['pathway:create', 'exercise:create']}>
         <ContentCreatorDashboard />
       </ProtectedComponent>
     )
   }
   ```

## Default Roles and Permissions

### System Roles (Cannot be deleted)

1. **super_admin** - Full system access
   - `system:manage`
   - `user:manage`
   - `role:manage`
   - `permission:manage`

2. **admin** - Administrative access
   - `user:manage`
   - `pathway:manage`
   - `module:manage`
   - `exercise:manage`
   - `submission:read`

3. **user** - Basic user access (default)
   - `user:read` (own profile)
   - `user:update` (own profile)
   - `pathway:read`
   - `module:read`
   - `exercise:read`
   - `submission:create`
   - `submission:read` (own submissions)

### Custom Roles (Can be created dynamically)

4. **content_creator** - Create and manage content
   - `pathway:create`
   - `pathway:update`
   - `module:create`
   - `module:update`
   - `exercise:create`
   - `exercise:update`

5. **moderator** - Moderate user content
   - `user:read`
   - `submission:read`
   - `submission:update`
   - `pathway:update`

## Permission Naming Convention

Permissions follow the format: `resource:action`

### Common Resources
- `user` - User accounts
- `role` - RBAC roles
- `permission` - RBAC permissions
- `pathway` - Learning pathways/courses
- `module` - Course modules
- `exercise` - Code exercises
- `submission` - Exercise submissions
- `system` - System configuration

### Common Actions
- `create` - Create new resource
- `read` - View resource
- `update` - Modify resource
- `delete` - Remove resource
- `manage` - Full control (create, read, update, delete)

## API Endpoints

### RBAC Management
- `POST /api/v1/rbac/roles` - Create role
- `GET /api/v1/rbac/roles` - List roles
- `GET /api/v1/rbac/roles/:id` - Get role with permissions
- `PUT /api/v1/rbac/roles/:id` - Update role
- `DELETE /api/v1/rbac/roles/:id` - Delete role

### Permission Management
- `POST /api/v1/rbac/permissions` - Create permission
- `GET /api/v1/rbac/permissions` - List permissions

### Role-Permission Management
- `POST /api/v1/rbac/roles/permissions/grant` - Grant permission to role
- `POST /api/v1/rbac/roles/permissions/revoke` - Revoke permission from role

### User-Role Management
- `POST /api/v1/rbac/users/roles/assign` - Assign role to user
- `POST /api/v1/rbac/users/roles/remove` - Remove role from user
- `GET /api/v1/rbac/users/:user_id/roles` - Get user roles
- `GET /api/v1/rbac/users/:user_id/permissions` - Get user permissions

### Permission Checking
- `POST /api/v1/rbac/permissions/check` - Check if user has permission

### Role Inheritance
- `POST /api/v1/rbac/roles/inheritance/add` - Add role inheritance
- `POST /api/v1/rbac/roles/inheritance/remove` - Remove role inheritance

## Middleware Usage Examples

### Basic Permission Check
```go
// Require single permission
router.GET("/admin/users", 
    middleware.RBACMiddleware("user:read", ""),
    userHandler.ListUsers)

// Require permission with resource ID
router.GET("/exercises/:id/submissions",
    middleware.RBACMiddleware("submission:read", "id"),
    submissionHandler.GetExerciseSubmissions)
```

### Multiple Permission Options
```go
// Require any of these permissions
router.POST("/content",
    middleware.RequireAnyPermission([]string{
        "pathway:create",
        "module:create",
        "exercise:create",
    }, ""),
    contentHandler.CreateContent)

// Require all of these permissions
router.GET("/admin/dashboard",
    middleware.RequireAllPermissions([]string{
        "user:read",
        "submission:read",
        "system:read",
    }, ""),
    adminHandler.GetDashboard)
```

### Role-Based Access
```go
// Require specific role
router.GET("/super-admin/audit",
    middleware.RequireRole("super_admin"),
    adminHandler.GetAuditLogs)

// Require any of these roles
router.GET("/moderation/tools",
    middleware.RequireAnyRole([]string{"admin", "moderator"}),
    moderationHandler.GetTools)
```

### Resource Ownership
```go
// Check if user owns the resource or has admin permission
router.PUT("/users/:id/profile",
    middleware.ResourceOwnerMiddleware("user", "id", getUserOwnerID),
    userHandler.UpdateProfile)

func getUserOwnerID(c *gin.Context, userID string) (uuid.UUID, error) {
    // Parse user ID and fetch from database
    id, err := uuid.Parse(userID)
    if err != nil {
        return uuid.Nil, err
    }
    
    user, err := userRepo.FindByID(id)
    if err != nil || user == nil {
        return uuid.Nil, fmt.Errorf("user not found")
    }
    
    return user.ID, nil
}
```

## Audit Logging

The system automatically logs:
- All permission checks (granted/denied)
- All role assignments and removals
- IP addresses and user agents
- Timestamps and metadata

View audit logs via:
```sql
SELECT * FROM permission_audit_log ORDER BY created_at DESC LIMIT 100;
SELECT * FROM role_audit_log ORDER BY created_at DESC LIMIT 100;
```

## Security Considerations

1. **Least Privilege**: Assign minimum necessary permissions
2. **Regular Audits**: Review permission assignments quarterly
3. **Dangerous Permissions**: Mark high-risk permissions with `is_dangerous = true`
4. **Temporary Roles**: Use `expires_at` for temporary access
5. **Input Validation**: Validate all role and permission names
6. **Rate Limiting**: Apply to RBAC management endpoints

## Migration Strategy

1. **Backup existing data** before running migrations
2. **Run in staging** first to test compatibility
3. **Assign default 'user' role** to all existing users
4. **Gradually migrate endpoints** to use RBAC middleware
5. **Monitor audit logs** for permission denials
6. **Train administrators** on RBAC management

## Testing

### Unit Tests
- Permission checking logic
- Role inheritance calculations
- Input validation

### Integration Tests
- API endpoints with various permission scenarios
- Database operations
- Middleware behavior

### End-to-End Tests
- Complete user workflows with permission checks
- Admin role management flows
- Permission denial scenarios

## Monitoring and Maintenance

1. **Monitor permission denials** for potential security issues
2. **Review role assignments** regularly
3. **Clean up expired role assignments**
4. **Archive old audit logs** (keep 90 days minimum)
5. **Update default permissions** as new features are added

## Troubleshooting

### Common Issues

1. **Permission denied but should be granted**
   - Check role assignments
   - Verify role inheritance
   - Check permission names match exactly

2. **Circular inheritance detected**
   - Review role inheritance chain
   - Remove circular references

3. **Cannot delete system role**
   - System roles are protected
   - Create custom roles instead

4. **Performance issues with permission checks**
   - Ensure proper indexes exist
   - Consider caching user permissions
   - Review recursive queries

## Next Steps

1. Implement the database schema
2. Integrate RBAC services into existing router
3. Migrate existing endpoints to use RBAC middleware
4. Create admin UI for role and permission management
5. Implement frontend RBAC context and components
6. Conduct security review and penetration testing

## Support

For questions or issues:
1. Check audit logs for permission denials
2. Review role and permission assignments
3. Test with the permission check endpoint
4. Consult this implementation guide