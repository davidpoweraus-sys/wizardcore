#!/bin/bash
set -e

echo "🚀 WizardCore Deployment from GitHub Release"
echo "=============================================="
echo ""

# Configuration
REPO="davidpoweraus-sys/wizardcore"
RELEASE_TAG="${1:-latest}"
SERVER="root@172.105.181.38"
COOLIFY_APP_DIR="/data/coolify/applications/d44co4gk48kok84wcg8o0os0"

echo "📦 Repository: $REPO"
echo "🏷️  Release: $RELEASE_TAG"
echo "🖥️  Server: $SERVER"
echo ""

# Get release info
echo "🔍 Fetching release information..."
if [ "$RELEASE_TAG" = "latest" ]; then
    RELEASE_URL="https://api.github.com/repos/$REPO/releases/latest"
else
    RELEASE_URL="https://api.github.com/repos/$REPO/releases/tags/$RELEASE_TAG"
fi

# Download URL for app package
DOWNLOAD_URL="https://github.com/$REPO/releases/download/$RELEASE_TAG/wizardcore-app.tar.gz"
if [ "$RELEASE_TAG" = "latest" ]; then
    DOWNLOAD_URL="https://github.com/$REPO/releases/latest/download/wizardcore-app.tar.gz"
fi

echo "📥 Download URL: $DOWNLOAD_URL"
echo ""

# Download to server
echo "📡 Downloading release to server..."
ssh $SERVER "curl -L '$DOWNLOAD_URL' -o /home/wizardcore-app.tar.gz"

# Check if download succeeded
ssh $SERVER "ls -lh /home/wizardcore-app.tar.gz"
echo ""

# Load images
echo "📦 Loading Docker images on server..."
ssh $SERVER "gunzip -c /home/wizardcore-app.tar.gz | docker load"
echo ""

# Restart services
echo "🔄 Restarting services with new images..."
ssh $SERVER "cd $COOLIFY_APP_DIR && docker compose up -d --force-recreate frontend backend"
echo ""

# Verify deployment
echo "✅ Verifying deployment..."
sleep 5
ssh $SERVER "docker ps --format 'table {{.Names}}\t{{.Image}}\t{{.Status}}' | grep -E 'frontend|backend|NAME'"
echo ""

# Cleanup
echo "🧹 Cleaning up..."
ssh $SERVER "rm -f /home/wizardcore-app.tar.gz"
echo ""

echo "✅ Deployment complete!"
echo ""
echo "🌐 Your app: https://offensivewizard.com"
echo "📊 Check logs: ssh $SERVER 'docker logs frontend-d44co4gk48kok84wcg8o0os0'"
