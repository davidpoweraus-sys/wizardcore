#!/bin/bash
# SSH Deployment Script for WizardCore - Simple and Direct
# Run this on your server via SSH to deploy everything

set -e

echo "🚀 WizardCore SSH Direct Deployment"
echo "==================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    echo -e "${YELLOW}⚠️  Please run as root or with sudo${NC}"
    exit 1
fi

# Function to check command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Install Docker if not present
if ! command_exists docker; then
    echo -e "${YELLOW}🐳 Docker not found. Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    echo -e "${GREEN}✅ Docker installed${NC}"
fi

# Install Docker Compose if not present
if ! command_exists docker-compose; then
    echo -e "${YELLOW}🐳 Docker Compose not found. Installing...${NC}"
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
fi

# Create deployment directory
DEPLOY_DIR="/opt/wizardcore"
echo -e "${YELLOW}📁 Setting up deployment directory at $DEPLOY_DIR...${NC}"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

# Create docker-compose.yml for direct deployment
echo -e "${YELLOW}📝 Creating docker-compose.yml...${NC}"
cat > $DEPLOY_DIR/docker-compose.yml << 'COMPOSE_EOF'
version: '3.8'

services:
  # Frontend - accessible at app.offensivewizard.com
  frontend:
    image: limpet/wizardcore-frontend:latest
    environment:
      - NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
      - NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com
      - NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
      - SUPABASE_INTERNAL_URL=http://auth:9999
    ports:
      - "3000:3000"
    depends_on:
      - backend
      - auth
    restart: unless-stopped
    networks:
      - wizardcore

  # Backend API - accessible at api.offensivewizard.com
  backend:
    image: limpet/wizardcore-backend:latest
    environment:
      # Database
      - DATABASE_URL=postgresql://wizardcore:wizardcore_password@postgres:5432/wizardcore?sslmode=disable
      - DATABASE_HOST=postgres
      - DATABASE_PORT=5432
      - DATABASE_USER=wizardcore
      - DATABASE_PASSWORD=wizardcore_password
      - DATABASE_NAME=wizardcore
      
      # Auth
      - SUPABASE_URL=http://auth:9999
      - SUPABASE_JWT_SECRET=aNBMlTAljHvKyR2dPu6R6nyggeW2398Na3R4XL1+oyebUDiuzSO61nZzoVmRi0h4
      
      # Redis
      - REDIS_URL=redis:6379
      - REDIS_PASSWORD=xZDxSwVXHxVf6FYksaeEoA==
      
      # Judge0
      - JUDGE0_API_URL=http://judge0:2358
      - JUDGE0_API_KEY=
      
      # Configuration
      - CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
      - PORT=8080
      - ENVIRONMENT=production
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
      - judge0
      - auth
    restart: unless-stopped
    networks:
      - wizardcore

  # Supabase Auth - accessible at auth.offensivewizard.com
  auth:
    image: supabase/gotrue:v2.184.0
    environment:
      - GOTRUE_API_HOST=0.0.0.0
      - GOTRUE_API_PORT=9999
      - GOTRUE_DB_DRIVER=postgres
      - GOTRUE_DB_DATABASE_URL=postgresql://supabase_auth_admin:password@auth-db:5432/supabase_auth?sslmode=disable
      - GOTRUE_SITE_URL=https://app.offensivewizard.com
      - API_EXTERNAL_URL=https://auth.offensivewizard.com
      - GOTRUE_CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
      - GOTRUE_CORS_ALLOWED_HEADERS=Authorization,Content-Type,X-Client-Info,X-Requested-With,apikey,x-client-info,x-supabase-api-version,Accept,Accept-Language,Content-Language
      - GOTRUE_CORS_ALLOWED_METHODS=GET,POST,PUT,PATCH,DELETE,OPTIONS
      - GOTRUE_CORS_EXPOSED_HEADERS=X-Total-Count
      - GOTRUE_CORS_ALLOW_CREDENTIALS=true
      - GOTRUE_DB_AUTOMIGRATE=true
    ports:
      - "9999:9999"
    depends_on:
      - auth-db
    restart: unless-stopped
    networks:
      - wizardcore

  # Judge0 API - accessible at judge0.offensivewizard.com
  judge0:
    image: judge0/judge0:latest
    environment:
      - REDIS_HOST=judge0-redis
      - REDIS_PORT=6379
      - POSTGRES_HOST=judge0-db
      - POSTGRES_PORT=5432
      - POSTGRES_USER=judge0
      - POSTGRES_PASSWORD=judge0
      - POSTGRES_DB=judge0
    ports:
      - "2358:2358"
    depends_on:
      - judge0-db
      - judge0-redis
    restart: unless-stopped
    networks:
      - wizardcore

  # Judge0 Worker
  judge0-worker:
    image: judge0/judge0:latest
    command: ["./scripts/workers"]
    environment:
      - REDIS_HOST=judge0-redis
      - REDIS_PORT=6379
      - POSTGRES_HOST=judge0-db
      - POSTGRES_PORT=5432
      - POSTGRES_USER=judge0
      - POSTGRES_PASSWORD=judge0
      - POSTGRES_DB=judge0
    depends_on:
      - judge0-db
      - judge0-redis
    restart: unless-stopped
    networks:
      - wizardcore

  # Database services (internal only)
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: wizardcore
      POSTGRES_PASSWORD: wizardcore_password
      POSTGRES_DB: wizardcore
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - wizardcore

  auth-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: supabase_auth_admin
      POSTGRES_PASSWORD: password
      POSTGRES_DB: supabase_auth
    volumes:
      - auth_db_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - wizardcore

  judge0-db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: judge0
      POSTGRES_PASSWORD: judge0
      POSTGRES_DB: judge0
    volumes:
      - judge0_db_data:/var/lib/postgresql/data
    restart: unless-stopped
    networks:
      - wizardcore

  redis:
    image: redis:7-alpine
    command: redis-server --requirepass xZDxSwVXHxVf6FYksaeEoA== --bind 0.0.0.0
    volumes:
      - redis_data:/data
    restart: unless-stopped
    networks:
      - wizardcore

  judge0-redis:
    image: redis:7.2-alpine
    command: redis-server --appendonly yes --bind 0.0.0.0
    volumes:
      - judge0_redis_data:/data
    restart: unless-stopped
    networks:
      - wizardcore

volumes:
  postgres_data:
  auth_db_data:
  judge0_db_data:
  redis_data:
  judge0_redis_data:

networks:
  wizardcore:
    driver: bridge
COMPOSE_EOF

# Create .env file
echo -e "${YELLOW}🔧 Creating environment file...${NC}"
cat > $DEPLOY_DIR/.env << 'ENV_EOF'
# Domain configuration
DOMAIN=offensivewizard.com
FRONTEND_URL=https://app.offensivewizard.com
BACKEND_URL=https://api.offensivewizard.com
AUTH_URL=https://auth.offensivewizard.com
JUDGE0_URL=https://judge0.offensivewizard.com

# Database passwords
DATABASE_PASSWORD=wizardcore_password
AUTH_DB_PASSWORD=password
JUDGE0_DB_PASSWORD=judge0

# Redis password
REDIS_PASSWORD=xZDxSwVXHxVf6FYksaeEoA==

# JWT secret
SUPABASE_JWT_SECRET=aNBMlTAljHvKyR2dPu6R6nyggeW2398Na3R4XL1+oyebUDiuzSO61nZzoVmRi0h4
ENV_EOF

# Pull Docker images
echo -e "${YELLOW}📥 Pulling Docker images...${NC}"
docker pull limpet/wizardcore-frontend:latest || echo -e "${YELLOW}⚠️  Could not pull frontend image${NC}"
docker pull limpet/wizardcore-backend:latest || echo -e "${YELLOW}⚠️  Could not pull backend image${NC}"
docker pull supabase/gotrue:v2.184.0
docker pull postgres:15-alpine
docker pull postgres:16-alpine
docker pull redis:7-alpine
docker pull redis:7.2-alpine
docker pull judge0/judge0:latest

# Stop any existing containers
echo -e "${YELLOW}🔄 Stopping any existing containers...${NC}"
cd $DEPLOY_DIR
docker-compose down --remove-orphans 2>/dev/null || true

# Start the stack
echo -e "${YELLOW}🚀 Starting WizardCore stack...${NC}"
docker-compose up -d --remove-orphans

# Wait for services to start
echo -e "${YELLOW}⏳ Waiting for services to start (30 seconds)...${NC}"
sleep 30

# Check service status
echo -e "${GREEN}📊 Service Status:${NC}"
echo "-----------------"
docker-compose ps

# Create management scripts
echo -e "${YELLOW}🔧 Creating management scripts...${NC}"

# Start script
cat > /usr/local/bin/wizardcore-start << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose up -d
EOF

# Stop script
cat > /usr/local/bin/wizardcore-stop << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose down
EOF

# Restart script
cat > /usr/local/bin/wizardcore-restart << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose restart
EOF

# Logs script
cat > /usr/local/bin/wizardcore-logs << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose logs -f
EOF

# Status script
cat > /usr/local/bin/wizardcore-status << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose ps
EOF

# Update script
cat > /usr/local/bin/wizardcore-update << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose pull
docker-compose down
docker-compose up -d
EOF

# Make scripts executable
chmod +x /usr/local/bin/wizardcore-*

echo ""
echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo ""
echo -e "${GREEN}📋 Access URLs:${NC}"
echo "   Frontend:    https://app.offensivewizard.com"
echo "   Backend API: https://api.offensivewizard.com"
echo "   Auth:        https://auth.offensivewizard.com"
echo "   Judge0:      https://judge0.offensivewizard.com"
echo ""
echo -e "${GREEN}🔧 Management commands:${NC}"
echo "   wizardcore-start    - Start all services"
echo "   wizardcore-stop     - Stop all services"
echo "   wizardcore-restart  - Restart all services"
echo "   wizardcore-logs     - View logs (follow mode)"
echo "   wizardcore-status   - Check service status"
echo "   wizardcore-update   - Update to latest images"
echo ""
echo -e "${GREEN}📊 Checking service health...${NC}"

# Check if services are responding
echo "Frontend:"
if curl -s -f http://localhost:3000 > /dev/null; then
    echo -e "  ✅ Running on port 3000"
else
    echo -e "  ❌ Not responding on port 3000"
fi

echo "Backend:"
if curl -s -f http://localhost:8080/health > /dev/null; then
    echo -e "  ✅ Running on port 8080"
else
    echo -e "  ❌ Not responding on port 8080"
fi

echo "Auth:"
if curl -s -f http://localhost:9999/health > /dev/null; then
    echo -e "  ✅ Running on port 9999"
else
    echo -e "  ❌ Not responding on port 9999"
fi

echo ""
echo -e "${GREEN}✅ Done! Now you need to:${NC}"
echo "1. Configure your reverse proxy (Nginx/Apache) to route:"
echo "   - app.offensivewizard.com → localhost:3000"
echo "   - api.offensivewizard.com → localhost:8080"
echo "   - auth.offensivewizard.com → localhost:9999"
echo "   - judge0.offensivewizard.com → localhost:2358"
echo ""
echo "2. Set up SSL certificates (Let's Encrypt)"
echo ""
echo "3. Run: wizardcore-logs to monitor startup"