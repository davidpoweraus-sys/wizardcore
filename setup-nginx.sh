#!/bin/bash
# Nginx Reverse Proxy Setup for WizardCore
# Run this after deploy-ssh.sh to configure Nginx

set -e

echo "🔧 Setting up Nginx reverse proxy for WizardCore"
echo "================================================"

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

# Install Nginx if not present
if ! command -v nginx &> /dev/null; then
    echo -e "${YELLOW}📦 Installing Nginx...${NC}"
    apt-get update
    apt-get install -y nginx
    echo -e "${GREEN}✅ Nginx installed${NC}"
fi

# Create Nginx configuration
echo -e "${YELLOW}📝 Creating Nginx configuration...${NC}"

cat > /etc/nginx/sites-available/wizardcore << 'NGINX_EOF'
# WizardCore Reverse Proxy Configuration
# Domain: offensivewizard.com

# Frontend - app.offensivewizard.com
server {
    listen 80;
    listen [::]:80;
    server_name app.offensivewizard.com;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        
        # Increase timeout for WebSocket connections
        proxy_read_timeout 86400;
        proxy_send_timeout 86400;
    }
}

# Backend API - api.offensivewizard.com
server {
    listen 80;
    listen [::]:80;
    server_name api.offensivewizard.com;
    
    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers
        add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type' always;
        add_header 'Access-Control-Allow-Credentials' 'true' always;
        
        # Handle preflight requests
        if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com';
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS';
            add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type';
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain charset=UTF-8';
            add_header 'Content-Length' 0;
            return 204;
        }
    }
}

# Supabase Auth - auth.offensivewizard.com
server {
    listen 80;
    listen [::]:80;
    server_name auth.offensivewizard.com;
    
    location / {
        proxy_pass http://localhost:9999;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers
        add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, apikey, x-client-info, x-supabase-api-version' always;
        add_header 'Access-Control-Allow-Credentials' 'true' always;
        add_header 'Access-Control-Expose-Headers' 'X-Total-Count' always;
        
        # Handle preflight requests
        if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com';
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS';
            add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, apikey, x-client-info, x-supabase-api-version';
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain charset=UTF-8';
            add_header 'Content-Length' 0;
            return 204;
        }
    }
}

# Judge0 API - judge0.offensivewizard.com
server {
    listen 80;
    listen [::]:80;
    server_name judge0.offensivewizard.com;
    
    location / {
        proxy_pass http://localhost:2358;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # CORS headers
        add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com' always;
        add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS' always;
        add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-RapidAPI-Key' always;
        
        # Handle preflight requests
        if ($request_method = 'OPTIONS') {
            add_header 'Access-Control-Allow-Origin' 'https://app.offensivewizard.com';
            add_header 'Access-Control-Allow-Methods' 'GET, POST, PUT, DELETE, OPTIONS';
            add_header 'Access-Control-Allow-Headers' 'Authorization, Content-Type, X-RapidAPI-Key';
            add_header 'Access-Control-Max-Age' 1728000;
            add_header 'Content-Type' 'text/plain charset=UTF-8';
            add_header 'Content-Length' 0;
            return 204;
        }
    }
}
NGINX_EOF

# Enable the site
echo -e "${YELLOW}🔗 Enabling Nginx site...${NC}"
ln -sf /etc/nginx/sites-available/wizardcore /etc/nginx/sites-enabled/
rm -f /etc/nginx/sites-enabled/default 2>/dev/null || true

# Test Nginx configuration
echo -e "${YELLOW}🧪 Testing Nginx configuration...${NC}"
if nginx -t; then
    echo -e "${GREEN}✅ Nginx configuration is valid${NC}"
else
    echo -e "${RED}❌ Nginx configuration has errors${NC}"
    exit 1
fi

# Restart Nginx
echo -e "${YELLOW}🔄 Restarting Nginx...${NC}"
systemctl restart nginx

# Set up SSL with Certbot (optional)
echo ""
echo -e "${YELLOW}🔐 SSL Certificate Setup (Optional)${NC}"
echo "To set up SSL certificates with Let's Encrypt, run:"
echo ""
echo "1. Install Certbot:"
echo "   apt-get install -y certbot python3-certbot-nginx"
echo ""
echo "2. Get certificates for all subdomains:"
echo "   certbot --nginx -d app.offensivewizard.com -d api.offensivewizard.com -d auth.offensivewizard.com -d judge0.offensivewizard.com"
echo ""
echo "3. Certbot will automatically update your Nginx configuration"
echo ""

# Create a simple SSL setup script
cat > /opt/wizardcore/setup-ssl.sh << 'SSL_EOF'
#!/bin/bash
# SSL Setup Script for WizardCore

echo "🔐 Setting up SSL certificates with Let's Encrypt"
echo "================================================="

# Install Certbot
apt-get update
apt-get install -y certbot python3-certbot-nginx

# Get certificates
certbot --nginx \
  -d app.offensivewizard.com \
  -d api.offensivewizard.com \
  -d auth.offensivewizard.com \
  -d judge0.offensivewizard.com \
  --non-interactive \
  --agree-tos \
  --email admin@offensivewizard.com \
  --redirect

echo ""
echo "✅ SSL certificates installed!"
echo "📋 Your sites are now accessible via HTTPS"
SSL_EOF

chmod +x /opt/wizardcore/setup-ssl.sh

echo -e "${GREEN}✅ Nginx reverse proxy configured!${NC}"
echo ""
echo -e "${GREEN}📋 Your services are now accessible at:${NC}"
echo "   HTTP → app.offensivewizard.com"
echo "   HTTP → api.offensivewizard.com"
echo "   HTTP → auth.offensivewizard.com"
echo "   HTTP → judge0.offensivewizard.com"
echo ""
echo -e "${YELLOW}⚠️  Important:${NC}"
echo "1. Make sure DNS records point to your server IP:"
echo "   - A record: app.offensivewizard.com"
echo "   - A record: api.offensivewizard.com"
echo "   - A record: auth.offensivewizard.com"
echo "   - A record: judge0.offensivewizard.com"
echo ""
echo "2. Run SSL setup when DNS is configured:"
echo "   cd /opt/wizardcore && ./setup-ssl.sh"
echo ""
echo "3. Check Nginx status: systemctl status nginx"
echo "4. Check WizardCore services: wizardcore-status"