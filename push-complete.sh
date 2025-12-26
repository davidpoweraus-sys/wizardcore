#!/bin/bash
# Push complete WizardCore deployment to server
# Usage: ./push-complete.sh your-server-ip

set -e

SERVER_IP="$1"
if [ -z "$SERVER_IP" ]; then
    echo "Usage: $0 <server-ip>"
    echo "Example: $0 192.168.1.100"
    echo "Example: $0 offensivewizard.com"
    exit 1
fi

echo "🚀 Pushing Complete WizardCore Deployment to $SERVER_IP"
echo "======================================================="

# Test SSH connection
echo "🔌 Testing SSH connection..."
if ! ssh -o ConnectTimeout=5 root@$SERVER_IP "echo 'SSH connection successful'" &>/dev/null; then
    echo "❌ Cannot connect to $SERVER_IP via SSH"
    echo "   Make sure:"
    echo "   1. SSH is enabled on the server"
    echo "   2. You have root access"
    echo "   3. Firewall allows port 22"
    exit 1
fi

echo "✅ SSH connection successful"

# Create directory on server
echo "📁 Creating directory on server..."
ssh root@$SERVER_IP "mkdir -p /opt/wizardcore"

# Transfer the complete deployment script
echo "📤 Transferring deployment script..."
scp deploy-complete.sh root@$SERVER_IP:/opt/wizardcore/

# Make it executable
echo "🔧 Making script executable..."
ssh root@$SERVER_IP "chmod +x /opt/wizardcore/deploy-complete.sh"

echo ""
echo "✅ Files transferred successfully!"
echo ""
echo "📋 Next steps:"
echo "1. SSH into server:"
echo "   ssh root@$SERVER_IP"
echo ""
echo "2. Run the complete deployment:"
echo "   cd /opt/wizardcore"
echo "   ./deploy-complete.sh"
echo ""
echo "📝 What this script does:"
echo "   ✅ Installs Docker and Docker Compose if needed"
echo "   ✅ Sets up all services with proper configuration"
echo "   ✅ Configures Traefik reverse proxy"
echo "   ✅ Sets up automatic SSL certificates (Let's Encrypt)"
echo "   ✅ Configures all subdomains:"
echo "      - app.offensivewizard.com (frontend)"
echo "      - api.offensivewizard.com (backend)"
echo "      - auth.offensivewizard.com (authentication)"
echo "      - judge0.offensivewizard.com (code execution)"
echo "   ✅ Creates management commands (wiz-*)"
echo "   ✅ Sets up persistent volumes for databases"
echo ""
echo "⚠️  IMPORTANT DNS CONFIGURATION:"
echo "   Before running, make sure these DNS A records point to $SERVER_IP:"
echo "   - app.offensivewizard.com"
echo "   - api.offensivewizard.com"
echo "   - auth.offensivewizard.com"
echo "   - judge0.offensivewizard.com"
echo ""
echo "🔧 After deployment, you can manage with:"
echo "   wiz-status    - Check service status"
echo "   wiz-logs      - View logs"
echo "   wiz-restart   - Restart services"
echo "   wiz-ssl-check - Check SSL certificates"
echo ""
echo "🌐 Your services will be accessible at:"
echo "   https://app.offensivewizard.com"
echo "   https://api.offensivewizard.com"
echo "   https://auth.offensivewizard.com"
echo "   https://judge0.offensivewizard.com"