'use client'

import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { createClient } from '@/lib/supabase/client'
import { useRouter } from 'next/navigation'

export default function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()
  const supabase = createClient()

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    // BLUE DIE TEST: Log to browser console for debugging
    console.log('🎲 BLUE DIE TEST - Login attempt starting')
    console.log('  Email:', email)
    console.log('  Supabase URL:', process.env.NEXT_PUBLIC_SUPABASE_URL)
    console.log('  Timestamp:', new Date().toISOString())

    const { error } = await supabase.auth.signInWithPassword({
      email,
      password,
    })

    if (error) {
      console.error('🎲 BLUE DIE TEST - Login error:', error.message)
      console.error('  Error details:', error)
      setError(error.message)
      setLoading(false)
      return
    }

    console.log('🎲 BLUE DIE TEST - Login successful - Starting post-login flow')
    
    try {
      // CRITICAL FIX: Wait a moment for cookies to be set before redirecting
      // This fixes the race condition where middleware doesn't see the auth cookie
      console.log('🎲 Step 1: Waiting 100ms for cookie propagation...')
      await new Promise(resolve => setTimeout(resolve, 100))
      console.log('🎲 Step 1: Complete')
      
      // Also refresh the session to ensure cookies are properly set
      console.log('🎲 Step 2: Getting session...')
      const { data: { session }, error: sessionError } = await supabase.auth.getSession()
      console.log('🎲 Step 2: Session result - session:', session ? 'present' : 'absent')
      if (sessionError) {
        console.error('🎲 Step 2: Session error:', sessionError)
      }
      
      // Check cookies
      console.log('🎲 Step 3: Checking browser cookies...')
      const cookies = document.cookie
      console.log('🎲 Step 3: Cookies:', cookies)
      const hasAuthCookie = cookies.includes('sb-')
      console.log('🎲 Step 3: Has auth cookie:', hasAuthCookie)
      
      console.log('🎲 Step 4: Redirecting to dashboard...')
      
      // CRITICAL FIX: Use window.location.href immediately to avoid Next.js RSC issues
      // The router.push() is failing with RSC payload fetch errors
      console.log('🎲 Step 4a: Using window.location.href (bypassing Next.js router)')
      window.location.href = '/dashboard'
    } catch (error) {
      console.error('🎲 ERROR in post-login flow:', error)
      setError('Login successful but redirect failed: ' + (error instanceof Error ? error.message : String(error)))
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="text-center mb-8">
          <div className="flex justify-center mb-6">
            <Image
              src="/wizard_logo.png"
              alt="WizardCore Logo"
              width={120}
              height={120}
              className="drop-shadow-[0_0_15px_rgba(138,43,226,0.5)]"
              priority
            />
          </div>
           <h1 className="text-3xl font-bold bg-gradient-to-r from-pink-500 to-cyan-500 bg-clip-text text-transparent">
            Welcome Back
          </h1>
          <p className="text-text-secondary mt-2">
            Sign in to continue your coding journey
          </p>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-8 shadow-2xl">
          <form onSubmit={handleLogin} className="space-y-6">
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-text-secondary mb-2">
                Email Address
              </label>
              <input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender focus:border-transparent transition"
                placeholder="you@example.com"
                required
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-medium text-text-secondary mb-2">
                Password
              </label>
              <input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="w-full px-4 py-3 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-neon-lavender focus:border-transparent transition"
                placeholder="••••••••"
                required
              />
            </div>

            {error && (
              <div className="p-3 bg-red-900/30 border border-red-700 rounded-lg text-red-300 text-sm">
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full py-3 px-4 bg-gradient-to-r from-pink-500 to-purple-600 text-white font-semibold rounded-lg hover:opacity-90 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Signing in...' : 'Sign In'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-text-tertiary text-sm">
              Don't have an account?{' '}
              <Link href="/register" className="text-neon-cyan hover:underline font-medium">
                Sign up
              </Link>
            </p>
          </div>

          <div className="mt-8 pt-6 border-t border-border-subtle">
            <p className="text-xs text-text-muted text-center">
              By signing in, you agree to our Terms of Service and Privacy Policy.
            </p>
          </div>
        </div>
      </div>
    </div>
  )
}