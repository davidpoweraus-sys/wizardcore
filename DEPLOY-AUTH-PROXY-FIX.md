# Quick Deployment Guide - Auth Proxy Fix

## Option 1: Deploy via Docker Hub (Recommended)

### Step 1: Build and Push Images

On your **local machine**:

```bash
# Build and push to Docker Hub
./build-and-push-to-dockerhub.sh
```

This will:
- Build frontend with the new `/api/auth` proxy
- Build backend
- Push both to `limpet/wizardcore-frontend:latest` and `limpet/wizardcore-backend:latest`

### Step 2: Setup Environment Variables

On your **production server**, the `.env.example` file already has the correct production values! Just copy it:

```bash
# If you don't have a .env file yet
cp .env.example .env

# Or use the setup script
./setup-env.sh
```

**That's it!** The `.env.example` already contains:
- ✅ Correct proxy URL: `https://app.offensivewizard.com/api/auth`
- ✅ Valid ANON key (derived from the JWT secret)
- ✅ Correct GoTrue URL for the proxy
- ✅ All production-ready values

**Only customize if you:**
- Use different domain names
- Need to change passwords (recommended for production!)
- Changed the JWT secret (then regenerate keys with `node scripts/generate-anon-key.js`)

### Step 3: Pull and Restart

On your **production server**:

```bash
# Pull the latest images
docker pull limpet/wizardcore-frontend:latest
docker pull limpet/wizardcore-backend:latest

# Restart the services
docker-compose down
docker-compose up -d

# Or if using specific docker-compose file
docker-compose -f docker-compose.prod.yml down
docker-compose -f docker-compose.prod.yml up -d
```

### Step 4: Verify

```bash
# Check if proxy is working
curl https://app.offensivewizard.com/api/auth/health

# Should return:
# {"version":"v2.184.0","name":"GoTrue"...}
```

---

## Option 2: Deploy via SCP (Manual Transfer)

If you prefer to transfer files directly to the server:

### Files to Transfer

**Essential deployment files:**

```bash
# On your local machine
SERVER_IP="your.server.ip"

# Create directory on server
ssh root@$SERVER_IP "mkdir -p /opt/wizardcore"

# Transfer docker-compose file
scp docker-compose.prod.yml root@$SERVER_IP:/opt/wizardcore/docker-compose.yml

# Transfer environment configuration (update with your values first!)
scp .env.example root@$SERVER_IP:/opt/wizardcore/.env

# Transfer the updated Dockerfile (if building on server)
scp Dockerfile.nextjs root@$SERVER_IP:/opt/wizardcore/

# Transfer build script
scp build-and-push-to-dockerhub.sh root@$SERVER_IP:/opt/wizardcore/

# Make script executable
ssh root@$SERVER_IP "chmod +x /opt/wizardcore/*.sh"
```

### Then on the server:

```bash
ssh root@$SERVER_IP

cd /opt/wizardcore

# Edit .env file with your production values
nano .env

# Update these critical values:
# NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth
# NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
# GOTRUE_URL=https://auth.offensivewizard.com

# Pull and restart services
docker-compose pull
docker-compose up -d
```

---

## Option 3: Full Project Transfer (Build on Server)

If you want to build on the server:

### Transfer entire project:

```bash
SERVER_IP="your.server.ip"

# Sync the entire project (excluding node_modules, .next, etc.)
rsync -avz --progress \
  --exclude 'node_modules' \
  --exclude '.next' \
  --exclude 'dist' \
  --exclude '.git' \
  --exclude 'wizardcore-backend/vendor' \
  /home/glbsi/Workbench/wizardcore/ \
  root@$SERVER_IP:/opt/wizardcore/
```

### Then on the server:

```bash
ssh root@$SERVER_IP

cd /opt/wizardcore

# Update .env with production values
nano .env

# Build and push (or just build locally)
./build-and-push-to-dockerhub.sh

# Start services
docker-compose up -d
```

---

## Recommended Approach

**I recommend Option 1 (Docker Hub)** because:

1. ✅ Faster deployment (no file transfer)
2. ✅ Images are tested on your local machine first
3. ✅ Server just pulls and runs
4. ✅ Easy rollback (pull previous image)
5. ✅ Consistent builds across environments

## Quick One-Liner

If you already have docker-compose on the server:

```bash
# On local machine: Build and push
./build-and-push-to-dockerhub.sh

# On server: Update env, pull, and restart
ssh root@YOUR_SERVER "cd /opt/wizardcore && \
  echo 'NEXT_PUBLIC_SUPABASE_URL=https://app.offensivewizard.com/api/auth' >> .env && \
  echo 'NEXT_PUBLIC_SUPABASE_ANON_KEY=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJyb2xlIjoiYW5vbiIsImlzcyI6InN1cGFiYXNlIiwiaWF0IjoxNzY2NzQ5NzMzLCJleHAiOjIwODIxMDk3MzN9.R7vaBwwIssuKBRIBN0jx7xvzs7rYxjeD3zcZXhF60eQ' >> .env && \
  echo 'GOTRUE_URL=https://auth.offensivewizard.com' >> .env && \
  docker-compose pull && \
  docker-compose up -d"
```

## Troubleshooting

### Issue: Docker pull fails

**Solution:** Login to Docker Hub on server:
```bash
docker login -u limpet
```

### Issue: Permission denied on scripts

**Solution:**
```bash
chmod +x *.sh
```

### Issue: Port conflicts

**Solution:** Check what's running:
```bash
docker ps
netstat -tulpn | grep :3000
```
