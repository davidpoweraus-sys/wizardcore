'use client'

import { useState } from 'react'
import Editor from '@monaco-editor/react'
import { Play, Save, Eye, Code, FileText, TestTube } from 'lucide-react'
import TestCaseManager, { TestCase } from './TestCaseManager'
import ExercisePreview from './ExercisePreview'
import { submitCode, LANGUAGES } from '@/lib/judge0/service'

interface ExerciseData {
  title: string
  difficulty: 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED'
  points: number
  time_limit_minutes: number
  objectives: string[]
  content: string
  description: string
  constraints: string[]
  hints: string[]
  starter_code: string
  solution_code: string
  language_id: number
  tags: string[]
  test_cases: TestCase[]
}

interface ExerciseBuilderProps {
  moduleId: string
  initialData?: Partial<ExerciseData>
  onSave?: (data: ExerciseData) => Promise<void>
  onPreview?: (data: ExerciseData) => void
}

export default function ExerciseBuilder({
  moduleId,
  initialData,
  onSave,
  onPreview,
}: ExerciseBuilderProps) {
  const [activeTab, setActiveTab] = useState<'details' | 'code' | 'tests' | 'preview'>('details')
  const [isSaving, setIsSaving] = useState(false)
  const [isTestingSolution, setIsTestingSolution] = useState(false)
  const [testResults, setTestResults] = useState<any>(null)
  const [showPreview, setShowPreview] = useState(false)

  const [formData, setFormData] = useState<ExerciseData>({
    title: '',
    difficulty: 'BEGINNER',
    points: 100,
    time_limit_minutes: 30,
    objectives: [''],
    content: '',
    description: '',
    constraints: [''],
    hints: [''],
    starter_code: '# Write your code here\n',
    solution_code: '# Solution code\n',
    language_id: 71, // Python default
    tags: [],
    test_cases: [],
    ...initialData,
  })

  const selectedLanguage = LANGUAGES.find(
    (lang) => lang.id === formData.language_id
  )

  const updateField = <K extends keyof ExerciseData>(
    field: K,
    value: ExerciseData[K]
  ) => {
    setFormData((prev) => ({ ...prev, [field]: value }))
  }

  const updateArrayField = (
    field: 'objectives' | 'constraints' | 'hints',
    index: number,
    value: string
  ) => {
    const newArray = [...formData[field]]
    newArray[index] = value
    updateField(field, newArray)
  }

  const addArrayItem = (field: 'objectives' | 'constraints' | 'hints') => {
    updateField(field, [...formData[field], ''])
  }

  const removeArrayItem = (
    field: 'objectives' | 'constraints' | 'hints',
    index: number
  ) => {
    updateField(
      field,
      formData[field].filter((_, i) => i !== index)
    )
  }

  const addTag = (tag: string) => {
    if (tag && !formData.tags.includes(tag)) {
      updateField('tags', [...formData.tags, tag])
    }
  }

  const removeTag = (tag: string) => {
    updateField(
      'tags',
      formData.tags.filter((t) => t !== tag)
    )
  }

  // Test solution code against all test cases
  const testSolution = async () => {
    if (formData.test_cases.length === 0) {
      alert('Please add at least one test case first')
      return
    }

    setIsTestingSolution(true)
    setTestResults(null)

    try {
      const results = await Promise.all(
        formData.test_cases.map(async (testCase) => {
          try {
            const result = await submitCode({
              source_code: formData.solution_code,
              language_id: formData.language_id,
              stdin: testCase.input,
            })

            const passed =
              result.stdout?.trim() === testCase.expected_output.trim()

            return {
              testCase,
              passed,
              actual_output: result.stdout || '',
              execution_time: result.time ? parseFloat(result.time) * 1000 : null,
              error: result.stderr || result.compile_output || null,
              status: result.status.description,
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

      setTestResults({
        results,
        passedCount,
        totalCount,
        allPassed: passedCount === totalCount,
      })
    } catch (error: any) {
      alert(`Testing failed: ${error.message}`)
    } finally {
      setIsTestingSolution(false)
    }
  }

  // Test a single test case
  const runSingleTest = async (testCase: TestCase) => {
    const result = await submitCode({
      source_code: formData.solution_code,
      language_id: formData.language_id,
      stdin: testCase.input,
    })

    const passed = result.stdout?.trim() === testCase.expected_output.trim()
    const executionTime = result.time ? parseFloat(result.time) * 1000 : undefined
    const errorMsg = result.stderr || result.compile_output || undefined

    return {
      passed,
      actual_output: result.stdout || '',
      execution_time: executionTime,
      error: errorMsg,
    }
  }

  const handleSave = async () => {
    // Validation
    if (!formData.title.trim()) {
      alert('Please enter a title')
      return
    }

    if (formData.test_cases.length === 0) {
      alert('Please add at least one test case')
      return
    }

    if (!formData.solution_code.trim()) {
      alert('Please provide solution code')
      return
    }

    setIsSaving(true)
    try {
      await onSave?.(formData)
    } catch (error: any) {
      alert(`Save failed: ${error.message}`)
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <>
      {/* Preview Modal */}
      <ExercisePreview
        exerciseData={formData}
        isOpen={showPreview}
        onClose={() => setShowPreview(false)}
      />

      <div className="min-h-screen bg-bg-primary">
      {/* Header */}
      <div className="border-b border-border-default bg-bg-elevated">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-2xl font-bold text-text-primary">
                {initialData?.title ? 'Edit Exercise' : 'Create New Exercise'}
              </h1>
              <p className="text-sm text-text-secondary mt-1">
                Build coding exercises with live Judge0 testing
              </p>
            </div>

            <div className="flex items-center gap-3">
              <button
                onClick={() => setShowPreview(true)}
                type="button"
                className="flex items-center gap-2 px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-text-primary"
              >
                <Eye className="w-4 h-4" />
                Preview
              </button>

              <button
                onClick={handleSave}
                disabled={isSaving}
                className="flex items-center gap-2 px-6 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors disabled:opacity-50"
              >
                <Save className="w-4 h-4" />
                {isSaving ? 'Saving...' : 'Save Exercise'}
              </button>
            </div>
          </div>

          {/* Tabs */}
          <div className="flex gap-2 mt-4">
            {[
              { id: 'details', label: 'Details', icon: FileText },
              { id: 'code', label: 'Code & Solution', icon: Code },
              { id: 'tests', label: 'Test Cases', icon: TestTube },
            ].map((tab) => {
              const Icon = tab.icon
              return (
                <button
                  key={tab.id}
                  onClick={() => setActiveTab(tab.id as any)}
                  className={`flex items-center gap-2 px-4 py-2 rounded-lg transition-colors ${
                    activeTab === tab.id
                      ? 'bg-accent-primary text-white'
                      : 'text-text-secondary hover:text-text-primary hover:bg-bg-hover'
                  }`}
                >
                  <Icon className="w-4 h-4" />
                  {tab.label}
                </button>
              )
            })}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="max-w-7xl mx-auto px-6 py-8">
        {/* Details Tab */}
        {activeTab === 'details' && (
          <div className="space-y-6">
            {/* Basic Info */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Basic Information
              </h2>

              <div className="space-y-4">
                {/* Title */}
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-2">
                    Exercise Title *
                  </label>
                  <input
                    type="text"
                    value={formData.title}
                    onChange={(e) => updateField('title', e.target.value)}
                    placeholder="e.g., Stack Buffer Overflow Exploitation"
                    className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                  />
                </div>

                {/* Description */}
                <div>
                  <label className="block text-sm font-medium text-text-secondary mb-2">
                    Short Description
                  </label>
                  <textarea
                    value={formData.description}
                    onChange={(e) => updateField('description', e.target.value)}
                    placeholder="Brief description of what students will learn..."
                    className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    rows={3}
                  />
                </div>

                {/* Grid: Difficulty, Points, Time Limit, Language */}
                <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-2">
                      Difficulty
                    </label>
                    <select
                      value={formData.difficulty}
                      onChange={(e) =>
                        updateField('difficulty', e.target.value as any)
                      }
                      className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    >
                      <option value="BEGINNER">Beginner</option>
                      <option value="INTERMEDIATE">Intermediate</option>
                      <option value="ADVANCED">Advanced</option>
                    </select>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-2">
                      Points
                    </label>
                    <input
                      type="number"
                      value={formData.points}
                      onChange={(e) =>
                        updateField('points', parseInt(e.target.value) || 0)
                      }
                      className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                      min="0"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-2">
                      Time Limit (min)
                    </label>
                    <input
                      type="number"
                      value={formData.time_limit_minutes}
                      onChange={(e) =>
                        updateField(
                          'time_limit_minutes',
                          parseInt(e.target.value) || 0
                        )
                      }
                      className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                      min="0"
                    />
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-text-secondary mb-2">
                      Language
                    </label>
                    <select
                      value={formData.language_id}
                      onChange={(e) =>
                        updateField('language_id', parseInt(e.target.value))
                      }
                      className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    >
                      {LANGUAGES.map((lang) => (
                        <option key={lang.id} value={lang.id}>
                          {lang.name}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>
              </div>
            </div>

            {/* Content (Markdown) */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Exercise Content (Markdown)
              </h2>
              <textarea
                value={formData.content}
                onChange={(e) => updateField('content', e.target.value)}
                placeholder="# Exercise Title&#10;&#10;## Problem Description&#10;Write detailed instructions here...&#10;&#10;## Examples&#10;..."
                className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary font-mono text-sm focus:outline-none focus:ring-2 focus:ring-accent-primary"
                rows={12}
              />
              <p className="text-xs text-text-tertiary mt-2">
                Use Markdown formatting. This will be displayed to students.
              </p>
            </div>

            {/* Objectives */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Learning Objectives
              </h2>
              <div className="space-y-2">
                {formData.objectives.map((obj, index) => (
                  <div key={index} className="flex gap-2">
                    <input
                      type="text"
                      value={obj}
                      onChange={(e) =>
                        updateArrayField('objectives', index, e.target.value)
                      }
                      placeholder="What students will learn..."
                      className="flex-1 px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    />
                    <button
                      onClick={() => removeArrayItem('objectives', index)}
                      className="px-3 py-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
                    >
                      Remove
                    </button>
                  </div>
                ))}
                <button
                  onClick={() => addArrayItem('objectives')}
                  className="px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-text-primary"
                >
                  + Add Objective
                </button>
              </div>
            </div>

            {/* Constraints */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Constraints
              </h2>
              <div className="space-y-2">
                {formData.constraints.map((constraint, index) => (
                  <div key={index} className="flex gap-2">
                    <input
                      type="text"
                      value={constraint}
                      onChange={(e) =>
                        updateArrayField('constraints', index, e.target.value)
                      }
                      placeholder="e.g., Input size <= 10^6"
                      className="flex-1 px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    />
                    <button
                      onClick={() => removeArrayItem('constraints', index)}
                      className="px-3 py-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
                    >
                      Remove
                    </button>
                  </div>
                ))}
                <button
                  onClick={() => addArrayItem('constraints')}
                  className="px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-text-primary"
                >
                  + Add Constraint
                </button>
              </div>
            </div>

            {/* Hints */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Hints (Progressive)
              </h2>
              <div className="space-y-2">
                {formData.hints.map((hint, index) => (
                  <div key={index} className="flex gap-2">
                    <span className="px-3 py-2 bg-bg-primary rounded-lg text-text-secondary">
                      #{index + 1}
                    </span>
                    <input
                      type="text"
                      value={hint}
                      onChange={(e) =>
                        updateArrayField('hints', index, e.target.value)
                      }
                      placeholder="Hint for students..."
                      className="flex-1 px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                    />
                    <button
                      onClick={() => removeArrayItem('hints', index)}
                      className="px-3 py-2 bg-red-500/20 text-red-400 rounded-lg hover:bg-red-500/30 transition-colors"
                    >
                      Remove
                    </button>
                  </div>
                ))}
                <button
                  onClick={() => addArrayItem('hints')}
                  className="px-4 py-2 bg-bg-primary border border-border-default rounded-lg hover:bg-bg-hover transition-colors text-text-primary"
                >
                  + Add Hint
                </button>
              </div>
            </div>

            {/* Tags */}
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Tags
              </h2>
              <div className="flex flex-wrap gap-2 mb-3">
                {formData.tags.map((tag) => (
                  <span
                    key={tag}
                    className="px-3 py-1 bg-accent-primary/20 text-accent-primary rounded-full text-sm flex items-center gap-2"
                  >
                    {tag}
                    <button
                      onClick={() => removeTag(tag)}
                      className="hover:text-red-400"
                    >
                      ×
                    </button>
                  </span>
                ))}
              </div>
              <input
                type="text"
                placeholder="Type a tag and press Enter..."
                className="w-full px-4 py-2 bg-bg-primary border border-border-default rounded-lg text-text-primary focus:outline-none focus:ring-2 focus:ring-accent-primary"
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    const input = e.currentTarget
                    addTag(input.value.trim())
                    input.value = ''
                  }
                }}
              />
            </div>
          </div>
        )}

        {/* Code Tab */}
        {activeTab === 'code' && (
          <div className="space-y-6">
            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <h2 className="text-lg font-semibold text-text-primary mb-4">
                Starter Code
              </h2>
              <p className="text-sm text-text-secondary mb-4">
                This code will be provided to students when they start the exercise
              </p>
              <Editor
                height="300px"
                language={selectedLanguage?.name.toLowerCase().split(' ')[0] || 'python'}
                value={formData.starter_code}
                onChange={(value) => updateField('starter_code', value || '')}
                theme="vs-dark"
                options={{
                  minimap: { enabled: false },
                  fontSize: 14,
                  scrollBeyondLastLine: false,
                }}
              />
            </div>

            <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
              <div className="flex items-center justify-between mb-4">
                <div>
                  <h2 className="text-lg font-semibold text-text-primary">
                    Solution Code
                  </h2>
                  <p className="text-sm text-text-secondary">
                    Your reference solution - will be tested against test cases
                  </p>
                </div>
                <button
                  onClick={testSolution}
                  disabled={isTestingSolution || formData.test_cases.length === 0}
                  className="flex items-center gap-2 px-4 py-2 bg-accent-primary text-white rounded-lg hover:bg-accent-secondary transition-colors disabled:opacity-50"
                >
                  <Play className="w-4 h-4" />
                  {isTestingSolution ? 'Testing...' : 'Test Solution'}
                </button>
              </div>

              <Editor
                height="400px"
                language={selectedLanguage?.name.toLowerCase().split(' ')[0] || 'python'}
                value={formData.solution_code}
                onChange={(value) => updateField('solution_code', value || '')}
                theme="vs-dark"
                options={{
                  minimap: { enabled: false },
                  fontSize: 14,
                  scrollBeyondLastLine: false,
                }}
              />

              {/* Test Results */}
              {testResults && (
                <div className="mt-4">
                  <div
                    className={`p-4 rounded-lg ${
                      testResults.allPassed
                        ? 'bg-green-500/10 border border-green-500/30'
                        : 'bg-red-500/10 border border-red-500/30'
                    }`}
                  >
                    <div className="flex items-center justify-between mb-3">
                      <span
                        className={`text-lg font-semibold ${
                          testResults.allPassed ? 'text-green-400' : 'text-red-400'
                        }`}
                      >
                        {testResults.allPassed
                          ? '✓ All Tests Passed!'
                          : `✗ ${testResults.passedCount}/${testResults.totalCount} Tests Passed`}
                      </span>
                    </div>

                    <div className="space-y-2">
                      {testResults.results.map((result: any, index: number) => (
                        <div
                          key={index}
                          className={`p-3 rounded-lg ${
                            result.passed
                              ? 'bg-green-500/20'
                              : 'bg-red-500/20'
                          }`}
                        >
                          <div className="flex items-center justify-between mb-1">
                            <span className="text-sm font-medium">
                              Test Case #{index + 1}
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
                            <div className="text-xs space-y-1 mt-2">
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
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        )}

        {/* Tests Tab */}
        {activeTab === 'tests' && (
          <div className="bg-bg-elevated border border-border-default rounded-lg p-6">
            <TestCaseManager
              testCases={formData.test_cases}
              onChange={(testCases) => updateField('test_cases', testCases)}
              onRunTest={runSingleTest}
              isRunningTest={isTestingSolution}
            />
          </div>
        )}
      </div>
    </div>
    </>
  )
}
