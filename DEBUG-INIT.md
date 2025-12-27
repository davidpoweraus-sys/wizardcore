# Debug Init Scripts Issue

## Check on Dokploy Server

```bash
# 1. Find the compose directory
ls -la /etc/dokploy/compose/

# 2. Check your project directory
ls -la /etc/dokploy/compose/offensivewizard-app-*/code/

# 3. Verify init-scripts exist
ls -la /etc/dokploy/compose/offensivewizard-app-*/code/init-scripts/

# 4. Check what's actually mounted in the container
docker exec offensivewizard-app-fomoox-supabase-init-1 ls -la /docker-entrypoint-initdb.d/

# 5. Check docker-compose context
docker inspect offensivewizard-app-fomoox-supabase-init-1 | grep -A 10 "Binds"
```

## Quick Fix

If the files aren't there, run this on the Dokploy server:

```bash
# Copy init scripts directly into the running container
docker cp /etc/dokploy/compose/offensivewizard-app-*/code/init-scripts/00-create-user.sql \
  offensivewizard-app-fomoox-supabase-init-1:/docker-entrypoint-initdb.d/

docker cp /etc/dokploy/compose/offensivewizard-app-*/code/init-scripts/01-create-auth-schema.sql \
  offensivewizard-app-fomoox-supabase-init-1:/docker-entrypoint-initdb.d/

# Restart the init container
docker restart offensivewizard-app-fomoox-supabase-init-1
```
