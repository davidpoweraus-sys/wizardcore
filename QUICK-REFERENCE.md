# Quick Reference - CORS Auth Fix

## 🚀 Quick Deploy

```bash
# 1. Test locally
npm run build

# 2. Commit and push
git add .
git commit -m "Fix CORS authentication issues"
git push origin main

# 3. Deploy in Coolify
# (Watch the deployment logs)

# 4. Test
./test-cors.sh
```

## 🔧 Quick Test

```bash
# CORS test
curl -v -X OPTIONS https://auth.offensivewizard.com/auth/v1/signup \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST"

# Expected: HTTP 204 with CORS headers
```

## 📋 What Changed

| File | Change | Why |
|------|--------|-----|
| `middleware.ts` | NEW | Session management, route protection |
| `next.config.ts` | Modified | CORS headers for frontend |
| `lib/supabase/client.ts` | Modified | Custom cookie handling |
| `docker-compose.prod.yml` | Modified | Fixed GoTrue CORS config |
| `app/auth/callback/route.ts` | NEW | OAuth callback handler |

## 🐛 Quick Troubleshooting

### CORS Errors
```bash
# Check auth service
curl https://auth.offensivewizard.com/health

# Verify CORS config in Coolify
# GOTRUE_CORS_ALLOWED_ORIGINS should be:
# "https://offensivewizard.com,..." (not "*")
```

### Cookies Not Set
```bash
# Check in browser DevTools → Application → Cookies
# Should have SameSite=Lax, Secure=true
```

### Random Logouts
```bash
# Verify middleware.ts exists in project root
ls -la middleware.ts
```

## 📍 Important URLs

- Frontend: `https://offensivewizard.com`
- Auth API: `https://auth.offensivewizard.com`
- Registration: `https://offensivewizard.com/register`
- Login: `https://offensivewizard.com/login`
- Dashboard: `https://offensivewizard.com/dashboard`

## 🔐 Environment Variables

```bash
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
```

## ✅ Success Checklist

- [ ] All containers running in Coolify
- [ ] `./test-cors.sh` passes
- [ ] Can create account at `/register`
- [ ] Can log in at `/login`
- [ ] No console errors
- [ ] Cookies are set
- [ ] Session persists on refresh

## 📚 Documentation

- **Full technical details**: `CORS-AUTH-FIX.md`
- **Deployment guide**: `DEPLOYMENT.md`
- **Changes summary**: `CHANGES-SUMMARY.md`
- **Pre-deployment checks**: `PRE-DEPLOYMENT-CHECKLIST.md`

## 🆘 Emergency Rollback

```bash
git log --oneline  # Find previous commit
git revert <commit-hash>
git push origin main
# Redeploy in Coolify
```

## 💡 Key Concepts

**CORS**: Cross-Origin Resource Sharing - allows frontend to talk to auth subdomain

**SameSite=Lax**: Cookie security setting that allows top-level navigation while preventing CSRF

**Middleware**: Runs on every request to manage sessions and protect routes

**Wildcard (*)**: Cannot be used with credentials per CORS spec

## 🧪 Test Commands

```bash
# TypeScript check
npx tsc --noEmit

# Build test
npm run build

# CORS test
./test-cors.sh

# Health check
curl https://auth.offensivewizard.com/health
```

## 📞 Support

If you encounter issues:
1. Check browser console for errors
2. Run `./test-cors.sh -v` for detailed output
3. Review `CORS-AUTH-FIX.md` troubleshooting section
4. Check Coolify logs for service errors
