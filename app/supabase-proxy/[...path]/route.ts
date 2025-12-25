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
    
    // STRIP /auth/v1 prefix - GoTrue serves endpoints at /* (root level)
    // IMPORTANT: Based on server diagnostics, GoTrue expects /signup NOT /auth/v1/signup
    targetPath = targetPath.replace(/^auth\/v1\//, '')
    
    // Ensure we don't get double slashes in URL
    const baseUrl = GOTRUE_URL.endsWith('/') ? GOTRUE_URL.slice(0, -1) : GOTRUE_URL
    const targetUrl = `${baseUrl}/${targetPath}${url.search}`

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
    console.log('✅ First 200 chars of response:', data.substring(0, 200))
    console.log('✅ Last 200 chars of response:', data.substring(Math.max(0, data.length - 200)))
    
    // Debug: Show character codes of first 50 chars to see what's causing JSON parse error
    console.log('🔍 First 50 character codes (to debug JSON parse error at column 5):')
    for (let i = 0; i < Math.min(50, data.length); i++) {
      const char = data[i]
      let displayChar = char
      if (char === '\n') displayChar = '\\n'
      else if (char === '\r') displayChar = '\\r'
      else if (char === '\t') displayChar = '\\t'
      else if (char === '"') displayChar = '\\"'
      else if (char === '\\') displayChar = '\\\\'
      else if (char.charCodeAt(0) < 32) displayChar = `\\x${char.charCodeAt(0).toString(16).padStart(2, '0')}`
      
      console.log(`  [${i}] '${displayChar}' = ${char.charCodeAt(0)} ${i === 4 ? '← COLUMN 5 (error location)' : ''}`)
    }
    
    // Also show the raw response as hex to see any hidden characters
    console.log('🔍 First 100 bytes as hex:')
    const hexArray = []
    for (let i = 0; i < Math.min(100, data.length); i++) {
      hexArray.push(data.charCodeAt(i).toString(16).padStart(2, '0'))
      if ((i + 1) % 16 === 0) {
        console.log(`  ${hexArray.slice(i - 15, i + 1).join(' ')}`)
      }
    }
    if (hexArray.length % 16 !== 0) {
      console.log(`  ${hexArray.slice(Math.floor(hexArray.length / 16) * 16).join(' ')}`)
    }
    
    // Remove UTF-8 BOM if present (common issue with some servers)
    let cleanedData = data
    if (cleanedData.length > 0 && cleanedData.charCodeAt(0) === 0xFEFF) {
      console.log('⚠️  Removing UTF-8 BOM from response')
      cleanedData = cleanedData.slice(1)
    }
    
    // Trim whitespace from response to fix JSON parsing issues
    const trimmedData = cleanedData.trim()
    console.log('✅ Trimmed response body length:', trimmedData.length)
    
    // Check if response is valid JSON
    let responseData = trimmedData
    let shouldReturnJson = false
    
    try {
      JSON.parse(responseData)
      console.log('✅ Response is valid JSON')
    } catch (e) {
      console.log('⚠️  Response is NOT valid JSON:', e instanceof Error ? e.message : String(e))
      // If it's not JSON but we got an error status, convert it to JSON
      // This is CRITICAL because Supabase client expects JSON responses
      if (response.status >= 400) {
        console.log('🔄 Converting non-JSON error response to JSON format for Supabase client')
        shouldReturnJson = true
        
        // Try to extract any meaningful error message from the response
        let errorMessage = 'Authentication service error'
        if (responseData.includes('Not Found') || responseData.includes('404')) {
          errorMessage = 'Authentication endpoint not found. Check if Supabase Auth is running.'
        } else if (responseData.length > 0) {
          // Take first 100 chars as error message
          errorMessage = responseData.substring(0, Math.min(100, responseData.length)).replace(/\s+/g, ' ')
        }
        
        responseData = JSON.stringify({
          error: 'auth_error',
          message: errorMessage,
          status: response.status,
          details: {
            original_status: response.status,
            original_status_text: response.statusText,
            content_type: response.headers.get('Content-Type') || 'unknown',
            response_preview: responseData.length > 200 ? responseData.substring(0, 200) + '...' : responseData
          }
        })
      }
    }

    // Return response with CORS headers
    // If we converted to JSON, use application/json content type
    const contentType = shouldReturnJson 
      ? 'application/json' 
      : response.headers.get('Content-Type') || 'application/json'
    
    return new NextResponse(responseData, {
      status: response.status,
      statusText: response.statusText,
      headers: {
        'Content-Type': contentType,
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
