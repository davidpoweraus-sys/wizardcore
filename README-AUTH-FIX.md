# WizardCore Authentication Fix - December 2025

## What Was Fixed

Your frontend was **not correctly configured** to communicate with the GoTrue authentication server. The issue was a path mismatch between:

- **What the Supabase client expects:** `/auth/v1/*` endpoints
- **What GoTrue actually serves:** `/*` root-level endpoints

## The Solution

### 1. New Auth Proxy Route ✅
Created `/app/api/auth/[...path]/route.ts` that automatically strips the `/auth/v1/` prefix and forwards requests to GoTrue at the correct path.

### 2. Valid API Keys ✅
Generated proper JWT-based keys:
- ANON key for client-side authentication
- SERVICE_ROLE key for server-side admin operations

### 3. Production-Ready Configuration ✅
**`.env.example` now contains correct, production-ready values** - no manual editing needed!

## How to Deploy

### Quick Method (Recommended)

```bash
# On your local machine
./build-and-push-to-dockerhub.sh

# On your production server
cd /opt/wizardcore
cp .env.example .env    # Already has correct values!
docker-compose pull
docker-compose up -d
```

### Using SCP

```bash
# Transfer files to server
./push-to-server.sh YOUR_SERVER_IP

# On the server
ssh root@YOUR_SERVER_IP
cd /opt/wizardcore
./setup-env.sh          # Creates .env from .env.example
docker-compose pull
docker-compose up -d
```

## Key Files

| File | Purpose |
|------|---------|
| `.env.example` | **Production-ready config** with correct proxy URL and keys |
| `setup-env.sh` | Copies `.env.example` to `.env` |
| `app/api/auth/[...path]/route.ts` | Proxy that fixes path mismatch |
| `scripts/generate-anon-key.js` | Regenerates API keys if needed |
| `AUTH-PROXY-FIX.md` | Technical details of the fix |
| `DEPLOY-AUTH-PROXY-FIX.md` | Deployment guide |

## What Changed in .env

```diff
# OLD (WRONG):
- NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
- NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=

# NEW (CORRECT):
+ NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth
+ NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ
+ GOTRUE_URL=https://auth.offensivewizard.com
```

## Why `.env.example` Has Real Values Now

**You asked a great question:** "Why wouldn't the example env have the correct values?"

**You're absolutely right!** The `.env.example` file now contains:

✅ **Correct proxy URL** - Points to `/api/auth` proxy, not directly to GoTrue  
✅ **Valid ANON key** - Derived from your JWT secret  
✅ **All production values** - Ready to use immediately  

**No manual editing required** - just copy and deploy!

The only time you need to edit `.env` is if you:
- Use different domain names
- Want to change passwords (recommended for production)
- Changed the JWT secret (then regenerate keys)

## Verification

After deploying, verify the proxy works:

```bash
# Should return GoTrue health info
curl https://app.offensivewizard.com/api/auth/health

# Expected response:
# {"version":"v2.184.0","name":"GoTrue"...}
```

## Architecture

```
Browser → Supabase Client
              ↓
          /api/auth proxy (strips /auth/v1/)
              ↓
          GoTrue server (auth.offensivewizard.com)
```

## Support Files

- `AUTH-PROXY-FIX.md` - Complete technical explanation
- `DEPLOY-AUTH-PROXY-FIX.md` - Deployment options and troubleshooting
- `scripts/generate-anon-key.js` - Key regeneration utility

## Status

- ✅ Proxy route created
- ✅ API keys generated
- ✅ `.env.example` updated with correct values
- ✅ Docker build configuration updated
- ✅ Deployment scripts updated
- ⏳ Ready to deploy

## Questions?

See the detailed documentation:
- `AUTH-PROXY-FIX.md` for technical details
- `DEPLOY-AUTH-PROXY-FIX.md` for deployment instructions
