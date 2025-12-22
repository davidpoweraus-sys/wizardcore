#!/bin/bash

# CORS Testing Script for WizardCore Authentication
# This script tests the CORS configuration between frontend and Supabase Auth

set -e

echo "=================================="
echo "WizardCore CORS Test Suite"
echo "=================================="
echo ""

# Configuration
FRONTEND_URL="${FRONTEND_URL:-https://offensivewizard.com}"
AUTH_URL="${AUTH_URL:-https://auth.offensivewizard.com}"

echo "Testing CORS between:"
echo "  Frontend: $FRONTEND_URL"
echo "  Auth API: $AUTH_URL"
echo ""

# Test 1: Health Check
echo "Test 1: Health Check"
echo "-------------------"
HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${AUTH_URL}/health" || echo "000")
if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "204" ]; then
    echo "✅ Auth service is responding (HTTP $HTTP_CODE)"
else
    echo "❌ Auth service not responding (HTTP $HTTP_CODE)"
    echo "   Check if supabase-auth container is running"
fi
echo ""

# Test 2: CORS Preflight for Signup
echo "Test 2: CORS Preflight (OPTIONS) for /auth/v1/signup"
echo "----------------------------------------------------"
CORS_RESPONSE=$(curl -s -i -X OPTIONS "${AUTH_URL}/auth/v1/signup" \
  -H "Origin: ${FRONTEND_URL}" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,apikey,x-client-info" \
  2>&1 || echo "FAILED")

if echo "$CORS_RESPONSE" | grep -q "access-control-allow-origin"; then
    echo "✅ CORS preflight successful"
    
    # Check specific headers
    if echo "$CORS_RESPONSE" | grep -qi "access-control-allow-origin.*${FRONTEND_URL}"; then
        echo "✅ Correct origin: $FRONTEND_URL"
    else
        echo "⚠️  Origin header present but may not match exactly"
        echo "$CORS_RESPONSE" | grep -i "access-control-allow-origin"
    fi
    
    if echo "$CORS_RESPONSE" | grep -q "access-control-allow-credentials"; then
        echo "✅ Credentials allowed"
    else
        echo "⚠️  Credentials header missing"
    fi
    
    if echo "$CORS_RESPONSE" | grep -q "access-control-allow-methods"; then
        echo "✅ Methods specified"
    else
        echo "⚠️  Methods header missing"
    fi
else
    echo "❌ CORS preflight failed"
    echo "   Response:"
    echo "$CORS_RESPONSE"
fi
echo ""

# Test 3: CORS Preflight for Token (Login)
echo "Test 3: CORS Preflight for /auth/v1/token"
echo "-----------------------------------------"
CORS_TOKEN=$(curl -s -i -X OPTIONS "${AUTH_URL}/auth/v1/token" \
  -H "Origin: ${FRONTEND_URL}" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,apikey" \
  2>&1 || echo "FAILED")

if echo "$CORS_TOKEN" | grep -q "access-control-allow-origin"; then
    echo "✅ CORS preflight for login endpoint successful"
else
    echo "❌ CORS preflight for login endpoint failed"
fi
echo ""

# Test 4: Test actual signup endpoint (will fail due to no body, but tests CORS)
echo "Test 4: Actual POST to signup (CORS test)"
echo "-----------------------------------------"
SIGNUP_TEST=$(curl -s -i -X POST "${AUTH_URL}/auth/v1/signup" \
  -H "Origin: ${FRONTEND_URL}" \
  -H "Content-Type: application/json" \
  -H "apikey: ${SUPABASE_ANON_KEY:-test-key}" \
  -d '{"email":"test@example.com","password":"testpass123"}' \
  2>&1 || echo "FAILED")

if echo "$SIGNUP_TEST" | grep -q "access-control-allow-origin"; then
    echo "✅ CORS headers present in POST response"
else
    echo "❌ CORS headers missing in POST response"
    echo "   This will cause actual signup to fail"
fi
echo ""

# Test 5: Check for wildcard with credentials (should not exist)
echo "Test 5: Security Check - No wildcard with credentials"
echo "-----------------------------------------------------"
if echo "$CORS_RESPONSE" | grep -q "access-control-allow-origin: \*" && \
   echo "$CORS_RESPONSE" | grep -q "access-control-allow-credentials: true"; then
    echo "❌ SECURITY ISSUE: Using wildcard (*) with credentials"
    echo "   This is invalid and will cause CORS errors"
    echo "   Update GOTRUE_CORS_ALLOWED_ORIGINS to specific domain"
else
    echo "✅ No security issues detected"
fi
echo ""

# Test 6: Frontend accessibility
echo "Test 6: Frontend Accessibility"
echo "------------------------------"
FRONTEND_CODE=$(curl -s -o /dev/null -w "%{http_code}" "${FRONTEND_URL}" || echo "000")
if [ "$FRONTEND_CODE" = "200" ]; then
    echo "✅ Frontend is accessible (HTTP $FRONTEND_CODE)"
else
    echo "❌ Frontend not accessible (HTTP $FRONTEND_CODE)"
fi
echo ""

# Summary
echo "=================================="
echo "Test Summary"
echo "=================================="
echo ""
echo "If all tests pass, your CORS configuration is correct."
echo ""
echo "Common issues:"
echo "  - Auth service not running: Check Docker containers"
echo "  - CORS headers missing: Update docker-compose.prod.yml"
echo "  - Wrong origin: Verify GOTRUE_CORS_ALLOWED_ORIGINS includes $FRONTEND_URL"
echo "  - Credentials issue: Cannot use wildcard (*) with credentials"
echo ""
echo "For more help, see CORS-AUTH-FIX.md"
echo ""

# Additional debugging info
if [ "$1" = "-v" ] || [ "$1" = "--verbose" ]; then
    echo "=================================="
    echo "Verbose Output"
    echo "=================================="
    echo ""
    echo "Full CORS Preflight Response:"
    echo "$CORS_RESPONSE"
fi
