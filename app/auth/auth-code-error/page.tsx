import Link from 'next/link'

export default function AuthCodeError() {
  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="bg-bg-elevated border border-border-default rounded-2xl p-8 shadow-2xl">
          <div className="text-center">
            <h1 className="text-2xl font-bold text-red-400 mb-4">
              Authentication Error
            </h1>
            <p className="text-text-secondary mb-6">
              Sorry, we couldn't authenticate your session. This could be due to an expired or invalid link.
            </p>
            <div className="space-y-3">
              <Link
                href="/register"
                className="block w-full py-3 px-4 bg-gradient-to-r from-neon-purple to-neon-pink text-white font-semibold rounded-lg hover:opacity-90 transition text-center"
              >
                Try Signing Up Again
              </Link>
              <Link
                href="/login"
                className="block w-full py-3 px-4 border border-border-default text-text-primary font-semibold rounded-lg hover:bg-bg-tertiary transition text-center"
              >
                Back to Login
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
