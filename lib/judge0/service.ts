/**
 * Judge0 integration service for code execution.
 * Assumes a self‑hosted Judge0 instance running at JUDGE0_API_URL.
 */

const JUDGE0_API_URL = process.env.NEXT_PUBLIC_JUDGE0_API_URL || 'http://localhost:2358'
const JUDGE0_API_KEY = process.env.JUDGE0_API_KEY || ''

export interface Submission {
  source_code: string
  language_id: number
  stdin?: string
  expected_output?: string
}

export interface SubmissionResult {
  stdout: string | null
  stderr: string | null
  compile_output: string | null
  message: string | null
  status: {
    id: number
    description: string
  }
  time: string
  memory: number
}

export interface Language {
  id: number
  name: string
}

// Common language IDs
export const LANGUAGES: Language[] = [
  { id: 71, name: 'Python (3.8.1)' },
  { id: 50, name: 'C (GCC 9.2.0)' },
  { id: 54, name: 'C++ (GCC 9.2.0)' },
  { id: 62, name: 'Java (OpenJDK 13.0.1)' },
  { id: 63, name: 'JavaScript (Node.js 12.14.0)' },
  { id: 82, name: 'SQL (SQLite 3.27.2)' },
]

/**
 * Submit code to Judge0 for execution.
 */
export async function submitCode(submission: Submission): Promise<SubmissionResult> {
  const url = `${JUDGE0_API_URL}/submissions?base64_encoded=false&wait=true`

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-RapidAPI-Key': JUDGE0_API_KEY,
    },
    body: JSON.stringify(submission),
  })

  if (!response.ok) {
    throw new Error(`Judge0 submission failed: ${response.statusText}`)
  }

  const data = await response.json()
  return data
}

/**
 * Create a batch submission for multiple test cases.
 */
export async function batchSubmit(
  source_code: string,
  language_id: number,
  test_cases: Array<{ stdin?: string; expected_output?: string }>
): Promise<SubmissionResult[]> {
  const submissions = test_cases.map((tc) => ({
    source_code,
    language_id,
    stdin: tc.stdin || '',
    expected_output: tc.expected_output || '',
  }))

  const url = `${JUDGE0_API_URL}/submissions/batch?base64_encoded=false`

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'X-RapidAPI-Key': JUDGE0_API_KEY,
    },
    body: JSON.stringify({ submissions }),
  })

  if (!response.ok) {
    throw new Error(`Judge0 batch submission failed: ${response.statusText}`)
  }

  const data = await response.json()
  return data
}

/**
 * Get submission result by token.
 */
export async function getSubmission(token: string): Promise<SubmissionResult> {
  const url = `${JUDGE0_API_URL}/submissions/${token}?base64_encoded=false`

  const response = await fetch(url, {
    headers: {
      'X-RapidAPI-Key': JUDGE0_API_KEY,
    },
  })

  if (!response.ok) {
    throw new Error(`Failed to fetch submission: ${response.statusText}`)
  }

  return response.json()
}

/**
 * Check if Judge0 is reachable.
 */
export async function healthCheck(): Promise<boolean> {
  try {
    const controller = new AbortController()
    const timeout = setTimeout(() => controller.abort(), 5000)
    const response = await fetch(`${JUDGE0_API_URL}/about`, { signal: controller.signal })
    clearTimeout(timeout)
    return response.ok
  } catch {
    return false
  }
}