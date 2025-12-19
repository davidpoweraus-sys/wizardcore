This is a [Next.js](https://nextjs.org) project bootstrapped with [`create-next-app`](https://nextjs.org/docs/app/api-reference/cli/create-next-app).

## Getting Started

First, run the development server:

```bash
npm run dev
# or
yarn dev
# or
pnpm dev
# or
bun dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

You can start editing the page by modifying `app/page.tsx`. The page auto-updates as you edit the file.

This project uses [`next/font`](https://nextjs.org/docs/app/building-your-application/optimizing/fonts) to automatically optimize and load [Geist](https://vercel.com/font), a new font family for Vercel.

## Learn More

To learn more about Next.js, take a look at the following resources:

- [Next.js Documentation](https://nextjs.org/docs) - learn about Next.js features and API.
- [Learn Next.js](https://nextjs.org/learn) - an interactive Next.js tutorial.

You can check out [the Next.js GitHub repository](https://github.com/vercel/next.js) - your feedback and contributions are welcome!

## Deploy on Vercel

The easiest way to deploy your Next.js app is to use the [Vercel Platform](https://vercel.com/new?utm_medium=default-template&filter=next.js&utm_source=create-next-app&utm_campaign=create-next-app-readme) from the creators of Next.js.

Check out our [Next.js deployment documentation](https://nextjs.org/docs/app/building-your-application/deploying) for more details.

## Supabase Auth Troubleshooting

If Supabase Auth fails to start with errors about missing "auth" schema, ensure the init container runs before the auth service.

The project includes an initialization script (`init-scripts/01-create-auth-schema.sql`) that creates the required schema. In both local and production Docker Compose configurations, a `supabase-init` service runs after the database is healthy and before `supabase-auth` starts.

### Common Issues

1. **Schema not created**: The init script may not have executed due to dependency cycles. Verify that `supabase-init` depends on `supabase-postgres` (condition: service_healthy) and `supabase-auth` depends on `supabase-init` (condition: service_started).

2. **Permissions**: The database user `supabase_auth_admin` must have privileges on the `auth` schema. The init script grants these.

3. **Production deployment**: When deploying with Coolify, ensure the environment variables `SUPABASE_JWT_SECRET` and `SUPABASE_URL` are set correctly.

### Fallback Mechanism

The backend includes a health check that can detect missing schema and attempt to create it programmatically (optional). However, the init container approach is preferred.

### Manual Fix

If the auth schema is missing, connect to the `supabase_auth` database and run:

```sql
CREATE SCHEMA IF NOT EXISTS auth;
GRANT ALL ON SCHEMA auth TO supabase_auth_admin;
```

### Debugging

Check logs of `supabase-auth` and `supabase-init` containers:

```bash
docker-compose -f docker-compose.local.yml logs supabase-auth
docker-compose -f docker-compose.local.yml logs supabase-init
```

For production, use Coolify's log viewer.
