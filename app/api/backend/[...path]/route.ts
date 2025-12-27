import { NextRequest, NextResponse } from 'next/server'
import { createClient } from '@/lib/supabase/server'

/**
 * Backend API Proxy
 * 
 * This proxy forwards requests to the backend API server and handles CORS.
 * It also adds authentication tokens from the session.
 */

const BACKEND_URL = process.env.BACKEND_URL || process.env.NEXT_PUBLIC_BACKEND_URL || 'https://api.offensivewizard.com'

// Handle OPTIONS preflight requests for CORS
export async function OPTIONS(request: NextRequest) {
  const origin = request.headers.get('origin') || '*'
  
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': origin,
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
    
    // Build the full target URL
    const baseUrl = BACKEND_URL.endsWith('/') ? BACKEND_URL.slice(0, -1) : BACKEND_URL
    const targetUrl = `${baseUrl}/${targetPath}${url.search}`

    // Log request details for debugging
    console.log('🔄 Backend Proxy:')
    console.log('  Method:', request.method)
    console.log('  Path:', targetPath)
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

    // Get origin for CORS
    const origin = request.headers.get('origin') || '*'

    // Return response with proper CORS headers
    return new NextResponse(responseData, {
      status: response.status,
      statusText: response.statusText,
      headers: {
        'Content-Type': response.headers.get('Content-Type') || 'application/json',
        'Access-Control-Allow-Origin': origin,
        'Access-Control-Allow-Credentials': 'true',
        'Access-Control-Expose-Headers': 'X-Total-Count',
      },
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

    // Return error response
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
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': request.headers.get('origin') || '*',
          'Access-Control-Allow-Credentials': 'true',
        },
      }
    )
  }
}