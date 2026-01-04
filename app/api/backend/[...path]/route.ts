import { NextRequest, NextResponse } from 'next/server'
import { createClient } from '@/lib/supabase/server'

/**
 * Backend API Proxy
 *
 * This proxy forwards requests to the backend API server and handles CORS.
 * It also adds authentication tokens from the session.
 */

const BACKEND_URL = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL || 'https://api.offensivewizard.com'

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
  // CRITICAL FIX: Allow null/empty origin for same-origin requests
  // Same-origin requests often don't include Origin header
  // We should allow these requests in production
  if (!origin) {
    // For same-origin requests (no Origin header), we need to determine
    // if this is a valid request. Since this is our own frontend making
    // requests to our own API, we should allow it.
    // Return a placeholder that indicates same-origin
    return 'same-origin'
  }
  
  if (origin === '*') {
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
  
  // Handle same-origin requests
  let corsOrigin = validatedOrigin
  if (validatedOrigin === 'same-origin') {
    // Determine actual origin from request
    const referer = request.headers.get('referer')
    if (referer) {
      try {
        const refererUrl = new URL(referer)
        corsOrigin = refererUrl.origin
      } catch {
        corsOrigin = 'https://app.offensivewizard.com'
      }
    } else {
      corsOrigin = 'https://app.offensivewizard.com'
    }
  }
  
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': corsOrigin,
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization, Content-Type, X-Client-Info, apikey',
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
 * Main proxy function that forwards requests to backend API
 */
async function proxyRequest(request: NextRequest, path: string[]) {
  try {
    const url = new URL(request.url)
    const targetPath = path.join('/')
    
    // Prepend /api to the path since backend routes are under /api/v1
    // Frontend calls /api/backend/v1/users -> backend expects /api/v1/users
    const backendPath = targetPath.startsWith('api/') ? targetPath : `api/${targetPath}`
    
    // Build the full target URL
    const baseUrl = BACKEND_URL.endsWith('/') ? BACKEND_URL.slice(0, -1) : BACKEND_URL
    const targetUrl = `${baseUrl}/${backendPath}${url.search}`

    // Log request details for debugging
    console.log('🔄 Backend Proxy:')
    console.log('  Method:', request.method)
    console.log('  Original Path:', targetPath)
    console.log('  Backend Path:', backendPath)
    console.log('  Target:', targetUrl)

    // Get auth token from session
    const supabase = await createClient()
    const { data: { session } } = await supabase.auth.getSession()

    // Copy headers from incoming request
    const headers = new Headers()
    request.headers.forEach((value, key) => {
      // Skip host header as we're proxying to a different host
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value)
      }
    })

    // Add auth token if available
    if (session?.access_token) {
      headers.set('Authorization', `Bearer ${session.access_token}`)
      console.log('  Auth: Added Bearer token')
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

    // Make request to backend with timeout
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
            allowed_origins: 'localhost:*, *.offensivewizard.com, offensivewizard.com, same-origin'
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

    // Handle same-origin requests
    let corsOrigin = validatedOrigin
    if (validatedOrigin === 'same-origin') {
      // Determine actual origin from request
      const host = request.headers.get('host')
      if (host && (host.includes('offensivewizard.com') || host.includes('localhost'))) {
        corsOrigin = `https://${host}`
      } else {
        corsOrigin = 'https://app.offensivewizard.com'
      }
    }

    // Return response with proper CORS headers
    const responseHeaders = new Headers()
    responseHeaders.set('Content-Type', response.headers.get('Content-Type') || 'application/json')
    responseHeaders.set('Access-Control-Allow-Origin', corsOrigin)
    responseHeaders.set('Access-Control-Allow-Credentials', 'true')
    responseHeaders.set('Access-Control-Expose-Headers', 'X-Total-Count')
    
    // Forward other important headers
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
    console.error('❌ Backend Proxy error:')
    console.error('  Error:', error instanceof Error ? error.message : String(error))
    console.error('  Target:', BACKEND_URL)
    
    // Check for common connection errors
    let errorMessage = error instanceof Error ? error.message : 'Unknown proxy error'
    let statusCode = 502
    
    if (error instanceof Error) {
      if (error.message.includes('ECONNREFUSED') || error.message.includes('ENOTFOUND')) {
        errorMessage = `Cannot connect to backend API at ${BACKEND_URL}. Check if the service is running and accessible.`
        console.error('  ⚠️ Network error - service may not be running or DNS not resolving')
      } else if (error.message.includes('fetch failed')) {
        errorMessage = `Failed to connect to backend API at ${BACKEND_URL}. This could be a network, DNS, or SSL issue.`
        console.error('  ⚠️ Fetch failed - check network connectivity and SSL certificates')
      } else if (error.message.includes('aborted')) {
        errorMessage = `Request to backend API timed out after 30 seconds.`
        statusCode = 504
      }
    }

    // Get and validate origin for error response
    const errorOrigin = request.headers.get('origin')
    const validatedErrorOrigin = validateOrigin(errorOrigin)
    
    // Return error response
    const errorHeaders = new Headers()
    errorHeaders.set('Content-Type', 'application/json')
    
    if (validatedErrorOrigin) {
      // Handle same-origin for error responses
      let errorCorsOrigin = validatedErrorOrigin
      if (validatedErrorOrigin === 'same-origin') {
        const host = request.headers.get('host')
        if (host && (host.includes('offensivewizard.com') || host.includes('localhost'))) {
          errorCorsOrigin = `https://${host}`
        } else {
          errorCorsOrigin = 'https://app.offensivewizard.com'
        }
      }
      errorHeaders.set('Access-Control-Allow-Origin', errorCorsOrigin)
      errorHeaders.set('Access-Control-Allow-Credentials', 'true')
    }
    
    return new NextResponse(
      JSON.stringify({
        error: 'proxy_error',
        message: errorMessage,
        details: {
          backend_url: BACKEND_URL,
          error_type: error instanceof Error ? error.constructor.name : typeof error,
          suggestion: 'Check if the backend API server is running and accessible from this container/network'
        }
      }),
      {
        status: statusCode,
        headers: errorHeaders,
      }
    )
  }
}