#!/bin/bash
set -e

# Test Supabase Auth connectivity and CORS headers
# Usage: ./scripts/test-auth.sh [AUTH_URL]
# Default AUTH_URL: http://localhost:9999

AUTH_URL="${1:-http://localhost:9999}"
echo "Testing Supabase Auth at $AUTH_URL"

# Check if service is reachable
if ! curl -s --fail --max-time 5 "$AUTH_URL/health"; then
    echo "ERROR: Cannot reach Supabase Auth health endpoint."
    echo "Make sure the service is running and accessible."
    exit 1
fi

echo "Health check passed."

# Test CORS headers on OPTIONS request (preflight)
echo "Testing CORS preflight..."
if ! curl -s --max-time 5 -X OPTIONS "$AUTH_URL/auth/v1/signup" \
    -H "Origin: https://offensivewizard.com" \
    -H "Access-Control-Request-Method: POST" \
    -H "Access-Control-Request-Headers: Authorization, Content-Type" \
    --verbose 2>&1 | grep -q "access-control-allow-origin"; then
    echo "WARNING: CORS headers may not be properly configured."
    echo "Check GOTRUE_CORS_ALLOWED_ORIGINS environment variable."
else
    echo "CORS preflight successful."
fi

# Test actual POST request (without body)
echo "Testing POST with Origin header..."
curl -s --max-time 5 -X POST "$AUTH_URL/auth/v1/signup" \
    -H "Origin: https://offensivewizard.com" \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"testpassword"}' \
    --verbose 2>&1 | head -20

echo "Test completed. If you see CORS errors in browser, ensure the frontend URL is allowed."