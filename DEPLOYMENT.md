# Deployment Guide - CORS Auth Fix

## Quick Deploy to Coolify

### Step 1: Verify Files
Ensure these new files are in your repository:
```bash
# In project root
ls -la middleware.ts                          # ✓ Should exist
ls -la app/auth/callback/route.ts            # ✓ Should exist  
ls -la app/auth/auth-code-error/page.tsx     # ✓ Should exist
ls -la CORS-AUTH-FIX.md                      # ✓ Should exist (documentation)
```

### Step 2: Review Changes
```bash
# Check modified files
git status

# Should show:
# - middleware.ts (new)
# - next.config.ts (modified)
# - lib/supabase/client.ts (modified)
# - docker-compose.prod.yml (modified)
# - app/auth/callback/route.ts (new)
# - app/auth/auth-code-error/page.tsx (new)
```

### Step 3: Test Locally (Optional but Recommended)

```bash
# Install dependencies
npm install

# Run development server
npm run dev

# In another terminal, test the build
npm run build

# Check for any TypeScript errors
npx tsc --noEmit
```

### Step 4: Commit and Push
```bash
git add .
git commit -m "Fix CORS authentication issues for account creation

- Add Next.js middleware for session management
- Update CORS configuration to use specific origins (not wildcard)
- Enhance Supabase client with proper cookie handling
- Add auth callback route for OAuth flow
- Add error page for auth failures
- Update docker-compose CORS settings for GoTrue"

git push origin main
```

### Step 5: Deploy in Coolify

1. **Go to Coolify Dashboard**
   - Navigate to your WizardCore application

2. **Check Environment Variables**
   Verify these are set correctly in Coolify:
   ```
   NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
   NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPms=
   ```

3. **Deploy**
   - Click "Deploy" or wait for auto-deploy
   - Watch the build logs for errors

4. **Monitor Deployment**
   - Check that all services start successfully
   - Look for these containers:
     - `supabase-auth` (GoTrue)
     - `frontend` (Next.js)
     - `supabase-postgres`

### Step 6: Verify Deployment

1. **Check Service Health**
   ```bash
   # Test auth service
   curl https://auth.offensivewizard.com/health
   
   # Should return: {"version":"...","name":"GoTrue"}
   ```

2. **Test CORS Configuration**
   ```bash
   # From your local machine or server
   ./test-cors.sh
   ```

3. **Manual Browser Test**
   - Go to `https://offensivewizard.com/register`
   - Open browser DevTools (F12)
   - Go to Network tab
   - Fill out registration form
   - Submit
   - Check for:
     - ✅ No CORS errors in console
     - ✅ POST request to auth.offensivewizard.com succeeds
     - ✅ Redirect to dashboard
     - ✅ Cookies are set

### Step 7: Test Authentication Flow

#### Test Signup:
1. Visit `https://offensivewizard.com/register`
2. Enter: `test@example.com` / `password123`
3. Submit form
4. Should redirect to `/dashboard` or show "Account created"

#### Test Login:
1. Visit `https://offensivewizard.com/login`
2. Enter credentials from above
3. Submit
4. Should redirect to `/dashboard`

#### Test Protected Routes:
1. Log out (if logged in)
2. Try to visit `https://offensivewizard.com/dashboard`
3. Should redirect to `/login`

#### Test Session Persistence:
1. Log in
2. Refresh the page
3. Should remain logged in
4. Close and reopen browser
5. Visit the site
6. Should still be logged in (if cookies persist)

## Troubleshooting Deployment

### Issue: Build Fails

**Check build logs for:**
```
Type error: ...
```

**Solution:** Run `npm run build` locally to catch TypeScript errors before deploying.

### Issue: Auth Service Not Starting

**Check container logs in Coolify:**
```bash
# In Coolify, view logs for supabase-auth container
```

**Look for:**
- Database connection errors
- Environment variable issues
- Port conflicts

**Solution:** 
- Verify `supabase-postgres` is healthy
- Check `GOTRUE_DB_DATABASE_URL` is correct
- Ensure port 9999 is not in use

### Issue: CORS Errors Persist

**In browser console, check error message:**
```
Access to fetch at 'https://auth.offensivewizard.com/...' from origin 'https://offensivewizard.com' 
has been blocked by CORS policy: Response to preflight request doesn't pass access control check.
```

**Solution:**
1. Verify in Coolify that `GOTRUE_CORS_ALLOWED_ORIGINS` includes `https://offensivewizard.com`
2. Check it's NOT set to `*` (wildcard with credentials is invalid)
3. Restart the `supabase-auth` container
4. Run `./test-cors.sh` to diagnose

### Issue: Cookies Not Being Set

**Check in DevTools → Application → Cookies:**
- No Supabase cookies appear

**Possible causes:**
1. Domain mismatch
2. SameSite policy blocking
3. Missing Secure flag in HTTPS

**Solution:**
- Check that both domains use HTTPS
- Verify cookie settings in `lib/supabase/client.ts`
- Check browser console for cookie warnings

### Issue: Random Logouts

**Symptoms:**
- User logs in successfully
- After refresh or navigation, user is logged out

**Cause:** Middleware not running or session not refreshing

**Solution:**
1. Verify `middleware.ts` exists in project root
2. Check middleware config matcher includes your routes
3. Ensure cookies are being set with correct MaxAge

## Rollback Plan

If deployment fails and you need to rollback:

```bash
# Revert to previous commit
git log --oneline  # Find previous commit hash
git revert <commit-hash>
git push origin main

# Or reset to previous state (CAUTION: loses changes)
git reset --hard <previous-commit-hash>
git push origin main --force
```

## Post-Deployment Checklist

- [ ] All containers running (check Coolify dashboard)
- [ ] Auth service health check passes
- [ ] CORS test script passes
- [ ] Can create new account
- [ ] Can log in with existing account
- [ ] Protected routes redirect correctly
- [ ] Session persists across page refreshes
- [ ] No console errors in browser
- [ ] Cookies are set correctly

## Performance Monitoring

After deployment, monitor:

1. **Auth Success Rate**
   - Watch for failed signup/login attempts
   - Check Supabase Auth logs for errors

2. **Response Times**
   - Auth requests should be < 500ms
   - Page loads should be < 2s

3. **Error Rate**
   - Monitor browser console for CORS errors
   - Check server logs for 500 errors

## Next Steps

Once authentication is working:

1. **Set up email confirmation** (optional)
   - Configure SMTP settings in GoTrue
   - Change `GOTRUE_MAILER_AUTOCONFIRM` to `false`

2. **Add OAuth providers** (optional)
   - Google, GitHub, etc.
   - Configure in docker-compose.prod.yml

3. **Implement password reset**
   - Add forgot password page
   - Configure email templates

4. **Add rate limiting**
   - Protect auth endpoints from brute force
   - Use Redis for rate limit tracking

## Support

If you encounter issues not covered here:

1. Check `CORS-AUTH-FIX.md` for detailed technical explanation
2. Review Supabase Auth logs in Coolify
3. Run `./test-cors.sh -v` for verbose debugging
4. Check browser Network tab for request/response details

## References

- [Supabase Auth Documentation](https://supabase.com/docs/guides/auth)
- [Next.js Middleware](https://nextjs.org/docs/app/building-your-application/routing/middleware)
- [CORS Specification](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)
