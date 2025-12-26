#!/bin/bash
# One-line deployment script for WizardCore
# Run this on any server with Docker installed

set -e

echo "🚀 WizardCore One-Line Deployment"
echo "================================="

# Check if Docker is installed
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check if Docker Compose is installed
if ! command -v docker-compose &> /dev/null; then
    echo "📦 Installing Docker Compose..."
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
fi

# Create deployment directory
mkdir -p /opt/wizardcore
cd /opt/wizardcore

# Download docker-compose.yml
echo "📥 Downloading configuration..."
curl -s -o docker-compose.yml https://raw.githubusercontent.com/davidpoweraus-sys/wizardcore/main/simple-docker-compose.yml

# Or use local file if available
if [ -f /home/glbsi/Workbench/wizardcore/simple-docker-compose.yml ]; then
    echo "📁 Using local configuration..."
    cp /home/glbsi/Workbench/wizardcore/simple-docker-compose.yml docker-compose.yml
fi

# Pull Docker images
echo "📥 Pulling Docker images..."
docker pull limpet/wizardcore-frontend:latest || echo "⚠️  Frontend image not found, will need to build"
docker pull limpet/wizardcore-backend:latest || echo "⚠️  Backend image not found, will need to build"
docker pull supabase/gotrue:v2.184.0
docker pull postgres:15-alpine
docker pull postgres:16-alpine
docker pull redis:7-alpine
docker pull redis:7.2-alpine
docker pull judge0/judge0:latest

# Stop any existing containers
echo "🔄 Stopping existing containers..."
docker-compose down --remove-orphans 2>/dev/null || true

# Start the stack
echo "🚀 Starting WizardCore stack..."
docker-compose up -d --remove-orphans

# Wait for services
echo "⏳ Waiting for services to start..."
sleep 20

# Check status
echo "📊 Service Status:"
docker-compose ps

# Create management scripts
echo "🔧 Creating management scripts..."

cat > /usr/local/bin/wiz-start << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose up -d
EOF

cat > /usr/local/bin/wiz-stop << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose down
EOF

cat > /usr/local/bin/wiz-restart << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose restart
EOF

cat > /usr/local/bin/wiz-logs << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose logs -f
EOF

cat > /usr/local/bin/wiz-status << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose ps
EOF

chmod +x /usr/local/bin/wiz-*

echo ""
echo "🎉 Deployment complete!"
echo ""
echo "📋 Services running on:"
echo "   Frontend:    http://localhost:3000"
echo "   Backend API: http://localhost:8080"
echo "   Auth:        http://localhost:9999"
echo "   Judge0:      http://localhost:2358"
echo ""
echo "🔧 Management commands:"
echo "   wiz-start    - Start all services"
echo "   wiz-stop     - Stop all services"
echo "   wiz-restart  - Restart all services"
echo "   wiz-logs     - View logs"
echo "   wiz-status   - Check status"
echo ""
echo "💡 To make accessible from outside:"
echo "   1. Configure firewall to allow ports 3000, 8080, 9999, 2358"
echo "   2. Set up Nginx reverse proxy for domain access"
echo "   3. Update environment variables for your domain"