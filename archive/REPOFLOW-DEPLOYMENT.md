# Repoflow Deployment Guide

## 📦 Docker Image Archive Created

Your Docker image has been saved and compressed for use with Repoflow or any container registry.

### Files Created

- **Uncompressed:** `wizardcore-frontend-image.tar` (206 MB)
- **Compressed:** `wizardcore-frontend-image.tar.gz` (68 MB) ⭐ **Use this one**

**Compression ratio:** 67% smaller (3x faster upload!)

---

## 🚀 Deployment Options

### **Option 1: Load on Coolify Server Directly**

If you can transfer the file to your Coolify server:

```bash
# 1. Copy the compressed tar to your server
scp wizardcore-frontend-image.tar.gz user@coolify-server:/tmp/

# 2. SSH into the server
ssh user@coolify-server

# 3. Load the image into Docker
gunzip -c /tmp/wizardcore-frontend-image.tar.gz | docker load

# 4. Verify it loaded
docker images | grep wizardcore-frontend

# 5. Redeploy in Coolify UI
# The image is now available locally, no registry pull needed!
```

---

### **Option 2: Upload to Repoflow**

If Repoflow is your container registry:

#### Step 1: Tag the Image for Repoflow

```bash
# Load the tar file first (if not already loaded)
docker load -i wizardcore-frontend-image.tar

# Tag it for your Repoflow registry
docker tag ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest \
  repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest
```

#### Step 2: Push to Repoflow

```bash
# Login to Repoflow
docker login repoflow.io
# Enter your Repoflow credentials

# Push the image
docker push repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest
```

#### Step 3: Update docker-compose.prod.yml

```yaml
frontend:
  image: repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest
  pull_policy: always
  # ... rest of config
```

#### Step 4: Commit and Deploy

```bash
git add docker-compose.prod.yml
git commit -m "Switch to Repoflow registry"
git push

# Redeploy in Coolify
```

---

### **Option 3: Private Registry (Self-Hosted)**

If you have your own private registry:

```bash
# Tag for your registry
docker tag ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest \
  registry.yourdomain.com/wizardcore-frontend:latest

# Push
docker push registry.yourdomain.com/wizardcore-frontend:latest
```

---

### **Option 4: Docker Hub**

```bash
# Login to Docker Hub
docker login

# Tag for Docker Hub
docker tag ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest \
  YOUR_DOCKERHUB_USERNAME/wizardcore-frontend:latest

# Push
docker push YOUR_DOCKERHUB_USERNAME/wizardcore-frontend:latest

# Update docker-compose.prod.yml
# image: YOUR_DOCKERHUB_USERNAME/wizardcore-frontend:latest
```

---

## 🔄 Automated Workflow for Repoflow

Want to automatically push to Repoflow via GitHub Actions?

### Create `.github/workflows/build-and-push-repoflow.yml`

```yaml
name: Build and Push to Repoflow

on:
  push:
    branches: [main]
    paths:
      - 'app/**'
      - 'components/**'
      - 'Dockerfile.nextjs'
      - 'package.json'

jobs:
  build-and-push:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3
      
      - name: Login to Repoflow
        uses: docker/login-action@v3
        with:
          registry: repoflow.io
          username: ${{ secrets.REPOFLOW_USERNAME }}
          password: ${{ secrets.REPOFLOW_TOKEN }}
      
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          file: ./Dockerfile.nextjs
          push: true
          tags: repoflow.io/YOUR_NAMESPACE/wizardcore-frontend:latest
          build-args: |
            NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
            NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPms=
            NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
            NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api
```

**Don't forget to add secrets to GitHub:**
- Go to: Settings → Secrets and variables → Actions
- Add: `REPOFLOW_USERNAME` and `REPOFLOW_TOKEN`

---

## 📋 Quick Reference: Loading Docker Images

### From Compressed Tar
```bash
gunzip -c wizardcore-frontend-image.tar.gz | docker load
```

### From Uncompressed Tar
```bash
docker load -i wizardcore-frontend-image.tar
```

### Verify Loaded
```bash
docker images | grep wizardcore-frontend
```

### Export Image Again (if needed)
```bash
docker save IMAGE_NAME:TAG | gzip > output.tar.gz
```

---

## 🔍 What's in the Tar File?

The tar archive contains:
- ✅ All Docker image layers
- ✅ Image metadata and configuration
- ✅ Built Next.js application
- ✅ Environment variables baked in
- ✅ Ready to run - just load and go!

---

## 💾 File Locations

```
/home/glbsi/Workbench/wizardcore/
├── wizardcore-frontend-image.tar       (206 MB - uncompressed)
└── wizardcore-frontend-image.tar.gz    (68 MB - compressed) ⭐
```

**Recommendation:** Use the `.tar.gz` file for faster transfers!

---

## 🎯 Recommended Approach for Coolify

**Fastest method:**

1. **Transfer compressed tar to server:**
   ```bash
   scp wizardcore-frontend-image.tar.gz coolify-server:/tmp/
   ```

2. **Load on server:**
   ```bash
   ssh coolify-server
   gunzip -c /tmp/wizardcore-frontend-image.tar.gz | docker load
   ```

3. **Update docker-compose.prod.yml to use local image:**
   ```yaml
   frontend:
     image: ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
     pull_policy: never  # Use local image, don't pull
   ```

4. **Redeploy in Coolify**

This avoids registry issues completely! 🚀

---

## 🛠️ Troubleshooting

### "No space left on device"
- The tar file is 206MB uncompressed
- Make sure you have at least 500MB free space

### "Cannot load image"
- Verify file isn't corrupted: `gunzip -t wizardcore-frontend-image.tar.gz`
- Check Docker is running: `docker ps`

### "Image not found after loading"
- Run: `docker images` to see all loaded images
- The image name is: `ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest`

---

## ✅ Next Steps

1. **Choose your deployment method** (direct load, Repoflow, Docker Hub, etc.)
2. **Transfer/upload the image**
3. **Update docker-compose.prod.yml** if changing registry
4. **Redeploy in Coolify**

The tar file is ready to use! 🎉
