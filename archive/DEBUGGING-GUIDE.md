# Debugging Supabase Auth CORS Error

## Problem
OPTIONS request to `https://auth.offensivewizard.com/auth/v1/signup` returns "undefined". Frontend at `https://offensivewizard.com` cannot authenticate.

## Quick Fix Attempted
1. Fixed Docker Compose deployment error ("no service selected") by removing `version` field
2. Set CORS to allow any origin (`"*"`) in `docker-compose.prod.yml`

## Steps to Diagnose

### 1. Check if Deployment Succeeded
- Go to Coolify dashboard
- Check if all services are running (green status)
- Look for any error messages in deployment logs

### 2. Check Supabase Auth Logs
```bash
# In Coolify, view logs for supabase-auth container
# Or via CLI if you have access:
docker logs <supabase-auth-container-id>
```

Look for:
- Startup errors
- Database connection errors
- CORS configuration logs

### 3. Test if Service is Accessible
```bash
# Test health endpoint
curl -v -X GET "https://auth.offensivewizard.com/auth/v1/health" --insecure

# Test OPTIONS request
curl -v -X OPTIONS "https://auth.offensivewizard.com/auth/v1/signup" \
  -H "Origin: https://offensivewizard.com" \
  -H "Access-Control-Request-Method: POST" \
  --insecure
```

Expected response: HTTP 200 or 204 with CORS headers.

### 4. Check Database Connection
Supabase Auth needs to connect to PostgreSQL. Check:
- Is `supabase-postgres` container running?
- Can `supabase-auth` connect to it?
- Are the credentials correct?

### 5. Check Environment Variables
In Coolify application settings, ensure:
- `NEXT_PUBLIC_SUPABASE_URL` = `https://auth.offensivewizard.com`
- `NEXT_PUBLIC_SUPABASE_ANON_KEY` = `uc8bo6Z4ZI4Fhu9XVgSz5LhDRWEQ0joGPMiZYroXPps=`

### 6. Check SSL Certificate
- Browser might reject self-signed certificate for `auth.offensivewizard.com`
- Test with `--insecure` flag in curl
- Consider using Let's Encrypt or valid SSL certificate

## If Service is Not Running

### Possible Causes:
1. **Database connection failed**: Check `supabase-postgres` logs
2. **Missing auth schema**: Check `supabase-init` logs
3. **Insufficient resources**: Check Coolify resource limits
4. **Network issues**: Check Docker network configuration

### Quick Workaround for Testing:
Temporarily modify `docker-compose.prod.yml` to remove healthcheck dependencies:
```yaml
depends_on:
  supabase-postgres:
    condition: service_started  # instead of service_healthy
```

## If Service is Running but CORS Fails

### Check CORS Headers:
Response should include:
```
Access-Control-Allow-Origin: https://offensivewizard.com
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Authorization, Content-Type, X-Client-Info, X-Requested-With, apikey, x-client-info, x-supabase-api-version
Access-Control-Allow-Credentials: true
```

### If using wildcard (*):
- Remove `Access-Control-Allow-Credentials: true` or set to `false`
- Or use specific origin instead of wildcard

## Immediate Action Plan

1. **Redeploy** with updated `docker-compose.prod.yml`
2. **Check logs** in Coolify for errors
3. **Test connectivity** with curl commands above
4. **Verify frontend configuration** - ensure `NEXT_PUBLIC_SUPABASE_URL` points to correct URL

## If Still Failing
1. Check browser console for exact error
2. Check network tab for request/response details
3. Consider temporary workaround: run auth on same domain as frontend (no CORS)
4. Or disable CORS in browser for testing (not recommended for production)