# Systematic CORS Fix - Step by Step

## 🎯 Goal
Fix CORS without breaking the site. Test each step before deploying.

## 📊 Current Status
- ✅ Site is running
- ❌ CORS error when signing up
- ❌ Using direct URL: `https://auth.offensivewizard.com`
- ✅ Proxy route code is deployed in git

## 🔬 Step 1: Verify Proxy Route Works (Local Test First)

### Test locally before touching production:

```bash
cd /home/glbsi/Workbench/wizardcore

# Install dependencies if needed
npm install

# Run dev server
npm run dev
```

Then in another terminal:
```bash
# Test if proxy route responds
curl http://localhost:3000/api/supabase-proxy/health

# Expected: Should proxy to Supabase and return health response
# If 404: Route not working
# If 500: Route has error
```

**If this doesn't work locally, fix it before deploying to production.**

## 🔬 Step 2: Verify Proxy Route Exists in Production

**In your browser console on https://offensivewizard.com:**

```javascript
// Test 1: Check if proxy endpoint exists
fetch('/api/supabase-proxy/health')
  .then(r => {
    console.log('Status:', r.status)
    if (r.status === 404) {
      console.log('❌ Proxy route NOT deployed')
    } else if (r.status === 200) {
      console.log('✅ Proxy route exists!')
    }
    return r.text()
  })
  .then(data => console.log('Response:', data))
  .catch(e => console.error('Error:', e))
```

**Results:**
- **404** = Proxy route not deployed → Need to redeploy code
- **200** = Proxy route works! → Can proceed to Step 3
- **500** = Proxy route has error → Check logs, fix code
- **CORS error** = That's fine, we'll fix that next

## 🔬 Step 3: Test Proxy with Direct Fetch (Before Changing Supabase Client)

**In browser console on https://offensivewizard.com:**

```javascript
// Test proxy with actual Supabase endpoint
fetch('/api/supabase-proxy/auth/v1/health', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json'
  }
})
  .then(r => {
    console.log('✅ Status:', r.status)
    return r.json()
  })
  .then(data => console.log('✅ Response:', data))
  .catch(e => console.error('❌ Error:', e))
```

**Expected result:**
```json
{"version":"...","name":"GoTrue"}
```

**If this works** → Proxy is functional, can proceed to Step 4
**If this fails** → Check runtime logs for proxy errors

## 🔬 Step 4: Create Test Registration Page (Isolated Test)

Let's create a separate test page that uses the proxy, without breaking the main registration.

**Create new file:** `app/test-signup/page.tsx`

```typescript
'use client'

import { useState } from 'react'
import { createBrowserClient } from '@supabase/ssr'

export default function TestSignupPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [status, setStatus] = useState('')

  // Create Supabase client with PROXY URL
  const supabase = createBrowserClient(
    'https://offensivewizard.com/api/supabase-proxy', // Using proxy
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
  )

  const handleTest = async (e: React.FormEvent) => {
    e.preventDefault()
    setStatus('Testing...')
    console.log('🧪 Testing signup with proxy URL')

    try {
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
      })

      if (error) {
        console.error('❌ Error:', error)
        setStatus(`Error: ${error.message}`)
      } else {
        console.log('✅ Success:', data)
        setStatus('Success! Check console for details.')
      }
    } catch (err) {
      console.error('❌ Exception:', err)
      setStatus(`Exception: ${err}`)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md p-8 bg-white rounded shadow">
        <h1 className="text-2xl font-bold mb-4">🧪 Test Signup (Proxy)</h1>
        
        <form onSubmit={handleTest} className="space-y-4">
          <div>
            <label className="block text-sm font-medium mb-2">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 border rounded"
              placeholder="test@example.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium mb-2">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border rounded"
              placeholder="password123"
              required
            />
          </div>

          <button
            type="submit"
            className="w-full py-2 px-4 bg-blue-600 text-white rounded hover:bg-blue-700"
          >
            Test Signup
          </button>
        </form>

        {status && (
          <div className="mt-4 p-3 bg-gray-100 rounded">
            <strong>Status:</strong> {status}
          </div>
        )}

        <div className="mt-4 text-sm text-gray-600">
          <p><strong>This is a test page.</strong></p>
          <p>Using proxy URL: /api/supabase-proxy</p>
          <p>Open browser console (F12) to see detailed logs.</p>
        </div>
      </div>
    </div>
  )
}
```

**Deploy this test page:**
```bash
git add app/test-signup/page.tsx
git commit -m "Add test signup page for proxy testing"
git push
```

**After deploy, test at:** `https://offensivewizard.com/test-signup`

**Expected results:**
- ✅ No CORS errors in console
- ✅ Signup works
- ✅ Can see proxy logs in server logs

**If test page works** → Proxy is good, can update main registration
**If test page fails** → Debug proxy without breaking main site

## 🔬 Step 5: Check Server Logs for Proxy Activity

**In Coolify → frontend service → Runtime Logs**

Look for these log messages when you test signup:
```
🔄 Proxying request to: http://supabase-auth:9999/auth/v1/signup
✅ Proxy response status: 200
```

**If you see these** → Proxy is working!
**If you don't see these** → Proxy route not being called

## 🔬 Step 6: Update Main Registration (Only After Test Page Works)

Once the test page works perfectly:

**Update:** `app/(auth)/register/page.tsx`

Change the Supabase client creation to use a flag:

```typescript
const USE_PROXY = true // Toggle this for testing

const supabaseUrl = USE_PROXY 
  ? 'https://offensivewizard.com/api/supabase-proxy'
  : process.env.NEXT_PUBLIC_SUPABASE_URL!

const supabase = createClient(supabaseUrl, process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!)
```

This way you can easily switch back if something breaks.

## 🔬 Step 7: Gradual Rollout Plan

### Phase 1: Test Page Only (Current)
- Test page uses proxy
- Main registration uses direct URL
- No risk to production

### Phase 2: Enable for Main Registration
- Update main registration to use proxy
- Keep test page for comparison
- Can quickly revert if issues

### Phase 3: Update Environment Variable
- Change `NEXT_PUBLIC_SUPABASE_URL` globally
- All pages use proxy
- Remove test code

## 🐛 Debugging Checklist

### If proxy returns 404:
- [ ] Check file exists: `app/api/supabase-proxy/[...path]/route.ts`
- [ ] Check it's in git: `git ls-files | grep supabase-proxy`
- [ ] Check it's deployed: Test `/api/supabase-proxy/health`
- [ ] Clear build cache and rebuild

### If proxy returns 500:
- [ ] Check runtime logs for error message
- [ ] Common issues:
  - Can't connect to `http://supabase-auth:9999`
  - Missing environment variable
  - TypeScript/syntax error

### If proxy returns 502/503:
- [ ] Is supabase-auth service running?
- [ ] Are frontend and supabase-auth on same network?
- [ ] Check: `docker network inspect wizardcore-network`

### If still getting CORS errors:
- [ ] Check you're actually using proxy URL (check Network tab)
- [ ] Verify proxy is adding CORS headers (check Response Headers)
- [ ] Make sure it's not cached (hard refresh: Ctrl+Shift+R)

## 📊 Success Criteria

Before declaring victory:
- ✅ Test page works without CORS errors
- ✅ Can see proxy logs in server logs
- ✅ Registration completes successfully
- ✅ User is created in database
- ✅ No errors in browser console
- ✅ Main registration page works
- ✅ Site remains stable after changes

## 🚀 Timeline

**Day 1 (Today):**
- ✅ Site is running
- 🔄 Test proxy route exists in production
- 🔄 Create and deploy test page
- 🔄 Verify proxy works in isolation

**Day 2:**
- 🔄 Update main registration to use proxy
- 🔄 Test thoroughly
- 🔄 Monitor for issues

**Day 3:**
- 🔄 Update environment variable globally
- 🔄 Remove test code
- 🔄 Clean up documentation

## 🎯 Next Immediate Steps

1. **Test if proxy route exists** (Step 2 above)
2. **Create test page** (Step 4 above)
3. **Deploy test page only**
4. **Test at /test-signup**
5. **Report results**

---

**This approach:**
- ✅ Doesn't break production
- ✅ Tests in isolation
- ✅ Easy to rollback
- ✅ Clear debugging steps
- ✅ Gradual rollout

Let's start with Step 2 - check if proxy exists in production!
