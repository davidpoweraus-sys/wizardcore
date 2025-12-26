# CapRover Quick Start for WizardCore

## **🚀 Simple 5-Step Deployment**

### **Step 1: Clean Up Coolify (if needed)**
```bash
# Stop and remove Coolify containers
docker stop coolify-proxy coolify-sentinel 2>/dev/null || true
docker rm coolify-proxy coolify-sentinel 2>/dev/null || true

# Stop WizardCore containers
docker-compose -f docker-compose.coolify.yml down 2>/dev/null || true
```

### **Step 2: Install CapRover (if not installed)**
```bash
# Install Docker (if not installed)
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install CapRover
docker run -p 80:80 -p 443:443 -p 3000:3000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v captain-data:/data \
  caprover/caprover
```

### **Step 3: Access CapRover Dashboard**
1. Open browser to: `http://your-server-ip:3000`
2. Complete setup wizard
3. Set root domain (e.g., `your-server.com`)

### **Step 4: Deploy WizardCore**

#### **Option A: Docker Compose (Recommended)**
1. In CapRover dashboard, create app: `wizardcore`
2. Go to **App Configs** → **Docker Compose Editor**
3. Paste contents of `docker-compose.caprover.yml`
4. Click **Save & Update**

#### **Option B: CLI Deployment**
```bash
# Install CapRover CLI
npm install -g caprover

# Login
caprover login --caproverUrl "https://captain.your-domain.com"

# Deploy
caprover deploy --appName wizardcore --imageName limpet/wizardcore-frontend:latest
```

### **Step 5: Configure Domain & Proxy**
1. In app settings, add custom domain: `offensivewizard.com`
2. Enable HTTPS
3. Add custom nginx configuration (see below)

## **🔧 Required Nginx Configuration**

In **App Configs** → **HTTP Settings** → **Custom Nginx Configuration**:

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

## **✅ Verification Commands**

```bash
# Check if services are running
docker ps | grep wizardcore

# Test frontend
curl https://offensivewizard.com

# Test backend
curl https://offensivewizard.com/api/health

# Test auth proxy
curl https://offensivewizard.com/supabase-proxy/health
```

## **🔍 Troubleshooting**

### **Issue: "Connection refused"**
```bash
# Check container logs
docker logs wizardcore-frontend
docker logs wizardcore-backend

# Check if containers are in same network
docker network inspect wizardcore_wizardcore-network
```

### **Issue: Database connection failed**
```bash
# Check PostgreSQL
docker exec wizardcore-postgres pg_isready -U wizardcore

# Check environment variables
docker exec wizardcore-backend env | grep DATABASE
```

### **Issue: Proxy not routing correctly**
1. Verify nginx configuration in CapRover
2. Check if services are reachable internally:
```bash
docker exec wizardcore-frontend curl http://backend:8080/health
```

## **📞 Need Help?**

1. Check CapRover documentation: https://caprover.com/docs/
2. View app logs in CapRover dashboard
3. Ensure all Docker images are available:
   - `limpet/wizardcore-frontend:latest`
   - `limpet/wizardcore-backend:latest`
   - `supabase/gotrue:v2.184.0`
   - `postgres:15-alpine`
   - `redis:7-alpine`
   - `judge0/judge0:latest`