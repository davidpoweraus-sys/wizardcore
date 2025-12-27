# Environment Variable Update

## Port Configuration

Your frontend is now running on **port 3001** (not 3000) because Dokploy uses port 3000 for its dashboard.

## Access URLs

### Development/Testing (Direct Port Access):
- **Frontend**: `http://your-server-ip:3001`
- **Backend**: Internal only (access via domains or frontend)
- **Supabase Auth**: `http://your-server-ip:9999`
- **Judge0**: `http://your-server-ip:2358`
- **Dokploy Dashboard**: `http://your-server-ip:3000`

### Production (Via Traefik/Domains):
- **Frontend**: `https://app.offensivewizard.com`
- **Backend API**: `https://api.offensivewizard.com`
- **Supabase Auth**: `https://auth.offensivewizard.com`
- **Judge0**: `https://judge0.offensivewizard.com`

## Testing Your Deployment

```bash
# 1. Check if frontend is accessible
curl http://your-server-ip:3001

# 2. Check backend health
docker exec offensivewizard-app-fomoox-backend-1 wget -qO- http://localhost:8080/health

# 3. Check supabase-auth health
curl http://your-server-ip:9999/health

# 4. Check all containers are healthy
docker ps --format "table {{.Names}}\t{{.Status}}" | grep offensivewizard-app
```

## Current Environment Variables (Verified)

✅ `DATABASE_URL` - Includes `?sslmode=disable`
✅ `REDIS_URL` - Correct format
✅ `GOTRUE_DB_DATABASE_URL` - Correct password
✅ All passwords match between services

## No Changes Needed

Your `.env` file is correct! The only thing that changed is the **exposed port** in docker-compose.yml:

```yaml
# OLD (conflicted with Dokploy):
ports:
  - "3000:3000"

# NEW (works with Dokploy):
ports:
  - "3001:3000"
```

This means:
- **Inside the container**: Frontend still runs on port 3000
- **On the host**: It's accessible via port 3001
- **Via Traefik**: Still uses app.offensivewizard.com (no port needed)
