-- ============================================
-- CREATE DEFAULT SYSTEM ADMINISTRATOR
-- ============================================
-- This script creates a default admin user that can be used to bootstrap the system
-- The user will be linked to Supabase Auth once they log in for the first time
-- 
-- IMPORTANT: Change these credentials in production!
-- Default credentials:
--   Email: admin@offensivewizard.com
--   The actual password is managed in Supabase Auth - you must create this user there first
-- ============================================

-- Wait for RBAC tables to exist (in case migrations haven't run yet)
DO $$
BEGIN
    -- Check if roles table exists
    IF NOT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'roles') THEN
        RAISE NOTICE 'RBAC tables not found - migrations may not have run yet. Skipping admin user creation.';
        RETURN;
    END IF;
    
    -- Check if required roles exist
    IF NOT EXISTS (SELECT 1 FROM roles WHERE name = 'super_admin') THEN
        RAISE NOTICE 'super_admin role not found - migrations may not have run yet. Skipping admin user creation.';
        RETURN;
    END IF;
    
    RAISE NOTICE 'RBAC tables found. Proceeding with admin user setup...';
END $$;

-- Function to create default admin user
-- This is idempotent - safe to run multiple times
CREATE OR REPLACE FUNCTION create_default_admin_user()
RETURNS VOID AS $$
DECLARE
    v_user_id UUID;
    v_super_admin_role_id UUID;
    v_admin_role_id UUID;
    v_supabase_user_id UUID;
BEGIN
    -- Generate a consistent UUID for the default admin (using a known seed)
    -- This ensures the same UUID is used if the script runs multiple times
    v_supabase_user_id := '00000000-0000-0000-0000-000000000001'::UUID;
    
    -- Get role IDs
    SELECT id INTO v_super_admin_role_id FROM roles WHERE name = 'super_admin';
    SELECT id INTO v_admin_role_id FROM roles WHERE name = 'admin';
    
    IF v_super_admin_role_id IS NULL THEN
        RAISE NOTICE 'super_admin role not found. Cannot create default admin user.';
        RETURN;
    END IF;
    
    -- Check if admin user already exists
    SELECT id INTO v_user_id FROM users WHERE email = 'admin@offensivewizard.com';
    
    IF v_user_id IS NOT NULL THEN
        RAISE NOTICE 'Admin user already exists with ID: %', v_user_id;
        
        -- Ensure the user has super_admin role
        INSERT INTO user_roles (user_id, role_id, assigned_at)
        VALUES (v_user_id, v_super_admin_role_id, CURRENT_TIMESTAMP)
        ON CONFLICT (user_id, role_id) DO NOTHING;
        
        RAISE NOTICE 'Verified super_admin role for existing admin user.';
        RETURN;
    END IF;
    
    -- Create the default admin user
    INSERT INTO users (
        id,
        supabase_user_id,
        email,
        display_name,
        role, -- Simple role column from migration 012
        is_active,
        total_xp,
        created_at,
        updated_at
    ) VALUES (
        gen_random_uuid(),
        v_supabase_user_id,
        'admin@offensivewizard.com',
        'System Administrator',
        'admin',
        true,
        0,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    )
    RETURNING id INTO v_user_id;
    
    RAISE NOTICE 'Created default admin user with ID: %', v_user_id;
    
    -- Assign super_admin role
    INSERT INTO user_roles (user_id, role_id, assigned_at)
    VALUES (v_user_id, v_super_admin_role_id, CURRENT_TIMESTAMP);
    
    RAISE NOTICE 'Assigned super_admin role to default admin user.';
    
    -- Also assign regular admin role for backwards compatibility
    IF v_admin_role_id IS NOT NULL THEN
        INSERT INTO user_roles (user_id, role_id, assigned_at)
        VALUES (v_user_id, v_admin_role_id, CURRENT_TIMESTAMP)
        ON CONFLICT (user_id, role_id) DO NOTHING;
        
        RAISE NOTICE 'Assigned admin role to default admin user.';
    END IF;
    
    -- Log the role assignment
    INSERT INTO role_audit_log (
        user_id,
        target_user_id,
        role_id,
        action,
        reason,
        performed_by,
        created_at
    ) VALUES (
        v_user_id,
        v_user_id,
        v_super_admin_role_id,
        'assigned',
        'Initial system setup - default admin user creation',
        v_user_id,
        CURRENT_TIMESTAMP
    );
    
    RAISE NOTICE '================================================';
    RAISE NOTICE 'DEFAULT ADMIN USER CREATED SUCCESSFULLY';
    RAISE NOTICE '================================================';
    RAISE NOTICE 'Email: admin@offensivewizard.com';
    RAISE NOTICE 'Supabase UUID: %', v_supabase_user_id;
    RAISE NOTICE 'Database UUID: %', v_user_id;
    RAISE NOTICE 'Roles: super_admin, admin';
    RAISE NOTICE '';
    RAISE NOTICE 'IMPORTANT: You must create this user in Supabase Auth with the same email!';
    RAISE NOTICE 'Use the Supabase dashboard or API to create the auth user.';
    RAISE NOTICE '================================================';
    
END;
$$ LANGUAGE plpgsql;

-- Execute the function (only if tables exist)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'roles') THEN
        PERFORM create_default_admin_user();
    ELSE
        RAISE NOTICE 'Skipping admin user creation - RBAC tables do not exist yet.';
        RAISE NOTICE 'This script will run again after migrations complete.';
    END IF;
END $$;

-- Clean up the function after use
DROP FUNCTION IF EXISTS create_default_admin_user();
