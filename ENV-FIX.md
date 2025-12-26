# Quick Fix for SSL Error

## Problem
Backend is failing with: `pq: SSL is not enabled on the server`

## Solution
Update your `.env` file in Dokploy:

### Change this line:
```
DATABASE_URL=postgresql://wizardcore:YKiDQeXatsFIVvILstuZrxYBCPUdBlFC@postgres:5432/wizardcore
```

### To this (add `?sslmode=disable`):
```
DATABASE_URL=postgresql://wizardcore:YKiDQeXatsFIVvILstuZrxYBCPUdBlFC@postgres:5432/wizardcore?sslmode=disable
```

## Why?
Your PostgreSQL container doesn't have SSL certificates configured. In production with internal Docker networks, SSL between containers isn't necessary since traffic doesn't leave the server.

## Apply the Fix

1. Go to Dokploy dashboard
2. Find your application
3. Click "Environment Variables" or "Settings"
4. Update the `DATABASE_URL` variable
5. Redeploy

The backend should start successfully after this change.
