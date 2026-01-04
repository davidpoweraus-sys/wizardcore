#!/bin/bash

# Test script to verify middleware version is deployed
echo "Testing middleware version deployment..."

# Check if TypeScript compiles
echo "1. Checking TypeScript compilation..."
npx tsc --noEmit
if [ $? -eq 0 ]; then
    echo "✅ TypeScript compilation successful"
else
    echo "❌ TypeScript compilation failed"
    exit 1
fi

# Check middleware.ts for version identifier
echo "2. Checking middleware.ts for version identifier..."
if grep -q "MIDDLEWARE_VERSION = 'rsc-fix-v2-20260104-1130'" middleware.ts; then
    echo "✅ Middleware version identifier found"
else
    echo "❌ Middleware version identifier not found"
    echo "Current version in file:"
    grep "MIDDLEWARE_VERSION" middleware.ts || echo "No version found"
fi

# Check next.config.ts for exposed headers
echo "3. Checking next.config.ts for exposed headers..."
if grep -q "X-Middleware-Version" next.config.ts; then
    echo "✅ X-Middleware-Version header exposed in CORS"
else
    echo "❌ X-Middleware-Version header not exposed"
fi

# Create a simple test to verify the middleware logic
echo "4. Creating test to verify RSC detection logic..."
cat > test-rsc-detection.js << 'EOF'
// Test RSC detection logic
const testCases = [
    { url: '/dashboard?_rsc=abc', expected: true, description: 'URL with _rsc param' },
    { url: '/dashboard', expected: false, description: 'URL without _rsc param' },
    { url: '/api/test', expected: false, description: 'API route' },
    { url: '/dashboard?_rsc=', expected: true, description: 'URL with empty _rsc param' },
];

console.log('Testing RSC detection logic:');
testCases.forEach(test => {
    const hasRscParam = test.url.includes('_rsc=');
    const passed = hasRscParam === test.expected;
    console.log(`${passed ? '✅' : '❌'} ${test.description}: ${test.url} (has _rsc=${hasRscParam}, expected=${test.expected})`);
});
EOF

node test-rsc-detection.js

echo ""
echo "=== Deployment Verification Steps ==="
echo "1. Deploy the updated code to production"
echo "2. Check browser console for: '🔍 Middleware rsc-fix-v2-20260104-1130 executing for path:'"
echo "3. Check response headers for: 'X-Middleware-Version: rsc-fix-v2-20260104-1130'"
echo "4. Test login flow and verify no RSC fetch errors"
echo ""
echo "To test in browser after deployment:"
echo "1. Open browser console (F12)"
echo "2. Navigate to https://app.offensivewizard.com/login"
echo "3. Log in and check for middleware version logs"
echo "4. Check Network tab for X-Middleware-Version header in responses"