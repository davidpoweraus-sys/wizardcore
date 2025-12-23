#!/bin/bash

# Script to download and load Docker images from GitHub releases
# This runs BEFORE docker-compose up to preload images

set -e

REPO="davidpoweraus-sys/wizardcore"
PACKAGE_TYPE="${PACKAGE_TYPE:-cors-fix}"
RELEASE_TAG="${RELEASE_TAG:-latest}"

echo "════════════════════════════════════════════════════════════════"
echo "🚀 Loading WizardCore Images from GitHub Release"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "📦 Package: $PACKAGE_TYPE"
echo "🏷️  Release: $RELEASE_TAG"
echo ""

# Determine package filename
if [ "$PACKAGE_TYPE" = "complete-stack" ]; then
    PACKAGE_NAME="wizardcore-complete-stack.tar.gz"
elif [ "$PACKAGE_TYPE" = "cors-fix" ]; then
    PACKAGE_NAME="wizardcore-cors-fix.tar.gz"
else
    echo "❌ Error: Invalid PACKAGE_TYPE. Use 'cors-fix' or 'complete-stack'"
    exit 1
fi

# Get download URL
echo "🔍 Fetching release info from GitHub..."

if [ "$RELEASE_TAG" = "latest" ]; then
    API_URL="https://api.github.com/repos/$REPO/releases/latest"
else
    API_URL="https://api.github.com/repos/$REPO/releases/tags/$RELEASE_TAG"
fi

# Fetch release data
RELEASE_DATA=$(curl -s "$API_URL")

if echo "$RELEASE_DATA" | grep -q "Not Found"; then
    echo "❌ Release not found: $RELEASE_TAG"
    exit 1
fi

# Extract download URL
DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep -o "https://github.com/$REPO/releases/download/[^\"]*/$PACKAGE_NAME" | head -1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "❌ Could not find $PACKAGE_NAME in release"
    exit 1
fi

TAG_NAME=$(echo "$RELEASE_DATA" | grep '"tag_name"' | cut -d'"' -f4)
echo "✅ Found release: $TAG_NAME"
echo "📥 Downloading: $DOWNLOAD_URL"
echo ""

# Check if images are already loaded (skip if already present)
if docker image inspect ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest >/dev/null 2>&1 && \
   docker image inspect supabase/gotrue:v2.184.0 >/dev/null 2>&1; then
    echo "✅ Images already loaded - skipping download"
    echo ""
    exit 0
fi

# Download tar file
TEMP_FILE="/tmp/wizardcore-images-$$.tar.gz"

echo "⬇️  Downloading images..."
if ! curl -L --progress-bar -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo "❌ Download failed"
    rm -f "$TEMP_FILE"
    exit 1
fi

echo ""
echo "✅ Download complete ($(du -h "$TEMP_FILE" | cut -f1))"
echo ""

# Load images into Docker
echo "📦 Loading images into Docker..."
if gunzip -c "$TEMP_FILE" | docker load; then
    echo ""
    echo "✅ Images loaded successfully!"
else
    echo ""
    echo "❌ Failed to load images"
    rm -f "$TEMP_FILE"
    exit 1
fi

# Cleanup
rm -f "$TEMP_FILE"

# Show what was loaded
echo ""
echo "📋 Loaded images:"
docker images | grep -E "(wizardcore|gotrue)" | head -10
echo ""
echo "════════════════════════════════════════════════════════════════"
echo "✅ Images ready! Proceeding with docker-compose..."
echo "════════════════════════════════════════════════════════════════"
