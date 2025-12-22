#!/bin/bash

# Build and push frontend to GitHub Container Registry
# This avoids memory issues on Coolify server

set -e

GITHUB_USERNAME="davidpoweraus-sys"
IMAGE_NAME="ghcr.io/${GITHUB_USERNAME}/wizardcore-frontend"
TAG="latest"

echo "🏗️  Building Frontend Docker Image..."
echo "This may take 5-10 minutes and requires 8GB+ RAM"
echo ""

# Build the image
docker build \
  -t "${IMAGE_NAME}:${TAG}" \
  -f Dockerfile.nextjs \
  --build-arg NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com \
  --build-arg NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps= \
  --build-arg NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com \
  --build-arg NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api \
  .

echo ""
echo "✅ Build complete!"
echo ""

# Check if logged in
if ! docker info 2>/dev/null | grep -q "ghcr.io"; then
    echo "⚠️  Not logged in to GitHub Container Registry"
    echo ""
    echo "Please run:"
    echo "  echo YOUR_GITHUB_TOKEN | docker login ghcr.io -u ${GITHUB_USERNAME} --password-stdin"
    echo ""
    echo "Get token from: https://github.com/settings/tokens"
    echo "Required scopes: write:packages, read:packages"
    echo ""
    exit 1
fi

echo "📤 Pushing to GitHub Container Registry..."
docker push "${IMAGE_NAME}:${TAG}"

echo ""
echo "✅ Success!"
echo ""
echo "Image pushed to: ${IMAGE_NAME}:${TAG}"
echo ""
echo "Next steps:"
echo "1. Commit docker-compose.prod.yml changes"
echo "2. Push to git"
echo "3. Redeploy in Coolify (will pull pre-built image)"
echo ""
