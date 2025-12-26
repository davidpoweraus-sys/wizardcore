-- Ensure supabase_auth_admin user exists with correct password
-- This runs BEFORE the main init script
-- Note: This is idempotent and safe to run multiple times

DO $$
BEGIN
    -- Check if the role exists
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'supabase_auth_admin') THEN
        -- Create the role if it doesn't exist
        RAISE NOTICE 'Creating supabase_auth_admin role...';
        CREATE ROLE supabase_auth_admin WITH LOGIN PASSWORD '0yvhLSetDKV4BlFOH6YeM5LCBe2jmV2B';
    ELSE
        -- Update the password if the role already exists
        RAISE NOTICE 'Updating supabase_auth_admin password...';
        ALTER ROLE supabase_auth_admin WITH LOGIN PASSWORD '0yvhLSetDKV4BlFOH6YeM5LCBe2jmV2B';
    END IF;
    
    -- Ensure the role has all necessary privileges
    ALTER ROLE supabase_auth_admin WITH SUPERUSER CREATEDB CREATEROLE;
END
$$;

-- Grant database connection privileges
GRANT ALL PRIVILEGES ON DATABASE supabase_auth TO supabase_auth_admin;
