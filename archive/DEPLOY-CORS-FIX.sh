#!/bin/bash

# Quick deployment script for CORS fix
# This script helps you deploy the Traefik CORS configuration

set -e

echo "=================================="
echo "🧙 WizardCore CORS Fix Deployment"
echo "=================================="
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Step 1: Checking git status...${NC}"
git status --short

echo ""
echo -e "${YELLOW}Step 2: Adding changed files...${NC}"
git add docker-compose.prod.yml
git add app/\(auth\)/register/page.tsx
git add COOLIFY-CORS-FIX.md
git add ACTUAL-CORS-FIX.md
git add DEPLOY-CORS-FIX.sh

echo -e "${GREEN}✓ Files staged${NC}"
echo ""

echo -e "${YELLOW}Step 3: Creating commit...${NC}"
git commit -m "Fix CORS by adding Traefik labels for Coolify reverse proxy

- Add Traefik CORS middleware labels to supabase-auth service
- Add detailed console logging to registration page
- Document Coolify-specific CORS configuration
- Provide deployment and troubleshooting guides"

echo -e "${GREEN}✓ Commit created${NC}"
echo ""

echo -e "${YELLOW}Step 4: Pushing to remote...${NC}"
git push

echo -e "${GREEN}✓ Pushed to remote${NC}"
echo ""

echo "=================================="
echo -e "${GREEN}✓ Code deployed to repository!${NC}"
echo "=================================="
echo ""

echo "📋 Next Steps (Manual):"
echo ""
echo "1. Go to Coolify Dashboard"
echo "   → Redeploy your WizardCore project"
echo ""
echo "2. Add router middleware label:"
echo "   → Go to supabase-auth service → Labels"
echo "   → Add: traefik.http.routers.<router-name>.middlewares = auth-cors"
echo "   → (Find router name in Traefik dashboard or logs)"
echo ""
echo "3. Restart supabase-auth service"
echo ""
echo "4. Test registration:"
echo "   → Open https://offensivewizard.com/register"
echo "   → Open browser DevTools (F12) → Console"
echo "   → Look for emoji logs: 🚀 📧 ✅"
echo ""
echo "📖 For detailed instructions, see:"
echo "   - ACTUAL-CORS-FIX.md (quick reference)"
echo "   - COOLIFY-CORS-FIX.md (detailed guide)"
echo ""
echo "🧪 Test CORS manually:"
echo "   curl -v -X OPTIONS \"https://auth.offensivewizard.com/auth/v1/signup\" \\"
echo "     -H \"Origin: https://offensivewizard.com\" \\"
echo "     -H \"Access-Control-Request-Method: POST\""
echo ""
