# Coolify Deployment Guide

## Option 1: Direct Git Deployment (RECOMMENDED)

### Setup in Coolify UI:

1. **Go to your Coolify dashboard** → Applications
2. **Edit your application** → Settings
3. **Configure Git Source:**
   - Repository: `https://github.com/davidpoweraus-sys/wizardcore`
   - Branch: `main`
   - Build Pack: `Dockerfile` or `Docker Compose`

4. **Set Build Configuration:**
   - **For Frontend:**
     - Dockerfile: `Dockerfile.nextjs`
     - Context: `/`
     - Build Args:
       ```
       NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/supabase-proxy
       NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
       NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
       NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api
       ```
   
   - **For Backend:**
     - Dockerfile: `wizardcore-backend/Dockerfile`
     - Context: `/wizardcore-backend`

5. **Enable Auto Deploy:**
   - ✅ Deploy on git push
   - ✅ Watch for changes on `main` branch

6. **Save and Deploy**

### Benefits:
- ✅ Automatic deployment on `git push`
- ✅ Builds from source (always latest code)
- ✅ No manual image transfer needed
- ✅ Coolify manages everything

---

## Option 2: Pre-built Images from GitHub Releases

If you want to use pre-built images from GitHub releases:

### 1. Update docker-compose.prod.yml

Change `pull_policy: never` to use GitHub releases:

```yaml
services:
  frontend:
    image: ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
    pull_policy: always  # Pull from registry
    # ... rest of config

  backend:
    image: ghcr.io/davidpoweraus-sys/wizardcore-backend:latest
    pull_policy: always  # Pull from registry
    # ... rest of config
```

### 2. Push Images to GitHub Container Registry (GHCR)

Update GitHub Actions to push to GHCR instead of creating tars:

```yaml
- name: Login to GHCR
  uses: docker/login-action@v3
  with:
    registry: ghcr.io
    username: ${{ github.actor }}
    password: ${{ secrets.GITHUB_TOKEN }}

- name: Build and Push Frontend
  uses: docker/build-push-action@v5
  with:
    context: .
    file: ./Dockerfile.nextjs
    push: true
    tags: ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
```

### 3. Configure Coolify to Pull from GHCR

In Coolify:
- Set image source to GHCR
- Enable auto-pull on deploy
- Coolify will pull latest images on each deployment

---

## Option 3: Hybrid - Load from GitHub Release (Current Method)

Keep using the manual deployment script:

```bash
# On your local machine
cd /home/glbsi/Workbench/wizardcore
./deploy-from-github.sh
```

This downloads the tar from GitHub releases and loads on the server.

---

## Recommended Approach

**Use Option 1** (Direct Git Deployment) because:
- ✅ Fully automated
- ✅ No manual steps
- ✅ Coolify's native workflow
- ✅ Builds are tracked and logged
- ✅ Easy rollbacks

**Avoid Option 2** for now because:
- ❌ GHCR authentication is complex with Coolify
- ❌ We had issues with this earlier

**Option 3** is fine for now but:
- ⚠️ Manual deployment required
- ⚠️ Not fully automated
