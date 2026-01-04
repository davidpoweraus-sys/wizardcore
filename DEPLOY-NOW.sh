#!/bin/bash

# DEPLOY-NOW.sh
# Immediate deployment of CORS fix for login issue

set -e

echo "🚀 DEPLOYING CORS FIX TO PRODUCTION"
echo "====================================="

echo "✅ Docker image already built and pushed:"
echo "   limpet/wizardcore-frontend:cors-fix-urgent"
echo ""

echo "📋 Fixes included in this image:"
echo "   1. CORS fix for same-origin requests in auth proxy"
echo "   2. CORS fix for same-origin requests in backend proxy"
echo "   3. Session refresh awareness in middleware"
echo "   4. Fixed Redis image (redis:7-alpine instead of 7.2-alpine)"
echo ""

echo "🔧 Deployment options:"
echo ""
echo "Option A: Update Dokploy deployment"
echo "-----------------------------------"
echo "1. Go to Dokploy Dashboard → WizardCore Application"
echo "2. Click 'Edit' or 'Settings'"
echo "3. Update image tag to: limpet/wizardcore-frontend:cors-fix-urgent"
echo "4. Click 'Save' and 'Redeploy'"
echo ""

echo "Option B: Deploy with Docker Compose"
echo "-------------------------------------"
echo "1. Update docker-compose.yml (already done)"
echo "2. Run: docker-compose pull frontend"
echo "3. Run: docker-compose up -d frontend"
echo ""

echo "Option C: Direct Docker commands"
echo "---------------------------------"
echo "1. Stop current frontend: docker stop [frontend-container-name]"
echo "2. Remove: docker rm [frontend-container-name]"
echo "3. Run new image:"
echo "   docker run -d \\"
echo "     --name wizardcore-frontend \\"
echo "     --network wizardcore_default \\"
echo "     -p 3001:3000 \\"
echo "     -e NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth \\"
echo "     -e NEXT_PUBLIC_SUPABASE_ANON_KEY=[your-key] \\"
echo "     -e NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com \\"
echo "     -e SUPABASE_INTERNAL_URL=http://supabase-auth:9999 \\"
echo "     -e GOTRUE_URL=http://supabase-auth:9999 \\"
echo "     -e BACKEND_URL=http://backend:8080 \\"
echo "     limpet/wizardcore-frontend:cors-fix-urgent"
echo ""

echo "🔍 Verification after deployment:"
echo "1. Clear browser cache and cookies"
echo "2. Log in at https://app.offensivewizard.com/login"
echo "3. Check browser console:"
echo "   - Should see: '🔍 Middleware session-refresh-fix-20260104-1159 executing'"
echo "   - /api/auth/auth/v1/user should return 200 (not 403)"
echo "   - User data should load (not 'guest' state)"
echo ""

echo "📞 If still having issues:"
echo "1. Check logs: docker logs [frontend-container] | grep -A5 -B5 'validateOrigin'"
echo "2. Verify CORS fix: docker exec [frontend-container] cat /app/app/api/auth/[...path]/route.ts | grep -A2 'function validateOrigin'"
echo ""

echo "🎯 The CORS fix resolves:"
echo "   - 403 'Origin not allowed' errors"
echo "   - 'Guest user' issue after login"
echo "   - Backend API calls failing with CORS errors"
echo ""

echo "🚨 URGENT: Deploy now to fix production login issue!"