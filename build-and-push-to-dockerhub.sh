#!/bin/bash

# Build and push to Docker Hub
set -e

DOCKER_USERNAME="limpet"
IMAGE_PREFIX="${DOCKER_USERNAME}/wizardcore"
TAG="latest"

echo "🔧 Building Backend Docker Image..."
cd wizardcore-backend
docker build -t "${IMAGE_PREFIX}-backend:${TAG}" .
cd ..

echo ""
echo "🔧 Building Frontend Docker Image..."
echo "This may take 5-10 minutes and requires 8GB+ RAM"
echo ""

# Build frontend with production environment variables
# CRITICAL: Use /api/auth proxy endpoint for Supabase URL
docker build \
  -t "${IMAGE_PREFIX}-frontend:${TAG}" \
  -f Dockerfile.nextjs \
  --build-arg NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth \
  --build-arg NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ \
  --build-arg NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com \
  --build-arg NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com \
  --build-arg NEXT_PUBLIC_SITE_URL=https://app.offensivewizard.com \
  .

echo ""
echo "✅ Builds complete!"
echo ""

# Check if logged in to Docker Hub
if ! docker info 2>/dev/null | grep -q "Username: ${DOCKER_USERNAME}"; then
    echo "🔐 Logging in to Docker Hub..."
    echo "Antony£)" | docker login -u "limpet" --password-stdin
fi

echo ""
echo "📤 Pushing Backend to Docker Hub..."
docker push "${IMAGE_PREFIX}-backend:${TAG}"

echo ""
echo "📤 Pushing Frontend to Docker Hub..."
docker push "${IMAGE_PREFIX}-frontend:${TAG}"

echo ""
echo "✅ Success!"
echo ""
echo "Images pushed to Docker Hub:"
echo "1. ${IMAGE_PREFIX}-backend:${TAG}"
echo "2. ${IMAGE_PREFIX}-frontend:${TAG}"
echo ""
echo "To use these images, update your docker-compose.yml:"
echo "  frontend:"
echo "    image: ${IMAGE_PREFIX}-frontend:${TAG}"
echo "  backend:"
echo "    image: ${IMAGE_PREFIX}-backend:${TAG}"