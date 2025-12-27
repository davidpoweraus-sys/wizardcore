# CORS Issue Resolution - Complete

## Problem
The dashboard was showing CORS errors when trying to fetch data from:
- `https://api.offensivewizard.com` (Backend API)
- `https://judge0.offensivewizard.com` (Judge0 service)

The errors indicated "CORS request did not succeed" which means the requests weren't even reaching the servers - this is a network-level failure, not a CORS header configuration issue.

## Root Cause
The frontend Next.js app was trying to make direct cross-origin requests to external domains, which can fail due to:
1. Network connectivity issues between services
2. SSL/certificate problems
3. Docker network isolation
4. DNS resolution failures
5. Backend services not running or not accessible

## Solution: API Proxy Pattern

Created Next.js API routes that act as proxies to forward requests from the frontend to the backend services. This is the same pattern used for the authentication proxy.

### What Was Changed

#### 1. Backend API Proxy
**File:** [`app/api/backend/[...path]/route.ts`](app/api/backend/[...path]/route.ts)
- Proxies all requests from frontend to backend API
- Adds authentication tokens from session
- Handles CORS headers properly
- Provides detailed error messages

#### 2. Judge0 API Proxy
**File:** [`app/api/judge0/[...path]/route.ts`](app/api/judge0/[...path]/route.ts)
- Proxies all requests from frontend to Judge0 service
- Handles code execution timeouts (60 seconds)
- Handles CORS headers properly
- Provides detailed error messages

#### 3. Updated API Client
**File:** [`lib/api.ts`](lib/api.ts:1)
- Changed from: `https://api.offensivewizard.com/v1/...`
- Changed to: `/api/backend/v1/...`
- Now uses the proxy instead of direct backend calls

#### 4. Updated Judge0 Service
**File:** [`lib/judge0/service.ts`](lib/judge0/service.ts:1)
- Changed from: `https://judge0.offensivewizard.com/...`
- Changed to: `/api/judge0/...`
- Now uses the proxy instead of direct Judge0 calls

## How It Works

### Before (Direct Requests - FAILING):
```
Frontend (browser) → https://api.offensivewizard.com ❌ CORS Error
Frontend (browser) → https://judge0.offensivewizard.com ❌ CORS Error
```

### After (Proxy Pattern - WORKING):
```
Frontend (browser) → /api/backend → Backend API ✅
Frontend (browser) → /api/judge0 → Judge0 Service ✅
```

All requests now stay within the same origin (your Next.js app), then the server-side proxy forwards them to the actual services.

## Benefits

1. **No CORS Issues**: All requests are same-origin from the browser's perspective
2. **Better Error Handling**: Detailed error messages when services are unreachable
3. **Authentication**: Session tokens automatically added to backend requests
4. **Security**: Backend API keys not exposed to frontend
5. **Flexibility**: Can switch between proxy and direct mode easily
6. **Consistent Pattern**: Same approach used for auth, backend, and Judge0

## Environment Variables

The proxy uses these environment variables:

```env
# Backend API proxy target
BACKEND_URL=https://api.offensivewizard.com

# Judge0 proxy target
JUDGE0_URL=https://judge0.offensivewizard.com

# These are still used for reference but requests now go through proxies
NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com
NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
```

## Testing

After deploying these changes:

1. Open the dashboard at `https://app.offensivewizard.com/dashboard`
2. Check browser console - CORS errors should be gone
3. You should see data loading or clear error messages if services are down
4. Check server logs for proxy request debugging info

## Troubleshooting

### If you still see errors:

1. **"Cannot connect to backend API"**
   - Backend service is not running or not accessible
   - Check `docker-compose ps` to see if backend container is running
   - Verify `BACKEND_URL` environment variable

2. **"Cannot connect to Judge0 server"**
   - Judge0 service is not running or not accessible
   - Check if Judge0 container is running
   - Verify `JUDGE0_URL` environment variable

3. **"Proxy timeout"**
   - Backend is taking too long to respond (>30s for API, >60s for Judge0)
   - Check backend logs for slow queries or issues

4. **Authentication issues**
   - Session might be expired - try logging out and back in
   - Check if auth tokens are being added to requests in proxy logs

## Migration Path

If you want to switch back to direct requests (not recommended):

In [`lib/api.ts`](lib/api.ts:5):
```typescript
const USE_PROXY = false  // Change to false
```

In [`lib/judge0/service.ts`](lib/judge0/service.ts:7):
```typescript
const USE_PROXY = false  // Change to false
```

## Related Files

- [`app/api/auth/[...path]/route.ts`](app/api/auth/[...path]/route.ts) - Auth proxy (same pattern)
- [`middleware.ts`](middleware.ts) - Session handling
- [`next.config.ts`](next.config.ts) - CORS headers configuration

## Next Steps

1. Deploy the changes
2. Verify the dashboard loads without CORS errors
3. If you still see errors, check if backend and Judge0 services are actually running and accessible
4. Review server logs to see detailed proxy request/response information