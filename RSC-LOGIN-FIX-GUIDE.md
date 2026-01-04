# RSC Login Fix Guide for Production

## Problem Summary
Users are getting stuck on the login page with the error: "Failed to fetch RSC payload for https://app.offensivewizard.com/dashboard. Falling back to browser navigation. TypeError: NetworkError when attempting to fetch resource."

## Root Cause Analysis
The issue occurs because:

1. **Next.js RSC (React Server Components) fetches** are made client-side when navigating to `/dashboard`
2. These RSC fetches don't include authentication cookies by default
3. The middleware intercepts these fetches and returns a 307 redirect to `/login`
4. The RSC fetch fails because it receives a redirect response instead of RSC data
5. Next.js falls back to browser navigation, but the redirect loop continues

## Changes Made

### 1. Updated Middleware (`middleware.ts`)
- Added detection for RSC fetch requests (checking for `_rsc` query param, `x-nextjs-data: 1`, `next-router-prefetch: 1`, `next-action: 1` headers)
- For RSC fetches to protected routes without auth cookies, returns **401 Unauthorized** instead of 307 redirect
- This allows the client to handle authentication errors properly
- Regular browser navigation still gets redirected as before

### 2. Enhanced CORS Configuration (`next.config.ts`)
- Added missing `Access-Control-Allow-Origin` header (set to `https://app.offensivewizard.com` in production)
- Added Next.js specific headers to `Access-Control-Allow-Headers`:
  - `Next-Router-Prefetch`
  - `Next-Router-State-Tree`
  - `Next-Url`
  - `RSC`
- Added `Access-Control-Expose-Headers` to expose Next.js response headers
- This ensures RSC fetches can properly communicate with the server

### 3. Maintained Existing Fixes
- Auth proxy remains unchanged (already handles CORS correctly)
- Login page continues to use `window.location.href` to bypass Next.js router issues
- Cookie configuration remains with `SameSite=Lax` for cross-domain support

## How the Fix Works

### Before Fix:
```
1. User logs in successfully
2. Login page calls `window.location.href = '/dashboard'`
3. Browser makes GET request to `/dashboard`
4. Middleware sees auth cookie, allows through
5. Dashboard page loads
6. Next.js client-side router tries to fetch RSC data for `/dashboard`
7. RSC fetch doesn't include cookies
8. Middleware returns 307 redirect to `/login`
9. RSC fetch fails with network error (can't follow redirect)
10. Next.js falls back to browser navigation
11. Loop continues
```

### After Fix:
```
1. User logs in successfully
2. Login page calls `window.location.href = '/dashboard'`
3. Browser makes GET request to `/dashboard`
4. Middleware sees auth cookie, allows through
5. Dashboard page loads
6. Next.js client-side router tries to fetch RSC data for `/dashboard`
7. RSC fetch doesn't include cookies
8. Middleware detects it's an RSC fetch, returns 401 Unauthorized
9. RSC fetch receives proper error response (not redirect)
10. Next.js handles 401 appropriately (may trigger re-auth)
11. User stays on dashboard
```

## Deployment Instructions

### Option 1: Update Existing Deployment (Dokploy)

1. **Build new Docker image:**
   ```bash
   docker build -f Dockerfile.nextjs -t limpet/wizardcore-frontend:rsc-fix-v1 .
   ```

2. **Push to Docker Hub:**
   ```bash
   docker push limpet/wizardcore-frontend:rsc-fix-v1
   ```

3. **Update Dokploy deployment:**
   - Go to Dokploy Dashboard → WizardCore Application
   - Click "Edit" or "Settings"
   - Update image tag to: `limpet/wizardcore-frontend:rsc-fix-v1`
   - Click "Save" and "Redeploy"

### Option 2: Rebuild from Source in Dokploy

1. **Ensure build arguments are correct:**
   - `NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth`
   - `NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=`
   - `NEXT_PUBLIC_BACKEND_URL=https://app.offensivewizard.com/api`
   - `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com`

2. **Trigger rebuild:**
   - In Dokploy, go to Application → Build
   - Click "Rebuild" or "Redeploy"

### Option 3: Manual Deployment

1. **Copy updated files to server:**
   ```bash
   scp middleware.ts next.config.ts server:/path/to/wizardcore/
   ```

2. **Restart Next.js service:**
   ```bash
   docker-compose restart frontend
   ```

## Verification Steps

After deployment:

1. **Clear browser cache and cookies** for `app.offensivewizard.com`
2. **Open browser console** (F12 → Console tab)
3. **Navigate to** `https://app.offensivewizard.com/login`
4. **Log in** with credentials
5. **Check console logs** for:
   - `🎲 BLUE DIE TEST - Login successful`
   - `🔍 Middleware executing for path: /dashboard`
   - `🔍 Has auth cookie: true`
   - No RSC fetch errors

6. **Check Network tab** for:
   - `/dashboard` request returns 200 (not 307)
   - Any `_rsc` requests return proper responses (200 or 401)

## Expected Behavior

1. User enters credentials and clicks "Sign In"
2. Button shows "Signing in..." while API call happens
3. Auth API returns 200 with access token
4. Login page detects success and redirects via `window.location.href`
5. Dashboard loads successfully
6. No RSC fetch errors in console
7. User stays logged in across page navigation

## Debugging

If issues persist:

### Check Dokploy logs:
```bash
# Look for middleware debug messages
docker logs wizardcore-frontend | grep -E "🔍|🎲|❌"
```

### Check browser Network tab:
- Filter for requests with `_rsc` in URL
- Check response status codes (should be 200 or 401, not 307)
- Check if cookies are being sent with RSC fetches

### Test RSC fetch manually:
```javascript
// In browser console on dashboard page
fetch('/dashboard?_rsc=test', {
  credentials: 'include'  // Important: include cookies
})
.then(r => console.log('Status:', r.status))
.then(d => console.log('Data:', d))
.catch(e => console.error('Error:', e))
```

## Rollback Plan

If the fix causes issues:

1. **Revert to previous image:**
   - Use tag: `limpet/wizardcore-frontend:login-fix-v6`
   
2. **Or revert code changes:**
   - Restore previous versions of `middleware.ts` and `next.config.ts`
   - Rebuild and redeploy

## Files Modified

1. `middleware.ts` - Added RSC fetch detection and 401 response
2. `next.config.ts` - Enhanced CORS headers for Next.js specific headers

## Technical Details

### RSC Fetch Detection
The middleware detects RSC fetches by checking:
- Query parameter: `_rsc=` in URL
- Header: `x-nextjs-data: 1`
- Header: `next-router-prefetch: 1`
- Header: `next-action: 1`

### CORS Headers Added
- `Access-Control-Allow-Origin: https://app.offensivewizard.com`
- `Access-Control-Allow-Headers: ... Next-Router-Prefetch, Next-Router-State-Tree, Next-Url, RSC`
- `Access-Control-Expose-Headers: X-NextJS-Data, X-RSC, X-Action-Redirect`

### Authentication Flow
- Regular page requests: 307 redirect to `/login` if no auth cookie
- RSC fetch requests: 401 Unauthorized if no auth cookie
- This distinction prevents redirect loops for API-like RSC fetches

## Support

For further issues:
1. Check Dokploy application logs for middleware debug messages
2. Examine browser console for RSC fetch errors
3. Verify CORS headers are being set correctly
4. Ensure auth proxy (`/api/auth`) is working correctly