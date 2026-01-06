/**
 * Middleware Configuration
 * Centralized configuration for middleware constants and types
 */

// ============================================
// ENVIRONMENT CONFIGURATION
// ============================================

export const MIDDLEWARE_VERSION = 'refactored-v1-20250106'

// Cookie configuration
export const AUTH_COOKIE_NAME = 'sb-app-auth-token'
export const AUTH_COOKIE_OPTIONS = {
  path: '/',
  httpOnly: true,
  secure: process.env.NODE_ENV === 'production',
  sameSite: 'lax' as const,
  maxAge: 60 * 60 * 24 * 7 // 7 days
}

// Route configuration
export const PROTECTED_ROUTE_PREFIXES = [
  '/dashboard',
  '/profile', 
  '/creator',
  '/settings'
] as const

export const AUTH_ROUTE_PREFIXES = [
  '/login',
  '/register',
  '/auth'
] as const

export const PUBLIC_ROUTE_PREFIXES = [
  '/',
  '/about',
  '/pricing',
  '/features',
  '/api/public'
] as const

export const API_ROUTE_PREFIX = '/api/'

// CORS configuration
export const ALLOWED_ORIGINS = (() => {
  const envOrigins = process.env.ALLOWED_ORIGINS?.split(',') || []
  const defaultOrigins = [
    'https://app.offensivewizard.com',
    'http://localhost:3000',
    'http://localhost:3001'
  ]
  
  return [...new Set([...envOrigins, ...defaultOrigins])]
})()

export const ALLOWED_METHODS = [
  'GET',
  'POST', 
  'PUT',
  'PATCH',
  'DELETE',
  'OPTIONS'
] as const

export const ALLOWED_HEADERS = [
  'Authorization',
  'Content-Type',
  'X-Client-Info',
  'apikey',
  'x-client-info',
  'x-supabase-api-version',
  'Next-Router-Prefetch',
  'Next-Router-State-Tree',
  'Next-Url',
  'RSC'
] as const

// Security headers configuration
export const SECURITY_HEADERS = {
  'X-Content-Type-Options': 'nosniff',
  'X-Frame-Options': 'DENY',
  'Referrer-Policy': 'strict-origin-when-cross-origin',
  'Permissions-Policy': 'camera=(), microphone=(), geolocation=()',
  'X-XSS-Protection': '1; mode=block'
} as const

// ============================================
// TYPE DEFINITIONS
// ============================================

export interface RouteAnalysis {
  pathname: string
  isProtected: boolean
  isAuthRoute: boolean
  isApiRoute: boolean
  isPublicRoute: boolean
}

export interface AuthAnalysis {
  hasAuthCookie: boolean
  requiresAuthentication: boolean
  shouldRedirectToLogin: boolean
  shouldRedirectToDashboard: boolean
}

export interface RSCAnalysis {
  isRSCRequest: boolean
  detectionMethod?: 'RSC-header' | 'accept-header' | 'query-param'
}

export interface RequestAnalysis {
  route: RouteAnalysis
  auth: AuthAnalysis
  rsc: RSCAnalysis
}

export interface MiddlewareResponseOptions {
  status?: number
  message?: string
  redirectTo?: string
  responseType?: 'json' | 'redirect' | 'next'
  headers?: Record<string, string>
}

export interface ErrorResponseData {
  error: string
  message: string
  code?: string
  path?: string
  timestamp?: string
}

// ============================================
// ERROR CODES AND MESSAGES
// ============================================

export const ERROR_CODES = {
  AUTH_REQUIRED: 'AUTH_REQUIRED',
  INVALID_SESSION: 'INVALID_SESSION',
  ACCESS_DENIED: 'ACCESS_DENIED',
  RATE_LIMITED: 'RATE_LIMITED',
  INTERNAL_ERROR: 'INTERNAL_ERROR'
} as const

export const ERROR_MESSAGES = {
  [ERROR_CODES.AUTH_REQUIRED]: 'Authentication required to access this resource',
  [ERROR_CODES.INVALID_SESSION]: 'Your session has expired or is invalid',
  [ERROR_CODES.ACCESS_DENIED]: 'You do not have permission to access this resource',
  [ERROR_CODES.RATE_LIMITED]: 'Too many requests, please try again later',
  [ERROR_CODES.INTERNAL_ERROR]: 'An internal server error occurred'
} as const

// ============================================
// RESPONSE FACTORY DEFAULTS
// ============================================

export const DEFAULT_RESPONSE_OPTIONS: MiddlewareResponseOptions = {
  status: 200,
  responseType: 'next',
  headers: {
    'X-Middleware-Version': MIDDLEWARE_VERSION
  }
} as const

export const UNAUTHORIZED_RESPONSE: MiddlewareResponseOptions = {
  status: 401,
  message: ERROR_MESSAGES[ERROR_CODES.AUTH_REQUIRED],
  responseType: 'json',
  headers: {
    'X-Middleware-Version': MIDDLEWARE_VERSION,
    'Content-Type': 'application/json'
  }
} as const

export const FORBIDDEN_RESPONSE: MiddlewareResponseOptions = {
  status: 403,
  message: ERROR_MESSAGES[ERROR_CODES.ACCESS_DENIED],
  responseType: 'json',
  headers: {
    'X-Middleware-Version': MIDDLEWARE_VERSION,
    'Content-Type': 'application/json'
  }
} as const