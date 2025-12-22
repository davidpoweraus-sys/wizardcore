#!/bin/bash

# Detailed CORS Testing Script for WizardCore Auth
# This script tests all CORS endpoints and configurations

echo "=================================="
echo "🧙 WizardCore CORS Diagnostic Test"
echo "=================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test configuration
AUTH_URL="https://auth.offensivewizard.com"
ORIGIN="https://offensivewizard.com"

echo -e "${BLUE}Testing CORS configuration...${NC}"
echo ""

# Test 1: OPTIONS preflight for /auth/v1/signup
echo -e "${YELLOW}Test 1: OPTIONS preflight to /auth/v1/signup${NC}"
echo "Request: OPTIONS ${AUTH_URL}/auth/v1/signup"
echo "Origin: ${ORIGIN}"
echo ""
RESPONSE=$(curl -s -i -X OPTIONS "${AUTH_URL}/auth/v1/signup" \
  -H "Origin: ${ORIGIN}" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: apikey,authorization,content-type,x-client-info,x-supabase-api-version" \
  2>&1)

if echo "$RESPONSE" | grep -q "HTTP.*204\|HTTP.*200"; then
    echo -e "${GREEN}✓ Status: OK${NC}"
else
    echo -e "${RED}✗ Status: Failed${NC}"
fi

if echo "$RESPONSE" | grep -qi "access-control-allow-origin.*${ORIGIN}"; then
    echo -e "${GREEN}✓ Access-Control-Allow-Origin: ${ORIGIN}${NC}"
else
    echo -e "${RED}✗ Missing or incorrect Access-Control-Allow-Origin${NC}"
fi

if echo "$RESPONSE" | grep -qi "access-control-allow-credentials.*true"; then
    echo -e "${GREEN}✓ Access-Control-Allow-Credentials: true${NC}"
else
    echo -e "${RED}✗ Missing Access-Control-Allow-Credentials${NC}"
fi

if echo "$RESPONSE" | grep -qi "access-control-allow-methods"; then
    METHODS=$(echo "$RESPONSE" | grep -i "access-control-allow-methods" | head -1)
    echo -e "${GREEN}✓ Allowed Methods: ${METHODS#*: }${NC}"
else
    echo -e "${RED}✗ Missing Access-Control-Allow-Methods${NC}"
fi

if echo "$RESPONSE" | grep -qi "access-control-allow-headers"; then
    HEADERS=$(echo "$RESPONSE" | grep -i "access-control-allow-headers" | head -1)
    echo -e "${GREEN}✓ Allowed Headers: ${HEADERS#*: }${NC}"
else
    echo -e "${RED}✗ Missing Access-Control-Allow-Headers${NC}"
fi

echo ""
echo "Raw Response:"
echo "---"
echo "$RESPONSE" | head -20
echo "---"
echo ""

# Test 2: GET to /health endpoint
echo -e "${YELLOW}Test 2: GET to /health endpoint${NC}"
HEALTH_RESPONSE=$(curl -s -i -X GET "${AUTH_URL}/health" \
  -H "Origin: ${ORIGIN}" \
  2>&1)

if echo "$HEALTH_RESPONSE" | grep -q "HTTP.*200"; then
    echo -e "${GREEN}✓ Health check passed${NC}"
else
    echo -e "${RED}✗ Health check failed${NC}"
fi
echo ""

# Test 3: Check if services are running
echo -e "${YELLOW}Test 3: Check local services${NC}"
if docker ps | grep -q "supabase-auth"; then
    echo -e "${GREEN}✓ Supabase Auth container is running${NC}"
else
    echo -e "${RED}✗ Supabase Auth container is NOT running${NC}"
fi

if docker ps | grep -q "nginx"; then
    echo -e "${GREEN}✓ Nginx container is running${NC}"
else
    echo -e "${RED}✗ Nginx container is NOT running${NC}"
fi
echo ""

# Test 4: Check nginx configuration
echo -e "${YELLOW}Test 4: Nginx Configuration${NC}"
if [ -f "nginx-config/nginx.conf" ]; then
    echo -e "${GREEN}✓ nginx.conf exists${NC}"
    if grep -q "map \$http_origin" nginx-config/nginx.conf; then
        echo -e "${GREEN}✓ Origin map is configured${NC}"
    else
        echo -e "${YELLOW}⚠ Origin map not found (may be using old config)${NC}"
    fi
else
    echo -e "${RED}✗ nginx.conf not found${NC}"
fi
echo ""

# Test 5: Environment variables
echo -e "${YELLOW}Test 5: Environment Variables${NC}"
if grep -q "NEXT_PUBLIC_SUPABASE_URL" docker-compose.prod.yml; then
    SUPABASE_URL=$(grep "NEXT_PUBLIC_SUPABASE_URL" docker-compose.prod.yml | head -1 | cut -d'=' -f2)
    echo -e "${GREEN}✓ NEXT_PUBLIC_SUPABASE_URL: ${SUPABASE_URL}${NC}"
else
    echo -e "${RED}✗ NEXT_PUBLIC_SUPABASE_URL not found${NC}"
fi

if grep -q "GOTRUE_CORS_ALLOWED_ORIGINS" docker-compose.prod.yml; then
    CORS_ORIGINS=$(grep "GOTRUE_CORS_ALLOWED_ORIGINS" docker-compose.prod.yml | head -1 | cut -d':' -f2)
    echo -e "${GREEN}✓ GOTRUE_CORS_ALLOWED_ORIGINS: ${CORS_ORIGINS}${NC}"
else
    echo -e "${RED}✗ GOTRUE_CORS_ALLOWED_ORIGINS not found${NC}"
fi
echo ""

# Test 6: SSL Certificate
echo -e "${YELLOW}Test 6: SSL Certificate${NC}"
CERT_CHECK=$(echo | openssl s_client -servername auth.offensivewizard.com -connect auth.offensivewizard.com:443 2>/dev/null | openssl x509 -noout -dates 2>/dev/null)
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ SSL Certificate is valid${NC}"
    echo "$CERT_CHECK"
else
    echo -e "${RED}✗ SSL Certificate issue or cannot connect${NC}"
fi
echo ""

# Summary
echo "=================================="
echo -e "${BLUE}🎯 Summary & Recommendations${NC}"
echo "=================================="
echo ""
echo "Common Issues:"
echo "1. If preflight fails: Check nginx CORS headers"
echo "2. If origin mismatch: Verify GOTRUE_CORS_ALLOWED_ORIGINS includes ${ORIGIN}"
echo "3. If DNS issues: Check /etc/hosts or DNS records"
echo "4. If SSL errors: Verify certificates are valid and trusted"
echo ""
echo "Next Steps:"
echo "- Check browser console for detailed error messages"
echo "- Verify Network tab shows the actual request/response headers"
echo "- Check docker logs: docker logs supabase-auth"
echo "- Check nginx logs: docker logs <nginx-container>"
echo ""
