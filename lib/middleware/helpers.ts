/**
 * Middleware Helper Functions
 * Reusable functions for middleware logic
 */

import { NextRequest } from 'next/server'
import {
  RouteAnalysis,
  AuthAnalysis,
  RSCAnalysis,
  RequestAnalysis,
  PROTECTED_ROUTE_PREFIXES,
  AUTH_ROUTE_PREFIXES,
  PUBLIC_ROUTE_PREFIXES,
  API_ROUTE_PREFIX,
  AUTH_COOKIE_NAME,
  ALLOWED_ORIGINS,
  ALLOWED_METHODS,
  ALLOWED_HEADERS,
  SECURITY_HEADERS,
  MIDDLEWARE_VERSION
} from './config'

// ============================================
// ROUTE ANALYSIS HELPERS
// ============================================

/**
 * Analyze route properties
 */
export function analyzeRoute(pathname: string): RouteAnalysis {
  const isProtected = PROTECTED_ROUTE_PREFIXES.some(prefix => 
    pathname.startsWith(prefix)
  )
  
  const isAuthRoute = AUTH_ROUTE_PREFIXES.some(prefix => 
    pathname.startsWith(prefix)
  )
  
  const isPublicRoute = PUBLIC_ROUTE_PREFIXES.some(prefix => 
    pathname.startsWith(prefix)
  )
  
  const isApiRoute = pathname.startsWith(API_ROUTE_PREFIX)
  
  return {
    pathname,
    isProtected,
    isAuthRoute,
    isApiRoute,
    isPublicRoute
  }
}

/**
 * Check if path should be excluded from middleware
 */
export function shouldExcludeFromMiddleware(pathname: string): boolean {
  // Exclude Next.js internal paths
  if (pathname.startsWith('/_next/')) {
    return true
  }
  
  // Exclude public assets
  if (pathname.startsWith('/public/')) {
    return true
  }
  
  // Exclude static file extensions
  const staticFileExtensions = [
    '.ico', '.png', '.jpg', '.jpeg', '.gif', '.svg',
    '.css', '.js', '.woff', '.woff2', '.ttf', '.eot'
  ]
  
  return staticFileExtensions.some(ext => pathname.endsWith(ext))
}

// ============================================
// AUTHENTICATION HELPERS
// ============================================

/**
 * Analyze authentication status
 */
export function analyzeAuth(
  request: NextRequest,
  routeAnalysis: RouteAnalysis
): AuthAnalysis {
  const hasAuthCookie = request.cookies.has(AUTH_COOKIE_NAME)
  
  // Determine if authentication is required
  const requiresAuthentication = routeAnalysis.isProtected && !hasAuthCookie
  
  // Determine redirect logic
  const shouldRedirectToLogin = requiresAuthentication && 
                               !routeAnalysis.isApiRoute && 
                               !routeAnalysis.isAuthRoute
  
  const shouldRedirectToDashboard = hasAuthCookie && routeAnalysis.isAuthRoute
  
  return {
    hasAuthCookie,
    requiresAuthentication,
    shouldRedirectToLogin,
    shouldRedirectToDashboard
  }
}

// ============================================
// RSC DETECTION HELPERS
// ============================================

/**
 * Detect RSC requests with simplified logic
 */
export function detectRSC(request: NextRequest): RSCAnalysis {
  // Method 1: RSC header (most reliable)
  const rscHeader = request.headers.get('RSC')
  if (rscHeader === '1') {
    return { isRSCRequest: true, detectionMethod: 'RSC-header' }
  }
  
  // Method 2: Accept header for text/x-component
  const acceptHeader = request.headers.get('accept') || ''
  if (acceptHeader.includes('text/x-component')) {
    return { isRSCRequest: true, detectionMethod: 'accept-header' }
  }
  
  // Method 3: _rsc query parameter (legacy/fallback)
  if (request.nextUrl.search.includes('_rsc=')) {
    return { isRSCRequest: true, detectionMethod: 'query-param' }
  }
  
  return { isRSCRequest: false }
}

// ============================================
// REQUEST ANALYSIS
// ============================================

/**
 * Comprehensive request analysis
 */
export function analyzeRequest(request: NextRequest): RequestAnalysis {
  const pathname = request.nextUrl.pathname
  const route = analyzeRoute(pathname)
  const auth = analyzeAuth(request, route)
  const rsc = detectRSC(request)
  
  return { route, auth, rsc }
}

// ============================================
// CORS HELPERS
// ============================================

/**
 * Check if origin is allowed
 */
export function isOriginAllowed(origin: string | null): boolean {
  if (!origin) return false
  
  return ALLOWED_ORIGINS.some(allowedOrigin => {
    // Exact match or subdomain match
    return origin === allowedOrigin || 
           origin.endsWith(`.${allowedOrigin.replace(/^https?:\/\//, '')}`)
  })
}

/**
 * Create CORS headers for response
 */
export function createCorsHeaders(
  request: NextRequest,
  existingHeaders: Record<string, string> = {}
): Record<string, string> {
  const origin = request.headers.get('origin')
  
  if (!isOriginAllowed(origin)) {
    return existingHeaders
  }
  
  return {
    ...existingHeaders,
    'Access-Control-Allow-Origin': origin!,
    'Access-Control-Allow-Credentials': 'true',
    'Access-Control-Allow-Methods': ALLOWED_METHODS.join(', '),
    'Access-Control-Allow-Headers': ALLOWED_HEADERS.join(', '),
    'Access-Control-Expose-Headers': 'X-Middleware-Version, X-Request-ID'
  }
}

// ============================================
// SECURITY HEADERS HELPERS
// ============================================

/**
 * Add security headers to response
 */
export function addSecurityHeaders(
  existingHeaders: Record<string, string> = {}
): Record<string, string> {
  return {
    ...existingHeaders,
    ...SECURITY_HEADERS
  }
}

// ============================================
// LOGGING HELPERS
// ============================================

/**
 * Structured logging for middleware
 */
export function logMiddlewareEvent(
  event: string,
  data: Record<string, unknown> = {},
  level: 'debug' | 'info' | 'warn' | 'error' = 'info'
): void {
  // Only log in development or if explicitly enabled
  if (process.env.NODE_ENV !== 'development' && 
      !process.env.MIDDLEWARE_LOGGING_ENABLED) {
    return
  }
  
  const logEntry = {
    timestamp: new Date().toISOString(),
    event,
    level,
    ...data,
    middlewareVersion: MIDDLEWARE_VERSION
  }
  
  // Use appropriate console method based on level
  const consoleMethod = {
    debug: console.debug,
    info: console.info,
    warn: console.warn,
    error: console.error
  }[level]
  
  consoleMethod(`[Middleware] ${event}:`, JSON.stringify(logEntry, null, 2))
}

// ============================================
// VALIDATION HELPERS
// ============================================

/**
 * Validate request for common issues
 */
export function validateRequest(request: NextRequest): string[] {
  const warnings: string[] = []
  
  // Check for missing host header
  if (!request.headers.get('host')) {
    warnings.push('Missing Host header')
  }
  
  // Check for suspicious user agent
  const userAgent = request.headers.get('user-agent') || ''
  if (userAgent.includes('curl') || userAgent.includes('wget')) {
    warnings.push('Non-browser user agent detected')
  }
  
  // Check for excessive query parameters
  const searchParams = request.nextUrl.searchParams
  if (searchParams.size > 20) {
    warnings.push('Excessive query parameters')
  }
  
  return warnings
}