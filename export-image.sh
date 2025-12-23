#!/bin/bash

# Export Docker image to tar file for Repoflow or direct server deployment
# This avoids registry issues by transferring the image directly

set -e

IMAGE_NAME="ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest"
OUTPUT_FILE="wizardcore-frontend-image"

echo "📦 Exporting Docker Image to Tar Archive"
echo "==========================================="
echo ""
echo "Image: ${IMAGE_NAME}"
echo ""

# Check if image exists
if ! docker image inspect "${IMAGE_NAME}" > /dev/null 2>&1; then
    echo "❌ Error: Image '${IMAGE_NAME}' not found!"
    echo ""
    echo "Build the image first by running:"
    echo "  ./build-and-push.sh"
    echo ""
    exit 1
fi

# Save to tar
echo "💾 Saving image to tar file..."
docker save "${IMAGE_NAME}" -o "${OUTPUT_FILE}.tar"

# Get tar file size
TAR_SIZE=$(du -h "${OUTPUT_FILE}.tar" | cut -f1)
echo "✅ Created: ${OUTPUT_FILE}.tar (${TAR_SIZE})"
echo ""

# Compress with gzip
echo "🗜️  Compressing with gzip..."
gzip -f -k "${OUTPUT_FILE}.tar"

# Get compressed size
GZ_SIZE=$(du -h "${OUTPUT_FILE}.tar.gz" | cut -f1)
echo "✅ Created: ${OUTPUT_FILE}.tar.gz (${GZ_SIZE})"
echo ""

# Show comparison
echo "📊 Compression Results:"
echo "-------------------------------------------"
echo "  Uncompressed: ${TAR_SIZE}"
echo "  Compressed:   ${GZ_SIZE}"
echo ""

# Calculate percentage (rough estimate)
TAR_BYTES=$(stat -c%s "${OUTPUT_FILE}.tar")
GZ_BYTES=$(stat -c%s "${OUTPUT_FILE}.tar.gz")
PERCENT=$((GZ_BYTES * 100 / TAR_BYTES))
echo "  Size reduction: ~$((100 - PERCENT))%"
echo "-------------------------------------------"
echo ""

echo "✅ Export Complete!"
echo ""
echo "📁 Files created:"
echo "  - ${OUTPUT_FILE}.tar       (uncompressed)"
echo "  - ${OUTPUT_FILE}.tar.gz    (compressed) ⭐"
echo ""
echo "🚀 Next steps:"
echo ""
echo "Option 1: Direct server deployment (fastest)"
echo "  scp ${OUTPUT_FILE}.tar.gz user@coolify-server:/tmp/"
echo "  ssh user@coolify-server"
echo "  gunzip -c /tmp/${OUTPUT_FILE}.tar.gz | docker load"
echo ""
echo "Option 2: Upload to Repoflow"
echo "  docker load -i ${OUTPUT_FILE}.tar"
echo "  docker tag ${IMAGE_NAME} repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest"
echo "  docker push repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest"
echo ""
echo "See REPOFLOW-DEPLOYMENT.md for more options!"
echo ""
