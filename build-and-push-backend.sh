#!/bin/bash
set -e

# ==============================================
# WIZARDCORE BACKEND - BUILD AND PUSH TO DOCKER HUB
# ==============================================

echo "🏗️  Building WizardCore Backend Docker Image..."
echo ""

# Docker Hub username
DOCKER_USERNAME="limpet"
IMAGE_NAME="wizardcore-backend"
TAG="latest"
FULL_IMAGE="${DOCKER_USERNAME}/${IMAGE_NAME}:${TAG}"

echo ""
echo "📦 Image: ${FULL_IMAGE}"
echo ""

# Build the Docker image with cache busting
echo "🔨 Building Docker image (with cache busting)..."
CACHE_BUST=$(date +%s)  # Unix timestamp for cache busting
cd wizardcore-backend
docker build \
  --build-arg CACHE_BUST="${CACHE_BUST}" \
  -t "${FULL_IMAGE}" \
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
echo "   - Includes: All backend API enhancements"
echo "   - Cache busting: ✅ (timestamp: ${CACHE_BUST})"
echo ""
echo "🎯 Next steps:"
echo "   1. Deploy the updated image to your server"
echo "   2. Restart the backend container to pick up the new image"
echo ""