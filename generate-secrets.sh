#!/bin/bash
# Generate secure passwords and secrets for WizardCore
# Run this to generate a secure .env file with proper passwords

set -e

echo "🔐 Generating Secure Passwords for WizardCore"
echo "=============================================="

# Function to generate random password
generate_password() {
    length=${1:-32}
    openssl rand -base64 $length | tr -d '/+=' | cut -c1-$length
}

# Function to generate JWT secret
generate_jwt_secret() {
    openssl rand -hex 64
}

echo "📝 Generating secure passwords..."

# Generate secure passwords
DB_PASSWORD=$(generate_password 32)
REDIS_PASSWORD=$(generate_password 32)
AUTH_DB_PASSWORD=$(generate_password 32)
JUDGE0_DB_PASSWORD=$(generate_password 32)
JWT_SECRET=$(generate_jwt_secret)
NEXTAUTH_SECRET=$(generate_password 64)

echo "✅ Passwords generated successfully!"
echo ""

# Create complete .env file with secure passwords
cat > .env << ENV_EOF
# ============================================
# WIZARDCORE SECURE ENVIRONMENT CONFIGURATION
# Generated on: $(date)
# ============================================

# DOMAIN CONFIGURATION
DOMAIN=offensivewizard.com
FRONTEND_DOMAIN=app.offensivewizard.com
BACKEND_DOMAIN=api.offensivewizard.com
AUTH_DOMAIN=auth.offensivewizard.com
JUDGE0_DOMAIN=judge0.offensivewizard.com

# FRONTEND CONFIGURATION
NEXT_PUBLIC_SUPABASE_URL=https://auth.offensivewizard.com
NEXT_PUBLIC_BACKEND_URL=https://api.offensivewizard.com
NEXT_PUBLIC_JUDGE0_API_URL=https://judge0.offensivewizard.com
NEXT_PUBLIC_SITE_URL=https://app.offensivewizard.com
NEXTAUTH_SECRET=${NEXTAUTH_SECRET}
NEXTAUTH_URL=https://app.offensivewizard.com

# SECRETS - GENERATED SECURELY
DATABASE_PASSWORD=${DB_PASSWORD}
SUPABASE_JWT_SECRET=${JWT_SECRET}
REDIS_PASSWORD=${REDIS_PASSWORD}

# BACKEND CONFIGURATION
CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
PORT=8080
ENVIRONMENT=production
LOG_LEVEL=info

# SUPABASE AUTH CONFIGURATION
GOTRUE_SITE_URL=https://app.offensivewizard.com
API_EXTERNAL_URL=https://auth.offensivewizard.com
GOTRUE_CORS_ALLOWED_ORIGINS=https://app.offensivewizard.com
GOTRUE_DB_AUTOMIGRATE=true
GOTRUE_JWT_SECRET=${JWT_SECRET}
GOTRUE_JWT_EXP=3600
GOTRUE_JWT_AUD=authenticated
GOTRUE_JWT_ISSUER=supabase

# DATABASE PASSWORDS
POSTGRES_PASSWORD=${DB_PASSWORD}
POSTGRES_USER=wizardcore
POSTGRES_DB=wizardcore

SUPABASE_DB_PASSWORD=${AUTH_DB_PASSWORD}
SUPABASE_DB_USER=supabase_auth_admin
SUPABASE_DB_NAME=supabase_auth

JUDGE0_POSTGRES_PASSWORD=${JUDGE0_DB_PASSWORD}
JUDGE0_POSTGRES_USER=judge0
JUDGE0_POSTGRES_DB=judge0

# REDIS PASSWORDS
JUDGE0_REDIS_PASSWORD=$(generate_password 32)

# SSL CERTIFICATE EMAIL
TRAEFIK_ACME_EMAIL=admin@offensivewizard.com

# ============================================
# GENERATED SECRETS SUMMARY
# ============================================
# Database Password: ${DB_PASSWORD}
# Redis Password: ${REDIS_PASSWORD}
# Auth DB Password: ${AUTH_DB_PASSWORD}
# Judge0 DB Password: ${JUDGE0_DB_PASSWORD}
# JWT Secret: ${JWT_SECRET}
# NextAuth Secret: ${NEXTAUTH_SECRET}
ENV_EOF

echo "📄 Created .env file with secure passwords"
echo ""
echo "🔐 Generated Secrets Summary:"
echo "============================="
echo "Database Password:     ${DB_PASSWORD}"
echo "Redis Password:        ${REDIS_PASSWORD}"
echo "Auth DB Password:      ${AUTH_DB_PASSWORD}"
echo "Judge0 DB Password:    ${JUDGE0_DB_PASSWORD}"
echo "JWT Secret:            ${JWT_SECRET}"
echo "NextAuth Secret:       ${NEXTAUTH_SECRET}"
echo ""
echo "📋 Next steps:"
echo "1. Review the generated .env file"
echo "2. Push to your server:"
echo "   scp .env root@your-server-ip:/opt/wizardcore/"
echo "3. Keep this information secure!"
echo ""
echo "⚠️  IMPORTANT: Save these passwords in a secure location!"
echo "   They cannot be recovered if lost."