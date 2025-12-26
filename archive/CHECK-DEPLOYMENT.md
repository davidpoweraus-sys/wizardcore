# Deployment Checklist - Bad Gateway Diagnosis

## 🔍 Check #1: Was Environment Variable Updated in Coolify?

**In Coolify Dashboard:**

1. Go to: **Projects → WizardCore → frontend service**
2. Click: **Environment Variables** tab
3. Find: `NEXT_PUBLIC_SUPABASE_URL`
4. **Current value should be:** `https://offensivewizard.com/api/supabase-proxy`

**If it still shows:** `https://auth.offensivewizard.com`
- ❌ **The env var was NOT updated**
- This means the build used the proxy code but with the wrong URL
- **Fix:** Update it now and redeploy

## 🔍 Check #2: Build Logs

**In Coolify:**

1. Frontend service → **Build Logs** (latest deployment)
2. Search for: `NEXT_PUBLIC_SUPABASE_URL`
3. **Should show:**
   ```
   NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase-proxy
   ```

**If it shows the old URL:**
- The build didn't pick up the new env var
- Need to update env var in Coolify and rebuild

**If build failed:**
- Look for error messages in build logs
- Common issues:
  - TypeScript errors
  - Memory issues (killed/OOM)
  - Missing dependencies

## 🔍 Check #3: Runtime Logs

**In Coolify:**

1. Frontend service → **Runtime Logs** (container logs)
2. Look at the last 20 lines
3. **Errors to check for:**

### Error: "Cannot find module"
```
Error: Cannot find module '/app/server.js'
```
**Means:** Build didn't complete successfully, missing output files
**Fix:** Force rebuild with clean cache

### Error: "ECONNREFUSED" or "ETIMEDOUT"
```
Error: connect ECONNREFUSED http://supabase-auth:9999
```
**Means:** Proxy can't reach supabase-auth (but app should still start)
**Check:** Is supabase-auth service running?

### Error: "Port in use"
```
Error: listen EADDRINUSE: address already in use :::3000
```
**Means:** Another process using port 3000
**Fix:** Restart frontend container

### Error: Syntax/Import errors
```
SyntaxError: Unexpected token
Error: Cannot find module 'next/server'
```
**Means:** Code or dependency issue
**Fix:** Check if route.ts has correct syntax

## 🔍 Check #4: Container Status

**In Coolify:**

1. Frontend service → Overview
2. **Status should be:** `Running` (green)
3. **If status is:**
   - `Stopped` → Container crashed, check logs
   - `Restarting` → Container keeps crashing
   - `Exited` → Container stopped with error

**If container keeps restarting:**
- There's a fatal error in the code
- Check runtime logs for the crash reason

## 🔍 Check #5: Test Proxy Route Directly

**If container IS running but site shows 502:**

This means Traefik can't reach the container. Check:

1. **In Coolify:**
   - Is the domain correctly configured?
   - Is Traefik routing to the right port (3000)?

2. **Test internal health:**
   - If you have SSH access, try:
   ```bash
   docker exec <frontend-container> wget -qO- http://localhost:3000/api/supabase-proxy/health
   ```
   - If this works, Traefik routing is broken
   - If this fails, app has issues

## ✅ Most Likely Issue

Based on the symptoms (Bad Gateway after env var change):

**Scenario 1: Env Var NOT Updated**
- Build has proxy code
- Build thinks URL is still `https://auth.offensivewizard.com`
- But docker-compose says use proxy
- **Result:** Mismatch causes issues

**Fix:** 
1. Update `NEXT_PUBLIC_SUPABASE_URL` in Coolify to proxy URL
2. Force rebuild
3. Deploy

**Scenario 2: Build Incomplete**
- Build ran out of memory (we saw this before)
- Build didn't finish creating all files
- Container tries to start but `server.js` missing
- **Result:** Container crashes immediately

**Fix:**
1. Check build logs for "killed" or "OOM"
2. Clear build cache
3. Force rebuild

**Scenario 3: Proxy Route Runtime Error**
- Build succeeded
- Container starts
- Proxy route has runtime error when accessed
- **Result:** App crashes on first request

**Fix:**
1. Check runtime logs for error
2. May need to fix proxy code

## 🚀 Recommended Action Plan

**Step 1:** Check Coolify environment variables
- Is `NEXT_PUBLIC_SUPABASE_URL` set to proxy URL?
- If NO → Update it

**Step 2:** Check latest build logs
- Did build complete successfully?
- Does it show correct env var?
- If NO → Force rebuild

**Step 3:** Check runtime logs
- Is container running?
- Any crash errors?
- If CRASHED → Share error here

**Step 4:** If still stuck
- Revert `NEXT_PUBLIC_SUPABASE_URL` to old URL temporarily
- Get site back online
- Debug proxy separately

---

## 📊 Quick Commands Reference

**Check what files were deployed:**
```bash
git log --oneline -5
git show HEAD:app/api/supabase-proxy/[...path]/route.ts
```

**Check if proxy route is in latest commit:**
```bash
git ls-tree -r HEAD --name-only | grep supabase-proxy
```

Should show:
```
app/api/supabase-proxy/[...path]/route.ts
```
