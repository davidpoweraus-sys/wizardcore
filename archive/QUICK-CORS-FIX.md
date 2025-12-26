# Quick CORS Fix Reference

## 🚨 The Root Cause

Your nginx configuration had **conflicting CORS headers**. When using `if` blocks in nginx with `add_header`, it discards all parent headers, causing incomplete CORS responses.

## ✅ What Was Fixed

### File: `nginx-config/nginx.conf`
- Added origin whitelist using `map` directive
- Fixed header handling in OPTIONS requests
- Removed header duplication
- Added all required headers including `x-supabase-api-version`

### File: `app/(auth)/register/page.tsx`  
- Added detailed console logging for debugging
- Better error reporting

## 🚀 Deploy the Fix

```bash
cd /home/glbsi/Workbench/wizardcore

# Restart services to apply changes
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d

# Watch logs
docker logs -f supabase-auth
```

## 🧪 Test It

### Browser Test:
1. Open https://offensivewizard.com/register
2. Open DevTools (F12) → Console tab
3. Fill form and submit
4. Look for logs starting with:
   - 🚀 Registration started
   - 📤 Calling Supabase signUp...
   - ✅ Registration successful! OR ❌ Error details

### Network Test:
1. DevTools → Network tab
2. Filter: "signup"
3. Look for OPTIONS request
4. Check Response Headers:
   - `access-control-allow-origin: https://offensivewizard.com` ✅
   - `access-control-allow-credentials: true` ✅
   - `access-control-allow-methods: ...` ✅

### Command Line Test:
```bash
./test-cors-detailed.sh
```

## 🐛 Still Not Working?

### Check 1: Is nginx using new config?
```bash
docker exec $(docker ps -qf "name=nginx") cat /etc/nginx/nginx.conf | head -20
```
Should show: `map $http_origin $cors_origin`

### Check 2: Are services healthy?
```bash
docker ps
```
All should show "Up" status.

### Check 3: Check logs
```bash
docker logs supabase-auth --tail 50
docker logs $(docker ps -qf "name=nginx") --tail 50
```

### Check 4: Browser console
Look for the detailed logs I added. They'll tell you exactly where it's failing.

## 📞 What to Share If Still Broken

If it's still not working, send:
1. Screenshot of browser Console tab (with my logs)
2. Screenshot of Network tab (showing the OPTIONS request/response headers)
3. Output of: `docker logs supabase-auth --tail 30`
4. Output of: `./test-cors-detailed.sh`

This will show exactly what's happening.

## 🎯 Expected Behavior After Fix

1. Browser sends OPTIONS preflight → nginx returns 204 with CORS headers
2. Browser sends POST to signup → GoTrue processes → Returns user data
3. Cookies are set with proper domain/SameSite
4. User redirects to /dashboard
5. No CORS errors in console

---

**Files changed:**
- ✅ nginx-config/nginx.conf (CORS fix)
- ✅ app/(auth)/register/page.tsx (logging)
- ✅ test-cors-detailed.sh (diagnostics)
- ✅ CORS-FIX-SUMMARY.md (detailed docs)
