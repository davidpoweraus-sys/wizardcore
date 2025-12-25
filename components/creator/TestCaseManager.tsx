'use client'

import { useState } from 'react'
import { Trash2, Plus, Eye, EyeOff, Play } from 'lucide-react'

export interface TestCase {
  id: string
  input: string
  expected_output: string
  is_hidden: boolean
  points: number
  sort_order: number
}

interface TestCaseManagerProps {
  testCases: TestCase[]
  onChange: (testCases: TestCase[]) => void
  onRunTest?: (testCase: TestCase) => Promise<{
    passed: boolean
    actual_output: string
    execution_time?: number
    error?: string
  }>
  isRunningTest?: boolean
}

export default function TestCaseManager({
  testCases,
  onChange,
  onRunTest,
  isRunningTest = false,
}: TestCaseManagerProps) {
  const [expandedTest, setExpandedTest] = useState<string | null>(null)
  const [testResults, setTestResults] = useState<Record<string, any>>({})

  const addTestCase = () => {
    const newTestCase: TestCase = {
      id: `test-${Date.now()}`,
      input: '',
      expected_output: '',
      is_hidden: false,
      points: 10,
      sort_order: testCases.length,
    }
    onChange([...testCases, newTestCase])
    setExpandedTest(newTestCase.id)
  }

  const updateTestCase = (id: string, updates: Partial<TestCase>) => {
    onChange(
      testCases.map((tc) => (tc.id === id ? { ...tc, ...updates } : tc))
    )
  }

  const deleteTestCase = (id: string) => {
    onChange(testCases.filter((tc) => tc.id !== id))
    if (expandedTest === id) {
      setExpandedTest(null)
    }
  }

  const moveTestCase = (id: string, direction: 'up' | 'down') => {
    const index = testCases.findIndex((tc) => tc.id === id)
    if (
      (direction === 'up' && index === 0) ||
      (direction === 'down' && index === testCases.length - 1)
    ) {
      return
    }

    const newTestCases = [...testCases]
    const newIndex = direction === 'up' ? index - 1 : index + 1
    ;[newTestCases[index], newTestCases[newIndex]] = [
      newTestCases[newIndex],
      newTestCases[index],
    ]

    // Update sort_order
    newTestCases.forEach((tc, idx) => {
      tc.sort_order = idx
    })

    onChange(newTestCases)
  }

  const handleRunTest = async (testCase: TestCase) => {
    if (!onRunTest) return

    setTestResults((prev) => ({ ...prev, [testCase.id]: { loading: true } }))

    try {
      const result = await onRunTest(testCase)
      setTestResults((prev) => ({ ...prev, [testCase.id]: result }))
    } catch (error: any) {
      setTestResults((prev) => ({
        ...prev,
        [testCase.id]: {
          passed: false,
          error: error.message,
        },
      }))
    }
  }

  const totalPoints = testCases.reduce((sum, tc) => sum + tc.points, 0)
  const visibleTests = testCases.filter((tc) => !tc.is_hidden).length
  const hiddenTests = testCases.filter((tc) => tc.is_hidden).length

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-semibold text-text-primary">Test Cases</h3>
          <p className="text-sm text-text-secondary">
            {testCases.length} total • {visibleTests} visible • {hiddenTests} hidden •{' '}
            {totalPoints} points
          </p>
        </div>
        <button
          onClick={addTestCase}
          className="flex items-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors"
        >
          <Plus className="w-4 h-4" />
          Add Test Case
        </button>
      </div>

      {/* Test Cases List */}
      {testCases.length === 0 ? (
        <div className="border-2 border-dashed border-border-default rounded-lg p-8 text-center">
          <p className="text-text-secondary mb-4">No test cases yet</p>
          <button
            onClick={addTestCase}
            className="px-4 py-2 bg-bg-elevated border border-border-default rounded-lg hover:bg-bg-hover transition-colors"
          >
            Add Your First Test Case
          </button>
        </div>
      ) : (
        <div className="space-y-2">
          {testCases.map((testCase, index) => {
            const isExpanded = expandedTest === testCase.id
            const result = testResults[testCase.id]

            return (
              <div
                key={testCase.id}
                className="border border-border-default rounded-lg bg-bg-elevated overflow-hidden"
              >
                {/* Test Case Header */}
                <div
                  className="flex items-center gap-4 p-4 cursor-pointer hover:bg-bg-hover transition-colors"
                  onClick={() =>
                    setExpandedTest(isExpanded ? null : testCase.id)
                  }
                >
                  <div className="flex items-center gap-2 flex-1">
                    <span className="text-sm font-medium text-text-secondary">
                      #{index + 1}
                    </span>
                    {testCase.is_hidden ? (
                      <EyeOff className="w-4 h-4 text-text-tertiary" />
                    ) : (
                      <Eye className="w-4 h-4 text-accent-primary" />
                    )}
                    <span className="text-sm text-text-primary font-medium">
                      {testCase.is_hidden ? 'Hidden' : 'Visible'} Test Case
                    </span>
                    <span className="text-xs text-text-tertiary">
                      {testCase.points} pts
                    </span>

                    {/* Result Indicator */}
                    {result && !result.loading && (
                      <span
                        className={`text-xs px-2 py-1 rounded ${
                          result.passed
                            ? 'bg-green-500/20 text-green-400'
                            : 'bg-red-500/20 text-red-400'
                        }`}
                      >
                        {result.passed ? '✓ Passed' : '✗ Failed'}
                      </span>
                    )}
                  </div>

                  <div className="flex items-center gap-2">
                    {onRunTest && (
                      <button
                        onClick={(e) => {
                          e.stopPropagation()
                          handleRunTest(testCase)
                        }}
                        disabled={isRunningTest || result?.loading}
                        className="p-2 hover:bg-bg-hover rounded-lg transition-colors disabled:opacity-50"
                        title="Run this test"
                      >
                        <Play className="w-4 h-4 text-accent-primary" />
                      </button>
                    )}

                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        moveTestCase(testCase.id, 'up')
                      }}
                      disabled={index === 0}
                      className="p-2 hover:bg-bg-hover rounded-lg transition-colors disabled:opacity-30"
                      title="Move up"
                    >
                      ↑
                    </button>

                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        moveTestCase(testCase.id, 'down')
                      }}
                      disabled={index === testCases.length - 1}
                      className="p-2 hover:bg-bg-hover rounded-lg transition-colors disabled:opacity-30"
                      title="Move down"
                    >
                      ↓
                    </button>

                    <button
                      onClick={(e) => {
                        e.stopPropagation()
                        deleteTestCase(testCase.id)
                      }}
                      className="p-2 hover:bg-red-500/20 rounded-lg transition-colors text-red-400"
                      title="Delete test case"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                {/* Expanded Test Case Form */}
                {isExpanded && (
                  <div className="border-t border-border-default p-4 space-y-4">
                    <div className="grid grid-cols-2 gap-4">
                      {/* Points */}
                      <div>
                        <label className="block text-sm font-medium text-text-secondary mb-2">
                          Points
                        </label>
                        <input
                          type="number"
                          value={testCase.points}
                          onChange={(e) =>
                            updateTestCase(testCase.id, {
                              points: parseInt(e.target.value) || 0,
                            })
                          }
                          className="w-full px-3 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                          min="0"
                        />
                      </div>

                      {/* Visibility */}
                      <div>
                        <label className="block text-sm font-medium text-text-secondary mb-2">
                          Visibility
                        </label>
                        <label className="flex items-center gap-2 cursor-pointer">
                          <input
                            type="checkbox"
                            checked={testCase.is_hidden}
                            onChange={(e) =>
                              updateTestCase(testCase.id, {
                                is_hidden: e.target.checked,
                              })
                            }
                            className="w-4 h-4 rounded border-border-default bg-bg-primary text-accent-primary focus:ring-2 focus:ring-accent-primary"
                          />
                          <span className="text-sm text-text-primary">
                            Hidden from students
                          </span>
                        </label>
                      </div>
                    </div>

                    {/* Input */}
                    <div>
                      <label className="block text-sm font-medium text-text-secondary mb-2">
                        Input (stdin)
                      </label>
                      <textarea
                        value={testCase.input}
                        onChange={(e) =>
                          updateTestCase(testCase.id, { input: e.target.value })
                        }
                        placeholder="Enter input for this test case..."
                        className="w-full px-3 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary font-mono text-sm focus:outline-none focus:ring-2 focus:ring-accent-primary"
                        rows={3}
                      />
                      <p className="text-xs text-text-tertiary mt-1">
                        Leave empty if no input is needed
                      </p>
                    </div>

                    {/* Expected Output */}
                    <div>
                      <label className="block text-sm font-medium text-text-secondary mb-2">
                        Expected Output
                      </label>
                      <textarea
                        value={testCase.expected_output}
                        onChange={(e) =>
                          updateTestCase(testCase.id, {
                            expected_output: e.target.value,
                          })
                        }
                        placeholder="Expected output from the program..."
                        className="w-full px-3 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary font-mono text-sm focus:outline-none focus:ring-2 focus:ring-accent-primary"
                        rows={3}
                      />
                    </div>

                    {/* Test Result Display */}
                    {result && !result.loading && (
                      <div
                        className={`p-4 rounded-lg ${
                          result.passed
                            ? 'bg-green-500/10 border border-green-500/30'
                            : 'bg-red-500/10 border border-red-500/30'
                        }`}
                      >
                        <div className="flex items-center justify-between mb-2">
                          <span
                            className={`font-semibold ${
                              result.passed ? 'text-green-400' : 'text-red-400'
                            }`}
                          >
                            {result.passed ? '✓ Test Passed' : '✗ Test Failed'}
                          </span>
                          {result.execution_time && (
                            <span className="text-xs text-text-tertiary">
                              {result.execution_time}ms
                            </span>
                          )}
                        </div>

                        {!result.passed && (
                          <div className="space-y-2">
                            <div>
                              <p className="text-xs font-semibold text-text-secondary mb-1">
                                Actual Output:
                              </p>
                              <pre className="text-xs font-mono bg-bg-primary p-2 rounded overflow-x-auto">
                                {result.actual_output || '(no output)'}
                              </pre>
                            </div>
                            {result.error && (
                              <div>
                                <p className="text-xs font-semibold text-red-400 mb-1">
                                  Error:
                                </p>
                                <pre className="text-xs font-mono bg-bg-primary p-2 rounded overflow-x-auto text-red-400">
                                  {result.error}
                                </pre>
                              </div>
                            )}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                )}
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}
