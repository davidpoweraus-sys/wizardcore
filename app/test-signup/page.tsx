'use client'

import { useState } from 'react'
import { createBrowserClient } from '@supabase/ssr'

export default function TestSignupPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [status, setStatus] = useState('')
  const [useProxy, setUseProxy] = useState(true)

  const handleTest = async (e: React.FormEvent) => {
    e.preventDefault()
    setStatus('🧪 Testing...')
    
    const supabaseUrl = useProxy 
      ? 'https://offensivewizard.com/api/supabase-proxy'
      : process.env.NEXT_PUBLIC_SUPABASE_URL!

    console.log('🧪 Test Configuration:')
    console.log('  Using Proxy:', useProxy)
    console.log('  Supabase URL:', supabaseUrl)
    console.log('  Email:', email)

    // Create Supabase client with selected URL
    const supabase = createBrowserClient(
      supabaseUrl,
      process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY!
    )

    try {
      console.log('📤 Attempting signup...')
      const { data, error } = await supabase.auth.signUp({
        email,
        password,
      })

      if (error) {
        console.error('❌ Supabase Error:', error)
        setStatus(`❌ Error: ${error.message}`)
      } else {
        console.log('✅ Success!', data)
        setStatus('✅ Success! Check console for details.')
      }
    } catch (err) {
      console.error('💥 Exception:', err)
      setStatus(`💥 Exception: ${err instanceof Error ? err.message : String(err)}`)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-gray-900 to-gray-800">
      <div className="w-full max-w-md p-8 bg-white rounded-lg shadow-xl">
        <h1 className="text-3xl font-bold mb-2 text-gray-900">🧪 Signup Test</h1>
        <p className="text-sm text-gray-600 mb-6">
          Test proxy vs direct connection
        </p>
        
        <div className="mb-6 p-4 bg-blue-50 border border-blue-200 rounded">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-700">Use Proxy:</span>
            <button
              type="button"
              onClick={() => setUseProxy(!useProxy)}
              className={`px-4 py-2 rounded font-medium transition ${
                useProxy 
                  ? 'bg-green-600 text-white' 
                  : 'bg-gray-300 text-gray-700'
              }`}
            >
              {useProxy ? 'ON ✅' : 'OFF'}
            </button>
          </div>
          <p className="text-xs text-gray-600">
            {useProxy 
              ? '✅ Using: /api/supabase-proxy (should work, no CORS)'
              : '❌ Using: auth.offensivewizard.com (will have CORS error)'
            }
          </p>
        </div>

        <form onSubmit={handleTest} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Email
            </label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900"
              placeholder="test@example.com"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Password
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900"
              placeholder="At least 6 characters"
              required
              minLength={6}
            />
          </div>

          <button
            type="submit"
            className="w-full py-3 px-4 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 transition"
          >
            🧪 Test Signup
          </button>
        </form>

        {status && (
          <div className="mt-4 p-4 bg-gray-100 rounded-lg border border-gray-300">
            <strong className="block text-sm text-gray-700 mb-1">Status:</strong>
            <p className="text-sm text-gray-900 whitespace-pre-wrap">{status}</p>
          </div>
        )}

        <div className="mt-6 p-4 bg-yellow-50 border border-yellow-200 rounded">
          <p className="text-xs font-semibold text-yellow-800 mb-2">📋 Instructions:</p>
          <ol className="text-xs text-yellow-700 space-y-1 list-decimal list-inside">
            <li>Open browser console (F12)</li>
            <li>Toggle "Use Proxy" ON/OFF to compare</li>
            <li>Try signup with both settings</li>
            <li>Check console for detailed logs</li>
            <li>Proxy should work, direct should fail with CORS</li>
          </ol>
        </div>

        <div className="mt-4 text-center">
          <a 
            href="/register" 
            className="text-sm text-blue-600 hover:underline"
          >
            ← Back to real registration
          </a>
        </div>
      </div>
    </div>
  )
}
