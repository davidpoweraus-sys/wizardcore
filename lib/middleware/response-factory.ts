/**
 * Middleware Response Factory
 * Creates consistent responses for middleware
 */

import { NextRequest, NextResponse } from 'next/server'
import {
  MiddlewareResponseOptions,
  DEFAULT_RESPONSE_OPTIONS,
  ERROR_CODES,
  ERROR_MESSAGES,
  MIDDLEWARE_VERSION
} from './config'
import { createCorsHeaders, addSecurityHeaders } from './helpers'

// ============================================
// RESPONSE FACTORY
// ============================================

export class MiddlewareResponseFactory {
  /**
   * Create a JSON error response
   */
  static createErrorResponse(
    request: NextRequest,
    status: number,
    errorCode: keyof typeof ERROR_CODES,
    additionalData: Record<string, unknown> = {}
  ): NextResponse {
    const errorData = {
      error: ERROR_CODES[errorCode],
      message: ERROR_MESSAGES[errorCode],
      code: errorCode,
      path: request.nextUrl.pathname,
      timestamp: new Date().toISOString(),
      ...additionalData
    }
    
    const response = new NextResponse(
      JSON.stringify(errorData),
      {
        status,
        headers: {
          'Content-Type': 'application/json',
          'X-Middleware-Version': MIDDLEWARE_VERSION,
          ...addSecurityHeaders()
        }
      }
    )
    
    // Add CORS headers if needed
    const corsHeaders = createCorsHeaders(request)
    Object.entries(corsHeaders).forEach(([key, value]) => {
      response.headers.set(key, value)
    })
    
    return response
  }
  
  /**
   * Create a redirect response
   */
  static createRedirectResponse(
    request: NextRequest,
    redirectTo: string,
    status: number = 307
  ): NextResponse {
    const url = request.nextUrl.clone()
    url.pathname = redirectTo
    
    // Add redirectedFrom parameter for login redirects
    if (redirectTo === '/login' && request.nextUrl.pathname !== '/login') {
      url.searchParams.set('redirectedFrom', request.nextUrl.pathname)
    }
    
    const response = NextResponse.redirect(url, status)
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    
    // Add security headers
    Object.entries(addSecurityHeaders()).forEach(([key, value]) => {
      response.headers.set(key, value)
    })
    
    return response
  }
  
  /**
   * Create a Next.js response (continue to next handler)
   */
  static createNextResponse(
    request: NextRequest,
    additionalHeaders: Record<string, string> = {}
  ): NextResponse {
    const response = NextResponse.next()
    
    // Set middleware version
    response.headers.set('X-Middleware-Version', MIDDLEWARE_VERSION)
    
    // Add security headers
    Object.entries(addSecurityHeaders()).forEach(([key, value]) => {
      response.headers.set(key, value)
    })
    
    // Add CORS headers if needed
    const corsHeaders = createCorsHeaders(request)
    Object.entries(corsHeaders).forEach(([key, value]) => {
      response.headers.set(key, value)
    })
    
    // Add any additional headers
    Object.entries(additionalHeaders).forEach(([key, value]) => {
      response.headers.set(key, value)
    })
    
    return response
  }
  
  /**
   * Create a response based on options
   */
  static createResponse(
    request: NextRequest,
    options: MiddlewareResponseOptions = {}
  ): NextResponse {
    const mergedOptions = { ...DEFAULT_RESPONSE_OPTIONS, ...options }
    
    switch (mergedOptions.responseType) {
      case 'json':
        return this.createErrorResponse(
          request,
          mergedOptions.status!,
          mergedOptions.message?.includes('Authentication') ? 'AUTH_REQUIRED' : 'INTERNAL_ERROR',
          { message: mergedOptions.message }
        )
        
      case 'redirect':
        if (!mergedOptions.redirectTo) {
          throw new Error('redirectTo is required for redirect response type')
        }
        return this.createRedirectResponse(
          request,
          mergedOptions.redirectTo,
          mergedOptions.status
        )
        
      case 'next':
      default:
        return this.createNextResponse(
          request,
          mergedOptions.headers
        )
    }
  }
  
  /**
   * Create unauthorized response (401)
   */
  static createUnauthorizedResponse(
    request: NextRequest,
    message?: string
  ): NextResponse {
    return this.createErrorResponse(
      request,
      401,
      'AUTH_REQUIRED',
      { message: message || ERROR_MESSAGES.AUTH_REQUIRED }
    )
  }
  
  /**
   * Create forbidden response (403)
   */
  static createForbiddenResponse(
    request: NextRequest,
    message?: string
  ): NextResponse {
    return this.createErrorResponse(
      request,
      403,
      'ACCESS_DENIED',
      { message: message || ERROR_MESSAGES.ACCESS_DENIED }
    )
  }
  
  /**
   * Create rate limited response (429)
   */
  static createRateLimitedResponse(
    request: NextRequest,
    retryAfter?: number
  ): NextResponse {
    const headers: Record<string, string> = {
      'X-Middleware-Version': MIDDLEWARE_VERSION,
      ...addSecurityHeaders()
    }
    
    if (retryAfter) {
      headers['Retry-After'] = retryAfter.toString()
    }
    
    return this.createErrorResponse(
      request,
      429,
      'RATE_LIMITED',
      { retryAfter }
    )
  }
  
  /**
   * Create internal error response (500)
   */
  static createInternalErrorResponse(
    request: NextRequest,
    error?: Error
  ): NextResponse {
    const errorData: Record<string, unknown> = {
      message: ERROR_MESSAGES.INTERNAL_ERROR
    }
    
    // Include error details in development
    if (process.env.NODE_ENV === 'development' && error) {
      errorData.debug = {
        name: error.name,
        message: error.message,
        stack: error.stack
      }
    }
    
    return this.createErrorResponse(
      request,
      500,
      'INTERNAL_ERROR',
      errorData
    )
  }
}

// ============================================
// CONVENIENCE FUNCTIONS
// ============================================

/**
 * Create response for authenticated requests to protected routes
 */
export function handleProtectedRoute(
  request: NextRequest,
  hasAuthCookie: boolean,
  isRSCRequest: boolean
): NextResponse | null {
  if (hasAuthCookie) {
    return null // Allow access
  }
  
  // Different response based on request type
  if (isRSCRequest || request.nextUrl.pathname.startsWith('/api/')) {
    return MiddlewareResponseFactory.createUnauthorizedResponse(request)
  } else {
    return MiddlewareResponseFactory.createRedirectResponse(request, '/login')
  }
}

/**
 * Create response for authenticated users accessing auth routes
 */
export function handleAuthRoute(
  request: NextRequest,
  hasAuthCookie: boolean
): NextResponse | null {
  if (!hasAuthCookie) {
    return null // Allow access (user needs to authenticate)
  }
  
  // Redirect authenticated users away from auth pages
  return MiddlewareResponseFactory.createRedirectResponse(request, '/dashboard')
}

/**
 * Create response for API routes with CORS
 */
export function handleApiRoute(
  request: NextRequest,
  response: NextResponse
): NextResponse {
  // Ensure CORS headers are set for API routes
  const corsHeaders = createCorsHeaders(request)
  Object.entries(corsHeaders).forEach(([key, value]) => {
    response.headers.set(key, value)
  })
  
  return response
}