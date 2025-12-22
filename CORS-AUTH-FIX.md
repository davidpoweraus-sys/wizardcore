# CORS Authentication Fix Documentation

## Problem Summary
Users were experiencing CORS errors when trying to create accounts. The frontend at `https://offensivewizard.com` could not communicate with the Supabase Auth service at `https://auth.offensivewizard.com` due to cross-origin restrictions.

## Root Causes Identified

1. **Missing Next.js Middleware**: No middleware to handle Supabase session refresh and cookie management
2. **Incorrect CORS Configuration**: Using wildcard (`*`) origin with credentials is not allowed per CORS specification
3. **Cookie Domain Mismatch**: Cookies weren't being shared properly between `offensivewizard.com` and `auth.offensivewizard.com`
4. **Missing Auth Callback Route**: No route to handle OAuth callbacks and token exchange

## Solutions Implemented

### 1. Next.js Middleware (`middleware.ts`)
Created a middleware file that:
- Manages Supabase session refresh on every request
- Handles cookie synchronization between browser and server
- Protects routes (redirects unauthenticated users from `/dashboard` to `/login`)
- Redirects authenticated users away from `/login` and `/register`

**Why this matters**: Without middleware, the auth session won't be refreshed automatically, leading to users being logged out randomly.

### 2. CORS Headers in Next.js Config (`next.config.ts`)
Added proper CORS headers to allow communication with the auth subdomain:
- `Access-Control-Allow-Credentials: true` - Allows cookies to be sent
- `Access-Control-Allow-Origin` - Set to the Supabase URL (not wildcard)
- `Access-Control-Allow-Headers` - Includes all necessary auth headers

### 3. Enhanced Supabase Client (`lib/supabase/client.ts`)
Updated the browser client with custom cookie handlers that:
- Properly encode/decode cookie values
- Set correct cookie attributes (`SameSite=Lax`, `Secure` in production)
- Handle cross-domain cookie sharing
- Support cookie removal for logout

**Why SameSite=Lax**: This allows cookies to be sent on top-level navigations (like redirects from auth subdomain), while still providing CSRF protection.

### 4. Fixed Docker Compose CORS Config (`docker-compose.prod.yml`)
Updated Supabase Auth (GoTrue) environment variables:
```yaml
# BEFORE (WRONG):
GOTRUE_CORS_ALLOWED_ORIGINS: "*"

# AFTER (CORRECT):
GOTRUE_CORS_ALLOWED_ORIGINS: "https://offensivewizard.com,http://localhost:3000,http://localhost"
GOTRUE_CORS_ALLOW_CREDENTIALS: "true"
```

**Critical**: When using `Access-Control-Allow-Credentials: true`, you MUST specify exact origins, not `*`.

### 5. Auth Callback Route (`app/auth/callback/route.ts`)
Created a server-side route to:
- Exchange authorization codes for session tokens
- Handle redirects after email verification
- Support both local development and production (with load balancer)

### 6. Auth Error Page (`app/auth/auth-code-error/page.tsx`)
User-friendly error page for when authentication fails.

## How Authentication Flow Works Now

### Registration Flow:
1. User fills out registration form at `https://offensivewizard.com/register`
2. Form submits to Supabase Auth at `https://auth.offensivewizard.com/auth/v1/signup`
3. Supabase Auth validates CORS (now succeeds with proper origin)
4. Auth service creates account and returns session tokens
5. Client sets cookies with proper domain/SameSite settings
6. User is redirected to `/dashboard` (or email confirmation page if enabled)
7. Middleware refreshes session on next request

### Login Flow:
1. User submits login form at `https://offensivewizard.com/login`
2. Form submits to `https://auth.offensivewizard.com/auth/v1/token?grant_type=password`
3. Auth service validates credentials and returns tokens
4. Cookies are set and user is redirected to `/dashboard`

### Session Management:
- Middleware runs on every request
- Checks cookie validity and refreshes if needed
- Protects routes from unauthenticated access
- Synchronizes session state between client and server

## Testing the Fix

### 1. Test CORS Preflight:
```bash
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,apikey"
```

**Expected response**:
```
HTTP/2 204
access-control-allow-origin: https://offensivewizard.com
access-control-allow-credentials: true
access-control-allow-methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
access-control-allow-headers: ...
```

### 2. Test Registration:
1. Go to `https://offensivewizard.com/register`
2. Fill out the form
3. Submit
4. Check browser console for errors (should be none)
5. Check Network tab - should see successful POST to auth.offensivewizard.com
6. Should redirect to dashboard

### 3. Test Cookie Persistence:
1. After login/register, check Application > Cookies in DevTools
2. Should see Supabase cookies with:
   - Domain: `.offensivewizard.com` or specific subdomain
   - SameSite: `Lax`
   - Secure: `true` (in production)
   - HttpOnly: varies by cookie type

### 4. Test Protected Routes:
1. While logged out, try visiting `/dashboard`
2. Should redirect to `/login`
3. After login, try visiting `/login`
4. Should redirect to `/dashboard`

## Common Issues and Troubleshooting

### Issue: "CORS policy: Response to preflight request doesn't pass access control check"
**Solution**: Verify that `GOTRUE_CORS_ALLOWED_ORIGINS` in docker-compose.prod.yml includes your exact frontend URL (not `*`).

### Issue: User gets logged out randomly
**Solution**: Ensure middleware.ts is in place and properly configured. The middleware refreshes sessions automatically.

### Issue: Cookies not being set
**Solution**: Check that:
- Cookies have correct `SameSite` attribute
- Using HTTPS in production (required for `Secure` cookies)
- Cookie domain matches your domain structure

### Issue: "Invalid API key" or "Missing authorization header"
**Solution**: Verify that `NEXT_PUBLIC_SUPABASE_ANON_KEY` is set correctly in:
- Dockerfile.nextjs build args
- docker-compose.prod.yml frontend environment
- Coolify environment variables (if applicable)

### Issue: Redirect after signup not working
**Solution**: Check that:
- `/auth/callback` route exists
- `emailRedirectTo` in register page matches your domain
- Email confirmation is disabled (`GOTRUE_MAILER_AUTOCONFIRM: 'true'`) if not using email service

## Deployment Checklist

Before deploying to production:

- [ ] Middleware file (`middleware.ts`) is in project root
- [ ] Next.js config has proper CORS headers
- [ ] Supabase client has custom cookie handlers
- [ ] Docker compose has correct CORS origins (no `*` with credentials)
- [ ] Auth callback route exists at `/app/auth/callback/route.ts`
- [ ] Environment variables are set correctly:
  - `NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com`
  - `NEXT_PUBLIC_SUPABASE_ANON_KEY=<your-key>`
- [ ] SSL certificates are valid for both domains
- [ ] Test signup/login flow in production

## Security Considerations

1. **ANON_KEY**: This is safe to expose in frontend code - it only allows public operations
2. **JWT_SECRET**: Keep this secret - never expose in frontend
3. **SameSite=Lax**: Provides CSRF protection while allowing auth flows
4. **Secure Flag**: Ensures cookies only sent over HTTPS
5. **Specific Origins**: Only allow exact domains you control in CORS

## Related Files

- `/middleware.ts` - Session management and route protection
- `/lib/supabase/client.ts` - Browser client with cookie handling
- `/lib/supabase/server.ts` - Server-side client (unchanged)
- `/next.config.ts` - CORS headers configuration
- `/docker-compose.prod.yml` - Supabase Auth CORS settings
- `/app/auth/callback/route.ts` - OAuth callback handler
- `/app/(auth)/register/page.tsx` - Registration form
- `/app/(auth)/login/page.tsx` - Login form

## References

- [Supabase Auth SSR Guide](https://supabase.com/docs/guides/auth/server-side/nextjs)
- [CORS with Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS#requests_with_credentials)
- [SameSite Cookies Explained](https://web.dev/articles/samesite-cookies-explained)
- [GoTrue CORS Configuration](https://github.com/supabase/gotrue#cors-configuration)
