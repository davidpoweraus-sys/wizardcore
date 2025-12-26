# Debug Proxy 404 Error

## 🔍 Current Status

- ✅ Proxy route exists (`/api/supabase-proxy`)
- ❌ Returns 404 when trying to reach supabase-auth
- ❌ Error: "JSON.parse: unexpected non-whitespace character"

This means the proxy **can't connect** to `http://supabase-auth:9999`

## 🐛 Possible Causes

### 1. Docker Network Issue (Most Likely)
Coolify might create **separate networks** for each service, preventing frontend from reaching supabase-auth.

### 2. Service Name Resolution
The hostname `supabase-auth` might not resolve within the frontend container.

### 3. Supabase Auth Not Running
The supabase-auth service might be stopped or crashed.

## 🔬 Diagnosis Steps

### Step 1: Check Frontend Logs (IMPORTANT!)

**In Coolify:**
1. Go to **frontend service** → **Runtime Logs**
2. Try signup on test page again
3. **Look for these log messages:**

```
🔄 Proxy Configuration:
  GOTRUE_URL: http://supabase-auth:9999
  Target Path: auth/v1/signup
  Full URL: http://supabase-auth:9999/auth/v1/signup
  Method: POST
📤 Making request to Supabase Auth...
```

**Then look for:**

**If you see:**
```
✅ Proxy response status: 200
```
→ It's working! The 404 is something else.

**If you see:**
```
❌ Proxy error details:
  Error message: fetch failed
  ⚠️  Cannot reach supabase-auth service
```
→ Network/DNS issue confirmed.

### Step 2: Check if Supabase Auth is Running

**In Coolify:**
1. Go to **supabase-auth service**
2. Check status - should be **"Running"** (green)
3. If not running, start it

### Step 3: Test Network Connectivity

**If you have SSH access to Coolify server:**

```bash
# Get container names
docker ps | grep -E "frontend|supabase-auth"

# Test DNS resolution from frontend container
docker exec <frontend-container-name> nslookup supabase-auth

# Test connectivity
docker exec <frontend-container-name> wget -qO- http://supabase-auth:9999/health

# Check if they're on the same network
docker inspect <frontend-container> | grep NetworkMode
docker inspect <supabase-auth-container> | grep NetworkMode
```

**Expected:**
- Both should be on `wizardcore-network` or Coolify's shared network
- `wget` should return: `{"version":"...","name":"GoTrue"}`

## ✅ Solutions

### Solution 1: Use Coolify's Internal Network

Coolify creates internal networks. We might need to use Coolify's network name instead.

**Check in Coolify:**
- Look at the network configuration for both services
- They should share a common network

### Solution 2: Use Direct IP Instead of Hostname

If DNS doesn't work, we can use the container IP:

```bash
# Find supabase-auth container IP
docker inspect <supabase-auth-container> | grep IPAddress
```

Then update the proxy URL to use that IP (not ideal, but works for testing).

### Solution 3: Use Backend Service as Proxy

Since your Go backend can already connect to supabase-auth (it uses it for JWT verification), we could proxy through the backend instead:

**Flow:**
```
Browser → Next.js → Go Backend → Supabase Auth
```

This would work because backend already has network access.

### Solution 4: Change to External URL (Temporary Fix)

For testing, we can make the proxy use the external URL:

```typescript
// In route.ts, change:
const GOTRUE_URL = process.env.SUPABASE_INTERNAL_URL || 'https://auth.offensivewizard.com'
```

This would make the proxy hit the external URL, which would still solve CORS (proxy adds headers), but adds extra latency.

## 🚀 Immediate Next Steps

### 1. Check Logs (Do This Now!)

After this deploys with enhanced logging:
1. Go to test page: https://offensivewizard.com/test-signup
2. Try signup with Proxy ON
3. **Go to Coolify → frontend → Runtime Logs**
4. **Share the log output here**

The logs will show exactly what's happening.

### 2. Verify Services Running

Check in Coolify dashboard:
- [ ] supabase-auth: Status = Running?
- [ ] frontend: Status = Running?
- [ ] backend: Status = Running?

### 3. Check Network Configuration

In Coolify:
- Go to supabase-auth service → Network tab
- Go to frontend service → Network tab
- **Do they share a common network name?**

## 💡 Alternative: Just Enable CORS on Supabase Auth

I noticed in the docker-compose that CORS is set to `*`:

```yaml
GOTRUE_CORS_ALLOWED_ORIGINS: "*"
```

But the issue is that with wildcard (`*`), you can't use `credentials: true`.

**Simpler fix:** Change this to specific origin:

```yaml
GOTRUE_CORS_ALLOWED_ORIGINS: "https://offensivewizard.com"
GOTRUE_CORS_ALLOW_CREDENTIALS: "true"
```

But wait... I thought we already set this! Let me check...

Actually, I see the issue now. In your current docker-compose, it says:

```yaml
GOTRUE_CORS_ALLOWED_ORIGINS: "*"
```

This should be:
```yaml
GOTRUE_CORS_ALLOWED_ORIGINS: "https://offensivewizard.com"
```

**But there's a problem:** Coolify's Traefik is probably blocking the preflight before it even reaches GoTrue.

So the proxy approach is still best, **IF** we can get the network connectivity working.

## 📊 Priority Actions

**Highest Priority:**
1. ✅ Check frontend runtime logs after trying signup
2. ✅ Share log output
3. ✅ Verify both services on same network

**Medium Priority:**
4. Test network connectivity between containers
5. Check if external URL works as workaround

**Low Priority:**
6. Consider backend proxy alternative
7. Investigate Coolify network configuration

---

**Next: After redeployment completes, test signup and check logs!**
