# RSC Login Fix v2 - Production Deployment Guide

## Current Status
The fix is still not deployed to production. The logs show:
- Cookie `sb-app-auth-token` is present in requests
- Middleware is still returning 307 redirect for `/dashboard`
- RSC fetch error persists: "Failed to fetch RSC payload for https://app.offensivewizard.com/dashboard"

## Version Tracking Added
To verify deployment, we've added version tracking:

### 1. Middleware Version Identifier
```typescript
const MIDDLEWARE_VERSION = 'rsc-fix-v2-20260104-1130'
```

### 2. Console Logs
All middleware logs now include the version:
```
🔍 Middleware rsc-fix-v2-20260104-1130 executing for path: /dashboard
```

### 3. Response Headers
All responses include:
```
X-Middleware-Version: rsc-fix-v2-20260104-1130
```

## Enhanced RSC Detection
Updated to detect more Next.js RSC request patterns:
- `_rsc=` query parameter
- `x-nextjs-data: 1` header
- `next-router-prefetch: 1` header  
- `next-action: 1` header
- `RSC: 1` header
- `Next-Router-State-Tree` header (any value)

## Files Modified

### 1. `middleware.ts`
- Added version identifier `MIDDLEWARE_VERSION`
- Enhanced RSC detection with more headers
- Added `X-Middleware-Version` header to all responses
- Returns 401 for RSC fetches without auth (instead of 307)

### 2. `next.config.ts`
- Added `X-Middleware-Version` to exposed CORS headers
- This allows browser to see the version header

## Deployment Instructions

### Step 1: Build New Docker Image
```bash
docker build -f Dockerfile.nextjs -t limpet/wizardcore-frontend:rsc-fix-v2 .
```

### Step 2: Push to Docker Hub
```bash
docker push limpet/wizardcore-frontend:rsc-fix-v2
```

### Step 3: Update Dokploy Deployment
1. Go to Dokploy Dashboard → WizardCore Application
2. Click "Edit" or "Settings"
3. Update image tag to: `limpet/wizardcore-frontend:rsc-fix-v2`
4. Click "Save" and "Redeploy"

## Verification Steps

### After Deployment:

1. **Clear browser cache and cookies** for `app.offensivewizard.com`
2. **Open browser console** (F12 → Console tab)
3. **Navigate to** `https://app.offensivewizard.com/login`
4. **Log in** with credentials
5. **Check console logs** for:
   - `🔍 Middleware rsc-fix-v2-20260104-1130 executing for path: /dashboard`
   - This confirms the new middleware is running

6. **Check Network tab**:
   - Select the `/dashboard` request
   - Look for `X-Middleware-Version: rsc-fix-v2-20260104-1130` in response headers
   - Status should be 200 (not 307)

7. **Verify no RSC errors**:
   - No "Failed to fetch RSC payload" errors in console
   - Dashboard loads successfully

## Testing the Fix

### Test Script
Run the verification script:
```bash
chmod +x test-middleware-version.sh
./test-middleware-version.sh
```

### Manual Test in Browser Console
```javascript
// Test if middleware is detecting RSC requests correctly
fetch('/dashboard', {
  headers: {
    'RSC': '1',
    'Next-Router-State-Tree': 'test'
  }
})
.then(r => {
  console.log('Status:', r.status)
  console.log('Middleware Version:', r.headers.get('X-Middleware-Version'))
  return r.text()
})
.then(d => console.log('Response:', d.slice(0, 100)))
.catch(e => console.error('Error:', e))
```

## Expected Behavior

### With Auth Cookie:
1. Regular GET `/dashboard` → 200 OK with dashboard page
2. RSC fetch to `/dashboard` → 200 OK with RSC data
3. No redirects, no RSC fetch errors

### Without Auth Cookie:
1. Regular GET `/dashboard` → 307 redirect to `/login`
2. RSC fetch to `/dashboard` → 401 Unauthorized (not redirect)
3. Next.js handles 401 appropriately

## Debugging

### If Still Getting 307 Redirect:
1. Check middleware logs for version identifier
2. Verify the new image is actually deployed
3. Check if RSC detection is working:
   ```javascript
   // In browser console after login
   fetch('/dashboard', { 
     credentials: 'include',
     headers: { 'RSC': '1' }
   })
   .then(r => console.log('Status:', r.status))
   .catch(e => console.error('Error:', e))
   ```

### Check Dokploy Logs:
```bash
# Look for middleware version in logs
docker logs wizardcore-frontend | grep -E "rsc-fix-v2|🔍 Middleware"
```

## Rollback Plan

If issues persist:
1. Revert to previous image: `limpet/wizardcore-frontend:login-fix-v6`
2. Or contact support for further debugging

## Root Cause Analysis

The issue is that Next.js makes two types of requests to `/dashboard`:

1. **Browser navigation**: Regular GET request, includes cookies, should get 307 if no auth
2. **RSC fetch**: Client-side fetch for React Server Components, may not include cookies, should get 401 if no auth (not 307)

The fix ensures RSC fetches get 401 (which Next.js can handle) instead of 307 (which causes fetch errors).