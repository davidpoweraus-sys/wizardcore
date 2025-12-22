import { createBrowserClient } from '@supabase/ssr'

export function createClient() {
  return createBrowserClient(
    process.env.NEXT_PUBLIC_SUPABASE_URL!,
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!,
    {
      cookies: {
        get(name: string) {
          // Get cookie from browser
          const cookie = document.cookie
            .split('; ')
            .find((row) => row.startsWith(`${name}=`))
          return cookie ? decodeURIComponent(cookie.split('=')[1]) : null
        },
        set(name: string, value: string, options: any) {
          // Set cookie in browser with proper options for cross-domain
          let cookieString = `${name}=${encodeURIComponent(value)}`
          
          if (options?.maxAge) {
            cookieString += `; Max-Age=${options.maxAge}`
          }
          
          // Set path to root
          cookieString += '; Path=/'
          
          // Allow cross-domain cookie sharing if domains match
          if (options?.domain) {
            cookieString += `; Domain=${options.domain}`
          }
          
          // Set SameSite for cross-domain requests
          cookieString += '; SameSite=Lax'
          
          // Use Secure in production
          if (window.location.protocol === 'https:') {
            cookieString += '; Secure'
          }
          
          document.cookie = cookieString
        },
        remove(name: string, options: any) {
          // Remove cookie by setting expiry to past
          let cookieString = `${name}=; Path=/; Max-Age=0`
          
          if (options?.domain) {
            cookieString += `; Domain=${options.domain}`
          }
          
          document.cookie = cookieString
        },
      },
    }
  )
}