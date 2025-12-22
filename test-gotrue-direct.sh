#!/bin/bash

# Test GoTrue endpoints directly
echo "Testing GoTrue API endpoints..."
echo ""

AUTH_URL="https://auth.offensivewizard.com"
ANON_KEY="aNBMlTAljHvKyR2dPu6R6nyggeW2398Na3R4XL1+oyebUDiuzSO61nZzoVmRi0h4"

echo "1. Test health endpoint:"
curl -s "${AUTH_URL}/health" | jq '.' || curl -s "${AUTH_URL}/health"
echo ""
echo ""

echo "2. Test signup endpoint (should work):"
curl -v -X POST "${AUTH_URL}/auth/v1/signup" \
  -H "Content-Type: application/json" \
  -H "apikey: ${ANON_KEY}" \
  -d '{
    "email": "test@example.com",
    "password": "test123456"
  }' 2>&1 | grep -E "HTTP|< " | head -20
echo ""
echo ""

echo "3. Try without /auth prefix:"
curl -v -X POST "${AUTH_URL}/signup" \
  -H "Content-Type: application/json" \
  -H "apikey: ${ANON_KEY}" \
  -d '{
    "email": "test2@example.com",
    "password": "test123456"
  }' 2>&1 | grep -E "HTTP|< " | head -20
echo ""
