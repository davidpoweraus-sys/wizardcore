#!/bin/sh

# Script to load Docker images for WizardCore
# Stage 1: Pull base images if missing (one-time)
# Stage 2: Download and load app images from GitHub (every deploy)

set -e

REPO="davidpoweraus-sys/wizardcore"
PACKAGE_TYPE="${PACKAGE_TYPE:-app}"
RELEASE_TAG="${RELEASE_TAG:-latest}"

echo "════════════════════════════════════════════════════════════════"
echo "🚀 WizardCore Image Loader"
echo "════════════════════════════════════════════════════════════════"
echo ""
echo "📦 Package: $PACKAGE_TYPE"
echo "🏷️  Release: $RELEASE_TAG"
echo ""

# ============================================================================
# STAGE 1: Ensure Base Images Exist (Pull if missing)
# ============================================================================

echo "📋 Stage 1: Checking base images..."
echo "────────────────────────────────────────────────────────────────"

BASE_IMAGES="
postgres:15-alpine
postgres:16-alpine
redis:7-alpine
redis:7.2-alpine
judge0/judge0:latest
"

MISSING_COUNT=0
EXISTING_COUNT=0

for img in $BASE_IMAGES; do
  if docker image inspect "$img" >/dev/null 2>&1; then
    echo "  ✓ $img (already exists)"
    EXISTING_COUNT=$((EXISTING_COUNT + 1))
  else
    echo "  ⬇ $img (pulling...)"
    docker pull "$img" &
    MISSING_COUNT=$((MISSING_COUNT + 1))
  fi
done

# Wait for all pulls to complete
if [ $MISSING_COUNT -gt 0 ]; then
  echo ""
  echo "⏳ Waiting for $MISSING_COUNT base image(s) to download..."
  wait
  echo "✅ All base images pulled!"
else
  echo ""
  echo "✅ All base images already exist (skipped pulling)"
fi

echo ""
echo "📊 Base images summary:"
echo "   Existing: $EXISTING_COUNT"
echo "   Downloaded: $MISSING_COUNT"
echo ""

# ============================================================================
# STAGE 2: Download and Load App Images from GitHub
# ============================================================================

echo "📋 Stage 2: Loading app images from GitHub..."
echo "────────────────────────────────────────────────────────────────"

# Determine package filename based on type
if [ "$PACKAGE_TYPE" = "complete-stack" ]; then
    PACKAGE_NAME="wizardcore-complete-stack.tar.gz"
elif [ "$PACKAGE_TYPE" = "app" ]; then
    PACKAGE_NAME="wizardcore-app.tar.gz"
else
    echo "❌ Error: Invalid PACKAGE_TYPE '$PACKAGE_TYPE'"
    echo "   Valid options: app, complete-stack"
    exit 1
fi

# Get download URL from GitHub
echo "🔍 Fetching release info from GitHub..."

if [ "$RELEASE_TAG" = "latest" ]; then
    API_URL="https://api.github.com/repos/$REPO/releases/latest"
else
    API_URL="https://api.github.com/repos/$REPO/releases/tags/$RELEASE_TAG"
fi

RELEASE_DATA=$(curl -s "$API_URL")

if echo "$RELEASE_DATA" | grep -q "Not Found"; then
    echo "❌ Release not found: $RELEASE_TAG"
    echo ""
    echo "💡 This is likely the first deployment."
    echo "   Base images are loaded. You can:"
    echo "   1. Wait for GitHub Actions to create a release"
    echo "   2. Or manually load app images"
    echo ""
    echo "✅ Continuing with base images only..."
    exit 0
fi

# Extract download URL
DOWNLOAD_URL=$(echo "$RELEASE_DATA" | grep -o "https://github.com/$REPO/releases/download/[^\"]*/$PACKAGE_NAME" | head -1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "⚠️  Could not find $PACKAGE_NAME in release"
    echo ""
    echo "Available assets:"
    echo "$RELEASE_DATA" | grep '"name"' | grep '.tar.gz' || echo "  None found"
    echo ""
    echo "✅ Continuing with base images only..."
    exit 0
fi

TAG_NAME=$(echo "$RELEASE_DATA" | grep '"tag_name"' | head -1 | cut -d'"' -f4)
echo "✅ Found release: $TAG_NAME"
echo "📥 Downloading: $PACKAGE_NAME"
echo ""

# Check if app images already loaded (skip download)
if docker image inspect ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest >/dev/null 2>&1 && \
   docker image inspect wizardcore-backend:latest >/dev/null 2>&1 && \
   docker image inspect supabase/gotrue:v2.184.0 >/dev/null 2>&1; then
    echo "✅ App images already loaded - skipping download"
    echo ""
    echo "📋 Loaded app images:"
    docker images | grep -E "(wizardcore-frontend|wizardcore-backend|gotrue)" | head -10
    echo ""
    echo "════════════════════════════════════════════════════════════════"
    echo "✅ All images ready!"
    echo "════════════════════════════════════════════════════════════════"
    exit 0
fi

# Download app package
TEMP_FILE="/tmp/wizardcore-app-$$.tar.gz"

echo "⬇️  Downloading app package..."
if ! curl -L --progress-bar -o "$TEMP_FILE" "$DOWNLOAD_URL"; then
    echo "❌ Download failed"
    rm -f "$TEMP_FILE"
    exit 1
fi

echo ""
echo "✅ Download complete ($(du -h "$TEMP_FILE" | cut -f1))"
echo ""

# Load images
echo "📦 Loading app images into Docker..."
if gunzip -c "$TEMP_FILE" | docker load; then
    echo ""
    echo "✅ App images loaded successfully!"
else
    echo ""
    echo "❌ Failed to load app images"
    rm -f "$TEMP_FILE"
    exit 1
fi

# Cleanup
rm -f "$TEMP_FILE"

# Show loaded images
echo ""
echo "📋 Loaded app images:"
docker images | grep -E "(wizardcore-frontend|wizardcore-backend|gotrue)" | head -10
echo ""

echo "════════════════════════════════════════════════════════════════"
echo "✅ All images ready! (Base + App)"
echo "════════════════════════════════════════════════════════════════"
