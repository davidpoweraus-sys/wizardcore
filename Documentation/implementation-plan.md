# Coding Platform Implementation Plan

## Project Overview
A full-stack web application for coding education with AI-generated exercises, code execution via Judge0, and automated grading.

**Tech Stack:** Next.js 14+, TypeScript, Tailwind CSS, PostgreSQL, Prisma, Judge0, OpenAI/Anthropic

**Timeline:** 4 weeks (adjustable based on experience level)

---

## Phase 1: Local Development Foundation (Week 1)
**Goal:** Get development environment running with basic Next.js app

### Task 1.1: Environment Setup (Day 1)
**Deliverables:**
- [ ] Node.js LTS installed (v18+ recommended)
- [ ] Git configured
- [ ] VS Code with extensions: ESLint, Prettier, Prisma

**Commands:**
```bash
node --version  # Verify 18+
npm --version
git --version
```

### Task 1.2: Create Next.js Project (Day 1)
**Deliverables:**
- [ ] Next.js app initialized with TypeScript
- [ ] Tailwind CSS configured
- [ ] App Router structure created
- [ ] Dev server runs successfully

**Commands:**
```bash
npx create-next-app@latest my-coding-platform \
  --typescript \
  --tailwind \
  --app \
  --no-turbopack \
  --import-alias "@/*"

cd my-coding-platform
npm run dev
```

**Verification:** Visit http://localhost:3000 and see default Next.js page

### Task 1.3: UI Framework Setup (Day 2)
**Deliverables:**
- [ ] ShadCN UI initialized
- [ ] Core components installed
- [ ] Theme configured (light/dark mode)

**Commands:**
```bash
npx shadcn@latest init
# Choose: Default style, Zinc base color, CSS variables

npx shadcn@latest add button card dialog textarea input label
npx shadcn@latest add tabs select badge alert
```

**Test Component:** Create a test page with all components to verify styling

### Task 1.4: Database Setup (Days 3-4)
**Deliverables:**
- [ ] PostgreSQL running (Docker recommended)
- [ ] Database connection verified
- [ ] Prisma ORM installed and initialized

**Docker Method (Recommended):**
```bash
# Install Docker Desktop first, then:
docker run --name code-platform-db \
  -e POSTGRES_PASSWORD=dev_password_123 \
  -e POSTGRES_DB=coding_platform \
  -p 5432:5432 \
  -d postgres:16-alpine

# Verify it's running
docker ps
```

**Prisma Setup:**
```bash
npm install prisma @prisma/client
npx prisma init
```

**Configure `.env`:**
```env
DATABASE_URL="postgresql://postgres:dev_password_123@localhost:5432/coding_platform?schema=public"
```

**Verification:** Run `npx prisma studio` and see empty database UI

### Task 1.5: Project Structure (Day 4)
**Deliverables:**
- [ ] Folder structure organized
- [ ] Basic routing setup
- [ ] Layout components created

**Directory Structure:**
```
my-coding-platform/
├── app/
│   ├── (auth)/
│   │   ├── login/
│   │   └── register/
│   ├── (dashboard)/
│   │   ├── page.tsx          # Main dashboard
│   │   └── layout.tsx
│   ├── exercise/
│   │   └── [id]/
│   │       └── page.tsx      # Exercise workspace
│   ├── admin/
│   │   └── page.tsx          # Admin panel
│   ├── api/
│   │   ├── exercises/
│   │   ├── submit/
│   │   └── generate/
│   └── layout.tsx
├── components/
│   ├── ui/                   # ShadCN components
│   ├── code-editor.tsx
│   ├── exercise-card.tsx
│   └── submission-result.tsx
├── lib/
│   ├── prisma.ts             # Prisma client singleton
│   ├── judge0.ts             # Judge0 API wrapper
│   └── ai.ts                 # AI service wrapper
└── prisma/
    └── schema.prisma
```

---

## Phase 2: Backend & Data Layer (Week 2)
**Goal:** Build database schema and core API routes

### Task 2.1: Database Schema Design (Days 5-6)
**Deliverables:**
- [ ] Complete Prisma schema
- [ ] Migration files created
- [ ] Database seeded with sample data

**Schema (`prisma/schema.prisma`):**
```prisma
generator client {
  provider = "prisma-client-js"
}

datasource db {
  provider = "postgresql"
  url      = env("DATABASE_URL")
}

model User {
  id            String       @id @default(cuid())
  email         String       @unique
  name          String?
  passwordHash  String
  role          Role         @default(STUDENT)
  createdAt     DateTime     @default(now())
  updatedAt     DateTime     @updatedAt
  submissions   Submission[]
  progress      Progress[]
}

enum Role {
  STUDENT
  TEACHER
  ADMIN
}

model Syllabus {
  id          String     @id @default(cuid())
  title       String
  description String?
  language    String     // "python", "javascript", etc.
  level       Level      @default(BEGINNER)
  order       Int
  exercises   Exercise[]
  createdAt   DateTime   @default(now())
  updatedAt   DateTime   @updatedAt
}

enum Level {
  BEGINNER
  INTERMEDIATE
  ADVANCED
}

model Exercise {
  id              String       @id @default(cuid())
  syllabusId      String
  syllabus        Syllabus     @relation(fields: [syllabusId], references: [id], onDelete: Cascade)
  title           String
  description     String       @db.Text
  difficulty      Level        @default(BEGINNER)
  language        String       // Programming language
  starterCode     String?      @db.Text
  solution        String       @db.Text
  testCases       Json         // Array of {input, expectedOutput}
  hints           Json?        // Array of hint strings
  order           Int
  points          Int          @default(100)
  timeLimit       Int          @default(2000) // milliseconds
  memoryLimit     Int          @default(128000) // kilobytes
  createdAt       DateTime     @default(now())
  updatedAt       DateTime     @updatedAt
  submissions     Submission[]
  
  @@index([syllabusId, order])
}

model Submission {
  id            String    @id @default(cuid())
  userId        String
  user          User      @relation(fields: [userId], references: [id])
  exerciseId    String
  exercise      Exercise  @relation(fields: [exerciseId], references: [id], onDelete: Cascade)
  code          String    @db.Text
  language      String
  status        SubmissionStatus
  output        String?   @db.Text
  error         String?   @db.Text
  executionTime Int?      // milliseconds
  memoryUsed    Int?      // kilobytes
  testResults   Json      // Array of test case results
  score         Int       @default(0)
  createdAt     DateTime  @default(now())
  
  @@index([userId, exerciseId])
  @@index([createdAt])
}

enum SubmissionStatus {
  PENDING
  RUNNING
  ACCEPTED
  WRONG_ANSWER
  TIME_LIMIT_EXCEEDED
  MEMORY_LIMIT_EXCEEDED
  RUNTIME_ERROR
  COMPILATION_ERROR
}

model Progress {
  id          String   @id @default(cuid())
  userId      String
  user        User     @relation(fields: [userId], references: [id])
  exerciseId  String
  completed   Boolean  @default(false)
  attempts    Int      @default(0)
  bestScore   Int      @default(0)
  lastAttempt DateTime @default(now())
  
  @@unique([userId, exerciseId])
}
```

**Run Migration:**
```bash
npx prisma migrate dev --name init
npx prisma generate
```

### Task 2.2: Prisma Client Setup (Day 6)
**Deliverables:**
- [ ] Singleton Prisma client created
- [ ] Type-safe database queries working

**Create `lib/prisma.ts`:**
```typescript
import { PrismaClient } from '@prisma/client'

const globalForPrisma = globalThis as unknown as {
  prisma: PrismaClient | undefined
}

export const prisma = globalForPrisma.prisma ?? new PrismaClient({
  log: process.env.NODE_ENV === 'development' ? ['query', 'error', 'warn'] : ['error'],
})

if (process.env.NODE_ENV !== 'production') globalForPrisma.prisma = prisma
```

### Task 2.3: Seed Database (Day 7)
**Deliverables:**
- [ ] Seed script created
- [ ] Sample data populated

**Create `prisma/seed.ts`:**
```typescript
import { PrismaClient, Level } from '@prisma/client'
import bcrypt from 'bcryptjs'

const prisma = new PrismaClient()

async function main() {
  // Create test user
  const hashedPassword = await bcrypt.hash('password123', 10)
  
  const user = await prisma.user.create({
    data: {
      email: 'student@test.com',
      name: 'Test Student',
      passwordHash: hashedPassword,
      role: 'STUDENT',
    },
  })

  // Create Python syllabus
  const pythonSyllabus = await prisma.syllabus.create({
    data: {
      title: 'Python Fundamentals',
      description: 'Learn Python programming from scratch',
      language: 'python',
      level: Level.BEGINNER,
      order: 1,
    },
  })

  // Create sample exercise
  await prisma.exercise.create({
    data: {
      syllabusId: pythonSyllabus.id,
      title: 'Hello World',
      description: 'Write a program that prints "Hello, World!"',
      difficulty: Level.BEGINNER,
      language: 'python',
      starterCode: '# Write your code here\n',
      solution: 'print("Hello, World!")',
      testCases: [
        { input: '', expectedOutput: 'Hello, World!\n' }
      ],
      order: 1,
      points: 100,
    },
  })

  console.log('Database seeded successfully!')
}

main()
  .catch((e) => {
    console.error(e)
    process.exit(1)
  })
  .finally(async () => {
    await prisma.$disconnect()
  })
```

**Add to `package.json`:**
```json
{
  "prisma": {
    "seed": "ts-node --compiler-options {\"module\":\"CommonJS\"} prisma/seed.ts"
  }
}
```

**Install dependencies and run:**
```bash
npm install -D ts-node
npm install bcryptjs
npm install -D @types/bcryptjs
npx prisma db seed
```

### Task 2.4: Exercises API Route (Days 8-9)
**Deliverables:**
- [ ] GET /api/exercises - List all exercises
- [ ] GET /api/exercises/[id] - Get single exercise
- [ ] POST /api/exercises - Create exercise (admin only)

**Create `app/api/exercises/route.ts`:**
```typescript
import { NextRequest, NextResponse } from 'next/server'
import { prisma } from '@/lib/prisma'

export async function GET(request: NextRequest) {
  try {
    const searchParams = request.nextUrl.searchParams
    const syllabusId = searchParams.get('syllabusId')
    const language = searchParams.get('language')

    const exercises = await prisma.exercise.findMany({
      where: {
        ...(syllabusId && { syllabusId }),
        ...(language && { language }),
      },
      include: {
        syllabus: true,
        _count: {
          select: { submissions: true },
        },
      },
      orderBy: [
        { syllabus: { order: 'asc' } },
        { order: 'asc' },
      ],
    })

    return NextResponse.json({ exercises })
  } catch (error) {
    console.error('Error fetching exercises:', error)
    return NextResponse.json(
      { error: 'Failed to fetch exercises' },
      { status: 500 }
    )
  }
}

export async function POST(request: NextRequest) {
  try {
    // TODO: Add authentication check here
    const body = await request.json()
    
    const exercise = await prisma.exercise.create({
      data: {
        syllabusId: body.syllabusId,
        title: body.title,
        description: body.description,
        difficulty: body.difficulty,
        language: body.language,
        starterCode: body.starterCode,
        solution: body.solution,
        testCases: body.testCases,
        hints: body.hints,
        order: body.order,
        points: body.points || 100,
        timeLimit: body.timeLimit || 2000,
        memoryLimit: body.memoryLimit || 128000,
      },
    })

    return NextResponse.json({ exercise }, { status: 201 })
  } catch (error) {
    console.error('Error creating exercise:', error)
    return NextResponse.json(
      { error: 'Failed to create exercise' },
      { status: 500 }
    )
  }
}
```

**Create `app/api/exercises/[id]/route.ts`:**
```typescript
import { NextRequest, NextResponse } from 'next/server'
import { prisma } from '@/lib/prisma'

export async function GET(
  request: NextRequest,
  { params }: { params: { id: string } }
) {
  try {
    const exercise = await prisma.exercise.findUnique({
      where: { id: params.id },
      include: {
        syllabus: true,
      },
    })

    if (!exercise) {
      return NextResponse.json(
        { error: 'Exercise not found' },
        { status: 404 }
      )
    }

    // Don't send solution to frontend
    const { solution, ...exerciseWithoutSolution } = exercise

    return NextResponse.json({ exercise: exerciseWithoutSolution })
  } catch (error) {
    console.error('Error fetching exercise:', error)
    return NextResponse.json(
      { error: 'Failed to fetch exercise' },
      { status: 500 }
    )
  }
}
```

### Task 2.5: Submit API Route (Days 9-10)
**Deliverables:**
- [ ] POST /api/submit - Core submission logic
- [ ] Test case execution
- [ ] Grading logic
- [ ] Database persistence

**Create `app/api/submit/route.ts`:**
```typescript
import { NextRequest, NextResponse } from 'next/server'
import { prisma } from '@/lib/prisma'
import { executeCode } from '@/lib/judge0'

interface TestCase {
  input: string
  expectedOutput: string
}

export async function POST(request: NextRequest) {
  try {
    const { exerciseId, code, userId } = await request.json()

    // Fetch exercise with test cases
    const exercise = await prisma.exercise.findUnique({
      where: { id: exerciseId },
    })

    if (!exercise) {
      return NextResponse.json(
        { error: 'Exercise not found' },
        { status: 404 }
      )
    }

    const testCases = exercise.testCases as TestCase[]
    const results = []
    let passedTests = 0

    // Execute code against each test case
    for (const testCase of testCases) {
      const result = await executeCode({
        code,
        language: exercise.language,
        stdin: testCase.input,
        timeLimit: exercise.timeLimit,
        memoryLimit: exercise.memoryLimit,
      })

      const passed = result.stdout?.trim() === testCase.expectedOutput.trim()
      if (passed) passedTests++

      results.push({
        input: testCase.input,
        expectedOutput: testCase.expectedOutput,
        actualOutput: result.stdout || '',
        passed,
        error: result.stderr || result.compile_output,
        executionTime: result.time,
        memoryUsed: result.memory,
      })
    }

    const score = Math.round((passedTests / testCases.length) * exercise.points)
    const allPassed = passedTests === testCases.length

    // Save submission
    const submission = await prisma.submission.create({
      data: {
        userId,
        exerciseId,
        code,
        language: exercise.language,
        status: allPassed ? 'ACCEPTED' : 'WRONG_ANSWER',
        testResults: results,
        score,
      },
    })

    // Update progress
    await prisma.progress.upsert({
      where: {
        userId_exerciseId: {
          userId,
          exerciseId,
        },
      },
      update: {
        attempts: { increment: 1 },
        bestScore: { set: score },
        completed: allPassed,
        lastAttempt: new Date(),
      },
      create: {
        userId,
        exerciseId,
        attempts: 1,
        bestScore: score,
        completed: allPassed,
      },
    })

    return NextResponse.json({
      submission: {
        id: submission.id,
        status: submission.status,
        score,
        testResults: results,
        allPassed,
      },
    })
  } catch (error) {
    console.error('Error processing submission:', error)
    return NextResponse.json(
      { error: 'Failed to process submission' },
      { status: 500 }
    )
  }
}
```

---

## Phase 3: External Services Integration (Week 3)
**Goal:** Connect Judge0 and AI services

### Task 3.1: Judge0 Self-Hosted Setup (Days 11-13)
**Goal:** Run your own Judge0 instance with full control

**Deliverables:**
- [ ] Judge0 backend running locally via Docker
- [ ] Judge0 IDE embedded in frontend
- [ ] API wrapper created
- [ ] Test execution working

**Prerequisites:**
- Docker Desktop installed and running
- At least 4GB RAM available
- Ports 2358 (API) and 3001 (IDE) available

#### Step 3.1.1: Clone and Configure Judge0 (Day 11)

**Clone the Repository:**
```bash
cd ~/projects  # or your preferred location
git clone https://github.com/judge0/judge0.git
cd judge0
```

**Review the Architecture:**
Judge0 consists of:
- **PostgreSQL**: Stores submissions and results
- **Redis**: Job queue management
- **Judge0 Server**: Executes code in isolated containers
- **Judge0 Workers**: Process submissions from the queue

**Configure Environment (Optional):**
Create a `.env.prod` file if you need custom settings:
```bash
# judge0/.env.prod
REDIS_HOST=redis
REDIS_PORT=6379
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=judge0
POSTGRES_USER=judge0
POSTGRES_PASSWORD=YourSecurePassword123

# Increase workers for better performance
WORKERS_COUNT=2

# Enable more languages if needed
ENABLE_WAIT_RESULT=true
ENABLE_COMPILER_OPTIONS=true
```

#### Step 3.1.2: Start Judge0 Services (Day 11)

**Start Core Services First:**
```bash
# Start database and Redis
docker-compose up -d db redis

# Wait 10 seconds for services to initialize
sleep 10

# Check they're running
docker-compose ps
```

**Start Judge0 Server and Workers:**
```bash
# Build and start all services
docker-compose up -d

# View logs to confirm it's working
docker-compose logs -f server

# You should see: "Judge0 v1.13.0 started"
```

**Verify Installation:**
```bash
# Test the API
curl http://localhost:2358/about

# Should return JSON with Judge0 version info
```

**Common Issues:**
```bash
# If port 2358 is in use:
lsof -ti:2358 | xargs kill -9

# If containers won't start:
docker-compose down -v  # Remove volumes
docker-compose up -d

# Check logs for specific service:
docker-compose logs server
docker-compose logs workers
```

#### Step 3.1.3: Set Up Judge0 IDE (Day 12)

**Option A: Use Judge0 IDE Embed (Simplest)**

Judge0 provides an embeddable IDE at https://ide.judge0.com that can connect to your self-hosted backend.

**Create `components/judge0-ide.tsx`:**
```typescript
'use client'

import { useEffect, useRef } from 'react'

interface Judge0IDEProps {
  language?: string
  code?: string
  onRun?: (result: any) => void
}

export default function Judge0IDE({ 
  language = 'python',
  code = '',
  onRun 
}: Judge0IDEProps) {
  const iframeRef = useRef<HTMLIFrameElement>(null)

  useEffect(() => {
    // Listen for messages from IDE
    const handleMessage = (event: MessageEvent) => {
      if (event.data.type === 'judge0-result') {
        onRun?.(event.data.result)
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [onRun])

  const iframeUrl = `https://ide.judge0.com/embed?` + new URLSearchParams({
    language,
    theme: 'dark',
    hideRun: 'false',
    // Point to your self-hosted Judge0
    judge0Url: 'http://localhost:2358'
  }).toString()

  return (
    <iframe
      ref={iframeRef}
      src={iframeUrl}
      className="w-full h-[600px] border rounded-lg"
      sandbox="allow-scripts allow-same-origin allow-forms"
    />
  )
}
```

**Option B: Self-Host Judge0 IDE (Full Control)**

```bash
# Clone IDE repository
cd ~/projects
git clone https://github.com/judge0/ide.git judge0-ide
cd judge0-ide

# Install dependencies
npm install

# Configure to use your Judge0 backend
# Edit public/config.js or create .env.local:
echo "REACT_APP_JUDGE0_API_URL=http://localhost:2358" > .env.local

# Start IDE
npm start
# Runs on http://localhost:3001
```

**Integrate Self-Hosted IDE:**
```typescript
// components/judge0-ide.tsx (self-hosted version)
'use client'

export default function Judge0IDE({ language = 'python', code = '' }) {
  return (
    <iframe
      src={`http://localhost:3001?language=${language}`}
      className="w-full h-[600px] border rounded-lg"
      allow="clipboard-read; clipboard-write"
    />
  )
}
```

**Option C: Build Custom Editor with Judge0 API (Most Flexible)**

Use Monaco Editor + direct Judge0 API calls for complete customization:

```typescript
// components/custom-judge0-editor.tsx
'use client'

import { useState } from 'react'
import Editor from '@monaco-editor/react'
import { Button } from '@/components/ui/button'
import { executeCode } from '@/lib/judge0'

export default function CustomJudge0Editor({ 
  language = 'python',
  initialCode = '',
  onSubmit 
}: any) {
  const [code, setCode] = useState(initialCode)
  const [output, setOutput] = useState('')
  const [isRunning, setIsRunning] = useState(false)

  const handleRun = async () => {
    setIsRunning(true)
    setOutput('')

    try {
      const result = await executeCode({
        code,
        language,
        stdin: '', // Can add input field for this
      })

      setOutput(result.stdout || result.stderr || 'No output')
    } catch (error) {
      setOutput(`Error: ${error}`)
    } finally {
      setIsRunning(false)
    }
  }

  return (
    <div className="space-y-4">
      <Editor
        height="400px"
        language={language}
        value={code}
        onChange={(value) => setCode(value || '')}
        theme="vs-dark"
        options={{
          minimap: { enabled: false },
          fontSize: 14,
        }}
      />
      
      <div className="flex gap-2">
        <Button onClick={handleRun} disabled={isRunning}>
          {isRunning ? 'Running...' : 'Run Code'}
        </Button>
        <Button onClick={() => onSubmit?.(code)} variant="default">
          Submit Solution
        </Button>
      </div>

      {output && (
        <div className="bg-black text-green-400 p-4 rounded font-mono text-sm">
          <pre>{output}</pre>
        </div>
      )}
    </div>
  )
}
```

#### Step 3.1.4: Create Judge0 API Wrapper (Day 13)

#### Step 3.1.4: Create Judge0 API Wrapper (Day 13)

**Enhanced `lib/judge0.ts` for Self-Hosted:**
```typescript
interface ExecuteCodeParams {
  code: string
  language: string
  stdin?: string
  timeLimit?: number // milliseconds
  memoryLimit?: number // kilobytes
  expectedOutput?: string
}

interface ExecutionResult {
  stdout: string | null
  stderr: string | null
  compile_output: string | null
  time: number | null // milliseconds
  memory: number | null // kilobytes
  status: {
    id: number
    description: string
  }
  token?: string
}

const LANGUAGE_IDS: Record<string, number> = {
  python: 71,       // Python 3.8.1
  javascript: 63,   // JavaScript (Node.js 12.14.0)
  typescript: 74,   // TypeScript 3.7.4
  java: 62,         // Java (OpenJDK 13.0.1)
  cpp: 54,          // C++ (GCC 9.2.0)
  c: 50,            // C (GCC 9.2.0)
  csharp: 51,       // C# (Mono 6.6.0.161)
  ruby: 72,         // Ruby 2.7.0
  go: 60,           // Go 1.13.5
  rust: 73,         // Rust 1.40.0
  php: 68,          // PHP 7.4.1
  swift: 83,        // Swift 5.2.3
  kotlin: 78,       // Kotlin 1.3.70
}

export async function executeCode(params: ExecuteCodeParams): Promise<ExecutionResult> {
  // Self-hosted Judge0 (no API key needed)
  const baseUrl = process.env.JUDGE0_BASE_URL || 'http://localhost:2358'
  
  const languageId = LANGUAGE_IDS[params.language.toLowerCase()]
  if (!languageId) {
    throw new Error(`Unsupported language: ${params.language}`)
  }

  // Create submission with wait=true for synchronous execution
  const createResponse = await fetch(`${baseUrl}/submissions?base64_encoded=false&wait=true`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      source_code: params.code,
      language_id: languageId,
      stdin: params.stdin || '',
      expected_output: params.expectedOutput || null,
      cpu_time_limit: (params.timeLimit || 2000) / 1000, // Convert to seconds
      memory_limit: params.memoryLimit || 256000, // 256MB default
    }),
  })

  if (!createResponse.ok) {
    const errorText = await createResponse.text()
    throw new Error(`Judge0 API error: ${createResponse.statusText} - ${errorText}`)
  }

  const result = await createResponse.json()
  
  return {
    stdout: result.stdout,
    stderr: result.stderr,
    compile_output: result.compile_output,
    time: result.time ? parseFloat(result.time) * 1000 : null, // Convert to ms
    memory: result.memory,
    status: result.status,
    token: result.token,
  }
}

// Batch execution for multiple test cases
export async function executeCodeBatch(
  code: string,
  language: string,
  testCases: Array<{ input: string; expectedOutput: string }>
): Promise<ExecutionResult[]> {
  const baseUrl = process.env.JUDGE0_BASE_URL || 'http://localhost:2358'
  const languageId = LANGUAGE_IDS[language.toLowerCase()]

  if (!languageId) {
    throw new Error(`Unsupported language: ${language}`)
  }

  // Create batch submission
  const submissions = testCases.map(testCase => ({
    source_code: code,
    language_id: languageId,
    stdin: testCase.input,
    expected_output: testCase.expectedOutput,
  }))

  const response = await fetch(`${baseUrl}/submissions/batch?base64_encoded=false`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ submissions }),
  })

  if (!response.ok) {
    throw new Error(`Batch submission failed: ${response.statusText}`)
  }

  const tokens = await response.json()
  
  // Poll for results
  const results: ExecutionResult[] = []
  for (const { token } of tokens) {
    const result = await pollSubmission(token)
    results.push(result)
  }

  return results
}

// Poll for submission result (for async execution)
async function pollSubmission(token: string, maxAttempts = 10): Promise<ExecutionResult> {
  const baseUrl = process.env.JUDGE0_BASE_URL || 'http://localhost:2358'
  
  for (let i = 0; i < maxAttempts; i++) {
    const response = await fetch(`${baseUrl}/submissions/${token}?base64_encoded=false`)
    const result = await response.json()

    // Status IDs: 1=In Queue, 2=Processing, 3=Accepted, 4+=Various errors
    if (result.status.id > 2) {
      return {
        stdout: result.stdout,
        stderr: result.stderr,
        compile_output: result.compile_output,
        time: result.time ? parseFloat(result.time) * 1000 : null,
        memory: result.memory,
        status: result.status,
        token: result.token,
      }
    }

    // Wait before polling again
    await new Promise(resolve => setTimeout(resolve, 500))
  }

  throw new Error('Submission timeout - execution took too long')
}

// Get supported languages
export async function getSupportedLanguages() {
  const baseUrl = process.env.JUDGE0_BASE_URL || 'http://localhost:2358'
  const response = await fetch(`${baseUrl}/languages`)
  return response.json()
}

// Get system info
export async function getSystemInfo() {
  const baseUrl = process.env.JUDGE0_BASE_URL || 'http://localhost:2358'
  const response = await fetch(`${baseUrl}/about`)
  return response.json()
}
```

#### Step 3.1.5: Test Judge0 Installation (Day 13)

**Create Test Script `scripts/test-judge0.ts`:**
```typescript
import { executeCode, getSupportedLanguages, getSystemInfo } from '../lib/judge0'

async function testJudge0() {
  console.log('🧪 Testing Judge0 Installation...\n')

  try {
    // Test 1: System Info
    console.log('📋 Test 1: Checking system info...')
    const info = await getSystemInfo()
    console.log('✅ Judge0 version:', info.version)
    console.log('')

    // Test 2: Supported Languages
    console.log('📋 Test 2: Getting supported languages...')
    const languages = await getSupportedLanguages()
    console.log('✅ Found', languages.length, 'languages')
    console.log('')

    // Test 3: Python Hello World
    console.log('📋 Test 3: Running Python code...')
    const pythonResult = await executeCode({
      code: 'print("Hello from Judge0!")',
      language: 'python',
    })
    console.log('Output:', pythonResult.stdout)
    console.log('Status:', pythonResult.status.description)
    console.log('Time:', pythonResult.time, 'ms')
    console.log('')

    // Test 4: JavaScript with Input
    console.log('📋 Test 4: Running JavaScript with input...')
    const jsResult = await executeCode({
      code: `
        const readline = require('readline');
        const rl = readline.createInterface({
          input: process.stdin,
          output: process.stdout
        });
        
        rl.on('line', (line) => {
          console.log('You said: ' + line);
          rl.close();
        });
      `,
      language: 'javascript',
      stdin: 'Hello Judge0!',
    })
    console.log('Output:', jsResult.stdout)
    console.log('')

    // Test 5: Compilation Error
    console.log('📋 Test 5: Testing error handling...')
    const errorResult = await executeCode({
      code: 'print("missing quote)',
      language: 'python',
    })
    console.log('Status:', errorResult.status.description)
    if (errorResult.stderr) {
      console.log('Error:', errorResult.stderr.substring(0, 100))
    }
    console.log('')

    console.log('✅ All tests passed! Judge0 is working correctly.')
  } catch (error) {
    console.error('❌ Test failed:', error)
    process.exit(1)
  }
}

testJudge0()
```

**Add to `package.json`:**
```json
{
  "scripts": {
    "test:judge0": "ts-node scripts/test-judge0.ts"
  }
}
```

**Run Tests:**
```bash
npm run test:judge0
```

**Expected Output:**
```
🧪 Testing Judge0 Installation...

📋 Test 1: Checking system info...
✅ Judge0 version: 1.13.0

📋 Test 2: Getting supported languages...
✅ Found 89 languages

📋 Test 3: Running Python code...
Output: Hello from Judge0!

Status: Accepted
Time: 0.023 ms

📋 Test 4: Running JavaScript with input...
Output: You said: Hello Judge0!

📋 Test 5: Testing error handling...
Status: Compilation Error
Error: SyntaxError: EOL while scanning string literal

✅ All tests passed! Judge0 is working correctly.
```

**Environment Variables (`.env.local`):**
```env
# Self-hosted Judge0
JUDGE0_BASE_URL=http://localhost:2358

# Optional: If using RapidAPI as fallback
JUDGE0_API_KEY=your_rapidapi_key_here
```

**Docker Management Commands:**
```bash
# Start Judge0
cd ~/projects/judge0
docker-compose up -d

# Stop Judge0
docker-compose down

# Stop and remove all data (fresh start)
docker-compose down -v

# View logs
docker-compose logs -f server
docker-compose logs -f workers

# Check resource usage
docker stats

# Restart if needed
docker-compose restart
```
```typescript
interface ExecuteCodeParams {
  code: string
  language: string
  stdin?: string
  timeLimit?: number
  memoryLimit?: number
}

interface ExecutionResult {
  stdout: string | null
  stderr: string | null
  compile_output: string | null
  time: number | null
  memory: number | null
  status: {
    id: number
    description: string
  }
}

const LANGUAGE_IDS: Record<string, number> = {
  python: 71,      // Python 3.8.1
  javascript: 63,  // JavaScript (Node.js 12.14.0)
  java: 62,        // Java (OpenJDK 13.0.1)
  cpp: 54,         // C++ (GCC 9.2.0)
  c: 50,           // C (GCC 9.2.0)
}

export async function executeCode(params: ExecuteCodeParams): Promise<ExecutionResult> {
  const baseUrl = process.env.JUDGE0_BASE_URL || 'https://judge0-ce.p.rapidapi.com'
  const apiKey = process.env.JUDGE0_API_KEY
  
  const languageId = LANGUAGE_IDS[params.language.toLowerCase()]
  if (!languageId) {
    throw new Error(`Unsupported language: ${params.language}`)
  }

  // Create submission
  const createResponse = await fetch(`${baseUrl}/submissions?wait=true`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      ...(apiKey && { 'X-RapidAPI-Key': apiKey }),
    },
    body: JSON.stringify({
      source_code: params.code,
      language_id: languageId,
      stdin: params.stdin || '',
      cpu_time_limit: (params.timeLimit || 2000) / 1000, // Convert to seconds
      memory_limit: params.memoryLimit || 128000,
    }),
  })

  if (!createResponse.ok) {
    throw new Error(`Judge0 API error: ${createResponse.statusText}`)
  }

  const result = await createResponse.json()
  
  return {
    stdout: result.stdout,
    stderr: result.stderr,
    compile_output: result.compile_output,
    time: result.time ? parseFloat(result.time) * 1000 : null, // Convert to ms
    memory: result.memory,
    status: result.status,
  }
}
```

**Environment Variables (`.env.local`):**
```env
# Option A: RapidAPI
JUDGE0_BASE_URL=https://judge0-ce.p.rapidapi.com
JUDGE0_API_KEY=your_rapidapi_key_here

# Option B: Self-hosted
JUDGE0_BASE_URL=http://localhost:2358
```

**Test Judge0:**
```bash
# Create test script: scripts/test-judge0.ts
npm run test:judge0
```

### Task 3.2: AI Service Setup (Days 14-15)
**Deliverables:**
- [ ] AI SDK installed
- [ ] API wrapper created
- [ ] Exercise generation prompt engineered

**Install Dependencies:**
```bash
npm install ai @ai-sdk/openai @ai-sdk/anthropic zod
```

**Create `lib/ai.ts`:**
```typescript
import { generateObject } from 'ai'
import { openai } from '@ai-sdk/openai'
import { z } from 'zod'

const exerciseSchema = z.object({
  title: z.string(),
  description: z.string(),
  difficulty: z.enum(['BEGINNER', 'INTERMEDIATE', 'ADVANCED']),
  language: z.string(),
  starterCode: z.string(),
  solution: z.string(),
  testCases: z.array(z.object({
    input: z.string(),
    expectedOutput: z.string(),
    description: z.string().optional(),
  })),
  hints: z.array(z.string()),
  points: z.number(),
})

export type GeneratedExercise = z.infer<typeof exerciseSchema>

interface GenerateExerciseParams {
  topic: string
  language: string
  difficulty: 'BEGINNER' | 'INTERMEDIATE' | 'ADVANCED'
  concepts?: string[]
}

export async function generateExercise(
  params: GenerateExerciseParams
): Promise<GeneratedExercise> {
  const prompt = `Generate a coding exercise with the following requirements:

Topic: ${params.topic}
Programming Language: ${params.language}
Difficulty: ${params.difficulty}
${params.concepts ? `Concepts to cover: ${params.concepts.join(', ')}` : ''}

Requirements:
1. Create a clear, educational problem statement
2. Include starter code with helpful comments
3. Provide a working solution
4. Create 3-5 test cases with inputs and expected outputs
5. Include 2-3 progressive hints
6. Make test cases comprehensive (edge cases, typical cases)

The exercise should be:
- Appropriate for ${params.difficulty.toLowerCase()} level students
- Have clear, unambiguous requirements
- Test the specified concepts thoroughly
- Include meaningful variable/function names
- Have properly formatted expected outputs (including newlines)`

  const { object } = await generateObject({
    model: openai('gpt-4-turbo'),
    schema: exerciseSchema,
    prompt,
  })

  return object
}
```

### Task 3.3: Generate API Route (Day 15)
**Deliverables:**
- [ ] POST /api/generate - AI exercise generation
- [ ] Automatic database insertion

**Create `app/api/generate/route.ts`:**
```typescript
import { NextRequest, NextResponse } from 'next/server'
import { prisma } from '@/lib/prisma'
import { generateExercise } from '@/lib/ai'

export async function POST(request: NextRequest) {
  try {
    // TODO: Add admin authentication
    const { topic, language, difficulty, syllabusId, concepts } = await request.json()

    // Generate exercise using AI
    const generatedExercise = await generateExercise({
      topic,
      language,
      difficulty,
      concepts,
    })

    // Get the highest order in this syllabus
    const lastExercise = await prisma.exercise.findFirst({
      where: { syllabusId },
      orderBy: { order: 'desc' },
    })

    const nextOrder = (lastExercise?.order || 0) + 1

    // Save to database
    const exercise = await prisma.exercise.create({
      data: {
        syllabusId,
        title: generatedExercise.title,
        description: generatedExercise.description,
        difficulty: generatedExercise.difficulty,
        language: generatedExercise.language,
        starterCode: generatedExercise.starterCode,
        solution: generatedExercise.solution,
        testCases: generatedExercise.testCases,
        hints: generatedExercise.hints,
        points: generatedExercise.points,
        order: nextOrder,
      },
    })

    return NextResponse.json({ exercise }, { status: 201 })
  } catch (error) {
    console.error('Error generating exercise:', error)
    return NextResponse.json(
      { error: 'Failed to generate exercise' },
      { status: 500 }
    )
  }
}
```

**Environment Variables:**
```env
OPENAI_API_KEY=sk-...
# OR
ANTHROPIC_API_KEY=sk-ant-...
```

---

## Phase 4: Frontend Development (Week 4)
**Goal:** Build user interfaces

### Task 4.1: Dashboard Page (Days 16-17)
**Deliverables:**
- [ ] Exercise list display
- [ ] Progress tracking
- [ ] Filter by language/difficulty

**Create `app/(dashboard)/page.tsx`:**
```typescript
import { prisma } from '@/lib/prisma'
import ExerciseCard from '@/components/exercise-card'

export default async function DashboardPage() {
  const syllabi = await prisma.syllabus.findMany({
    include: {
      exercises: {
        orderBy: { order: 'asc' },
      },
    },
    orderBy: { order: 'asc' },
  })

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-4xl font-bold mb-8">Coding Platform</h1>
      
      {syllabi.map((syllabus) => (
        <div key={syllabus.id} className="mb-12">
          <h2 className="text-2xl font-semibold mb-4">{syllabus.title}</h2>
          <p className="text-muted-foreground mb-6">{syllabus.description}</p>
          
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {syllabus.exercises.map((exercise) => (
              <ExerciseCard key={exercise.id} exercise={exercise} />
            ))}
          </div>
        </div>
      ))}
    </div>
  )
}
```

**Create `components/exercise-card.tsx`:**
```typescript
import Link from 'next/link'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'

export default function ExerciseCard({ exercise }: { exercise: any }) {
  return (
    <Link href={`/exercise/${exercise.id}`}>
      <Card className="hover:shadow-lg transition-shadow cursor-pointer">
        <CardHeader>
          <div className="flex justify-between items-start">
            <CardTitle className="text-lg">{exercise.title}</CardTitle>
            <Badge variant={
              exercise.difficulty === 'BEGINNER' ? 'default' :
              exercise.difficulty === 'INTERMEDIATE' ? 'secondary' : 'destructive'
            }>
              {exercise.difficulty}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-sm text-muted-foreground line-clamp-2">
            {exercise.description}
          </p>
          <div className="flex items-center gap-2 mt-4">
            <Badge variant="outline">{exercise.language}</Badge>
            <span className="text-sm text-muted-foreground">
              {exercise.points} points
            </span>
          </div>
        </CardContent>
      </Card>
    </Link>
  )
}
```

### Task 4.2: Exercise Workspace with Judge0 IDE (Days 18-19)
**Deliverables:**
- [ ] Judge0 IDE embedded in exercise page
- [ ] Submit button and logic
- [ ] Test results display
- [ ] Hints system

**Choose Your Editor Approach:**

**Option 1: Embedded Judge0 IDE (Recommended for Beginners)**
Uses the official Judge0 IDE with minimal setup.

**Option 2: Custom Monaco Editor (Recommended for Production)**
Full control over UI/UX with Monaco Editor + Judge0 API.

**Option 3: Self-Hosted Judge0 IDE (Full Control)**
Run your own IDE frontend.

---

#### Option 1: Using Embedded Judge0 IDE

**Create `app/exercise/[id]/page.tsx`:**
```typescript
import { prisma } from '@/lib/prisma'
import { notFound } from 'next/navigation'
import ExerciseWorkspaceWithIDE from '@/components/exercise-workspace-ide'

export default async function ExercisePage({
  params,
}: {
  params: { id: string }
}) {
  const exercise = await prisma.exercise.findUnique({
    where: { id: params.id },
    include: { syllabus: true },
  })

  if (!exercise) {
    notFound()
  }

  // Don't send solution to client
  const { solution, ...exerciseData } = exercise

  return <ExerciseWorkspaceWithIDE exercise={exerciseData} />
}
```

**Create `components/exercise-workspace-ide.tsx`:**
```typescript
'use client'

import { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import SubmissionResult from './submission-result'

export default function ExerciseWorkspaceWithIDE({ exercise }: { exercise: any }) {
  const [code, setCode] = useState(exercise.starterCode || '')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [result, setResult] = useState<any>(null)
  const [showHints, setShowHints] = useState<number>(0)
  const [ideOutput, setIdeOutput] = useState<string>('')

  // Listen for messages from Judge0 IDE iframe
  useEffect(() => {
    const handleMessage = (event: MessageEvent) => {
      // Judge0 IDE sends execution results
      if (event.data.type === 'judge0-execution') {
        setIdeOutput(event.data.stdout || event.data.stderr || 'No output')
        setCode(event.data.code) // Update code state
      }
    }

    window.addEventListener('message', handleMessage)
    return () => window.removeEventListener('message', handleMessage)
  }, [])

  const handleSubmit = async () => {
    setIsSubmitting(true)
    setResult(null)

    try {
      const response = await fetch('/api/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          exerciseId: exercise.id,
          code,
          userId: 'temp-user-id', // TODO: Get from auth
        }),
      })

      const data = await response.json()
      setResult(data.submission)
    } catch (error) {
      console.error('Submission error:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const hints = exercise.hints as string[] || []

  // Build Judge0 IDE URL
  const ideUrl = `https://ide.judge0.com/?` + new URLSearchParams({
    source: exercise.starterCode || '',
    language: exercise.language,
    theme: 'dark',
  }).toString()

  return (
    <div className="container mx-auto py-8">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left: Problem Description */}
        <div className="space-y-4">
          <Card className="p-6">
            <div className="flex items-start justify-between mb-4">
              <h1 className="text-2xl font-bold">{exercise.title}</h1>
              <Badge variant={
                exercise.difficulty === 'BEGINNER' ? 'default' :
                exercise.difficulty === 'INTERMEDIATE' ? 'secondary' : 'destructive'
              }>
                {exercise.difficulty}
              </Badge>
            </div>
            
            <div className="prose prose-sm max-w-none mb-4">
              <p className="whitespace-pre-wrap">{exercise.description}</p>
            </div>

            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span>🏆 {exercise.points} points</span>
              <span>⏱️ {exercise.timeLimit}ms limit</span>
              <span>💾 {Math.round(exercise.memoryLimit / 1024)}MB memory</span>
            </div>
          </Card>

          {/* Hints */}
          {hints.length > 0 && (
            <Card className="p-6">
              <h3 className="font-semibold mb-2">💡 Hints</h3>
              <div className="space-y-2">
                {hints.slice(0, showHints).map((hint, idx) => (
                  <div key={idx} className="p-3 bg-muted rounded text-sm">
                    <strong>Hint {idx + 1}:</strong> {hint}
                  </div>
                ))}
                {showHints < hints.length && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowHints(showHints + 1)}
                  >
                    Show Next Hint ({showHints + 1}/{hints.length})
                  </Button>
                )}
              </div>
            </Card>
          )}

          {/* IDE Output Preview */}
          {ideOutput && (
            <Card className="p-4 bg-black">
              <h3 className="text-sm font-semibold text-green-400 mb-2">
                Test Run Output:
              </h3>
              <pre className="text-green-400 text-sm font-mono whitespace-pre-wrap">
                {ideOutput}
              </pre>
            </Card>
          )}
        </div>

        {/* Right: Code Editor */}
        <div className="space-y-4">
          <Card className="p-0 overflow-hidden">
            <Tabs defaultValue="editor">
              <div className="bg-muted p-2">
                <TabsList>
                  <TabsTrigger value="editor">Code Editor</TabsTrigger>
                  <TabsTrigger value="result" disabled={!result}>
                    Submission Result
                  </TabsTrigger>
                </TabsList>
              </div>

              <TabsContent value="editor" className="m-0">
                {/* Judge0 IDE Embed */}
                <iframe
                  src={ideUrl}
                  className="w-full h-[600px] border-0"
                  sandbox="allow-scripts allow-same-origin allow-forms"
                  title="Judge0 Code Editor"
                />
                
                <div className="p-4 space-y-2">
                  <Button
                    onClick={handleSubmit}
                    disabled={isSubmitting}
                    className="w-full"
                    size="lg"
                  >
                    {isSubmitting ? 'Submitting...' : 'Submit Solution'}
                  </Button>
                  <p className="text-xs text-muted-foreground text-center">
                    Use "Run" above to test your code, then "Submit Solution" to grade it
                  </p>
                </div>
              </TabsContent>

              <TabsContent value="result" className="p-4">
                {result && <SubmissionResult result={result} />}
              </TabsContent>
            </Tabs>
          </Card>
        </div>
      </div>
    </div>
  )
}
```

---

#### Option 2: Custom Monaco Editor (Production Ready)

**Create `components/exercise-workspace-custom.tsx`:**
#### Option 2: Custom Monaco Editor (Production Ready)

**Install Monaco Editor:**
```bash
npm install @monaco-editor/react
```

**Create `components/exercise-workspace-custom.tsx`:**
```typescript
'use client'

import { useState } from 'react'
import Editor from '@monaco-editor/react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { PlayCircle, Send } from 'lucide-react'
import SubmissionResult from './submission-result'

export default function ExerciseWorkspaceCustom({ exercise }: { exercise: any }) {
  const [code, setCode] = useState(exercise.starterCode || '')
  const [isRunning, setIsRunning] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [testOutput, setTestOutput] = useState<any>(null)
  const [result, setResult] = useState<any>(null)
  const [showHints, setShowHints] = useState<number>(0)
  const [activeTab, setActiveTab] = useState('editor')

  // Run code for testing (not graded)
  const handleRun = async () => {
    setIsRunning(true)
    setTestOutput(null)

    try {
      const response = await fetch('/api/run', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          code,
          language: exercise.language,
          stdin: '', // Could add input field
        }),
      })

      const data = await response.json()
      setTestOutput(data)
    } catch (error) {
      setTestOutput({ error: 'Failed to run code' })
    } finally {
      setIsRunning(false)
    }
  }

  // Submit for grading
  const handleSubmit = async () => {
    setIsSubmitting(true)
    setResult(null)

    try {
      const response = await fetch('/api/submit', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          exerciseId: exercise.id,
          code,
          userId: 'temp-user-id',
        }),
      })

      const data = await response.json()
      setResult(data.submission)
      setActiveTab('result')
    } catch (error) {
      console.error('Submission error:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const hints = exercise.hints as string[] || []

  // Monaco language mapping
  const getMonacoLanguage = (lang: string) => {
    const map: Record<string, string> = {
      python: 'python',
      javascript: 'javascript',
      typescript: 'typescript',
      java: 'java',
      cpp: 'cpp',
      c: 'c',
      csharp: 'csharp',
      go: 'go',
      rust: 'rust',
      ruby: 'ruby',
      php: 'php',
    }
    return map[lang.toLowerCase()] || 'plaintext'
  }

  return (
    <div className="container mx-auto py-8">
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Left: Problem Description */}
        <div className="space-y-4">
          <Card className="p-6">
            <div className="flex items-start justify-between mb-4">
              <h1 className="text-2xl font-bold">{exercise.title}</h1>
              <Badge variant={
                exercise.difficulty === 'BEGINNER' ? 'default' :
                exercise.difficulty === 'INTERMEDIATE' ? 'secondary' : 'destructive'
              }>
                {exercise.difficulty}
              </Badge>
            </div>
            
            <div className="prose prose-sm max-w-none mb-4">
              <p className="whitespace-pre-wrap">{exercise.description}</p>
            </div>

            <div className="flex items-center gap-4 text-sm text-muted-foreground">
              <span>🏆 {exercise.points} points</span>
              <span>⏱️ {exercise.timeLimit}ms limit</span>
              <span>💾 {Math.round(exercise.memoryLimit / 1024)}MB memory</span>
            </div>
          </Card>

          {/* Hints */}
          {hints.length > 0 && (
            <Card className="p-6">
              <h3 className="font-semibold mb-2">💡 Hints</h3>
              <div className="space-y-2">
                {hints.slice(0, showHints).map((hint, idx) => (
                  <div key={idx} className="p-3 bg-muted rounded text-sm">
                    <strong>Hint {idx + 1}:</strong> {hint}
                  </div>
                ))}
                {showHints < hints.length && (
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => setShowHints(showHints + 1)}
                  >
                    Show Next Hint ({showHints + 1}/{hints.length})
                  </Button>
                )}
              </div>
            </Card>
          )}

          {/* Test Output */}
          {testOutput && (
            <Card className="p-4 bg-black">
              <h3 className="text-sm font-semibold text-green-400 mb-2">
                Test Run Output:
              </h3>
              <pre className="text-green-400 text-sm font-mono whitespace-pre-wrap">
                {testOutput.stdout || testOutput.stderr || testOutput.error || 'No output'}
              </pre>
              {testOutput.time && (
                <p className="text-xs text-muted-foreground mt-2">
                  Execution time: {testOutput.time}ms
                </p>
              )}
            </Card>
          )}
        </div>

        {/* Right: Code Editor */}
        <div className="space-y-4">
          <Card className="overflow-hidden">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
              <div className="bg-muted p-2">
                <TabsList>
                  <TabsTrigger value="editor">Code Editor</TabsTrigger>
                  <TabsTrigger value="result" disabled={!result}>
                    Submission Result
                  </TabsTrigger>
                </TabsList>
              </div>

              <TabsContent value="editor" className="m-0">
                <Editor
                  height="500px"
                  language={getMonacoLanguage(exercise.language)}
                  value={code}
                  onChange={(value) => setCode(value || '')}
                  theme="vs-dark"
                  options={{
                    minimap: { enabled: false },
                    fontSize: 14,
                    lineNumbers: 'on',
                    scrollBeyondLastLine: false,
                    automaticLayout: true,
                    tabSize: 2,
                    wordWrap: 'on',
                  }}
                />
                
                <div className="p-4 space-y-2 bg-background">
                  <div className="flex gap-2">
                    <Button
                      onClick={handleRun}
                      disabled={isRunning}
                      variant="outline"
                      className="flex-1"
                    >
                      <PlayCircle className="w-4 h-4 mr-2" />
                      {isRunning ? 'Running...' : 'Test Run'}
                    </Button>
                    <Button
                      onClick={handleSubmit}
                      disabled={isSubmitting}
                      className="flex-1"
                    >
                      <Send className="w-4 h-4 mr-2" />
                      {isSubmitting ? 'Submitting...' : 'Submit Solution'}
                    </Button>
                  </div>
                  <p className="text-xs text-muted-foreground text-center">
                    Test your code first, then submit for grading
                  </p>
                </div>
              </TabsContent>

              <TabsContent value="result" className="p-4">
                {result && <SubmissionResult result={result} />}
              </TabsContent>
            </Tabs>
          </Card>
        </div>
      </div>
    </div>
  )
}
```

**Create `/api/run` endpoint for test runs:**
```typescript
// app/api/run/route.ts
import { NextRequest, NextResponse } from 'next/server'
import { executeCode } from '@/lib/judge0'

export async function POST(request: NextRequest) {
  try {
    const { code, language, stdin } = await request.json()

    const result = await executeCode({
      code,
      language,
      stdin: stdin || '',
    })

    return NextResponse.json({
      stdout: result.stdout,
      stderr: result.stderr,
      compile_output: result.compile_output,
      status: result.status,
      time: result.time,
      memory: result.memory,
    })
  } catch (error) {
    console.error('Run error:', error)
    return NextResponse.json(
      { error: 'Failed to run code' },
      { status: 500 }
    )
  }
}
```

**Choose which workspace to use in your exercise page:**
```typescript
// app/exercise/[id]/page.tsx
import ExerciseWorkspaceCustom from '@/components/exercise-workspace-custom'
// OR
import ExerciseWorkspaceWithIDE from '@/components/exercise-workspace-ide'

export default async function ExercisePage({ params }: { params: { id: string } }) {
  // ... fetch exercise ...
  
  // Use custom Monaco editor (recommended)
  return <ExerciseWorkspaceCustom exercise={exerciseData} />
  
  // OR use Judge0 IDE embed
  // return <ExerciseWorkspaceWithIDE exercise={exerciseData} />
}
```

---

#### Comparison of Editor Options

| Feature | Judge0 IDE Embed | Custom Monaco | Self-Hosted IDE |
|---------|------------------|---------------|-----------------|
| Setup Difficulty | ⭐ Easy | ⭐⭐ Medium | ⭐⭐⭐ Complex |
| Customization | ❌ Limited | ✅ Full Control | ✅ Full Control |
| Offline Support | ❌ No | ✅ Yes | ✅ Yes |
| Maintenance | ✅ Zero | ⭐ Minimal | ⭐⭐ Moderate |
| Best For | Quick MVP | Production | Enterprise |

**Recommendation:** Start with **Custom Monaco Editor** (Option 2) for the best balance of control and simplicity.

**Create `components/submission-result.tsx`:**
```typescript
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { CheckCircle2, XCircle } from 'lucide-react'

export default function SubmissionResult({ result }: { result: any }) {
  return (
    <div className="space-y-4">
      <Card className="p-4">
        <div className="flex items-center justify-between mb-4">
          <h3 className="font-semibold">
            {result.allPassed ? '✅ All Tests Passed!' : '❌ Some Tests Failed'}
          </h3>
          <Badge variant={result.allPassed ? 'default' : 'destructive'}>
            Score: {result.score}
          </Badge>
        </div>

        <div className="space-y-3">
          {result.testResults.map((test: any, idx: number) => (
            <Card key={idx} className="p-3">
              <div className="flex items-start gap-2">
                {test.passed ? (
                  <CheckCircle2 className="w-5 h-5 text-green-500 flex-shrink-0 mt-0.5" />
                ) : (
                  <XCircle className="w-5 h-5 text-red-500 flex-shrink-0 mt-0.5" />
                )}
                <div className="flex-1 space-y-2 text-sm">
                  <div>
                    <span className="font-medium">Test {idx + 1}:</span>
                    {test.input && (
                      <div className="text-muted-foreground">
                        Input: <code className="bg-muted px-1">{test.input}</code>
                      </div>
                    )}
                  </div>
                  {!test.passed && (
                    <>
                      <div>
                        Expected: <code className="bg-muted px-1">{test.expectedOutput}</code>
                      </div>
                      <div>
                        Got: <code className="bg-muted px-1">{test.actualOutput}</code>
                      </div>
                      {test.error && (
                        <div className="text-red-600">
                          Error: {test.error}
                        </div>
                      )}
                    </>
                  )}
                  {test.executionTime && (
                    <div className="text-xs text-muted-foreground">
                      Time: {test.executionTime}ms
                    </div>
                  )}
                </div>
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  )
}
```

### Task 4.3: Admin Panel (Days 19-20)
**Deliverables:**
- [ ] Exercise generation UI
- [ ] Syllabus management
- [ ] Exercise editing

**Create `app/admin/page.tsx`:**
```typescript
'use client'

import { useState } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Card } from '@/components/ui/card'
import { Textarea } from '@/components/ui/textarea'

export default function AdminPage() {
  const [topic, setTopic] = useState('')
  const [language, setLanguage] = useState('python')
  const [difficulty, setDifficulty] = useState('BEGINNER')
  const [syllabusId, setSyllabusId] = useState('')
  const [concepts, setConcepts] = useState('')
  const [isGenerating, setIsGenerating] = useState(false)
  const [result, setResult] = useState<any>(null)

  const handleGenerate = async () => {
    setIsGenerating(true)
    setResult(null)

    try {
      const response = await fetch('/api/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          topic,
          language,
          difficulty,
          syllabusId,
          concepts: concepts.split(',').map(c => c.trim()).filter(Boolean),
        }),
      })

      const data = await response.json()
      setResult(data.exercise)
      
      // Reset form
      setTopic('')
      setConcepts('')
    } catch (error) {
      console.error('Generation error:', error)
    } finally {
      setIsGenerating(false)
    }
  }

  return (
    <div className="container mx-auto py-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-8">Admin Panel</h1>

      <Card className="p-6 space-y-4">
        <h2 className="text-xl font-semibold">Generate Exercise with AI</h2>

        <div className="space-y-2">
          <Label htmlFor="topic">Topic / Exercise Name</Label>
          <Input
            id="topic"
            value={topic}
            onChange={(e) => setTopic(e.target.value)}
            placeholder="e.g., FizzBuzz, Two Sum, Fibonacci"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="language">Programming Language</Label>
          <Select value={language} onValueChange={setLanguage}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="python">Python</SelectItem>
              <SelectItem value="javascript">JavaScript</SelectItem>
              <SelectItem value="java">Java</SelectItem>
              <SelectItem value="cpp">C++</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="difficulty">Difficulty</Label>
          <Select value={difficulty} onValueChange={setDifficulty}>
            <SelectTrigger>
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="BEGINNER">Beginner</SelectItem>
              <SelectItem value="INTERMEDIATE">Intermediate</SelectItem>
              <SelectItem value="ADVANCED">Advanced</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2">
          <Label htmlFor="syllabusId">Syllabus ID</Label>
          <Input
            id="syllabusId"
            value={syllabusId}
            onChange={(e) => setSyllabusId(e.target.value)}
            placeholder="Paste syllabus ID from database"
          />
        </div>

        <div className="space-y-2">
          <Label htmlFor="concepts">Concepts (comma-separated)</Label>
          <Textarea
            id="concepts"
            value={concepts}
            onChange={(e) => setConcepts(e.target.value)}
            placeholder="e.g., loops, conditionals, arrays"
          />
        </div>

        <Button
          onClick={handleGenerate}
          disabled={isGenerating || !topic || !syllabusId}
          className="w-full"
          size="lg"
        >
          {isGenerating ? 'Generating...' : 'Generate Exercise'}
        </Button>

        {result && (
          <Card className="p-4 mt-4 bg-green-50">
            <h3 className="font-semibold text-green-800 mb-2">✓ Exercise Created!</h3>
            <p className="text-sm text-green-700">
              {result.title} - {result.points} points
            </p>
          </Card>
        )}
      </Card>
    </div>
  )
}
```

### Task 4.4: Enhanced Code Editor (Day 20)
**Deliverables:**
- [ ] Syntax highlighting
- [ ] Line numbers
- [ ] Better UX

**Optional: Integrate Monaco Editor**
```bash
npm install @monaco-editor/react
```

**Update `components/exercise-workspace.tsx`:**
```typescript
import Editor from '@monaco-editor/react'

// Replace Textarea with:
<Editor
  height="400px"
  defaultLanguage={exercise.language}
  value={code}
  onChange={(value) => setCode(value || '')}
  theme="vs-dark"
  options={{
    minimap: { enabled: false },
    fontSize: 14,
    lineNumbers: 'on',
    scrollBeyondLastLine: false,
  }}
/>
```

---

## Phase 5: Polish & Deployment (Week 5+)
**Goal:** Production readiness

### Task 5.1: Authentication (Days 21-23)
**Options:**
- NextAuth.js with credentials
- Clerk
- Supabase Auth

### Task 5.2: Testing (Days 23-24)
- Unit tests with Jest
- Integration tests with Playwright
- API route testing

### Task 5.3: Production Deployment with Self-Hosted Judge0 (Days 24-25)

**Architecture Overview:**
```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   Vercel    │────▶│  Your VPS    │────▶│  PostgreSQL DB  │
│  (Frontend) │     │  (Judge0)    │     │  (Neon/Railway) │
└─────────────┘     └──────────────┘     └─────────────────┘
```

#### Deployment Option 1: VPS (Best for Control)

**Recommended Providers:**
- DigitalOcean ($12/month - 2GB RAM)
- Hetzner ($5/month - 2GB RAM, best value)
- Linode ($12/month - 2GB RAM)

**Step 1: Set Up VPS**

```bash
# SSH into your server
ssh root@your-server-ip

# Update system
apt update && apt upgrade -y

# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Install Docker Compose
apt install docker-compose -y

# Create directory for Judge0
mkdir -p /opt/judge0
cd /opt/judge0
```

**Step 2: Deploy Judge0 on VPS**

**Create `docker-compose.prod.yml`:**
```yaml
version: '3.7'

services:
  db:
    image: postgres:13
    env_file: .env.prod
    volumes:
      - postgres-data:/var/lib/postgresql/data
    restart: always

  redis:
    image: redis:6
    command: redis-server --appendonly yes
    volumes:
      - redis-data:/data
    restart: always

  server:
    image: judge0/judge0:1.13.0
    env_file: .env.prod
    volumes:
      - ./judge0.conf:/judge0.conf:ro
    ports:
      - "2358:2358"
    privileged: true
    restart: always
    depends_on:
      - db
      - redis

  workers:
    image: judge0/judge0:1.13.0
    command: ["./scripts/workers"]
    env_file: .env.prod
    volumes:
      - ./judge0.conf:/judge0.conf:ro
    privileged: true
    restart: always
    depends_on:
      - db
      - redis

volumes:
  postgres-data:
  redis-data:
```

**Create `.env.prod`:**
```bash
# PostgreSQL
POSTGRES_HOST=db
POSTGRES_PORT=5432
POSTGRES_DB=judge0
POSTGRES_USER=judge0
POSTGRES_PASSWORD=CHANGE_THIS_TO_SECURE_PASSWORD

# Redis
REDIS_HOST=redis
REDIS_PORT=6379

# Judge0 Configuration
WORKERS_COUNT=2
MAX_QUEUE_SIZE=100
ENABLE_WAIT_RESULT=true
ENABLE_COMPILER_OPTIONS=true

# Security
JUDGE0_HOMEPAGE=https://your-domain.com
AUTHENTICATION_TOKEN=CHANGE_THIS_TO_SECURE_TOKEN
```

**Create `judge0.conf`:**
```ruby
Rails.application.config.judge do
  # Maximum execution time for submissions (in seconds)
  max_cpu_time_limit = 15
  max_wall_time_limit = 30
  
  # Memory limits
  max_memory_limit = 512000  # 512 MB
  
  # Allow multiple submissions
  enable_batched_submissions = true
  
  # Submission queue limits
  max_queue_size = 100
end
```

**Step 3: Start Judge0**
```bash
# Generate secure passwords
openssl rand -base64 32  # For POSTGRES_PASSWORD
openssl rand -base64 32  # For AUTHENTICATION_TOKEN

# Update .env.prod with these values

# Start services
docker-compose -f docker-compose.prod.yml up -d

# Check logs
docker-compose -f docker-compose.prod.yml logs -f

# Test the API
curl http://localhost:2358/about
```

**Step 4: Set Up Nginx Reverse Proxy**

```bash
apt install nginx certbot python3-certbot-nginx -y
```

**Create `/etc/nginx/sites-available/judge0`:**
```nginx
server {
    listen 80;
    server_name judge0.your-domain.com;

    location / {
        proxy_pass http://localhost:2358;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        
        # Timeouts for long-running code execution
        proxy_connect_timeout 300s;
        proxy_send_timeout 300s;
        proxy_read_timeout 300s;
    }
}
```

**Enable site and get SSL:**
```bash
ln -s /etc/nginx/sites-available/judge0 /etc/nginx/sites-enabled/
nginx -t
systemctl reload nginx

# Get SSL certificate
certbot --nginx -d judge0.your-domain.com
```

**Step 5: Configure Firewall**
```bash
ufw allow 22/tcp   # SSH
ufw allow 80/tcp   # HTTP
ufw allow 443/tcp  # HTTPS
ufw enable
```

---

#### Deployment Option 2: Railway (Easier Setup)

Railway can host Judge0 with a PostgreSQL addon.

**Step 1: Create Railway Project**
```bash
# Install Railway CLI
npm i -g @railway/cli

# Login
railway login

# Create new project
railway init
```

**Step 2: Add Judge0 Service**

Create `railway.json`:
```json
{
  "$schema": "https://railway.app/railway.schema.json",
  "build": {
    "builder": "DOCKERFILE",
    "dockerfilePath": "./Dockerfile"
  },
  "deploy": {
    "startCommand": "/judge0/scripts/server",
    "restartPolicyType": "ON_FAILURE"
  }
}
```

**Create `Dockerfile`:**
```dockerfile
FROM judge0/judge0:1.13.0

# Expose port
EXPOSE 2358

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:2358/about || exit 1

# Start server
CMD ["/judge0/scripts/server"]
```

**Deploy:**
```bash
railway up
railway domain  # Get your public URL
```

---

#### Deployment Option 3: Fly.io (Cost-Effective)

**Step 1: Install Fly CLI**
```bash
curl -L https://fly.io/install.sh | sh
fly auth login
```

**Step 2: Create fly.toml**
```toml
app = "your-judge0-app"
primary_region = "lhr"  # London

[build]
  image = "judge0/judge0:1.13.0"

[env]
  REDIS_HOST = "redis.internal"
  POSTGRES_HOST = "postgres.internal"
  WORKERS_COUNT = "2"

[http_service]
  internal_port = 2358
  force_https = true
  auto_stop_machines = false
  auto_start_machines = true
  min_machines_running = 1

[[services]]
  protocol = "tcp"
  internal_port = 2358

  [[services.ports]]
    port = 80
    handlers = ["http"]

  [[services.ports]]
    port = 443
    handlers = ["tls", "http"]

[[vm]]
  cpu_kind = "shared"
  cpus = 1
  memory_mb = 1024
```

**Deploy:**
```bash
fly launch  # Follow prompts
fly postgres create  # Create database
fly redis create     # Create Redis
fly deploy
```

---

### Next.js Frontend Deployment (Vercel)

**Step 1: Prepare Environment Variables**

Create `.env.production`:
```env
# Your deployed Judge0 URL
JUDGE0_BASE_URL=https://judge0.your-domain.com

# Database (Neon, Railway, or Supabase)
DATABASE_URL=postgresql://user:pass@host:5432/db

# AI Service
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...

# Auth (NextAuth)
NEXTAUTH_SECRET=generated-secret
NEXTAUTH_URL=https://your-app.vercel.app

# Judge0 Authentication (if you set one)
JUDGE0_AUTH_TOKEN=your-secure-token
```

**Step 2: Deploy to Vercel**
```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel

# Set environment variables
vercel env add JUDGE0_BASE_URL production
vercel env add DATABASE_URL production
# ... repeat for all vars

# Deploy to production
vercel --prod
```

**Step 3: Update Next.js Config**

Add to `next.config.js`:
```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  async headers() {
    return [
      {
        source: '/api/:path*',
        headers: [
          { key: 'Access-Control-Allow-Origin', value: '*' },
          { key: 'Access-Control-Allow-Methods', value: 'GET,POST,OPTIONS' },
          { key: 'Access-Control-Allow-Headers', value: 'Content-Type' },
        ],
      },
    ]
  },
  // If using self-hosted Judge0 IDE
  async rewrites() {
    return [
      {
        source: '/judge0/:path*',
        destination: 'https://judge0.your-domain.com/:path*',
      },
    ]
  },
}

module.exports = nextConfig
```

---

### Database Deployment Options

**Option 1: Neon (Recommended)**
- Free tier: 0.5GB storage
- Serverless PostgreSQL
- Auto-scaling
- https://neon.tech

**Option 2: Railway**
- $5/month for 1GB
- Built-in backups
- Easy setup
- https://railway.app

**Option 3: Supabase**
- Free tier: 500MB
- Includes auth system
- Real-time features
- https://supabase.com

**Migration to Production DB:**
```bash
# Export current schema
npx prisma migrate dev --create-only

# Update DATABASE_URL in .env.production
# Deploy migrations
npx prisma migrate deploy

# Seed production data
npx prisma db seed
```

---

### Production Checklist

**Pre-Deployment:**
- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] Judge0 tested and accessible
- [ ] SSL certificates installed
- [ ] Firewall configured
- [ ] Backups enabled

**Security:**
- [ ] Judge0 authentication token set
- [ ] CORS configured properly
- [ ] Rate limiting implemented
- [ ] User authentication working
- [ ] API routes protected
- [ ] Input validation on all endpoints

**Monitoring:**
- [ ] Error tracking (Sentry)
- [ ] Uptime monitoring (UptimeRobot)
- [ ] Log aggregation (Better Stack)
- [ ] Performance monitoring (Vercel Analytics)

**Testing:**
- [ ] Test code execution end-to-end
- [ ] Verify all languages work
- [ ] Test submission grading
- [ ] Check AI generation
- [ ] Load test with multiple submissions

---

### Cost Estimates

**Monthly Costs (Minimal Setup):**
- Frontend (Vercel): $0 (Hobby tier)
- Database (Neon): $0 (Free tier)
- Judge0 VPS (Hetzner): $5
- Domain: $10/year (~$1/month)
- SSL: $0 (Let's Encrypt)
**Total: ~$6/month**

**Monthly Costs (Production):**
- Frontend (Vercel Pro): $20
- Database (Railway): $10
- Judge0 VPS (DigitalOcean): $12
- CDN (Cloudflare): $0
- Monitoring (Sentry): $0 (Developer tier)
**Total: ~$42/month**

### Task 5.4: Additional Features
- [ ] User profiles
- [ ] Leaderboards
- [ ] Social features (comments, likes)
- [ ] Progress analytics
- [ ] Email notifications
- [ ] Certificate generation

---

## Critical Success Milestones

### Milestone 1: "Hello World" ✅
**Definition:** Click button → sends code → gets result back
**Test:** Hardcoded Python `print("Hello")` executes successfully

### Milestone 2: "Full Loop" ✅
**Definition:** Load exercise → edit code → submit → see grading
**Test:** Complete one exercise end-to-end

### Milestone 3: "AI Generation" ✅
**Definition:** Generate exercise via AI → save to DB → appears in UI
**Test:** Create 3 exercises with different difficulties

### Milestone 4: "Production Ready" ✅
**Definition:** Deployed, authenticated, tested
**Test:** External user can register and complete exercises

---

## Troubleshooting Guide

### Common Issues & Solutions

**Database Connection Failed**
```bash
# Check if PostgreSQL is running
docker ps
# Reset database
npx prisma migrate reset
```

**Judge0 Not Responding**
```bash
# Check Judge0 logs
docker logs judge0-server
# Test with curl
curl -X POST http://localhost:2358/submissions
```

**Build Errors**
```bash
# Clear cache
rm -rf .next
npm run build
```

**Prisma Issues**
```bash
# Regenerate client
npx prisma generate
# View database
npx prisma studio
```

---

## Environment Variables Checklist

```env
# Database
DATABASE_URL="postgresql://..."

# Judge0 (choose one)
JUDGE0_BASE_URL="http://localhost:2358"
JUDGE0_API_KEY="your_rapidapi_key"

# AI (choose one)
OPENAI_API_KEY="sk-..."
ANTHROPIC_API_KEY="sk-ant-..."

# Auth (if using NextAuth)
NEXTAUTH_SECRET="generate_with_openssl"
NEXTAUTH_URL="http://localhost:3000"

# App
NODE_ENV="development"
```

---

## Daily Checklist Template

**Before Starting Work:**
- [ ] Pull latest code
- [ ] Start database (`docker start code-platform-db`)
- [ ] Start dev server (`npm run dev`)
- [ ] Open Prisma Studio (`npx prisma studio`)

**After Completing Work:**
- [ ] Commit changes
- [ ] Update documentation
- [ ] Test in browser
- [ ] Push to Git

---

## Resources & Documentation

**Official Docs:**
- Next.js: https://nextjs.org/docs
- Prisma: https://www.prisma.io/docs
- Judge0: https://ce.judge0.com
- Tailwind CSS: https://tailwindcss.com/docs
- ShadCN UI: https://ui.shadcn.com

**Community:**
- Judge0 GitHub Issues
- Next.js Discord
- Stack Overflow

---

## Next Steps After Completion

1. **Performance Optimization**
   - Implement caching
   - Code splitting
   - Database indexing

2. **Advanced Features**
   - Real-time collaboration
   - Video explanations
   - Peer code reviews

3. **Monetization**
   - Premium content
   - Certification programs
   - Enterprise features

---

**Remember:** Start small, test often, and iterate. Your first goal is getting one exercise working end-to-end, not building the entire platform at once.