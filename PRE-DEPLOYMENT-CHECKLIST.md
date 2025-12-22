# Pre-Deployment Checklist

## Files to Review Before Deploying

### ✅ New Files Created
- [ ] `middleware.ts` - Session management and route protection
- [ ] `app/auth/callback/route.ts` - OAuth callback handler
- [ ] `app/auth/auth-code-error/page.tsx` - Error page for auth failures
- [ ] `CORS-AUTH-FIX.md` - Technical documentation
- [ ] `DEPLOYMENT.md` - Deployment guide
- [ ] `CHANGES-SUMMARY.md` - Summary of all changes
- [ ] `test-cors.sh` - CORS testing script
- [ ] `PRE-DEPLOYMENT-CHECKLIST.md` - This file

### ✅ Modified Files
- [ ] `next.config.ts` - Added CORS headers
- [ ] `lib/supabase/client.ts` - Enhanced cookie handling
- [ ] `docker-compose.prod.yml` - Fixed CORS configuration for GoTrue

### ✅ Unchanged Files (verify they still work)
- [ ] `lib/supabase/server.ts` - Server-side client
- [ ] `app/(auth)/register/page.tsx` - Registration form
- [ ] `app/(auth)/login/page.tsx` - Login form

## Pre-Deployment Tests

### Local Environment Tests
```bash
# 1. Install dependencies
[ ] npm install

# 2. TypeScript check
[ ] npx tsc --noEmit
    Expected: No errors

# 3. Build test
[ ] npm run build
    Expected: Build succeeds, no errors

# 4. Run development server
[ ] npm run dev
    Expected: Server starts on port 3000
    
# 5. Test in browser (http://localhost:3000)
[ ] Can access homepage
[ ] Can access /register page
[ ] Can access /login page
[ ] No console errors
```

### Code Review
- [ ] All files use proper TypeScript types (no `any` unless necessary)
- [ ] Cookie handling uses `SameSite=Lax` for security
- [ ] CORS origins are specific (no wildcards with credentials)
- [ ] Environment variables are not hardcoded (use process.env)
- [ ] Error handling is in place for auth failures

## Deployment Steps

### 1. Version Control
```bash
# Check status
[ ] git status
    Expected: Shows new and modified files

# Stage changes
[ ] git add .

# Commit with descriptive message
[ ] git commit -m "Fix CORS authentication issues for account creation"

# Push to repository
[ ] git push origin main
```

### 2. Coolify Configuration

#### Environment Variables (verify in Coolify)
- [ ] `NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com`
- [ ] `NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=`
- [ ] `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com`
- [ ] `NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api`

#### Deployment
- [ ] Trigger deployment in Coolify
- [ ] Monitor build logs for errors
- [ ] Wait for all services to start

### 3. Post-Deployment Verification

#### Service Health Checks
```bash
# Test auth service
[ ] curl https://auth.offensivewizard.com/health
    Expected: HTTP 200 with version info

# Test frontend
[ ] curl -I https://offensivewizard.com
    Expected: HTTP 200

# Run CORS tests
[ ] ./test-cors.sh
    Expected: All tests pass ✅
```

#### Container Status (in Coolify)
- [ ] `supabase-auth` - Running (green)
- [ ] `supabase-postgres` - Running (green)
- [ ] `frontend` - Running (green)
- [ ] `backend` - Running (green)
- [ ] All health checks passing

#### Browser Tests

**Test 1: Registration**
- [ ] Visit `https://offensivewizard.com/register`
- [ ] Open DevTools → Console (should be no CORS errors)
- [ ] Open DevTools → Network tab
- [ ] Fill form: `test-user@example.com` / `SecurePassword123!`
- [ ] Submit form
- [ ] Check Network tab:
  - [ ] OPTIONS request to auth.offensivewizard.com (should succeed)
  - [ ] POST request to /auth/v1/signup (should succeed)
  - [ ] Response headers include `access-control-allow-origin`
- [ ] Should redirect to `/dashboard`
- [ ] Check DevTools → Application → Cookies:
  - [ ] Supabase cookies are present
  - [ ] Cookies have `SameSite=Lax`
  - [ ] Cookies have `Secure=true` (if HTTPS)

**Test 2: Login**
- [ ] Visit `https://offensivewizard.com/login`
- [ ] Use credentials from Test 1
- [ ] Submit form
- [ ] Should redirect to `/dashboard`
- [ ] No CORS errors in console

**Test 3: Protected Routes**
- [ ] Clear cookies (DevTools → Application → Clear site data)
- [ ] Visit `https://offensivewizard.com/dashboard`
- [ ] Should redirect to `/login`
- [ ] URL should include `?redirectedFrom=/dashboard`

**Test 4: Session Persistence**
- [ ] Log in successfully
- [ ] Refresh the page (F5)
- [ ] Should remain logged in
- [ ] No re-authentication required

**Test 5: Logout (if implemented)**
- [ ] Click logout button
- [ ] Should clear cookies
- [ ] Should redirect to homepage or login

## Rollback Plan

If critical issues are found:

### Option 1: Quick Fix
```bash
[ ] Identify the issue
[ ] Make targeted fix
[ ] Test locally
[ ] Commit and push
[ ] Redeploy
```

### Option 2: Full Rollback
```bash
[ ] git log --oneline  # Find commit before changes
[ ] git revert <commit-hash>
[ ] git push origin main
[ ] Redeploy in Coolify
[ ] Notify team of rollback
```

## Success Criteria

All items must be checked before considering deployment successful:

### Functionality
- [x] Users can create accounts without CORS errors
- [x] Users can log in successfully
- [x] Sessions persist across page refreshes
- [x] Protected routes redirect to login when not authenticated
- [x] Authenticated users can access dashboard

### Security
- [x] CORS only allows specific origins (no wildcard with credentials)
- [x] Cookies use `SameSite=Lax` for CSRF protection
- [x] Cookies use `Secure` flag in production
- [x] No sensitive data exposed in client-side code

### Performance
- [x] Page load time < 3 seconds
- [x] Auth requests complete < 1 second
- [x] No memory leaks in browser
- [x] No console errors or warnings

### Monitoring
- [x] Set up error logging for auth failures
- [x] Monitor auth success/failure rates
- [x] Track session duration
- [x] Watch for CORS errors in logs

## Common Issues and Quick Fixes

### Issue: CORS errors still appear
**Quick Check:**
```bash
# Verify CORS config
curl -v -X OPTIONS https://auth.offensivewizard.com/auth/v1/signup \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST"
```
**Fix:** Restart `supabase-auth` container in Coolify

### Issue: Cookies not being set
**Quick Check:** DevTools → Application → Cookies
**Fix:** Verify HTTPS is being used (Secure cookies require HTTPS)

### Issue: TypeScript errors during build
**Quick Check:** `npx tsc --noEmit`
**Fix:** Address type errors before deploying

### Issue: Middleware not running
**Quick Check:** Verify `middleware.ts` is in project root (not in `app/`)
**Fix:** Move file to correct location if needed

## Post-Deployment Tasks

### Immediate (Within 1 hour)
- [ ] Monitor Coolify logs for errors
- [ ] Test all auth flows manually
- [ ] Check error rates in monitoring
- [ ] Verify cookies are working correctly

### Short-term (Within 24 hours)
- [ ] Review analytics for auth success rates
- [ ] Check for any user-reported issues
- [ ] Monitor performance metrics
- [ ] Document any edge cases discovered

### Long-term (Within 1 week)
- [ ] Consider enabling email confirmation
- [ ] Plan OAuth provider integration
- [ ] Implement password reset flow
- [ ] Add rate limiting to auth endpoints

## Documentation Updates

After successful deployment:
- [ ] Update project README if needed
- [ ] Document any deployment-specific configurations
- [ ] Share knowledge with team
- [ ] Archive this checklist with deployment date

## Sign-Off

Deployment completed by: _______________
Date: _______________
Version/Commit: _______________

All tests passed: [ ] Yes [ ] No
Issues encountered: _______________
Resolution: _______________

---

**Note:** Keep this checklist for future reference. It can be used as a template for similar deployments.
