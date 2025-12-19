#!/bin/bash
set -e

# Ensure auth schema exists in Supabase Auth database
# Usage: ./scripts/ensure-auth-schema.sh [DB_HOST] [DB_PORT] [DB_USER] [DB_PASSWORD] [DB_NAME]
# Defaults match docker-compose.prod.yml

DB_HOST="${1:-supabase-postgres}"
DB_PORT="${2:-5432}"
DB_USER="${3:-supabase_auth_admin}"
DB_PASSWORD="${4:-password}"
DB_NAME="${5:-supabase_auth}"

export PGPASSWORD="$DB_PASSWORD"

echo "Checking auth schema in database $DB_NAME on $DB_HOST:$DB_PORT..."

# Check if auth schema exists
schema_exists=$(psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -t -c "SELECT 1 FROM information_schema.schemata WHERE schema_name = 'auth';" 2>/dev/null || echo "0")

if [ "$schema_exists" -eq 1 ]; then
    echo "Auth schema already exists."
else
    echo "Auth schema missing. Creating..."
    # Create schema and necessary objects inline
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'EOF'
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

-- Create enum types required by Supabase Auth migrations
-- factor_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'factor_type' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.factor_type AS ENUM ('totp', 'webauthn');
        ALTER TYPE auth.factor_type ADD VALUE 'phone';
    END IF;
END
$$;

-- code_challenge_method
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'code_challenge_method' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.code_challenge_method AS ENUM ('S256', 'plain');
    END IF;
END
$$;

-- oauth_authorization_status
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_authorization_status' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.oauth_authorization_status AS ENUM ('pending', 'approved', 'denied', 'expired');
    END IF;
END
$$;

-- oauth_response_type
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_response_type' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.oauth_response_type AS ENUM ('code');
    END IF;
END
$$;

-- Grant necessary privileges to supabase_auth_admin
GRANT ALL ON SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA auth TO supabase_auth_admin;
GRANT ALL PRIVILEGES ON ALL FUNCTIONS IN SCHEMA auth TO supabase_auth_admin;
EOF
    if [ $? -eq 0 ]; then
        echo "Auth schema created successfully."
    else
        echo "Failed to create auth schema." >&2
        exit 1
    fi
fi

# Additionally ensure the postgres role exists (already done above, but double-check)
echo "Ensuring postgres role exists..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "DO \$\$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'postgres') THEN
        CREATE ROLE postgres WITH SUPERUSER LOGIN;
    END IF;
END
\$\$;" 2>&1 | grep -v "NOTICE"

# Ensure enum types exist (they should have been created, but if schema existed without them, add them)
echo "Ensuring enum types exist in auth schema..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" <<'EOF'
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'factor_type' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.factor_type AS ENUM ('totp', 'webauthn');
        ALTER TYPE auth.factor_type ADD VALUE 'phone';
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'code_challenge_method' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.code_challenge_method AS ENUM ('S256', 'plain');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_authorization_status' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.oauth_authorization_status AS ENUM ('pending', 'approved', 'denied', 'expired');
    END IF;
END
$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'oauth_response_type' AND typnamespace = (SELECT oid FROM pg_namespace WHERE nspname = 'auth')) THEN
        CREATE TYPE auth.oauth_response_type AS ENUM ('code');
    END IF;
END
$$;
EOF

echo "Schema verification complete."