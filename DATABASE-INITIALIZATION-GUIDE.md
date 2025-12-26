# Database Initialization Guide

## Overview

Your Wizardcore deployment uses **THREE separate PostgreSQL databases**:

1. **Supabase Auth Database** (`supabase_auth`) - For authentication
2. **Main Wizardcore Database** (`wizardcore`) - For application data
3. **Judge0 Database** (`judge0`) - For code execution service

## Database Initialization Flow

### 1. Supabase Auth Database (`supabase-postgres`)

**Initialization Method**: Manual init scripts via `supabase-init` container

**Startup Sequence**:
```
1. supabase-postgres container starts
   └─> Creates database & user from POSTGRES_USER env var
   └─> Healthcheck: pg_isready succeeds

2. supabase-init container starts (waits for postgres healthy)
   └─> Runs: 00-create-user.sql (creates/updates auth user)
   └─> Runs: 01-create-auth-schema.sql (creates auth schema & types)
   └─> Writes status file
   └─> Healthcheck: checks status file exists

3. supabase-auth container starts (waits for init healthy)
   └─> GoTrue connects to database
   └─> Runs internal GoTrue migrations automatically
   └─> Creates auth.users, auth.sessions, etc.
```

**Init Scripts**:
- `init-scripts/00-create-user.sql` - Creates `supabase_auth_admin` user
- `init-scripts/01-create-auth-schema.sql` - Creates auth schema, enum types, grants permissions

**Important**: These scripts run **on every deployment** (not just first time), making them idempotent for password updates.

---

### 2. Main Wizardcore Database (`postgres`)

**Initialization Method**: Backend runs migrations on startup

**Startup Sequence**:
```
1. postgres container starts
   └─> Creates database & user from POSTGRES_USER env var
   └─> Healthcheck: pg_isready succeeds
   └─> NO init scripts run (no volume mounted)

2. backend container starts (waits for postgres healthy)
   └─> Backend application launches
   └─> database.RunMigrations() executes
   └─> Applies migrations from internal/database/migrations/ in order:
       - 001_users.up.sql
       - 002_pathways.up.sql
       - 003_exercises.up.sql
       - 004_submissions.up.sql
       - 005_progress.up.sql
       - 006_achievements.up.sql
       - 007_leaderboard.up.sql
       - 008_matches.up.sql
       - 009_notifications.up.sql
       - 010_deadlines.up.sql
       - 011_certificates.up.sql
       - 012_content_creators.up.sql
   └─> Application starts serving requests
```

**Current Migrations Create**:
- ✅ `users` table with supabase_user_id, email, XP, streaks
- ✅ `pathways`, `modules`, `lessons` (learning content)
- ✅ `exercises` with test cases
- ✅ `submissions` (code submissions)
- ✅ `user_progress` tracking
- ✅ `achievements` system
- ✅ `leaderboard` rankings
- ✅ `matches` (competitive coding)
- ✅ `notifications` system
- ✅ `deadlines` and `certificates`
- ✅ `content_creator_profiles` with simple role column

---

### 3. Judge0 Database (`judge0-postgres`)

**Initialization Method**: Judge0 service runs its own migrations

**Startup Sequence**:
```
1. judge0-postgres container starts
   └─> Creates database & user from JUDGE0_POSTGRES_USER env var
   └─> Healthcheck: pg_isready succeeds

2. judge0 container starts (waits for postgres healthy)
   └─> Judge0 service connects
   └─> Runs internal migrations automatically
   └─> Creates submissions, languages, etc.
```

**No custom init scripts needed** - Judge0 handles its own schema.

---

## RBAC Schema Situation

### Current Status: ⚠️ NOT DEPLOYED

You have two RBAC-related files in the project root:
- `/wizardcore/rbac-schema.sql` - Full RBAC system (roles, permissions, inheritance)
- `/wizardcore/create-default-permissions.sql` - Default permissions & role assignments

**These files are NOT currently being run** during deployment!

### What Exists Instead

Migration `012_content_creators.up.sql` added a **simple role column**:
```sql
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'student' 
CHECK (role IN ('student', 'content_creator', 'admin'));
```

This is a **basic role system**, not the full RBAC with:
- ❌ Permission granularity
- ❌ Role inheritance
- ❌ Permission categories
- ❌ Audit logging
- ❌ Dynamic permission checks

---

## Recommendations

### Option 1: Deploy Full RBAC System (Recommended for Production)

**Create migration 013:**

```bash
# Create new migration files
cd wizardcore-backend/internal/database/migrations
```

**File: `013_rbac_system.up.sql`**
```sql
-- Copy contents from /wizardcore/rbac-schema.sql
-- Then add contents from /wizardcore/create-default-permissions.sql
```

**File: `013_rbac_system.down.sql`**
```sql
DROP TABLE IF EXISTS permission_audit_log CASCADE;
DROP TABLE IF EXISTS role_audit_log CASCADE;
DROP TABLE IF EXISTS role_inheritance CASCADE;
DROP TABLE IF EXISTS user_roles CASCADE;
DROP TABLE IF EXISTS role_permissions CASCADE;
DROP TABLE IF EXISTS permissions CASCADE;
DROP TABLE IF EXISTS permission_categories CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
DROP FUNCTION IF EXISTS check_user_permission CASCADE;
DROP FUNCTION IF EXISTS get_user_permissions CASCADE;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;
```

**Then redeploy** - backend will automatically run the new migration.

---

### Option 2: Keep Simple Role System (Current State)

If you don't need complex permissions yet:
- ✅ Users have one role: student, content_creator, or admin
- ✅ Simple to implement in backend code
- ✅ Already deployed via migration 012
- ❌ No fine-grained permissions
- ❌ No audit trail
- ❌ Hard to extend later

---

### Option 3: Add RBAC via Init Scripts (Not Recommended)

You could mount the RBAC scripts to postgres:

```yaml
postgres:
  volumes:
    - postgres_data:/var/lib/postgresql/data
    - ./wizardcore-backend/rbac-init:/docker-entrypoint-initdb.d
```

**Problems**:
- Init scripts only run on FIRST database creation
- Won't run if volume already exists
- Doesn't fit the migration-based approach
- Can cause version conflicts

---

## What Happens on Fresh Deployment (Clean Database)

### First Time Deployment:

1. **Supabase Auth DB**:
   - Creates empty database
   - Runs init scripts (user + schema)
   - GoTrue adds auth tables
   - ✅ Ready for authentication

2. **Wizardcore DB**:
   - Creates empty database
   - Backend runs all 12 migrations
   - Creates all application tables
   - ✅ Ready with simple role system
   - ❌ No RBAC tables (unless you add migration 013)

3. **Judge0 DB**:
   - Creates empty database
   - Judge0 adds its tables
   - ✅ Ready for code execution

### Subsequent Deployments:

1. **Supabase Auth DB**:
   - Existing database persists (volume)
   - Init scripts run again (idempotent)
   - Updates passwords if changed
   - GoTrue checks migrations (no new changes)

2. **Wizardcore DB**:
   - Existing database persists (volume)
   - Backend checks for new migrations
   - Only runs NEW migrations (e.g., if you add 013)
   - Existing data preserved

3. **Judge0 DB**:
   - Existing database persists
   - Judge0 checks migrations
   - Updates if needed

---

## Summary Answer to Your Question

> "Does that spawn with the data needed to launch the app or do we need a script for that as well?"

**Answer**: 

**For Supabase Auth**: ✅ **YES** - The init scripts we created handle everything

**For Main Wizardcore DB**: ⚠️ **PARTIALLY**
- ✅ Backend auto-runs migrations for core tables
- ✅ Application can launch and run
- ❌ RBAC tables (`roles`, `permissions`, etc.) are NOT created unless you add migration 013
- The simple `role` column in `users` table works, but you won't have the full permission system

**For Judge0**: ✅ **YES** - Judge0 handles its own schema

---

## Recommended Next Steps

1. **Decide on RBAC approach** (simple roles vs full RBAC)

2. **If using full RBAC**, create migration 013:
   ```bash
   # Combine rbac-schema.sql + create-default-permissions.sql
   # into 013_rbac_system.up.sql
   ```

3. **Test locally** before deploying:
   ```bash
   docker compose down -v  # Clean slate
   docker compose up -d
   ./diagnose-deployment.sh
   ```

4. **Verify all services start** and migrations complete

5. **Deploy to Dokploy**

---

## Files Reference

```
wizardcore/
├── init-scripts/              # For supabase-postgres only
│   ├── 00-create-user.sql    # ✅ Creates auth user
│   └── 01-create-auth-schema.sql  # ✅ Creates auth schema
│
├── rbac-schema.sql           # ⚠️ NOT CURRENTLY USED
├── create-default-permissions.sql  # ⚠️ NOT CURRENTLY USED
│
└── wizardcore-backend/
    └── internal/database/migrations/  # For main postgres
        ├── 001_users.up.sql       # ✅ Run by backend
        ├── 002_pathways.up.sql    # ✅ Run by backend
        ├── ...
        ├── 012_content_creators.up.sql  # ✅ Run by backend (simple roles)
        └── 013_rbac_system.up.sql      # ❌ NEEDS TO BE CREATED
```
