# Dokploy Deployment Fix Guide

## Problem
The `supabase-auth` container was failing with password authentication errors:
```
FATAL: password authentication failed for user "supabase_auth_admin"
```

## Root Cause
The Supabase Postgres volume on Dokploy contains old data with mismatched credentials. When PostgreSQL initializes from an existing volume, the `/docker-entrypoint-initdb.d` scripts don't run again, leaving the old credentials in place.

## Solution Applied

### 1. Created User Initialization Script
**File**: `init-scripts/00-create-user.sql`
- Creates or updates the `supabase_auth_admin` user with the correct password
- Grants necessary privileges (SUPERUSER, CREATEDB, CREATEROLE)
- This script is idempotent and safe to run multiple times

### 2. Updated Schema Initialization Script
**File**: `init-scripts/01-create-auth-schema.sql`
- Fixed enum type creation (removed ALTER TYPE inside transaction)
- Creates all necessary enum types for Supabase Auth
- Grants proper privileges to the auth schema

### 3. Enhanced Docker Compose Configuration
**File**: `docker-compose.yml`
- Improved `supabase-auth` healthcheck timings:
  - `start_period`: 30s → 60s
  - `retries`: 5 → 10
  - `interval`: 10s → 15s
  - `timeout`: 5s → 10s
- Updated `supabase-init` to run both initialization scripts in order
- Added `GOTRUE_LOG_LEVEL` for better debugging

### 4. Created Helper Scripts
- **`fix-dokploy-volumes.sh`**: Cleans up and resets the Supabase Postgres volume
- **`diagnose-deployment.sh`**: Diagnoses deployment issues

## How It Works Now

### Correct Startup Sequence

1. **`supabase-postgres`** starts
   - PostgreSQL initializes (creates user from `POSTGRES_USER` env var)
   - Healthcheck waits for database to be ready
   - Status: ✅ HEALTHY

2. **`supabase-init`** starts (waits for postgres healthy)
   - Runs `00-create-user.sql` - creates/updates user with correct password
   - Runs `01-create-auth-schema.sql` - creates auth schema and types
   - Writes "READY" status file
   - **Healthcheck succeeds only after scripts complete**
   - Status: ✅ HEALTHY

3. **`supabase-auth`** starts (waits for postgres AND init healthy)
   - Only starts **AFTER** initialization scripts are complete
   - Connects to database with correct credentials
   - Runs GoTrue migrations successfully
   - Status: ✅ HEALTHY

The key fix: `supabase-auth` now depends on `supabase-init: service_healthy` instead of `service_started`, ensuring scripts finish before auth tries to connect.

## Deployment Steps

### Option A: Clean Deployment (Recommended)

If you can afford to lose the Supabase Auth database data:

1. **SSH into your Dokploy server**

2. **Run the volume cleanup script**:
   ```bash
   cd /path/to/wizardcore
   ./fix-dokploy-volumes.sh
   ```

3. **Redeploy in Dokploy**:
   - Go to your Dokploy dashboard
   - Navigate to your application
   - Click "Redeploy" or trigger a new deployment

4. **Monitor the deployment**:
   - Watch the logs for `supabase-postgres` - should show database initialization
   - Watch the logs for `supabase-init` - should show script execution
   - Watch the logs for `supabase-auth` - should successfully connect and run migrations

### Option B: Keep Existing Data

If you need to keep existing auth data:

1. **Commit and push your code changes** to trigger a Dokploy redeploy

2. **The new init scripts will run** and update the user password automatically

3. **Monitor the deployment** to ensure `supabase-auth` starts successfully

## Verification

After deployment, check that all services are healthy:

```bash
# Check container status
docker compose ps

# Check supabase-auth logs
docker compose logs supabase-auth --tail=50

# Test the health endpoint
curl http://localhost:9999/health

# Run full diagnostics
./diagnose-deployment.sh
```

## Expected Success Indicators

✅ `supabase-postgres` is healthy  
✅ `supabase-init` runs scripts without errors  
✅ `supabase-auth` starts and passes healthcheck  
✅ GoTrue logs show successful database migrations  
✅ Health endpoint returns 200 OK  

## If Issues Persist

1. **Check environment variables in Dokploy**:
   - Ensure `SUPABASE_POSTGRES_PASSWORD` matches the password in the init scripts
   - Ensure `GOTRUE_DB_DATABASE_URL` uses the correct password
   - Verify `GOTRUE_JWT_SECRET` is set

2. **Check the logs**:
   ```bash
   docker compose logs supabase-postgres --tail=100
   docker compose logs supabase-init --tail=100
   docker compose logs supabase-auth --tail=100
   ```

3. **Manually verify the database**:
   ```bash
   docker compose exec supabase-postgres psql -U supabase_auth_admin -d supabase_auth -c "\du"
   ```

4. **Clean slate approach**:
   - Delete all volumes: `docker volume prune`
   - Redeploy from scratch

## Files Modified

- ✅ `init-scripts/00-create-user.sql` (NEW)
- ✅ `init-scripts/01-create-auth-schema.sql` (UPDATED)
- ✅ `docker-compose.yml` (UPDATED)
- ✅ `fix-dokploy-volumes.sh` (NEW)
- ✅ `diagnose-deployment.sh` (NEW)

## Security Note

⚠️ The password `0yvhLSetDKV4BlFOH6YeM5LCBe2jmV2B` is hardcoded in:
- `00-create-user.sql`
- Your `.env` file as `SUPABASE_POSTGRES_PASSWORD`
- The `GOTRUE_DB_DATABASE_URL` connection string

Ensure these all match exactly!
