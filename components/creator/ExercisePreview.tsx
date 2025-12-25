'use client'

import { useState } from 'react'
import Editor from '@monaco-editor/react'
import { X, Play, Send, Eye, EyeOff, Lightbulb } from 'lucide-react'
import { submitCode, LANGUAGES } from '@/lib/judge0/service'

interface ExercisePreviewProps {
  exerciseData: any
  onClose: () => void
  isOpen: boolean
}

export default function ExercisePreview({
  exerciseData,
  onClose,
  isOpen,
}: ExercisePreviewProps) {
  const [code, setCode] = useState(exerciseData.starter_code || '')
  const [isRunning, setIsRunning] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [testOutput, setTestOutput] = useState<any>(null)
  const [submissionResult, setSubmissionResult] = useState<any>(null)
  const [showHints, setShowHints] = useState(0)
  const [activeTab, setActiveTab] = useState<'description' | 'editor' | 'results'>('description')

  if (!isOpen) return null

  const language = LANGUAGES.find((lang) => lang.id === exerciseData.language_id)
  const monacoLang = language?.name.toLowerCase().split(' ')[0] || 'python'
  const visibleTestCases = exerciseData.test_cases?.filter((tc: any) => !tc.is_hidden) || []

  const handleRun = async () => {
    setIsRunning(true)
    setTestOutput(null)

    try {
      const result = await submitCode({
        source_code: code,
        language_id: exerciseData.language_id,
        stdin: '',
      })

      setTestOutput({
        stdout: result.stdout,
        stderr: result.stderr,
        time: result.time ? parseFloat(result.time) * 1000 : null,
        status: result.status.description,
      })
      setActiveTab('results')
    } catch (error: any) {
      setTestOutput({ error: error.message })
      setActiveTab('results')
    } finally {
      setIsRunning(false)
    }
  }

  const handleSubmit = async () => {
    if (!exerciseData.test_cases || exerciseData.test_cases.length === 0) {
      alert('No test cases configured for this exercise')
      return
    }

    setIsSubmitting(true)
    setSubmissionResult(null)

    try {
      const results = await Promise.all(
        exerciseData.test_cases.map(async (testCase: any) => {
          try {
            const result = await submitCode({
              source_code: code,
              language_id: exerciseData.language_id,
              stdin: testCase.input || '',
            })

            const passed = result.stdout?.trim() === testCase.expected_output.trim()

            return {
              testCase,
              passed,
              actual_output: result.stdout || '',
              error: result.stderr || result.compile_output || null,
            }
          } catch (error: any) {
            return {
              testCase,
              passed: false,
              error: error.message,
            }
          }
        })
      )

      const passedCount = results.filter((r) => r.passed).length
      const totalCount = results.length
      const visibleResults = results.filter((r) => !r.testCase.is_hidden)

      const score = Math.round(
        (passedCount / totalCount) * (exerciseData.points || 100)
      )

      setSubmissionResult({
        results,
        visibleResults,
        passedCount,
        totalCount,
        allPassed: passedCount === totalCount,
        score,
      })
      setActiveTab('results')
    } catch (error: any) {
      alert(`Submission failed: ${error.message}`)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/70 backdrop-blur-sm">
      <div className="w-full max-w-7xl h-[90vh] bg-bg-elevated border border-border-default rounded-2xl shadow-2xl flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b border-border-default">
          <div>
            <div className="flex items-center gap-3">
              <Eye className="w-5 h-5 text-accent-primary" />
              <h2 className="text-xl font-bold text-text-primary">Exercise Preview</h2>
            </div>
            <p className="text-sm text-text-secondary mt-1">
              This is how students will see your exercise
            </p>
          </div>
          <button
            onClick={onClose}
            className="p-2 hover:bg-bg-hover rounded-lg transition-colors"
          >
            <X className="w-5 h-5 text-text-secondary" />
          </button>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-hidden flex">
          {/* Left Panel - Exercise Description */}
          <div className="w-1/2 border-r border-border-default overflow-y-auto p-6 space-y-6">
            {/* Title & Metadata */}
            <div>
              <div className="flex items-center gap-3 mb-2">
                <span
                  className={`px-3 py-1 rounded-full text-xs font-medium ${
                    exerciseData.difficulty === 'BEGINNER'
                      ? 'bg-green-500/20 text-green-400'
                      : exerciseData.difficulty === 'INTERMEDIATE'
                      ? 'bg-yellow-500/20 text-yellow-400'
                      : 'bg-red-500/20 text-red-400'
                  }`}
                >
                  {exerciseData.difficulty}
                </span>
                <span className="text-sm text-text-tertiary">
                  {exerciseData.points || 100} points
                </span>
                {exerciseData.time_limit_minutes && (
                  <span className="text-sm text-text-tertiary">
                    {exerciseData.time_limit_minutes} min
                  </span>
                )}
              </div>
              <h1 className="text-2xl font-bold text-text-primary">
                {exerciseData.title || 'Untitled Exercise'}
              </h1>
              {exerciseData.description && (
                <p className="text-text-secondary mt-2">{exerciseData.description}</p>
              )}
            </div>

            {/* Objectives */}
            {exerciseData.objectives && exerciseData.objectives.filter((o: string) => o.trim()).length > 0 && (
              <div>
                <h3 className="text-lg font-semibold text-text-primary mb-3">
                  Learning Objectives
                </h3>
                <ul className="space-y-2">
                  {exerciseData.objectives
                    .filter((o: string) => o.trim())
                    .map((objective: string, index: number) => (
                      <li key={index} className="flex items-start gap-2 text-text-secondary">
                        <span className="text-accent-primary mt-1">•</span>
                        <span>{objective}</span>
                      </li>
                    ))}
                </ul>
              </div>
            )}

            {/* Content (Markdown) */}
            {exerciseData.content && (
              <div>
                <h3 className="text-lg font-semibold text-text-primary mb-3">
                  Problem Description
                </h3>
                <div className="prose prose-invert max-w-none">
                  <pre className="whitespace-pre-wrap text-text-secondary">
                    {exerciseData.content}
                  </pre>
                </div>
              </div>
            )}

            {/* Constraints */}
            {exerciseData.constraints && exerciseData.constraints.filter((c: string) => c.trim()).length > 0 && (
              <div>
                <h3 className="text-lg font-semibold text-text-primary mb-3">
                  Constraints
                </h3>
                <ul className="space-y-2">
                  {exerciseData.constraints
                    .filter((c: string) => c.trim())
                    .map((constraint: string, index: number) => (
                      <li key={index} className="flex items-start gap-2 text-text-secondary">
                        <span className="text-accent-primary mt-1">•</span>
                        <span>{constraint}</span>
                      </li>
                    ))}
                </ul>
              </div>
            )}

            {/* Test Cases (Visible Only) */}
            {visibleTestCases.length > 0 && (
              <div>
                <h3 className="text-lg font-semibold text-text-primary mb-3">
                  Example Test Cases
                </h3>
                <div className="space-y-3">
                  {visibleTestCases.map((testCase: any, index: number) => (
                    <div
                      key={index}
                      className="p-3 bg-bg-primary border border-border-default rounded-lg"
                    >
                      <div className="flex items-center gap-2 mb-2">
                        <Eye className="w-4 h-4 text-accent-primary" />
                        <span className="text-sm font-medium text-text-primary">
                          Test Case {index + 1}
                        </span>
                        <span className="text-xs text-text-tertiary">
                          ({testCase.points} pts)
                        </span>
                      </div>
                      {testCase.input && (
                        <div className="mb-2">
                          <span className="text-xs font-medium text-text-secondary">Input:</span>
                          <pre className="text-xs bg-bg-elevated p-2 rounded mt-1 overflow-x-auto">
                            {testCase.input}
                          </pre>
                        </div>
                      )}
                      <div>
                        <span className="text-xs font-medium text-text-secondary">
                          Expected Output:
                        </span>
                        <pre className="text-xs bg-bg-elevated p-2 rounded mt-1 overflow-x-auto">
                          {testCase.expected_output}
                        </pre>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Hints */}
            {exerciseData.hints && exerciseData.hints.filter((h: string) => h.trim()).length > 0 && (
              <div>
                <h3 className="text-lg font-semibold text-text-primary mb-3 flex items-center gap-2">
                  <Lightbulb className="w-5 h-5" />
                  Hints
                </h3>
                <div className="space-y-2">
                  {exerciseData.hints
                    .filter((h: string) => h.trim())
                    .slice(0, showHints)
                    .map((hint: string, index: number) => (
                      <div
                        key={index}
                        className="p-3 bg-yellow-500/10 border border-yellow-500/30 rounded-lg"
                      >
                        <div className="flex items-start gap-2">
                          <span className="text-yellow-400 font-medium">#{index + 1}</span>
                          <span className="text-text-secondary">{hint}</span>
                        </div>
                      </div>
                    ))}
                  {showHints < exerciseData.hints.filter((h: string) => h.trim()).length && (
                    <button
                      onClick={() => setShowHints(showHints + 1)}
                      className="px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-sm"
                    >
                      Show Next Hint ({showHints + 1}/
                      {exerciseData.hints.filter((h: string) => h.trim()).length})
                    </button>
                  )}
                </div>
              </div>
            )}
          </div>

          {/* Right Panel - Code Editor & Results */}
          <div className="w-1/2 flex flex-col">
            {/* Tabs */}
            <div className="flex border-b border-border-default">
              <button
                onClick={() => setActiveTab('editor')}
                className={`px-4 py-3 text-sm font-medium transition-colors ${
                  activeTab === 'editor'
                    ? 'text-accent-primary border-b-2 border-accent-primary'
                    : 'text-text-secondary hover:text-text-primary'
                }`}
              >
                Code Editor
              </button>
              <button
                onClick={() => setActiveTab('results')}
                className={`px-4 py-3 text-sm font-medium transition-colors ${
                  activeTab === 'results'
                    ? 'text-accent-primary border-b-2 border-accent-primary'
                    : 'text-text-secondary hover:text-text-primary'
                }`}
              >
                Results
                {submissionResult && (
                  <span className="ml-2 px-2 py-0.5 rounded-full text-xs bg-accent-primary/20 text-accent-primary">
                    {submissionResult.passedCount}/{submissionResult.totalCount}
                  </span>
                )}
              </button>
            </div>

            {/* Editor Tab */}
            {activeTab === 'editor' && (
              <div className="flex-1 flex flex-col">
                <Editor
                  height="100%"
                  language={monacoLang}
                  value={code}
                  onChange={(value) => setCode(value || '')}
                  theme="vs-dark"
                  options={{
                    minimap: { enabled: false },
                    fontSize: 14,
                    scrollBeyondLastLine: false,
                  }}
                />
                <div className="p-4 border-t border-border-default space-y-2">
                  <div className="flex gap-2">
                    <button
                      onClick={handleRun}
                      disabled={isRunning}
                      className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors disabled:opacity-50"
                    >
                      <Play className="w-4 h-4" />
                      {isRunning ? 'Running...' : 'Test Run'}
                    </button>
                    <button
                      onClick={handleSubmit}
                      disabled={isSubmitting}
                      className="flex-1 flex items-center justify-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors disabled:opacity-50"
                    >
                      <Send className="w-4 h-4" />
                      {isSubmitting ? 'Submitting...' : 'Submit Solution'}
                    </button>
                  </div>
                  <p className="text-xs text-text-tertiary text-center">
                    Test Run executes your code. Submit Solution grades it against test cases.
                  </p>
                </div>
              </div>
            )}

            {/* Results Tab */}
            {activeTab === 'results' && (
              <div className="flex-1 overflow-y-auto p-6 space-y-4">
                {/* Test Output */}
                {testOutput && (
                  <div>
                    <h3 className="text-lg font-semibold text-text-primary mb-3">
                      Test Run Output
                    </h3>
                    <div className="bg-bg-primary border border-border-default rounded-lg p-4">
                      {testOutput.error ? (
                        <pre className="text-red-400 text-sm font-mono whitespace-pre-wrap">
                          {testOutput.error}
                        </pre>
                      ) : (
                        <>
                          <pre className="text-green-400 text-sm font-mono whitespace-pre-wrap">
                            {testOutput.stdout || '(no output)'}
                          </pre>
                          {testOutput.stderr && (
                            <pre className="text-red-400 text-sm font-mono whitespace-pre-wrap mt-2">
                              {testOutput.stderr}
                            </pre>
                          )}
                          <div className="mt-2 text-xs text-text-tertiary">
                            Status: {testOutput.status}
                            {testOutput.time && ` • ${testOutput.time.toFixed(2)}ms`}
                          </div>
                        </>
                      )}
                    </div>
                  </div>
                )}

                {/* Submission Results */}
                {submissionResult && (
                  <div>
                    <h3 className="text-lg font-semibold text-text-primary mb-3">
                      Submission Results
                    </h3>
                    
                    {/* Score Card */}
                    <div
                      className={`p-4 rounded-lg mb-4 ${
                        submissionResult.allPassed
                          ? 'bg-green-500/10 border border-green-500/30'
                          : 'bg-red-500/10 border border-red-500/30'
                      }`}
                    >
                      <div className="flex items-center justify-between">
                        <div>
                          <div
                            className={`text-lg font-semibold ${
                              submissionResult.allPassed ? 'text-green-400' : 'text-red-400'
                            }`}
                          >
                            {submissionResult.allPassed
                              ? '✓ All Tests Passed!'
                              : `${submissionResult.passedCount}/${submissionResult.totalCount} Tests Passed`}
                          </div>
                          <div className="text-sm text-text-secondary mt-1">
                            Score: {submissionResult.score} / {exerciseData.points || 100} points
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Test Case Results (Visible Only) */}
                    <div className="space-y-3">
                      {submissionResult.visibleResults.map((result: any, index: number) => (
                        <div
                          key={index}
                          className={`p-3 rounded-lg ${
                            result.passed
                              ? 'bg-green-500/10 border border-green-500/30'
                              : 'bg-red-500/10 border border-red-500/30'
                          }`}
                        >
                          <div className="flex items-center justify-between mb-2">
                            <span className="text-sm font-medium">
                              Visible Test Case {index + 1}
                            </span>
                            <span
                              className={`text-sm ${
                                result.passed ? 'text-green-400' : 'text-red-400'
                              }`}
                            >
                              {result.passed ? '✓ Passed' : '✗ Failed'}
                            </span>
                          </div>
                          {!result.passed && (
                            <div className="text-xs space-y-2">
                              <div>
                                <span className="text-text-tertiary">Expected:</span>
                                <pre className="bg-bg-primary p-2 rounded mt-1 overflow-x-auto">
                                  {result.testCase.expected_output}
                                </pre>
                              </div>
                              <div>
                                <span className="text-text-tertiary">Got:</span>
                                <pre className="bg-bg-primary p-2 rounded mt-1 overflow-x-auto">
                                  {result.actual_output || '(no output)'}
                                </pre>
                              </div>
                              {result.error && (
                                <div>
                                  <span className="text-red-400">Error:</span>
                                  <pre className="bg-bg-primary p-2 rounded mt-1 overflow-x-auto text-red-400">
                                    {result.error}
                                  </pre>
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      ))}
                      
                      {/* Hidden Tests Summary */}
                      {submissionResult.results.length > submissionResult.visibleResults.length && (
                        <div className="p-3 bg-bg-primary border border-border-default rounded-lg">
                          <div className="flex items-center gap-2 text-sm text-text-secondary">
                            <EyeOff className="w-4 h-4" />
                            <span>
                              {submissionResult.results.length - submissionResult.visibleResults.length} hidden test
                              case(s) - results not shown
                            </span>
                          </div>
                        </div>
                      )}
                    </div>
                  </div>
                )}

                {!testOutput && !submissionResult && (
                  <div className="text-center text-text-tertiary py-12">
                    <p>No results yet</p>
                    <p className="text-sm mt-2">Run your code to see output</p>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
