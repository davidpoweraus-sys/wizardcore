import { NextResponse, type NextRequest } from 'next/server'

// Constants
const MIDDLEWARE_VERSION = 'refactored-20250106'
const AUTH_COOKIE_NAME = 'sb-app-auth-token'
const ALLOWED_ORIGINS = process.env.ALLOWED_ORIGINS?.split(',') || [
  'https://app.offensivewizard.com',
  'http://localhost:3000'
]

// Route patterns
const PROTECTED_ROUTES = ['/dashboard', '/profile', '/creator']
const AUTH_ROUTES = ['/login', '/register', '/auth']
const PUBLIC_ROUTES = ['/', '/about', '/pricing', '/api/public']

// Helper functions
function isProtectedRoute(pathname: string): boolean {
  return PROTECTED_ROUTES.some(route => pathname.startsWith(route))
}

function isAuthRoute(pathname: string): boolean {
  return AUTH_ROUTES.some(route => pathname.startsWith(route))
}

function isPublicRoute(pathname: string): boolean {
  return PUBLIC_ROUTES.some(route => pathname.startsWith(route))
}

function isApiRoute(pathname: string): boolean {
  return pathname.startsWith('/api/')
}

function isRSCRequest(request: NextRequest): boolean {
  // Simplified RSC detection based on Next.js patterns
  const acceptHeader = request.headers.get('accept') || ''
  const rscHeader = request.headers.get('RSC')
  const nextAction = request.headers.get('next-action')
  
  return acceptHeader.includes('text/x-component') || 
         rscHeader === '1' || 
         nextAction === '1' ||
         request.nextUrl.search.includes('_rsc=')
}

function addSecurityHeaders(response: NextResponse): NextResponse {
  // Add security headers to all responses
  response.headers.set('X-Content-Type-Options', 'nosniff')
  response.headers.set('X-Frame-Options', 'DENY')
  response.headers.set('Referrer-Policy', 'strict-origin-when-cross-origin')
  response.headers.set('Permissions-Policy', 'camera=(), microphone=(), geolocation=()')
  
  // Add middleware version for debugging
  response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
  
  return response
}

function addCorsHeaders(request: NextRequest, response: NextResponse): NextResponse {
  if (!isApiRoute(request.nextUrl.pathname)) {
    return response
  }
  
  const origin = request.headers.get('origin')
  if (origin && ALLOWED_ORIGINS.some(allowed => origin.includes(allowed))) {
    response.headers.set('Access-Control-Allow-Origin', origin)
    response.headers.set('Access-Control-Allow-Credentials', 'true')
    response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, PATCH, DELETE, OPTIONS')
    response.headers.set('Access-Control-Allow-Headers', 
      'Authorization, Content-Type, X-Client-Info, apikey, x-client-info, x-supabase-api-version, Next-Router-Prefetch, Next-Router-State-Tree, Next-Url, RSC'
    )
  }
  
  return response
}

function createUnauthorizedResponse(request: NextRequest): NextResponse {
  const isRSC = isRSCRequest(request)
  
  if (isRSC || isApiRoute(request.nextUrl.pathname)) {
    // JSON response for API/RSC requests
    return new NextResponse(
      JSON.stringify({
        error: 'Unauthorized',
        message: 'Authentication required',
        code: 'AUTH_REQUIRED'
      }),
      {
        status: 401,
        headers: {
          'Content-Type': 'application/json',
          'X-Middleware-Version': MIDDLEWARE_VERSION,
        },
      }
    )
  } else {
    // Redirect for browser requests
    const url = request.nextUrl.clone()
    url.pathname = '/login'
    url.searchParams.set('redirectedFrom', request.nextUrl.pathname)
    return NextResponse.redirect(url)
  }
}

function createAlreadyAuthenticatedResponse(request: NextRequest): NextResponse {
  const url = request.nextUrl.clone()
  url.pathname = '/dashboard'
  return NextResponse.redirect(url)
}

export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname
  const hasAuthCookie = request.cookies.has(AUTH_COOKIE_NAME)
  
  // Skip middleware for static files and public routes
  if (isPublicRoute(pathname) && !isProtectedRoute(pathname)) {
    const response = NextResponse.next()
    return addSecurityHeaders(response)
  }
  
  // Handle API routes with CORS
  if (isApiRoute(pathname)) {
    const response = NextResponse.next()
    const withCors = addCorsHeaders(request, response)
    return addSecurityHeaders(withCors)
  }
  
  // Check authentication for protected routes
  if (isProtectedRoute(pathname) && !hasAuthCookie) {
    return createUnauthorizedResponse(request)
  }
  
  // Redirect authenticated users away from auth pages
  if (isAuthRoute(pathname) && hasAuthCookie) {
    return createAlreadyAuthenticatedResponse(request)
  }
  
  // Default: allow request with security headers
  const response = NextResponse.next()
  return addSecurityHeaders(response)
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     * - api/public (public API routes)
     */
    '/((?!_next/static|_next/image|favicon.ico|public/|api/public/).*)',
  ],
}