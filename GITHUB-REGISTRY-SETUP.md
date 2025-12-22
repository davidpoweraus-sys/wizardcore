# GitHub Container Registry Setup

## ✅ Why GitHub Container Registry?

- ✅ **100% Free** - No limits on private images
- ✅ **Unlimited** - No rate limiting like Docker Hub
- ✅ **Integrated** - Works seamlessly with GitHub repos
- ✅ **Fast** - Good CDN and download speeds

## 🚀 One-Time Setup

### Step 1: Create GitHub Token

1. Go to: https://github.com/settings/tokens/new
2. Token name: `Coolify Docker Registry`
3. Expiration: `No expiration` (or 1 year)
4. Select scopes:
   - ✅ `write:packages`
   - ✅ `read:packages`
   - ✅ `delete:packages` (optional)
5. Click **"Generate token"**
6. **Copy the token** (starts with `ghp_...`)

### Step 2: Login to GitHub Registry

On your **local machine** (where you'll build):

```bash
# Login to GitHub Container Registry
echo "YOUR_TOKEN_HERE" | docker login ghcr.io -u davidpoweraus-sys --password-stdin
```

You should see: `Login Succeeded`

### Step 3: Build and Push Frontend Image

```bash
cd /home/glbsi/Workbench/wizardcore

# Run the build script
./build-and-push.sh
```

This will:
1. Build the Next.js app (takes 5-10 min, needs 8GB RAM)
2. Create Docker image
3. Push to `ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest`

### Step 4: Configure Coolify to Pull Image

1. **In Coolify Dashboard:**
2. Go to **frontend service** → **Settings**
3. Add **Registry Credentials** (if private):
   - Registry: `ghcr.io`
   - Username: `davidpoweraus-sys`
   - Password: `YOUR_GITHUB_TOKEN`
4. **Save**

### Step 5: Make Image Public (Optional but Recommended)

To avoid needing credentials in Coolify:

1. Go to: https://github.com/users/davidpoweraus-sys/packages/container/wizardcore-frontend/settings
2. Scroll to **"Danger Zone"**
3. Click **"Change visibility"**
4. Select **"Public"**
5. Confirm

Now Coolify can pull without credentials!

## 🔄 Workflow: Making Changes

When you update your frontend code:

### 1. Make Code Changes Locally
```bash
# Edit files in your IDE
code app/some-page.tsx
```

### 2. Build and Push New Image
```bash
cd /home/glbsi/Workbench/wizardcore

# Build with latest code and push
./build-and-push.sh
```

### 3. Deploy to Coolify
```bash
# Commit docker-compose changes (if any)
git add .
git commit -m "Update frontend"
git push

# In Coolify: Click "Redeploy"
# Or wait for auto-deploy
```

Coolify will:
1. ✅ Pull latest image from GitHub (fast, no build!)
2. ✅ Start new container
3. ✅ No memory issues!

## 📊 How It Works

**Before (failing):**
```
Coolify Server (2GB RAM)
  → Download code
  → Build Next.js ❌ OUT OF MEMORY
  → Never finishes
```

**After (working):**
```
Your Computer (8GB+ RAM)
  → Build Next.js ✅
  → Push to ghcr.io ✅

Coolify Server
  → Pull pre-built image from ghcr.io ✅
  → Run container ✅
```

## 🛠️ Advanced: Automated Builds with GitHub Actions

Want to build automatically on every push?

Create `.github/workflows/build-frontend.yml`:

```yaml
name: Build and Push Frontend

on:
  push:
    branches: [main]
    paths:
      - 'app/**'
      - 'components/**'
      - 'lib/**'
      - 'public/**'
      - 'Dockerfile.nextjs'
      - 'package.json'

jobs:
  build:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write

    steps:
      - uses: actions/checkout@v3

      - name: Login to GitHub Container Registry
        uses: docker/login-action@v2
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          file: ./Dockerfile.nextjs
          push: true
          tags: ghcr.io/${{ github.repository_owner }}/wizardcore-frontend:latest
          build-args: |
            NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
            NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
            NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
            NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api
```

Then:
1. Push code to GitHub
2. GitHub Actions builds automatically
3. Image pushed to registry
4. Coolify pulls and deploys

Fully automated! 🎉

## ✅ Benefits

1. **No more memory issues** - Build on powerful machine
2. **Faster deployments** - Just pull image, don't build
3. **Consistent builds** - Same image everywhere
4. **Version control** - Tag images by version/commit
5. **Rollback easily** - Keep old images, switch tags

## 🔍 Troubleshooting

### "unauthorized: unauthenticated"
- Token expired or wrong
- Run `docker login ghcr.io` again

### "manifest unknown"
- Image doesn't exist yet
- Run `./build-and-push.sh` first

### "denied: permission_denied"
- Token needs `write:packages` scope
- Create new token with correct scopes

### Coolify can't pull image
- Make package public, or
- Add registry credentials in Coolify

---

**Ready to build?**

```bash
./build-and-push.sh
```

Then redeploy in Coolify!
