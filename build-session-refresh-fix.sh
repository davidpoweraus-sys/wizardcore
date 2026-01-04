#!/bin/bash

# Build script for session refresh & CORS fix
# This builds the frontend Docker image with the login fix

set -e

echo "🔨 Building session refresh & CORS fix Docker image..."

# Build the image
docker build -t limpet/wizardcore-frontend:session-refresh-cors-fix -f Dockerfile.nextjs .

echo "✅ Image built successfully: limpet/wizardcore-frontend:session-refresh-cors-fix"

echo "📤 Pushing to Docker Hub..."
docker push limpet/wizardcore-frontend:session-refresh-cors-fix

echo "🎉 Image pushed successfully!"
echo ""
echo "📋 Next steps:"
echo "1. Update your deployment to use: limpet/wizardcore-frontend:session-refresh-cors-fix"
echo "2. Redeploy the frontend service"
echo "3. Clear browser cache and test login"
echo ""
echo "📝 Fix includes:"
echo "   - CORS fix for same-origin requests in auth/backend proxies"
echo "   - Session refresh awareness in middleware"
echo "   - Version: session-refresh-fix-20260104-1159"