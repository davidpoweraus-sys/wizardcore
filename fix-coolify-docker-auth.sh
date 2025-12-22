#!/bin/bash

# Script to fix Docker authentication cache on Coolify server
# Run this ON YOUR COOLIFY SERVER via SSH

echo "🔧 Fixing Docker authentication cache for GHCR..."
echo ""

# Step 1: Logout from GHCR to clear cached credentials
echo "1️⃣ Logging out from ghcr.io to clear cache..."
docker logout ghcr.io 2>/dev/null || echo "Already logged out"

# Step 2: Remove any cached authentication
echo "2️⃣ Clearing Docker auth cache..."
rm -f ~/.docker/config.json.lock 2>/dev/null || true

# Step 3: Try pulling the public image to verify access
echo "3️⃣ Testing pull from ghcr.io (should work without credentials)..."
if docker pull ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest; then
    echo "✅ SUCCESS! Image pulled successfully"
    echo ""
    echo "🎯 Next step: Go to Coolify and click 'Redeploy'"
else
    echo "❌ Pull failed. The package might not be public yet."
    echo ""
    echo "📝 Check package visibility:"
    echo "   https://github.com/davidpoweraus-sys/wizardcore/pkgs/container/wizardcore-frontend"
    echo ""
    echo "   Make sure it says 'Public' not 'Private'"
fi

echo ""
echo "✅ Done!"
