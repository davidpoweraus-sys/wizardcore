#!/bin/bash

set -e

echo "=========================================="
echo "Starting WizardCore Full Stack Application"
echo "=========================================="

# Check for Docker and Docker Compose
if ! command -v docker &> /dev/null; then
    echo "Docker not found. Please install Docker."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "docker-compose not found. Please install Docker Compose."
    exit 1
fi

# Ensure we are in the project root
cd "$(dirname "$0")"

# Start all services defined in docker-compose.local.yml
echo "Starting Docker services (PostgreSQL, Redis, Judge0, Supabase Auth, Backend)..."
docker-compose -f docker-compose.local.yml up -d

echo "Waiting for services to become healthy (30 seconds)..."
sleep 30

# Check if PostgreSQL is ready
echo "Checking PostgreSQL readiness..."
until docker-compose -f docker-compose.local.yml exec -T postgres pg_isready -U wizardcore; do
    echo "PostgreSQL is not ready yet. Waiting 5 seconds..."
    sleep 5
done

echo "PostgreSQL is ready."

# Run database migrations
echo "Running database migrations..."
cd wizardcore-backend
go run cmd/migrate/main.go
cd ..

# Optionally seed the database (uncomment if needed)
# echo "Seeding database with sample data..."
# cd wizardcore-backend
# go run cmd/seed/main.go
# cd ..

# Start frontend Next.js dev server in the background
echo "Starting Next.js frontend development server..."
cd "$(dirname "$0")"
npm install > /dev/null 2>&1 || echo "npm install failed, but continuing..."
npm run dev &
FRONTEND_PID=$!

echo "Frontend started with PID $FRONTEND_PID"

# Display service URLs
echo ""
echo "=========================================="
echo "Services are now running:"
echo "  - Frontend:      http://localhost:3000"
echo "  - Backend API:   http://localhost:8080"
echo "  - Judge0:        http://localhost:2358"
echo "  - Supabase Auth: http://localhost:9999"
echo "  - PostgreSQL:    localhost:5432"
echo "  - Redis:         localhost:6379"
echo ""
echo "To view logs, run:"
echo "  docker-compose -f docker-compose.local.yml logs -f"
echo ""
echo "To stop all services, run:"
echo "  docker-compose -f docker-compose.local.yml down"
echo "  kill $FRONTEND_PID"
echo "=========================================="

# Wait for user interrupt
trap "echo 'Stopping services...'; docker-compose -f docker-compose.local.yml down; kill $FRONTEND_PID 2>/dev/null; exit 0" INT TERM

echo "Press Ctrl+C to stop all services."
wait $FRONTEND_PID