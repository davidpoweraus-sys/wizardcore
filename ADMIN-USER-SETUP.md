## Default Admin User Setup Guide

## Overview

Your Wizardcore deployment now automatically creates a default administrator account during initialization. This account has full `super_admin` privileges and is used to bootstrap the system.

## Default Admin Credentials

**Email**: `admin@offensivewizard.com`  
**Password**: Set in Supabase Auth (see setup instructions below)  
**Roles**: `super_admin`, `admin`  
**Permissions**: ALL (full system access)

## CRITICAL: Two-Part Setup Required

Your admin user exists in **TWO separate databases**:

1. **Wizardcore Database** (`postgres`) - Application user profile with roles
2. **Supabase Auth Database** (`supabase-postgres`) - Authentication credentials

### The postgres-init script creates #1 automatically
### You MUST manually create #2 in Supabase

---

## Setup Instructions

### Step 1: Deploy Your Application

When you deploy Wizardcore to Dokploy, the following happens automatically:

```
1. Supabase Auth DB initializes
2. Main Wizardcore DB initializes
3. Backend runs migrations (including 013_rbac_system.up.sql)
   └─> Creates RBAC tables (roles, permissions, user_roles, etc.)
4. postgres-init service starts
   └─> Waits 10 seconds for migrations to complete
   └─> Runs 01-create-default-admin.sql
   └─> Creates admin user in database with UUID: 00000000-0000-0000-0000-000000000001
   └─> Assigns super_admin and admin roles
```

### Step 2: Create Admin User in Supabase Auth

You have **TWO options** to create the Supabase Auth user:

---

#### **Option A: Use Supabase Dashboard (Easiest)**

1. **Access your Supabase Auth dashboard**:
   ```
   https://auth.offensivewizard.com
   ```

2. **Navigate to**: Authentication → Users

3. **Click**: "Invite User" or "Create User"

4. **Enter**:
   - Email: `admin@offensivewizard.com`
   - Password: Choose a strong password
   - User ID (UUID): `00000000-0000-0000-0000-000000000001`

5. **Confirm** the user (or use auto-confirm if enabled)

---

#### **Option B: Use SQL in Supabase Auth Database**

If you don't have dashboard access, connect to the Supabase Auth database and run:

```sql
-- Connect to supabase_auth database
-- psql -h supabase-postgres -U supabase_auth_admin -d supabase_auth

-- Create the admin user in auth.users table
INSERT INTO auth.users (
    id,
    instance_id,
    aud,
    role,
    email,
    encrypted_password,
    email_confirmed_at,
    created_at,
    updated_at,
    confirmation_token,
    recovery_token
) VALUES (
    '00000000-0000-0000-0000-000000000001'::UUID,
    '00000000-0000-0000-0000-000000000000'::UUID,
    'authenticated',
    'authenticated',
    'admin@offensivewizard.com',
    crypt('YOUR_SECURE_PASSWORD', gen_salt('bf')), -- Replace with your password
    NOW(),
    NOW(),
    NOW(),
    '',
    ''
);

-- Create identity record
INSERT INTO auth.identities (
    id,
    user_id,
    identity_data,
    provider,
    last_sign_in_at,
    created_at,
    updated_at
) VALUES (
    gen_random_uuid(),
    '00000000-0000-0000-0000-000000000001'::UUID,
    jsonb_build_object('sub', '00000000-0000-0000-0000-000000000001', 'email', 'admin@offensivewizard.com'),
    'email',
    NOW(),
    NOW(),
    NOW()
);
```

---

### Step 3: Verify Admin User Setup

1. **Navigate to your frontend**:
   ```
   https://app.offensivewizard.com/login
   ```

2. **Log in with**:
   - Email: `admin@offensivewizard.com`
   - Password: (the password you set in Step 2)

3. **Verify permissions**:
   - You should have full access to all features
   - Check user profile shows "super_admin" role
   - Verify you can access admin panels

4. **Check the database**:
   ```sql
   -- Connect to wizardcore database
   SELECT 
       u.email, 
       u.display_name,
       u.role as simple_role,
       string_agg(r.name, ', ') as rbac_roles
   FROM users u
   LEFT JOIN user_roles ur ON u.id = ur.user_id
   LEFT JOIN roles r ON ur.role_id = r.id
   WHERE u.email = 'admin@offensivewizard.com'
   GROUP BY u.id, u.email, u.display_name, u.role;
   ```

   Expected output:
   ```
   email                      | display_name           | simple_role | rbac_roles
   ---------------------------+------------------------+-------------+------------------------
   admin@offensivewizard.com  | System Administrator   | admin       | super_admin, admin
   ```

---

## How User Roles Work

### Simple Role Column (Legacy)
From migration 012, users have a `role` column:
- Values: `student`, `content_creator`, `admin`
- Used for backwards compatibility
- Quick access without joins

### RBAC System (Full System)
From migration 013, users get roles via `user_roles` table:
- Multiple roles per user
- Role inheritance support
- Fine-grained permissions
- Audit logging

### Permission Checking

Use the built-in function to check permissions:

```sql
-- Check if user has specific permission
SELECT check_user_permission(
    '00000000-0000-0000-0000-000000000001'::UUID,  -- admin user ID
    'system:admin'  -- permission name
);
-- Returns: true

-- Get all permissions for a user
SELECT * FROM get_user_permissions('00000000-0000-0000-0000-000000000001'::UUID);
```

---

## Available Permissions

The default admin has ALL of these permissions:

### User Management
- `user:read` - Read user profiles
- `user:update` - Update user profiles
- `user:delete` - Delete users ⚠️
- `user:manage` - Manage all users ⚠️

### Content Management
- `content:read` - Read content
- `content:create` - Create content
- `content:update` - Update content
- `content:delete` - Delete content ⚠️
- `content:manage` - Manage all content ⚠️

### Exercise Management
- `exercise:read` - Read exercises
- `exercise:create` - Create exercises
- `exercise:update` - Update exercises
- `exercise:delete` - Delete exercises ⚠️
- `exercise:submit` - Submit solutions

### Pathway Management
- `pathway:read` - Read pathways
- `pathway:create` - Create pathways
- `pathway:update` - Update pathways
- `pathway:delete` - Delete pathways ⚠️
- `pathway:enroll` - Enroll in pathways

### System Administration
- `system:admin` - Full system administration ⚠️
- `system:config` - Configure settings ⚠️
- `rbac:manage` - Manage roles & permissions ⚠️

### Analytics
- `analytics:view` - View analytics
- `analytics:export` - Export data

### Moderation
- `moderation:review` - Review content
- `moderation:action` - Take moderation actions ⚠️

⚠️ = Dangerous permission (flagged in database)

---

## Creating Additional Admin Users

### Method 1: Via Application UI (Future)
Once logged in as super_admin:
1. Navigate to Users → Create User
2. Assign roles via the RBAC management interface

### Method 2: Via Database

```sql
-- First, create the user in Supabase Auth (see Option B above)
-- Then run this in the wizardcore database:

DO $$
DECLARE
    v_user_id UUID;
    v_admin_role_id UUID;
BEGIN
    -- Create the user
    INSERT INTO users (
        supabase_user_id,
        email,
        display_name,
        role,
        is_active
    ) VALUES (
        'SUPABASE_USER_UUID_HERE'::UUID,
        'newadmin@offensivewizard.com',
        'New Administrator',
        'admin',
        true
    )
    RETURNING id INTO v_user_id;
    
    -- Get admin role ID
    SELECT id INTO v_admin_role_id FROM roles WHERE name = 'admin';
    
    -- Assign admin role
    INSERT INTO user_roles (user_id, role_id)
    VALUES (v_user_id, v_admin_role_id);
    
    RAISE NOTICE 'Created admin user with ID: %', v_user_id;
END $$;
```

---

## Security Recommendations

### 🔒 PRODUCTION SECURITY CHECKLIST

- [ ] **Change default admin email** to something non-obvious
- [ ] **Use a strong password** (20+ characters, mixed case, numbers, symbols)
- [ ] **Enable 2FA** for the admin account (if available)
- [ ] **Create individual admin accounts** for each team member
- [ ] **Disable the default admin** after creating your own
- [ ] **Regularly audit** the `role_audit_log` table
- [ ] **Monitor** the `permission_audit_log` for suspicious activity
- [ ] **Rotate admin passwords** every 90 days
- [ ] **Use service role key** only in backend, never expose to frontend

### Disable Default Admin (After Setup)

```sql
UPDATE users 
SET is_active = false 
WHERE email = 'admin@offensivewizard.com';
```

---

## Troubleshooting

### "User not found" when logging in
**Problem**: Admin user exists in wizardcore DB but not in Supabase Auth  
**Solution**: Complete Step 2 above to create the Supabase Auth user

### "Access denied" despite being admin
**Problem**: User has simple `role='admin'` but no RBAC roles  
**Solution**: Assign roles via `user_roles` table:
```sql
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id 
FROM users u, roles r 
WHERE u.email = 'your@email.com' 
AND r.name = 'super_admin';
```

### Admin user not created during deployment
**Problem**: Migrations might not have completed before postgres-init ran  
**Solution**: Manually run the script:
```bash
docker compose exec postgres psql -U wizardcore -d wizardcore -f /scripts/01-create-default-admin.sql
```

### Want to reset everything
```bash
# Stop containers
docker compose down

# Remove volumes (WARNING: Deletes all data!)
docker volume rm offensivewizard-app-dqesjh_postgres_data
docker volume rm offensivewizard-app-dqesjh_supabase_postgres_data

# Redeploy
docker compose up -d
```

---

## Files Reference

```
wizardcore/
├── postgres-init-scripts/
│   └── 01-create-default-admin.sql    # Creates default admin user
│
└── wizardcore-backend/
    └── internal/database/migrations/
        ├── 013_rbac_system.up.sql     # Creates RBAC tables & permissions
        └── 013_rbac_system.down.sql   # Rollback script
```

---

## Next Steps

1. ✅ Deploy your application
2. ✅ Create admin user in Supabase Auth
3. ✅ Log in and verify access
4. 🔒 Change default credentials
5. 👥 Create individual admin accounts for your team
6. 🚀 Start using your platform!
