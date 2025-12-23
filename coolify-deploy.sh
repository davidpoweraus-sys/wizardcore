#!/bin/bash

# One-liner deployment script for Coolify servers
# Downloads latest release from GitHub and loads images

set -e

PACKAGE_TYPE="${1:-cors-fix}"

echo "🚀 WizardCore Quick Deploy"
echo ""

# Download and execute the main deploy script
curl -fsSL https://raw.githubusercontent.com/davidpoweraus-sys/wizardcore/main/deploy-from-release.sh | bash -s "$PACKAGE_TYPE" latest

echo ""
echo "✅ Done! Now redeploy in Coolify UI or run:"
echo "   docker-compose -f /path/to/wizardcore/docker-compose.prod.yml up -d"
