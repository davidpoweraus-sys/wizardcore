#!/bin/bash

echo "Testing CORS fix for Supabase Auth..."
echo "====================================="

# Test 1: Check if local Supabase Auth is running
echo -e "\n1. Testing local Supabase Auth (localhost:9999)..."
curl -s -o /dev/null -w "%{http_code}" -X OPTIONS "http://localhost:9999/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST"

echo " - OPTIONS request status"

# Test 2: Check headers
echo -e "\n2. Checking CORS headers from local Supabase Auth..."
curl -s -I -X OPTIONS "http://localhost:9999/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" | grep -i "access-control"

# Test 3: Test with actual POST (simulated)
echo -e "\n3. Testing if we can reach the auth endpoint..."
curl -s -o /dev/null -w "%{http_code}" -X GET "http://localhost:9999/auth/v1/health"
echo " - Health check status"

# Test 4: Check nginx config syntax
echo -e "\n4. Checking nginx configuration..."
if command -v nginx &> /dev/null; then
  nginx -t -c "$(pwd)/nginx-config/nginx.conf" 2>&1 | grep -E "(successful|failed)"
else
  echo "nginx not installed, skipping config test"
fi

echo -e "\n====================================="
echo "INSTRUCTIONS:"
echo "1. If testing locally, add to /etc/hosts:"
echo "   127.0.0.1 auth.offensivewizard.com"
echo "2. Restart services:"
echo "   docker-compose -f docker-compose.local.yml down"
echo "   docker-compose -f docker-compose.local.yml up -d"
echo "3. Test with:"
echo "   curl -v -X OPTIONS 'https://auth.offensivewizard.com/auth/v1/signup' \\"
echo "     -H 'Origin: https://offensivewizard.com' \\"
echo "     -H 'Access-Control-Request-Method: POST' \\"
echo "     --insecure"
echo ""
echo "The fix includes:"
echo "- Updated nginx CORS headers with 'always' flag"
echo "- Added missing headers (apikey, x-client-info, x-supabase-api-version)"
echo "- Set GoTrue CORS to allow all origins (*)"
echo "- Fixed nginx header inheritance issue"