import { NextResponse, type NextRequest } from 'next/server'

export async function middleware(request: NextRequest) {
  // Debug logging for production login issue
  console.log('🔍 Middleware executing for path:', request.nextUrl.pathname)
  console.log('🔍 Request cookies:', request.cookies.getAll().map(c => c.name).join(', '))
  
  // Check for Supabase auth cookie
  const hasAuthCookie = request.cookies.has('sb-app-auth-token')
  console.log('🔍 Has auth cookie:', hasAuthCookie)
  
  // CRITICAL FIX: Check if this is an RSC fetch request
  // Next.js RSC fetches include `_rsc` query parameter or specific headers
  const isRSCFetch = request.nextUrl.search.includes('_rsc=') ||
                     request.headers.get('x-nextjs-data') === '1' ||
                     request.headers.get('next-router-prefetch') === '1' ||
                     request.headers.get('next-action') === '1'
  
  if (isRSCFetch) {
    console.log('🔍 RSC fetch detected, checking authentication')
    
    // For RSC fetches to protected routes, we need to check if user is authenticated
    // If not authenticated, we should return a 401 or 403 instead of redirecting
    // This allows the client to handle the authentication error properly
    const isProtectedRoute = request.nextUrl.pathname.startsWith('/dashboard') ||
                             request.nextUrl.pathname.startsWith('/profile')
    
    if (isProtectedRoute && !hasAuthCookie) {
      console.log('🔍 RSC fetch to protected route without auth cookie, returning 401')
      // Return 401 Unauthorized for RSC fetches without auth
      // This is better than redirecting for API-like requests
      return new NextResponse(
        JSON.stringify({ error: 'Unauthorized', message: 'Authentication required' }),
        {
          status: 401,
          headers: {
            'Content-Type': 'application/json',
          },
        }
      )
    }
    
    // If authenticated or not a protected route, allow through
    console.log('🔍 RSC fetch allowed through')
    return NextResponse.next()
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
    return NextResponse.redirect(url)
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
    return NextResponse.redirect(url)
  }

  return NextResponse.next()
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
