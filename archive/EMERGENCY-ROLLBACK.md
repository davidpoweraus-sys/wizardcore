# Emergency Rollback Guide

## 🚨 Bad Gateway Error - Frontend Container Issue

### Quick Diagnosis in Coolify

1. **Check Container Status**
   - Go to Coolify Dashboard → frontend service
   - Check if container is running (should show "Running")
   - If stopped/crashed, check logs

2. **Check Logs**
   - Click on frontend service → "Logs"
   - Look for errors at the bottom
   - Common issues:
     - "Cannot find module" - build incomplete
     - "ECONNREFUSED" - can't connect to backend
     - Port already in use
     - Memory issues

3. **Check Previous Deployment**
   - Look for the last successful deployment
   - Note the commit hash or build number

## 🔄 Option 1: Quick Rollback (FASTEST)

### In Coolify:

1. **Go to Deployments History**
   - Frontend service → Deployments tab
   - Find last successful deployment (before the bad gateway)

2. **Click "Redeploy" on Previous Version**
   - This will revert to working state
   - Should restore site immediately

## 🔄 Option 2: Rollback via Git

### If you know the last working commit:

```bash
# Find recent commits
git log --oneline -10

# Rollback to previous commit (find the one before our changes)
# Let's say the last working commit was "028e034"
git revert HEAD --no-edit
git push

# Or hard reset (WARNING: loses recent commits)
git reset --hard 028e034
git push --force

# Then redeploy in Coolify
```

## 🔧 Option 3: Fix Current Deployment

### Likely Issues:

#### Issue 1: Proxy Route Not Found at Runtime

The proxy route might have build issues. **Quick fix:**

Check if this file exists in deployment:
- `app/api/supabase-proxy/[...path]/route.ts`

If it's missing or has errors, the app won't start.

#### Issue 2: Environment Variable Issue

The `SUPABASE_INTERNAL_URL` might be causing connection issues.

**Fix:** In Coolify environment variables, temporarily remove:
```
SUPABASE_INTERNAL_URL
```

Or set it to empty string if it can't be removed.

#### Issue 3: Build Artifacts Missing

The build might have completed partially but missing files.

**Fix:** Force a clean rebuild:
1. Coolify → frontend service
2. Clear build cache (if available)
3. Click "Force Rebuild"

## 🚨 IMMEDIATE ACTION: Revert Environment Variable

Since the issue started after env var changes:

1. **Go to Coolify → frontend → Environment Variables**
2. **Change back to:**
   ```
   NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
   ```
3. **Remove (if you added it):**
   ```
   SUPABASE_INTERNAL_URL
   ```
4. **Save and Redeploy**

This should restore the site (with CORS errors, but at least it loads).

## 📊 Diagnostic Commands

If you have SSH access to Coolify server:

```bash
# Check if frontend container is running
docker ps | grep frontend

# If not running, check why
docker ps -a | grep frontend

# Check container logs
docker logs <frontend-container-name> --tail 50

# Check if container can start
docker start <frontend-container-name>
docker logs -f <frontend-container-name>
```

## 🔍 What to Look For in Logs

### Good (server starting):
```
ready - started server on 0.0.0.0:3000
```

### Bad (errors):
```
Error: Cannot find module
Error: ECONNREFUSED
Error: listen EADDRINUSE
ModuleNotFoundError
SyntaxError
```

## ✅ Recovery Steps

1. **First:** Revert `NEXT_PUBLIC_SUPABASE_URL` back to `https://auth.offensivewizard.com`
2. **Second:** Remove any new env vars you added (`SUPABASE_INTERNAL_URL`)
3. **Third:** Redeploy
4. **Fourth:** Site should load (with CORS errors)
5. **Then:** We can troubleshoot the proxy approach differently

## 🎯 Alternative Approach (After Recovery)

Instead of changing the env var, we can:

1. Keep using `https://auth.offensivewizard.com`
2. Add Traefik/Caddy CORS headers via Coolify UI
3. Or use Cloudflare Workers as CORS proxy
4. Or enable CORS on Supabase Auth directly

The proxy approach is correct, but needs careful deployment testing.

---

**Priority:** Get site back online first, then fix CORS properly!
