#!/bin/bash

# Build and push to Docker Hub with CORRECT values for Coolify
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

# Build frontend with CORRECT environment variables for Coolify
# Using proxy URL: https://offensivewizard.com/supabase-proxy
docker build \
  -t "${IMAGE_PREFIX}-frontend:${TAG}" \
  -f Dockerfile.nextjs \
  --build-arg NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/supabase-proxy \
  --build-arg NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps= \
  --build-arg NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com \
  --build-arg NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api \
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
echo ""
echo "📋 IMPORTANT: These images include the proxy fix for whitespace/JSON parsing errors!"
echo "   - Proxy no longer strips /auth/v1 prefix"
echo "   - Added UTF-8 BOM removal"
echo "   - Enhanced error handling for non-JSON responses"
echo "   - Fixed URL construction to prevent double slashes"