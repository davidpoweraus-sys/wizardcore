# Environment Configuration Guide

## Overview

WizardCore uses **ONE** `.env` file for all configuration. This guide explains each variable and how they work together.

## File Structure

```
wizardcore/
├── .env                  # YOUR actual configuration (git-ignored)
├── .env.example         # Template with all variables documented
└── docker-compose.yml   # References these variables
```

## Critical Concept: Public vs Internal URLs

### Public URLs (NEXT_PUBLIC_*)
- Used by the **browser** (client-side)
- Must be accessible from the internet
- Use HTTPS in production
- Example: `https://auth.offensivewizard.com`

### Internal URLs (no NEXT_PUBLIC_ prefix)
- Used by **server-side** code and Docker containers
- Use Docker network names
- Use HTTP (containers communicate internally)
- Example: `http://supabase-auth:9999`

## How the Proxy System Works

### 1. Auth Proxy
```
Browser → https://auth.offensivewizard.com
         ↓ (via Next.js route)
Server → http://supabase-auth:9999 (internal Docker network)
```

**Variables:**
- `NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com` (browser uses this)
- `GOTRUE_URL=http://supabase-auth:9999` (server proxy uses this)

### 2. Backend API Proxy
```
Browser → /api/backend/v1/users/me/stats
         ↓ (via Next.js route)
Server → http://backend:8080/v1/users/me/stats (internal Docker network)
```

**Variables:**
- `NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com` (reference only)
- `BACKEND_URL=http://backend:8080` (server proxy uses this)

### 3. Judge0 Proxy
```
Browser → /api/judge0/submissions
         ↓ (via Next.js route)
Server → http://judge0:2358/submissions (internal Docker network)
```

**Variables:**
- `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com` (reference only)
- `JUDGE0_URL=http://judge0:2358` (server proxy uses this)

## Complete Environment Variables

### Frontend Configuration
```bash
# Public URLs (used by browser)
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbG...  # JWT token
NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com
NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
NEXT_PUBLIC_SITE_URL=https://app.offensivewizard.com

# Internal URLs (used by Next.js server-side proxies)
GOTRUE_URL=http://supabase-auth:9999
SUPABASE_INTERNAL_URL=http://supabase-auth:9999
BACKEND_URL=http://backend:8080
JUDGE0_URL=http://judge0:2358
```

### Backend Configuration
```bash
PORT=8080
ENVIRONMENT=production
DATABASE_URL=postgresql://wizardcore:PASSWORD@postgres:5432/wizardcore?sslmode=disable
REDIS_URL=redis://:PASSWORD@redis:6379
CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
```

### GoTrue/Auth Configuration
```bash
# JWT Secret (CRITICAL - must match on frontend and backend!)
SUPABASE_JWT_SECRET=your-secret-here
GOTRUE_JWT_SECRET=your-secret-here  # Must be identical to SUPABASE_JWT_SECRET

# Service role key (server-side only, never expose to frontend)
SUPABASE_SERVICE_ROLE_KEY=eyJhbG...

# GoTrue settings
GOTRUE_SITE_URL=https://app.offensivewizard.com
API_EXTERNAL_URL=https://auth.offensivewizard.com
GOTRUE_CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
GOTRUE_DB_AUTOMIGRATE=true
```

### Database Configuration
```bash
# Main Wizardcore database
POSTGRES_USER=wizardcore
POSTGRES_PASSWORD=your-password
POSTGRES_DB=wizardcore

# Supabase auth database
SUPABASE_POSTGRES_USER=supabase_auth_admin
SUPABASE_POSTGRES_PASSWORD=your-password
SUPABASE_POSTGRES_DB=supabase_auth

# Judge0 database
JUDGE0_POSTGRES_USER=judge0
JUDGE0_POSTGRES_PASSWORD=your-password
JUDGE0_POSTGRES_DB=judge0
```

### Redis Configuration
```bash
REDIS_PASSWORD=your-password
JUDGE0_REDIS_PASSWORD=your-password
```

### Judge0 Configuration
```bash
JUDGE0_API_KEY=your-api-key
JUDGE0_API_URL=http://judge0:2358  # Internal Docker network
```

## Common Issues

### 1. CORS Errors
**Symptom:** Browser console shows "CORS policy" errors

**Solution:** Make sure:
- `CORS_ALLOWED_ORIGINS` includes your frontend domain
- `GOTRUE_CORS_ALLOWED_ORIGINS` includes your frontend domain
- Internal URLs (`BACKEND_URL`, `JUDGE0_URL`, `GOTRUE_URL`) use Docker network names

### 2. "Cannot connect to backend"
**Symptom:** Dashboard fails to load data

**Cause:** Server-side proxy can't reach backend service

**Solution:**
```bash
# Check if backend is running
docker-compose ps backend

# Verify the backend URL is correct
docker-compose exec frontend env | grep BACKEND_URL
# Should show: BACKEND_URL=http://backend:8080

# Test connectivity
docker-compose exec frontend ping backend
```

### 3. Auth not working
**Symptom:** Can't log in or session expires immediately

**Cause:** JWT secrets don't match

**Solution:**
```bash
# Verify JWT secrets match
grep GOTRUE_JWT_SECRET .env
grep SUPABASE_JWT_SECRET .env
# Both should be identical!
```

## Environment Variable Priority

When a variable is set in multiple places, Docker uses this priority:
1. Docker Compose command line (`-e` flag)
2. `.env` file in project root
3. Default value in `docker-compose.yml`

**Best practice:** Put everything in `.env` file for consistency.

## Security Checklist

- [ ] `NEXT_PUBLIC_*` variables are safe to expose (they're public)
- [ ] `SUPABASE_SERVICE_ROLE_KEY` is NEVER in frontend code
- [ ] All passwords are strong and unique
- [ ] JWT secret is at least 64 characters
- [ ] `.env` file is in `.gitignore`
- [ ] Production uses HTTPS for all public URLs

## Testing Your Configuration

### 1. Verify environment variables are loaded
```bash
docker-compose config | grep NEXT_PUBLIC_SUPABASE_URL
```

### 2. Check frontend can reach backend (inside container)
```bash
docker-compose exec frontend wget -O- http://backend:8080/health
```

### 3. Check auth proxy
```bash
docker-compose exec frontend wget -O- http://supabase-auth:9999/health
```

### 4. Check from browser
Open DevTools Console and run:
```javascript
console.log({
  SUPABASE_URL: process.env.NEXT_PUBLIC_SUPABASE_URL,
  BACKEND_URL: process.env.NEXT_PUBLIC_BACKEND_URL,
  JUDGE0_URL: process.env.NEXT_PUBLIC_JUDGE0_API_URL
});
```

## Migration from Multiple .env Files

If you have multiple `.env*` files:

1. **Consolidate into one `.env` file:**
   ```bash
   cat .env.local .env.production > .env.consolidated
   # Review and remove duplicates
   mv .env.consolidated .env
   ```

2. **Remove old files:**
   ```bash
   rm .env.local .env.production .env.development
   ```

3. **Update `.env` with values from `.env.example`**

4. **Rebuild and deploy:**
   ```bash
   ./build-and-push.sh
   ```

## Quick Reference

| Variable | Type | Used By | Example |
|----------|------|---------|---------|
| `NEXT_PUBLIC_SUPABASE_URL` | Public | Browser | `https://auth.offensivewizard.com` |
| `GOTRUE_URL` | Internal | Server proxy | `http://supabase-auth:9999` |
| `NEXT_PUBLIC_BACKEND_URL` | Public | Browser (ref) | `https://api.offensivewizard.com` |
| `BACKEND_URL` | Internal | Server proxy | `http://backend:8080` |
| `NEXT_PUBLIC_JUDGE0_API_URL` | Public | Browser (ref) | `https://judge0.offensivewizard.com` |
| `JUDGE0_URL` | Internal | Server proxy | `http://judge0:2358` |

## Support

If you're still having issues:
1. Check [`CORS-FIX-COMPLETE.md`](CORS-FIX-COMPLETE.md) for technical details
2. Review [`DEPLOYMENT-CHECKLIST.md`](DEPLOYMENT-CHECKLIST.md) for deployment steps
3. Compare your `.env` with `.env.example`