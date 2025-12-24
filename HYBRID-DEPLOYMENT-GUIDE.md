# Hybrid Deployment Guide - Best of Both Worlds!

## 🎯 Strategy

**Base Images** (Judge0, Postgres, Redis) → Pulled directly by docker-compose (one-time)  
**App Images** (Frontend, Backend, Auth) → Downloaded from GitHub releases (~90 MB)

This solves:
- ✅ GitHub Actions disk space issues
- ✅ Large file transfer problems
- ✅ Slow deployments

---

## 🚀 First Deployment (One-Time Setup)

### Step 1: Deploy in Coolify

Just click "Redeploy" in Coolify!

**What happens:**
1. `image-loader` service starts
2. Checks for base images (postgres, redis, judge0)
3. **Pulls them if missing** (takes 20-30 min, one-time only)
4. Tries to download app package from GitHub
5. If no release yet, continues with base images
6. Services start

**Expected on first run:**
- Base images pull automatically
- App images: "Release not found" (normal - GitHub Actions hasn't run yet)
- Services start (but frontend/backend will fail until app images loaded)

### Step 2: Wait for GitHub Actions

GitHub Actions is building the app package now:

1. Go to: https://github.com/davidpoweraus-sys/wizardcore/actions
2. Wait for "Build and Release Docker Images as Tar" (~5 min)
3. Check release created: https://github.com/davidpoweraus-sys/wizardcore/releases

### Step 3: Redeploy Again

Click "Redeploy" in Coolify again:

**What happens:**
1. Base images already exist (skip pulling)
2. Downloads app package from GitHub (~90 MB, fast!)
3. Loads frontend, backend, auth
4. All services start
5. **Working!** 🎉

---

## 🔄 Future Deployments

After first deployment, it's automatic:

```bash
# 1. Make code changes
git add .
git commit -m "Update feature"
git push

# 2. Wait 5 min for GitHub Actions

# 3. Click "Redeploy" in Coolify

# Done! Only downloads 90 MB app package
```

**Timeline:**
- GitHub Actions build: ~5 min (app only, not full stack)
- Download app package: ~30 sec (90 MB)
- Total: **Under 6 minutes from push to deployed!**

---

## 📦 What Gets Downloaded Each Time

| Item | Size | Frequency |
|------|------|-----------|
| **Base Images** | 10.9 GB | Once (first deployment) |
| **App Package** | 90 MB | Every deployment |

After first deployment: **37x less data** to transfer! (90 MB vs 3.4 GB)

---

## 🔧 Manual Override (Optional)

If you need to force re-pull base images:

```bash
# SSH into Coolify server
ssh root@your-coolify-server

# Delete specific base image
docker rmi judge0/judge0:latest

# Redeploy - will re-pull that image
```

Or set environment variable in Coolify:
```
PACKAGE_TYPE=complete-stack
```
This forces download of ALL images from GitHub (old behavior).

---

## 📋 Package Types

Set `PACKAGE_TYPE` environment variable in Coolify:

### `app` (Default - Recommended)
- Downloads: Frontend + Backend + Auth (~90 MB)
- Pulls base images if missing
- Fast updates
- Use for: Regular deployments

### `complete-stack` (Fallback)
- Downloads: All 8 services (~3.4 GB)
- Doesn't pull base images
- Slow but self-contained
- Use for: Troubleshooting, offline servers

---

## 🎨 Architecture

```
┌─────────────────────────────────────────────────────────────┐
│ First Deployment                                            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Click "Redeploy"                                          │
│         ↓                                                   │
│  image-loader starts                                        │
│         ↓                                                   │
│  ┌──────────────────────┐     ┌──────────────────────┐    │
│  │ Check Base Images    │     │ Download App Package │    │
│  │                      │     │ from GitHub          │    │
│  │ Missing?             │     │                      │    │
│  │   → Pull from Docker │     │ Release not found?   │    │
│  │     Hub (10.9 GB)    │     │   → Skip (ok)       │    │
│  │                      │     │                      │    │
│  │ Existing?            │     │ Found?               │    │
│  │   → Skip            │     │   → Load (90 MB)     │    │
│  └──────────────────────┘     └──────────────────────┘    │
│         ↓                              ↓                   │
│  All images ready                                          │
│         ↓                                                   │
│  Services start                                            │
│                                                             │
└─────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────┐
│ Subsequent Deployments (After GitHub Actions builds)       │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  Click "Redeploy"                                          │
│         ↓                                                   │
│  image-loader starts                                        │
│         ↓                                                   │
│  ┌──────────────────────┐     ┌──────────────────────┐    │
│  │ Check Base Images    │     │ Download App Package │    │
│  │                      │     │ from GitHub          │    │
│  │ ✓ All exist         │     │                      │    │
│  │ → Skip pulling       │     │ ✓ Found latest       │    │
│  │   (instant!)         │     │ → Load (90 MB)      │    │
│  │                      │     │   (30 seconds)       │    │
│  └──────────────────────┘     └──────────────────────┘    │
│         ↓                              ↓                   │
│  All images ready                                          │
│         ↓                                                   │
│  Services restart with new app images                      │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## ✅ Benefits of This Approach

### GitHub Actions
- ✅ Builds succeed (only 300 MB to process, not 11 GB)
- ✅ Fast builds (~5 min vs 10-15 min)
- ✅ No disk space issues
- ✅ Cheap bandwidth (90 MB releases)

### Coolify Server
- ✅ Base images pulled directly (their bandwidth)
- ✅ Fast updates (90 MB vs 3.4 GB)
- ✅ First deployment: pulls base images automatically
- ✅ Future deployments: only download app changes

### Developer Experience
- ✅ Push code → Auto build → Click redeploy → Done
- ✅ Under 6 minutes total
- ✅ No manual transfers
- ✅ No SSH needed
- ✅ Reliable (no timeouts, no "no space" errors)

---

## 🔍 Troubleshooting

### "Base image not found"
**Cause:** Image-loader couldn't pull base image

**Solution:**
```bash
# SSH into server
ssh root@your-coolify-server

# Manually pull the missing image
docker pull judge0/judge0:latest

# Redeploy
```

### "App images not found"
**Cause:** GitHub Actions hasn't created a release yet

**Solution:**
1. Check: https://github.com/davidpoweraus-sys/wizardcore/actions
2. Wait for workflow to complete
3. Redeploy in Coolify

### "Download failed"
**Cause:** Network issue or GitHub down

**Solution:**
Try again in a few minutes, or set:
```
RELEASE_TAG=v2024.12.23-XXXX  # Use specific older release
```

---

## 📊 Size Comparison

### Old Way (Complete Stack)
```
First deployment:  Transfer 3.4 GB → Load → Deploy
Every update:      Transfer 3.4 GB → Load → Deploy
```

### New Way (Hybrid)
```
First deployment:  Pull 10.9 GB (server) + Download 90 MB → Load → Deploy
Every update:      Download 90 MB → Load → Deploy
```

**Savings per update:** 97% less data transferred! (90 MB vs 3.4 GB)

---

## 🎯 Summary

**First deployment:**
- Click "Redeploy" → Pulls base images (20-30 min) → Done

**GitHub Actions finishes:**
- Creates app package release (~5 min after your push)

**Click "Redeploy" again:**
- Downloads app package (30 sec) → Working!

**Future deployments:**
- Push code → Wait 5 min → Click "Redeploy" → Done (30 sec)

**Total time from code push to deployed: Under 6 minutes!** 🚀
