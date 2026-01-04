# Session Refresh & CORS Fix Deployment Guide

## Summary
Fixed the production issue where users could log in but were shown as "guest users" and couldn't access their data. The problem was caused by:

1. **CORS issues in auth/backend proxies**: Same-origin requests (without Origin header) were being rejected with 403 "Origin not allowed"
2. **Missing session refresh in middleware**: The middleware was only checking for cookie presence, not refreshing sessions

## Changes Made

### 1. **Updated Middleware (`middleware.ts`)**
- Added session refresh awareness for authenticated users
- Added CORS headers for API routes
- Version: `session-refresh-fix-20260104-1159`

### 2. **Fixed Auth Proxy CORS (`app/api/auth/[...path]/route.ts`)**
- **Critical fix**: Updated `validateOrigin` function to allow `null`/empty origins for same-origin requests
- Same-origin requests often don't include Origin header
- Returns `'same-origin'` placeholder which is then resolved to actual origin
- Updated OPTIONS handler and main proxy to handle `'same-origin'` value

### 3. **Fixed Backend Proxy CORS (`app/api/backend/[...path]/route.ts`)**
- Same fix as auth proxy: allow `null` origins for same-origin requests
- Updated `validateOrigin`, OPTIONS handler, and main proxy

## Root Cause Analysis

### The Problem:
1. User logs in successfully (auth cookie `sb-app-auth-token` is set)
2. Middleware sees cookie and allows access to `/dashboard`
3. Frontend tries to fetch user data via `/api/auth/auth/v1/user`
4. **Auth proxy rejects request with 403** because Origin header is `null` (same-origin request)
5. Supabase client thinks there's no user, shows "guest" state
6. All backend API calls also fail with 403 CORS errors

### The Solution:
- Allow `null` origins in CORS validation for same-origin requests
- When origin is `null`, treat it as `'same-origin'` and resolve to actual domain
- This allows the frontend to make requests to its own API endpoints

## New Docker Image
Image: `limpet/wizardcore-frontend:session-refresh-cors-fix`

Built with:
- Updated middleware with session refresh awareness
- Fixed CORS validation in auth and backend proxies
- All environment variables from production

## How to Deploy

### Option 1: Update Image Tag in Dokploy (Quickest)
1. Go to Dokploy Dashboard → WizardCore Application
2. Click "Edit" or "Settings"
3. Update image tag to: **`limpet/wizardcore-frontend:session-refresh-cors-fix`**
4. Click "Save" and "Redeploy"

### Option 2: Update docker-compose.yml
```yaml
frontend:
  image: limpet/wizardcore-frontend:session-refresh-cors-fix
  # ... rest of configuration
```

### Option 3: Rebuild from Source
1. Ensure all changes are committed:
   - `middleware.ts` (updated)
   - `app/api/auth/[...path]/route.ts` (updated)
   - `app/api/backend/[...path]/route.ts` (updated)
2. Build and push new image:
   ```bash
   docker build -t limpet/wizardcore-frontend:session-refresh-cors-fix -f Dockerfile.nextjs .
   docker push limpet/wizardcore-frontend:session-refresh-cors-fix
   ```
3. Update deployment to use new image tag

## Verification Steps

After deployment:

1. **Clear browser cache and cookies** for `app.offensivewizard.com`
2. **Log in** and check browser console for:
   ```
   🔍 Middleware session-refresh-fix-20260104-1159 executing for path: /dashboard
   🔍 Has auth cookie: true
   🔍 Attempting to refresh session for authenticated user
   ```

3. **Check Network tab** for:
   - `/api/auth/auth/v1/user` should return 200 (not 403)
   - `/api/backend/v1/users/me/*` calls should return 200 (not 403)
   - Response headers should include `X-Middleware-Version: session-refresh-fix-20260104-1159`

4. **Verify user is recognized**:
   - Dashboard should show user-specific data (not "guest" or generic content)
   - User profile should be accessible
   - Stats and progress data should load

## Expected Behavior After Fix

1. **Login flow**:
   - User enters credentials → auth succeeds → cookie set
   - Redirect to `/dashboard` succeeds

2. **Session recognition**:
   - Middleware sees cookie and allows access
   - Frontend fetches user data successfully
   - User is shown as logged in (not guest)

3. **Data loading**:
   - Backend API calls succeed with proper CORS headers
   - User stats, progress, activities load correctly
   - No "Failed to fetch" errors in console

## Debugging

If issues persist:

1. **Check middleware logs**:
   ```
   docker compose logs frontend | grep -A5 -B5 "Middleware"
   ```

2. **Check auth proxy logs**:
   ```
   docker compose logs frontend | grep -A5 -B5 "GoTrue Proxy"
   ```

3. **Check browser console** for:
   - CORS errors (should be gone)
   - Network request failures
   - Console logs from middleware

4. **Verify cookies**:
   - `sb-app-auth-token` should be present after login
   - Cookie should be `HttpOnly`, `Secure`, `SameSite=Lax`

## Rollback

If issues persist, rollback to previous version:
- Use image tag: `limpet/wizardcore-frontend:rsc-fix-v2`
- Or: `limpet/wizardcore-frontend:login-fix-v6`

## Files Modified

1. `middleware.ts` - Added session refresh awareness and CORS headers
2. `app/api/auth/[...path]/route.ts` - Fixed CORS validation for same-origin requests
3. `app/api/backend/[...path]/route.ts` - Fixed CORS validation for same-origin requests

## Technical Details

### CORS Fix Details:
- **Before**: `validateOrigin(null)` returned `null` in production → 403 error
- **After**: `validateOrigin(null)` returns `'same-origin'` → resolved to actual domain → request allowed

### Session Refresh:
- Middleware now logs session refresh attempts
- Sets proper CORS headers for API routes
- Ensures cookies are forwarded properly

This fix resolves the "guest user" issue by allowing same-origin API requests and ensuring sessions are properly recognized.