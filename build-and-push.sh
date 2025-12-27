#!/bin/bash
set -e

# ==============================================
# WIZARDCORE FRONTEND - BUILD AND PUSH TO DOCKER HUB
# ==============================================

echo "🏗️  Building WizardCore Frontend Docker Image..."
echo ""

# Load environment variables
if [ -f .env ]; then
    echo "📋 Loading environment variables from .env..."
    export $(grep -v '^#' .env | xargs)
else
    echo "⚠️  No .env file found. Using .env.example as reference..."
    echo "   Make sure you have the following build args set:"
    echo "   - NEXT_PUBLIC_SUPABASE_URL"
    echo "   - NEXT_PUBLIC_SUPABASE_ANON_KEY"
    echo "   - NEXT_PUBLIC_BACKEND_URL"
    echo "   - NEXT_PUBLIC_JUDGE0_API_URL"
fi

# Docker Hub username
DOCKER_USERNAME="limpet"
IMAGE_NAME="wizardcore-frontend"
TAG="latest"
FULL_IMAGE="${DOCKER_USERNAME}/${IMAGE_NAME}:${TAG}"

echo ""
echo "📦 Image: ${FULL_IMAGE}"
echo ""

# Build the Docker image with build args
echo "🔨 Building Docker image..."
docker build \
  --build-arg NEXT_PUBLIC_SUPABASE_URL="${NEXT_PUBLIC_SUPABASE_URL}" \
  --build-arg NEXT_PUBLIC_SUPABASE_ANON_KEY="${NEXT_PUBLIC_SUPABASE_ANON_KEY}" \
  --build-arg NEXT_PUBLIC_BACKEND_URL="${NEXT_PUBLIC_BACKEND_URL}" \
  --build-arg NEXT_PUBLIC_JUDGE0_API_URL="${NEXT_PUBLIC_JUDGE0_API_URL}" \
  -t "${FULL_IMAGE}" \
  -f Dockerfile \
  .

echo ""
echo "✅ Build complete!"
echo ""

# Push to Docker Hub
echo "🚀 Pushing to Docker Hub..."
docker push "${FULL_IMAGE}"

echo ""
echo "✅ Push complete!"
echo ""
echo "📝 Summary:"
echo "   - Image: ${FULL_IMAGE}"
echo "   - Includes: API proxies for backend and Judge0"
echo "   - CORS fixes: ✅"
echo ""
echo "🎯 Next steps:"
echo "   1. Deploy the updated image to your server"
echo "   2. Make sure these environment variables are set on the server:"
echo "      - BACKEND_URL (internal backend URL, e.g., http://backend:8080)"
echo "      - JUDGE0_URL (internal Judge0 URL, e.g., http://judge0:2358)"
echo "      - NEXT_PUBLIC_SUPABASE_URL (public auth URL)"
echo "      - NEXT_PUBLIC_BACKEND_URL (public backend URL - for reference only)"
echo "      - NEXT_PUBLIC_JUDGE0_API_URL (public Judge0 URL - for reference only)"
echo ""
echo "   3. Restart the frontend container to pick up the new image"
echo ""