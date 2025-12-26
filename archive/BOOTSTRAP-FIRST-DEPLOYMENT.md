# Bootstrap First Deployment

## 🚨 One-Time Setup

The automated system needs a GitHub Release to exist before it can download images. Here's how to bootstrap the first deployment:

---

## Option A: Wait for GitHub Actions (Easiest)

GitHub Actions is building RIGHT NOW because we pushed `docker-compose.prod.yml`.

### Steps:
1. **Check GitHub Actions status:**
   ```
   https://github.com/davidpoweraus-sys/wizardcore/actions
   ```

2. **Wait for "Build and Release Docker Images as Tar" to complete** (~8-10 minutes)

3. **Verify release was created:**
   ```
   https://github.com/davidpoweraus-sys/wizardcore/releases
   ```
   Should see a release like `v2024.12.23-XXXX` with two tar files

4. **Redeploy in Coolify:**
   - Click "Redeploy" button
   - image-loader will download from the new release
   - Everything will work!

---

## Option B: Manual Bootstrap (If You Can't Wait)

Use the local tar file we already created:

### Steps:
1. **Transfer the tar to Coolify server:**
   ```bash
   scp /home/glbsi/Workbench/wizardcore/wizardcore-cors-fix.tar \
     user@coolify-server:/tmp/
   ```

2. **SSH into Coolify server:**
   ```bash
   ssh user@coolify-server
   ```

3. **Load the images:**
   ```bash
   docker load -i /tmp/wizardcore-cors-fix.tar
   ```

4. **Verify images loaded:**
   ```bash
   docker images | grep -E "(wizardcore-frontend|gotrue)"
   ```
   Should see:
   - `ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest`
   - `supabase/gotrue:v2.184.0`

5. **Redeploy in Coolify:**
   - Click "Redeploy" button
   - image-loader will see images already exist and skip download
   - Services will start using pre-loaded images

---

## Option C: Trigger GitHub Actions Manually (Alternative)

If GitHub Actions didn't trigger automatically:

### Steps:
1. **Go to Actions tab:**
   ```
   https://github.com/davidpoweraus-sys/wizardcore/actions
   ```

2. **Click "Build and Release Docker Images as Tar"**

3. **Click "Run workflow" button**

4. **Select branch: main**

5. **Click "Run workflow"**

6. **Wait 8-10 minutes**

7. **Redeploy in Coolify**

---

## ✅ After First Successful Deployment

Once the first deployment works (via any method above), **all future deployments are automatic**:

```bash
# Just push code
git push

# Wait 8-10 min for GitHub Actions

# Click Redeploy in Coolify

# Done!
```

The image-loader will automatically:
1. Check if images exist locally
2. If not, download latest release from GitHub
3. Load images
4. Exit
5. Other services start

---

## 🔍 Troubleshooting

### "Release not found"
**Cause:** GitHub Actions hasn't created a release yet

**Solution:** Wait for Actions to complete or use Option B (manual bootstrap)

### "Could not find wizardcore-cors-fix.tar.gz in release"
**Cause:** Release exists but tar files weren't uploaded

**Solution:** 
- Check Actions logs for errors
- Try re-running the workflow
- Or use Option B (manual bootstrap)

### "Images already loaded - skipping download"
**Not an error!** This means the bootstrap worked and images are cached.

### Check image-loader logs
```bash
docker logs wizardcore-image-loader
```

Should show download progress or "images already loaded" message.

---

## 📊 Current Status

**GitHub Actions Workflow:**
- ✅ Configured
- ✅ Watching for changes to app/, components/, docker-compose.prod.yml
- ⏳ Building now (after our push)

**Local Tar Files:**
- ✅ `wizardcore-cors-fix.tar` (254 MB) - Ready for manual bootstrap
- ✅ `wizardcore-complete-stack.tar` (11 GB) - For full deployments

**Docker Compose:**
- ✅ image-loader service configured
- ✅ Dependencies set correctly
- ✅ Will download from GitHub or use cached images

---

## 🎯 Recommended Path

**For first deployment: Option A (Wait for GitHub Actions)**

This creates the release properly and tests the full automated workflow.

**For emergency/testing: Option B (Manual bootstrap)**

Gets you running immediately while GitHub Actions runs in background.

---

## ✨ After Bootstrap

Once bootstrapped, your workflow is:

1. Make code changes
2. `git push`
3. Wait 8-10 min (GitHub builds)
4. Click "Redeploy" in Coolify
5. Done!

No manual tar transfers ever again! 🎉
