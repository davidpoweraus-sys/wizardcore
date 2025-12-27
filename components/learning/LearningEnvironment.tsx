'use client'

import { useState, useEffect } from 'react'
import { Menu, BookOpen, Code2, CheckCircle2, Lightbulb, ChevronRight, Play, Send, Target, Users, Clock, Loader2, Save } from 'lucide-react'
import { submitCode } from '@/lib/judge0/service'
import { api } from '@/lib/api'

interface ExerciseWithTests {
  exercise: {
    id: string
    module_id: string
    title: string
    difficulty: string
    points: number
    time_limit_minutes?: number
    sort_order: number
    objectives: string[]
    content?: string
    examples?: any[]
    description?: string
    constraints: string[]
    hints: string[]
    starter_code?: string
    solution_code?: string
    language_id: number
    tags: string[]
    concurrent_solvers: number
    total_submissions: number
    total_completions: number
    average_completion_time?: number
    created_at: string
    updated_at: string
  }
  test_cases: Array<{
    id: string
    exercise_id: string
    input?: string
    expected_output: string
    is_hidden: boolean
    points: number
    sort_order: number
    created_at: string
  }>
}

interface Submission {
  id: string
  user_id: string
  exercise_id: string
  source_code: string
  language_id: number
  status: string
  test_cases_passed: number
  test_cases_total: number
  points_earned: number
  is_correct: boolean
  created_at: string
  updated_at: string
}

interface LearningEnvironmentProps {
  exerciseId: string
}

export default function LearningEnvironment({ exerciseId }: LearningEnvironmentProps) {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [activeTab, setActiveTab] = useState<'lesson' | 'code' | 'tests' | 'hints'>('lesson')
  const [code, setCode] = useState('')
  const [exerciseData, setExerciseData] = useState<ExerciseWithTests | null>(null)
  const [, setLatestSubmission] = useState<Submission | null>(null)
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [output, setOutput] = useState<string>('')
  const [isRunning, setIsRunning] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Fetch exercise data and latest submission
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)
        setError(null)
        
        // Fetch exercise data
        const exerciseResponse = await api.get<ExerciseWithTests>(`/exercises/${exerciseId}`)
        setExerciseData(exerciseResponse)
        
        // Set starter code from exercise or latest submission
        if (exerciseResponse.exercise.starter_code) {
          setCode(exerciseResponse.exercise.starter_code)
        }
        
        // Fetch latest submission
        try {
          const submissionResponse = await api.get<{ submission: Submission }>(`/submissions/latest/${exerciseId}`)
          if (submissionResponse.submission && submissionResponse.submission.source_code) {
            setCode(submissionResponse.submission.source_code)
            setLatestSubmission(submissionResponse.submission)
          }
        } catch (submissionError) {
          // It's okay if no submission exists yet
          console.log('No previous submission found')
        }
      } catch (err: any) {
        console.error('Failed to fetch exercise data:', err)
        setError(err.message || 'Failed to load exercise')
      } finally {
        setLoading(false)
      }
    }
    
    if (exerciseId) {
      fetchData()
    }
  }, [exerciseId])

  // Auto-save functionality
  useEffect(() => {
    if (!exerciseId || !code || loading) return
    
    const autoSave = async () => {
      try {
        setSaving(true)
        await api.post(`/submissions/save-draft/${exerciseId}`, {
          source_code: code,
          language_id: exerciseData?.exercise.language_id || 71 // Default to Python
        })
      } catch (err) {
        console.error('Auto-save failed:', err)
      } finally {
        setSaving(false)
      }
    }
    
    const timer = setTimeout(autoSave, 30000) // Auto-save every 30 seconds
    
    return () => clearTimeout(timer)
  }, [code, exerciseId, exerciseData, loading])

  const handleRunCode = async () => {
    if (!code.trim()) {
      setOutput('> Error: No code to run')
      return
    }
    
    setIsRunning(true)
    setOutput('> Running code...\n')
    try {
      const result = await submitCode({
        source_code: code,
        language_id: exerciseData?.exercise.language_id || 71, // Default to Python
        stdin: '',
        expected_output: '',
      })
      const outputText = result.stdout || result.stderr || result.compile_output || result.message || 'No output'
      setOutput(`> Execution completed (${result.status.description})\n${outputText}`)
    } catch (error) {
      setOutput(`> Error: ${error instanceof Error ? error.message : 'Unknown error'}`)
    } finally {
      setIsRunning(false)
    }
  }

  const handleSubmit = async () => {
    if (!code.trim()) {
      setOutput('> Error: No code to submit')
      return
    }
    
    if (!exerciseData) {
      setOutput('> Error: Exercise data not loaded')
      return
    }
    
    setSubmitting(true)
    setOutput('> Submitting solution...\n')
    
    try {
      const response = await api.post<{ submission: Submission }>('/submissions', {
        exercise_id: exerciseId,
        source_code: code,
        language_id: exerciseData.exercise.language_id,
      })
      
      const submission = response.submission
      setLatestSubmission(submission)
      
      if (submission.is_correct) {
        setOutput(`> ✅ Submission accepted! You earned ${submission.points_earned} XP\n> Test cases passed: ${submission.test_cases_passed}/${submission.test_cases_total}`)
      } else {
        setOutput(`> ❌ Submission rejected\n> Test cases passed: ${submission.test_cases_passed}/${submission.test_cases_total}\n> Status: ${submission.status}`)
      }
    } catch (err: any) {
      setOutput(`> Error submitting solution: ${err.message || 'Unknown error'}`)
    } finally {
      setSubmitting(false)
    }
  }

  const handleSaveDraft = async () => {
    if (!code.trim() || !exerciseData) return
    
    try {
      setSaving(true)
      await api.post(`/submissions/save-draft/${exerciseId}`, {
        source_code: code,
        language_id: exerciseData.exercise.language_id,
      })
      setOutput('> Draft saved successfully')
    } catch (err: any) {
      setOutput(`> Error saving draft: ${err.message || 'Unknown error'}`)
    } finally {
      setSaving(false)
    }
  }

  const handleResetCode = () => {
    if (exerciseData?.exercise.starter_code) {
      setCode(exerciseData.exercise.starter_code)
      setOutput('> Code reset to starter template')
    }
  }

  if (loading) {
    return (
      <div className="h-full flex items-center justify-center bg-bg-primary">
        <div className="text-center">
          <Loader2 className="w-8 h-8 animate-spin text-neon-cyan mx-auto mb-4" />
          <p className="text-text-secondary">Loading exercise...</p>
        </div>
      </div>
    )
  }

  if (error || !exerciseData) {
    return (
      <div className="h-full flex items-center justify-center bg-bg-primary">
        <div className="text-center">
          <p className="text-red-400 mb-2">Error loading exercise</p>
          <p className="text-text-secondary">{error || 'Exercise not found'}</p>
        </div>
      </div>
    )
  }

  const exercise = exerciseData.exercise
  const testCases = exerciseData.test_cases.filter(tc => !tc.is_hidden) // Show only non-hidden test cases

  // Format examples from exercise data
  const examples = exercise.examples ? (Array.isArray(exercise.examples) ? exercise.examples : []) : []

  return (
    <div className="h-full flex flex-col bg-bg-primary">
      {/* Header */}
      <header className="h-16 border-b border-border-default px-6 flex items-center justify-between bg-bg-elevated">
        <div className="flex items-center gap-4">
          <button
            onClick={() => setSidebarOpen(!sidebarOpen)}
            className="p-2 rounded-lg hover:bg-bg-tertiary text-text-secondary"
          >
            <Menu className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-lg font-bold text-text-primary">{exercise.title}</h1>
            <p className="text-sm text-text-secondary">Module: {exercise.module_id}</p>
          </div>
        </div>
        
        <div className="flex items-center gap-6">
          <div className="hidden md:flex items-center gap-4">
            {exercise.time_limit_minutes && (
              <div className="text-sm text-text-secondary">
                <Clock className="w-4 h-4 inline mr-2" />
                {exercise.time_limit_minutes} min remaining
              </div>
            )}
            <div className="flex items-center gap-2">
              <Users className="w-4 h-4 text-text-muted" />
              <span className="text-sm text-text-secondary">{exercise.concurrent_solvers} solving</span>
            </div>
          </div>
          <div className="flex items-center gap-2">
            <div className="px-3 py-1 rounded-full bg-gradient-to-r from-green-900/30 to-green-800/30 border border-green-700 text-green-300 text-sm font-medium">
              {exercise.difficulty}
            </div>
            <div className="px-3 py-1 rounded-full bg-bg-tertiary text-text-secondary text-sm">
              {exercise.points} XP
            </div>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Collapsible Sidebar */}
        <aside
          className={`border-r border-border-default bg-bg-elevated transition-all duration-300 overflow-y-auto ${sidebarOpen ? 'w-96' : 'w-0'}`}
        >
          {sidebarOpen && (
            <div className="p-6">
              <div className="flex items-center justify-between mb-6">
                <h2 className="text-lg font-bold text-text-primary flex items-center gap-2">
                  <BookOpen className="w-5 h-5" />
                  Lesson Content
                </h2>
                <button className="text-sm text-neon-cyan hover:underline">
                  Pin
                </button>
              </div>

              {/* Learning Objectives */}
              {exercise.objectives && exercise.objectives.length > 0 && (
                <div className="mb-8 p-4 rounded-xl bg-gradient-to-r from-neon-cyan/10 to-neon-lavender/10 border border-neon-cyan/20">
                  <h3 className="font-semibold text-text-primary mb-3 flex items-center gap-2">
                    <Target className="w-4 h-4" />
                    Learning Objectives
                  </h3>
                  <ul className="space-y-2">
                    {exercise.objectives.map((obj, idx) => (
                      <li key={idx} className="text-sm text-text-secondary flex items-start gap-2">
                        <ChevronRight className="w-4 h-4 text-neon-cyan mt-0.5 flex-shrink-0" />
                        {obj}
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              {/* Lesson Content */}
              {exercise.content && (
                <div className="prose prose-sm max-w-none text-text-primary">
                  <h3 className="text-xl font-bold mb-4">{exercise.title}</h3>
                  <div className="whitespace-pre-line text-text-secondary mb-6">
                    {exercise.content}
                  </div>
                </div>
              )}

              {/* Examples */}
              {examples.length > 0 && (
                <div className="mb-8">
                  <h3 className="font-semibold text-text-primary mb-3 flex items-center gap-2">
                    <Code2 className="w-4 h-4" />
                    Examples
                  </h3>
                  <div className="space-y-4">
                    {examples.map((example: any, idx) => (
                      <div key={idx} className="bg-bg-tertiary border border-border-subtle rounded-lg p-4">
                        <pre className="text-sm font-mono text-green-400 overflow-x-auto">
                          {example.code || JSON.stringify(example, null, 2)}
                        </pre>
                        {example.output && (
                          <div className="mt-3">
                            <p className="text-xs text-text-muted mb-1">Output:</p>
                            <pre className="text-sm font-mono bg-black/30 p-2 rounded">
                              {example.output}
                            </pre>
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              {/* Exercise Instructions */}
              <div className="border-t border-border-subtle pt-6">
                <h3 className="font-bold text-text-primary mb-3">Exercise Instructions</h3>
                <p className="text-text-secondary mb-4">{exercise.description || 'Complete the exercise as described.'}</p>
                
                {exercise.constraints && exercise.constraints.length > 0 && (
                  <div className="mb-4">
                    <h4 className="text-sm font-semibold text-text-primary mb-2">Constraints:</h4>
                    <ul className="text-sm text-text-secondary space-y-1">
                      {exercise.constraints.map((constraint, idx) => (
                        <li key={idx} className="flex items-start gap-2">
                          <span className="text-neon-cyan">•</span>
                          {constraint}
                        </li>
                      ))}
                    </ul>
                  </div>
                )}

                {testCases.length > 0 && (
                  <div>
                    <h4 className="text-sm font-semibold text-text-primary mb-2">Example Test Cases:</h4>
                    {testCases.slice(0, 2).map((test, idx) => (
                      <div key={idx} className="text-sm mb-3 bg-bg-tertiary p-3 rounded border border-border-subtle">
                        <div className="font-mono">
                          <div><span className="text-text-muted">Input:</span> {test.input || '(none)'}</div>
                          <div><span className="text-text-muted">Expected Output:</span> {test.expected_output}</div>
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            </div>
          )}
        </aside>

        {/* Main Editor Area */}
        <main className="flex-1 flex flex-col overflow-hidden">
          {/* Tabs */}
          <div className="border-b border-border-default">
            <div className="flex">
              {(['lesson', 'code', 'tests', 'hints'] as const).map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab)}
                  className={`px-6 py-3 text-sm font-medium border-b-2 transition ${activeTab === tab
                      ? 'border-neon-cyan text-neon-cyan'
                      : 'border-transparent text-text-secondary hover:text-text-primary'
                    }`}
                >
                  {tab === 'lesson' && <BookOpen className="w-4 h-4 inline mr-2" />}
                  {tab === 'code' && <Code2 className="w-4 h-4 inline mr-2" />}
                  {tab === 'tests' && <CheckCircle2 className="w-4 h-4 inline mr-2" />}
                  {tab === 'hints' && <Lightbulb className="w-4 h-4 inline mr-2" />}
                  {tab.charAt(0).toUpperCase() + tab.slice(1)}
                </button>
              ))}
            </div>
          </div>

          {/* Tab Content */}
          <div className="flex-1 overflow-hidden">
            {activeTab === 'lesson' && (
              <div className="h-full overflow-y-auto p-6">
                <div className="max-w-3xl">
                  <h2 className="text-2xl font-bold text-text-primary mb-4">Lesson Overview</h2>
                  <div className="text-text-secondary space-y-4">
                    {exercise.content ? (
                      <div className="prose prose-invert max-w-none">
                        <div dangerouslySetInnerHTML={{ __html: exercise.content.replace(/\n/g, '<br/>') }} />
                      </div>
                    ) : (
                      <p>No lesson content available for this exercise.</p>
                    )}
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'code' && (
              <div className="h-full flex flex-col">
                {/* Editor */}
                <div className="flex-1 p-4">
                  <div className="h-full border border-border-subtle rounded-lg overflow-hidden">
                    <div className="bg-bg-tertiary px-4 py-2 border-b border-border-subtle text-sm text-text-secondary font-mono flex justify-between items-center">
                      <span>exercise.py</span>
                      <div className="flex items-center gap-2">
                        {saving && (
                          <span className="text-xs text-text-muted flex items-center gap-1">
                            <Loader2 className="w-3 h-3 animate-spin" />
                            Saving...
                          </span>
                        )}
                      </div>
                    </div>
                    <textarea
                      value={code}
                      onChange={(e) => setCode(e.target.value)}
                      className="w-full h-full bg-bg-primary text-text-primary font-mono text-sm p-4 resize-none focus:outline-none"
                      spellCheck="false"
                      placeholder="Write your code here..."
                    />
                  </div>
                </div>
                
                {/* Action Buttons */}
                <div className="border-t border-border-subtle p-4 flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <button
                      onClick={handleRunCode}
                      disabled={isRunning}
                      className="px-4 py-2 bg-neon-cyan text-black font-medium rounded-lg hover:bg-neon-cyan/90 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                    >
                      <Play className="w-4 h-4" />
                      {isRunning ? 'Running...' : 'Run Code'}
                    </button>
                    <button
                      onClick={handleSubmit}
                      disabled={submitting}
                      className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-black font-medium rounded-lg hover:opacity-90 disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
                    >
                      <Send className="w-4 h-4" />
                      {submitting ? 'Submitting...' : 'Submit Solution'}
                    </button>
                  </div>
                  
                  <div className="flex items-center gap-2">
                    <button
                      onClick={handleSaveDraft}
                      disabled={saving}
                      className="px-3 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary disabled:opacity-50 flex items-center gap-2"
                    >
                      <Save className="w-4 h-4" />
                      {saving ? 'Saving...' : 'Save Draft'}
                    </button>
                    <button
                      onClick={handleResetCode}
                      className="px-3 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary"
                    >
                      Reset
                    </button>
                  </div>
                </div>
                
                {/* Output Panel */}
                <div className="border-t border-border-subtle">
                  <div className="p-4">
                    <h3 className="text-sm font-semibold text-text-primary mb-2">Output</h3>
                    <pre className="bg-bg-tertiary border border-border-subtle rounded-lg p-4 font-mono text-sm text-text-primary whitespace-pre-wrap min-h-[100px] max-h-[300px] overflow-y-auto">
                      {output || '> Output will appear here...'}
                    </pre>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'tests' && (
              <div className="h-full overflow-y-auto p-6">
                <h2 className="text-2xl font-bold text-text-primary mb-6">Test Cases</h2>
                
                {testCases.length > 0 ? (
                  <div className="space-y-4">
                    {testCases.map((test, idx) => (
                      <div key={test.id} className="border border-border-subtle rounded-lg p-4 bg-bg-tertiary">
                        <div className="flex items-center justify-between mb-3">
                          <h3 className="font-semibold text-text-primary">Test Case {idx + 1}</h3>
                          <div className="text-sm text-text-secondary">
                            {test.points} points
                          </div>
                        </div>
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                          <div>
                            <h4 className="text-sm font-medium text-text-muted mb-1">Input</h4>
                            <pre className="bg-bg-primary p-3 rounded border border-border-default text-sm font-mono text-text-primary overflow-x-auto">
                              {test.input || '(none)'}
                            </pre>
                          </div>
                          <div>
                            <h4 className="text-sm font-medium text-text-muted mb-1">Expected Output</h4>
                            <pre className="bg-bg-primary p-3 rounded border border-border-default text-sm font-mono text-green-400 overflow-x-auto">
                              {test.expected_output}
                            </pre>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <CheckCircle2 className="w-12 h-12 text-text-muted mx-auto mb-4" />
                    <p className="text-text-secondary">No test cases available for this exercise.</p>
                  </div>
                )}
              </div>
            )}

            {activeTab === 'hints' && (
              <div className="h-full overflow-y-auto p-6">
                <h2 className="text-2xl font-bold text-text-primary mb-6">Hints</h2>
                
                {exercise.hints && exercise.hints.length > 0 ? (
                  <div className="space-y-6">
                    {exercise.hints.map((hint, idx) => (
                      <div key={idx} className="border border-neon-lavender/30 rounded-lg p-5 bg-gradient-to-r from-neon-lavender/10 to-transparent">
                        <div className="flex items-start gap-3">
                          <Lightbulb className="w-5 h-5 text-neon-lavender mt-0.5 flex-shrink-0" />
                          <div>
                            <h3 className="font-semibold text-text-primary mb-2">Hint #{idx + 1}</h3>
                            <p className="text-text-secondary">{hint}</p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="text-center py-12">
                    <Lightbulb className="w-12 h-12 text-text-muted mx-auto mb-4" />
                    <p className="text-text-secondary">No hints available for this exercise.</p>
                    <p className="text-sm text-text-muted mt-2">Try working through the problem step by step!</p>
                  </div>
                )}
                
                {/* Additional Resources */}
                <div className="mt-8 pt-6 border-t border-border-subtle">
                  <h3 className="font-semibold text-text-primary mb-4">Additional Resources</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <a href="#" className="block p-4 border border-border-subtle rounded-lg hover:bg-bg-tertiary transition">
                      <h4 className="font-medium text-text-primary mb-2">Python Documentation</h4>
                      <p className="text-sm text-text-secondary">Official Python language reference</p>
                    </a>
                    <a href="#" className="block p-4 border border-border-subtle rounded-lg hover:bg-bg-tertiary transition">
                      <h4 className="font-medium text-text-primary mb-2">Stack Overflow</h4>
                      <p className="text-sm text-text-secondary">Community Q&A for programming questions</p>
                    </a>
                  </div>
                </div>
              </div>
            )}
          </div>
        </main>
      </div>
    </div>
  )
}