#!/bin/bash

# ============================================
# WIZARDCORE ENVIRONMENT VERIFICATION SCRIPT
# ============================================
# Verifies all required environment variables are set correctly

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Load .env file
if [ ! -f .env ]; then
    echo -e "${RED}❌ Error: .env file not found${NC}"
    exit 1
fi

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   WIZARDCORE ENV VERIFICATION${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Source the .env file
set -a
source .env
set +a

# Counter for issues
ERRORS=0
WARNINGS=0

# Function to check if variable is set
check_required() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo -e "${RED}❌ MISSING: $var_name${NC}"
        ((ERRORS++))
        return 1
    else
        echo -e "${GREEN}✅ SET: $var_name${NC}"
        return 0
    fi
}

# Function to check if variable is set (warning only)
check_optional() {
    local var_name=$1
    local var_value="${!var_name}"
    
    if [ -z "$var_value" ]; then
        echo -e "${YELLOW}⚠️  OPTIONAL: $var_name (not set)${NC}"
        ((WARNINGS++))
        return 1
    else
        echo -e "${GREEN}✅ SET: $var_name${NC}"
        return 0
    fi
}

# Function to verify two variables match
check_match() {
    local var1_name=$1
    local var2_name=$2
    local var1_value="${!var1_name}"
    local var2_value="${!var2_name}"
    
    if [ "$var1_value" != "$var2_value" ]; then
        echo -e "${RED}❌ MISMATCH: $var1_name and $var2_name must be identical!${NC}"
        echo -e "   $var1_name: ${var1_value:0:50}..."
        echo -e "   $var2_name: ${var2_value:0:50}..."
        ((ERRORS++))
        return 1
    else
        echo -e "${GREEN}✅ MATCH: $var1_name == $var2_name${NC}"
        return 0
    fi
}

echo -e "${BLUE}Checking Critical Variables...${NC}"
echo ""

# Critical variables that MUST be set
check_required "GOTRUE_JWT_SECRET"
check_required "SUPABASE_JWT_SECRET"
check_required "DATABASE_URL"
check_required "REDIS_URL"
check_required "JUDGE0_API_KEY"
check_required "GOTRUE_DB_DATABASE_URL"

echo ""
echo -e "${BLUE}Checking JWT Secret Consistency...${NC}"
echo ""

# Verify JWT secrets match
check_match "GOTRUE_JWT_SECRET" "SUPABASE_JWT_SECRET"

echo ""
echo -e "${BLUE}Checking Supabase Configuration...${NC}"
echo ""

check_required "NEXT_PUBLIC_SUPABASE_URL"
check_required "NEXT_PUBLIC_SUPABASE_ANON_KEY"
check_required "SUPABASE_SERVICE_ROLE_KEY"
check_required "SUPABASE_POSTGRES_USER"
check_required "SUPABASE_POSTGRES_PASSWORD"
check_required "SUPABASE_POSTGRES_DB"

echo ""
echo -e "${BLUE}Checking Database Configuration...${NC}"
echo ""

check_required "POSTGRES_USER"
check_required "POSTGRES_PASSWORD"
check_required "POSTGRES_DB"

echo ""
echo -e "${BLUE}Checking Redis Configuration...${NC}"
echo ""

check_required "REDIS_PASSWORD"

echo ""
echo -e "${BLUE}Checking Judge0 Configuration...${NC}"
echo ""

check_required "JUDGE0_API_URL"
check_required "JUDGE0_POSTGRES_USER"
check_required "JUDGE0_POSTGRES_PASSWORD"
check_required "JUDGE0_POSTGRES_DB"

echo ""
echo -e "${BLUE}Checking Frontend Configuration...${NC}"
echo ""

check_required "NEXT_PUBLIC_BACKEND_URL"
check_required "NEXT_PUBLIC_JUDGE0_API_URL"
check_required "NEXT_PUBLIC_SITE_URL"

echo ""
echo -e "${BLUE}Checking Backend Configuration...${NC}"
echo ""

check_required "PORT"
check_required "ENVIRONMENT"
check_optional "LOG_LEVEL"
check_optional "CORS_ALLOWED_ORIGINS"

echo ""
echo -e "${BLUE}Checking GoTrue Configuration...${NC}"
echo ""

check_required "GOTRUE_SITE_URL"
check_required "API_EXTERNAL_URL"
check_required "GOTRUE_CORS_ALLOWED_ORIGINS"
check_optional "GOTRUE_DB_AUTOMIGRATE"
check_optional "GOTRUE_API_HOST"
check_optional "GOTRUE_API_PORT"

echo ""
echo -e "${BLUE}Checking Domain Configuration...${NC}"
echo ""

check_required "DOMAIN"
check_optional "FRONTEND_DOMAIN"
check_optional "BACKEND_DOMAIN"
check_optional "AUTH_DOMAIN"
check_optional "JUDGE0_DOMAIN"

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   VERIFICATION SUMMARY${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}✅ All critical variables are set correctly!${NC}"
else
    echo -e "${RED}❌ Found $ERRORS critical error(s)${NC}"
fi

if [ $WARNINGS -gt 0 ]; then
    echo -e "${YELLOW}⚠️  Found $WARNINGS optional variable(s) not set${NC}"
fi

echo ""

# Additional validation
echo -e "${BLUE}Additional Validations...${NC}"
echo ""

# Check DATABASE_URL format
if [[ "$DATABASE_URL" =~ ^postgresql:// ]]; then
    echo -e "${GREEN}✅ DATABASE_URL format is valid (PostgreSQL)${NC}"
else
    echo -e "${RED}❌ DATABASE_URL format is invalid (should start with postgresql://)${NC}"
    ((ERRORS++))
fi

# Check REDIS_URL format
if [[ "$REDIS_URL" =~ ^redis:// ]]; then
    echo -e "${GREEN}✅ REDIS_URL format is valid (Redis)${NC}"
else
    echo -e "${RED}❌ REDIS_URL format is invalid (should start with redis://)${NC}"
    ((ERRORS++))
fi

# Check if GOTRUE_DB_DATABASE_URL is a valid PostgreSQL connection string
if [[ "$GOTRUE_DB_DATABASE_URL" =~ ^postgresql:// ]]; then
    echo -e "${GREEN}✅ GOTRUE_DB_DATABASE_URL format is valid (PostgreSQL)${NC}"
else
    echo -e "${RED}❌ GOTRUE_DB_DATABASE_URL format is invalid (should start with postgresql://)${NC}"
    ((ERRORS++))
fi

# Check JWT secret length (should be at least 32 characters)
if [ ${#GOTRUE_JWT_SECRET} -ge 32 ]; then
    echo -e "${GREEN}✅ GOTRUE_JWT_SECRET length is sufficient (${#GOTRUE_JWT_SECRET} chars)${NC}"
else
    echo -e "${RED}❌ GOTRUE_JWT_SECRET is too short (${#GOTRUE_JWT_SECRET} chars, should be >= 32)${NC}"
    ((ERRORS++))
fi

# Check ANON key format (should be a JWT)
if [[ "$NEXT_PUBLIC_SUPABASE_ANON_KEY" =~ ^eyJ ]]; then
    echo -e "${GREEN}✅ NEXT_PUBLIC_SUPABASE_ANON_KEY appears to be a valid JWT${NC}"
else
    echo -e "${RED}❌ NEXT_PUBLIC_SUPABASE_ANON_KEY does not appear to be a valid JWT${NC}"
    ((ERRORS++))
fi

# Check SERVICE_ROLE key format (should be a JWT)
if [[ "$SUPABASE_SERVICE_ROLE_KEY" =~ ^eyJ ]]; then
    echo -e "${GREEN}✅ SUPABASE_SERVICE_ROLE_KEY appears to be a valid JWT${NC}"
else
    echo -e "${RED}❌ SUPABASE_SERVICE_ROLE_KEY does not appear to be a valid JWT${NC}"
    ((ERRORS++))
fi

echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   FINAL RESULT${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

if [ $ERRORS -eq 0 ]; then
    echo -e "${GREEN}🎉 SUCCESS! Your environment is properly configured.${NC}"
    echo -e "${GREEN}You can now deploy your application:${NC}"
    echo -e "${YELLOW}   docker compose up -d --build${NC}"
    echo ""
    exit 0
else
    echo -e "${RED}❌ FAILED! Found $ERRORS error(s) in your configuration.${NC}"
    echo -e "${RED}Please fix the errors above before deploying.${NC}"
    echo ""
    exit 1
fi
