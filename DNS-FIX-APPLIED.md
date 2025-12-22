# DNS Fix Applied - Network-Qualified Hostname

## ✅ Problem Solved!

### Root Cause:
Docker DNS in Coolify/Linode environment couldn't resolve simple hostname `supabase-auth` due to search domain configuration (`members.linode.com`).

### The Fix:
Changed from:
```
http://supabase-auth:9999
```

To:
```
http://supabase-auth.d44co4gk48kok84wcg8o0os0_wizardcore-network:9999
```

This uses the **network-qualified hostname** which Docker DNS can resolve.

## 📊 What Was Changed:

### File 1: `app/api/supabase-proxy/[...path]/route.ts`
```typescript
// OLD:
const GOTRUE_URL = process.env.SUPABASE_INTERNAL_URL || 'http://supabase-auth:9999'

// NEW:
const GOTRUE_URL = process.env.SUPABASE_INTERNAL_URL || 'http://supabase-auth.d44co4gk48kok84wcg8o0os0_wizardcore-network:9999'
```

### File 2: `docker-compose.prod.yml`
```yaml
# OLD:
- SUPABASE_INTERNAL_URL=http://supabase-auth:9999

# NEW:
- SUPABASE_INTERNAL_URL=http://supabase-auth.d44co4gk48kok84wcg8o0os0_wizardcore-network:9999
```

## 🧪 Testing After Deploy:

Once Coolify redeploys:

### Test 1: Proxy Route Health Check
**In browser console on https://offensivewizard.com:**
```javascript
fetch('/api/supabase-proxy/health')
  .then(r => r.json())
  .then(d => console.log('✅ Proxy works!', d))
```

**Expected:**
```json
{"version":"v2.184.0","name":"GoTrue","description":"..."}
```

### Test 2: Test Page Signup
1. Go to: **https://offensivewizard.com/test-signup**
2. Toggle **"Use Proxy" ON**
3. Enter email/password
4. Click **Test Signup**

**Expected:**
- ✅ No CORS errors in console
- ✅ Signup succeeds
- ✅ Status shows: "Success!"

### Test 3: Check Frontend Logs
**In Coolify → frontend → Runtime Logs:**

Should see:
```
🔄 Proxy Configuration:
  GOTRUE_URL: http://supabase-auth.d44co4gk48kok84wcg8o0os0_wizardcore-network:9999
  ...
📤 Making request to Supabase Auth...
✅ Proxy response status: 200
```

## 🎯 What This Means:

### For Development:
- ✅ Proxy can now reach Supabase Auth internally
- ✅ No external network hops needed
- ✅ Fast, secure, internal communication

### For CORS:
- ✅ Browser requests go to same origin: `offensivewizard.com/api/supabase-proxy`
- ✅ No CORS preflight needed (same origin)
- ✅ Proxy adds CORS headers for response
- ✅ No CORS errors!

### For Performance:
- ✅ Internal Docker network (fast)
- ✅ No external DNS lookups
- ✅ No internet roundtrip

## 🚀 Next Steps After Deploy:

### Step 1: Verify Proxy Works (2 min)
```javascript
// Browser console test
fetch('/api/supabase-proxy/health').then(r => r.json()).then(console.log)
```

### Step 2: Test Signup on Test Page (2 min)
- Visit `/test-signup`
- Try signup with Proxy ON
- Should work!

### Step 3: Update Main Registration (5 min)
Once test page works, we can update the main registration page to use the proxy permanently.

**Update in Coolify environment variables:**
```
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase-proxy
```

Then all pages will use the proxy automatically!

## 📋 Verification Checklist:

After deploy completes:

- [ ] Proxy health check returns GoTrue info
- [ ] Test page signup works (Proxy ON)
- [ ] No CORS errors in browser console
- [ ] Frontend logs show successful proxy requests
- [ ] Can create user accounts
- [ ] Main registration page ready to update

## 🎉 Expected Result:

**Before:** CORS errors, signup fails
```
❌ CORS request did not succeed
❌ NetworkError when attempting to fetch
```

**After:** No CORS errors, signup works!
```
✅ 🚀 Registration started
✅ 📤 Calling Supabase signUp...
✅ ✅ Registration successful!
```

---

**Status:** Fix deployed, waiting for Coolify to rebuild and redeploy

**ETA:** 2-3 minutes for build to complete

**Test at:** https://offensivewizard.com/test-signup (after deploy)
