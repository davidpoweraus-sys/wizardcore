# 🚀 Quick Start - Wizardcore Deployment

## TL;DR

1. **Deploy to Dokploy** (all scripts run automatically)
2. **Create admin in Supabase Auth**:
   - Email: `admin@offensivewizard.com`
   - UUID: `00000000-0000-0000-0000-000000000001`
   - Password: (your choice)
3. **Login** at `https://app.offensivewizard.com/login`

---

## What Happens Automatically

✅ Supabase Auth database initializes  
✅ Main database initializes  
✅ RBAC system deploys (30+ permissions, 5 roles)  
✅ Default admin user created in database  
✅ All services start and become healthy  

---

## What You Must Do Manually

⚠️ **Create admin user in Supabase Auth**

See `ADMIN-USER-SETUP.md` for detailed instructions.

---

## Files You Created

```
├── init-scripts/
│   ├── 00-create-user.sql              # Supabase user
│   └── 01-create-auth-schema.sql       # Supabase schema
│
├── postgres-init-scripts/
│   └── 01-create-default-admin.sql     # Admin user
│
├── wizardcore-backend/.../migrations/
│   ├── 013_rbac_system.up.sql          # RBAC tables
│   └── 013_rbac_system.down.sql        # Rollback
│
└── docker-compose.yml                   # Updated with postgres-init
```

---

## Default Credentials

**Admin Email**: `admin@offensivewizard.com`  
**Password**: Set in Supabase Auth (you choose)  
**Roles**: `super_admin`, `admin`  
**Permissions**: ALL  

⚠️ **Change these in production!**

---

## Troubleshooting

**Can't login?** → Create Supabase Auth user (see `ADMIN-USER-SETUP.md`)  
**Auth errors?** → Run `./fix-dokploy-volumes.sh`  
**Check status**: Run `./diagnose-deployment.sh`  

---

## Full Documentation

- `DEPLOYMENT-SUMMARY.md` - Complete overview
- `ADMIN-USER-SETUP.md` - Admin setup guide
- `DOKPLOY-FIX-GUIDE.md` - Troubleshooting
- `DATABASE-INITIALIZATION-GUIDE.md` - How it works
