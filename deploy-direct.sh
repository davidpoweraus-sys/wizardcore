#!/bin/bash
# Direct deployment script for WizardCore
# Run this on your server to deploy the entire stack

set -e

echo "🚀 WizardCore Direct Server Deployment"
echo "======================================"

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo "⚠️  Please run as root or with sudo"
    exit 1
fi

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Installing Docker..."
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Create deployment directory
DEPLOY_DIR="/opt/wizardcore"
echo "📁 Setting up deployment directory at $DEPLOY_DIR..."
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

# Clone or update repository
if [ -d "$DEPLOY_DIR/.git" ]; then
    echo "📥 Updating existing repository..."
    git pull origin main
else
    echo "📥 Cloning repository..."
    git clone https://github.com/davidpoweraus-sys/wizardcore.git .
fi

# Set environment variables
echo "🔧 Configuring environment variables..."
cat > $DEPLOY_DIR/.env << 'ENV_EOF'
# Database
DATABASE_PASSWORD=wizardcore_password
POSTGRES_USER=wizardcore
POSTGRES_PASSWORD=wizardcore_password
POSTGRES_DB=wizardcore

# Supabase Auth
SUPABASE_JWT_SECRET=aNBMlTAljHvKyR2dPu6R6nyggeW2398Na3R4XL1+oyebUDiuzSO61nZzoVmRi0h4
GOTRUE_API_EXTERNAL_URL=https://auth.offensivewizard.com
GOTRUE_CORS_ALLOWED_ORIGINS=https://offensivewizard.com
GOTRUE_DB_AUTOMIGRATE=true

# Frontend
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/supabase-proxy
NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api

# Backend
SUPABASE_INTERNAL_URL=http://supabase-auth:9999
SUPABASE_URL=http://supabase-auth:9999
JUDGE0_API_URL=http://judge0:2358
REDIS_URL=redis:6379
REDIS_PASSWORD=xZDxSwVXHxVf6FYksaeEoA==
CORS_ALLOWED_ORIGINS=https://offensivewizard.com
PORT=8080
ENVIRONMENT=production

# Judge0
JUDGE0_REDIS_PASSWORD=judge0
JUDGE0_POSTGRES_PASSWORD=judge0
ENV_EOF

# Pull Docker images
echo "📥 Pulling Docker images..."
docker pull limpet/wizardcore-frontend:latest || echo "⚠️  Could not pull frontend image"
docker pull limpet/wizardcore-backend:latest || echo "⚠️  Could not pull backend image"
docker pull supabase/gotrue:v2.184.0
docker pull postgres:15-alpine
docker pull postgres:16-alpine
docker pull redis:7-alpine
docker pull redis:7.2-alpine
docker pull judge0/judge0:latest

# Build images if not available
if ! docker image inspect limpet/wizardcore-frontend:latest > /dev/null 2>&1; then
    echo "🔨 Building frontend from source..."
    docker build -f Dockerfile.nextjs -t limpet/wizardcore-frontend:latest .
fi

if ! docker image inspect limpet/wizardcore-backend:latest > /dev/null 2>&1; then
    echo "🔨 Building backend from source..."
    docker build -f wizardcore-backend/Dockerfile -t limpet/wizardcore-backend:latest wizardcore-backend/
fi

# Stop existing containers
echo "🔄 Stopping existing containers..."
docker-compose -f docker-compose.coolify.yml down --remove-orphans 2>/dev/null || true

# Start the stack
echo "🚀 Starting WizardCore stack..."
docker-compose -f docker-compose.coolify.yml up -d --remove-orphans

# Wait for services to start
echo "⏳ Waiting for services to be healthy..."
sleep 15

# Check service status
echo "📊 Service Status:"
echo "-----------------"
docker-compose -f docker-compose.coolify.yml ps

# Create systemd service for auto-start
echo "🔧 Creating systemd service..."
cat > /etc/systemd/system/wizardcore.service << 'SERVICE_EOF'
[Unit]
Description=WizardCore Application Stack
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/wizardcore
ExecStart=/usr/local/bin/docker-compose -f docker-compose.coolify.yml up -d
ExecStop=/usr/local/bin/docker-compose -f docker-compose.coolify.yml down
TimeoutStartSec=0

[Install]
WantedBy=multi-user.target
SERVICE_EOF

# Create management scripts
echo "🔧 Creating management scripts..."
cat > /usr/local/bin/wizardcore-start << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose -f docker-compose.coolify.yml up -d
EOF

cat > /usr/local/bin/wizardcore-stop << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose -f docker-compose.coolify.yml down
EOF

cat > /usr/local/bin/wizardcore-restart << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose -f docker-compose.coolify.yml restart
EOF

cat > /usr/local/bin/wizardcore-logs << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose -f docker-compose.coolify.yml logs -f
EOF

cat > /usr/local/bin/wizardcore-status << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose -f docker-compose.coolify.yml ps
EOF

cat > /usr/local/bin/wizardcore-update << 'EOF'
#!/bin/bash
cd /opt/wizardcore
git pull origin main
docker pull limpet/wizardcore-frontend:latest 2>/dev/null || true
docker pull limpet/wizardcore-backend:latest 2>/dev/null || true
docker-compose -f docker-compose.coolify.yml down
docker-compose -f docker-compose.coolify.yml up -d
EOF

chmod +x /usr/local/bin/wizardcore-*

# Enable and start systemd service
echo "🔧 Enabling auto-start..."
systemctl daemon-reload
systemctl enable wizardcore.service
systemctl start wizardcore.service

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "📋 Access URLs:"
echo "   Frontend: https://offensivewizard.com"
echo "   Backend API: https://offensivewizard.com/api"
echo "   Supabase Auth: https://offensivewizard.com/supabase-proxy"
echo ""
echo "🔧 Management commands:"
echo "   wizardcore-start    - Start services"
echo "   wizardcore-stop     - Stop services"
echo "   wizardcore-restart  - Restart services"
echo "   wizardcore-logs     - View logs"
echo "   wizardcore-status   - Check status"
echo "   wizardcore-update   - Update to latest version"
echo ""
echo "📊 Current status:"
docker-compose -f docker-compose.coolify.yml ps