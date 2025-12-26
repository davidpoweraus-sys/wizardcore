# Dokploy Deployment Guide

Complete guide to deploying WizardCore on Dokploy.

## Prerequisites

- Dokploy installed on your VPS
- Domain name configured
- GitHub repository access

## Deployment Options

### Option 1: Docker Compose (Recommended - Full Stack)

Deploy the entire stack (frontend, backend, databases, auth) with one configuration.

#### Steps:

1. **In Dokploy Dashboard:**
   - Click "Create Project" → Enter "WizardCore"
   - Click "Create Application" → Select "Compose"
   
2. **Configure Source:**
   - **Provider**: GitHub
   - **Repository**: `https://github.com/yourusername/wizardcore.git`
   - **Branch**: `main`
   - **Compose File Path**: `docker-compose.yml` (auto-detected)

3. **Set Environment Variables:**
   
   Click "Environment" tab and add these variables:
   
   ```bash
   # Database Passwords (GENERATE STRONG VALUES!)
   POSTGRES_PASSWORD=<your-strong-password-here>
   SUPABASE_POSTGRES_PASSWORD=<your-strong-password-here>
   JUDGE0_POSTGRES_PASSWORD=<your-strong-password-here>
   
   # Redis Password
   REDIS_PASSWORD=<your-strong-password-here>
   
   # JWT Secret (generate with: openssl rand -base64 64)
   SUPABASE_JWT_SECRET=<your-generated-jwt-secret>
   
   # Judge0 API Key
   JUDGE0_API_KEY=<your-generated-api-key>
   
   # URLs (replace yourdomain.com with your actual domain)
   FRONTEND_URL=https://yourdomain.com
   NEXT_PUBLIC_SUPABASE_URL=https://auth.yourdomain.com
   NEXT_PUBLIC_BACKEND_URL=https://yourdomain.com/api
   NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.yourdomain.com
   API_EXTERNAL_URL=https://auth.yourdomain.com
   GOTRUE_SITE_URL=https://yourdomain.com
   GOTRUE_CORS_ALLOWED_ORIGINS=https://yourdomain.com
   
   # Keep this value as-is
   NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
   ```

4. **Configure Domains (in Dokploy):**
   - **Frontend**: `yourdomain.com` → port `3000`
   - **Auth**: `auth.yourdomain.com` → port `9999`
   - **Judge0**: `judge0.yourdomain.com` → port `2358`
   - **Backend**: Configure as path `/api` on main domain

5. **Deploy:**
   - Click "Deploy"
   - Monitor logs for successful startup
   - All services will start automatically

6. **Verify:**
   ```bash
   curl https://yourdomain.com
   curl https://auth.yourdomain.com/health
   ```

---

### Option 2: Single Dockerfile (Frontend Only)

Deploy just the Next.js frontend. You'll need to manage backend/databases separately.

#### Steps:

1. **In Dokploy Dashboard:**
   - Click "Create Application" → Select "Docker"
   
2. **Configure Source:**
   - **Provider**: GitHub
   - **Repository**: `https://github.com/yourusername/wizardcore.git`
   - **Branch**: `main`
   - **Dockerfile Path**: `Dockerfile`

3. **Set Build Arguments:**
   ```bash
   NEXT_PUBLIC_SUPABASE_URL=https://auth.yourdomain.com
   NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
   NEXT_PUBLIC_BACKEND_URL=https://yourdomain.com/api
   NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.yourdomain.com
   ```

4. **Configure Domain:**
   - **Domain**: `yourdomain.com`
   - **Port**: `3000`

5. **Deploy**

**Note:** This only deploys the frontend. You need to deploy backend, databases, and auth separately.

---

## Docker Compose vs Dockerfile

| Feature | Docker Compose | Single Dockerfile |
|---------|----------------|-------------------|
| **Services** | All (8 containers) | Frontend only |
| **Complexity** | Medium | Low |
| **Setup Time** | 10 minutes | 5 minutes |
| **Maintenance** | Single deployment | Multiple deployments |
| **Scaling** | Swarm-ready | Manual |
| **Best For** | Production | Development/Testing |

**Recommendation**: Use **Docker Compose** for production deployments.

---

## Environment Variable Reference

### Required Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `POSTGRES_PASSWORD` | Main database password | Generate with `openssl rand -hex 32` |
| `SUPABASE_POSTGRES_PASSWORD` | Auth database password | Generate with `openssl rand -hex 32` |
| `SUPABASE_JWT_SECRET` | JWT signing secret | Generate with `openssl rand -base64 64` |
| `REDIS_PASSWORD` | Redis cache password | Generate with `openssl rand -hex 32` |
| `JUDGE0_POSTGRES_PASSWORD` | Judge0 database password | Generate with `openssl rand -hex 32` |
| `JUDGE0_API_KEY` | Judge0 API key | Generate with `openssl rand -hex 32` |

### Domain Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `FRONTEND_URL` | Main application URL | `https://yourdomain.com` |
| `NEXT_PUBLIC_SUPABASE_URL` | Auth service URL | `https://auth.yourdomain.com` |
| `NEXT_PUBLIC_BACKEND_URL` | Backend API URL | `https://yourdomain.com/api` |
| `API_EXTERNAL_URL` | Auth external URL | `https://auth.yourdomain.com` |
| `GOTRUE_SITE_URL` | Auth site URL | `https://yourdomain.com` |

---

## Troubleshooting

### Error: "Building with Dockerfile but docker-compose.yml found"

**Solution**: You selected "Docker" type instead of "Compose"
- Delete the application
- Create new one with "Compose" type
- Select `docker-compose.yml`

### Error: "Service won't start / Health check failed"

**Check logs in Dokploy:**
1. Go to Application → Logs
2. Look for database connection errors
3. Verify environment variables are set correctly

**Common fixes:**
- Ensure all passwords are set (no `CHANGE_THIS` values)
- Check database services started before dependent services
- Verify domain names match in all variables

### Error: "CORS errors in browser"

**Fix:**
1. Ensure `GOTRUE_CORS_ALLOWED_ORIGINS` matches your frontend domain exactly
2. Check `GOTRUE_SITE_URL` is set correctly
3. Verify domains use HTTPS (not HTTP)

### Services taking too long to start

**Normal on first deploy:**
- PostgreSQL initialization: 30-60 seconds
- Next.js build: 2-5 minutes (if building from source)
- Judge0 setup: 30-60 seconds

**Check:**
```bash
docker compose ps
docker compose logs -f
```

---

## Scaling with Swarm

The `docker-compose.yml` is configured for Docker Swarm with replicas:

- **Frontend**: 2 replicas (load balanced)
- **Backend**: 2 replicas (load balanced)
- **Judge0 Workers**: 2 replicas
- **Databases**: 1 replica (pinned to manager node)

### Enable Swarm in Dokploy:

1. **Single Node** (default):
   - Dokploy automatically manages Swarm mode
   - Replicas run on same server
   - Good for < 10K users

2. **Multi-Node** (scaling):
   - Add servers in Dokploy settings
   - Dokploy joins them to Swarm cluster
   - Replicas distribute across nodes
   - Good for > 10K users

---

## Post-Deployment Checklist

- [ ] All containers running (check in Dokploy)
- [ ] Frontend accessible at `https://yourdomain.com`
- [ ] Auth service responds at `https://auth.yourdomain.com/health`
- [ ] Can create user account
- [ ] Can log in
- [ ] Database backups configured in Dokploy
- [ ] SSL certificates auto-renewing
- [ ] Monitoring enabled

---

## Updating the Application

### Via Git Push (Automatic):

1. Make changes locally
2. Commit and push to GitHub
3. In Dokploy, enable "Auto Deploy" for the application
4. Dokploy automatically rebuilds and redeploys on git push

### Manual Deploy:

1. Push changes to GitHub
2. In Dokploy dashboard, click "Redeploy"
3. Monitor logs for successful deployment

### Zero-Downtime Updates (Swarm):

Swarm automatically does rolling updates:
- Updates one replica at a time
- Waits for health check
- Continues to next replica
- No downtime!

---

## Backup Strategy

### Databases (Automatic in Dokploy):

1. Go to Application → Database Backup
2. Enable automatic backups
3. Set schedule (daily recommended)
4. Configure retention (7 days minimum)

### Manual Backup:

```bash
# Backup main database
docker exec postgres pg_dump -U wizardcore wizardcore > backup.sql

# Backup auth database
docker exec supabase-postgres pg_dump -U supabase_auth_admin supabase_auth > auth-backup.sql
```

---

## Performance Tuning

### For Small VPS (2GB RAM):
- Set frontend replicas to 1
- Set backend replicas to 1
- Disable Judge0 workers (use only main instance)

Edit `docker-compose.yml`:
```yaml
frontend:
  deploy:
    replicas: 1

backend:
  deploy:
    replicas: 1
```

### For Large VPS (8GB+ RAM):
- Keep default replicas (2 each)
- Add more Judge0 workers if needed
- Enable Swarm multi-node scaling

---

## Support

- **Dokploy Docs**: https://docs.dokploy.com
- **Issues**: GitHub Issues
- **Logs**: Check Dokploy dashboard → Application → Logs

---

## Quick Reference Commands

```bash
# Generate secrets
openssl rand -base64 64  # For JWT_SECRET
openssl rand -hex 32     # For passwords

# Check service status
docker compose ps

# View logs
docker compose logs -f

# Restart a service
docker compose restart frontend

# Scale a service
docker service scale wizardcore_frontend=3
```

---

## Summary

1. ✅ Choose **Docker Compose** for full-stack deployment
2. ✅ Set all environment variables in Dokploy
3. ✅ Generate strong passwords for production
4. ✅ Configure domains with SSL
5. ✅ Enable automatic backups
6. ✅ Monitor logs on first deploy

Your WizardCore platform will be live in ~10 minutes! 🚀
