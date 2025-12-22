# ACTUAL CORS Fix - The Real Root Cause

## 🔴 The REAL Problem

Your CORS error is happening because:

1. **Coolify's Traefik reverse proxy** is handling `auth.offensivewizard.com`
2. **Traefik has NO CORS configuration** 
3. When your browser sends OPTIONS preflight → Traefik blocks it (no CORS headers)
4. Request **never reaches GoTrue** (your auth service)
5. Browser sees failed preflight → Shows CORS error

**The nginx.conf we edited earlier is NOT being used** because you're using Coolify!

## ✅ The Solution

I've added **Traefik labels** to your `docker-compose.prod.yml` that will configure CORS at the reverse proxy level.

### What Changed:

**File:** `docker-compose.prod.yml`

Added these labels to `supabase-auth` service:
```yaml
labels:
  - "traefik.http.middlewares.auth-cors.headers.accessControlAllowOriginList=https://offensivewizard.com"
  - "traefik.http.middlewares.auth-cors.headers.accessControlAllowCredentials=true"
  - "traefik.http.middlewares.auth-cors.headers.accessControlAllowMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS"
  - "traefik.http.middlewares.auth-cors.headers.accessControlAllowHeaders=Authorization,Content-Type,X-Client-Info,X-Requested-With,apikey,x-client-info,x-supabase-api-version,Accept,Accept-Language,Content-Language"
  - "traefik.http.middlewares.auth-cors.headers.accessControlExposeHeaders=X-Total-Count"
  - "traefik.http.middlewares.auth-cors.headers.accessControlMaxAge=86400"
```

## 🚀 Deploy Steps

### Step 1: Commit and Push Changes

```bash
cd /home/glbsi/Workbench/wizardcore

git add docker-compose.prod.yml
git add app/\(auth\)/register/page.tsx  # Has logging
git add COOLIFY-CORS-FIX.md
git add ACTUAL-CORS-FIX.md

git commit -m "Fix CORS by adding Traefik labels for Coolify reverse proxy

- Add Traefik CORS middleware labels to supabase-auth service
- Add detailed logging to registration page
- Document Coolify-specific CORS configuration"

git push
```

### Step 2: Redeploy in Coolify

1. Go to **Coolify Dashboard**
2. Find your **WizardCore** project
3. Click **"Redeploy"** or wait for auto-deploy

### Step 3: Apply Middleware (IMPORTANT!)

The labels define the CORS middleware, but you need to **apply it to the router**.

**Option A: In Coolify UI (Recommended)**

1. Go to **supabase-auth** service settings
2. Find **"Labels"** or **"Advanced"** section
3. Add one more label:
   - **Key:** `traefik.http.routers.supabase-auth.middlewares`
   - **Value:** `auth-cors`
   
   *(Replace `supabase-auth` with the actual router name Coolify generated)*

4. **Save** and **Restart** the service

**Option B: Check Auto-Generated Router Name**

If Coolify auto-generated a router name, you need to find it:

```bash
# SSH into Coolify server
docker inspect <supabase-auth-container> | grep traefik.http.routers
```

Look for a label like: `traefik.http.routers.wizardcore-supabase-auth-xyz`

Then in Coolify, add label:
- **Key:** `traefik.http.routers.wizardcore-supabase-auth-xyz.middlewares`
- **Value:** `auth-cors`

### Step 4: Verify Deployment

```bash
# Test CORS preflight
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: apikey,content-type"
```

**Expected Response:**
```
< HTTP/2 204
< access-control-allow-origin: https://offensivewizard.com
< access-control-allow-credentials: true
< access-control-allow-methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
< access-control-allow-headers: ...
```

### Step 5: Test Registration

1. Open **https://offensivewizard.com/register**
2. Open **DevTools** (F12) → **Console** tab
3. Fill form and submit
4. Check for my emoji logs:
   ```
   🚀 Registration started
   📧 Email: your@email.com
   📤 Calling Supabase signUp...
   ✅ Registration successful!
   ```

## 🐛 If It Still Doesn't Work

### Issue: Labels not being recognized

**Cause:** Coolify might override labels

**Solution:** Add labels **directly in Coolify UI** instead of docker-compose:

1. Coolify → supabase-auth → Labels
2. Add each label manually
3. Restart service

### Issue: Middleware not applied

**Symptom:** Headers still show no CORS headers

**Debug:**
```bash
# Check Traefik config
docker exec <traefik-container> cat /etc/traefik/traefik.yml

# Check if middleware exists
docker logs <traefik-container> | grep auth-cors
```

**Solution:** The middleware label might need to be applied to the router. See Step 3 above.

### Issue: Wrong router name

**Symptom:** Middleware created but not applied

**Solution:** Find the actual router name:

```bash
# List all Traefik routers
docker exec <traefik-container> wget -qO- http://localhost:8080/api/http/routers | jq
```

Look for the router handling `auth.offensivewizard.com`, note its name, then apply middleware to that router.

## 🔧 Alternative: Next.js Proxy (If Traefik Fix Doesn't Work)

If you can't get Traefik labels working, use this workaround:

**Create:** `/home/glbsi/Workbench/wizardcore/app/api/supabase/[...path]/route.ts`

```typescript
import { NextRequest, NextResponse } from 'next/server'

const GOTRUE_URL = 'http://supabase-auth:9999'

export async function OPTIONS(request: NextRequest) {
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': 'https://offensivewizard.com',
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET,POST,PUT,PATCH,DELETE,OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization,Content-Type,apikey,x-client-info,x-supabase-api-version',
      'Access-Control-Max-Age': '86400',
    },
  })
}

export async function GET(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path)
}

export async function POST(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path)
}

async function proxy(req: NextRequest, path: string[]) {
  const url = new URL(req.url)
  const target = `${GOTRUE_URL}/${path.join('/')}${url.search}`
  
  const headers = new Headers(req.headers)
  headers.delete('host')
  
  const res = await fetch(target, {
    method: req.method,
    headers,
    body: req.body,
  })
  
  const data = await res.text()
  
  return new NextResponse(data, {
    status: res.status,
    headers: {
      'Content-Type': res.headers.get('Content-Type') || 'application/json',
      'Access-Control-Allow-Origin': 'https://offensivewizard.com',
      'Access-Control-Allow-Credentials': 'true',
    },
  })
}
```

**Then update in Coolify environment:**
```
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase
```

This proxies all Supabase requests through Next.js, avoiding CORS entirely.

## 📋 Summary of Files Changed

1. ✅ `docker-compose.prod.yml` - Added Traefik CORS labels
2. ✅ `app/(auth)/register/page.tsx` - Added debug logging
3. ✅ `COOLIFY-CORS-FIX.md` - Detailed Coolify instructions
4. ✅ `ACTUAL-CORS-FIX.md` - This quick reference

## ✅ Success Criteria

After deploying, you should see:

1. ✅ Browser console shows my emoji logs
2. ✅ No CORS errors in console
3. ✅ Network tab shows OPTIONS request returns 204
4. ✅ Registration completes successfully
5. ✅ User redirected to dashboard

---

**Current Status:** Ready to deploy

**Next Step:** Push to git, redeploy in Coolify, add router middleware label
