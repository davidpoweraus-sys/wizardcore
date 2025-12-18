'use client'

import { useState } from 'react'
import { submitCode, LANGUAGES } from '@/lib/judge0/service'

interface Judge0IntegrationProps {
  initialCode?: string
  languageId?: number
  onResult?: (result: any) => void
}

export default function Judge0Integration({
  initialCode = '',
  languageId = 71, // Python
  onResult,
}: Judge0IntegrationProps) {
  const [code, setCode] = useState(initialCode)
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<any>(null)
  const [error, setError] = useState<string | null>(null)

  const handleRun = async () => {
    setLoading(true)
    setError(null)
    try {
      const submission = {
        source_code: code,
        language_id: languageId,
        stdin: '',
      }
      const data = await submitCode(submission)
      setResult(data)
      if (onResult) onResult(data)
    } catch (err: any) {
      setError(err.message || 'Failed to execute code')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <select
          className="px-3 py-2 bg-bg-tertiary border border-border-subtle rounded-lg text-text-primary"
          value={languageId}
          onChange={(e) => console.log('Language changed', e.target.value)}
        >
          {LANGUAGES.map((lang) => (
            <option key={lang.id} value={lang.id}>
              {lang.name}
            </option>
          ))}
        </select>
        <button
          onClick={handleRun}
          disabled={loading}
          className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg font-medium hover:opacity-90 disabled:opacity-50"
        >
          {loading ? 'Running...' : 'Run with Judge0'}
        </button>
      </div>

      {error && (
        <div className="p-3 bg-red-900/30 border border-red-700 rounded-lg text-red-300">
          {error}
        </div>
      )}

      {result && (
        <div className="p-4 bg-black/50 border border-border-subtle rounded-lg">
          <h4 className="font-medium text-text-primary mb-2">Execution Result</h4>
          <pre className="text-sm font-mono text-green-400 overflow-x-auto">
            {result.stdout || result.stderr || result.compile_output || 'No output'}
          </pre>
          <div className="mt-2 text-xs text-text-secondary">
            Status: {result.status?.description}
          </div>
        </div>
      )}
    </div>
  )
}