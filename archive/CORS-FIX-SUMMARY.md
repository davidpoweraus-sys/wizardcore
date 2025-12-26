# CORS Error Fix - Root Cause Analysis & Resolution

## 🔍 Root Cause Identified

Your CORS errors were caused by **problematic nginx configuration** that was handling CORS headers incorrectly. Specifically:

### The Problems:

1. **Header Duplication/Conflict** (nginx.conf:31-48)
   - CORS headers were being added twice: once for all requests and again in the `if` block for OPTIONS
   - When using `add_header` inside an `if` block, nginx **discards** all headers from parent contexts
   - This caused the OPTIONS preflight responses to have incomplete or conflicting headers

2. **Insecure Origin Handling**
   - Using `$http_origin` directly without validation is a security risk
   - Any origin could trigger CORS responses, defeating the purpose of CORS

3. **Missing Request Headers**
   - The browser was sending `x-supabase-api-version` header (per your request log)
   - This wasn't explicitly listed in the allowed headers

## 🔧 What I Fixed

### 1. **Nginx Configuration** (`nginx-config/nginx.conf`)

#### Before:
```nginx
# CORS headers for all responses
add_header Access-Control-Allow-Origin $http_origin always;
# ... more headers

# Handle preflight requests
if ($request_method = 'OPTIONS') {
    add_header Access-Control-Allow-Origin $http_origin always;
    # ... duplicate headers - this REPLACES all parent headers!
    return 204;
}
```

#### After:
```nginx
# Secure origin whitelist using map
map $http_origin $cors_origin {
    default "";
    "~^https://offensivewizard\.com$" "$http_origin";
    "~^http://localhost(:\d+)?$" "$http_origin";
}

location / {
    # Handle OPTIONS at the top - returns immediately
    if ($request_method = 'OPTIONS') {
        add_header Access-Control-Allow-Origin $cors_origin always;
        add_header Access-Control-Allow-Credentials 'true' always;
        add_header Access-Control-Allow-Methods 'GET, POST, PUT, PATCH, DELETE, OPTIONS' always;
        add_header Access-Control-Allow-Headers 'Authorization, Content-Type, X-Client-Info, X-Requested-With, apikey, x-client-info, x-supabase-api-version, Accept, Accept-Language, Content-Language' always;
        add_header Access-Control-Max-Age '86400' always;
        add_header Content-Type 'text/plain charset=UTF-8' always;
        add_header Content-Length 0 always;
        return 204;
    }

    # CORS for actual requests (after OPTIONS)
    add_header Access-Control-Allow-Origin $cors_origin always;
    add_header Access-Control-Allow-Credentials 'true' always;
    add_header Access-Control-Expose-Headers 'X-Total-Count' always;
    
    proxy_pass http://supabase_auth;
    # ... proxy headers
}
```

**Key improvements:**
- ✅ Origin whitelist via `map` directive for security
- ✅ Single source of truth for CORS headers
- ✅ OPTIONS requests handled first to avoid header conflicts
- ✅ All required headers included (including `x-supabase-api-version`)
- ✅ Proper `always` flag ensures headers are added even on error responses

### 2. **Enhanced Logging** (`app/(auth)/register/page.tsx`)

Added comprehensive console logging to track:
- Request initiation
- Environment variables being used
- Supabase API responses
- Error details (status, message, name)
- Success confirmations

This will help you see exactly where the request fails in your browser console.

## 📊 How to Test

### Option 1: Automated Test Script

Run the diagnostic script I created:
```bash
cd /home/glbsi/Workbench/wizardcore
./test-cors-detailed.sh
```

This will check:
- ✅ CORS preflight responses
- ✅ Required headers presence
- ✅ Container status
- ✅ Configuration files
- ✅ SSL certificates

### Option 2: Manual Browser Test

1. **Rebuild and deploy** your services:
   ```bash
   cd /home/glbsi/Workbench/wizardcore
   docker-compose -f docker-compose.prod.yml down
   docker-compose -f docker-compose.prod.yml up -d --build
   ```

2. **Open browser** to https://offensivewizard.com/register

3. **Open DevTools** (F12):
   - Go to **Console** tab (to see my new logs)
   - Go to **Network** tab (to see actual headers)
   - Filter by "signup" to see the request

4. **Fill out registration form** and submit

5. **Check Console** for detailed logs:
   ```
   🚀 Registration started
   📧 Email: test@example.com
   🌐 Supabase URL: https://auth.offensivewizard.com
   🔑 API Key present: true
   🔗 Redirect URL: https://offensivewizard.com/auth/callback
   📤 Calling Supabase signUp...
   ```

6. **Check Network tab**:
   - Look for OPTIONS request to `/auth/v1/signup`
   - Response should be **204 No Content**
   - Response headers should include:
     ```
     access-control-allow-origin: https://offensivewizard.com
     access-control-allow-credentials: true
     access-control-allow-methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
     access-control-allow-headers: Authorization, Content-Type, ...
     ```

### Option 3: cURL Test

Test the preflight manually:
```bash
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: apikey,authorization,content-type,x-client-info,x-supabase-api-version"
```

**Expected output:**
```
< HTTP/2 204
< access-control-allow-origin: https://offensivewizard.com
< access-control-allow-credentials: true
< access-control-allow-methods: GET, POST, PUT, PATCH, DELETE, OPTIONS
< access-control-allow-headers: Authorization, Content-Type, X-Client-Info, ...
```

## 🚀 Deployment Steps

1. **Update nginx container** (if using separate nginx container):
   ```bash
   docker-compose -f docker-compose.prod.yml up -d --force-recreate nginx
   ```
   Or if nginx is configured via Coolify, update the config there.

2. **Restart all services** to ensure changes take effect:
   ```bash
   docker-compose -f docker-compose.prod.yml restart
   ```

3. **Verify services are healthy**:
   ```bash
   docker ps
   docker logs supabase-auth --tail 50
   # Check for any error messages
   ```

4. **Test registration** via browser

## 🐛 Troubleshooting

### Issue: Still seeing CORS errors after update

**Check:**
1. Is nginx using the new config?
   ```bash
   docker exec <nginx-container> cat /etc/nginx/nginx.conf | grep "map \$http_origin"
   ```
   Should see the origin map directive.

2. Did you restart nginx after config change?
   ```bash
   docker-compose -f docker-compose.prod.yml restart
   ```

3. Browser cache? Try in incognito mode.

### Issue: "Origin not allowed"

**Check:** 
- The origin map in nginx.conf includes your domain
- GoTrue GOTRUE_CORS_ALLOWED_ORIGINS includes your domain
- You're using the exact origin (check for trailing slashes, http vs https)

### Issue: Headers missing

**Check nginx logs:**
```bash
docker logs <nginx-container> | tail -50
```

Look for errors like:
- `upstream prematurely closed connection`
- `no live upstreams`

This means GoTrue isn't responding. Check:
```bash
docker logs supabase-auth | tail -50
```

### Issue: Self-signed certificate errors

In production, you should use **real SSL certificates** (Let's Encrypt via Coolify or certbot).

For testing, you can bypass in browser:
- Chrome/Edge: Type `thisisunsafe` when on the warning page
- Firefox: Click "Advanced" → "Accept Risk and Continue"

## 📝 Files Changed

1. ✅ `/home/glbsi/Workbench/wizardcore/nginx-config/nginx.conf` - Fixed CORS handling
2. ✅ `/home/glbsi/Workbench/wizardcore/app/(auth)/register/page.tsx` - Added logging
3. ✅ `/home/glbsi/Workbench/wizardcore/test-cors-detailed.sh` - New diagnostic script

## 🔐 Security Notes

- ✅ Origin whitelist prevents unauthorized domains
- ✅ Credentials flag requires specific origins (not `*`)
- ✅ ANON_KEY is safe to expose (read-only public operations)
- ✅ JWT_SECRET remains server-side only
- ✅ SameSite=Lax provides CSRF protection

## 📚 Additional Resources

- [MDN CORS Guide](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
- [Nginx CORS Configuration](https://enable-cors.org/server_nginx.html)
- [Supabase Auth SSR](https://supabase.com/docs/guides/auth/server-side/nextjs)
- [GoTrue CORS Docs](https://github.com/supabase/gotrue#cors-configuration)

## ✅ Checklist Before Going Live

- [ ] nginx.conf updated with origin map
- [ ] docker-compose.prod.yml has correct GOTRUE_CORS_ALLOWED_ORIGINS
- [ ] All services restarted
- [ ] Test signup flow in browser
- [ ] Check browser console for logs
- [ ] Verify Network tab shows correct headers
- [ ] Test from production domain (not localhost)
- [ ] SSL certificates valid and trusted
- [ ] Backend can verify JWT tokens from auth service

---

**Status:** Ready for deployment and testing

If you still encounter issues after deploying these changes, the diagnostic logs will show exactly where the request is failing, which will help us narrow down the problem further.
