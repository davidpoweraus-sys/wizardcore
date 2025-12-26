# WizardCore Complete Stack Deployment Guide

## 📦 What's Included

A complete, self-contained deployment package with **ALL** Docker images needed to run WizardCore:

### Services Included (8 images total)

1. **Frontend** - Next.js application (`ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest`)
2. **Backend** - Go API server (`wizardcore-backend:latest`)
3. **Supabase Auth** - GoTrue authentication (`supabase/gotrue:v2.184.0`)
4. **PostgreSQL 15** - Supabase database (`postgres:15-alpine`)
5. **PostgreSQL 16** - Judge0 database (`postgres:16-alpine`)
6. **Redis 7** - Main cache (`redis:7-alpine`)
7. **Redis 7.2** - Judge0 cache (`redis:7.2-alpine`)
8. **Judge0** - Code execution engine (`judge0/judge0:latest`)

### Package Details

- **Uncompressed:** ~11 GB
- **Compressed:** ~3.4 GB (69% compression ratio)
- **Format:** `.tar.gz` (gzip compressed tar archive)
- **File:** `wizardcore-complete-stack.tar.gz`

---

## 🚀 Quick Start Deployment

### Prerequisites

- Docker installed on target server
- At least 15GB free disk space
- Docker Compose installed (optional, for easy orchestration)

### Step 1: Transfer Stack to Server

```bash
# From your local machine
scp wizardcore-complete-stack.tar.gz user@your-server:/tmp/

# Should take 5-15 minutes depending on your connection
```

### Step 2: Load Images on Server

```bash
# SSH into your server
ssh user@your-server

# Navigate to the file
cd /tmp

# Load all images (takes 5-10 minutes)
gunzip -c wizardcore-complete-stack.tar.gz | docker load

# Verify images loaded
docker images | grep -E "(wizardcore|supabase|postgres|redis|judge0)"
```

### Step 3: Clone Repository (if needed)

```bash
# Clone your repo
git clone https://github.com/davidpoweraus-sys/wizardcore.git
cd wizardcore

# Or pull latest changes if already cloned
git pull origin main
```

### Step 4: Deploy with Docker Compose

```bash
# Start all services
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f
```

### Step 5: Verify Deployment

```bash
# Check all containers are running
docker ps

# Test frontend
curl http://localhost:3000

# Test backend
curl http://localhost:8080/health

# Test auth
curl http://localhost:9999/health
```

---

## 🛠️ Using the Helper Scripts

Two scripts are included to make this easier:

### Export Script (Run Locally)

```bash
# Build and export complete stack
./export-complete-stack.sh

# This will:
# 1. Build backend image
# 2. Pull all required images
# 3. Export everything to tar
# 4. Compress to tar.gz
```

### Load Script (Run on Server)

```bash
# After transferring wizardcore-complete-stack.tar.gz to server
./load-stack.sh

# This will:
# 1. Verify file exists
# 2. Load all images into Docker
# 3. Show next steps
```

---

## 📋 Detailed Deployment Steps

### For Coolify

1. **Transfer the stack file:**
   ```bash
   scp wizardcore-complete-stack.tar.gz user@coolify-server:/tmp/
   ```

2. **SSH into Coolify server:**
   ```bash
   ssh user@coolify-server
   ```

3. **Load images:**
   ```bash
   cd /tmp
   gunzip -c wizardcore-complete-stack.tar.gz | docker load
   ```

4. **Update docker-compose.prod.yml** to disable pulling:
   ```yaml
   frontend:
     image: ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest
     pull_policy: never  # Use local image, don't pull from registry
   
   backend:
     build: ./wizardcore-backend
     pull_policy: never
   ```

5. **Deploy in Coolify:**
   - Push updated docker-compose to git
   - Click "Redeploy" in Coolify UI
   - All images will be used from local cache!

### For Regular VPS/Server

1. **Transfer and load** (same as above)

2. **Set up environment:**
   ```bash
   cd ~/wizardcore
   cp .env.example .env
   # Edit .env with your settings
   ```

3. **Deploy:**
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

4. **Set up reverse proxy** (nginx, Caddy, etc.)

---

## 🔄 Updating the Stack

### When You Make Code Changes

**Option 1: Rebuild and Export Entire Stack**

```bash
# On your local machine
./export-complete-stack.sh

# Transfer new file to server
scp wizardcore-complete-stack.tar.gz user@server:/tmp/

# On server, reload
gunzip -c /tmp/wizardcore-complete-stack.tar.gz | docker load

# Redeploy
docker-compose -f docker-compose.prod.yml up -d
```

**Option 2: Export Only Changed Services**

```bash
# If only frontend changed
docker save ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest | gzip > frontend-update.tar.gz

# Transfer and load
scp frontend-update.tar.gz user@server:/tmp/
ssh user@server "gunzip -c /tmp/frontend-update.tar.gz | docker load"

# Restart just frontend
docker-compose -f docker-compose.prod.yml up -d frontend
```

---

## 💾 Disk Space Requirements

### During Load Process
- **Compressed file:** ~3.4 GB
- **Uncompressed (loading):** ~11 GB
- **Docker images:** ~11 GB
- **Total needed:** ~25 GB (temporary, then ~15 GB after cleanup)

### After Cleanup
```bash
# Remove the tar.gz file after loading
rm /tmp/wizardcore-complete-stack.tar.gz

# Final disk usage: ~11-12 GB for images
```

### Running Containers
- **Images:** ~11 GB
- **Volumes (databases):** 1-5 GB (grows over time)
- **Logs:** 100 MB - 1 GB
- **Total recommended:** 20-30 GB free space

---

## 🔍 Troubleshooting

### "No space left on device"
```bash
# Check disk space
df -h

# Clean up old Docker resources
docker system prune -a

# Remove unused volumes
docker volume prune
```

### "Cannot load image"
```bash
# Verify file isn't corrupted
gunzip -t wizardcore-complete-stack.tar.gz

# If corrupted, re-transfer with rsync
rsync -avz --progress wizardcore-complete-stack.tar.gz user@server:/tmp/
```

### Images loaded but services won't start
```bash
# Check Docker logs
docker-compose -f docker-compose.prod.yml logs

# Verify all images are present
docker images

# Check for port conflicts
netstat -tulpn | grep -E "(3000|8080|9999|5432|6379|2358)"
```

### "pull access denied" or "unauthorized"
```bash
# This shouldn't happen with pre-loaded images
# If it does, add to docker-compose:
pull_policy: never  # For each service
```

---

## 📊 Image Manifest

Complete list of images in the stack:

| Image | Size | Purpose |
|-------|------|---------|
| `ghcr.io/davidpoweraus-sys/wizardcore-frontend:latest` | 212 MB | Next.js frontend |
| `wizardcore-backend:latest` | 47 MB | Go API backend |
| `supabase/gotrue:v2.184.0` | 50 MB | Authentication service |
| `postgres:15-alpine` | 274 MB | Supabase database |
| `postgres:16-alpine` | 276 MB | Judge0 database |
| `redis:7-alpine` | 41 MB | Main cache |
| `redis:7.2-alpine` | 41 MB | Judge0 cache |
| `judge0/judge0:latest` | 10.5 GB | Code execution |
| **Total** | **~11 GB** | |

---

## ✅ Verification Checklist

After deployment, verify:

- [ ] All 8 images loaded: `docker images | grep -E "(wizardcore|supabase|postgres|redis|judge0)" | wc -l` (should be 8+)
- [ ] All containers running: `docker-compose ps` (all should be "Up")
- [ ] Frontend accessible: `curl http://localhost:3000`
- [ ] Backend API working: `curl http://localhost:8080/health`
- [ ] Auth service healthy: `curl http://localhost:9999/health`
- [ ] Database connections work: Check logs for "connected to database"
- [ ] No error logs: `docker-compose logs --tail=50` (check for errors)

---

## 🎯 Benefits of This Approach

### ✅ Advantages
1. **No registry issues** - Everything is self-contained
2. **Faster deployment** - No pulling from external registries
3. **Offline capable** - Deploy without internet (after transfer)
4. **Version locked** - Exact same images every time
5. **No authentication** - No tokens or credentials needed
6. **Portable** - Transfer between servers easily

### ❌ Considerations
1. **Large file size** - 3.4 GB to transfer
2. **Manual updates** - Need to re-export and re-transfer for updates
3. **Disk space** - Requires ~25 GB during load, ~15 GB after

---

## 📝 Quick Command Reference

```bash
# Export (local)
./export-complete-stack.sh

# Transfer
scp wizardcore-complete-stack.tar.gz user@server:/tmp/

# Load (server)
gunzip -c /tmp/wizardcore-complete-stack.tar.gz | docker load

# Deploy (server)
cd ~/wizardcore
docker-compose -f docker-compose.prod.yml up -d

# Check status
docker-compose -f docker-compose.prod.yml ps

# View logs
docker-compose -f docker-compose.prod.yml logs -f

# Stop all
docker-compose -f docker-compose.prod.yml down

# Stop and remove volumes (⚠️ deletes data!)
docker-compose -f docker-compose.prod.yml down -v
```

---

## 🚀 Ready to Deploy!

Your complete WizardCore stack is packaged and ready. Just transfer, load, and run!

**Files you have:**
- `wizardcore-complete-stack.tar.gz` - Complete stack (3.4 GB)
- `export-complete-stack.sh` - Re-export script
- `load-stack.sh` - Server load script
- `docker-compose.prod.yml` - Deployment configuration

**For support or questions, check the logs:**
```bash
docker-compose -f docker-compose.prod.yml logs -f
```

Happy deploying! 🎉
