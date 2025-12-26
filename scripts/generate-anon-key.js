#!/usr/bin/env node

/**
 * Generate a Supabase ANON key from JWT secret
 * 
 * The ANON key is a JWT token with specific claims that allows
 * public access to GoTrue endpoints. It's signed with your JWT secret.
 */

const crypto = require('crypto');

// Read JWT secret from environment or command line
const JWT_SECRET = process.env.SUPABASE_JWT_SECRET || process.argv[2];

if (!JWT_SECRET) {
  console.error('❌ Error: JWT_SECRET is required');
  console.error('');
  console.error('Usage:');
  console.error('  SUPABASE_JWT_SECRET=your-secret node scripts/generate-anon-key.js');
  console.error('  OR');
  console.error('  node scripts/generate-anon-key.js your-secret');
  process.exit(1);
}

/**
 * Base64URL encode
 */
function base64urlEncode(str) {
  return Buffer.from(str)
    .toString('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');
}

/**
 * Create JWT token
 */
function createJWT(secret, payload) {
  // JWT Header
  const header = {
    alg: 'HS256',
    typ: 'JWT'
  };

  // Encode header and payload
  const encodedHeader = base64urlEncode(JSON.stringify(header));
  const encodedPayload = base64urlEncode(JSON.stringify(payload));

  // Create signature
  const signatureInput = `${encodedHeader}.${encodedPayload}`;
  const signature = crypto
    .createHmac('sha256', secret)
    .update(signatureInput)
    .digest('base64')
    .replace(/\+/g, '-')
    .replace(/\//g, '_')
    .replace(/=/g, '');

  // Combine into JWT
  return `${encodedHeader}.${encodedPayload}.${signature}`;
}

/**
 * Generate ANON key
 */
function generateAnonKey(secret) {
  const now = Math.floor(Date.now() / 1000);
  
  // ANON key payload - standard claims for public access
  const payload = {
    role: 'anon',
    iss: 'supabase',
    iat: now,
    exp: now + (10 * 365 * 24 * 60 * 60), // 10 years from now
  };

  return createJWT(secret, payload);
}

/**
 * Generate SERVICE_ROLE key (for admin operations)
 */
function generateServiceRoleKey(secret) {
  const now = Math.floor(Date.now() / 1000);
  
  const payload = {
    role: 'service_role',
    iss: 'supabase',
    iat: now,
    exp: now + (10 * 365 * 24 * 60 * 60), // 10 years from now
  };

  return createJWT(secret, payload);
}

// Generate keys
console.log('🔑 Generating Supabase API Keys...\n');

const anonKey = generateAnonKey(JWT_SECRET);
const serviceRoleKey = generateServiceRoleKey(JWT_SECRET);

console.log('✅ Keys generated successfully!\n');
console.log('Add these to your .env file:\n');
console.log('# Public ANON key (safe to expose in frontend)');
console.log(`NEXT_PUBLIC_SUPABASE_ANON_KEY=${anonKey}`);
console.log('');
console.log('# Service role key (KEEP SECRET - server-side only)');
console.log(`SUPABASE_SERVICE_ROLE_KEY=${serviceRoleKey}`);
console.log('');

// Verify the keys
console.log('📋 Key Information:');
console.log('  ANON Key length:', anonKey.length);
console.log('  Service Role Key length:', serviceRoleKey.length);
console.log('');
console.log('⚠️  IMPORTANT:');
console.log('  - ANON key is for public/client-side use');
console.log('  - Service Role key must be kept SECRET (server-side only)');
console.log('  - Both keys are valid for 10 years from now');
