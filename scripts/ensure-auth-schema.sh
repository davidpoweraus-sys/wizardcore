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
    # Run the schema creation script
    psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f /docker-entrypoint-initdb.d/01-create-auth-schema.sql 2>&1
    if [ $? -eq 0 ]; then
        echo "Auth schema created successfully."
    else
        echo "Failed to create auth schema." >&2
        exit 1
    fi
fi

# Additionally ensure the postgres role exists (required by Supabase Auth migrations)
echo "Ensuring postgres role exists..."
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "DO \$\$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'postgres') THEN
        CREATE ROLE postgres WITH SUPERUSER LOGIN;
    END IF;
END
\$\$;" 2>&1 | grep -v "NOTICE"

echo "Schema verification complete."