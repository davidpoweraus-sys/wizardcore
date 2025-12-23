# GitHub-Automated Deployment System

## 🎯 Overview

Fully automated deployment system that:
1. **Builds** Docker images on every push to main
2. **Packages** them as tar files
3. **Uploads** to GitHub Releases
4. **Downloads** and loads automatically when you deploy in Coolify

**No manual tar transfers needed!**

---

## 🚀 How It Works

```
┌─────────────┐
│  Push Code  │
│   to main   │
└──────┬──────┘
       │
       ▼
┌──────────────────┐
│ GitHub Actions   │ ← Builds images automatically
│ - Build Frontend │
│ - Build Backend  │
│ - Pull Base imgs │
│ - Create Tars    │
│ - Upload Release │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ GitHub Releases  │ ← Stores tar files
│ wizardcore-cors- │
│ fix.tar.gz (90MB)│
│                  │
│ wizardcore-      │
│ complete-stack.  │
│ tar.gz (3.4GB)   │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Coolify Deploys  │ ← Redeploy in UI
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ image-loader svc │ ← Runs automatically
│ - Downloads tar  │
│ - Loads images   │
│ - Exits          │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Services Start   │ ← Uses loaded images
│ - Frontend       │
│ - Supabase Auth  │
│ - Backend        │
│ - Databases      │
│ etc...           │
└──────────────────┘
```

---

## ⚙️ Setup (One-Time)

### 1. GitHub Actions Workflow

Already configured! See `.github/workflows/build-and-release-tar.yml`

**Triggers on:**
- Push to `main` branch
- Changes to frontend, backend, or Dockerfile
- Manual trigger from Actions tab

**What it does:**
- Builds all Docker images
- Creates CORS fix tar (frontend + auth)
- Creates complete stack tar (all 8 services)
- Compresses with gzip
- Uploads to GitHub Releases with timestamp tag

### 2. Docker Compose Configuration

Already configured! The `docker-compose.prod.yml` includes:

```yaml
services:
  image-loader:
    # Downloads and loads images from GitHub releases
    # Runs FIRST before all other services
    image: docker:24-cli
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - ./load-images-from-github.sh:/load-images.sh
    environment:
      - PACKAGE_TYPE=cors-fix  # or "complete-stack"
      - RELEASE_TAG=latest     # or specific version
    entrypoint: ["/bin/sh", "/load-images.sh"]
    restart: "no"

  frontend:
    image: ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
    pull_policy: never  # Use local image loaded by image-loader
    depends_on:
      image-loader:
        condition: service_completed_successfully
    # ... rest of config

  supabase-auth:
    image: supabase/gotrue:v2.184.0
    pull_policy: never  # Use local image loaded by image-loader
    depends_on:
      image-loader:
        condition: service_completed_successfully
    # ... rest of config
```

---

## 🎯 Deployment Workflow

### For Coolify Users (Recommended)

**Step 1: Code Changes**
```bash
# Make your changes
git add .
git commit -m "Update frontend"
git push origin main
```

**Step 2: Wait for GitHub Actions**
- Go to: https://github.com/davidpoweraus-sys/wizardcore/actions
- Wait for build to complete (~8-10 minutes)
- Check releases: https://github.com/davidpoweraus-sys/wizardcore/releases

**Step 3: Redeploy in Coolify**
- Go to Coolify dashboard
- Click "Redeploy" on your WizardCore project

**What happens automatically:**
1. ✅ Coolify pulls latest `docker-compose.prod.yml` from git
2. ✅ `image-loader` service starts first
3. ✅ Downloads latest tar from GitHub Releases
4. ✅ Loads images into Docker
5. ✅ Exits successfully
6. ✅ Frontend and supabase-auth start using loaded images
7. ✅ Other services start normally

**No SSH needed! No manual file transfers!**

---

## 🎛️ Configuration Options

### Environment Variables

Set these in Coolify or your `.env` file:

```bash
# Which package to download
PACKAGE_TYPE=cors-fix           # Options: cors-fix, complete-stack

# Which release version
RELEASE_TAG=latest              # Or: v2024.12.23-1430 (specific version)
```

**CORS Fix (Default)**
- Size: ~90 MB
- Services: Frontend + Supabase Auth
- Use for: Quick deployments, CORS fixes, frontend updates

**Complete Stack**
- Size: ~3.4 GB
- Services: All 8 (frontend, backend, auth, databases, redis, judge0)
- Use for: Full deployments, new servers, major updates

---

## 📋 Available Releases

Releases are automatically created with timestamp tags:

```
v2024.12.23-1430   ← Latest
v2024.12.23-1200
v2024.12.22-1800
```

Each release contains:
- `wizardcore-cors-fix.tar.gz` (~90 MB)
- `wizardcore-complete-stack.tar.gz` (~3.4 GB)

View all releases: https://github.com/davidpoweraus-sys/wizardcore/releases

---

## 🔧 Manual Deployment (Alternative)

If you want to deploy outside of Coolify:

### Option 1: One-Liner
```bash
curl -fsSL https://raw.githubusercontent.com/davidpoweraus-sys/wizardcore/main/coolify-deploy.sh | bash -s cors-fix
```

### Option 2: Full Script
```bash
# Download deployment script
curl -fsSL -o deploy.sh https://raw.githubusercontent.com/davidpoweraus-sys/wizardcore/main/deploy-from-release.sh
chmod +x deploy.sh

# Deploy CORS fix
./deploy.sh cors-fix latest

# Or deploy complete stack
./deploy.sh complete-stack latest

# Or specific version
./deploy.sh cors-fix v2024.12.23-1430
```

### Option 3: Docker Compose
```bash
# Clone repo
git clone https://github.com/davidpoweraus-sys/wizardcore.git
cd wizardcore

# Set environment
export PACKAGE_TYPE=cors-fix
export RELEASE_TAG=latest

# Deploy
docker-compose -f docker-compose.prod.yml up -d
```

The `image-loader` service will automatically download and load images before other services start.

---

## 🎨 Customization

### Use Different Release

```bash
# In Coolify environment variables or .env
RELEASE_TAG=v2024.12.23-1200  # Use specific version instead of latest
```

### Skip Image Loading (Use Pre-Loaded)

```bash
# In docker-compose.prod.yml, comment out the image-loader service
# Or remove it from depends_on for frontend/supabase-auth
```

### Force Re-Download

```bash
# Delete local images to force fresh download
docker rmi ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
docker rmi supabase/gotrue:v2.184.0

# Then redeploy
```

---

## 🐛 Troubleshooting

### "Release not found"

**Cause:** No GitHub releases exist yet

**Solution:** 
1. Go to Actions tab
2. Run "Build and Release Docker Images as Tar" manually
3. Wait for completion
4. Check Releases tab

### "Download failed"

**Cause:** Network issue or GitHub API rate limit

**Solution:**
```bash
# Check GitHub status
curl -I https://api.github.com/rate_limit

# Try again in a few minutes
# Or use specific release tag instead of "latest"
```

### "Failed to load images"

**Cause:** Corrupted tar file or disk space issue

**Solution:**
```bash
# Check disk space
df -h

# Clear Docker cache
docker system prune -a

# Try redeploying
```

### "Images already loaded - skipping download"

**Not an error!** This means images are cached and ready to use.

To force re-download:
```bash
docker rmi ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
docker rmi supabase/gotrue:v2.184.0
```

### image-loader service fails

**Check logs:**
```bash
docker-compose -f docker-compose.prod.yml logs image-loader
```

**Common issues:**
- No internet connection
- GitHub rate limit
- Disk full
- Docker socket permission denied

---

## ✅ Benefits

| Feature | Before | After |
|---------|--------|-------|
| **Build Location** | Coolify server (limited RAM) | GitHub Actions (unlimited) |
| **Transfer** | Manual SCP | Automatic download |
| **Registry Auth** | Required (fails often) | Not needed |
| **Deployment Speed** | Slow (builds + fails) | Fast (just loads) |
| **Disk Usage** | Build artifacts | Just final images |
| **Memory Required** | 8GB+ for build | Minimal |
| **Internet Required** | For pulling base images | For downloading tar |
| **Version Control** | Manual tagging | Automatic timestamp |
| **Rollback** | Hard | Easy (use old release) |

---

## 📊 Workflow Summary

### Development Cycle

```bash
# 1. Make changes locally
vim app/page.tsx

# 2. Test locally
npm run dev

# 3. Commit and push
git add .
git commit -m "Update homepage"
git push origin main

# 4. Wait 8-10 minutes for GitHub Actions

# 5. Click "Redeploy" in Coolify

# 6. Done! CORS fix automatically applied
```

**Total time:** ~10 minutes from push to deployed
**Manual steps:** 2 (commit + click redeploy)
**SSH required:** None
**File transfers:** None

---

## 🎉 That's It!

Your deployment is now fully automated:

1. ✅ GitHub Actions builds images
2. ✅ Releases uploaded automatically
3. ✅ Coolify downloads and loads automatically
4. ✅ No manual tar transfers
5. ✅ No registry authentication issues
6. ✅ Version-controlled releases
7. ✅ Easy rollbacks

**Just push to main and click redeploy!** 🚀

---

## 📚 Related Files

- `.github/workflows/build-and-release-tar.yml` - GitHub Actions workflow
- `load-images-from-github.sh` - Image loader script
- `deploy-from-release.sh` - Manual deployment script
- `coolify-deploy.sh` - One-liner deployment
- `docker-compose.prod.yml` - Production configuration

---

**Questions?** Check the logs:
```bash
docker-compose -f docker-compose.prod.yml logs image-loader
```
