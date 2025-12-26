#!/bin/bash
# Setup environment file with production-ready defaults
# This script copies .env.example to .env if .env doesn't exist

set -e

ENV_FILE=".env"
EXAMPLE_FILE=".env.example"

echo "🔧 WizardCore Environment Setup"
echo "================================"
echo ""

# Check if .env already exists
if [ -f "$ENV_FILE" ]; then
    echo "⚠️  .env file already exists!"
    echo ""
    read -p "Do you want to overwrite it? (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "❌ Cancelled. Keeping existing .env file."
        exit 0
    fi
fi

# Check if .env.example exists
if [ ! -f "$EXAMPLE_FILE" ]; then
    echo "❌ Error: $EXAMPLE_FILE not found!"
    exit 1
fi

# Copy .env.example to .env
echo "📋 Copying $EXAMPLE_FILE to $ENV_FILE..."
cp "$EXAMPLE_FILE" "$ENV_FILE"

echo ""
echo "✅ Environment file created successfully!"
echo ""
echo "📝 The .env file now contains production-ready values."
echo ""
echo "⚠️  IMPORTANT: Review and customize these values:"
echo "   - Domain names (if not using offensivewizard.com)"
echo "   - Database passwords (CHANGE IN PRODUCTION!)"
echo "   - JWT secrets (CHANGE IN PRODUCTION!)"
echo "   - Email settings (if using SMTP)"
echo ""
echo "🔐 SECURITY NOTES:"
echo "   - NEXT_PUBLIC_SUPABASE_ANON_KEY is derived from SUPABASE_JWT_SECRET"
echo "   - If you change SUPABASE_JWT_SECRET, you MUST regenerate the ANON key:"
echo "     node scripts/generate-anon-key.js"
echo ""
echo "📚 For more information, see:"
echo "   - AUTH-PROXY-FIX.md"
echo "   - DEPLOY-AUTH-PROXY-FIX.md"
echo ""
