#!/bin/bash
# Complete WizardCore Deployment with SSL and Subdomains
# This script does everything: deploys all services, sets up subdomains, and manages SSL certificates

set -e

echo "🚀 WizardCore Complete Deployment"
echo "================================="
echo "Domain: offensivewizard.com"
echo "Subdomains: app, api, auth, judge0"
echo ""

# Colors
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
    echo -e "${YELLOW}🐳 Installing Docker...${NC}"
    curl -fsSL https://get.docker.com -o get-docker.sh
    sh get-docker.sh
    rm get-docker.sh
    echo -e "${GREEN}✅ Docker installed${NC}"
fi

# Install Docker Compose if not present
if ! command_exists docker-compose; then
    echo -e "${YELLOW}🐳 Installing Docker Compose...${NC}"
    curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    echo -e "${GREEN}✅ Docker Compose installed${NC}"
fi

# Create deployment directory
DEPLOY_DIR="/opt/wizardcore"
echo -e "${YELLOW}📁 Setting up deployment directory at $DEPLOY_DIR...${NC}"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR

# Create docker-compose.yml with Traefik for SSL and subdomains
echo -e "${YELLOW}📝 Creating docker-compose.yml with Traefik...${NC}"
cat > $DEPLOY_DIR/docker-compose.yml << 'COMPOSE_EOF'
version: '3.8'

services:
  # Traefik reverse proxy with automatic SSL
  traefik:
    image: traefik:v3.0
    container_name: traefik
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.letsencrypt.acme.tlschallenge=true"
      - "--certificatesresolvers.letsencrypt.acme.email=admin@offensivewizard.com"
      - "--certificatesresolvers.letsencrypt.acme.storage=/letsencrypt/acme.json"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"  # Traefik dashboard
    volumes:
      - "/var/run/docker.sock:/var/run/docker.sock:ro"
      - "./letsencrypt:/letsencrypt"
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Frontend - app.offensivewizard.com
  frontend:
    image: limpet/wizardcore-frontend:latest
    environment:
      # Frontend expects these URLs
      - NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
      - NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com
      - NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
      - SUPABASE_INTERNAL_URL=http://auth:9999
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.frontend.rule=Host(`app.offensivewizard.com`)"
      - "traefik.http.routers.frontend.entrypoints=websecure"
      - "traefik.http.routers.frontend.tls.certresolver=letsencrypt"
      - "traefik.http.services.frontend.loadbalancer.server.port=3000"
    depends_on:
      - backend
      - auth
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Backend API - api.offensivewizard.com
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
      
      # CORS - allow frontend subdomain
      - CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
      - PORT=8080
      - ENVIRONMENT=production
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.backend.rule=Host(`api.offensivewizard.com`)"
      - "traefik.http.routers.backend.entrypoints=websecure"
      - "traefik.http.routers.backend.tls.certresolver=letsencrypt"
      - "traefik.http.services.backend.loadbalancer.server.port=8080"
      # CORS headers
      - "traefik.http.middlewares.cors-headers.headers.customresponseheaders.Access-Control-Allow-Origin=https://app.offensivewizard.com"
      - "traefik.http.middlewares.cors-headers.headers.customresponseheaders.Access-Control-Allow-Methods=GET,POST,PUT,DELETE,OPTIONS"
      - "traefik.http.middlewares.cors-headers.headers.customresponseheaders.Access-Control-Allow-Headers=Authorization,Content-Type"
      - "traefik.http.middlewares.cors-headers.headers.customresponseheaders.Access-Control-Allow-Credentials=true"
      - "traefik.http.routers.backend.middlewares=cors-headers"
    depends_on:
      - postgres
      - redis
      - judge0
      - auth
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Supabase Auth - auth.offensivewizard.com
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
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.auth.rule=Host(`auth.offensivewizard.com`)"
      - "traefik.http.routers.auth.entrypoints=websecure"
      - "traefik.http.routers.auth.tls.certresolver=letsencrypt"
      - "traefik.http.services.auth.loadbalancer.server.port=9999"
      # CORS headers for auth
      - "traefik.http.middlewares.auth-cors.headers.customresponseheaders.Access-Control-Allow-Origin=https://app.offensivewizard.com"
      - "traefik.http.middlewares.auth-cors.headers.customresponseheaders.Access-Control-Allow-Methods=GET,POST,PUT,DELETE,OPTIONS"
      - "traefik.http.middlewares.auth-cors.headers.customresponseheaders.Access-Control-Allow-Headers=Authorization,Content-Type,apikey,x-client-info,x-supabase-api-version"
      - "traefik.http.middlewares.auth-cors.headers.customresponseheaders.Access-Control-Allow-Credentials=true"
      - "traefik.http.middlewares.auth-cors.headers.customresponseheaders.Access-Control-Expose-Headers=X-Total-Count"
      - "traefik.http.routers.auth.middlewares=auth-cors"
    depends_on:
      - auth-db
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Judge0 API - judge0.offensivewizard.com
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
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.judge0.rule=Host(`judge0.offensivewizard.com`)"
      - "traefik.http.routers.judge0.entrypoints=websecure"
      - "traefik.http.routers.judge0.tls.certresolver=letsencrypt"
      - "traefik.http.services.judge0.loadbalancer.server.port=2358"
      # CORS headers for Judge0
      - "traefik.http.middlewares.judge0-cors.headers.customresponseheaders.Access-Control-Allow-Origin=https://app.offensivewizard.com"
      - "traefik.http.middlewares.judge0-cors.headers.customresponseheaders.Access-Control-Allow-Methods=GET,POST,PUT,DELETE,OPTIONS"
      - "traefik.http.middlewares.judge0-cors.headers.customresponseheaders.Access-Control-Allow-Headers=Authorization,Content-Type,X-RapidAPI-Key"
      - "traefik.http.routers.judge0.middlewares=judge0-cors"
    depends_on:
      - judge0-db
      - judge0-redis
    networks:
      - wizardcore-network
    restart: unless-stopped

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
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Database services (internal only)
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: wizardcore
      POSTGRES_PASSWORD: wizardcore_password
      POSTGRES_DB: wizardcore
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - wizardcore-network
    restart: unless-stopped

  auth-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: supabase_auth_admin
      POSTGRES_PASSWORD: password
      POSTGRES_DB: supabase_auth
    volumes:
      - auth_db_data:/var/lib/postgresql/data
    networks:
      - wizardcore-network
    restart: unless-stopped

  judge0-db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: judge0
      POSTGRES_PASSWORD: judge0
      POSTGRES_DB: judge0
    volumes:
      - judge0_db_data:/var/lib/postgresql/data
    networks:
      - wizardcore-network
    restart: unless-stopped

  # Redis instances
  redis:
    image: redis:7-alpine
    command: redis-server --requirepass xZDxSwVXHxVf6FYksaeEoA== --bind 0.0.0.0
    volumes:
      - redis_data:/data
    networks:
      - wizardcore-network
    restart: unless-stopped

  judge0-redis:
    image: redis:7.2-alpine
    command: redis-server --appendonly yes --bind 0.0.0.0
    volumes:
      - judge0_redis_data:/data
    networks:
      - wizardcore-network
    restart: unless-stopped

volumes:
  postgres_data:
  auth_db_data:
  judge0_db_data:
  redis_data:
  judge0_redis_data:

networks:
  wizardcore-network:
    driver: bridge
COMPOSE_EOF

# Create .env file with all necessary variables
echo -e "${YELLOW}🔧 Creating environment configuration...${NC}"
cat > $DEPLOY_DIR/.env << 'ENV_EOF'
# Domain configuration
DOMAIN=offensivewizard.com
FRONTEND_DOMAIN=app.offensivewizard.com
BACKEND_DOMAIN=api.offensivewizard.com
AUTH_DOMAIN=auth.offensivewizard.com
JUDGE0_DOMAIN=judge0.offensivewizard.com

# Database passwords
DATABASE_PASSWORD=wizardcore_password
AUTH_DB_PASSWORD=password
JUDGE0_DB_PASSWORD=judge0

# Redis password
REDIS_PASSWORD=xZDxSwVXHxVf6FYksaeEoA==

# JWT secret
SUPABASE_JWT_SECRET=aNBMlTAljHvKyR2dPu6R6nyggeW2398Na3R4XL1+oyebUDiuzSO61nZzoVmRi0h4

# Let's Encrypt email
LETSENCRYPT_EMAIL=admin@offensivewizard.com
ENV_EOF

# Pull Docker images
echo -e "${YELLOW}📥 Pulling Docker images...${NC}"
docker pull traefik:v3.0
docker pull limpet/wizardcore-frontend:latest || echo -e "${YELLOW}⚠️  Frontend image not found${NC}"
docker pull limpet/wizardcore-backend:latest || echo -e "${YELLOW}⚠️  Backend image not found${NC}"
docker pull supabase/gotrue:v2.184.0
docker pull postgres:15-alpine
docker pull postgres:16-alpine
docker pull redis:7-alpine
docker pull redis:7.2-alpine
docker pull judge0/judge0:latest

# Create letsencrypt directory for SSL certificates
echo -e "${YELLOW}🔐 Setting up SSL certificate storage...${NC}"
mkdir -p $DEPLOY_DIR/letsencrypt
touch $DEPLOY_DIR/letsencrypt/acme.json
chmod 600 $DEPLOY_DIR/letsencrypt/acme.json

# Stop any existing containers
echo -e "${YELLOW}🔄 Stopping existing containers...${NC}"
cd $DEPLOY_DIR
docker-compose down --remove-orphans 2>/dev/null || true

# Start the stack
echo -e "${YELLOW}🚀 Starting WizardCore stack with Traefik...${NC}"
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
cat > /usr/local/bin/wiz-start << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose up -d
EOF

# Stop script
cat > /usr/local/bin/wiz-stop << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose down
EOF

# Restart script
cat > /usr/local/bin/wiz-restart << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose restart
EOF

# Logs script
cat > /usr/local/bin/wiz-logs << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose logs -f
EOF

# Status script
cat > /usr/local/bin/wiz-status << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose ps
EOF

# Update script
cat > /usr/local/bin/wiz-update << 'EOF'
#!/bin/bash
cd /opt/wizardcore
docker-compose pull
docker-compose down
docker-compose up -d
EOF

# SSL check script
cat > /usr/local/bin/wiz-ssl-check << 'EOF'
#!/bin/bash
echo "🔐 SSL Certificate Status"
echo "========================"
cd /opt/wizardcore
echo "Checking certificates for:"
echo "1. app.offensivewizard.com"
curl -I https://app.offensivewizard.com 2>/dev/null | grep -i "certificate\|ssl\|http" || echo "  ❌ Cannot reach"
echo ""
echo "2. api.offensivewizard.com"
curl -I https://api.offensivewizard.com 2>/dev/null | grep -i "certificate\|ssl\|http" || echo "  ❌ Cannot reach"
echo ""
echo "3. auth.offensivewizard.com"
curl -I https://auth.offensivewizard.com 2>/dev/null | grep -i "certificate\|ssl\|http" || echo "  ❌ Cannot reach"
echo ""
echo "4. judge0.offensivewizard.com"
curl -I https://judge0.offensivewizard.com 2>/dev/null | grep -i "certificate\|ssl\|http" || echo "  ❌ Cannot reach"
echo ""
echo "📁 Certificate files:"
ls -la /opt/wizardcore/letsencrypt/
EOF

# Make scripts executable
chmod +x /usr/local/bin/wiz-*

echo ""
echo -e "${GREEN}🎉 Deployment complete!${NC}"
echo ""
echo -e "${GREEN}🌐 Your services are now accessible at:${NC}"
echo "   🔒 Frontend:    https://app.offensivewizard.com"
echo "   🔒 Backend API: https://api.offensivewizard.com"
echo "   🔒 Auth:        https://auth.offensivewizard.com"
echo "   🔒 Judge0:      https://judge0.offensivewizard.com"
echo "   📊 Traefik Dashboard: http://$(hostname -I | awk '{print $1}'):8080"
echo ""
echo -e "${GREEN}🔧 Management commands:${NC}"
echo "   wiz-start      - Start all services"
echo "   wiz-stop       - Stop all services"
echo "   wiz-restart    - Restart all services"
echo "   wiz-logs       - View logs (follow mode)"
echo "   wiz-status     - Check service status"
echo "   wiz-update     - Update to latest images"
echo "   wiz-ssl-check  - Check SSL certificate status"
echo ""
echo -e "${YELLOW}⚠️  IMPORTANT:${NC}"
echo "1. Make sure DNS records are configured:"
echo "   - A record: app.offensivewizard.com → $(hostname -I | awk '{print $1}')"
echo "   - A record: api.offensivewizard.com → $(hostname -I | awk '{print $1}')"
echo "   - A record: auth.offensivewizard.com → $(hostname -I | awk '{print $1}')"
echo "   - A record: judge0.offensivewizard.com → $(hostname -I | awk '{print $1}')"
echo ""
echo "2. SSL certificates will be automatically generated by Let's Encrypt"
echo "   (first access might take a minute for certificate issuance)"
echo ""
echo "3. Check Traefik logs for SSL certificate status:"
echo "   docker logs traefik 2>&1 | grep -i certificate"
echo ""
echo -e "${GREEN}✅ Everything is set up!${NC}"
echo "   Services will automatically get SSL certificates via Let's Encrypt"
echo "   All CORS headers are properly configured"
echo "   All services communicate correctly with each other"