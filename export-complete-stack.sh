#!/bin/bash

# Export Complete WizardCore Stack
# Saves all Docker images needed for deployment to a single tar.gz file

set -e

OUTPUT_FILE="wizardcore-complete-stack"

echo "📦 WizardCore Complete Stack Exporter"
echo "=========================================="
echo ""

# Build backend if needed
echo "🏗️  Building backend image..."
cd wizardcore-backend
docker build -t wizardcore-backend:latest .
cd ..
echo "✅ Backend built"
echo ""

# Pull all required images
echo "📥 Pulling required base images..."
docker pull supabase/gotrue:v2.184.0
docker pull postgres:15-alpine
docker pull postgres:16-alpine  
docker pull redis:7-alpine
docker pull redis:7.2-alpine
docker pull judge0/judge0:latest
echo "✅ All base images pulled"
echo ""

# Export all images
echo "💾 Exporting complete stack..."
echo ""
echo "   Including:"
echo "   - Frontend (Next.js)"
echo "   - Backend (Go API)"
echo "   - Supabase Auth (GoTrue)"
echo "   - PostgreSQL 15 & 16"
echo "   - Redis 7 & 7.2"
echo "   - Judge0 + Worker"
echo ""
echo "⏱️  This may take 5-10 minutes..."
echo ""

docker save \
  ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest \
  wizardcore-backend:latest \
  supabase/gotrue:v2.184.0 \
  postgres:15-alpine \
  postgres:16-alpine \
  redis:7-alpine \
  redis:7.2-alpine \
  judge0/judge0:latest \
  -o "${OUTPUT_FILE}.tar"

TAR_SIZE=$(du -h "${OUTPUT_FILE}.tar" | cut -f1)
echo ""
echo "✅ Tar created: ${OUTPUT_FILE}.tar (${TAR_SIZE})"
echo ""

# Compress
echo "🗜️  Compressing (this takes 10-15 minutes)..."
gzip -v -f "${OUTPUT_FILE}.tar"

GZ_SIZE=$(du -h "${OUTPUT_FILE}.tar.gz" | cut -f1)
echo ""
echo "✅ Compressed: ${OUTPUT_FILE}.tar.gz (${GZ_SIZE})"
echo ""

# Summary
echo "📊 Export Summary:"
echo "-------------------------------------------"
echo "  File: ${OUTPUT_FILE}.tar.gz"
echo "  Size: ${GZ_SIZE}"
echo "  Images: 8 total"
echo "-------------------------------------------"
echo ""

echo "✅ Export Complete!"
echo ""
echo "🚀 Next steps:"
echo ""
echo "1. Transfer to your server:"
echo "   scp ${OUTPUT_FILE}.tar.gz user@coolify-server:/tmp/"
echo ""
echo "2. SSH into server:"
echo "   ssh user@coolify-server"
echo ""
echo "3. Load the stack:"
echo "   cd /tmp"
echo "   gunzip -c ${OUTPUT_FILE}.tar.gz | docker load"
echo ""
echo "4. Deploy:"
echo "   cd /path/to/wizardcore"
echo "   docker-compose -f docker-compose.prod.yml up -d"
echo ""
echo "See COMPLETE-STACK-DEPLOYMENT.md for detailed instructions!"
echo ""
