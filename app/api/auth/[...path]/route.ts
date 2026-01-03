import { NextRequest, NextResponse } from 'next/server'

/**
 * Supabase Auth Proxy for GoTrue
 *
 * This proxy solves the path mismatch between:
 * - @supabase/supabase-js client (expects /auth/v1/* endpoints)
 * - Standalone GoTrue server (serves at /* root endpoints)
 *
 * The Supabase client will make requests to:
 *   https://yourdomain.com/api/auth/auth/v1/signup
 *
 * This proxy strips "/auth/v1" and forwards to:
 *   https://auth.offensivewizard.com/signup
 */

// GoTrue server URL - must point to your standalone GoTrue instance
// In Docker: use internal container name
// Outside Docker: use public URL
const GOTRUE_URL = process.env.GOTRUE_URL || process.env.SUPABASE_INTERNAL_URL || 'http://supabase-auth:9999'

// BLUE DIE TEST: Clear identifier for debugging login issues
// This will appear in logs and response headers to verify the correct image is running
const PROXY_VERSION = 'login-fix-v2-20260103-1851'

/**
 * Validate and normalize CORS origin
 *
 * Enterprise CORS handling principles:
 * 1. Validate against allowed origins list
 * 2. Normalize origin (strip path, port, etc.)
 * 3. Support wildcard patterns
 * 4. Environment-specific configurations
 *
 * This implementation allows:
 * - Any subdomain of offensivewizard.com (*.offensivewizard.com)
 * - Localhost with any port (localhost:*)
 * - Specific development ports
 */
function validateOrigin(origin: string | null): string | null {
  if (!origin || origin === '*') {
    // In production, avoid wildcard when using credentials
    if (process.env.NODE_ENV === 'production') {
      return null
    }
    return '*'
  }

  try {
    const url = new URL(origin)
    const hostname = url.hostname
    
    // Enterprise pattern: Allow any offensivewizard.com subdomain
    if (hostname === 'offensivewizard.com' ||
        hostname === 'www.offensivewizard.com' ||
        hostname.endsWith('.offensivewizard.com')) {
      // Return the exact origin (browser will send correct origin without path)
      return origin
    }
    
    // Development: Allow localhost with any port
    if (hostname === 'localhost') {
      return origin
    }
    
    // Check specific allowed origins from environment
    const allowedOrigins = process.env.ALLOWED_ORIGINS?.split(',') || []
    if (allowedOrigins.includes(origin)) {
      return origin
    }
    
    return null
  } catch (error) {
    console.warn('Invalid origin format:', origin)
    return null
  }
}

// Handle OPTIONS preflight requests for CORS
export async function OPTIONS(request: NextRequest) {
  const origin = request.headers.get('origin')
  const validatedOrigin = validateOrigin(origin)
  
  // If no valid origin, return 403 for preflight
  if (!validatedOrigin) {
    return new NextResponse(null, {
      status: 403,
      headers: {
        'Content-Type': 'application/json',
      },
    })
  }
  
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': validatedOrigin,
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization, Content-Type, X-Client-Info, apikey, x-client-info, x-supabase-api-version',
      'Access-Control-Max-Age': '86400',
      // BLUE DIE TEST: Add version header to preflight
      'X-Auth-Proxy-Version': PROXY_VERSION,
      'X-Login-Fix': 'active',
    },
  })
}

// Handle all HTTP methods
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

/**
 * Main proxy function that forwards requests to GoTrue
 */
async function proxyRequest(request: NextRequest, path: string[]) {
  try {
    // BLUE DIE TEST: Log version identifier
    console.log('🎲 BLUE DIE TEST - Auth Proxy Version:', PROXY_VERSION)
    console.log('🎲 This log confirms the login-fix image is running')
    
    const url = new URL(request.url)
    let targetPath = path.join('/')
    
    // CRITICAL: Strip /auth/v1 prefix that Supabase client adds
    // GoTrue serves endpoints at root level, not under /auth/v1/
    //
    // Examples:
    //   /auth/v1/signup    → /signup
    //   /auth/v1/token     → /token
    //   /auth/v1/user      → /user
    //   /auth/v1/logout    → /logout
    targetPath = targetPath.replace(/^auth\/v1\//, '')
    
    // Build the full target URL
    const baseUrl = GOTRUE_URL.endsWith('/') ? GOTRUE_URL.slice(0, -1) : GOTRUE_URL
    const targetUrl = `${baseUrl}/${targetPath}${url.search}`

    // Log request details for debugging - ALWAYS log in production for auth issues
    console.log('🔄 GoTrue Proxy:')
    console.log('  Method:', request.method)
    console.log('  Original:', path.join('/'))
    console.log('  Stripped:', targetPath)
    console.log('  Target:', targetUrl)
    console.log('  GOTRUE_URL:', GOTRUE_URL)
    console.log('  Query:', url.search)
    console.log('  Version:', PROXY_VERSION)
    
    // Log headers for debugging (redact sensitive values)
    const authHeader = request.headers.get('authorization') || request.headers.get('Authorization')
    if (authHeader) {
      console.log('  Has Authorization header:', authHeader.substring(0, 20) + '...')
    }
    console.log('  Has apikey header:', !!request.headers.get('apikey'))

    // Copy headers from incoming request
    const headers = new Headers()
    request.headers.forEach((value, key) => {
      // Skip host header as we're proxying to a different host
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value)
      }
    })

    // CRITICAL FIX: For password grant login requests, remove Authorization header and auth cookies
    // The Supabase client incorrectly adds Authorization header with anon key JWT
    // for login requests, which can cause GoTrue to hang or reject the request
    // Also remove auth cookies to prevent session conflicts during login
    if (targetPath === 'token' && url.search.includes('grant_type=password')) {
      console.log('🔐 Removing Authorization header and auth cookies for password grant login')
      headers.delete('authorization')
      headers.delete('Authorization')
      
      // Also remove auth-related cookies to prevent conflicts
      const cookieHeader = headers.get('cookie')
      if (cookieHeader) {
        // Remove sb-* cookies (Supabase auth cookies)
        const cleanedCookies = cookieHeader.split(';').filter(cookie => {
          const cookieName = cookie.trim().split('=')[0]
          return !cookieName.startsWith('sb-')
        }).join('; ')
        
        if (cleanedCookies !== cookieHeader) {
          console.log('  Removed auth cookies')
          if (cleanedCookies) {
            headers.set('cookie', cleanedCookies)
          } else {
            headers.delete('cookie')
          }
        }
      }
    }

    // Read request body if present
    let body = null
    if (request.method !== 'GET' && request.method !== 'HEAD') {
      try {
        body = await request.text()
      } catch (e) {
        // Body may not exist or already consumed
      }
    }

    // Make request to GoTrue with timeout
    const controller = new AbortController()
    const timeoutId = setTimeout(() => {
      console.error('⏰ Proxy timeout after 30 seconds for:', targetUrl)
      console.error('  Method:', request.method)
      console.error('  Headers:', Object.fromEntries(headers.entries()))
      controller.abort()
    }, 30000) // 30 second timeout

    const startTime = Date.now()
    const response = await fetch(targetUrl, {
      method: request.method,
      headers: headers,
      body: body,
      signal: controller.signal,
    })
    const endTime = Date.now()
    const duration = endTime - startTime

    clearTimeout(timeoutId)
    
    // Log slow requests
    if (duration > 5000) {
      console.warn(`🐌 Slow request: ${duration}ms for ${targetUrl}`)
    }

    // Get response body
    const responseData = await response.text()
    
    // Log response for debugging - ALWAYS log in production for auth issues
    console.log('✅ Response:', response.status, response.statusText)
    if (response.status >= 400) {
      console.log('  Error response:', responseData.substring(0, 200))
    }

    // Get and validate origin for CORS
    const origin = request.headers.get('origin')
    const validatedOrigin = validateOrigin(origin)

    // If no valid origin, return 403
    if (!validatedOrigin) {
      return new NextResponse(
        JSON.stringify({
          error: 'cors_error',
          message: 'Origin not allowed',
          details: {
            origin,
            allowed_origins: 'localhost:*, *.offensivewizard.com, offensivewizard.com'
          }
        }),
        {
          status: 403,
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )
    }

    // Return response with proper CORS headers
    const responseHeaders = new Headers()
    responseHeaders.set('Content-Type', response.headers.get('Content-Type') || 'application/json')
    responseHeaders.set('Access-Control-Allow-Origin', validatedOrigin)
    responseHeaders.set('Access-Control-Allow-Credentials', 'true')
    responseHeaders.set('Access-Control-Expose-Headers', 'X-Total-Count')
    
    // BLUE DIE TEST: Add version header for browser inspection
    responseHeaders.set('X-Auth-Proxy-Version', PROXY_VERSION)
    responseHeaders.set('X-Login-Fix', 'active')
    
    // Forward ALL Set-Cookie headers from GoTrue (multiple cookies)
    const setCookieHeaders = response.headers.getSetCookie()
    for (const cookie of setCookieHeaders) {
      responseHeaders.append('Set-Cookie', cookie)
    }
    
    // Also forward other important headers
    const headersToForward = ['X-Total-Count', 'Cache-Control', 'ETag']
    for (const header of headersToForward) {
      const value = response.headers.get(header)
      if (value) {
        responseHeaders.set(header, value)
      }
    }

    return new NextResponse(responseData, {
      status: response.status,
      statusText: response.statusText,
      headers: responseHeaders,
    })
  } catch (error) {
    console.error('❌ Proxy error:')
    console.error('  Error:', error instanceof Error ? error.message : String(error))
    console.error('  Target:', GOTRUE_URL)
    console.error('  Error type:', error instanceof Error ? error.constructor.name : typeof(error))
    
    // Check for common connection errors
    let errorMessage = error instanceof Error ? error.message : 'Unknown proxy error'
    let statusCode = 502
    
    if (error instanceof Error) {
      if (error.message.includes('ECONNREFUSED') || error.message.includes('ENOTFOUND')) {
        errorMessage = `Cannot connect to GoTrue server at ${GOTRUE_URL}. Check if the service is running and accessible.`
        console.error('  ⚠️ Network error - service may not be running or DNS not resolving')
      } else if (error.message.includes('fetch failed')) {
        errorMessage = `Failed to connect to GoTrue server at ${GOTRUE_URL}. This could be a network, DNS, or SSL issue.`
        console.error('  ⚠️ Fetch failed - check network connectivity and SSL certificates')
      }
    }

    // Get and validate origin for error response
    const errorOrigin = request.headers.get('origin')
    const validatedErrorOrigin = validateOrigin(errorOrigin)
    
    // Return error response
    const errorHeaders = new Headers()
    errorHeaders.set('Content-Type', 'application/json')
    
    // BLUE DIE TEST: Add version header even for errors
    errorHeaders.set('X-Auth-Proxy-Version', PROXY_VERSION)
    errorHeaders.set('X-Login-Fix', 'active')
    
    if (validatedErrorOrigin) {
      errorHeaders.set('Access-Control-Allow-Origin', validatedErrorOrigin)
      errorHeaders.set('Access-Control-Allow-Credentials', 'true')
    }
    
    return new NextResponse(
      JSON.stringify({
        error: 'proxy_error',
        message: errorMessage,
        details: {
          gotrue_url: GOTRUE_URL,
          error_type: error instanceof Error ? error.constructor.name : typeof error,
          suggestion: 'If running in Docker, ensure SUPABASE_INTERNAL_URL is set to http://supabase-auth:9999',
          proxy_version: PROXY_VERSION
        }
      }),
      {
        status: statusCode,
        headers: errorHeaders,
      }
    )
  }
}
