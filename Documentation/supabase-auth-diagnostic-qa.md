# Supabase Auth Diagnostic Q&A for AI Agent

This document provides a structured set of questions and answers to help diagnose Supabase Auth startup failures. The AI agent can use these to interrogate the system state and identify root causes.

## Primary Symptoms
- Supabase Auth container fails to start
- Error logs mention missing "auth" schema (SQLSTATE 3F000)
- Startup scripts may have run but schema not created

## Diagnostic Questions

### 1. Log Inspection
**Question:** What is the exact error message from the `supabase-auth` container logs?
- **Answer A:** "schema \"auth\" does not exist" (SQLSTATE 3F000)
  - **Diagnosis:** The auth schema is missing. Proceed to question 2.
- **Answer B:** "permission denied for schema auth"
  - **Diagnosis:** The database user lacks privileges. Check grants.
- **Answer C:** "connection refused" or "database not ready"
  - **Diagnosis:** Database dependency issue. Check if `supabase-postgres` is healthy.
- **Answer D:** Other error (provide exact text)
  - **Diagnosis:** Need further analysis.

### 2. Schema Existence Check
**Question:** Is the `auth` schema present in the `supabase_auth` database?
- **Answer A:** Yes, schema exists.
  - **Diagnosis:** The error may be about permissions or something else. Check question 3.
- **Answer B:** No, schema does not exist.
  - **Diagnosis:** The init container (`supabase-init`) may have failed or not run. Check question 4.

### 3. Database User Privileges
**Question:** Does the user `supabase_auth_admin` have ALL privileges on the `auth` schema?
- **Answer A:** Yes, confirmed via `\dn+` and `\z`.
  - **Diagnosis:** Privileges are fine. Look for other issues (e.g., network, configuration).
- **Answer B:** No, missing some or all privileges.
  - **Diagnosis:** Run the grant statements from `init-scripts/01-create-auth-schema.sql`.

### 4. Init Container Status
**Question:** Did the `supabase-init` container run successfully?
- **Answer A:** Yes, logs show "Auth schema initialized."
  - **Diagnosis:** The schema should exist. Verify with direct database query.
- **Answer B:** No, container failed or never started.
  - **Diagnosis:** Check Docker Compose dependencies and health checks.

### 5. Dependency Chain
**Question:** Are the services started in the correct order?
- **Expected order:** `supabase-postgres` (healthy) → `supabase-init` (started) → `supabase-auth` (started)
- **Check:** Run `docker-compose ps` and examine `STATUS` and `DEPENDS ON`.
- **Answer A:** All services are up and healthy.
  - **Diagnosis:** The issue may be transient; restart `supabase-auth`.
- **Answer B:** `supabase-postgres` is unhealthy.
  - **Diagnosis:** Fix database startup (check logs, disk space, credentials).
- **Answer C:** `supabase-init` is missing or exited.
  - **Diagnosis:** Review the init container's command and logs.

### 6. Configuration Variables
**Question:** Are all required environment variables set correctly for `supabase-auth`?
- **Critical variables:**
  - `GOTRUE_DB_DATABASE_URL`
  - `GOTRUE_JWT_SECRET`
  - `GOTRUE_SITE_URL`
- **Answer A:** All are set and match the database connection.
  - **Diagnosis:** Configuration is fine.
- **Answer B:** One or more are missing or incorrect.
  - **Diagnosis:** Update environment variables and restart.

### 7. Network Connectivity
**Question:** Can the `supabase-auth` container reach the `supabase-postgres` container on port 5432?
- **Answer A:** Yes, tested with `nc -zv supabase-postgres 5432` from within the auth container.
  - **Diagnosis:** Network is fine.
- **Answer B:** No, connection refused.
  - **Diagnosis:** Check if PostgreSQL is listening on the correct interface and if firewall rules allow.

### 8. Volume Mounts
**Question:** Is the init script mounted correctly in `supabase-init`?
- **Check:** `docker exec supabase-init ls -la /docker-entrypoint-initdb.d/`
- **Answer A:** File `01-create-auth-schema.sql` is present.
  - **Diagnosis:** Mount is okay.
- **Answer B:** File missing or directory empty.
  - **Diagnosis:** Fix volume mount in Docker Compose.

### 9. Production vs Local
**Question:** Is this happening in production or local environment?
- **Answer A:** Production (Coolify)
  - **Diagnosis:** Ensure `docker-compose.prod.yml` is used and environment variables are set in Coolify.
- **Answer B:** Local (docker-compose.local.yml)
  - **Diagnosis:** Ensure local Docker daemon is running and ports are free.

## Quick Fixes Based on Answers

### If schema missing:
1. Connect to the database:
   ```bash
   docker exec -it wizardcore_supabase-postgres_1 psql -U supabase_auth_admin -d supabase_auth
   ```
2. Create schema manually:
   ```sql
   CREATE SCHEMA IF NOT EXISTS auth;
   GRANT ALL ON SCHEMA auth TO supabase_auth_admin;
   ```
3. Restart `supabase-auth`.

### If init container failed:
1. Check logs:
   ```bash
   docker-compose logs supabase-init
   ```
2. Fix any syntax errors in the command (missing semicolons, etc.).
3. Recreate the init container:
   ```bash
   docker-compose up -d --force-recreate supabase-init
   ```

### If dependencies wrong:
1. Edit `docker-compose.yml` to ensure:
   ```yaml
   supabase-init:
     depends_on:
       supabase-postgres:
         condition: service_healthy
   supabase-auth:
     depends_on:
       supabase-init:
         condition: service_started
   ```
2. Redeploy.

## Automated Diagnostic Script (Optional)
If the AI agent has shell access, it can run `scripts/ensure-auth-schema.sh` (already created) to verify and fix the schema.

## Conclusion
Use this Q&A tree to systematically identify the cause of Supabase Auth startup failures. Most issues are resolved by ensuring the `auth` schema exists and the init container runs before the auth service.