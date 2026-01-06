# Local Build & Deploy System

## Problem Solved
Instead of building images locally and pushing to Docker Hub, we now build directly on the production server. This eliminates:
- Registry dependency issues
- Image push/pull complexity  
- Version mismatch problems
- Cache busting challenges

## New Architecture

### 1. **External Images** (Pulled from Docker Hub)
- `postgres:15-alpine` - WizardCore database
- `postgres:16-alpine` - Judge0 database  
- `redis:7-alpine` - Redis cache
- `supabase/gotrue:v2.184.0` - Authentication
- `judge0/judge0:latest` - Code execution

### 2. **Local Images** (Built on server)
- `wizardcore-backend:local` - Go backend with version logging
- `wizardcore-frontend:local` - Next.js frontend with null fixes

## Files Created

### `build-and-deploy.sh`
Complete build and deployment script with options:
```bash
# Full build and deploy
./build-and-deploy.sh full

# Pull external images only
./build-and-deploy.sh pull-only

# Build local images only  
./build-and-deploy.sh build-only

# Deploy only (assumes images built)
./build-and-deploy.sh deploy-only

# Show logs
./build-and-deploy.sh logs
```

### Updated `docker-compose.yml`
- Uses `wizardcore-backend:local` and `wizardcore-frontend:local`
- No `pull_policy: always` (not needed for local images)
- External images remain unchanged

### Updated Backend `Dockerfile`
- Accepts `BUILD_TIMESTAMP` as build argument
- Creates `.build-info` file with build metadata
- Cache busting with timestamp

### Updated `internal/version/version.go`
- Reads `BUILD_TIMESTAMP` from environment
- Dynamic version information at runtime
- Falls back to defaults if not set

## How It Works

### Build Process
1. **Pull external images** from Docker Hub
2. **Build backend** with timestamp: `wizardcore-backend:local`
3. **Build frontend** with timestamp: `wizardcore-frontend:local`  
4. **Update docker-compose.yml** to use local images
5. **Deploy** with `docker-compose up -d`

### Version Logging
Backend logs now include:
```json
{
  "version": "1.0.0",
  "build_time": "2026-01-06T06:38:30Z",  # From BUILD_TIMESTAMP
  "git_commit": "unknown",              # Can be set via GIT_COMMIT env
  "environment": "production"           # From ENVIRONMENT env
}
```

## Deployment Instructions

### On Production Server
```bash
# Make script executable
chmod +x build-and-deploy.sh

# Full build and deploy
./build-and-deploy.sh full

# Or step by step
./build-and-deploy.sh pull-only
./build-and-deploy.sh build-only
./build-and-deploy.sh deploy-only
```

### Verify Deployment
```bash
# Check services are running
docker-compose ps

# Check backend logs for version info
docker-compose logs backend | grep -E "(version|Starting server)"

# Expected output:
# Starting server version=1.0.0 build_time=2026-01-06T06:38:30Z
```

## Benefits

### 1. **Simplified Deployment**
- No Docker Hub account required
- No image tagging/pushing
- Direct build on target server

### 2. **Better Cache Control**
- Local Docker cache used efficiently
- Timestamp-based cache busting
- No registry cache issues

### 3. **Version Transparency**
- Build timestamp in logs
- Easy to identify which build is running
- No confusion between registry versions

### 4. **Faster Iteration**
- Build once, run immediately
- No push/pull latency
- Direct feedback on build errors

## GitHub Actions Integration (Optional)

If you want to trigger builds from GitHub:

```yaml
# .github/workflows/build-on-server.yml
name: Build on Server
on:
  push:
    branches: [main]

jobs:
  trigger-build:
    runs-on: ubuntu-latest
    steps:
      - name: Trigger server build via SSH
        uses: appleboy/ssh-action@v1.0.3
        with:
          host: ${{ secrets.SERVER_HOST }}
          username: ${{ secrets.SERVER_USER }}
          key: ${{ secrets.SERVER_SSH_KEY }}
          script: |
            cd /home/glbsi/Workbench/wizardcore
            git pull origin main
            ./build-and-deploy.sh full
```

## Troubleshooting

### 1. **Build fails with "Docker not running"**
```bash
# Start Docker daemon
sudo systemctl start docker
sudo systemctl enable docker
```

### 2. **Out of disk space**
```bash
# Clean up old images
docker system prune -a
```

### 3. **Port conflicts**
```bash
# Check what's using ports 3001, 8080, 9999, 2358
sudo lsof -i :3001
sudo lsof -i :8080

# Stop conflicting services
sudo systemctl stop conflicting-service
```

### 4. **Database migration issues**
```bash
# Reset databases (WARNING: deletes all data)
docker-compose down -v
docker-compose up -d
```

## Rollback Procedure

If something goes wrong:
```bash
# Stop services
docker-compose down

# Restore backup
cp docker-compose.yml.backup.$TIMESTAMP docker-compose.yml

# Start with previous configuration
docker-compose up -d
```

## Monitoring

### Check Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs backend -f
docker-compose logs frontend -f
```

### Check Health
```bash
# Service status
docker-compose ps

# Health checks
docker-compose exec backend wget -q -O- http://localhost:8080/health
docker-compose exec frontend wget -q -O- http://localhost:3000
```

## Next Steps

1. **Test the build script** on a staging server first
2. **Monitor logs** for version information
3. **Consider adding** Git commit hash to builds
4. **Set up monitoring** for build failures
5. **Create rollback** automation if needed

This system provides a robust, simple deployment process that builds directly on the target server, eliminating registry complexity while maintaining version tracking.