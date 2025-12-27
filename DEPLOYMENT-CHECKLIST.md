# Deployment Checklist - CORS Fix Update

## Pre-Deployment

- [ ] Review all changes:
  - [`app/api/backend/[...path]/route.ts`](app/api/backend/[...path]/route.ts) - Backend API proxy
  - [`app/api/judge0/[...path]/route.ts`](app/api/judge0/[...path]/route.ts) - Judge0 API proxy
  - [`lib/api.ts`](lib/api.ts) - Updated to use proxy
  - [`lib/judge0/service.ts`](lib/judge0/service.ts) - Updated to use proxy
  - [`docker-compose.yml`](docker-compose.yml) - Added proxy environment variables
  - [`.env.example`](/.env.example) - Updated documentation

- [ ] Ensure you have a `.env` file with all required variables (copy from `.env.example` if needed)

- [ ] Make sure you're logged into Docker Hub:
  ```bash
  docker login
  ```

## Build and Push

1. **Build the Docker image with all changes:**
   ```bash
   ./build-and-push.sh
   ```
   
   This script will:
   - Load environment variables from `.env`
   - Build the Docker image with the CORS fixes included
   - Tag it as `limpet/wizardcore-frontend:latest`
   - Push to Docker Hub

2. **Verify the build includes the new files:**
   ```bash
   docker run --rm limpet/wizardcore-frontend:latest ls -la /app/app/api/
   ```
   
   You should see:
   - `backend/` directory
   - `judge0/` directory
   - `auth/` directory

## Deploy to Server

1. **SSH into your server:**
   ```bash
   ssh your-server
   ```

2. **Navigate to your project directory:**
   ```bash
   cd /path/to/wizardcore
   ```

3. **Update your `.env` file to include the new proxy URLs:**
   ```bash
   # Add these lines to your .env file
   BACKEND_URL=http://backend:8080
   JUDGE0_URL=http://judge0:2358
   GOTRUE_URL=http://supabase-auth:9999
   ```

4. **Pull the latest image:**
   ```bash
   docker-compose pull frontend
   ```

5. **Restart the frontend service:**
   ```bash
   docker-compose up -d frontend
   ```

6. **Wait for the service to be healthy:**
   ```bash
   docker-compose ps
   ```
   
   Wait until the frontend shows as "healthy"

7. **Check the logs to ensure no errors:**
   ```bash
   docker-compose logs -f frontend
   ```
   
   Look for:
   - `✅ Compiled successfully` (Next.js compilation)
   - No errors related to missing files
   - Proxy request logs when you access the dashboard

## Post-Deployment Verification

1. **Open the dashboard:**
   ```
   https://app.offensivewizard.com/dashboard
   ```

2. **Open browser developer console (F12) and check:**
   - [ ] No CORS errors in console
   - [ ] Stats cards load successfully
   - [ ] Pathway progress displays
   - [ ] Weekly activity chart shows
   - [ ] Recent activities load

3. **Check the Network tab:**
   - [ ] Requests to `/api/backend/v1/*` return 200 OK
   - [ ] Requests to `/api/judge0/*` return 200 OK (if testing code execution)
   - [ ] No failed requests to `https://api.offensivewizard.com` (should all go through proxy)

4. **Check server logs for proxy activity:**
   ```bash
   docker-compose logs -f frontend | grep "Proxy"
   ```
   
   You should see:
   ```
   🔄 Backend Proxy:
     Method: GET
     Path: v1/users/me/stats
     Target: http://backend:8080/v1/users/me/stats
   ✅ Response: 200 OK
   ```

## Troubleshooting

### If you see "Cannot connect to backend API"

1. **Check if backend is running:**
   ```bash
   docker-compose ps backend
   ```

2. **Check backend logs:**
   ```bash
   docker-compose logs backend
   ```

3. **Verify network connectivity:**
   ```bash
   docker-compose exec frontend ping backend
   ```

### If you see "Cannot connect to Judge0 server"

1. **Check if Judge0 is running:**
   ```bash
   docker-compose ps judge0
   ```

2. **Check Judge0 logs:**
   ```bash
   docker-compose logs judge0
   ```

### If dashboard shows old CORS errors

1. **Hard refresh the browser:**
   - Chrome/Edge: `Ctrl+Shift+R` (Windows) or `Cmd+Shift+R` (Mac)
   - Firefox: `Ctrl+F5` (Windows) or `Cmd+Shift+R` (Mac)

2. **Clear browser cache and reload**

3. **Verify the new image was pulled:**
   ```bash
   docker-compose exec frontend cat /app/app/api/backend/[...path]/route.ts
   ```
   
   If this file doesn't exist, the old image is still running

## Rollback Plan

If something goes wrong and you need to rollback:

1. **Stop the frontend:**
   ```bash
   docker-compose stop frontend
   ```

2. **Remove the new environment variables from `.env`:**
   ```bash
   # Remove these lines:
   # BACKEND_URL=http://backend:8080
   # JUDGE0_URL=http://judge0:2358
   ```

3. **Pull the previous image version:**
   ```bash
   docker pull limpet/wizardcore-frontend:previous-tag
   ```

4. **Update docker-compose.yml to use the previous tag:**
   ```yaml
   frontend:
     image: limpet/wizardcore-frontend:previous-tag
   ```

5. **Restart:**
   ```bash
   docker-compose up -d frontend
   ```

## Success Criteria

✅ **Deployment is successful when:**
- [ ] No CORS errors in browser console
- [ ] Dashboard loads with all stats, charts, and activities
- [ ] Backend proxy logs show successful requests
- [ ] Judge0 code execution works (if tested)
- [ ] No connectivity errors in server logs

## Documentation

- Full details: [`CORS-FIX-COMPLETE.md`](CORS-FIX-COMPLETE.md)
- Environment setup: [`.env.example`](.env.example)
- Docker configuration: [`docker-compose.yml`](docker-compose.yml)