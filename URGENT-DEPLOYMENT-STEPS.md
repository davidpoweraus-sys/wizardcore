# URGENT: Deployment Steps for CORS Fix

## Current Status
The CORS fixes are **NOT deployed** to production. You're still getting 403 "Origin not allowed" errors because:

1. `GET https://app.offensivewizard.com/api/auth/auth/v1/user` → 403
2. All `/api/backend/v1/users/me/*` calls → 403

## Immediate Action Required

### Step 1: Verify Current Deployment
Check what's currently deployed:

```bash
# SSH into your server
ssh your-server

# Check running containers
docker ps | grep frontend

# Check logs for middleware version
docker logs offensivewizard-app-wcilze-frontend-1 | grep "Middleware" | head -5
```

### Step 2: Deploy the CORS Fixes

**Option A: If using source-based deployment (Dokploy builds from source)**
1. **Commit all CORS fixes**:
   - `middleware.ts` (updated)
   - `app/api/auth/[...path]/route.ts` (CORS fix applied)
   - `app/api/backend/[...path]/route.ts` (CORS fix applied)
2. **Push to Git repository**
3. **Trigger redeploy in Dokploy**

**Option B: If using pre-built Docker images**
1. **Build new image with fixes**:
   ```bash
   docker build -t limpet/wizardcore-frontend:cors-fix-urgent -f Dockerfile.nextjs .
   docker push limpet/wizardcore-frontend:cors-fix-urgent
   ```
2. **Update deployment to use new image**:
   - In Dokploy: Change image to `limpet/wizardcore-frontend:cors-fix-urgent`
   - Or update `docker-compose.yml`: `image: limpet/wizardcore-frontend:cors-fix-urgent`
3. **Redeploy**

### Step 3: Verify Fix is Deployed
After deployment:

1. **Check middleware version** in logs:
   ```
   docker logs [frontend-container] | grep "Middleware"
   ```
   Should show: `🔍 Middleware session-refresh-fix-20260104-1159 executing`

2. **Test login**:
   - Clear browser cache/cookies
   - Log in at `https://app.offensivewizard.com/login`
   - Check browser console for 403 errors (should be gone)

## Critical: The CORS Fix Explained

The issue is in the `validateOrigin` function in both proxy files:

**Before (causing 403)**:
```typescript
if (!origin || origin === '*') {
  // In production, avoid wildcard when using credentials
  if (process.env.NODE_ENV === 'production') {
    return null  // ← This causes 403 for same-origin requests!
  }
  return '*'
}
```

**After (fixed)**:
```typescript
if (!origin) {
  // Allow null/empty origin for same-origin requests
  return 'same-origin'
}
if (origin === '*') {
  if (process.env.NODE_ENV === 'production') {
    return null
  }
  return '*'
}
```

## Why This Matters

- **Same-origin requests** (frontend → `app.offensivewizard.com/api/auth/*`) often don't include Origin header
- **Old code**: `validateOrigin(null)` returns `null` → 403 error
- **Fixed code**: `validateOrigin(null)` returns `'same-origin'` → request allowed

## Quick Verification Script

Create a test script to verify the fix:

```bash
#!/bin/bash
# test-cors-fix.sh
echo "Testing CORS fix..."
curl -v -H "Cookie: sb-app-auth-token=test" \
  https://app.offensivewizard.com/api/auth/auth/v1/user 2>&1 | grep "HTTP/"
```

Should return `HTTP/3 200` (not `HTTP/3 403`)

## If Still Failing After Deployment

1. **Check proxy logs**:
   ```bash
   docker logs [frontend-container] | grep -A5 -B5 "validateOrigin"
   ```

2. **Verify files are updated**:
   ```bash
   docker exec [frontend-container] cat /app/app/api/auth/[...path]/route.ts | grep -A10 "function validateOrigin"
   ```

3. **Restart frontend**:
   ```bash
   docker-compose restart frontend
   ```

## Summary

The CORS fixes are **code changes** that need to be **deployed to production**. Until they're deployed, you'll continue to get 403 errors and be shown as a "guest user".

**Deploy the fixes NOW to resolve the login issue.**