#!/bin/bash
# CapRover Deployment Script for WizardCore

set -e

echo "🚀 WizardCore CapRover Deployment"
echo "================================="

# Check if CapRover CLI is installed
if ! command -v caprover &> /dev/null; then
    echo "❌ CapRover CLI not installed. Installing..."
    npm install -g caprover
fi

# Configuration
APP_NAME="wizardcore"
DOMAIN="offensivewizard.com"
CAPROVER_URL="captain.your-domain.com"  # Change this to your CapRover URL
CAPROVER_PASSWORD=""  # Set your CapRover password

# Check if we have the password
if [ -z "$CAPROVER_PASSWORD" ]; then
    echo "⚠️  CapRover password not set. Please set CAPROVER_PASSWORD environment variable."
    echo "   export CAPROVER_PASSWORD='your-password'"
    exit 1
fi

echo "📋 Configuration:"
echo "  App Name: $APP_NAME"
echo "  Domain: $DOMAIN"
echo "  CapRover URL: $CAPROVER_URL"

# Login to CapRover
echo "🔐 Logging into CapRover..."
caprover login --caproverUrl "https://$CAPROVER_URL" --password "$CAPROVER_PASSWORD"

# Create app if it doesn't exist
echo "📦 Creating app '$APP_NAME'..."
caprover app create --appName "$APP_NAME" || echo "App may already exist"

# Set app domain
echo "🌐 Configuring domain '$DOMAIN'..."
caprover app addcustomdomain --appName "$APP_NAME" --customDomain "$DOMAIN"

# Enable HTTPS
echo "🔒 Enabling HTTPS..."
caprover app enablebasedomainssl --appName "$APP_NAME"

# Deploy using Docker Compose
echo "🐳 Deploying with Docker Compose..."
if [ -f "docker-compose.caprover.yml" ]; then
    echo "📄 Using docker-compose.caprover.yml"
    
    # Create a temporary directory for deployment
    TEMP_DIR=$(mktemp -d)
    cp docker-compose.caprover.yml "$TEMP_DIR/docker-compose.yml"
    cp -r init-scripts "$TEMP_DIR/" 2>/dev/null || true
    
    # Create captain-definition file
    cat > "$TEMP_DIR/captain-definition" << EOF
{
  "schemaVersion": 2,
  "dockerfileLines": [
    "FROM alpine:latest",
    "RUN apk add --no-cache docker docker-compose",
    "COPY . .",
    "CMD [\"docker-compose\", \"up\"]"
  ]
}
EOF
    
    # Deploy to CapRover
    cd "$TEMP_DIR"
    caprover deploy --appName "$APP_NAME" --tarPath .
    
    # Cleanup
    cd -
    rm -rf "$TEMP_DIR"
else
    echo "❌ docker-compose.caprover.yml not found"
    echo "📥 Deploying from GitHub instead..."
    caprover deploy --appName "$APP_NAME" --imageName "limpet/wizardcore-frontend:latest"
fi

# Configure environment variables
echo "🔧 Setting environment variables..."
caprover app setenv --appName "$APP_NAME" -e "NEXT_PUBLIC_SUPABASE_URL=https://$DOMAIN/supabase-proxy"
caprover app setenv --appName "$APP_NAME" -e "NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps="
caprover app setenv --appName "$APP_NAME" -e "NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.$DOMAIN"
caprover app setenv --appName "$APP_NAME" -e "NEXT_PUBLIC_BACKEND_URL=https://$DOMAIN/api"

# Configure nginx for proxy routing
echo "🔄 Configuring nginx proxy routes..."
cat > /tmp/nginx-custom.conf << 'NGINX_EOF'
location /api/ {
    proxy_pass http://backend:8080/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /supabase-proxy/ {
    proxy_pass http://supabase-auth:9999/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    
    # Strip /supabase-proxy prefix
    rewrite ^/supabase-proxy/(.*)$ /$1 break;
}
NGINX_EOF

# Upload nginx config (this requires CapRover API)
echo "📤 Uploading nginx configuration..."
# Note: This step may require manual configuration in CapRover dashboard

echo ""
echo "🎉 Deployment initiated!"
echo ""
echo "📋 Next steps:"
echo "1. Go to CapRover dashboard: https://$CAPROVER_URL"
echo "2. Check app '$APP_NAME' status"
echo "3. Configure custom nginx in App Configs → HTTP Settings"
echo "4. Add the nginx configuration from above"
echo ""
echo "🔍 Check deployment status:"
echo "   caprover app status --appName $APP_NAME"
echo ""
echo "📊 View logs:"
echo "   caprover app logs --appName $APP_NAME"
echo ""
echo "🌐 Test your deployment:"
echo "   curl https://$DOMAIN"
echo "   curl https://$DOMAIN/api/health"