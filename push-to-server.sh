#!/bin/bash
# Push WizardCore files to server
# Usage: ./push-to-server.sh your-server-ip

set -e

SERVER_IP="$1"
if [ -z "$SERVER_IP" ]; then
    echo "Usage: $0 <server-ip>"
    echo "Example: $0 192.168.1.100"
    exit 1
fi

echo "🚀 Pushing WizardCore files to $SERVER_IP"
echo "=========================================="

# Create directory on server
echo "📁 Creating directory on server..."
ssh root@$SERVER_IP "mkdir -p /opt/wizardcore"

# Transfer essential files
echo "📤 Transferring files..."
scp simple-docker-compose.yml root@$SERVER_IP:/opt/wizardcore/
scp .env.example root@$SERVER_IP:/opt/wizardcore/
scp setup-env.sh root@$SERVER_IP:/opt/wizardcore/
scp one-line-deploy.sh root@$SERVER_IP:/opt/wizardcore/
scp deploy-ssh.sh root@$SERVER_IP:/opt/wizardcore/
scp setup-nginx.sh root@$SERVER_IP:/opt/wizardcore/
scp SSH-DEPLOYMENT-GUIDE.md root@$SERVER_IP:/opt/wizardcore/
scp AUTH-PROXY-FIX.md root@$SERVER_IP:/opt/wizardcore/
scp DEPLOY-AUTH-PROXY-FIX.md root@$SERVER_IP:/opt/wizardcore/

# Make scripts executable on server
echo "🔧 Making scripts executable..."
ssh root@$SERVER_IP "chmod +x /opt/wizardcore/*.sh"

echo ""
echo "✅ Files transferred successfully!"
echo ""
echo "📋 Next steps:"
echo "1. SSH into server:"
echo "   ssh root@$SERVER_IP"
echo ""
echo "2. Setup environment (production-ready values included!):"
echo "   cd /opt/wizardcore"
echo "   ./setup-env.sh"
echo ""
echo "3. Run deployment:"
echo "   ./one-line-deploy.sh"
echo ""
echo "4. Or for full setup with Nginx:"
echo "   ./deploy-ssh.sh"
echo "   ./setup-nginx.sh"
echo ""
echo "📁 Files on server:"
echo "   /opt/wizardcore/.env.example (production-ready!)"
echo "   /opt/wizardcore/setup-env.sh"
echo "   /opt/wizardcore/simple-docker-compose.yml"
echo "   /opt/wizardcore/one-line-deploy.sh"
echo "   /opt/wizardcore/deploy-ssh.sh"
echo "   /opt/wizardcore/setup-nginx.sh"
echo "   /opt/wizardcore/AUTH-PROXY-FIX.md"
echo "   /opt/wizardcore/DEPLOY-AUTH-PROXY-FIX.md"