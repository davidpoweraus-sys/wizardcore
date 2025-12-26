#!/bin/bash
# Deployment Diagnostics Script for Wizardcore
# Run this script to check the health of all services

echo "=========================================="
echo "WIZARDCORE DEPLOYMENT DIAGNOSTICS"
echo "=========================================="
echo ""

# Check if docker-compose is running
echo "1. Checking Docker Compose services..."
docker compose ps
echo ""

# Check supabase-postgres health
echo "2. Checking Supabase Postgres..."
docker compose exec -T supabase-postgres pg_isready -U supabase_auth_admin -d supabase_auth
echo ""

# Check if auth schema was created
echo "3. Checking auth schema existence..."
docker compose exec -T supabase-postgres psql -U supabase_auth_admin -d supabase_auth -c "\dn auth"
echo ""

# Check supabase-auth logs
echo "4. Last 20 lines of supabase-auth logs..."
docker compose logs --tail=20 supabase-auth
echo ""

# Check supabase-auth health endpoint
echo "5. Testing supabase-auth health endpoint..."
curl -s http://localhost:9999/health || echo "Health endpoint not reachable"
echo ""

# Check main postgres
echo "6. Checking main Wizardcore Postgres..."
docker compose exec -T postgres pg_isready -U wizardcore -d wizardcore
echo ""

# Check Redis
echo "7. Checking Redis..."
docker compose exec -T redis redis-cli -a "$REDIS_PASSWORD" ping 2>/dev/null || echo "Redis check failed"
echo ""

# Check Judge0 Postgres
echo "8. Checking Judge0 Postgres..."
docker compose exec -T judge0-postgres pg_isready -U judge0 -d judge0
echo ""

# Check Judge0 Redis
echo "9. Checking Judge0 Redis..."
docker compose exec -T judge0-redis redis-cli ping
echo ""

# Network connectivity test
echo "10. Network connectivity tests..."
echo "  - Backend can reach supabase-auth:"
docker compose exec -T backend wget -q --spider http://supabase-auth:9999/health && echo "    ✓ Success" || echo "    ✗ Failed"
echo ""

echo "=========================================="
echo "DIAGNOSTICS COMPLETE"
echo "=========================================="
