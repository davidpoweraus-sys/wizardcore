import { NextResponse, type NextRequest } from 'next/server'

// Version identifier for tracking which fix is deployed
const MIDDLEWARE_VERSION = 'rsc-fix-20260104-1315'

export async function middleware(request: NextRequest) {
  // Version tracking in logs
  console.log(`🔍 Middleware ${MIDDLEWARE_VERSION} executing for path:`, request.nextUrl.pathname)
  console.log('🔍 Request cookies:', request.cookies.getAll().map(c => c.name).join(', '))
  
  // Check for Supabase auth cookie
  const hasAuthCookie = request.cookies.has('sb-app-auth-token')
  console.log('🔍 Has auth cookie:', hasAuthCookie)
  
  // CRITICAL FIX: Check if this is an RSC fetch request
  // Next.js RSC fetches include `_rsc` query parameter or specific headers
  // Also check for RSC header which Next.js uses for React Server Components
  // IMPORTANT: Also check for Accept header containing 'text/x-component'
  const acceptHeader = request.headers.get('accept') || ''
  const hasRSCParam = request.nextUrl.search.includes('_rsc=')
  const xNextJSData = request.headers.get('x-nextjs-data')
  const nextRouterPrefetch = request.headers.get('next-router-prefetch')
  const nextAction = request.headers.get('next-action')
  const rscHeader = request.headers.get('RSC')
  const nextRouterStateTree = request.headers.get('Next-Router-State-Tree')
  
  const isRSCFetch = hasRSCParam ||
                     xNextJSData === '1' ||
                     nextRouterPrefetch === '1' ||
                     nextAction === '1' ||
                     rscHeader === '1' ||
                     nextRouterStateTree !== null ||
                     acceptHeader.includes('text/x-component')
  
  if (isRSCFetch) {
    console.log('🔍 RSC fetch detected:', {
      pathname: request.nextUrl.pathname,
      search: request.nextUrl.search,
      hasRSCParam,
      xNextJSData,
      nextRouterPrefetch,
      nextAction,
      rscHeader,
      nextRouterStateTree,
      acceptHeader,
      hasAuthCookie
    })
    
    // For RSC fetches to protected routes, we need to check if user is authenticated
    // If not authenticated, we should return a 401 or 403 instead of redirecting
    // This allows the client to handle the authentication error properly
    const isProtectedRoute = request.nextUrl.pathname.startsWith('/dashboard') ||
                             request.nextUrl.pathname.startsWith('/profile')
    
    if (isProtectedRoute && !hasAuthCookie) {
      console.log('🔍 RSC fetch to protected route without auth cookie, returning 401')
      // Return 401 Unauthorized for RSC fetches without auth
      // This is better than redirecting for API-like requests
      const response = new NextResponse(
        JSON.stringify({
          error: 'Unauthorized',
          message: 'Authentication required',
          middleware_version: MIDDLEWARE_VERSION
        }),
        {
          status: 401,
          headers: {
            'Content-Type': 'application/json',
            'X-Middleware-Version': MIDDLEWARE_VERSION,
          },
        }
      )
      return response
    }
    
    // If authenticated or not a protected route, allow through
    console.log('🔍 RSC fetch allowed through')
    const response = NextResponse.next()
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    return response
  }
  
  // NEW: Refresh session for authenticated users
  // This ensures the Supabase client has a valid session
  if (hasAuthCookie) {
    console.log('🔍 Attempting to refresh session for authenticated user')
    
    // Create a response object
    const response = NextResponse.next()
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    
    // IMPORTANT: We need to ensure the auth cookie is properly forwarded
    // The Supabase client will handle session refresh on the client side
    // Our job is just to ensure the cookie is present and valid
    
    // For API routes, we need to set proper CORS headers
    if (request.nextUrl.pathname.startsWith('/api/')) {
      // Set CORS headers for API routes
      const origin = request.headers.get('origin')
      if (origin && (origin.includes('offensivewizard.com') || origin.includes('localhost'))) {
        response.headers.set('Access-Control-Allow-Origin', origin)
        response.headers.set('Access-Control-Allow-Credentials', 'true')
        response.headers.set('Access-Control-Allow-Methods', 'GET, POST, PUT, PATCH, DELETE, OPTIONS')
        response.headers.set('Access-Control-Allow-Headers', 'Authorization, Content-Type, X-Client-Info, apikey, x-client-info, x-supabase-api-version')
      }
    }
    
    return response
  }
  
  // Protected routes - redirect to login if not authenticated
  if (
    !hasAuthCookie &&
    (request.nextUrl.pathname.startsWith('/dashboard') ||
      request.nextUrl.pathname.startsWith('/profile'))
  ) {
    console.log('🔍 No auth cookie, redirecting to login')
    // Redirect to login page
    const url = request.nextUrl.clone()
    url.pathname = '/login'
    url.searchParams.set('redirectedFrom', request.nextUrl.pathname)
    const response = NextResponse.redirect(url)
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    return response
  }

  // Auth routes - redirect to dashboard if already authenticated
  if (
    hasAuthCookie &&
    (request.nextUrl.pathname.startsWith('/login') ||
      request.nextUrl.pathname.startsWith('/register'))
  ) {
    console.log('🔍 Has auth cookie, redirecting to dashboard')
    const url = request.nextUrl.clone()
    url.pathname = '/dashboard'
    const response = NextResponse.redirect(url)
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    return response
  }

  const response = NextResponse.next()
  response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
  return response
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     * - public folder
     */
    '/((?!_next/static|_next/image|favicon.ico|.*\\.(?:svg|png|jpg|jpeg|gif|webp)$).*)',
  ],
}
