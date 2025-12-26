# Authentication Proxy Fix - GoTrue API Routing

## Problem Summary

The frontend was not correctly configured to communicate with the standalone GoTrue authentication server due to a path mismatch:

- **Supabase Client Expected:** `/auth/v1/*` endpoints (designed for full Supabase)
- **GoTrue Actually Serves:** `/*` root-level endpoints (standalone GoTrue)

## Solution Implemented

### 1. New Proxy Route: `/api/auth`

Created a new proxy at `/app/api/auth/[...path]/route.ts` that:
- Intercepts all requests from the Supabase client
- Strips the `/auth/v1/` prefix that the client adds
- Forwards requests to the standalone GoTrue server at the correct path

**Example Flow:**
```
Client calls: supabase.auth.signInWithPassword()
    ↓
Supabase client makes request to: https://app.offensivewizard.com/api/auth/auth/v1/token?grant_type=password
    ↓
Proxy strips /auth/v1/: https://auth.offensivewizard.com/token?grant_type=password
    ↓
GoTrue processes request ✅
```

### 2. Generated Valid API Keys

Created `scripts/generate-anon-key.js` to generate proper JWT-based API keys:

- **ANON Key:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ`
- **SERVICE_ROLE Key:** `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoic2VydmljZV9yb2xlIiwiaXNzIjoic3VwYWJhc2UiLCJpYXQiOjE3NjY3NDk3MzMsImV4cCI6MjA4MjEwOTczM30.-ZbBuT-rC0B2Uwpxf44wtGKTPTzKmwdSkNNAPtEXjqo`

Both keys are signed with your `SUPABASE_JWT_SECRET` and valid for 10 years.

### 3. Updated Environment Variables

**CRITICAL CHANGE:** The Supabase URL now points to the proxy, not directly to GoTrue:

```env
# OLD (WRONG):
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com

# NEW (CORRECT):
NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4. Updated Docker Configuration

Updated both `Dockerfile.nextjs` and `build-and-push-to-dockerhub.sh` to use the correct environment variables during build.

## Files Changed

1. **New Files:**
   - `/app/api/auth/[...path]/route.ts` - New proxy route
   - `/scripts/generate-anon-key.js` - Key generation script
   - `/AUTH-PROXY-FIX.md` - This documentation

2. **Updated Files:**
   - `/Dockerfile.nextjs` - Updated build args and defaults
   - `/build-and-push-to-dockerhub.sh` - Updated build command
   - `/.env.example` - Updated with correct configuration

3. **Deprecated (can be removed):**
   - `/app/supabase-proxy/[...path]/route.ts` - Old proxy (not used anymore)

## Testing the Fix

### Local Testing

1. **Update your `.env` file:**
   ```bash
   NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth
   NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ
   ```

2. **Run the development server:**
   ```bash
   npm run dev
   ```

3. **Test authentication:**
   - Go to `/login`
   - Try signing in
   - Check browser console for proxy logs (in development mode)

### Production Testing

After deploying the Docker image:

1. **Verify proxy is accessible:**
   ```bash
   curl https://app.offensivewizard.com/api/auth/health
   # Should return: {"version":"v2.184.0","name":"GoTrue"...}
   ```

2. **Test signup flow:**
   ```bash
   curl -X POST https://app.offensivewizard.com/api/auth/auth/v1/signup \
     -H "Content-Type: application/json" \
     -H "apikey: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
     -d '{"email":"test@test.com","password":"password123"}'
   ```

3. **Test in browser:**
   - Open DevTools Network tab
   - Go to login page
   - Attempt login
   - Verify requests go to `/api/auth/auth/v1/*` endpoints
   - Check for 200/400 responses (not 404)

## Deployment Instructions

### Deploy to Docker Hub

```bash
# Build and push both images
./build-and-push-to-dockerhub.sh
```

This will:
1. Build backend image: `limpet/wizardcore-backend:latest`
2. Build frontend image with correct env vars: `limpet/wizardcore-frontend:latest`
3. Push both to Docker Hub

### Update Production Environment

After deploying, ensure your production environment variables are set:

```env
# In your deployment platform (Coolify, CapRover, etc.)
NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ

# GoTrue server URL (for proxy to use)
GOTRUE_URL=https://auth.offensivewizard.com

# JWT secret (must match what was used to generate keys)
SUPABASE_JWT_SECRET=cd3e9a8225704a2b5d628a76409464bbbb7f08554ce9da8e6f350a16c7e809ccc1fec2db56c45f87abdfbfcf5b0526770a7bac12ef6e552eeebcb4f63a34501b
```

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                        Browser                               │
│                                                              │
│  Supabase Client (@supabase/supabase-js)                    │
│  - Configured with: app.offensivewizard.com/api/auth        │
│  - Automatically adds /auth/v1/ to all requests             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ HTTPS
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    Next.js Frontend                          │
│              (app.offensivewizard.com)                       │
│                                                              │
│  /api/auth/[...path]/route.ts (Proxy)                       │
│  - Receives: /api/auth/auth/v1/signup                       │
│  - Strips: /auth/v1/                                        │
│  - Forwards: /signup                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       │ HTTPS
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    GoTrue Server                             │
│              (auth.offensivewizard.com)                      │
│                                                              │
│  Endpoints at root:                                         │
│  - /signup                                                  │
│  - /token?grant_type=password                               │
│  - /user                                                    │
│  - /logout                                                  │
│  - /health                                                  │
└─────────────────────────────────────────────────────────────┘
```

## Key Benefits

1. **No Code Changes Required:** Frontend code using `@supabase/supabase-js` works without modification
2. **Proper Path Handling:** Proxy automatically handles path translation
3. **CORS Resolved:** Same-origin requests from frontend to proxy
4. **Standard JWT Auth:** Valid ANON and SERVICE_ROLE keys for GoTrue
5. **Future-Proof:** Can swap GoTrue for full Supabase later with minimal changes

## Troubleshooting

### Issue: 404 errors on auth endpoints

**Solution:** Verify proxy is working:
```bash
curl https://app.offensivewizard.com/api/auth/health
```

### Issue: Authentication fails silently

**Check:**
1. ANON key is correctly set in environment
2. JWT secret matches between key generation and GoTrue
3. Browser console for network errors
4. GoTrue server logs for auth errors

### Issue: CORS errors

**Solution:** Ensure proxy is being used (URL should be `/api/auth`, not direct to `auth.offensivewizard.com`)

## Regenerating Keys

If you need to regenerate the API keys (e.g., JWT secret changed):

```bash
# Using environment variable
SUPABASE_JWT_SECRET="your-new-secret" node scripts/generate-anon-key.js

# Or pass as argument
node scripts/generate-anon-key.js "your-new-secret"
```

Then update both:
1. `NEXT_PUBLIC_SUPABASE_ANON_KEY` in `.env`
2. Build args in `build-and-push-to-dockerhub.sh`
3. Defaults in `Dockerfile.nextjs`

## Next Steps

1. ✅ Proxy route created at `/api/auth`
2. ✅ Valid API keys generated
3. ✅ Environment variables updated
4. ✅ Docker configuration updated
5. 🔄 Build and push Docker images
6. ⏳ Deploy to production
7. ⏳ Test authentication flow
8. ⏳ Monitor for errors
