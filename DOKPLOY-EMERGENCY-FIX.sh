#!/bin/bash
# EMERGENCY FIX FOR DOKPLOY PASSWORD MISMATCH
# Run this on your Dokploy server via SSH

set -e

echo "=========================================="
echo "EMERGENCY DOKPLOY VOLUME CLEANUP"
echo "=========================================="
echo ""
echo "This will:"
echo "1. Stop ALL containers for your app"
echo "2. Remove the supabase_postgres_data volume (deletes auth DB)"
echo "3. Restart containers with correct passwords"
echo ""
echo "⚠️  WARNING: This will delete all authentication data!"
echo "⚠️  Users will need to re-register!"
echo ""
read -p "Type 'DELETE' to continue: " confirmation

if [ "$confirmation" != "DELETE" ]; then
    echo "Aborted."
    exit 0
fi

PROJECT_NAME="offensivewizard-app-dqesjh"

echo ""
echo "Step 1: Stopping all containers..."
docker compose down || docker stop $(docker ps -a -q --filter name=${PROJECT_NAME})

echo ""
echo "Step 2: Removing supabase_postgres_data volume..."
docker volume rm ${PROJECT_NAME}_supabase_postgres_data || \
docker volume rm supabase_postgres_data || \
echo "Volume not found with expected name, trying to find it..."

echo ""
echo "Step 3: Listing all volumes (to find the correct one)..."
docker volume ls | grep postgres

echo ""
echo "If you see a volume above that looks like your supabase postgres,"
echo "run: docker volume rm <volume_name>"
echo ""
echo "Step 4: Restarting deployment..."
echo "Go to Dokploy dashboard and click 'Redeploy'"
echo ""
echo "=========================================="
echo "CLEANUP COMPLETE"
echo "=========================================="
