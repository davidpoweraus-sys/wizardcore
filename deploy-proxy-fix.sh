#!/bin/bash

# Deploy Next.js Proxy Solution for CORS
set -e

echo "=================================="
echo "🧙 Deploy Next.js Proxy CORS Fix"
echo "=================================="
echo ""

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}This will bypass Coolify CORS issues by routing through Next.js${NC}"
echo ""

echo -e "${YELLOW}Step 1: Checking files...${NC}"
if [ -f "app/api/supabase-proxy/[...path]/route.ts" ]; then
    echo -e "${GREEN}✓ Proxy route exists${NC}"
else
    echo -e "${RED}✗ Proxy route missing!${NC}"
    exit 1
fi

echo ""
echo -e "${YELLOW}Step 2: Staging changes...${NC}"
git add app/api/supabase-proxy/
git add docker-compose.prod.yml
git add NEXTJS-PROXY-SOLUTION.md
git add deploy-proxy-fix.sh

echo -e "${GREEN}✓ Files staged${NC}"
echo ""

echo -e "${YELLOW}Step 3: Committing...${NC}"
git commit -m "Add Next.js proxy to bypass Coolify CORS issues

- Create /api/supabase-proxy route to handle Supabase requests
- Update NEXT_PUBLIC_SUPABASE_URL to use proxy
- Add SUPABASE_INTERNAL_URL for internal Docker network
- Bypass Coolify Traefik CORS limitations with server-side proxy"

echo -e "${GREEN}✓ Committed${NC}"
echo ""

echo -e "${YELLOW}Step 4: Pushing to remote...${NC}"
git push

echo -e "${GREEN}✓ Pushed!${NC}"
echo ""

echo "=================================="
echo -e "${GREEN}✓ Deployment Complete!${NC}"
echo "=================================="
echo ""

echo "📋 Next Steps:"
echo ""
echo "1. Coolify will auto-deploy (or click Redeploy in UI)"
echo ""
echo "2. Wait for build to complete (~2-5 minutes)"
echo ""
echo "3. Test registration:"
echo "   → https://offensivewizard.com/register"
echo "   → Open DevTools (F12) → Console"
echo "   → Look for: 🚀 Registration started"
echo ""
echo "4. Verify in Network tab:"
echo "   → Request goes to /api/supabase-proxy/auth/v1/signup"
echo "   → Status: 200 OK"
echo "   → NO CORS errors!"
echo ""
echo "5. Check Next.js logs in Coolify:"
echo "   → Should see: 🔄 Proxying request to..."
echo ""
echo "📖 See NEXTJS-PROXY-SOLUTION.md for details"
echo ""
