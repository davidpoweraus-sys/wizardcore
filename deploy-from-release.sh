#!/bin/bash

# WizardCore Deployment from GitHub Release
# Downloads and loads Docker images from GitHub releases

set -e

REPO="davidpoweraus-sys/wizardcore"
PACKAGE_TYPE="${1:-cors-fix}"  # cors-fix or complete-stack
RELEASE_TAG="${2:-latest}"     # Release tag or "latest"

echo "╔══════════════════════════════════════════════════════════════════════════════╗"
echo "║              WizardCore Deployment from GitHub Release                      ║"
echo "╚══════════════════════════════════════════════════════════════════════════════╝"
echo ""

# Determine package filename
if [ "$PACKAGE_TYPE" = "complete-stack" ]; then
    PACKAGE_NAME="wizardcore-complete-stack.tar.gz"
    SERVICES="all 8 services"
elif [ "$PACKAGE_TYPE" = "cors-fix" ]; then
    PACKAGE_NAME="wizardcore-cors-fix.tar.gz"
    SERVICES="frontend + supabase-auth"
else
    echo "❌ Error: Invalid package type '$PACKAGE_TYPE'"
    echo ""
    echo "Usage: $0 [cors-fix|complete-stack] [release-tag]"
    echo ""
    echo "Examples:"
    echo "  $0 cors-fix latest              # Download latest CORS fix"
    echo "  $0 complete-stack latest        # Download latest complete stack"
    echo "  $0 cors-fix v2024.12.23-1430   # Download specific version"
    exit 1
fi

echo "📦 Package: $PACKAGE_NAME"
echo "🏷️  Release: $RELEASE_TAG"
echo "📋 Services: $SERVICES"
echo ""

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "❌ Error: Docker is not installed or not in PATH"
    exit 1
fi

# Check Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker daemon is not running"
    echo "   Try: sudo systemctl start docker"
    exit 1
fi

# Get the download URL
echo "🔍 Fetching release information..."

if [ "$RELEASE_TAG" = "latest" ]; then
    API_URL="https://api.github.com/repos/$REPO/releases/latest"
else
    API_URL="https://api.github.com/repos/$REPO/releases/tags/$RELEASE_TAG"
fi

# Fetch release data
RELEASE_DATA=$(curl -s "$API_URL")

# Check if release exists
if echo "$RELEASE_DATA" | grep -q "Not Found"; then
    echo "❌ Error: Release '$RELEASE_TAG' not found"
    echo ""
    echo "Available releases:"
    curl -s "https://api.github.com/repos/$REPO/releases" | grep '"tag_name"' | head -5
    exit 1
fi

# Extract download URL for the package
DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep -o "https://github.com/$REPO/releases/download/[^\"]*/$PACKAGE_NAME" | head -1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Error: Could not find $PACKAGE_NAME in release"
    echo ""
    echo "Available assets in this release:"
    echo "$RELEASE_DATA" | grep '"name"' | grep '.tar.gz'
    exit 1
fi

echo "✅ Found release: $(echo "$RELEASE_DATA" | grep '"tag_name"' | cut -d'"' -f4)"
echo "📥 Download URL: $DOWNLOAD_URL"
echo ""

# Download the package
TEMP_DIR=$(mktemp -d)
DOWNLOAD_PATH="$TEMP_DIR/$PACKAGE_NAME"

echo "⬇️  Downloading $PACKAGE_NAME..."
echo "   This may take a few minutes depending on package size..."
echo ""

if curl -L --progress-bar -o "$DOWNLOAD_PATH" "$DOWNLOAD_URL"; then
    echo ""
    echo "✅ Download complete!"
    ls -lh "$DOWNLOAD_PATH"
else
    echo ""
    echo "❌ Error: Download failed"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Decompress and load
echo ""
echo "📦 Loading Docker images..."
echo "   This will take 5-10 minutes for complete stack, ~1 minute for CORS fix..."
echo ""

if gunzip -c "$DOWNLOAD_PATH" | docker load; then
    echo ""
    echo "✅ Images loaded successfully!"
else
    echo ""
    echo "❌ Error: Failed to load images"
    rm -rf "$TEMP_DIR"
    exit 1
fi

# Cleanup
rm -rf "$TEMP_DIR"

# Show loaded images
echo ""
echo "📋 Loaded images:"
echo "═══════════════════════════════════════════════════════════════"
docker images | grep -E "(wizardcore|supabase|postgres|redis|judge0)" | head -20
echo "═══════════════════════════════════════════════════════════════"
echo ""

# Next steps
echo "✅ Deployment Complete!"
echo ""
echo "🎯 Next steps:"
echo ""

if [ "$PACKAGE_TYPE" = "cors-fix" ]; then
    echo "1. Restart affected services:"
    echo "   docker-compose -f docker-compose.prod.yml up -d frontend supabase-auth"
    echo ""
    echo "2. Or redeploy in Coolify UI (recommended)"
    echo ""
    echo "3. Test CORS:"
    echo "   Open https://offensivewizard.com and try to login"
else
    echo "1. Deploy all services:"
    echo "   docker-compose -f docker-compose.prod.yml up -d"
    echo ""
    echo "2. Or redeploy in Coolify UI (recommended)"
    echo ""
    echo "3. Check status:"
    echo "   docker-compose -f docker-compose.prod.yml ps"
fi

echo ""
echo "═══════════════════════════════════════════════════════════════"
echo "🎉 Images are ready! Redeploy when ready."
echo "═══════════════════════════════════════════════════════════════"
