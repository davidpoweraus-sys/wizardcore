# Coolify CORS Configuration Fix

## 🚨 The Real Problem

You're using **Coolify** which has its own reverse proxy (Traefik/Caddy). The nginx.conf file we edited is **not being used**. 

The CORS preflight requests are being **blocked by Coolify's reverse proxy** before they even reach your GoTrue service.

## ✅ Solution: Configure Coolify Reverse Proxy

### Option 1: Add CORS Middleware in Coolify (RECOMMENDED)

If Coolify uses **Traefik** (most likely):

1. **Go to Coolify Dashboard**
2. **Navigate to your `supabase-auth` service** (the one on `auth.offensivewizard.com`)
3. **Find "Labels" or "Custom Configuration"** section
4. **Add these Traefik labels:**

```yaml
# In Coolify UI, add these as custom labels:
traefik.http.middlewares.auth-cors.headers.accessControlAllowOriginList=https://offensivewizard.com
traefik.http.middlewares.auth-cors.headers.accessControlAllowCredentials=true
traefik.http.middlewares.auth-cors.headers.accessControlAllowMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS
traefik.http.middlewares.auth-cors.headers.accessControlAllowHeaders=Authorization,Content-Type,X-Client-Info,X-Requested-With,apikey,x-client-info,x-supabase-api-version,Accept,Accept-Language,Content-Language
traefik.http.middlewares.auth-cors.headers.accessControlExposeHeaders=X-Total-Count
traefik.http.middlewares.auth-cors.headers.accessControlMaxAge=86400

# Apply the middleware to your auth service
traefik.http.routers.auth-offensivewizard.middlewares=auth-cors
```

5. **Save and redeploy** the service

### Option 2: Add Labels to docker-compose.prod.yml

Add labels directly to the `supabase-auth` service:

```yaml
services:
  supabase-auth:
    image: supabase/gotrue:v2.184.0
    labels:
      # CORS middleware
      - "traefik.http.middlewares.auth-cors.headers.accessControlAllowOriginList=https://offensivewizard.com"
      - "traefik.http.middlewares.auth-cors.headers.accessControlAllowCredentials=true"
      - "traefik.http.middlewares.auth-cors.headers.accessControlAllowMethods=GET,POST,PUT,PATCH,DELETE,OPTIONS"
      - "traefik.http.middlewares.auth-cors.headers.accessControlAllowHeaders=Authorization,Content-Type,X-Client-Info,X-Requested-With,apikey,x-client-info,x-supabase-api-version,Accept,Accept-Language,Content-Language"
      - "traefik.http.middlewares.auth-cors.headers.accessControlExposeHeaders=X-Total-Count"
      - "traefik.http.middlewares.auth-cors.headers.accessControlMaxAge=86400"
      # Apply middleware to router (replace 'auth-offensivewizard' with your actual router name)
      - "traefik.http.routers.${ROUTER_NAME:-auth}.middlewares=auth-cors"
    environment:
      # ... existing environment variables
```

### Option 3: If Coolify Uses Caddy

If Coolify uses Caddy instead, you need to add a `Caddyfile` configuration:

1. In Coolify dashboard, find the **Caddy configuration** for `auth.offensivewizard.com`
2. Add CORS headers:

```caddy
auth.offensivewizard.com {
    header {
        Access-Control-Allow-Origin "https://offensivewizard.com"
        Access-Control-Allow-Credentials "true"
        Access-Control-Allow-Methods "GET, POST, PUT, PATCH, DELETE, OPTIONS"
        Access-Control-Allow-Headers "Authorization, Content-Type, X-Client-Info, X-Requested-With, apikey, x-client-info, x-supabase-api-version, Accept, Accept-Language, Content-Language"
        Access-Control-Expose-Headers "X-Total-Count"
        Access-Control-Max-Age "86400"
    }
    
    # Handle OPTIONS preflight
    @options {
        method OPTIONS
    }
    handle @options {
        respond 204
    }
    
    reverse_proxy supabase-auth:9999
}
```

## 🔍 How to Check Which Proxy Coolify Uses

### Check Traefik:
```bash
# SSH into your Coolify server
docker ps | grep traefik
```

If you see Traefik container → Use Option 1 or 2

### Check Caddy:
```bash
# SSH into your Coolify server  
docker ps | grep caddy
```

If you see Caddy container → Use Option 3

## 🧪 Testing After Configuration

### 1. From your browser:
```javascript
// Open browser console on https://offensivewizard.com
fetch('https://auth.offensivewizard.com/health', {
  method: 'GET',
  headers: {
    'Origin': 'https://offensivewizard.com'
  }
}).then(r => r.json()).then(console.log)
```

Should show: `{version: "...", name: "GoTrue"}` with no CORS errors

### 2. Test OPTIONS preflight:
```bash
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: apikey,content-type"
```

**Expected response:**
```
< HTTP/2 204
< access-control-allow-origin: https://offensivewizard.com
< access-control-allow-credentials: true
< access-control-allow-methods: GET,POST,PUT,PATCH,DELETE,OPTIONS
```

## 🚀 Quick Fix Without Coolify Access

If you can't easily modify Coolify's proxy config, you have two options:

### Workaround A: Use GoTrue Directly (No Proxy)

**Temporarily expose GoTrue directly:**

1. In `docker-compose.prod.yml`, change supabase-auth ports:
```yaml
supabase-auth:
  ports:
    - "9999:9999"  # Already exposed
```

2. Point DNS `auth.offensivewizard.com` directly to `your-server-ip:9999`

3. **⚠️ Not recommended for production** (no SSL termination, no load balancing)

### Workaround B: Add CORS Proxy Middleware in Next.js

Add a Next.js API route that proxies requests to GoTrue:

**Create:** `app/api/auth/[...path]/route.ts`

```typescript
import { NextRequest, NextResponse } from 'next/server'

const GOTRUE_URL = process.env.SUPABASE_URL || 'http://supabase-auth:9999'

export async function GET(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

export async function POST(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

export async function OPTIONS(request: NextRequest) {
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': 'https://offensivewizard.com',
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET,POST,PUT,PATCH,DELETE,OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization,Content-Type,X-Client-Info,apikey,x-supabase-api-version',
      'Access-Control-Max-Age': '86400',
    },
  })
}

async function proxyRequest(request: NextRequest, path: string[]) {
  const url = new URL(request.url)
  const targetUrl = `${GOTRUE_URL}/${path.join('/')}${url.search}`

  const headers = new Headers(request.headers)
  headers.delete('host')

  const response = await fetch(targetUrl, {
    method: request.method,
    headers: headers,
    body: request.body,
  })

  const data = await response.text()

  return new NextResponse(data, {
    status: response.status,
    headers: {
      'Content-Type': response.headers.get('Content-Type') || 'application/json',
      'Access-Control-Allow-Origin': 'https://offensivewizard.com',
      'Access-Control-Allow-Credentials': 'true',
    },
  })
}
```

**Then update** `NEXT_PUBLIC_SUPABASE_URL`:
```env
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/auth
```

## 📋 Recommended Approach

**Best → Worst:**

1. ✅ **Configure Coolify's reverse proxy** (Option 1) - Proper, production-ready
2. ⚠️ **Add Next.js proxy** (Workaround B) - Works but adds latency
3. ❌ **Expose GoTrue directly** (Workaround A) - Not secure for production

## 🔧 Step-by-Step: Traefik Labels in Coolify

Since most Coolify instances use Traefik, here's the exact process:

1. **Login to Coolify** dashboard
2. **Go to:** Projects → Your Project → `supabase-auth` service
3. **Scroll to:** "Advanced" or "Labels" section
4. **Click:** "Add Label"
5. **Add each label** one by one:

| Key | Value |
|-----|-------|
| `traefik.http.middlewares.auth-cors.headers.accessControlAllowOriginList` | `https://offensivewizard.com` |
| `traefik.http.middlewares.auth-cors.headers.accessControlAllowCredentials` | `true` |
| `traefik.http.middlewares.auth-cors.headers.accessControlAllowMethods` | `GET,POST,PUT,PATCH,DELETE,OPTIONS` |
| `traefik.http.middlewares.auth-cors.headers.accessControlAllowHeaders` | `Authorization,Content-Type,X-Client-Info,X-Requested-With,apikey,x-client-info,x-supabase-api-version` |

6. **Apply middleware:** Find the router label (usually auto-generated like `traefik.http.routers.<service-name>.middlewares`)
   - If it doesn't exist, create: `traefik.http.routers.auth.middlewares` = `auth-cors`
   - If it exists with other middlewares: append `,auth-cors` to the value

7. **Save** and **Restart** the service

## ✅ Verification

After applying the fix, run this in your browser console on `https://offensivewizard.com`:

```javascript
// Test 1: Simple fetch
fetch('https://auth.offensivewizard.com/health')
  .then(r => r.json())
  .then(data => console.log('✅ Health check:', data))
  .catch(err => console.error('❌ Failed:', err))

// Test 2: CORS with credentials
fetch('https://auth.offensivewizard.com/health', {
  credentials: 'include',
  headers: {
    'Content-Type': 'application/json'
  }
})
  .then(r => r.json())
  .then(data => console.log('✅ CORS test:', data))
  .catch(err => console.error('❌ CORS failed:', err))
```

Both should succeed without CORS errors!

## 📞 Need Help?

If you're still stuck, provide:
1. Screenshot of Coolify service configuration
2. Output of: `docker ps | grep -E "traefik|caddy"`
3. Browser Network tab showing the failed OPTIONS request
