# GitHub Actions Automated Build Workflow

## ✅ What Was Set Up

A GitHub Actions workflow has been created at `.github/workflows/build-and-push-frontend.yml` that **automatically builds and pushes** your frontend Docker image to GitHub Container Registry.

## 🚀 How It Works

### Automatic Triggers
The workflow runs automatically when you:
1. **Push to main branch** AND
2. **Change any of these files:**
   - `app/**` (any page or component)
   - `components/**`
   - `lib/**`
   - `public/**`
   - `Dockerfile.nextjs`
   - `package.json` or `package-lock.json`
   - `next.config.ts` or `tsconfig.json`
   - The workflow file itself

### Manual Trigger
You can also trigger it manually:
1. Go to: https://github.com/davidpoweraus-sys/wizardcore/actions
2. Click **"Build and Push Frontend to GHCR"**
3. Click **"Run workflow"** → **"Run workflow"**

## 📋 What Happens During Build

1. ✅ **Checks out your code** from the repository
2. ✅ **Sets up Docker Buildx** (for efficient builds)
3. ✅ **Logs in to GHCR** using GitHub's built-in token (no manual token needed!)
4. ✅ **Builds the Docker image** with all your environment variables
5. ✅ **Pushes to GHCR** with two tags:
   - `latest` - Always points to the newest build
   - `main-<git-sha>` - Specific commit version (e.g., `main-48c026e`)
6. ✅ **Caches layers** to speed up future builds

## 🔄 New Workflow for Making Changes

### Before (Manual):
```bash
# Edit code
code app/page.tsx

# Build locally (takes 5-10 min, needs 8GB RAM)
./build-and-push.sh

# Redeploy in Coolify
```

### After (Automated):
```bash
# Edit code
code app/page.tsx

# Commit and push
git add .
git commit -m "Update homepage"
git push

# GitHub Actions builds automatically (4-6 min)
# Then redeploy in Coolify (or set up auto-deploy webhook!)
```

## 🎯 Benefits

1. ✅ **No local builds needed** - GitHub's servers handle it
2. ✅ **Faster** - GitHub Actions runners have better hardware
3. ✅ **Consistent** - Same build environment every time
4. ✅ **Version tracking** - Every commit gets its own tagged image
5. ✅ **Free** - 2,000 free minutes/month for private repos, unlimited for public
6. ✅ **No manual tokens** - Uses GitHub's built-in `GITHUB_TOKEN`

## 📊 Monitoring Builds

### Check Build Status
1. Go to: https://github.com/davidpoweraus-sys/wizardcore/actions
2. See all workflow runs with status (✅ success, ❌ failed, 🟡 in progress)
3. Click on any run to see detailed logs

### Build Duration
- **First build**: ~8-10 minutes (no cache)
- **Subsequent builds**: ~4-6 minutes (with cache)

## 🔧 Configuration Details

### Environment Variables
Built into the image at build time:
- `NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=`
- `NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com`
- `NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api`

### Image Tags
- `ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest`
- `ghcr.io/davidpoweraus-sys/wizardcore-frontend:main-<commit-sha>`

### Caching
The workflow uses Docker layer caching to speed up builds:
- Cache stored at: `ghcr.io/davidpoweraus-sys/wizardcore-frontend:buildcache`
- Significantly reduces build time for unchanged dependencies

## 🚀 Next Steps: Auto-Deploy to Coolify (Optional)

Want Coolify to automatically redeploy when GitHub Actions finishes? You can set up a webhook:

### Option 1: Coolify Webhook (Easiest)
1. In Coolify → Your frontend service → **Webhooks**
2. Copy the webhook URL
3. In GitHub → Settings → Webhooks → Add webhook
4. Paste URL, select **Package published** event
5. Now Coolify auto-deploys when image is pushed!

### Option 2: GitHub Actions Deploy Step
Add this as the final step in the workflow:

```yaml
- name: Trigger Coolify Deployment
  run: |
    curl -X POST "YOUR_COOLIFY_WEBHOOK_URL"
```

## 📝 Workflow File Location

`.github/workflows/build-and-push-frontend.yml`

## 🔍 Troubleshooting

### Workflow doesn't trigger
- Check if your file changes match the `paths:` filter
- Ensure you pushed to the `main` branch

### Build fails
- Check the Actions tab for detailed logs
- Common issues: syntax errors in code, missing dependencies

### Permission denied
- Workflow has `packages: write` permission by default
- If it fails, check repository Settings → Actions → General → Workflow permissions

### Need to update environment variables
Edit the workflow file and change the `build-args:` section

---

**Your workflow is now live!** 🎉

Every push to main will automatically build and push a new Docker image.

Check it out: https://github.com/davidpoweraus-sys/wizardcore/actions
