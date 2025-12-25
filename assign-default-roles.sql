-- Assign default 'user' role to all existing users who don't have any role
-- This should be run after the RBAC migration

DO $$
DECLARE
    default_role_id UUID;
    user_record RECORD;
BEGIN
    -- Get the default 'user' role ID
    SELECT id INTO default_role_id FROM roles WHERE is_default = true AND name = 'user';
    
    IF default_role_id IS NULL THEN
        RAISE NOTICE 'Default user role not found';
        RETURN;
    END IF;
    
    RAISE NOTICE 'Default user role ID: %', default_role_id;
    
    -- Assign default role to all users who don't have any role
    FOR user_record IN 
        SELECT u.id 
        FROM users u
        WHERE NOT EXISTS (
            SELECT 1 FROM user_roles ur WHERE ur.user_id = u.id
        )
    LOOP
        BEGIN
            INSERT INTO user_roles (user_id, role_id, assigned_at, assigned_by)
            VALUES (user_record.id, default_role_id, NOW(), NULL);
            
            RAISE NOTICE 'Assigned default role to user: %', user_record.id;
        EXCEPTION
            WHEN unique_violation THEN
                RAISE NOTICE 'User % already has role', user_record.id;
            WHEN OTHERS THEN
                RAISE NOTICE 'Error assigning role to user %: %', user_record.id, SQLERRM;
        END;
    END LOOP;
    
    RAISE NOTICE 'Default role assignment completed';
END $$;