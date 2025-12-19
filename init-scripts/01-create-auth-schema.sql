-- Create auth schema for Supabase GoTrue
CREATE SCHEMA IF NOT EXISTS auth;

-- Ensure postgres role exists (required by Supabase Auth migrations)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'postgres') THEN
        CREATE ROLE postgres WITH SUPERUSER LOGIN;
    END IF;
END
$$;

-- Grant necessary privileges to supabase_auth_admin
GRANT ALL ON SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA auth TO supabase_auth_admin;