# Learning Environment Page Layouts

A comprehensive guide to designing effective coding education interfaces with lesson content and interactive code editors.

---

## Design Principles for Learning Environments

### Core Requirements
1. **Reduce Cognitive Load**: Students shouldn't have to choose between reading and coding
2. **Progressive Disclosure**: Show information as needed, not all at once
3. **Clear Visual Hierarchy**: Lesson → Example → Exercise → Solution
4. **Responsive Design**: Works on tablets and smaller screens
5. **Distraction-Free**: Focus on learning, minimize UI chrome

### Key Elements to Include
- 📖 Lesson content (theory, explanations)
- 💡 Code examples (read-only demonstrations)
- ✍️ Interactive exercises (editable code)
- ✅ Test cases and validation
- 🎯 Learning objectives
- 💭 Hints and tips
- 📊 Progress indicators

---

## Layout Option 1: Split View (Classic LeetCode Style)

**Best For:** Traditional coding challenges, algorithmic problems

```
┌─────────────────────────────────────────────────────────────┐
│  Header: Exercise Title, Difficulty, Progress              │
├──────────────────────┬──────────────────────────────────────┤
│                      │                                      │
│   LESSON PANEL       │      CODE EDITOR PANEL              │
│   (Scrollable)       │                                      │
│                      │   ┌──────────────────────────────┐  │
│  • Description       │   │                              │  │
│  • Examples          │   │  def solution():             │  │
│  • Constraints       │   │      # Your code here        │  │
│  • Hints             │   │      pass                    │  │
│                      │   │                              │  │
│                      │   └──────────────────────────────┘  │
│                      │                                      │
│                      │   [Run Code]  [Submit Solution]     │
│                      │                                      │
│                      │   ┌──────────────────────────────┐  │
│                      │   │  Console Output / Results    │  │
│                      │   └──────────────────────────────┘  │
│                      │                                      │
└──────────────────────┴──────────────────────────────────────┘
```

**Pros:**
- Both panels always visible
- Good for quick reference while coding
- Industry-standard pattern

**Cons:**
- Can feel cramped on smaller screens
- May require horizontal scrolling
- Content fighting for attention

### Implementation:

```typescript
// components/layouts/split-view-layout.tsx
'use client'

import { useState } from 'react'
import { Panel, PanelGroup, PanelResizeHandle } from 'react-resizable-panels'
import { GripVertical } from 'lucide-react'
import LessonContent from './lesson-content'
import CodeEditor from './code-editor'

export default function SplitViewLayout({ lesson, exercise }: any) {
  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <header className="h-16 border-b px-6 flex items-center justify-between bg-background">
        <div>
          <h1 className="text-xl font-bold">{exercise.title}</h1>
          <p className="text-sm text-muted-foreground">{lesson.title}</p>
        </div>
        <div className="flex items-center gap-4">
          <Badge>{exercise.difficulty}</Badge>
          <Button variant="ghost" size="sm">
            <BookOpen className="w-4 h-4 mr-2" />
            View Lesson
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <PanelGroup direction="horizontal" className="flex-1">
        {/* Left Panel: Lesson Content */}
        <Panel defaultSize={40} minSize={30}>
          <div className="h-full overflow-y-auto">
            <LessonContent lesson={lesson} exercise={exercise} />
          </div>
        </Panel>

        {/* Resize Handle */}
        <PanelResizeHandle className="w-2 bg-border hover:bg-primary/20 transition-colors relative group">
          <div className="absolute inset-y-0 left-1/2 -translate-x-1/2 w-1 bg-border group-hover:bg-primary transition-colors" />
          <GripVertical className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground opacity-0 group-hover:opacity-100" />
        </PanelResizeHandle>

        {/* Right Panel: Code Editor */}
        <Panel defaultSize={60} minSize={40}>
          <div className="h-full">
            <CodeEditor exercise={exercise} />
          </div>
        </Panel>
      </PanelGroup>
    </div>
  )
}
```

---

## Layout Option 2: Tabbed View (Mobile-First)

**Best For:** Tutorial-style learning, step-by-step progression

```
┌─────────────────────────────────────────────────────────────┐
│  Header: Exercise Title                                     │
├─────────────────────────────────────────────────────────────┤
│  [📖 Lesson]  [💻 Code]  [✅ Tests]  [💡 Hints]            │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│                                                             │
│               ACTIVE TAB CONTENT                            │
│               (Lesson OR Code OR Tests)                     │
│                                                             │
│                                                             │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

**Pros:**
- Works perfectly on mobile/tablets
- Clear mental model (one thing at a time)
- Larger workspace for each view

**Cons:**
- Can't see lesson while coding
- More tab switching required
- Less efficient for experienced users

### Implementation:

```typescript
// components/layouts/tabbed-view-layout.tsx
'use client'

import { useState } from 'react'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BookOpen, Code2, CheckCircle2, Lightbulb } from 'lucide-react'
import LessonContent from './lesson-content'
import CodeEditor from './code-editor'
import TestResults from './test-results'
import HintsPanel from './hints-panel'

export default function TabbedViewLayout({ lesson, exercise }: any) {
  const [activeTab, setActiveTab] = useState('lesson')

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <header className="h-16 border-b px-6 flex items-center justify-between">
        <h1 className="text-xl font-bold">{exercise.title}</h1>
        <Badge>{exercise.difficulty}</Badge>
      </header>

      {/* Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab} className="flex-1 flex flex-col">
        <TabsList className="w-full justify-start border-b rounded-none h-12 px-6">
          <TabsTrigger value="lesson" className="gap-2">
            <BookOpen className="w-4 h-4" />
            Lesson
          </TabsTrigger>
          <TabsTrigger value="code" className="gap-2">
            <Code2 className="w-4 h-4" />
            Code
          </TabsTrigger>
          <TabsTrigger value="tests" className="gap-2">
            <CheckCircle2 className="w-4 h-4" />
            Tests
          </TabsTrigger>
          <TabsTrigger value="hints" className="gap-2">
            <Lightbulb className="w-4 h-4" />
            Hints
          </TabsTrigger>
        </TabsList>

        <div className="flex-1 overflow-hidden">
          <TabsContent value="lesson" className="h-full m-0">
            <div className="h-full overflow-y-auto p-6">
              <LessonContent lesson={lesson} exercise={exercise} />
            </div>
          </TabsContent>

          <TabsContent value="code" className="h-full m-0">
            <CodeEditor exercise={exercise} onComplete={() => setActiveTab('tests')} />
          </TabsContent>

          <TabsContent value="tests" className="h-full m-0 p-6">
            <TestResults exercise={exercise} />
          </TabsContent>

          <TabsContent value="hints" className="h-full m-0 p-6">
            <HintsPanel hints={exercise.hints} />
          </TabsContent>
        </div>
      </Tabs>
    </div>
  )
}
```

---

## Layout Option 3: Collapsible Sidebar (Recommended)

**Best For:** Interactive tutorials, comprehensive learning

```
┌─────────────────────────────────────────────────────────────┐
│  Header: Course Path, Progress Bar                          │
├──┬──────────────────────────────────────────────────────────┤
│☰ │                                                          │
│📖│              CODE EDITOR (Main Focus)                    │
│  │                                                          │
│L │   ┌──────────────────────────────────────────────────┐  │
│E │   │                                                  │  │
│S │   │  def calculate_sum(numbers):                     │  │
│S │   │      """                                         │  │
│O │   │      Calculate the sum of a list                │  │
│N │   │      """                                         │  │
│  │   │      # Write your solution here                 │  │
│  │   │                                                  │  │
│  │   └──────────────────────────────────────────────────┘  │
│  │                                                          │
│  │   [▶ Run Code]  [📤 Submit]  [💡 Show Hint]            │
│  │                                                          │
│  │   Console Output:                                        │
│  │   ┌──────────────────────────────────────────────────┐  │
│  │   │ >>> calculate_sum([1, 2, 3])                     │  │
│  │   │ 6                                                │  │
│  │   └──────────────────────────────────────────────────┘  │
└──┴──────────────────────────────────────────────────────────┘
      ↑
   Click to expand lesson sidebar
```

**Pros:**
- Code editor is primary focus
- Lesson accessible but not intrusive
- Best of both worlds
- Great for progressive learning

**Cons:**
- Requires interaction to view lesson
- Slightly more complex UX

### Implementation:

```typescript
// components/layouts/collapsible-sidebar-layout.tsx
'use client'

import { useState } from 'react'
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from '@/components/ui/sheet'
import { Button } from '@/components/ui/button'
import { Menu, BookOpen, ChevronRight } from 'lucide-react'
import { cn } from '@/lib/utils'
import LessonContent from './lesson-content'
import CodeEditor from './code-editor'

export default function CollapsibleSidebarLayout({ lesson, exercise }: any) {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const [isPinned, setIsPinned] = useState(false)

  return (
    <div className="h-screen flex flex-col">
      {/* Header */}
      <header className="h-16 border-b px-6 flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button
            variant="ghost"
            size="icon"
            onClick={() => setSidebarOpen(!sidebarOpen)}
          >
            <Menu className="w-5 h-5" />
          </Button>
          <div>
            <h1 className="text-lg font-bold">{exercise.title}</h1>
            <p className="text-sm text-muted-foreground">{lesson.title}</p>
          </div>
        </div>
        
        {/* Progress Indicator */}
        <div className="flex items-center gap-4">
          <div className="text-sm text-muted-foreground">
            Step 3 of 12
          </div>
          <div className="w-32 h-2 bg-muted rounded-full overflow-hidden">
            <div className="h-full bg-primary w-1/4" />
          </div>
          <Badge>{exercise.difficulty}</Badge>
        </div>
      </header>

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Collapsible Sidebar */}
        <aside
          className={cn(
            "border-r bg-muted/30 transition-all duration-300 overflow-y-auto",
            sidebarOpen ? "w-96" : "w-0"
          )}
        >
          {sidebarOpen && (
            <div className="p-6">
              <div className="flex items-center justify-between mb-4">
                <h2 className="text-lg font-semibold flex items-center gap-2">
                  <BookOpen className="w-5 h-5" />
                  Lesson Content
                </h2>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={() => setIsPinned(!isPinned)}
                >
                  {isPinned ? '📌' : '📍'} {isPinned ? 'Pinned' : 'Pin'}
                </Button>
              </div>
              <LessonContent lesson={lesson} exercise={exercise} />
            </div>
          )}
        </aside>

        {/* Main Editor Area */}
        <main className="flex-1 overflow-hidden">
          <CodeEditor 
            exercise={exercise} 
            onNeedHelp={() => setSidebarOpen(true)}
          />
        </main>
      </div>
    </div>
  )
}
```

---

## Layout Option 4: Storytelling Flow (Duolingo-Style)

**Best For:** Absolute beginners, gamified learning

```
┌─────────────────────────────────────────────────────────────┐
│  [←]  Lesson 3: Working with Lists          Progress: 67%   │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│         ┌───────────────────────────────┐                   │
│         │  📖 Let's Learn About Lists   │                   │
│         │                               │                   │
│         │  In Python, a list is a       │                   │
│         │  collection of items...       │                   │
│         └───────────────────────────────┘                   │
│                      ↓                                      │
│         ┌───────────────────────────────┐                   │
│         │  💡 Example                   │                   │
│         │  fruits = ["apple", "banana"] │                   │
│         └───────────────────────────────┘                   │
│                      ↓                                      │
│         ┌───────────────────────────────┐                   │
│         │  ✍️ Now You Try!              │                   │
│         │  Create a list of 3 colors:   │                   │
│         │                               │                   │
│         │  [  Your Code Here  ]         │                   │
│         │                               │                   │
│         │  [Check Answer]               │                   │
│         └───────────────────────────────┘                   │
│                                                             │
│                           [Continue →]                      │
└─────────────────────────────────────────────────────────────┘
```

**Pros:**
- Extremely beginner-friendly
- Clear progression
- Feels less intimidating
- Great engagement

**Cons:**
- Not efficient for complex problems
- Can feel slow for advanced users
- Requires more content creation

### Implementation:

```typescript
// components/layouts/storytelling-flow-layout.tsx
'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Progress } from '@/components/ui/progress'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

interface Step {
  type: 'explanation' | 'example' | 'exercise'
  content: any
}

export default function StorytellingFlowLayout({ 
  lesson, 
  steps 
}: { 
  lesson: any
  steps: Step[] 
}) {
  const [currentStep, setCurrentStep] = useState(0)
  const [userAnswer, setUserAnswer] = useState('')

  const step = steps[currentStep]
  const progress = ((currentStep + 1) / steps.length) * 100

  const handleNext = () => {
    if (currentStep < steps.length - 1) {
      setCurrentStep(currentStep + 1)
      setUserAnswer('')
    }
  }

  const handlePrev = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1)
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-50 dark:from-gray-900 dark:to-gray-800">
      {/* Header */}
      <header className="border-b bg-background/95 backdrop-blur">
        <div className="container mx-auto px-6 py-4 flex items-center justify-between">
          <Button variant="ghost" size="icon" onClick={handlePrev}>
            <ChevronLeft className="w-5 h-5" />
          </Button>
          
          <div className="flex-1 max-w-md mx-8">
            <div className="flex items-center justify-between mb-2">
              <span className="text-sm font-medium">{lesson.title}</span>
              <span className="text-sm text-muted-foreground">
                {currentStep + 1} / {steps.length}
              </span>
            </div>
            <Progress value={progress} className="h-2" />
          </div>

          <Button variant="ghost" size="icon">
            <span className="text-xl">✕</span>
          </Button>
        </div>
      </header>

      {/* Main Content */}
      <main className="container mx-auto px-6 py-12">
        <AnimatePresence mode="wait">
          <motion.div
            key={currentStep}
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
            transition={{ duration: 0.3 }}
            className="max-w-3xl mx-auto"
          >
            {step.type === 'explanation' && (
              <ExplanationStep content={step.content} />
            )}
            
            {step.type === 'example' && (
              <ExampleStep content={step.content} />
            )}
            
            {step.type === 'exercise' && (
              <ExerciseStep 
                content={step.content}
                answer={userAnswer}
                onAnswerChange={setUserAnswer}
              />
            )}
          </motion.div>
        </AnimatePresence>
      </main>

      {/* Footer Navigation */}
      <footer className="fixed bottom-0 left-0 right-0 border-t bg-background/95 backdrop-blur">
        <div className="container mx-auto px-6 py-4 flex justify-between items-center">
          <Button
            variant="outline"
            onClick={handlePrev}
            disabled={currentStep === 0}
          >
            <ChevronLeft className="w-4 h-4 mr-2" />
            Back
          </Button>

          <div className="text-sm text-muted-foreground">
            Press Enter or click Continue ↵
          </div>

          <Button
            onClick={handleNext}
            disabled={currentStep === steps.length - 1}
          >
            Continue
            <ChevronRight className="w-4 h-4 ml-2" />
          </Button>
        </div>
      </footer>
    </div>
  )
}

// Sub-components for different step types
function ExplanationStep({ content }: { content: any }) {
  return (
    <Card className="p-12 text-center">
      <div className="text-6xl mb-6">{content.icon}</div>
      <h2 className="text-3xl font-bold mb-4">{content.title}</h2>
      <p className="text-lg text-muted-foreground leading-relaxed max-w-2xl mx-auto">
        {content.description}
      </p>
    </Card>
  )
}

function ExampleStep({ content }: { content: any }) {
  return (
    <Card className="p-8">
      <div className="flex items-center gap-3 mb-6">
        <div className="text-3xl">💡</div>
        <h3 className="text-xl font-semibold">Example</h3>
      </div>
      <div className="bg-black text-green-400 p-6 rounded-lg font-mono text-sm mb-4">
        <pre>{content.code}</pre>
      </div>
      <p className="text-muted-foreground">{content.explanation}</p>
    </Card>
  )
}

function ExerciseStep({ 
  content, 
  answer, 
  onAnswerChange 
}: { 
  content: any
  answer: string
  onAnswerChange: (val: string) => void
}) {
  const [showFeedback, setShowFeedback] = useState(false)
  const [isCorrect, setIsCorrect] = useState(false)

  const handleCheck = () => {
    const correct = answer.trim() === content.solution.trim()
    setIsCorrect(correct)
    setShowFeedback(true)
  }

  return (
    <Card className="p-8">
      <div className="flex items-center gap-3 mb-6">
        <div className="text-3xl">✍️</div>
        <h3 className="text-xl font-semibold">Your Turn!</h3>
      </div>
      
      <p className="text-lg mb-6">{content.prompt}</p>

      <textarea
        value={answer}
        onChange={(e) => onAnswerChange(e.target.value)}
        className="w-full h-32 p-4 font-mono text-sm border rounded-lg mb-4"
        placeholder="Type your code here..."
      />

      <Button onClick={handleCheck} className="w-full" size="lg">
        Check Answer
      </Button>

      {showFeedback && (
        <motion.div
          initial={{ opacity: 0, y: 10 }}
          animate={{ opacity: 1, y: 0 }}
          className={cn(
            "mt-4 p-4 rounded-lg",
            isCorrect ? "bg-green-50 border-green-200" : "bg-red-50 border-red-200"
          )}
        >
          <div className="flex items-center gap-2">
            <span className="text-2xl">{isCorrect ? '✅' : '❌'}</span>
            <span className="font-semibold">
              {isCorrect ? 'Correct! Well done!' : 'Not quite right. Try again!'}
            </span>
          </div>
          {!isCorrect && content.hint && (
            <p className="text-sm text-muted-foreground mt-2">
              💡 Hint: {content.hint}
            </p>
          )}
        </motion.div>
      )}
    </Card>
  )
}
```

---

## Layout Option 5: Video-Integrated Layout

**Best For:** Video courses, instructor-led learning

```
┌─────────────────────────────────────────────────────────────┐
│  Header: Course Title, Instructor                           │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌─────────────────────────────────────┐                   │
│  │                                     │                   │
│  │         VIDEO PLAYER                │                   │
│  │      (Can be minimized)             │                   │
│  │                                     │                   │
│  └─────────────────────────────────────┘                   │
│                                                             │
│  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━ 3:24 / 12:45        │
│                                                             │
│  ┌─────────────────────────────────────────────────────┐   │
│  │ 📝 Follow Along:                                    │   │
│  │                                                     │   │
│  │  [Code Editor matching video]                      │   │
│  │                                                     │   │
│  └─────────────────────────────────────────────────────┘   │
│                                                             │
│  [⏮ Previous]  [▶ Play/Pause]  [⏭ Next Section]          │
└─────────────────────────────────────────────────────────────┘
```

---

## Complete Lesson Content Component

```typescript
// components/lesson-content.tsx
'use client'

import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { BookOpen, Target, Lightbulb, Code2 } from 'lucide-react'

export default function LessonContent({ lesson, exercise }: any) {
  return (
    <div className="space-y-6">
      {/* Learning Objectives */}
      <Card className="p-6 bg-blue-50 dark:bg-blue-950 border-blue-200">
        <div className="flex items-center gap-2 mb-4">
          <Target className="w-5 h-5 text-blue-600" />
          <h3 className="font-semibold text-blue-900 dark:text-blue-100">
            Learning Objectives
          </h3>
        </div>
        <ul className="space-y-2">
          {lesson.objectives?.map((obj: string, idx: number) => (
            <li key={idx} className="flex items-start gap-2 text-sm">
              <span className="text-blue-600">✓</span>
              <span>{obj}</span>
            </li>
          ))}
        </ul>
      </Card>

      {/* Main Content */}
      <div className="prose prose-sm max-w-none">
        <h2 className="text-2xl font-bold mb-4">{lesson.title}</h2>
        <div dangerouslySetInnerHTML={{ __html: lesson.content }} />
      </div>

      {/* Code Examples */}
      {lesson.examples && lesson.examples.length > 0 && (
        <Card className="p-6">
          <div className="flex items-center gap-2 mb-4">
            <Code2 className="w-5 h-5" />
            <h3 className="font-semibold">Examples</h3>
          </div>
          <Tabs defaultValue="example-0">
            <TabsList>
              {lesson.examples.map((ex: any, idx: number) => (
                <TabsTrigger key={idx} value={`example-${idx}`}>
                  Example {idx + 1}
                </TabsTrigger>
              ))}
            </TabsList>
            {lesson.examples.map((example: any, idx: number) => (
              <TabsContent key={idx} value={`example-${idx}`} className="space-y-3">
                <p className="text-sm text-muted-foreground">{example.description}</p>
                <pre className="bg-black text-green-400 p-4 rounded-lg overflow-x-auto">
                  <code>{example.code}</code>
                </pre>
                {example.output && (
                  <div>
                    <p className="text-sm font-medium mb-2">Output:</p>
                    <pre className="bg-muted p-3 rounded text-sm">
                      {example.output}
                    </pre>
                  </div>
                )}
              </TabsContent>
            ))}
          </Tabs>
        </Card>
      )}

      {/* Key Concepts */}
      {lesson.keyConcepts && (
        <Card className="p-6 bg-amber-50 dark:bg-amber-950 border-amber-200">
          <div className="flex items-center gap-2 mb-4">
            <Lightbulb className="w-5 h-5 text-amber-600" />
            <h3 className="font-semibold text-amber-900 dark:text-amber-100">
              Key Concepts
            </h3>
          </div>
          <ul className="space-y-2">
            {lesson.keyConcepts.map((concept: string, idx: number) => (
              <li key={idx} className="text-sm flex items-start gap-2">
                <span className="text-amber-600">💡</span>
                <span>{concept}</span>
              </li>
            ))}
          </ul>
        </Card>
      )}

      <Separator />

      {/* Exercise Instructions */}
      <div>
        <h3 className="text-xl font-bold mb-4">Exercise: {exercise.title}</h3>
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <Badge variant={
              exercise.difficulty === 'BEGINNER' ? 'default' :
              exercise.difficulty === 'INTERMEDIATE' ? 'secondary' : 'destructive'
            }>
              {exercise.difficulty}
            </Badge>
            <span className="text-sm text-muted-foreground">
              {exercise.points} points
            </span>
          </div>
          
          <p className="text-muted-foreground whitespace-pre-wrap">
            {exercise.description}
          </p>

          {/* Constraints */}
          {exercise.constraints && (
            <div className="border-l-4 border-blue-500 pl-4">
              <h4 className="font-semibold text-sm mb-2">Constraints:</h4>
              <ul className="text-sm space-y-1">
                {exercise.constraints.map((constraint: string, idx: number) => (
                  <li key={idx} className="text-muted-foreground">• {constraint}</li>
                ))}
              </ul>
            </div>
          )}

          {/* Test Cases Preview */}
          {exercise.testCases && (
            <div className="border-l-4 border-green-500 pl-4">
              <h4 className="font-semibold text-sm mb-2">Example Test Cases:</h4>
              {exercise.testCases.slice(0, 2).map((test: any, idx: number) => (
                <div key={idx} className="text-sm mb-2">
                  <div className="font-mono text-xs bg-muted p-2 rounded">
                    <div><strong>Input:</strong> {test.input || '(none)'}</div>
                    <div><strong>Expected Output:</strong> {test.expectedOutput}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
```

---

## Responsive Breakpoints Strategy

```typescript
// tailwind.config.ts
export default {
  theme: {
    screens: {
      'xs': '475px',    // Small phones
      'sm': '640px',    // Phones landscape
      'md': '768px',    // Tablets
      'lg': '1024px',   // Laptop
      'xl': '1280px',   // Desktop
      '2xl': '1536px',  // Large Desktop
    }
  }
}

// Responsive behavior:
// xs-sm: Tabbed view (mobile)
// md-lg: Collapsible sidebar
// xl+: Split view with both panels
```

---

## Recommended Layout by Use Case

| Use Case | Recommended Layout | Alternative |
|----------|-------------------|-------------|
| **Algorithmic Challenges** | Split View | Collapsible Sidebar |
| **Tutorial Learning** | Storytelling Flow | Tabbed View |
| **Video Courses** | Video-Integrated | Collapsible Sidebar |
| **Interactive Coding** | Collapsible Sidebar | Split View |
| **Mobile Learning** | Tabbed View | Storytelling Flow |
| **Assessment/Testing** | Split View | Tabbed View |

---

## Accessibility Considerations

```typescript
// Add keyboard shortcuts
useEffect(() => {
  const handleKeyPress = (e: KeyboardEvent) => {
    // Ctrl/Cmd + Enter: Run code
    if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
      handleRunCode()
    }
    // Ctrl/Cmd + S: Submit
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
      e.preventDefault()
      handleSubmit()
    }
    // Ctrl/Cmd + H: Toggle hints
    if ((e.ctrlKey || e.metaKey) && e.key === 'h') {
      toggleHints()
    }
  }

  window.addEventListener('keydown', handleKeyPress)
  return () => window.removeEventListener('keydown', handleKeyPress)
}, [])

// Screen reader announcements
const announceResult = (result: any) => {
  const announcement = result.allPassed 
    ? `Success! All ${result.testResults.length} tests passed.`
    : `${result.testResults.filter((t: any) => t.passed).length} of ${result.testResults.length} tests passed.`
  
  // Create live region for screen readers
  const liveRegion = document.createElement('div')
  liveRegion.setAttribute('role', 'status')
  liveRegion.setAttribute('aria-live', 'polite')
  liveRegion.textContent = announcement
  document.body.appendChild(liveRegion)
  setTimeout(() => liveRegion.remove(), 1000)
}
```

---

## My Recommendation: Collapsible Sidebar Layout

**Why it's best for your platform:**

1. ✅ **Code-first approach**: Students spend 80% of time coding
2. ✅ **Lesson always accessible**: One click away, not hidden
3. ✅ **Progressive disclosure**: Read → Code → Test → Review
4. ✅ **Responsive**: Works on tablets and up
5. ✅ **Industry standard**: Familiar pattern (VS Code sidebar)
6. ✅ **Flexible**: Can pin open for reference

**Implementation Priority:**
1. Start with Collapsible Sidebar (main experience)
2. Add Tabbed View for mobile
3. Later: Add Storytelling Flow for absolute beginners
4. Optional: Split View for power users

Would you like me to create the complete implementation code for the Collapsible Sidebar layout with all the polish and features?
