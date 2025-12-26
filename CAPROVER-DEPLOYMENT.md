# CapRover Deployment Guide for WizardCore

## **📋 Prerequisites**

1. **CapRover installed** on your server
2. **Domain configured** with DNS pointing to your CapRover instance
3. **GitHub repository access** (or Docker Hub access)

## **🚀 Deployment Methods**

### **Method 1: Docker Compose Deployment (Recommended)**

#### **Step 1: Prepare Your Repository**
1. Ensure you have `docker-compose.caprover.yml` in your repository root
2. Make sure all images are available on Docker Hub

#### **Step 2: Deploy via CapRover Dashboard**
1. **Login to CapRover** dashboard (`captain.your-domain.com`)
2. **Create a new app**: `wizardcore` (or your preferred name)
3. **Enable HTTPS** and assign your domain: `offensivewizard.com`
4. **Go to App Configs** → **Docker Compose Editor**
5. **Paste the contents** of `docker-compose.caprover.yml`
6. **Click Save & Update**

#### **Step 3: Configure Environment Variables**
CapRover will automatically use:
- `${CAPROVER_APP_DOMAIN}` - Your app's domain (e.g., `offensivewizard.com`)
- `${CAPROVER_ROOT_DOMAIN}` - CapRover's root domain (e.g., `your-server.com`)

### **Method 2: Git Deployment**

#### **Step 1: Configure CapRover Git Deployment**
1. In CapRover app settings, go to **Deployment**
2. Enable **Git Deployment**
3. Add your GitHub repository: `https://github.com/davidpoweraus-sys/wizardcore`
4. Set branch to `main`

#### **Step 2: Configure Build Settings**
1. Go to **App Configs** → **Build & Deploy**
2. Set **Dockerfile Path**: `Dockerfile.caprover` (create this file)
3. Set **Context**: `/`

#### **Step 3: Create Dockerfile.caprover**
```dockerfile
# Multi-stage build for CapRover
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY . .
RUN npm ci
RUN npm run build

FROM golang:1.24-alpine AS backend-builder
WORKDIR /app
COPY wizardcore-backend .
RUN go build -o main ./cmd/api

FROM alpine:latest
RUN apk add --no-cache nginx
COPY --from=frontend-builder /app/.next/standalone /usr/share/nginx/html
COPY --from=frontend-builder /app/.next/static /usr/share/nginx/html/.next/static
COPY --from=backend-builder /app/main /usr/local/bin/backend
COPY nginx.conf /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

### **Method 3: One-Click App Template**

Create a `captain-definition` file for One-Click Apps:

```json
{
  "schemaVersion": 2,
  "templateId": "wizardcore",
  "templateVersion": "1.0.0",
  "displayName": "WizardCore Learning Platform",
  "description": "Complete cybersecurity learning platform with code execution",
  "dockerCompose": "docker-compose.caprover.yml",
  "variables": [
    {
      "id": "DOMAIN",
      "label": "Application Domain",
      "description": "Domain for your WizardCore instance",
      "defaultValue": "offensivewizard.com"
    }
  ]
}
```

## **🔧 Post-Deployment Configuration**

### **1. Configure Reverse Proxy Routes**
CapRover's nginx needs to route:
- `/` → `frontend:3000`
- `/api/*` → `backend:8080`
- `/supabase-proxy/*` → `supabase-auth:9999`

Add this to **App Configs** → **HTTP Settings** → **Custom Nginx Configuration**:

```nginx
location /api/ {
    proxy_pass http://backend:8080/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
}

location /supabase-proxy/ {
    proxy_pass http://supabase-auth:9999/;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    
    # Strip /supabase-proxy prefix
    rewrite ^/supabase-proxy/(.*)$ /$1 break;
}
```

### **2. Environment Variables**
In **App Configs** → **Environment Variables**, add:

```
NEXT_PUBLIC_SUPABASE_URL=https://offensivewizard.com/supabase-proxy
NEXT_PUBLIC_SUPABASE_ANON_KEY=uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=
NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.your-server.com
NEXT_PUBLIC_BACKEND_URL=https://offensivewizard.com/api
```

### **3. Persistent Volumes**
Ensure these volumes are configured in CapRover:
- `supabase_postgres_data` (Supabase Auth database)
- `postgres_data` (Main application database)
- `redis_data` (Redis cache)
- `judge0_postgres_data` (Judge0 database)
- `judge0_redis_data` (Judge0 Redis)

## **🔍 Health Checks**

### **1. Database Health**
```bash
# PostgreSQL health
docker exec wizardcore-postgres pg_isready -U wizardcore -d wizardcore

# Redis health
docker exec wizardcore-redis redis-cli ping
```

### **2. Service Health**
```bash
# Frontend
curl https://offensivewizard.com

# Backend API
curl https://offensivewizard.com/api/health

# Supabase Auth
curl https://offensivewizard.com/supabase-proxy/health
```

## **🚨 Troubleshooting**

### **Issue: Services can't communicate**
**Solution**: Ensure all services are in the same network (`wizardcore-network`)

### **Issue: Database connection failures**
**Solution**: Check environment variables and ensure PostgreSQL is healthy

### **Issue: Proxy routing not working**
**Solution**: Verify custom nginx configuration in CapRover

### **Issue: SSL/HTTPS issues**
**Solution**: Ensure domain is properly configured in CapRover and DNS is propagated

## **📊 Monitoring**

### **1. CapRover Dashboard**
- Check app status in CapRover dashboard
- View logs for each service
- Monitor resource usage

### **2. Custom Monitoring**
```bash
# Check all containers
docker ps --filter "name=wizardcore"

# View logs
docker logs wizardcore-frontend
docker logs wizardcore-backend

# Check resource usage
docker stats wizardcore-frontend wizardcore-backend wizardcore-postgres
```

## **🔄 Updates & Maintenance**

### **1. Update Images**
```bash
# Pull latest images
docker pull limpet/wizardcore-frontend:latest
docker pull limpet/wizardcore-backend:latest

# Restart services via CapRover dashboard
```

### **2. Backup Databases**
```bash
# Backup PostgreSQL
docker exec wizardcore-postgres pg_dump -U wizardcore wizardcore > backup.sql

# Backup volumes are automatically managed by CapRover
```

### **3. Scale Services**
In CapRover dashboard:
- Go to **App Configs** → **Horizontal Scaling**
- Adjust instance count for each service

## **🎯 Quick Start Command**

If you want to deploy quickly:

```bash
# 1. Install CapRover (if not already installed)
curl -sSL https://get.caprover.com | bash

# 2. Deploy WizardCore
# Use the CapRover dashboard or CLI:
caprover deploy --app wizardcore --image limpet/wizardcore-frontend:latest
```

## **📞 Support**

For issues:
1. Check CapRover logs
2. Verify DNS configuration
3. Ensure ports 80, 443, 3000 are open
4. Check Docker image availability on Docker Hub