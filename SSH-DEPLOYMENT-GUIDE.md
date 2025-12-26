# SSH Deployment Guide for WizardCore

This guide provides a simple, direct deployment method using SSH instead of CapRover's One-Click App.

## Quick Start

### 1. **Copy deployment scripts to your server:**

```bash
# On your local machine
scp deploy-ssh.sh setup-nginx.sh root@your-server-ip:/tmp/
```

### 2. **SSH into your server and run:**

```bash
ssh root@your-server-ip

# Make scripts executable
chmod +x /tmp/deploy-ssh.sh /tmp/setup-nginx.sh

# Run deployment
/tmp/deploy-ssh.sh
```

### 3. **Set up Nginx reverse proxy:**

```bash
/tmp/setup-nginx.sh
```

### 4. **Configure DNS:**

Add these A records pointing to your server IP:
- `app.offensivewizard.com`
- `api.offensivewizard.com`
- `auth.offensivewizard.com`
- `judge0.offensivewizard.com`

### 5. **Set up SSL (after DNS propagates):**

```bash
cd /opt/wizardcore
./setup-ssl.sh
```

## What Gets Deployed

### Services:
1. **Frontend** (`app.offensivewizard.com`) - Port 3000
2. **Backend API** (`api.offensivewizard.com`) - Port 8080
3. **Supabase Auth** (`auth.offensivewizard.com`) - Port 9999
4. **Judge0** (`judge0.offensivewizard.com`) - Port 2358
5. **PostgreSQL databases** (3 instances) - Internal only
6. **Redis caches** (2 instances) - Internal only

### Management Commands:
- `wizardcore-start` - Start all services
- `wizardcore-stop` - Stop all services
- `wizardcore-restart` - Restart all services
- `wizardcore-logs` - View logs (follow mode)
- `wizardcore-status` - Check service status
- `wizardcore-update` - Update to latest images

## Directory Structure

```
/opt/wizardcore/
├── docker-compose.yml    # Main deployment configuration
├── .env                  # Environment variables
├── setup-ssl.sh         # SSL certificate setup
└── data/                # Persistent data volumes
```

## Troubleshooting

### Check service status:
```bash
wizardcore-status
```

### View logs:
```bash
wizardcore-logs
```

### Test individual services:
```bash
# Frontend
curl http://localhost:3000

# Backend
curl http://localhost:8080/health

# Auth
curl http://localhost:9999/health

# Judge0
curl http://localhost:2358/about
```

### Check Nginx:
```bash
# Test configuration
nginx -t

# Check logs
tail -f /var/log/nginx/error.log
tail -f /var/log/nginx/access.log

# Restart Nginx
systemctl restart nginx
```

### Common Issues:

1. **Ports already in use:**
   ```bash
   netstat -tulpn | grep :3000
   # Kill process or change port in docker-compose.yml
   ```

2. **Database connection errors:**
   ```bash
   # Check if databases are running
   docker-compose ps | grep postgres
   
   # View database logs
   docker logs wizardcore-postgres-1
   ```

3. **CORS errors:**
   - Ensure Nginx CORS headers are correct
   - Check environment variables match your domain

4. **SSL certificate issues:**
   ```bash
   # Renew certificates
   certbot renew
   
   # Check certificate status
   certbot certificates
   ```

## Updating

### Update to latest version:
```bash
wizardcore-update
```

### Manual update:
```bash
cd /opt/wizardcore
git pull origin main  # If using git
docker-compose pull
docker-compose down
docker-compose up -d
```

## Backup and Restore

### Backup databases:
```bash
cd /opt/wizardcore

# Main database
docker exec wizardcore-postgres-1 pg_dump -U wizardcore wizardcore > backup-main.sql

# Auth database
docker exec wizardcore-auth-db-1 pg_dump -U supabase_auth_admin supabase_auth > backup-auth.sql

# Judge0 database
docker exec wizardcore-judge0-db-1 pg_dump -U judge0 judge0 > backup-judge0.sql
```

### Restore databases:
```bash
cd /opt/wizardcore

# Main database
cat backup-main.sql | docker exec -i wizardcore-postgres-1 psql -U wizardcore wizardcore

# Auth database
cat backup-auth.sql | docker exec -i wizardcore-auth-db-1 psql -U supabase_auth_admin supabase_auth

# Judge0 database
cat backup-judge0.sql | docker exec -i wizardcore-judge0-db-1 psql -U judge0 judge0
```

## Monitoring

### Check resource usage:
```bash
# Docker container stats
docker stats

# Disk usage
df -h /opt/wizardcore

# Memory usage
free -h
```

### Set up monitoring (optional):
```bash
# Install monitoring tools
apt-get install -y htop nmon

# View real-time stats
htop
```

## Security Notes

1. **Firewall configuration:**
   ```bash
   # Allow necessary ports
   ufw allow 22/tcp    # SSH
   ufw allow 80/tcp    # HTTP
   ufw allow 443/tcp   # HTTPS
   ufw enable
   ```

2. **Regular updates:**
   ```bash
   apt-get update && apt-get upgrade -y
   ```

3. **Backup strategy:**
   - Schedule daily database backups
   - Store backups off-site
   - Test restore procedures regularly

## Support

If you encounter issues:
1. Check logs: `wizardcore-logs`
2. Verify DNS configuration
3. Ensure ports are not blocked by firewall
4. Check disk space and memory

For persistent issues, check the `/opt/wizardcore/docker-compose.yml` configuration and adjust as needed for your environment.