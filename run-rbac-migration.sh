#!/bin/bash
set -e

echo "=== RBAC Database Migration ==="

# Check if PostgreSQL is running
if ! docker ps | grep -q "wizardcore-postgres\|postgres"; then
    echo "Starting PostgreSQL container..."
    cd /home/glbsi/Workbench/wizardcore
    docker-compose -f docker-compose.prod.yml up -d postgres 2>/dev/null || true
    
    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    sleep 5
fi

# Get database connection details from environment or use defaults
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-wizardcore}
DB_USER=${DB_USER:-wizardcore}
DB_PASSWORD=${DB_PASSWORD:-/8I+qc0afMvSHmL987Vk6A==}

# URL encode the password (replace = with %3D)
ENCODED_PASSWORD=$(echo "$DB_PASSWORD" | sed 's/=/\%3D/g')

# Construct connection string
CONNECTION_STRING="postgresql://${DB_USER}:${ENCODED_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Database: ${DB_NAME}"
echo "User: ${DB_USER}"
echo "Host: ${DB_HOST}:${DB_PORT}"

# Check if we can connect
if ! PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT 1;" >/dev/null 2>&1; then
    echo "ERROR: Cannot connect to PostgreSQL database"
    echo "Trying to connect via Docker exec..."
    
    # Try connecting via Docker exec
    CONTAINER_ID=$(docker ps -q --filter "name=wizardcore-postgres" --filter "name=postgres")
    if [ -n "$CONTAINER_ID" ]; then
        echo "Found PostgreSQL container: $CONTAINER_ID"
        docker exec "$CONTAINER_ID" psql -U "${DB_USER}" -d "${DB_NAME}" -c "SELECT 1;" >/dev/null 2>&1
        if [ $? -eq 0 ]; then
            echo "Connected via Docker exec"
            # Run migration via Docker exec
            echo "Running RBAC migration..."
            docker exec -i "$CONTAINER_ID" psql -U "${DB_USER}" -d "${DB_NAME}" < /home/glbsi/Workbench/wizardcore/rbac-schema.sql
            echo "Migration completed successfully!"
            exit 0
        fi
    fi
    
    echo "Please ensure PostgreSQL is running and accessible"
    exit 1
fi

# Run the migration
echo "Running RBAC migration..."
PGPASSWORD="${DB_PASSWORD}" psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USER}" -d "${DB_NAME}" -f /home/glbsi/Workbench/wizardcore/rbac-schema.sql

echo "Migration completed successfully!"