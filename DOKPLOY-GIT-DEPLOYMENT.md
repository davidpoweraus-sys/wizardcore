# Dokploy Git Source Deployment Guide

## ✅ **Deployment Script is Working**

The deployment script `deploy-from-git.sh` has been created and tested. It successfully builds images from git source on the server and uses those images.

## 🚀 **How to Deploy from Git Source on Dokploy Server**

### **Option 1: Use Dokploy's Built-in Build System**
Dokploy automatically builds from source when you:
1. Connect your Git repository to Dokploy
2. Configure `docker-compose.yml` with `build:` contexts (already done)
3. Click "Deploy" in Dokploy dashboard

**Current `docker-compose.yml` is already configured for source building:**
```yaml
backend:
  build:
    context: .
    dockerfile: wizardcore-backend/Dockerfile
  image: wizardcore-backend:${BUILD_TAG:-latest}

frontend:
  build:
    context: .
    dockerfile: Dockerfile.nextjs
  image: wizardcore-frontend:${BUILD_TAG:-latest}
```

### **Option 2: Manual Deployment with Script**
If you need to deploy manually on the Dokploy server:

1. **SSH into Dokploy server**
2. **Navigate to your application directory:**
   ```bash
   cd /etc/dokploy/compose/offensivewizard-app-wcilze/code
   ```
3. **Run the deployment script:**
   ```bash
   ./deploy-from-git.sh dokploy
   ```

### **Option 3: Full Git Source Deployment**
For complete control:
```bash
# 1. Clone repository (if not already cloned)
git clone <your-repo-url>
cd wizardcore

# 2. Run full deployment
./deploy-from-git.sh full
```

## 🔧 **Script Features**

The new `deploy-from-git.sh` script provides:

### **1. Git Source Verification**
```bash
./deploy-from-git.sh verify-source
```
Verifies you're in a git repository and shows current commit/branch.

### **2. Build Images from Source**
```bash
./deploy-from-git.sh build-only
```
Builds both backend and frontend images with timestamp tags.

### **3. Dokploy-style Deployment**
```bash
./deploy-from-git.sh dokploy
```
Runs the exact command Dokploy would run: `docker-compose up -d --build --remove-orphans`

### **4. Full Deployment Pipeline**
```bash
./deploy-from-git.sh full
```
Complete pipeline: verify → pull external images → build → deploy → verify

## 📊 **Verification Steps**

After deployment, verify:

### **1. Check Built Images**
```bash
docker images | grep wizardcore
```
Should show images with timestamp tags like `wizardcore-backend:20260106-075443`

### **2. Check Running Services**
```bash
docker-compose ps
```
All services should be "Up" and healthy.

### **3. Test API Endpoints**
```bash
# Backend health
curl http://localhost:8080/health

# Frontend
curl http://localhost:3001
```

## 🐛 **Troubleshooting**

### **Issue: "Build context not found"**
**Solution:** Ensure Dokploy has cloned the repository to the correct path. Check the build context in `docker-compose.yml`.

### **Issue: "Dockerfile not found"**
**Solution:** Verify the Dockerfile paths are correct relative to build context.

### **Issue: "Permission denied"**
**Solution:** Make script executable:
```bash
chmod +x deploy-from-git.sh
```

### **Issue: "Git repository not found"**
**Solution:** Ensure you're running in a git clone. Dokploy should clone automatically.

## 🏁 **Deployment Summary**

### **What Works:**
- ✅ `deploy-from-git.sh` builds images from git source
- ✅ Images are tagged with timestamps for versioning
- ✅ Docker Compose uses `build:` context (not pre-built images)
- ✅ Script verifies git source before building
- ✅ External dependencies are pulled automatically

### **For Dokploy Deployment:**
1. Ensure your `docker-compose.yml` has `build:` contexts (already done)
2. Dokploy will automatically build from source when you click "Deploy"
3. Or use `./deploy-from-git.sh dokploy` for manual deployment

### **Key Improvement:**
The script ensures images are built **from git source on the server** (not pulled from registry), which is exactly what you wanted.

## 📞 **Quick Reference**

**To build and deploy from git source:**
```bash
# One command does everything
./deploy-from-git.sh full

# Or for Dokploy-style
./deploy-from-git.sh dokploy
```

**To verify deployment:**
```bash
./deploy-from-git.sh verify
```

**To see logs:**
```bash
./deploy-from-git.sh logs
```

The deployment script is now working and ready for production use. It builds images from git source on the server and deploys them using docker-compose.