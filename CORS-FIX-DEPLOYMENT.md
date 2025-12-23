# CORS Fix - Minimal Deployment Guide

## 🎯 Problem

You're getting CORS errors when the frontend tries to authenticate with Supabase Auth (GoTrue).

## ✅ Solution

Deploy **ONLY** the two services that handle CORS:
1. **Frontend** - Updated with correct CORS configuration
2. **Supabase Auth** - Updated with CORS headers

**You DON'T need to redeploy:**
- Backend (internal network only)
- PostgreSQL (no CORS involved)
- Redis (no CORS involved)
- Judge0 (separate domain, different CORS)

---

## 📦 What's Included

**File:** `wizardcore-cors-fix.tar.gz` (92 MB)

**Services:**
- Frontend (Next.js) - 212 MB
- Supabase Auth (GoTrue) - 50 MB

**Total:** 262 MB → Compressed to 92 MB (65% savings!)

---

## 🚀 Deployment Steps

### Option 1: Quick Fix on Coolify (Recommended)

```bash
# 1. Transfer minimal package to server
scp wizardcore-cors-fix.tar.gz user@coolify-server:/tmp/

# 2. SSH into server
ssh user@coolify-server

# 3. Load only the updated images
gunzip -c /tmp/wizardcore-cors-fix.tar.gz | docker load

# 4. Verify images loaded
docker images | grep -E "(wizardcore-frontend|gotrue)"

# 5. In Coolify UI: Restart ONLY these services
#    - Frontend
#    - Supabase-auth
```

### Option 2: Docker Compose Selective Restart

```bash
# After loading images on server
cd /path/to/wizardcore

# Restart only frontend and auth
docker-compose -f docker-compose.prod.yml up -d frontend supabase-auth

# Other services will continue running unchanged
```

---

## 🔍 Verify CORS Fix

After deployment, test CORS:

```bash
# Test Supabase Auth CORS headers
curl -I -X OPTIONS https://auth.offensivewizard.com/health \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST"

# Should see:
# Access-Control-Allow-Origin: https://offensivewizard.com
# Access-Control-Allow-Methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
# Access-Control-Allow-Credentials: true
```

Or test in browser console:

```javascript
// Test auth endpoint
fetch('https://auth.offensivewizard.com/health', {
  method: 'GET',
  credentials: 'include'
})
.then(r => r.json())
.then(d => console.log('✅ CORS working!', d))
.catch(e => console.error('❌ CORS failed:', e))
```

---

## 📋 What Changed

### Frontend
- Updated environment variables
- Latest build with correct Supabase URL

### Supabase Auth (GoTrue)
- CORS headers configured:
  ```yaml
  GOTRUE_CORS_ALLOWED_ORIGINS: "https://offensivewizard.com"
  GOTRUE_CORS_ALLOWED_HEADERS: "Authorization,Content-Type,X-Client-Info,..."
  GOTRUE_CORS_ALLOWED_METHODS: "GET,POST,PUT,PATCH,DELETE,OPTIONS"
  GOTRUE_CORS_ALLOW_CREDENTIALS: "true"
  ```

---

## ⚡ Why This is Better Than Full Stack

| Aspect | Minimal (CORS Fix) | Full Stack |
|--------|-------------------|------------|
| **File Size** | 92 MB | 3.4 GB |
| **Transfer Time** | ~1 min | ~5-15 min |
| **Load Time** | ~30 sec | ~5-10 min |
| **Services Affected** | 2 | 8 |
| **Downtime Risk** | Minimal | Higher |
| **Rollback** | Easy | Complex |

---

## 🛡️ Safety

**Other services keep running:**
- ✅ Backend continues serving API requests
- ✅ Databases maintain connections and data
- ✅ Redis keeps cache warm
- ✅ Judge0 continues executing code

**Only restarted:**
- 🔄 Frontend (fast restart)
- 🔄 Supabase Auth (fast restart)

---

## 🔄 If CORS Still Fails

### Check 1: Verify Environment Variables in Coolify

In Coolify UI, check that `supabase-auth` service has:
```
GOTRUE_CORS_ALLOWED_ORIGINS=https://offensivewizard.com
GOTRUE_CORS_ALLOW_CREDENTIALS=true
```

### Check 2: Verify Frontend is Using Correct URL

Frontend should use:
```
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
```

NOT a proxy URL.

### Check 3: Network Inspection

In browser DevTools → Network tab:
- Look for OPTIONS request before POST/GET
- Check response headers for `Access-Control-Allow-*`
- If missing, GoTrue isn't configured correctly

---

## 💡 Alternative: Use Next.js Proxy (No CORS Needed)

If CORS continues to be problematic, you can proxy Supabase through Next.js:

**Frontend Environment:**
```env
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase-proxy
```

**Add to app/api/supabase-proxy/[...path]/route.ts:**
```typescript
// Already exists in your codebase
// Proxies all Supabase requests through Next.js
// Avoids CORS entirely
```

This avoids CORS because the browser talks to the same domain (offensivewizard.com).

---

## 📊 Comparison: Deploy Strategy

### Minimal CORS Fix (This Guide)
```
✅ Pros:
  - Fast (92 MB)
  - Low risk
  - Targeted fix
  - Other services unaffected

❌ Cons:
  - Only fixes CORS issue
  - If other services need updates, requires another deployment
```

### Full Stack (wizardcore-complete-stack.tar.gz)
```
✅ Pros:
  - Everything updated
  - One-time deployment
  - All images version-locked

❌ Cons:
  - Large file (3.4 GB)
  - Longer downtime
  - Overkill for CORS fix
  - Restarts all services
```

---

## 🎯 Recommendation

**For CORS fix only:** Use this minimal package (92 MB)

**For full deployment or major updates:** Use complete stack (3.4 GB)

---

## ✅ Quick Command Reference

```bash
# Transfer minimal package
scp wizardcore-cors-fix.tar.gz user@server:/tmp/

# Load images
ssh user@server
gunzip -c /tmp/wizardcore-cors-fix.tar.gz | docker load

# Restart services
docker-compose -f docker-compose.prod.yml up -d frontend supabase-auth

# Test CORS
curl -I -X OPTIONS https://auth.offensivewizard.com/health \
  -H "Origin: https://offensivewizard.com"
```

---

## 🎉 Expected Result

After deployment:
- ✅ No more CORS errors in browser console
- ✅ Login/signup works
- ✅ Authentication flows complete
- ✅ All other services continue running

Total downtime: **~30 seconds** (just frontend + auth restart)

---

**Ready to fix CORS? Transfer `wizardcore-cors-fix.tar.gz` and follow the steps above!**
