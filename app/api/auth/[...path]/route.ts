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

// Handle OPTIONS preflight requests for CORS
export async function OPTIONS(request: NextRequest) {
  const origin = request.headers.get('origin') || '*'
  
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': origin,
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization, Content-Type, X-Client-Info, apikey, x-client-info, x-supabase-api-version',
      'Access-Control-Max-Age': '86400',
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

    // Log request details for debugging
    console.log('🔄 GoTrue Proxy:')
    console.log('  Method:', request.method)
    console.log('  Original:', path.join('/'))
    console.log('  Stripped:', targetPath)
    console.log('  Target:', targetUrl)
    console.log('  GOTRUE_URL:', GOTRUE_URL)

    // Copy headers from incoming request
    const headers = new Headers()
    request.headers.forEach((value, key) => {
      // Skip host header as we're proxying to a different host
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value)
      }
    })

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
    const timeoutId = setTimeout(() => controller.abort(), 30000) // 30 second timeout

    const response = await fetch(targetUrl, {
      method: request.method,
      headers: headers,
      body: body,
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    // Get response body
    const responseData = await response.text()
    
    // Log response for debugging
    if (process.env.NODE_ENV === 'development') {
      console.log('✅ Response:', response.status, response.statusText)
    }

    // Get origin for CORS
    const origin = request.headers.get('origin') || '*'

    // Return response with proper CORS headers
    const responseHeaders = new Headers()
    responseHeaders.set('Content-Type', response.headers.get('Content-Type') || 'application/json')
    responseHeaders.set('Access-Control-Allow-Origin', origin)
    responseHeaders.set('Access-Control-Allow-Credentials', 'true')
    responseHeaders.set('Access-Control-Expose-Headers', 'X-Total-Count')
    
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

    // Return error response
    return new NextResponse(
      JSON.stringify({
        error: 'proxy_error',
        message: errorMessage,
        details: {
          gotrue_url: GOTRUE_URL,
          error_type: error instanceof Error ? error.constructor.name : typeof error,
          suggestion: 'If running in Docker, ensure SUPABASE_INTERNAL_URL is set to http://supabase-auth:9999'
        }
      }),
      {
        status: statusCode,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': request.headers.get('origin') || '*',
          'Access-Control-Allow-Credentials': 'true',
        },
      }
    )
  }
}
