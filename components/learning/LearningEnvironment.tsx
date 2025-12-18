'use client'

import { useState } from 'react'
import { Menu, BookOpen, Code2, CheckCircle2, Lightbulb, ChevronRight, Play, Send, Zap, Target, Users, Clock, Loader2 } from 'lucide-react'
import { submitCode, LANGUAGES } from '@/lib/judge0/service'

export default function LearningEnvironment() {
  const [sidebarOpen, setSidebarOpen] = useState(true)
  const [activeTab, setActiveTab] = useState<'lesson' | 'code' | 'tests' | 'hints'>('lesson')
  const [code, setCode] = useState(`def calculate_sum(numbers):
    """
    Calculate the sum of a list of numbers.
    
    Args:
        numbers (list): List of integers
    
    Returns:
        int: Sum of the numbers
    """
    # Write your solution here
    total = 0
    for num in numbers:
        total += num
    return total`)

  const lesson = {
    title: 'Python Functions: Sum Calculation',
    module: 'Module 1: The Hacker\'s Toolkit',
    objectives: [
      'Understand function syntax in Python',
      'Learn to iterate over lists',
      'Practice returning values from functions',
      'Apply to security tooling scenarios',
    ],
    content: `In offensive security, you'll often need to process lists of data—IP addresses, ports, payloads, etc. Being able to quickly sum, filter, and transform lists is essential.

## Why This Matters
- Log analysis: Summing request counts
- Payload generation: Calculating checksums
- Data exfiltration: Aggregating stolen data sizes

## Key Concepts
- Functions encapsulate reusable logic
- Loops iterate over collections
- Return statements provide output`,
    examples: [
      {
        code: `# Basic list iteration
ports = [80, 443, 22, 8080]
for port in ports:
    print(f"Scanning port {port}")`,
        output: 'Scanning port 80\nScanning port 443\nScanning port 22\nScanning port 8080',
      },
      {
        code: `# Function with return
def check_port_open(port):
    return port in open_ports`,
        output: '',
      },
    ],
  }

  const exercise = {
    title: 'Sum of Port List',
    difficulty: 'BEGINNER',
    points: 100,
    description: `In network scanning, you often need to calculate the total number of ports scanned. Write a function that sums a list of port numbers.

Your function should:
1. Accept a list of integers (port numbers)
2. Return the sum of all numbers in the list
3. Handle empty lists (return 0)`,
    constraints: [
      'Do not use the built-in sum() function',
      'Time complexity must be O(n)',
      'Function name must be calculate_sum',
    ],
    testCases: [
      { input: '[1, 2, 3]', expectedOutput: '6' },
      { input: '[]', expectedOutput: '0' },
      { input: '[80, 443, 22]', expectedOutput: '545' },
    ],
    hints: [
      'Use a for loop to iterate through the list',
      'Initialize a variable to store the total',
      'Add each element to the total',
      'Return the total after the loop',
    ],
  }

  const [output, setOutput] = useState<string>('')
  const [isRunning, setIsRunning] = useState(false)

  const handleRunCode = async () => {
    setIsRunning(true)
    setOutput('> Running code...\n')
    try {
      const result = await submitCode({
        source_code: code,
        language_id: 71, // Python 3.8.1
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

  const handleSubmit = () => {
    console.log('Submitting solution')
    // Validate and submit to backend
  }

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
            <p className="text-sm text-text-secondary">{lesson.module}</p>
          </div>
        </div>
        
        <div className="flex items-center gap-6">
          <div className="hidden md:flex items-center gap-4">
            <div className="text-sm text-text-secondary">
              <Clock className="w-4 h-4 inline mr-2" />
              15 min remaining
            </div>
            <div className="flex items-center gap-2">
              <Users className="w-4 h-4 text-text-muted" />
              <span className="text-sm text-text-secondary">1,245 solving</span>
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
              <div className="mb-8 p-4 rounded-xl bg-gradient-to-r from-neon-cyan/10 to-neon-lavender/10 border border-neon-cyan/20">
                <h3 className="font-semibold text-text-primary mb-3 flex items-center gap-2">
                  <Target className="w-4 h-4" />
                  Learning Objectives
                </h3>
                <ul className="space-y-2">
                  {lesson.objectives.map((obj, idx) => (
                    <li key={idx} className="text-sm text-text-secondary flex items-start gap-2">
                      <ChevronRight className="w-4 h-4 text-neon-cyan mt-0.5 flex-shrink-0" />
                      {obj}
                    </li>
                  ))}
                </ul>
              </div>

              {/* Lesson Content */}
              <div className="prose prose-sm max-w-none text-text-primary">
                <h3 className="text-xl font-bold mb-4">{lesson.title}</h3>
                <div className="whitespace-pre-line text-text-secondary mb-6">
                  {lesson.content}
                </div>
              </div>

              {/* Examples */}
              <div className="mb-8">
                <h3 className="font-semibold text-text-primary mb-3 flex items-center gap-2">
                  <Code2 className="w-4 h-4" />
                  Examples
                </h3>
                <div className="space-y-4">
                  {lesson.examples.map((example, idx) => (
                    <div key={idx} className="bg-bg-tertiary border border-border-subtle rounded-lg p-4">
                      <pre className="text-sm font-mono text-green-400 overflow-x-auto">
                        {example.code}
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

              {/* Exercise Instructions */}
              <div className="border-t border-border-subtle pt-6">
                <h3 className="font-bold text-text-primary mb-3">Exercise Instructions</h3>
                <p className="text-text-secondary mb-4">{exercise.description}</p>
                
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

                <div>
                  <h4 className="text-sm font-semibold text-text-primary mb-2">Example Test Cases:</h4>
                  {exercise.testCases.slice(0, 2).map((test, idx) => (
                    <div key={idx} className="text-sm mb-3 bg-bg-tertiary p-3 rounded border border-border-subtle">
                      <div className="font-mono">
                        <div><span className="text-text-muted">Input:</span> {test.input}</div>
                        <div><span className="text-text-muted">Expected Output:</span> {test.expectedOutput}</div>
                      </div>
                    </div>
                  ))}
                </div>
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
                    <p>This lesson teaches you how to write Python functions for security tooling.</p>
                    <p>You'll apply this knowledge to build a port scanner aggregator in the next module.</p>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'code' && (
              <div className="h-full flex flex-col">
                {/* Editor */}
                <div className="flex-1 p-4">
                  <div className="h-full border border-border-subtle rounded-lg overflow-hidden">
                    <div className="bg-bg-tertiary px-4 py-2 border-b border-border-subtle text-sm text-text-secondary font-mono">
                      exercise.py
                    </div>
                    <textarea
                      value={code}
                      onChange={(e) => setCode(e.target.value)}
                      className="w-full h-full p-4 font-mono text-sm bg-bg-tertiary text-text-primary resize-none focus:outline-none"
                      spellCheck="false"
                    />
                  </div>
                </div>

                {/* Actions & Output */}
                <div className="border-t border-border-default p-4">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex gap-3">
                      <button
                        onClick={handleRunCode}
                        disabled={isRunning}
                        className="px-4 py-2 bg-gradient-to-r from-neon-cyan to-neon-lavender text-white rounded-lg font-medium hover:opacity-90 transition flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        {isRunning ? (
                          <>
                            <Loader2 className="w-4 h-4 animate-spin" />
                            Running...
                          </>
                        ) : (
                          <>
                            <Play className="w-4 h-4" />
                            Run Code
                          </>
                        )}
                      </button>
                      <button
                        onClick={handleSubmit}
                        className="px-4 py-2 bg-gradient-to-r from-neon-pink to-neon-purple text-white rounded-lg font-medium hover:opacity-90 transition flex items-center gap-2"
                      >
                        <Send className="w-4 h-4" />
                        Submit Solution
                      </button>
                      <button className="px-4 py-2 border border-border-subtle text-text-secondary rounded-lg hover:bg-bg-tertiary transition">
                        Reset
                      </button>
                    </div>
                    <div className="flex items-center gap-2 text-sm text-text-muted">
                      <Zap className="w-4 h-4" />
                      <span>Powered by Judge0</span>
                    </div>
                  </div>

                  {/* Output Panel */}
                  <div className="bg-black/50 border border-border-subtle rounded-lg p-4">
                    <div className="text-sm text-text-secondary mb-2">Output:</div>
                    <pre className="text-sm font-mono text-green-400 whitespace-pre-wrap">
                      {output || '> No output yet. Click "Run Code" to execute.'}
                    </pre>
                  </div>
                </div>
              </div>
            )}

            {activeTab === 'tests' && (
              <div className="h-full overflow-y-auto p-6">
                <h2 className="text-2xl font-bold text-text-primary mb-6">Test Results</h2>
                <div className="space-y-4">
                  {exercise.testCases.map((test, idx) => (
                    <div key={idx} className="border border-border-subtle rounded-lg p-4">
                      <div className="flex items-center justify-between mb-2">
                        <div className="flex items-center gap-2">
                          <div className={`w-3 h-3 rounded-full ${idx < 2 ? 'bg-green-500' : 'bg-red-500'}`} />
                          <span className="font-medium text-text-primary">Test Case {idx + 1}</span>
                        </div>
                        <span className="text-sm text-text-secondary">{idx < 2 ? 'Passed' : 'Failed'}</span>
                      </div>
                      <div className="grid grid-cols-2 gap-4 text-sm">
                        <div>
                          <div className="text-text-muted">Input</div>
                          <div className="font-mono bg-bg-tertiary p-2 rounded">{test.input}</div>
                        </div>
                        <div>
                          <div className="text-text-muted">Expected</div>
                          <div className="font-mono bg-bg-tertiary p-2 rounded">{test.expectedOutput}</div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {activeTab === 'hints' && (
              <div className="h-full overflow-y-auto p-6">
                <h2 className="text-2xl font-bold text-text-primary mb-6">Hints</h2>
                <div className="space-y-4">
                  {exercise.hints.map((hint, idx) => (
                    <div key={idx} className="border border-border-subtle rounded-lg p-4">
                      <div className="flex items-start gap-3">
                        <div className="p-2 rounded-lg bg-gradient-to-r from-neon-cyan/20 to-neon-lavender/20">
                          <Lightbulb className="w-5 h-5 text-neon-cyan" />
                        </div>
                        <div>
                          <h3 className="font-medium text-text-primary mb-1">Hint {idx + 1}</h3>
                          <p className="text-text-secondary">{hint}</p>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            )}
          </div>
        </main>
      </div>
    </div>
  )
}