# Coolify Environment Variable Update Required

## 🚨 Issue

The build is using the **old Supabase URL** because Coolify has environment variables that override the docker-compose.prod.yml values.

**Current (wrong):**
```
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
```

**Should be:**
```
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase-proxy
```

## ✅ Solution: Update Coolify Environment Variables

### Step-by-Step:

1. **Login to Coolify Dashboard**
   - Go to your Coolify instance

2. **Navigate to Frontend Service**
   - Projects → WizardCore → **frontend** service
   - Click on "Environment Variables" or "Configuration"

3. **Find and Update This Variable:**
   ```
   Variable: NEXT_PUBLIC_SUPABASE_URL
   Old Value: https://auth.offensivewizard.com
   New Value: https://offensivewizard.com/api/supabase-proxy
   ```

4. **Also Add (if not present):**
   ```
   Variable: SUPABASE_INTERNAL_URL
   Value: http://supabase-auth:9999
   ```

5. **Save Changes**

6. **Redeploy** the frontend service

## 🔍 Why This Is Needed

Coolify environment variables have **higher priority** than docker-compose.prod.yml values during build time. The Dockerfile uses:

```dockerfile
ARG NEXT_PUBLIC_SUPABASE_URL
```

Coolify passes this ARG from its environment variables, not from docker-compose.

## ✅ After Update

The build logs should show:
```
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/api/supabase-proxy
```

And your app will route Supabase requests through the Next.js proxy, bypassing CORS issues!

## 🐛 If Build Still Fails (Memory Issue)

If the build gets killed during "Creating an optimized production build...", it's running out of memory.

**Solutions:**

### Option A: Increase Build Memory in Coolify
1. Go to Server settings
2. Increase Docker build memory limit
3. Or add swap space to your server

### Option B: Optimize Build in Dockerfile
Add memory limits to the build command:

```dockerfile
# In Dockerfile.nextjs, modify the build command:
RUN NODE_OPTIONS="--max-old-space-size=2048" npm run build
```

This limits Node.js to 2GB RAM during build.

### Option C: Build Locally and Push Image
If server resources are limited:

```bash
# Build locally
docker buildx build --platform linux/amd64 -t yourregistry/wizardcore-frontend:latest -f Dockerfile.nextjs .

# Push to registry
docker push yourregistry/wizardcore-frontend:latest

# Update Coolify to use pre-built image
```

## 📋 Checklist

- [ ] Update `NEXT_PUBLIC_SUPABASE_URL` in Coolify UI
- [ ] Add `SUPABASE_INTERNAL_URL` in Coolify UI (if not present)
- [ ] Save changes
- [ ] Redeploy frontend service
- [ ] Check build logs show correct URL
- [ ] Test registration at https://offensivewizard.com/register
- [ ] Verify no CORS errors in browser console

---

**Priority:** HIGH - Required for CORS fix to work!
