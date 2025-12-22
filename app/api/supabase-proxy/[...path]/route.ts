import { NextRequest, NextResponse } from 'next/server'

// Internal Supabase Auth URL (within Docker network)
const GOTRUE_URL = process.env.SUPABASE_INTERNAL_URL || 'http://supabase-auth:9999'

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
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

// Handle POST requests
export async function POST(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

// Handle PUT requests
export async function PUT(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

// Handle PATCH requests
export async function PATCH(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

// Handle DELETE requests
export async function DELETE(
  request: NextRequest,
  { params }: { params: { path: string[] } }
) {
  return proxyRequest(request, params.path)
}

async function proxyRequest(request: NextRequest, path: string[]) {
  try {
    const url = new URL(request.url)
    const targetPath = path.join('/')
    const targetUrl = `${GOTRUE_URL}/${targetPath}${url.search}`

    console.log('🔄 Proxying request to:', targetUrl)

    // Copy headers from incoming request
    const headers = new Headers()
    request.headers.forEach((value, key) => {
      // Skip host header as we're changing the target
      if (key.toLowerCase() !== 'host') {
        headers.set(key, value)
      }
    })

    // Make request to internal Supabase Auth
    const response = await fetch(targetUrl, {
      method: request.method,
      headers: headers,
      body: request.body,
      // @ts-ignore - duplex is needed for streaming
      duplex: 'half',
    })

    // Get response body
    const data = await response.text()

    console.log('✅ Proxy response status:', response.status)

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
    console.error('❌ Proxy error:', error)
    return new NextResponse(
      JSON.stringify({
        error: 'Proxy error',
        message: error instanceof Error ? error.message : 'Unknown error',
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
