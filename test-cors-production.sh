#!/bin/bash
echo "Testing CORS configuration for production deployment..."
echo "======================================================"

echo "1. Testing if auth.offensivewizard.com is accessible..."
curl -s -o /dev/null -w "HTTP Status: %{http_code}\n" -I "https://auth.offensivewizard.com/auth/v1/health" --insecure 2>&1 | head -1

echo ""
echo "2. Testing OPTIONS preflight request..."
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup?redirect_to=https%3A%2F%2Foffensivewizard.com%2Fauth%2Fcallback" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: apikey,authorization,content-type,x-client-info,x-supabase-api-version" \
  --insecure 2>&1 | grep -E "(< HTTP|< Access-Control|OPTIONS)"

echo ""
echo "======================================================"
echo "If you see HTTP 200 or 204 with proper CORS headers, the fix works."
echo "If you see HTTP 4xx/5xx or missing CORS headers, check:"
echo "1. Is supabase-auth service running in Coolify?"
echo "2. Is Coolify routing auth.offensivewizard.com correctly?"
echo "3. Check Coolify application settings for CORS configuration"
echo "4. Check Coolify SSL certificate configuration"