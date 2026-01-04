# Immediate Fix Guide for Login Issue

## Current Status
Deployment failed because the Docker image `limpet/wizardcore-frontend:session-refresh-cors-fix` doesn't exist yet.

## Quick Solution
Revert to working image `limpet/wizardcore-frontend:rsc-fix-v2` and apply the CORS fixes manually.

## Step 1: Update docker-compose.yml
Already done - reverted to `limpet/wizardcore-frontend:rsc-fix-v2`

## Step 2: Apply CORS Fixes to Current Deployment

If you're deploying from source (not using pre-built image), you need to apply these fixes:

### Fix 1: Update Auth Proxy CORS
Replace the `validateOrigin` function in `app/api/auth/[...path]/route.ts`:

```typescript
function validateOrigin(origin: string | null): string | null {
  // CRITICAL FIX: Allow null/empty origin for same-origin requests
  if (!origin) {
    return 'same-origin'
  }
  
  if (origin === '*') {
    if (process.env.NODE_ENV === 'production') {
      return null
    }
    return '*'
  }

  try {
    const url = new URL(origin)
    const hostname = url.hostname
    
    if (hostname === 'offensivewizard.com' ||
        hostname === 'www.offensivewizard.com' ||
        hostname.endsWith('.offensivewizard.com')) {
      return origin
    }
    
    if (hostname === 'localhost') {
      return origin
    }
    
    const allowedOrigins = process.env.ALLOWED_ORIGINS?.split(',') || []
    if (allowedOrigins.includes(origin)) {
      return origin
    }
    
    return null
  } catch (error) {
    console.warn('Invalid origin format:', origin)
    return null
  }
}
```

### Fix 2: Update Backend Proxy CORS
Same fix in `app/api/backend/[...path]/route.ts`

### Fix 3: Update Middleware
Update `middleware.ts` to version `session-refresh-fix-20260104-1159` (already done)

## Step 3: Redeploy
```bash
# If using Docker Compose
docker-compose pull frontend
docker-compose up -d frontend

# If using Dokploy
# Just redeploy with current code (image: limpet/wizardcore-frontend:rsc-fix-v2)
```

## Step 4: Verify Fix
1. Clear browser cache and cookies
2. Log in at `https://app.offensivewizard.com/login`
3. Check browser console for:
   - No 403 errors on `/api/auth/auth/v1/user`
   - User data loads successfully
   - No "guest user" state

## Alternative: Build New Image
If you want to build the new image:

```bash
# 1. Apply all fixes (middleware.ts, auth proxy, backend proxy)
# 2. Build and push
./build-session-refresh-fix.sh

# 3. Update docker-compose.yml to use:
#    image: limpet/wizardcore-frontend:session-refresh-cors-fix

# 4. Redeploy
```

## Expected Result
After applying CORS fixes:
- ✅ Same-origin requests allowed (no 403 "Origin not allowed")
- ✅ User data loads after login
- ✅ No "guest user" state
- ✅ Backend API calls work

The root issue was CORS validation rejecting same-origin requests. The fix allows `null` origins for same-origin requests.