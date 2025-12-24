import { NextRequest, NextResponse } from 'next/server'

// Supabase Auth URL
// Using external URL to bypass Docker DNS issues in Coolify
// This still solves CORS (proxy adds headers from same origin)
const GOTRUE_URL = process.env.SUPABASE_INTERNAL_URL || 'https://auth.offensivewizard.com'

// Handle OPTIONS preflight requests
export async function OPTIONS(request: NextRequest) {
  return new NextResponse(null, {
    status: 204,
    headers: {
      'Access-Control-Allow-Origin': 'https://offensivewizard.com',
      'Access-Control-Allow-Credentials': 'true',
      'Access-Control-Allow-Methods': 'GET, POST, PUT, PATCH, DELETE, OPTIONS',
      'Access-Control-Allow-Headers': 'Authorization, Content-Type, X-Client-Info, X-Requested-With, apikey, x-client-info, x-supabase-api-version, Accept, Accept-Language, Content-Language',
      'Access-Control-Max-Age': '86400',
      'Content-Type': 'text/plain charset=UTF-8',
      'Content-Length': '0',
    },
  })
}

// Handle GET requests
export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

// Handle POST requests
export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

// Handle PUT requests
export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

// Handle PATCH requests
export async function PATCH(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

// Handle DELETE requests
export async function DELETE(
  request: NextRequest,
  { params }: { params: Promise<{ path: string[] }> }
) {
  const { path } = await params
  return proxyRequest(request, path)
}

async function proxyRequest(request: NextRequest, path: string[]) {
  try {
    const url = new URL(request.url)
    let targetPath = path.join('/')
    
    // Strip /auth/v1 prefix if present (Supabase client adds this automatically)
    targetPath = targetPath.replace(/^auth\/v1\//, '')
    
    const targetUrl = `${GOTRUE_URL}/${targetPath}${url.search}`

    console.log('🔄 Proxy Configuration:')
    console.log('  GOTRUE_URL:', GOTRUE_URL)
    console.log('  Original Path:', path.join('/'))
    console.log('  Stripped Path:', targetPath)
    console.log('  Full URL:', targetUrl)
    console.log('  Method:', request.method)

    // Copy headers from incoming request
    const headers = new Headers()
    request.headers.forEach((value, key) => {
      // Skip host header as we're changing the target
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value)
      }
    })

    console.log('📤 Making request to Supabase Auth...')

    // Read body once if it exists
    let body = null
    if (request.method !== 'GET' && request.method !== 'HEAD') {
      try {
        body = await request.text()
        console.log('📦 Request body length:', body.length)
      } catch (e) {
        console.log('⚠️  No body or already consumed')
      }
    }

    // Make request to Supabase Auth with timeout
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), 15000) // 15 second timeout

    const response = await fetch(targetUrl, {
      method: request.method,
      headers: headers,
      body: body,
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    console.log('✅ Proxy response status:', response.status)
    console.log('✅ Response headers:', Object.fromEntries(response.headers.entries()))

    // Get response body
    const data = await response.text()
    
    console.log('✅ Response body length:', data.length)

    // Return response with CORS headers
    return new NextResponse(data, {
      status: response.status,
      statusText: response.statusText,
      headers: {
        'Content-Type': response.headers.get('Content-Type') || 'application/json',
        'Access-Control-Allow-Origin': 'https://offensivewizard.com',
        'Access-Control-Allow-Credentials': 'true',
        'Access-Control-Expose-Headers': 'X-Total-Count',
      },
    })
  } catch (error) {
    console.error('❌ Proxy error details:')
    console.error('  Error type:', error instanceof Error ? error.constructor.name : typeof error)
    console.error('  Error message:', error instanceof Error ? error.message : String(error))
    console.error('  GOTRUE_URL was:', GOTRUE_URL)
    
    // Check if it's a network error
    if (error instanceof Error && (error.message.includes('ECONNREFUSED') || error.message.includes('ENOTFOUND'))) {
      console.error('  ⚠️  Cannot reach supabase-auth service - check Docker network configuration')
    }

    return new NextResponse(
      JSON.stringify({
        error: 'Proxy error',
        message: error instanceof Error ? error.message : 'Unknown error',
        details: {
          targetUrl: GOTRUE_URL,
          errorType: error instanceof Error ? error.constructor.name : typeof error,
        }
      }),
      {
        status: 500,
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': 'https://offensivewizard.com',
          'Access-Control-Allow-Credentials': 'true',
        },
      }
    )
  }
}
