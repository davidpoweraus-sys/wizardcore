'use client'

import { useState } from 'react'
import Link from 'next/link'
import Image from 'next/image'
import { createClient } from '@/lib/supabase/client'
import { useRouter } from 'next/navigation'

export default function RegisterPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const router = useRouter()
  const supabase = createClient()

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError(null)

    console.log('🚀 Registration started')
    console.log('📧 Email:', email)
    console.log('🌐 Supabase URL:', process.env.NEXT_PUBLIC_SUPABASE_URL)
    console.log('🔑 API Key present:', !!process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY)
    console.log('🔗 Redirect URL:', `${window.location.origin}/auth/callback`)

    if (password !== confirmPassword) {
      console.error('❌ Passwords do not match')
      setError('Passwords do not match')
      setLoading(false)
      return
    }

    if (password.length < 6) {
      console.error('❌ Password too short')
      setError('Password must be at least 6 characters')
      setLoading(false)
      return
    }

    try {
      console.log('📤 Calling Supabase signUp...')
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
        options: {
          emailRedirectTo: `${window.location.origin}/auth/callback`,
        },
      })

      if (error) {
        console.error('❌ Supabase error:', error)
        console.error('Error details:', {
          message: error.message,
          status: error.status,
          name: error.name,
        })
        setError(error.message)
        setLoading(false)
        return
      }

      console.log('✅ Registration successful!', data)
      console.log('👤 User:', data.user)
      console.log('🎫 Session:', data.session)

      // Create user in wizardcore database
      if (data.user) {
        try {
          console.log('📤 Creating user in wizardcore database...')
          const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_URL}/v1/users`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${data.session?.access_token}`,
            },
            body: JSON.stringify({
              supabase_user_id: data.user.id,
              email: data.user.email,
              display_name: data.user.email?.split('@')[0] || 'User',
            }),
          })

          if (!response.ok) {
            console.error('⚠️ Failed to create user in wizardcore database:', await response.text())
            // Continue anyway - user can still use the app, data will be missing
          } else {
            console.log('✅ User created in wizardcore database')
          }
        } catch (err) {
          console.error('⚠️ Error creating user in wizardcore database:', err)
          // Continue anyway
        }
      }
      
      router.push('/dashboard?registered=true')
    } catch (err) {
      console.error('💥 Unexpected error during registration:', err)
      setError(err instanceof Error ? err.message : 'An unexpected error occurred')
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
          <h1 className="text-3xl font-bold bg-gradient-to-r from-neon-cyan to-neon-lavender bg-clip-text text-transparent">
            Join WizardCore
          </h1>
          <p className="text-text-secondary mt-2">
            Start your journey to master programming
          </p>
        </div>

        <div className="bg-bg-elevated border border-border-default rounded-2xl p-8 shadow-2xl">
          <form onSubmit={handleRegister} className="space-y-6">
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
              <p className="text-xs text-text-muted mt-1">
                At least 6 characters
              </p>
            </div>

            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-text-secondary mb-2">
                Confirm Password
              </label>
              <input
                id="confirmPassword"
                type="password"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
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
              className="w-full py-3 px-4 bg-gradient-to-r from-neon-purple to-neon-pink text-white font-semibold rounded-lg hover:opacity-90 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Creating account...' : 'Create Account'}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-text-tertiary text-sm">
              Already have an account?{' '}
              <Link href="/login" className="text-neon-cyan hover:underline font-medium">
                Sign in
              </Link>
            </p>
          </div>

          <div className="mt-8 pt-6 border-t border-border-subtle">
            <div className="space-y-3">
              <p className="text-xs text-text-muted text-center">
                By registering, you agree to our{' '}
                <a href="#" className="text-neon-lavender hover:underline">Terms of Service</a>{' '}
                and{' '}
                <a href="#" className="text-neon-lavender hover:underline">Privacy Policy</a>.
              </p>
              <p className="text-xs text-text-muted text-center">
                You'll receive a confirmation email to verify your account.
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}