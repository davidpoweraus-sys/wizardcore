# Next.js Proxy Solution for CORS

## 🎯 The Solution

Since Coolify's Traefik proxy is blocking CORS, we're **bypassing it entirely** by routing Supabase requests through Next.js.

### How It Works:

```
Browser (offensivewizard.com)
    ↓
Next.js /api/supabase-proxy/*
    ↓ (internal Docker network - no CORS needed)
Supabase Auth (supabase-auth:9999)
    ↓
Next.js adds CORS headers
    ↓
Browser ✅
```

## ✅ Changes Made

### 1. Created Proxy Route
**File:** `app/api/supabase-proxy/[...path]/route.ts`

This Next.js API route:
- Handles OPTIONS preflight requests (returns 204 with CORS headers)
- Proxies all requests to internal Supabase Auth
- Adds proper CORS headers to responses
- Logs all requests for debugging

### 2. Updated Environment Variables
**File:** `docker-compose.prod.yml`

Changed:
```yaml
# OLD:
NEXT_PUBLIC_SUPABASE_URL: https://auth.offensivewizard.com

# NEW:
NEXT_PUBLIC_SUPABASE_URL: https://offensivewizard.com/api/supabase-proxy
SUPABASE_INTERNAL_URL: http://supabase-auth:9999
```

## 🚀 Deploy

```bash
cd /home/glbsi/Workbench/wizardcore

# Commit changes
git add app/api/supabase-proxy/
git add docker-compose.prod.yml
git add NEXTJS-PROXY-SOLUTION.md

git commit -m "Add Next.js proxy to bypass Coolify CORS issues

- Create API proxy route for Supabase requests
- Update environment variables to use proxy
- Bypass Coolify Traefik CORS limitations"

git push

# Coolify will auto-deploy or click Redeploy in UI
```

## ✅ Benefits

1. ✅ **No Coolify configuration needed** - works out of the box
2. ✅ **No Traefik label complexity** - pure Next.js solution
3. ✅ **Full CORS control** - we manage all headers
4. ✅ **Works immediately** - no waiting for Coolify fixes
5. ✅ **Better logging** - see requests in Next.js logs

## 🧪 Testing

### 1. After deployment, check browser console:

```
🚀 Registration started
📧 Email: test@example.com
🌐 Supabase URL: https://offensivewizard.com/api/supabase-proxy
🔑 API Key present: true
📤 Calling Supabase signUp...
✅ Registration successful!
```

### 2. Check Next.js logs in Coolify:

You should see:
```
🔄 Proxying request to: http://supabase-auth:9999/auth/v1/signup
✅ Proxy response status: 200
```

### 3. Network tab should show:

- Request URL: `https://offensivewizard.com/api/supabase-proxy/auth/v1/signup`
- Status: 200 OK
- Response headers include: `access-control-allow-origin: https://offensivewizard.com`

## 🔧 Troubleshooting

### Issue: 404 on /api/supabase-proxy

**Cause:** Route file not deployed or in wrong location

**Solution:** 
```bash
# Verify file exists
ls -la app/api/supabase-proxy/\[...path\]/route.ts

# Should show the file. If not, recreate it.
```

### Issue: 500 Internal Server Error

**Cause:** Can't connect to supabase-auth from Next.js

**Debug:**
```bash
# Check if both containers are on same network
docker network inspect wizardcore-network

# Should show both 'frontend' and 'supabase-auth'
```

**Solution:** Verify `SUPABASE_INTERNAL_URL` env var is set correctly in Coolify

### Issue: Still getting CORS errors

**Cause:** Browser cached old requests

**Solution:**
1. Hard refresh: Ctrl+Shift+R (or Cmd+Shift+R on Mac)
2. Clear site data: DevTools → Application → Clear storage
3. Try incognito mode

## 📊 Performance Impact

**Slight latency increase:**
- Direct: Browser → Traefik → GoTrue (~50-100ms)
- Proxied: Browser → Traefik → Next.js → GoTrue (~100-150ms)

The extra ~50ms is negligible for auth operations and worth it for reliability.

## 🔄 Reverting to Direct Connection

If Coolify CORS is fixed later, revert by:

```yaml
# In docker-compose.prod.yml:
NEXT_PUBLIC_SUPABASE_URL: https://auth.offensivewizard.com
```

And optionally delete `app/api/supabase-proxy/` directory.

## 🎯 Why This Works

1. **Same-origin requests:** Browser → offensivewizard.com (no CORS check)
2. **Server-to-server:** Next.js → supabase-auth (internal network, no CORS)
3. **We control headers:** Next.js adds CORS headers to response
4. **No external dependencies:** Works without Coolify/Traefik cooperation

## ✅ Success Criteria

After deploying, you should:
- ✅ See NO CORS errors in console
- ✅ Successfully create accounts
- ✅ See proxy logs in Next.js console
- ✅ Complete registration flow

---

**This is the most reliable solution** given Coolify's Traefik CORS limitations!
