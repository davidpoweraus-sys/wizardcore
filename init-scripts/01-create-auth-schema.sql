-- Create auth schema for Supabase GoTrue
CREATE SCHEMA IF NOT EXISTS auth;

-- Grant necessary privileges to supabase_auth_admin
GRANT ALL ON SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA auth TO supabase_auth_admin;