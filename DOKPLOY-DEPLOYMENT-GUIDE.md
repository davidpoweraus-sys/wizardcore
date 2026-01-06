# Dokploy Deployment Guide for WizardCore

## ✅ **Fixes Applied for Pathways Null Issue**

The following fixes have been implemented and are ready for deployment:

### **1. Backend Repository Fixes**
- **Pathway Repository**: Fixed `FindAll()`, `FindAllWithEnrollment()`, and `FindEnrollmentsByUserID()` to return empty slices `[]` instead of `null`
- **Achievement Repository**: Fixed `FindAll()` and `GetUserAchievements()` to return empty slices `[]` instead of `null`
- **All repository functions** now initialize slices with `make([]Type, 0)` instead of `var slice []Type`

### **2. Docker Compose Fixes**
- **Updated `docker-compose.yml`**: Changed from `image: wizardcore-backend:local` to `build:` context with proper paths
- **Fixed build contexts**: 
  - Backend: `context: .` with `dockerfile: wizardcore-backend/Dockerfile`
  - Frontend: `context: .` with `dockerfile: Dockerfile.nextjs`
- **Added `.env.dokploy`**: Default environment variables for Dokploy deployment

### **3. Expected API Response Changes**
**Before Fix:**
```json
{"pathways":null}
```

**After Fix:**
```json
{"pathways":[]}
```

## 🚀 **Deployment Instructions for Dokploy**

### **Option 1: Automatic Deployment (Recommended)**
1. **Go to Dokploy dashboard** → Your app `offensivewizard-app-wcilze`
2. **Click "Deploy"** - Dokploy will:
   - Clone the repository
   - Run `docker compose -f ./docker-compose.yml up -d --build --remove-orphans`
   - Build images locally on Dokploy server
   - Start all services

### **Option 2: Manual Configuration**
If automatic deployment fails:

1. **Update Environment Variables in Dokploy:**
   - Go to your app settings in Dokploy
   - Add environment variables from `.env.dokploy`
   - **Important**: Update `GOTRUE_JWT_SECRET` and `SUPABASE_JWT_SECRET` with secure values

2. **Verify Build Contexts:**
   - Dokploy clones repo to: `/etc/dokploy/compose/offensivewizard-app-wcilze/code`
   - Build contexts are relative to this directory
   - Current configuration uses `context: .` (repository root)

3. **Monitor Deployment Logs:**
   - Check Dokploy deployment logs for errors
   - Look for "Building wizardcore-backend" and "Building wizardcore-frontend" messages

## 🔧 **Troubleshooting**

### **Issue: "pull access denied for wizardcore-backend"**
**Cause**: Docker Compose trying to pull non-existent images
**Solution**: The updated `docker-compose.yml` uses `build:` context instead of `image:` references

### **Issue: Build failures**
**Check**:
1. Dockerfile paths are correct
2. Build context has all necessary files
3. Environment variables are set

### **Issue: Services not starting**
**Check**:
1. Health checks in `docker-compose.yml`
2. Database connections
3. Network configuration

## 📊 **Verification Steps After Deployment**

### **1. Check Pathways API**
```bash
curl -H "Authorization: Bearer <token>" \
     https://app.offensivewizard.com/api/backend/v1/pathways
```
**Expected**: `{"pathways":[]}` (empty array, not null)

### **2. Check Frontend JavaScript Errors**
1. Open browser console (F12)
2. Navigate to `https://app.offensivewizard.com/dashboard/pathways`
3. **Verify**: No `can't access property "map", h is null` errors

### **3. Test Complete User Flow**
1. Register new user
2. Login
3. Navigate to dashboard
4. Check pathways page
5. Verify no JavaScript errors

## 🏁 **Deployment Summary**

### **Changes Ready for Deployment:**
- ✅ Backend null-safety fixes for pathways and achievements
- ✅ Updated docker-compose.yml with proper build contexts
- ✅ Environment variables configuration
- ✅ Deployment guide for Dokploy

### **Expected Outcome:**
- JavaScript error `can't access property "map", h is null` will be resolved
- Pathways API returns empty arrays `[]` instead of `null`
- All list endpoints return proper empty arrays
- Frontend components work without null checks

## ⚡ **Quick Deployment Command**

If Dokploy supports custom deployment commands, use:
```bash
# Set build tag
export BUILD_TAG=$(date +%Y%m%d-%H%M%S)

# Deploy with build
docker compose -f ./docker-compose.yml up -d --build --remove-orphans

# Check status
docker compose ps
```

## 📞 **Support**

If deployment fails:
1. Check Dokploy logs
2. Verify environment variables
3. Ensure Docker build contexts are correct
4. Contact development team if issues persist

**The pathways null fix is complete and ready for deployment. Once deployed, the JavaScript errors will be resolved and the application will function correctly.**