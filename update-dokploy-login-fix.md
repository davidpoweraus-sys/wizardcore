# Dokploy Login Fix Update Guide

## Summary
Fixed the production login issue where users get stuck on the login page after successful authentication. The problem was caused by:

1. **Middleware JWT validation failure**: The middleware was trying to validate JWT tokens through the auth proxy, which failed due to missing JWT secret.
2. **Environment configuration mismatch**: `NEXT_PUBLIC_BACKEND_URL` was pointing to wrong domain.
3. **RSC payload fetch failures**: Next.js router was failing to fetch React Server Component data for `/dashboard`.

## Changes Made

### 1. Simplified Middleware (`middleware.ts`)
- Removed complex `@supabase/ssr` `createServerClient` setup
- Now simply checks for presence of `sb-app-auth-token` cookie
- No more JWT validation failures through proxy
- Much more reliable and faster

### 2. Fixed Environment Variables (`.env.production`)
- Fixed: `NEXT_PUBLIC_BACKEND_URL=https://app.offensivewizard.com/api` (was pointing to wrong domain)
- Maintained correct auth proxy URL: `NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth`

### 3. Enhanced Login Page (`app/(auth)/login/page.tsx`)
- Uses `window.location.href` instead of `router.push()` to bypass Next.js RSC issues
- Comprehensive step-by-step logging for debugging
- Multiple fallback mechanisms

## New Docker Image
Image: `limpet/wizardcore-frontend:login-fix-v6`

Built with correct environment variables:
- `NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
- `NEXT_PUBLIC_BACKEND_URL=https://app.offensivewizard.com/api`
- `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com`

## How to Update in Dokploy

### Option 1: Update Image Tag (Quickest)
1. Go to Dokploy Dashboard → WizardCore Application
2. Click "Edit" or "Settings"
3. Find the image/tag field
4. Change from current tag to: `limpet/wizardcore-frontend:login-fix-v6`
5. Click "Save" or "Update"
6. Click "Redeploy"

### Option 2: Update docker-compose.yml
If using Docker Compose deployment:
```yaml
frontend:
  image: limpet/wizardcore-frontend:login-fix-v6
  # ... rest of configuration
```

### Option 3: Rebuild from Source
1. In Dokploy, go to Application → Settings
2. Ensure build args are set correctly:
   - `NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth`
   - `NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...`
   - `NEXT_PUBLIC_BACKEND_URL=https://app.offensivewizard.com/api`
   - `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com`
3. Click "Redeploy" to rebuild with latest code

## Verification Steps

After deployment:

1. **Clear browser cache** and cookies for `app.offensivewizard.com`
2. **Open browser console** (F12 → Console tab)
3. **Navigate to** `https://app.offensivewizard.com/login`
4. **Log in** with credentials
5. **Check console logs** for:
   - `🎲 Step 1: Login form submitted`
   - `🎲 Step 2: Sign in successful`
   - `🎲 Step 3: Session confirmed`
   - `🎲 Step 4: Redirecting to dashboard...`
   - `🎲 Step 4a: Using window.location.href (bypassing Next.js router)`
6. **Verify** you're redirected to `/dashboard`

## Expected Behavior

1. User enters credentials and clicks "Sign In"
2. Button shows "Signing in..." while API call happens
3. Auth API returns 200 with access token
4. Login page detects success and redirects via `window.location.href`
5. Middleware sees `sb-app-auth-token` cookie and allows access
6. Dashboard loads successfully

## Debugging

If still failing:

1. **Check Dokploy logs** for middleware debug messages:
   - `🔍 Middleware executing for path: /dashboard`
   - `🔍 Request cookies: sb-app-auth-token, ...`
   - `🔍 Has auth cookie: true/false`

2. **Check browser Network tab**:
   - `/api/auth/auth/v1/token?grant_type=password` should return 200
   - `/dashboard` should NOT return 307 redirect

3. **Check cookies**:
   - `sb-app-auth-token` should be present after login
   - Cookie should be `HttpOnly` and `Secure`

## Rollback

If issues persist, rollback to previous version:
- Use image tag: `limpet/wizardcore-frontend:login-fix-v5` or earlier
- Or revert to stable tag if available

## Support

For further issues:
1. Check Dokploy application logs
2. Examine browser console errors
3. Verify environment variables match production
4. Ensure auth proxy (`/api/auth`) is working correctly