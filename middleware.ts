/**
 * Refactored Middleware
 * Addressing problems 4, 5, 6, 9, 10, 11
 * 
 * Problem 4: Overly Complex RSC Detection - Simplified
 * Problem 5: Inconsistent Error Handling - Standardized
 * Problem 6: Session Refresh Logic is Misplaced - Separated
 * Problem 9: Magic Strings - Extracted to constants
 * Problem 10: Missing Type Safety - Added TypeScript interfaces
 * Problem 11: Logic Duplication - Extracted to helper functions
 */

import { NextRequest } from 'next/server'
import {
  analyzeRequest,
  shouldExcludeFromMiddleware,
  logMiddlewareEvent
} from '@/lib/middleware/helpers'
import {
  handleProtectedRoute,
  handleAuthRoute,
  handleApiRoute
} from '@/lib/middleware/response-factory'
import { MiddlewareResponseFactory } from '@/lib/middleware/response-factory'

export async function middleware(request: NextRequest) {
  const pathname = request.nextUrl.pathname
  
  // Skip middleware for excluded paths (static files, etc.)
  if (shouldExcludeFromMiddleware(pathname)) {
    logMiddlewareEvent('path-excluded', { pathname }, 'debug')
    return MiddlewareResponseFactory.createNextResponse(request)
  }
  
  try {
    // Analyze the request
    const analysis = analyzeRequest(request)
    
    logMiddlewareEvent('request-analyzed', {
      pathname,
      isProtected: analysis.route.isProtected,
      isAuthRoute: analysis.route.isAuthRoute,
      isApiRoute: analysis.route.isApiRoute,
      hasAuthCookie: analysis.auth.hasAuthCookie,
      isRSCRequest: analysis.rsc.isRSCRequest
    }, 'debug')
    
    // Handle protected routes (require authentication)
    if (analysis.route.isProtected) {
      const protectedResponse = handleProtectedRoute(
        request,
        analysis.auth.hasAuthCookie,
        analysis.rsc.isRSCRequest
      )
      
      if (protectedResponse) {
        logMiddlewareEvent('protected-route-handled', {
          pathname,
          action: 'redirect-or-deny'
        }, 'info')
        return protectedResponse
      }
    }
    
    // Handle auth routes (redirect authenticated users)
    if (analysis.route.isAuthRoute) {
      const authRouteResponse = handleAuthRoute(
        request,
        analysis.auth.hasAuthCookie
      )
      
      if (authRouteResponse) {
        logMiddlewareEvent('auth-route-handled', {
          pathname,
          action: 'redirect-to-dashboard'
        }, 'info')
        return authRouteResponse
      }
    }
    
    // Create base response
    let response = MiddlewareResponseFactory.createNextResponse(request)
    
    // Handle API routes (add CORS headers)
    if (analysis.route.isApiRoute) {
      response = handleApiRoute(request, response)
      logMiddlewareEvent('api-route-handled', { pathname }, 'debug')
    }
    
    // Log successful request handling
    logMiddlewareEvent('request-completed', {
      pathname,
      status: 'allowed'
    }, 'info')
    
    return response
    
  } catch (error) {
    // Handle middleware errors gracefully
    logMiddlewareEvent('middleware-error', {
      pathname,
      error: error instanceof Error ? error.message : 'Unknown error',
      stack: error instanceof Error ? error.stack : undefined
    }, 'error')
    
    // Don't block requests on middleware errors in production
    if (process.env.NODE_ENV === 'production') {
      return MiddlewareResponseFactory.createNextResponse(request)
    } else {
      return MiddlewareResponseFactory.createInternalErrorResponse(
        request,
        error instanceof Error ? error : undefined
      )
    }
  }
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