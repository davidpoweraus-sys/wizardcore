# Quick Deploy Guide - Wizardcore on Dokploy

## ✅ Environment Configuration Complete

Your `.env` file has been updated with all required variables. All critical environment variables have been verified and are properly configured.

## 🚀 Deploy Now

### Option 1: Using Dokploy UI
1. Go to your Dokploy dashboard
2. Navigate to your application
3. Click "Deploy" or "Redeploy"
4. Dokploy will automatically use the updated `.env` file

### Option 2: Using Docker Compose Locally
```bash
cd /home/glbsi/Workbench/wizardcore
docker compose up -d --build
```

## 📊 Monitor Deployment

### Check Container Status
```bash
docker compose ps
```

Expected output - all containers should show "Up" and "healthy":
```
NAME                                          STATUS
offensivewizard-app-supabase-postgres-1      Up (healthy)
offensivewizard-app-supabase-auth-1          Up (healthy)  ← This was failing before
offensivewizard-app-postgres-1               Up (healthy)
offensivewizard-app-redis-1                  Up (healthy)
offensivewizard-app-backend-1                Up (healthy)
offensivewizard-app-backend-2                Up (healthy)
offensivewizard-app-frontend-1               Up (healthy)
offensivewizard-app-frontend-2               Up (healthy)
offensivewizard-app-judge0-postgres-1        Up (healthy)
offensivewizard-app-judge0-redis-1           Up (healthy)
offensivewizard-app-judge0-1                 Up
offensivewizard-app-judge0-worker-1          Up
offensivewizard-app-judge0-worker-2          Up
```

### Check Specific Container Logs
```bash
# Check supabase-auth (this was the failing container)
docker compose logs -f supabase-auth

# Check backend
docker compose logs -f backend

# Check frontend
docker compose logs -f frontend
```

## 🔍 Verify Services are Working

### 1. Check Supabase Auth Health
```bash
curl http://localhost:9999/health
```
Expected: `{"version":"...","name":"GoTrue"}`

### 2. Check Backend Health
```bash
curl http://localhost:8080/health
```
Expected: HTTP 200 OK

### 3. Check Frontend
```bash
curl http://localhost:3000
```
Expected: HTML response from Next.js

### 4. Check Judge0
```bash
curl http://localhost:2358/about
```
Expected: Judge0 version information

## 🐛 Troubleshooting

### Issue: supabase-auth still unhealthy

**Solution 1: Check logs for specific error**
```bash
docker compose logs supabase-auth | tail -50
```

**Solution 2: Restart the container**
```bash
docker compose restart supabase-auth
docker compose logs -f supabase-auth
```

**Solution 3: Recreate from scratch**
```bash
docker compose stop supabase-auth
docker compose rm -f supabase-auth
docker compose up -d supabase-auth
```

### Issue: Backend can't connect to database

**Verify database is running:**
```bash
docker compose exec postgres pg_isready -U wizardcore
```

**Test connection manually:**
```bash
docker compose exec postgres psql -U wizardcore -d wizardcore -c "SELECT version();"
```

**Check DATABASE_URL is correct:**
```bash
grep DATABASE_URL .env
# Should output: postgresql://wizardcore:YKiDQeXatsFIVvILstuZrxYBCPUdBlFC@postgres:5432/wizardcore
```

### Issue: Frontend can't authenticate

**Verify Supabase auth is accessible:**
```bash
curl http://localhost:9999/health
```

**Check JWT secrets match:**
```bash
grep -E "JWT_SECRET" .env
# GOTRUE_JWT_SECRET and SUPABASE_JWT_SECRET should be identical
```

**Verify ANON key is set:**
```bash
grep NEXT_PUBLIC_SUPABASE_ANON_KEY .env
```

### Issue: Judge0 not working

**Check Judge0 logs:**
```bash
docker compose logs judge0
```

**Verify API key is set:**
```bash
grep JUDGE0_API_KEY .env
```

**Test Judge0 API:**
```bash
curl -X POST http://localhost:2358/submissions \
  -H "Content-Type: application/json" \
  -d '{"source_code":"print(42)","language_id":71}'
```

## 📝 What Was Fixed

The following critical issues were resolved:

1. ✅ **GOTRUE_JWT_SECRET** - Added (required for Supabase Auth JWT signing)
2. ✅ **DATABASE_URL** - Added complete PostgreSQL connection string
3. ✅ **REDIS_URL** - Added complete Redis connection string
4. ✅ **JUDGE0_API_KEY** - Generated and added
5. ✅ **GOTRUE_DB_DATABASE_URL** - Added complete connection string for GoTrue
6. ✅ **Supabase API Keys** - Regenerated from JWT secret
7. ✅ **GoTrue Configuration** - Added all required GoTrue-specific variables

## 📂 Files Created/Modified

- ✅ `.env` - Updated with complete configuration
- ✅ `.env.backup-TIMESTAMP` - Backup of your original configuration
- ✅ `.env.new` - New configuration template
- ✅ `verify-env.sh` - Environment verification script
- ✅ `ENV-FIX-SUMMARY.md` - Detailed fix documentation
- ✅ `QUICK-DEPLOY-GUIDE.md` - This file

## 🔐 Security Reminders

**PUBLIC (safe to expose):**
- `NEXT_PUBLIC_SUPABASE_ANON_KEY`
- `NEXT_PUBLIC_SUPABASE_URL`
- `NEXT_PUBLIC_BACKEND_URL`
- `NEXT_PUBLIC_JUDGE0_API_URL`

**SECRET (never expose):**
- `SUPABASE_SERVICE_ROLE_KEY`
- `SUPABASE_JWT_SECRET` / `GOTRUE_JWT_SECRET`
- `POSTGRES_PASSWORD`
- `REDIS_PASSWORD`
- `JUDGE0_API_KEY`
- All database passwords

## ✨ Success Indicators

Your deployment is successful when:

1. ✅ All containers show "healthy" status
2. ✅ `supabase-auth` container starts without errors
3. ✅ Backend can connect to database and Redis
4. ✅ Frontend loads at https://app.offensivewizard.com
5. ✅ Users can register and login
6. ✅ Judge0 can execute code submissions
7. ✅ No warning messages about missing environment variables

## 🎯 Next Steps After Deployment

1. **Test User Registration:**
   - Go to https://app.offensivewizard.com/register
   - Create a test account
   - Verify you can login

2. **Test Backend API:**
   - Check that API endpoints are accessible
   - Verify authentication tokens work

3. **Test Judge0 Integration:**
   - Submit a code execution request
   - Verify results are returned correctly

4. **Monitor Logs:**
   - Keep an eye on container logs for any errors
   - Check for any performance issues

5. **Set Up Backups:**
   - Configure database backups
   - Back up `.env` file securely

---

**Date:** December 26, 2025  
**Status:** ✅ Ready to Deploy  
**Verification:** All environment variables validated  

🚀 **You're ready to deploy!** 🚀
