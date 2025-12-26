#!/bin/bash
# Script to clean up and reset Supabase Postgres volumes on Dokploy
# Run this if you're getting password authentication errors

set -e

PROJECT_NAME="offensivewizard-app-dqesjh"

echo "=========================================="
echo "DOKPLOY VOLUME CLEANUP SCRIPT"
echo "=========================================="
echo ""
echo "This script will:"
echo "1. Stop the supabase-auth and supabase-postgres containers"
echo "2. Remove the supabase_postgres_data volume (this will delete the database!)"
echo "3. Restart the containers to reinitialize with correct credentials"
echo ""
read -p "Are you sure you want to continue? (yes/no): " confirmation

if [ "$confirmation" != "yes" ]; then
    echo "Aborted."
    exit 0
fi

echo ""
echo "Step 1: Stopping supabase-auth container..."
docker stop ${PROJECT_NAME}-supabase-auth-1 2>/dev/null || echo "Container not running"

echo "Step 2: Stopping supabase-init container..."
docker stop ${PROJECT_NAME}-supabase-init-1 2>/dev/null || echo "Container not running"

echo "Step 3: Stopping supabase-postgres container..."
docker stop ${PROJECT_NAME}-supabase-postgres-1 2>/dev/null || echo "Container not running"

echo "Step 4: Removing supabase-auth container..."
docker rm ${PROJECT_NAME}-supabase-auth-1 2>/dev/null || echo "Container not found"

echo "Step 5: Removing supabase-init container..."
docker rm ${PROJECT_NAME}-supabase-init-1 2>/dev/null || echo "Container not found"

echo "Step 6: Removing supabase-postgres container..."
docker rm ${PROJECT_NAME}-supabase-postgres-1 2>/dev/null || echo "Container not found"

echo "Step 7: Removing supabase_postgres_data volume..."
docker volume rm ${PROJECT_NAME}_supabase_postgres_data 2>/dev/null || echo "Volume not found or still in use"

echo ""
echo "=========================================="
echo "Cleanup complete!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Redeploy your application in Dokploy"
echo "2. The supabase-postgres container will initialize with fresh data"
echo "3. The init scripts will create the user and schema correctly"
echo ""
