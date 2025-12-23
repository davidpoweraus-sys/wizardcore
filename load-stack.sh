#!/bin/bash

# Load WizardCore Complete Stack from tar.gz archive
# Run this script on your target server (Coolify, VPS, etc.)

set -e

STACK_FILE="wizardcore-complete-stack.tar.gz"

echo "🚀 WizardCore Stack Loader"
echo "=========================================="
echo ""

# Check if file exists
if [ ! -f "${STACK_FILE}" ]; then
    echo "❌ Error: ${STACK_FILE} not found!"
    echo ""
    echo "Please copy the stack file to this directory first:"
    echo "  scp ${STACK_FILE} user@server:~/"
    echo ""
    exit 1
fi

# Check Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Error: Docker is not running or you don't have permission"
    echo ""
    echo "Try:"
    echo "  sudo systemctl start docker"
    echo "  sudo usermod -aG docker $USER"
    echo ""
    exit 1
fi

# Show file info
FILE_SIZE=$(du -h "${STACK_FILE}" | cut -f1)
echo "📦 Stack archive: ${STACK_FILE} (${FILE_SIZE})"
echo ""

# Load the images
echo "⏱️  Loading all Docker images..."
echo "   This will take 5-10 minutes..."
echo ""

gunzip -c "${STACK_FILE}" | docker load

echo ""
echo "✅ All images loaded successfully!"
echo ""

# Show loaded images
echo "📋 Loaded images:"
echo "-------------------------------------------"
docker images | grep -E "(wizardcore|supabase|postgres|redis|judge0)" | head -20
echo "-------------------------------------------"
echo ""

echo "✅ Stack Load Complete!"
echo ""
echo "🎯 Next steps:"
echo ""
echo "1. Clone your repository (if not already done):"
echo "   git clone https://github.com/davidpoweraus-sys/wizardcore.git"
echo "   cd wizardcore"
echo ""
echo "2. Deploy with docker-compose:"
echo "   docker-compose -f docker-compose.prod.yml up -d"
echo ""
echo "3. Check status:"
echo "   docker-compose -f docker-compose.prod.yml ps"
echo ""
echo "4. View logs:"
echo "   docker-compose -f docker-compose.prod.yml logs -f"
echo ""
echo "🎉 Your complete WizardCore stack is ready to deploy!"
echo ""
