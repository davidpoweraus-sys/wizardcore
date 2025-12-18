WizardCore: Complete Mock Data Audit & Go Backend Implementation Plan
PART 1: FILES CONTAINING MOCK DATA/PLACEHOLDERS
Frontend Files Requiring Backend Integration
1. app/dashboard/page.tsx
Mock Data:
- Lines 10-13: Hardcoded welcome message "Welcome back, Coder!"
- Lines 23-50: Static pathway cards with hardcoded progress (75%, 30%, 15%, 50%)
- Lines 67-79: Fake upcoming deadlines section
Required Changes:
- Fetch user profile from API: GET /api/v1/users/me
- Fetch enrolled pathways with progress: GET /api/v1/users/me/pathways
- Fetch upcoming deadlines: GET /api/v1/users/me/deadlines
- Replace hardcoded data with API responses
- Add loading states and error handling
---
2. components/dashboard/QuickStats.tsx
Mock Data:
- Lines 3-32: All statistics hardcoded (Active Courses: 4, Completion Rate: 78%, Study Time: 42h, XP: 2,850)
Required Changes:
- Fetch stats from: GET /api/v1/users/me/stats
- API should return:
    {
    active_courses: 4,
    active_courses_change: +1 this week,
    completion_rate: 78,
    completion_rate_change: +5% from last month,
    study_time_hours: 42,
    study_time_week: 12,
    xp_total: 2850,
    xp_today: 320
  }
  - Convert to client component with SWR/React Query
- Add real-time updates
---
3. components/dashboard/ProgressChart.tsx
Mock Data:
- Lines 6-14: Hardcoded weekly progress data (Mon-Sun percentages)
- Lines 48, 52, 56: Fake metrics (Avg Daily Time: 2h 15m, Completion Rate: 78%, Streak: 14 days)
- Line 27: Fake trend "+12%"
Required Changes:
- Fetch weekly activity: GET /api/v1/users/me/activity/weekly
- API response:
    {
    weekly_data: [
      {day: Mon, value: 40, hours: 2.5},
      {day: Tue, value: 60, hours: 3.0},
      ...
    ],
    avg_daily_time_minutes: 135,
    completion_rate: 78,
    current_streak: 14,
    trend_percentage: 12
  }
  - Replace static data array with API data
- Add date range selector
---
4. components/dashboard/RecentActivity.tsx
Mock Data:
- Lines 3-49: Entire activity feed hardcoded (5 fake activities)
Required Changes:
- Fetch activities: GET /api/v1/users/me/activities?limit=5
- API response:
    {
    activities: [
      {
        id: uuid,
        type: completion|practice|reading|achievement|streak,
        title: Completed Python Module 3,
        description: Functions and Lambda Expressions,
        timestamp: 2025-12-17T14:30:00Z,
        icon: CheckCircle,
        color: text-green-400
      }
    ]
  }
  - Map activity types to icons dynamically
- Add "View All" pagination
---
5. components/dashboard/PathwayCard.tsx
Mock Data:
- Component receives props but buttons are non-functional (lines 50-55)
Required Changes:
- Add onClick handlers:
  - "Continue" → Navigate to: /dashboard/learning?pathway={id}&resume=true
  - "Details" → Navigate to: /dashboard/pathways/{id}
- Pass pathway ID via props
- Add navigation logic with Next.js router
- Track button clicks: POST /api/v1/analytics/events
---
6. components/dashboard/Header.tsx
Mock Data:
- Lines 8-17: Fetches user from Supabase but uses basic email display
- Lines 24-29: Search input with no functionality
- Lines 34-37: Notification bell with fake badge
Required Changes:
- Fetch full user profile: GET /api/v1/users/me
- Implement search: GET /api/v1/search?q={query}&type=course,module,help
- Fetch notifications: GET /api/v1/users/me/notifications?unread=true
- Add search results dropdown
- Add notification panel with real data
- Mark notifications read: PATCH /api/v1/notifications/{id}/read
---
7. components/dashboard/Sidebar.tsx
Mock Data:
- Lines 20-30: Navigation items (functional, but could show counts)
- Lines 37-40: Sign out (functional via Supabase)
Required Changes:
- Add badge counts to nav items from: GET /api/v1/users/me/nav-counts
    {
    achievements: 2,
    notifications: 5,
    practice_challenges: 3
  }
  - Display unread counts as badges
- Invalidate session on backend: POST /api/v1/auth/logout
---
8. app/dashboard/achievements/page.tsx + components/achievements/AchievementsDisplay.tsx
Mock Data:
- Lines 5-86 (AchievementsDisplay): All 8 badges hardcoded
- Lines 88-97: Fake leaderboard with 8 users
- Lines 101, 126, 141: Fake stats (badges earned, global rank, total XP)
Required Changes:
- Fetch badges: GET /api/v1/users/me/achievements
    {
    badges: [
      {
        id: uuid,
        title: First Steps,
        description: Complete your first exercise,
        icon: Star,
        color: from-yellow-400 to-orange-400,
        earned: true,
        earned_date: 2025-11-10,
        rarity: Common,
        progress: 100
      }
    ],
    stats: {
      earned_count: 4,
      total_count: 8,
      new_this_month: 2,
      global_rank: 7,
      total_xp: 2850,
      next_rank_xp: 3000
    }
  }
  - Fetch leaderboard: GET /api/v1/leaderboard?limit=10&include_self=true
- Replace all static arrays
- Add achievement detail modal
- Implement badge sharing feature
---
9. app/dashboard/leaderboard/page.tsx + components/leaderboard/LeaderboardTable.tsx
Mock Data:
- Lines 5-16 (LeaderboardTable): 10 hardcoded users with full stats
- Lines 18-24: Fake timeframe filters (non-functional)
- Lines 35, 47, 63, 76: Fake aggregate stats
Required Changes:
- Fetch leaderboard: GET /api/v1/leaderboard?timeframe={all|month|week}&pathway={id}&limit=10&offset=0
    {
    leaderboard: [
      {
        rank: 1,
        user_id: uuid,
        username: Alex Chen,
        avatar_url: https://...,
        xp: 12540,
        streak: 21,
        badges: 12,
        country_code: US,
        trend: up|down|same,
        is_current_user: false
      }
    ],
    stats: {
      total_learners: 2847,
      current_user_rank: 7,
      current_user_change: 2,
      top_xp: 12540,
      top_username: Alex Chen,
      country_count: 64
    },
    pagination: {
      total: 2847,
      page: 1,
      per_page: 10
    }
  }
  - Implement timeframe filters with API calls
- Add pagination with page state
- Add "View Profile" modal/page
- Implement country flags with country code lookup
---
10. app/dashboard/learning/page.tsx + components/learning/LearningEnvironment.tsx
Mock Data:
- Lines 26-61 (LearningEnvironment): Hardcoded lesson content
- Lines 63-89: Hardcoded exercise with test cases
- Lines 136-144: Fake timer, solver count, difficulty, XP
- Lines 113-116: Submit solution just logs to console
Required Changes:
- Fetch lesson/exercise: GET /api/v1/exercises/{id}
    {
    exercise: {
      id: uuid,
      title: Python Functions: Sum Calculation,
      module_id: uuid,
      module_title: Module 1: The Hacker's Toolkit,
      difficulty: BEGINNER,
      points: 100,
      time_limit_minutes: 15,
      concurrent_solvers: 1245,
      objectives: [...],
      content: markdown content,
      examples: [...],
      description: ...,
      constraints: [...],
      test_cases: [
        {
          id: uuid,
          input: [1, 2, 3],
          expected_output: 6,
          is_hidden: false,
          points: 20
        }
      ],
      hints: [...],
      starter_code: def calculate_sum(numbers):\n    ...
    }
  }
  - Get user's saved code: GET /api/v1/submissions/{exercise_id}/latest
- Submit solution: POST /api/v1/submissions
    {
    exercise_id: uuid,
    source_code: ...,
    language_id: 71
  }
  - Backend should:
  1. Store submission in database
  2. Execute via Judge0 (call judge0 from Go backend)
  3. Validate against test cases
  4. Award XP if passed
  5. Update user progress
- Poll for result: GET /api/v1/submissions/{submission_id}
- Fetch real-time solver count: WebSocket or polling GET /api/v1/exercises/{id}/stats
- Implement auto-save: POST /api/v1/submissions/{exercise_id}/save-draft
---
11. app/dashboard/pathways/page.tsx + components/pathways/PathwaySelection.tsx
Mock Data:
- Lines 6-97 (PathwaySelection): All 6 pathways hardcoded with full details
- Lines 166, 189: "Enroll Now" and "Take Assessment" buttons non-functional
Required Changes:
- Fetch pathways: GET /api/v1/pathways
    {
    pathways: [
      {
        id: uuid,
        title: Python for Offensive Security,
        subtitle: The Hacker's Swiss Army Knife,
        description: ...,
        level: Beginner,
        duration_weeks: 8,
        student_count: 1200,
        rating: 4.8,
        module_count: 5,
        color_gradient: from-green-400 to-cyan-400,
        icon: 🐍,
        is_locked: false,
        progress: 0,
        is_enrolled: false,
        prerequisites: []
      }
    ]
  }
  - Enroll in pathway: POST /api/v1/pathways/{id}/enroll
- Start assessment: GET /api/v1/assessment/start → Navigate to assessment page
- Check prerequisites before enrollment
- Show enrollment confirmation modal
---
12. app/dashboard/practice/page.tsx + components/practice/PracticeArena.tsx
Mock Data:
- Lines 6-43 (PracticeArena): 4 challenge types hardcoded
- Lines 45-52: 6 practice areas with fake completion stats
- Lines 54-59: 4 fake recent matches
- Lines 78, 94, 109, 124: All practice stats fake
Required Changes:
- Fetch practice stats: GET /api/v1/users/me/practice/stats
    {
    practice_score: 1850,
    rank_percentile: 15,
    duels_won: 24,
    duels_total: 35,
    win_rate: 68,
    avg_time_seconds: 134,
    live_opponents: 12
  }
  - Fetch challenge types: GET /api/v1/practice/challenges
- Fetch practice areas: GET /api/v1/practice/areas
    {
    areas: [
      {
        name: Python,
        exercise_count: 42,
        completed_count: 28,
        color_gradient: from-green-400 to-cyan-400
      }
    ]
  }
  - Fetch recent matches: GET /api/v1/users/me/matches?limit=5
- Start challenge: POST /api/v1/practice/challenges/{type}/start → Returns exercise ID
- Start duel: POST /api/v1/practice/duels/matchmake → WebSocket connection
- Implement real-time duel system with WebSocket
---
13. app/dashboard/profile/page.tsx + components/profile/ProfileManager.tsx
Mock Data:
- Lines 7-15 (ProfileManager): Hardcoded profile data (Alex Chen)
- Lines 17-24: Preferences in local state only
- Lines 34-38: Save function just logs to console
- Lines 147-178: Security buttons non-functional
Required Changes:
- Fetch profile: GET /api/v1/users/me/profile
    {
    id: uuid,
    name: Alex Chen,
    email: alex@example.com,
    avatar_url: https://...,
    bio: ...,
    location: San Francisco, CA,
    website: https://alexchen.dev,
    github_username: alexchen,
    twitter_username: @alexchen,
    created_at: 2025-11-01T00:00:00Z
  }
  - Fetch preferences: GET /api/v1/users/me/preferences
- Update profile: PUT /api/v1/users/me/profile
- Update preferences: PUT /api/v1/users/me/preferences
- Change password: POST /api/v1/auth/change-password
- Enable 2FA: POST /api/v1/auth/2fa/enable → Returns QR code
- Manage connected accounts: GET /api/v1/auth/connections
- Upload avatar: POST /api/v1/users/me/avatar (multipart/form-data)
---
14. app/dashboard/progress/page.tsx + components/progress/ProgressTracker.tsx
Mock Data:
- Lines 5-66 (ProgressTracker): All 6 pathways with fake progress
- Lines 68-74: 5 hardcoded milestones
- Lines 235-242: Fake weekly activity hours
Required Changes:
- Fetch progress: GET /api/v1/users/me/progress
    {
    pathways: [
      {
        pathway_id: uuid,
        title: Python for Offensive Security,
        progress_percentage: 75,
        completed_modules: 3,
        total_modules: 5,
        xp_earned: 1250,
        streak_days: 7,
        last_activity: 2025-12-17T12:30:00Z
      }
    ],
    totals: {
      total_xp: 2850,
      xp_this_week: 320,
      overall_progress: 42,
      current_streak: 7,
      modules_completed: 7,
      modules_total: 33
    }
  }
  - Fetch milestones: GET /api/v1/users/me/milestones
- Fetch weekly activity: GET /api/v1/users/me/activity/weekly-hours
- Add filters for date ranges
- Add export progress feature: GET /api/v1/users/me/progress/export
---
15. app/dashboard/settings/page.tsx + components/settings/SettingsPanel.tsx
Mock Data:
- Lines 7-17 (SettingsPanel): Settings in local state only
- Lines 27-31: Hardcoded certificates
- Lines 252-260: Danger zone buttons non-functional
Required Changes:
- Fetch settings: GET /api/v1/users/me/settings
- Update settings: PUT /api/v1/users/me/settings
- Fetch certificates: GET /api/v1/users/me/certificates
    {
    certificates: [
      {
        id: uuid,
        title: Python Fundamentals,
        issued_date: 2025-11-15,
        verification_url: https://...,
        download_url: https://...,
        is_verified: true
      }
    ]
  }
  - Download certificate: GET /api/v1/certificates/{id}/download
- Delete account: DELETE /api/v1/users/me (with confirmation)
- Reset progress: POST /api/v1/users/me/progress/reset
- Export data: GET /api/v1/users/me/export (GDPR compliance)
---
16. components/learning/Judge0Integration.tsx
Mock Data:
- Not exactly mock, but incomplete integration
- Lines 12-39: Direct Judge0 calls from frontend
Required Changes:
- DO NOT call Judge0 directly from frontend
- All Judge0 calls should go through Go backend
- Backend proxy endpoints:
  - POST /api/v1/code/execute → Calls Judge0, stores submission
  - GET /api/v1/code/submissions/{id} → Returns cached result
- Frontend should only call Go backend
- Benefits:
  - Hide Judge0 credentials
  - Rate limiting
  - Usage tracking
  - Caching results
  - Better error handling
---
17. lib/judge0/service.ts
Issues:
- Lines 70-99: batchSubmit() defined but never used
- Lines 104-118: getSubmission() defined but never used
- Lines 123-133: healthCheck() defined but never used
- Direct Judge0 calls expose API URL to frontend
Required Changes:
- Remove this file entirely OR convert to types-only
- Move all Judge0 logic to Go backend
- Create TypeScript types for API responses:
    // types/api.ts
  export interface CodeExecutionRequest {
    source_code: string
    language_id: number
    stdin?: string
  }
  
  export interface CodeExecutionResult {
    submission_id: string
    stdout: string | null
    stderr: string | null
    compile_output: string | null
    status: {
      id: number
      description: string
    }
    time: string
    memory: number
  }
  
---
18. lib/supabase/client.ts + lib/supabase/server.ts
Current State: Properly implemented for auth
Issues:
- No user profile creation after signup
- No route protection
Required Changes:
- Keep Supabase for authentication only
- After Supabase signup, create user profile in Go backend:
    // In register flow
  const { data, error } = await supabase.auth.signUp({ email, password })
  if (!error && data.user) {
    await fetch('/api/v1/users', {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${data.session.access_token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        supabase_user_id: data.user.id,
        email: data.user.email
      })
    })
  }
  - Add middleware to validate Supabase JWT with Go backend
---
19. app/(auth)/register/page.tsx
Issues:
- Line 38: References non-existent /auth/callback route
- Line 48: Redirects to dashboard without creating profile
Required Changes:
1. Create callback route: app/auth/callback/route.ts
      import { createClient } from '@/lib/supabase/server'
   import { NextResponse } from 'next/server'
   export async function GET(request: Request) {
     const { searchParams } = new URL(request.url)
     const code = searchParams.get('code')
     
     if (code) {
       const supabase = await createClient()
       await supabase.auth.exchangeCodeForSession(code)
     }
     
     return NextResponse.redirect(new URL('/dashboard', request.url))
   }
   
2. After signup, create backend profile:
      const { error } = await supabase.auth.signUp({
     email,
     password,
     options: {
       emailRedirectTo: `${window.location.origin}/auth/callback`,
     },
   })
   if (!error) {
     // Get session token
     const { data: { session } } = await supabase.auth.getSession()
     
     // Create profile in Go backend
     await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/users`, {
       method: 'POST',
       headers: {
         'Authorization': `Bearer ${session?.access_token}`,
         'Content-Type': 'application/json'
       },
       body: JSON.stringify({
         email: email,
         display_name: email.split('@')[0]
       })
     })
     
     router.push('/dashboard')
   }
   
---
20. app/(auth)/login/page.tsx
Current State: Mostly functional
Required Changes:
- After login, sync with backend:
    const { error } = await supabase.auth.signInWithPassword({ email, password })
  
  if (!error) {
    const { data: { session } } = await supabase.auth.getSession()
    
    // Sync session with Go backend (optional, for session tracking)
    await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/session`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${session?.access_token}`,
      }
    })
    
    router.push('/dashboard')
  }
  
---
PART 2: GO BACKEND ARCHITECTURE & IMPLEMENTATION
Project Structure
wizardcore-backend/
├── cmd/
│   └── api/
│       └── main.go                    # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management
│   ├── database/
│   │   ├── postgres.go                # PostgreSQL connection
│   │   └── migrations/                # SQL migrations
│   │       ├── 001_users.sql
│   │       ├── 002_pathways.sql
│   │       ├── 003_exercises.sql
│   │       ├── 004_submissions.sql
│   │       ├── 005_progress.sql
│   │       ├── 006_achievements.sql
│   │       ├── 007_leaderboard.sql
│   │       └── 008_matches.sql
│   ├── models/
│   │   ├── user.go
│   │   ├── pathway.go
│   │   ├── exercise.go
│   │   ├── submission.go
│   │   ├── achievement.go
│   │   └── match.go
│   ├── repositories/
│   │   ├── user_repository.go
│   │   ├── pathway_repository.go
│   │   ├── exercise_repository.go
│   │   ├── submission_repository.go
│   │   ├── achievement_repository.go
│   │   └── match_repository.go
│   ├── services/
│   │   ├── auth_service.go            # Supabase JWT validation
│   │   ├── user_service.go
│   │   ├── pathway_service.go
│   │   ├── exercise_service.go
│   │   ├── submission_service.go      # Judge0 integration
│   │   ├── achievement_service.go     # Badge logic
│   │   ├── leaderboard_service.go
│   │   ├── progress_service.go
│   │   ├── match_service.go           # Duel matchmaking
│   │   └── notification_service.go
│   ├── handlers/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   ├── pathway_handler.go
│   │   ├── exercise_handler.go
│   │   ├── submission_handler.go
│   │   ├── achievement_handler.go
│   │   ├── leaderboard_handler.go
│   │   ├── practice_handler.go
│   │   └── websocket_handler.go
│   ├── middleware/
│   │   ├── auth.go                    # JWT validation
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── ratelimit.go
│   ├── router/
│   │   └── router.go                  # Route definitions
│   ├── websocket/
│   │   ├── hub.go                     # WebSocket hub
│   │   ├── client.go                  # WebSocket client
│   │   └── message.go                 # Message types
│   └── utils/
│       ├── errors.go
│       ├── pagination.go
│       └── validator.go
├── pkg/
│   ├── judge0/
│   │   └── client.go                  # Judge0 API client
│   └── supabase/
│       └── validator.go               # Supabase JWT validation
├── docs/
│   ├── api.md                         # API documentation
│   └── database_schema.md
├── scripts/
│   ├── migrate.sh
│   └── seed.sh
├── .env.example
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yml
└── README.md
---
Detailed Implementation Guide
Step 1: Initialize Go Project
mkdir wizardcore-backend
cd wizardcore-backend
go mod init github.com/yourusername/wizardcore-backend
# Install dependencies
go get -u github.com/gin-gonic/gin
go get -u github.com/lib/pq
go get -u github.com/golang-jwt/jwt/v5
go get -u github.com/joho/godotenv
go get -u github.com/go-playground/validator/v10
go get -u github.com/golang-migrate/migrate/v4
go get -u github.com/gorilla/websocket
go get -u github.com/go-redis/redis/v8
go get -u github.com/rs/cors
go get -u go.uber.org/zap
---
Step 2: Database Schema (PostgreSQL)
internal/database/migrations/001_users.sql
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    supabase_user_id UUID UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    display_name VARCHAR(100),
    avatar_url TEXT,
    bio TEXT,
    location VARCHAR(100),
    website VARCHAR(255),
    github_username VARCHAR(100),
    twitter_username VARCHAR(100),
    total_xp INTEGER DEFAULT 0,
    practice_score INTEGER DEFAULT 0,
    global_rank INTEGER,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_activity_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_users_supabase_id ON users(supabase_user_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_total_xp ON users(total_xp DESC);
-- Create user preferences table
CREATE TABLE IF NOT EXISTS user_preferences (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    theme VARCHAR(20) DEFAULT 'dark',
    language VARCHAR(10) DEFAULT 'en',
    email_notifications BOOLEAN DEFAULT true,
    push_notifications BOOLEAN DEFAULT false,
    public_profile BOOLEAN DEFAULT true,
    show_progress BOOLEAN DEFAULT true,
    auto_save BOOLEAN DEFAULT true,
    sound_effects BOOLEAN DEFAULT true,
    two_factor_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create user activity log
CREATE TABLE IF NOT EXISTS user_activities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    activity_type VARCHAR(50) NOT NULL, -- 'completion', 'practice', 'reading', 'achievement', 'streak'
    title VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color VARCHAR(50),
    metadata JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_activities_user_id ON user_activities(user_id);
CREATE INDEX idx_activities_created_at ON user_activities(created_at DESC);
---
internal/database/migrations/002_pathways.sql
-- Create pathways table
CREATE TABLE IF NOT EXISTS pathways (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    subtitle VARCHAR(255),
    description TEXT,
    level VARCHAR(50) NOT NULL, -- 'Beginner', 'Intermediate', 'Advanced', 'Expert'
    duration_weeks INTEGER NOT NULL,
    student_count INTEGER DEFAULT 0,
    rating DECIMAL(3,2) DEFAULT 0.0,
    module_count INTEGER DEFAULT 0,
    color_gradient VARCHAR(100),
    icon VARCHAR(10),
    is_locked BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    prerequisites UUID[], -- Array of pathway IDs
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create modules table
CREATE TABLE IF NOT EXISTS modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pathway_id UUID REFERENCES pathways(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    sort_order INTEGER NOT NULL,
    estimated_hours INTEGER,
    xp_reward INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_modules_pathway_id ON modules(pathway_id);
-- Create user pathway enrollments
CREATE TABLE IF NOT EXISTS user_pathway_enrollments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pathway_id UUID REFERENCES pathways(id) ON DELETE CASCADE,
    progress_percentage INTEGER DEFAULT 0,
    completed_modules INTEGER DEFAULT 0,
    xp_earned INTEGER DEFAULT 0,
    streak_days INTEGER DEFAULT 0,
    last_activity_at TIMESTAMP,
    enrolled_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    UNIQUE(user_id, pathway_id)
);
CREATE INDEX idx_enrollments_user_id ON user_pathway_enrollments(user_id);
CREATE INDEX idx_enrollments_pathway_id ON user_pathway_enrollments(pathway_id);
---
internal/database/migrations/003_exercises.sql
-- Create exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    difficulty VARCHAR(50) NOT NULL, -- 'BEGINNER', 'INTERMEDIATE', 'ADVANCED'
    points INTEGER DEFAULT 0,
    time_limit_minutes INTEGER,
    sort_order INTEGER NOT NULL,
    
    -- Lesson content
    objectives TEXT[],
    content TEXT, -- Markdown
    examples JSONB,
    
    -- Exercise details
    description TEXT,
    constraints TEXT[],
    hints TEXT[],
    starter_code TEXT,
    solution_code TEXT,
    
    -- Language
    language_id INTEGER NOT NULL, -- Judge0 language ID
    
    -- Metadata
    tags TEXT[],
    concurrent_solvers INTEGER DEFAULT 0,
    total_submissions INTEGER DEFAULT 0,
    total_completions INTEGER DEFAULT 0,
    average_completion_time INTEGER, -- in seconds
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_exercises_module_id ON exercises(module_id);
CREATE INDEX idx_exercises_difficulty ON exercises(difficulty);
-- Create test cases table
CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    input TEXT,
    expected_output TEXT,
    is_hidden BOOLEAN DEFAULT false,
    points INTEGER DEFAULT 0,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_test_cases_exercise_id ON test_cases(exercise_id);
---
internal/database/migrations/004_submissions.sql
-- Create submissions table
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    source_code TEXT NOT NULL,
    language_id INTEGER NOT NULL,
    
    -- Judge0 results
    judge0_token VARCHAR(255),
    status VARCHAR(50), -- 'pending', 'processing', 'accepted', 'wrong_answer', 'time_limit_exceeded', 'runtime_error', 'compilation_error'
    stdout TEXT,
    stderr TEXT,
    compile_output TEXT,
    execution_time DECIMAL(10,3), -- in seconds
    memory_used INTEGER, -- in KB
    
    -- Scoring
    test_cases_passed INTEGER DEFAULT 0,
    test_cases_total INTEGER DEFAULT 0,
    points_earned INTEGER DEFAULT 0,
    is_correct BOOLEAN DEFAULT false,
    
    -- Metadata
    submission_type VARCHAR(50) DEFAULT 'solution', -- 'draft', 'solution', 'practice'
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_exercise_id ON submissions(exercise_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_submissions_created_at ON submissions(created_at DESC);
-- Create submission test results
CREATE TABLE IF NOT EXISTS submission_test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID REFERENCES submissions(id) ON DELETE CASCADE,
    test_case_id UUID REFERENCES test_cases(id) ON DELETE CASCADE,
    passed BOOLEAN NOT NULL,
    actual_output TEXT,
    execution_time DECIMAL(10,3),
    memory_used INTEGER,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_test_results_submission_id ON submission_test_results(submission_id);
---
internal/database/migrations/005_progress.sql
-- Create user module progress
CREATE TABLE IF NOT EXISTS user_module_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    pathway_id UUID REFERENCES pathways(id) ON DELETE CASCADE,
    
    progress_percentage INTEGER DEFAULT 0,
    completed_exercises INTEGER DEFAULT 0,
    total_exercises INTEGER DEFAULT 0,
    xp_earned INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0,
    
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP,
    last_activity_at TIMESTAMP,
    
    UNIQUE(user_id, module_id)
);
CREATE INDEX idx_module_progress_user_id ON user_module_progress(user_id);
CREATE INDEX idx_module_progress_module_id ON user_module_progress(module_id);
-- Create user daily activity
CREATE TABLE IF NOT EXISTS user_daily_activity (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    activity_date DATE NOT NULL,
    exercises_completed INTEGER DEFAULT 0,
    xp_earned INTEGER DEFAULT 0,
    time_spent_minutes INTEGER DEFAULT 0,
    submissions_count INTEGER DEFAULT 0,
    streak_maintained BOOLEAN DEFAULT false,
    UNIQUE(user_id, activity_date)
);
CREATE INDEX idx_daily_activity_user_id ON user_daily_activity(user_id);
CREATE INDEX idx_daily_activity_date ON user_daily_activity(activity_date DESC);
-- Create user milestones
CREATE TABLE IF NOT EXISTS user_milestones (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    milestone_type VARCHAR(50) NOT NULL, -- 'exercise', 'module', 'pathway', 'streak', 'achievement'
    xp_awarded INTEGER DEFAULT 0,
    achieved_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_milestones_user_id ON user_milestones(user_id);
CREATE INDEX idx_milestones_achieved_at ON user_milestones(achieved_at DESC);
---
internal/database/migrations/006_achievements.sql
-- Create achievements (badges) table
CREATE TABLE IF NOT EXISTS achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    icon VARCHAR(50),
    color_gradient VARCHAR(100),
    rarity VARCHAR(50), -- 'Common', 'Uncommon', 'Rare', 'Epic', 'Legendary', 'Mythic'
    xp_reward INTEGER DEFAULT 0,
    
    -- Unlock criteria
    criteria_type VARCHAR(50) NOT NULL, -- 'exercise_count', 'pathway_complete', 'streak', 'speed', 'perfect_score', 'custom'
    criteria_value INTEGER,
    criteria_metadata JSONB,
    
    is_hidden BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create user achievements
CREATE TABLE IF NOT EXISTS user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    achievement_id UUID REFERENCES achievements(id) ON DELETE CASCADE,
    progress INTEGER DEFAULT 0, -- For achievements with progress tracking
    earned_at TIMESTAMP,
    UNIQUE(user_id, achievement_id)
);
CREATE INDEX idx_user_achievements_user_id ON user_achievements(user_id);
CREATE INDEX idx_user_achievements_earned_at ON user_achievements(earned_at DESC);
---
internal/database/migrations/007_leaderboard.sql
-- Create leaderboard entries (denormalized for performance)
CREATE TABLE IF NOT EXISTS leaderboard_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    timeframe VARCHAR(20) NOT NULL, -- 'all', 'month', 'week'
    pathway_id UUID REFERENCES pathways(id) ON DELETE SET NULL, -- NULL means global
    
    rank INTEGER NOT NULL,
    previous_rank INTEGER,
    xp INTEGER NOT NULL,
    streak_days INTEGER DEFAULT 0,
    badge_count INTEGER DEFAULT 0,
    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id, timeframe, pathway_id)
);
CREATE INDEX idx_leaderboard_timeframe ON leaderboard_entries(timeframe, rank);
CREATE INDEX idx_leaderboard_pathway ON leaderboard_entries(pathway_id, rank);
CREATE INDEX idx_leaderboard_user_id ON leaderboard_entries(user_id);
-- Create leaderboard update log (for tracking rank changes)
CREATE TABLE IF NOT EXISTS leaderboard_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    timeframe VARCHAR(20) NOT NULL,
    pathway_id UUID,
    rank INTEGER NOT NULL,
    xp INTEGER NOT NULL,
    recorded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_leaderboard_history_user_id ON leaderboard_history(user_id);
CREATE INDEX idx_leaderboard_history_recorded_at ON leaderboard_history(recorded_at DESC);
---
internal/database/migrations/008_matches.sql
-- Create practice matches (duels)
CREATE TABLE IF NOT EXISTS practice_matches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_type VARCHAR(50) NOT NULL, -- 'duel', 'speed_run', 'random'
    status VARCHAR(50) NOT NULL, -- 'pending', 'active', 'completed', 'cancelled'
    
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    time_limit_minutes INTEGER,
    
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Create match participants
CREATE TABLE IF NOT EXISTS match_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    match_id UUID REFERENCES practice_matches(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    submission_id UUID REFERENCES submissions(id) ON DELETE SET NULL,
    
    score INTEGER DEFAULT 0,
    rank INTEGER,
    result VARCHAR(20), -- 'win', 'loss', 'draw'
    xp_earned INTEGER DEFAULT 0,
    
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    finished_at TIMESTAMP
);
CREATE INDEX idx_match_participants_match_id ON match_participants(match_id);
CREATE INDEX idx_match_participants_user_id ON match_participants(user_id);
-- Create practice statistics
CREATE TABLE IF NOT EXISTS user_practice_stats (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    duels_total INTEGER DEFAULT 0,
    duels_won INTEGER DEFAULT 0,
    duels_lost INTEGER DEFAULT 0,
    duels_draw INTEGER DEFAULT 0,
    
    speed_runs_completed INTEGER DEFAULT 0,
    best_speed_run_time INTEGER, -- in seconds
    
    random_challenges_completed INTEGER DEFAULT 0,
    
    total_practice_xp INTEGER DEFAULT 0,
    practice_score INTEGER DEFAULT 0,
    practice_rank INTEGER,
    avg_completion_time INTEGER, -- in seconds
    
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
---
internal/database/migrations/009_notifications.sql
-- Create notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    type VARCHAR(50) NOT NULL, -- 'achievement', 'match_invite', 'deadline', 'streak', 'milestone'
    title VARCHAR(255) NOT NULL,
    message TEXT,
    icon VARCHAR(50),
    action_url TEXT,
    
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at DESC);
---
internal/database/migrations/010_deadlines.sql
-- Create user deadlines table
CREATE TABLE IF NOT EXISTS user_deadlines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    
    title VARCHAR(255) NOT NULL,
    description TEXT,
    deadline_type VARCHAR(50) NOT NULL, -- 'project', 'quiz', 'lab', 'assignment'
    
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    pathway_id UUID REFERENCES pathways(id) ON DELETE SET NULL,
    module_id UUID REFERENCES modules(id) ON DELETE SET NULL,
    
    due_date TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    is_completed BOOLEAN DEFAULT false,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_deadlines_user_id ON user_deadlines(user_id);
CREATE INDEX idx_deadlines_due_date ON user_deadlines(due_date);
CREATE INDEX idx_deadlines_is_completed ON user_deadlines(is_completed);
---
internal/database/migrations/011_certificates.sql
-- Create certificates table
CREATE TABLE IF NOT EXISTS certificates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    pathway_id UUID REFERENCES pathways(id) ON DELETE CASCADE,
    
    title VARCHAR(255) NOT NULL,
    description TEXT,
    certificate_number VARCHAR(100) UNIQUE NOT NULL,
    verification_url TEXT,
    download_url TEXT,
    
    is_verified BOOLEAN DEFAULT false,
    verified_at TIMESTAMP,
    
    issued_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_certificates_user_id ON certificates(user_id);
CREATE INDEX idx_certificates_certificate_number ON certificates(certificate_number);
---
Step 3: Go Models
internal/models/user.go
package models
import (
    "time"
    "github.com/google/uuid"
)
type User struct {
    ID              uuid.UUID  `json:"id" db:"id"`
    SupabaseUserID  uuid.UUID  `json:"supabase_user_id" db:"supabase_user_id"`
    Email           string     `json:"email" db:"email" validate:"required,email"`
    DisplayName     *string    `json:"display_name,omitempty" db:"display_name"`
    AvatarURL       *string    `json:"avatar_url,omitempty" db:"avatar_url"`
    Bio             *string    `json:"bio,omitempty" db:"bio"`
    Location        *string    `json:"location,omitempty" db:"location"`
    Website         *string    `json:"website,omitempty" db:"website"`
    GithubUsername  *string    `json:"github_username,omitempty" db:"github_username"`
    TwitterUsername *string    `json:"twitter_username,omitempty" db:"twitter_username"`
    TotalXP         int        `json:"total_xp" db:"total_xp"`
    PracticeScore   int        `json:"practice_score" db:"practice_score"`
    GlobalRank      *int       `json:"global_rank,omitempty" db:"global_rank"`
    CurrentStreak   int        `json:"current_streak" db:"current_streak"`
    LongestStreak   int        `json:"longest_streak" db:"longest_streak"`
    LastActivityDate *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
    CreatedAt       time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}
type UserPreferences struct {
    UserID             uuid.UUID `json:"user_id" db:"user_id"`
    Theme              string    `json:"theme" db:"theme"`
    Language           string    `json:"language" db:"language"`
    EmailNotifications bool      `json:"email_notifications" db:"email_notifications"`
    PushNotifications  bool      `json:"push_notifications" db:"push_notifications"`
    PublicProfile      bool      `json:"public_profile" db:"public_profile"`
    ShowProgress       bool      `json:"show_progress" db:"show_progress"`
    AutoSave           bool      `json:"auto_save" db:"auto_save"`
    SoundEffects       bool      `json:"sound_effects" db:"sound_effects"`
    TwoFactorEnabled   bool      `json:"two_factor_enabled" db:"two_factor_enabled"`
    CreatedAt          time.Time `json:"created_at" db:"created_at"`
    UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}
type UserActivity struct {
    ID           uuid.UUID              `json:"id" db:"id"`
    UserID       uuid.UUID              `json:"user_id" db:"user_id"`
    ActivityType string                 `json:"activity_type" db:"activity_type"`
    Title        string                 `json:"title" db:"title"`
    Description  *string                `json:"description,omitempty" db:"description"`
    Icon         *string                `json:"icon,omitempty" db:"icon"`
    Color        *string                `json:"color,omitempty" db:"color"`
    Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
    CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}
type UserStats struct {
    ActiveCourses          int    `json:"active_courses"`
    ActiveCoursesChange    string `json:"active_courses_change"`
    CompletionRate         int    `json:"completion_rate"`
    CompletionRateChange   string `json:"completion_rate_change"`
    StudyTimeHours         int    `json:"study_time_hours"`
    StudyTimeWeek          int    `json:"study_time_week"`
    XPTotal                int    `json:"xp_total"`
    XPToday                int    `json:"xp_today"`
}
type CreateUserRequest struct {
    SupabaseUserID uuid.UUID `json:"supabase_user_id" validate:"required"`
    Email          string    `json:"email" validate:"required,email"`
    DisplayName    string    `json:"display_name"`
}
type UpdateUserProfileRequest struct {
    DisplayName     *string `json:"display_name"`
    Bio             *string `json:"bio"`
    Location        *string `json:"location"`
    Website         *string `json:"website"`
    GithubUsername  *string `json:"github_username"`
    TwitterUsername *string `json:"twitter_username"`
}
type UpdatePreferencesRequest struct {
    Theme              *string `json:"theme"`
    Language           *string `json:"language"`
    EmailNotifications *bool   `json:"email_notifications"`
    PushNotifications  *bool   `json:"push_notifications"`
    PublicProfile      *bool   `json:"public_profile"`
    ShowProgress       *bool   `json:"show_progress"`
    AutoSave           *bool   `json:"auto_save"`
    SoundEffects       *bool   `json:"sound_effects"`
}
---
internal/models/pathway.go
package models
import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)
type Pathway struct {
    ID             uuid.UUID      `json:"id" db:"id"`
    Title          string         `json:"title" db:"title"`
    Subtitle       *string        `json:"subtitle,omitempty" db:"subtitle"`
    Description    *string        `json:"description,omitempty" db:"description"`
    Level          string         `json:"level" db:"level"`
    DurationWeeks  int            `json:"duration_weeks" db:"duration_weeks"`
    StudentCount   int            `json:"student_count" db:"student_count"`
    Rating         float64        `json:"rating" db:"rating"`
    ModuleCount    int            `json:"module_count" db:"module_count"`
    ColorGradient  *string        `json:"color_gradient,omitempty" db:"color_gradient"`
    Icon           *string        `json:"icon,omitempty" db:"icon"`
    IsLocked       bool           `json:"is_locked" db:"is_locked"`
    SortOrder      int            `json:"sort_order" db:"sort_order"`
    Prerequisites  pq.StringArray `json:"prerequisites" db:"prerequisites"`
    CreatedAt      time.Time      `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}
type PathwayWithEnrollment struct {
    Pathway
    IsEnrolled bool `json:"is_enrolled"`
    Progress   int  `json:"progress"`
}
type Module struct {
    ID             uuid.UUID  `json:"id" db:"id"`
    PathwayID      uuid.UUID  `json:"pathway_id" db:"pathway_id"`
    Title          string     `json:"title" db:"title"`
    Description    *string    `json:"description,omitempty" db:"description"`
    SortOrder      int        `json:"sort_order" db:"sort_order"`
    EstimatedHours *int       `json:"estimated_hours,omitempty" db:"estimated_hours"`
    XPReward       int        `json:"xp_reward" db:"xp_reward"`
    CreatedAt      time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}
type UserPathwayEnrollment struct {
    ID                uuid.UUID  `json:"id" db:"id"`
    UserID            uuid.UUID  `json:"user_id" db:"user_id"`
    PathwayID         uuid.UUID  `json:"pathway_id" db:"pathway_id"`
    ProgressPercentage int       `json:"progress_percentage" db:"progress_percentage"`
    CompletedModules  int        `json:"completed_modules" db:"completed_modules"`
    XPEarned          int        `json:"xp_earned" db:"xp_earned"`
    StreakDays        int        `json:"streak_days" db:"streak_days"`
    LastActivityAt    *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
    EnrolledAt        time.Time  `json:"enrolled_at" db:"enrolled_at"`
    CompletedAt       *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}
type PathwayProgress struct {
    PathwayID         uuid.UUID  `json:"pathway_id"`
    Title             string     `json:"title"`
    Progress          int        `json:"progress_percentage"`
    CompletedModules  int        `json:"completed_modules"`
    TotalModules      int        `json:"total_modules"`
    XPEarned          int        `json:"xp_earned"`
    StreakDays        int        `json:"streak_days"`
    LastActivity      *time.Time `json:"last_activity,omitempty"`
}
---
internal/models/exercise.go
package models
import (
    "time"
    "github.com/google/uuid"
    "github.com/lib/pq"
)
type Exercise struct {
    ID                     uuid.UUID              `json:"id" db:"id"`
    ModuleID               uuid.UUID              `json:"module_id" db:"module_id"`
    Title                  string                 `json:"title" db:"title"`
    Difficulty             string                 `json:"difficulty" db:"difficulty"`
    Points                 int                    `json:"points" db:"points"`
    TimeLimitMinutes       *int                   `json:"time_limit_minutes,omitempty" db:"time_limit_minutes"`
    SortOrder              int                    `json:"sort_order" db:"sort_order"`
    Objectives             pq.StringArray         `json:"objectives" db:"objectives"`
    Content                *string                `json:"content,omitempty" db:"content"`
    Examples               map[string]interface{} `json:"examples,omitempty" db:"examples"`
    Description            *string                `json:"description,omitempty" db:"description"`
    Constraints            pq.StringArray         `json:"constraints" db:"constraints"`
    Hints                  pq.StringArray         `json:"hints" db:"hints"`
    StarterCode            *string                `json:"starter_code,omitempty" db:"starter_code"`
    SolutionCode           *string                `json:"solution_code,omitempty" db:"solution_code"`
    LanguageID             int                    `json:"language_id" db:"language_id"`
    Tags                   pq.StringArray         `json:"tags" db:"tags"`
    ConcurrentSolvers      int                    `json:"concurrent_solvers" db:"concurrent_solvers"`
    TotalSubmissions       int                    `json:"total_submissions" db:"total_submissions"`
    TotalCompletions       int                    `json:"total_completions" db:"total_completions"`
    AvgCompletionTime      *int                   `json:"average_completion_time,omitempty" db:"average_completion_time"`
    CreatedAt              time.Time              `json:"created_at" db:"created_at"`
    UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}
type TestCase struct {
    ID             uuid.UUID `json:"id" db:"id"`
    ExerciseID     uuid.UUID `json:"exercise_id" db:"exercise_id"`
    Input          *string   `json:"input,omitempty" db:"input"`
    ExpectedOutput string    `json:"expected_output" db:"expected_output"`
    IsHidden       bool      `json:"is_hidden" db:"is_hidden"`
    Points         int       `json:"points" db:"points"`
    SortOrder      int       `json:"sort_order" db:"sort_order"`
    CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
type ExerciseWithTests struct {
    Exercise
    TestCases []TestCase `json:"test_cases"`
}
type ExerciseStats struct {
    ConcurrentSolvers int `json:"concurrent_solvers"`
    TotalSubmissions  int `json:"total_submissions"`
    CompletionRate    int `json:"completion_rate"`
}
---
internal/models/submission.go
package models
import (
    "time"
    "github.com/google/uuid"
)
type Submission struct {
    ID               uuid.UUID  `json:"id" db:"id"`
    UserID           uuid.UUID  `json:"user_id" db:"user_id"`
    ExerciseID       uuid.UUID  `json:"exercise_id" db:"exercise_id"`
    SourceCode       string     `json:"source_code" db:"source_code"`
    LanguageID       int        `json:"language_id" db:"language_id"`
    Judge0Token      *string    `json:"judge0_token,omitempty" db:"judge0_token"`
    Status           string     `json:"status" db:"status"`
    Stdout           *string    `json:"stdout,omitempty" db:"stdout"`
    Stderr           *string    `json:"stderr,omitempty" db:"stderr"`
    CompileOutput    *string    `json:"compile_output,omitempty" db:"compile_output"`
    ExecutionTime    *float64   `json:"execution_time,omitempty" db:"execution_time"`
    MemoryUsed       *int       `json:"memory_used,omitempty" db:"memory_used"`
    TestCasesPassed  int        `json:"test_cases_passed" db:"test_cases_passed"`
    TestCasesTotal   int        `json:"test_cases_total" db:"test_cases_total"`
    PointsEarned     int        `json:"points_earned" db:"points_earned"`
    IsCorrect        bool       `json:"is_correct" db:"is_correct"`
    SubmissionType   string     `json:"submission_type" db:"submission_type"`
    IPAddress        *string    `json:"ip_address,omitempty" db:"ip_address"`
    UserAgent        *string    `json:"user_agent,omitempty" db:"user_agent"`
    CreatedAt        time.Time  `json:"created_at" db:"created_at"`
    UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}
type SubmissionTestResult struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    SubmissionID  uuid.UUID  `json:"submission_id" db:"submission_id"`
    TestCaseID    uuid.UUID  `json:"test_case_id" db:"test_case_id"`
    Passed        bool       `json:"passed" db:"passed"`
    ActualOutput  *string    `json:"actual_output,omitempty" db:"actual_output"`
    ExecutionTime *float64   `json:"execution_time,omitempty" db:"execution_time"`
    MemoryUsed    *int       `json:"memory_used,omitempty" db:"memory_used"`
    ErrorMessage  *string    `json:"error_message,omitempty" db:"error_message"`
    CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}
type CreateSubmissionRequest struct {
    ExerciseID uuid.UUID `json:"exercise_id" validate:"required"`
    SourceCode string    `json:"source_code" validate:"required"`
    LanguageID int       `json:"language_id" validate:"required"`
}
type SubmissionResponse struct {
    Submission
    TestResults []SubmissionTestResult `json:"test_results,omitempty"`
}
---
internal/models/achievement.go
package models
import (
    "time"
    "github.com/google/uuid"
)
type Achievement struct {
    ID               uuid.UUID              `json:"id" db:"id"`
    Title            string                 `json:"title" db:"title"`
    Description      *string                `json:"description,omitempty" db:"description"`
    Icon             *string                `json:"icon,omitempty" db:"icon"`
    ColorGradient    *string                `json:"color_gradient,omitempty" db:"color_gradient"`
    Rarity           string                 `json:"rarity" db:"rarity"`
    XPReward         int                    `json:"xp_reward" db:"xp_reward"`
    CriteriaType     string                 `json:"criteria_type" db:"criteria_type"`
    CriteriaValue    *int                   `json:"criteria_value,omitempty" db:"criteria_value"`
    CriteriaMetadata map[string]interface{} `json:"criteria_metadata,omitempty" db:"criteria_metadata"`
    IsHidden         bool                   `json:"is_hidden" db:"is_hidden"`
    SortOrder        int                    `json:"sort_order" db:"sort_order"`
    CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}
type UserAchievement struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    UserID        uuid.UUID  `json:"user_id" db:"user_id"`
    AchievementID uuid.UUID  `json:"achievement_id" db:"achievement_id"`
    Progress      int        `json:"progress" db:"progress"`
    EarnedAt      *time.Time `json:"earned_at,omitempty" db:"earned_at"`
}
type AchievementWithProgress struct {
    Achievement
    Earned     bool       `json:"earned"`
    Progress   int        `json:"progress"`
    EarnedDate *time.Time `json:"earned_date,omitempty"`
}
---
internal/models/leaderboard.go
package models
import (
    "time"
    "github.com/google/uuid"
)
type LeaderboardEntry struct {
    ID            uuid.UUID  `json:"id" db:"id"`
    UserID        uuid.UUID  `json:"user_id" db:"user_id"`
    Username      string     `json:"username" db:"username"`
    AvatarURL     *string    `json:"avatar_url,omitempty" db:"avatar_url"`
    Timeframe     string     `json:"timeframe" db:"timeframe"`
    PathwayID     *uuid.UUID `json:"pathway_id,omitempty" db:"pathway_id"`
    Rank          int        `json:"rank" db:"rank"`
    PreviousRank  *int       `json:"previous_rank,omitempty" db:"previous_rank"`
    XP            int        `json:"xp" db:"xp"`
    StreakDays    int        `json:"streak_days" db:"streak_days"`
    BadgeCount    int        `json:"badge_count" db:"badge_count"`
    CountryCode   *string    `json:"country_code,omitempty"`
    Trend         string     `json:"trend"` // up, down, same
    IsCurrentUser bool       `json:"is_current_user"`
    UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}
type LeaderboardResponse struct {
    Leaderboard []LeaderboardEntry `json:"leaderboard"`
    Stats       LeaderboardStats   `json:"stats"`
    Pagination  Pagination         `json:"pagination"`
}
type LeaderboardStats struct {
    TotalLearners      int     `json:"total_learners"`
    CurrentUserRank    int     `json:"current_user_rank"`
    CurrentUserChange  int     `json:"current_user_change"`
    TopXP              int     `json:"top_xp"`
    TopUsername        string  `json:"top_username"`
    CountryCount       int     `json:"country_count"`
}
type Pagination struct {
    Total   int `json:"total"`
    Page    int `json:"page"`
    PerPage int `json:"per_page"`
}
---
Step 4: Core Go Service Files
cmd/api/main.go
package main
import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
    "github.com/yourusername/wizardcore-backend/internal/config"
    "github.com/yourusername/wizardcore-backend/internal/database"
    "github.com/yourusername/wizardcore-backend/internal/router"
    
    "github.com/joho/godotenv"
    "go.uber.org/zap"
)
func main() {
    // Load environment variables
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
    // Initialize logger
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to initialize logger: %v", err)
    }
    defer logger.Sync()
    // Load configuration
    cfg, err := config.Load()
    if err != nil {
        logger.Fatal("Failed to load configuration", zap.Error(err))
    }
    // Connect to database
    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        logger.Fatal("Failed to connect to database", zap.Error(err))
    }
    defer db.Close()
    // Run migrations
    if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
        logger.Fatal("Failed to run migrations", zap.Error(err))
    }
    // Initialize router
    r := router.Setup(db, cfg, logger)
    // Create HTTP server
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%d", cfg.Port),
        Handler:      r,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    // Start server in a goroutine
    go func() {
        logger.Info("Starting server", zap.Int("port", cfg.Port))
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed to start", zap.Error(err))
        }
    }()
    // Wait for interrupt signal to gracefully shutdown the server
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit
    logger.Info("Shutting down server...")
    // Graceful shutdown with 5 second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Server forced to shutdown", zap.Error(err))
    }
    logger.Info("Server exited")
}
---
internal/config/config.go
package config
import (
    "fmt"
    "os"
    "strconv"
)
type Config struct {
    Port                int
    DatabaseURL         string
    SupabaseURL         string
    SupabaseJWTSecret   string
    Judge0APIURL        string
    Judge0APIKey        string
    RedisURL            string
    CORSAllowedOrigins  []string
    Environment         string
    LogLevel            string
}
func Load() (*Config, error) {
    port, err := strconv.Atoi(getEnv("PORT", "8080"))
    if err != nil {
        return nil, fmt.Errorf("invalid PORT: %w", err)
    }
    cfg := &Config{
        Port:                port,
        DatabaseURL:         getEnv("DATABASE_URL", ""),
        SupabaseURL:         getEnv("SUPABASE_URL", ""),
        SupabaseJWTSecret:   getEnv("SUPABASE_JWT_SECRET", ""),
        Judge0APIURL:        getEnv("JUDGE0_API_URL", "http://localhost:2358"),
        Judge0APIKey:        getEnv("JUDGE0_API_KEY", ""),
        RedisURL:            getEnv("REDIS_URL", "localhost:6379"),
        CORSAllowedOrigins:  []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
        Environment:         getEnv("ENVIRONMENT", "development"),
        LogLevel:            getEnv("LOG_LEVEL", "info"),
    }
    if cfg.DatabaseURL == "" {
        return nil, fmt.Errorf("DATABASE_URL is required")
    }
    if cfg.SupabaseJWTSecret == "" {
        return nil, fmt.Errorf("SUPABASE_JWT_SECRET is required")
    }
    return cfg, nil
}
func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
---
internal/database/postgres.go
package database
import (
    "database/sql"
    "fmt"
    "time"
    _ "github.com/lib/pq"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)
func Connect(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }
    // Configure connection pool
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    db.SetConnMaxLifetime(5 * time.Minute)
    // Test connection
    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }
    return db, nil
}
func RunMigrations(databaseURL string) error {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return fmt.Errorf("failed to open database for migrations: %w", err)
    }
    defer db.Close()
    driver, err := postgres.WithInstance(db, &postgres.Config{})
    if err != nil {
        return fmt.Errorf("failed to create migration driver: %w", err)
    }
    m, err := migrate.NewWithDatabaseInstance(
        "file://internal/database/migrations",
        "postgres",
        driver,
    )
    if err != nil {
        return fmt.Errorf("failed to create migration instance: %w", err)
    }
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("failed to run migrations: %w", err)
    }
    return nil
}
---
internal/middleware/auth.go
package middleware
import (
    "context"
    "net/http"
    "strings"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)
type contextKey string
const UserIDKey contextKey = "user_id"
const SupabaseUserIDKey contextKey = "supabase_user_id"
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
            c.Abort()
            return
        }
        // Extract token from "Bearer <token>"
        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
            c.Abort()
            return
        }
        tokenString := parts[1]
        // Parse and validate JWT
        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }
        // Extract claims
        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }
        // Get Supabase user ID from claims
        subStr, ok := claims["sub"].(string)
        if !ok {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
            c.Abort()
            return
        }
        supabaseUserID, err := uuid.Parse(subStr)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID format"})
            c.Abort()
            return
        }
        // Store in context
        ctx := context.WithValue(c.Request.Context(), SupabaseUserIDKey, supabaseUserID)
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}
func GetSupabaseUserID(c *gin.Context) (uuid.UUID, bool) {
    userID, ok := c.Request.Context().Value(SupabaseUserIDKey).(uuid.UUID)
    return userID, ok
}
---
internal/router/router.go
package router
import (
    "database/sql"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/yourusername/wizardcore-backend/internal/config"
    "github.com/yourusername/wizardcore-backend/internal/handlers"
    "github.com/yourusername/wizardcore-backend/internal/middleware"
    "github.com/yourusername/wizardcore-backend/internal/repositories"
    "github.com/yourusername/wizardcore-backend/internal/services"
    "github.com/yourusername/wizardcore-backend/pkg/judge0"
    
    "go.uber.org/zap"
)
func Setup(db *sql.DB, cfg *config.Config, logger *zap.Logger) *gin.Engine {
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    r := gin.Default()
    // CORS
    r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))
    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })
    // Initialize repositories
    userRepo := repositories.NewUserRepository(db)
    pathwayRepo := repositories.NewPathwayRepository(db)
    exerciseRepo := repositories.NewExerciseRepository(db)
    submissionRepo := repositories.NewSubmissionRepository(db)
    achievementRepo := repositories.NewAchievementRepository(db)
    // Initialize Judge0 client
    judge0Client := judge0.NewClient(cfg.Judge0APIURL, cfg.Judge0APIKey)
    // Initialize services
    userService := services.NewUserService(userRepo)
    pathwayService := services.NewPathwayService(pathwayRepo, userRepo)
    exerciseService := services.NewExerciseService(exerciseRepo)
    submissionService := services.NewSubmissionService(submissionRepo, exerciseRepo, userRepo, judge0Client)
    achievementService := services.NewAchievementService(achievementRepo, userRepo)
    leaderboardService := services.NewLeaderboardService(userRepo)
    progressService := services.NewProgressService(userRepo, pathwayRepo, exerciseRepo)
    // Initialize handlers
    authHandler := handlers.NewAuthHandler(userService, logger)
    userHandler := handlers.NewUserHandler(userService, logger)
    pathwayHandler := handlers.NewPathwayHandler(pathwayService, logger)
    exerciseHandler := handlers.NewExerciseHandler(exerciseService, logger)
    submissionHandler := handlers.NewSubmissionHandler(submissionService, logger)
    achievementHandler := handlers.NewAchievementHandler(achievementService, logger)
    leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService, logger)
    practiceHandler := handlers.NewPracticeHandler(logger)
    // API routes
    api := r.Group("/api/v1")
    {
        // Public routes
        api.POST("/users", authHandler.CreateUser)
        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware(cfg.SupabaseJWTSecret))
        {
            // User routes
            protected.GET("/users/me", userHandler.GetCurrentUser)
            protected.PUT("/users/me/profile", userHandler.UpdateProfile)
            protected.GET("/users/me/preferences", userHandler.GetPreferences)
            protected.PUT("/users/me/preferences", userHandler.UpdatePreferences)
            protected.GET("/users/me/stats", userHandler.GetStats)
            protected.GET("/users/me/activities", userHandler.GetActivities)
            protected.GET("/users/me/nav-counts", userHandler.GetNavCounts)
            protected.DELETE("/users/me", userHandler.DeleteAccount)
            protected.GET("/users/me/export", userHandler.ExportData)
            // Pathway routes
            protected.GET("/pathways", pathwayHandler.GetAllPathways)
            protected.GET("/pathways/:id", pathwayHandler.GetPathway)
            protected.POST("/pathways/:id/enroll", pathwayHandler.EnrollPathway)
            protected.GET("/users/me/pathways", pathwayHandler.GetUserPathways)
            protected.GET("/users/me/deadlines", pathwayHandler.GetDeadlines)
            // Exercise routes
            protected.GET("/exercises/:id", exerciseHandler.GetExercise)
            protected.GET("/exercises/:id/stats", exerciseHandler.GetExerciseStats)
            // Submission routes
            protected.POST("/submissions", submissionHandler.CreateSubmission)
            protected.GET("/submissions/:id", submissionHandler.GetSubmission)
            protected.GET("/submissions/:exercise_id/latest", submissionHandler.GetLatestSubmission)
            protected.POST("/submissions/:exercise_id/save-draft", submissionHandler.SaveDraft)
            // Achievement routes
            protected.GET("/users/me/achievements", achievementHandler.GetUserAchievements)
            // Leaderboard routes
            protected.GET("/leaderboard", leaderboardHandler.GetLeaderboard)
            // Progress routes
            protected.GET("/users/me/progress", progressService.GetUserProgress)
            protected.GET("/users/me/milestones", progressService.GetMilestones)
            protected.GET("/users/me/activity/weekly", progressService.GetWeeklyActivity)
            protected.GET("/users/me/activity/weekly-hours", progressService.GetWeeklyHours)
            // Practice routes
            protected.GET("/practice/challenges", practiceHandler.GetChallenges)
            protected.GET("/practice/areas", practiceHandler.GetAreas)
            protected.GET("/users/me/practice/stats", practiceHandler.GetStats)
            protected.GET("/users/me/matches", practiceHandler.GetRecentMatches)
            protected.POST("/practice/challenges/:type/start", practiceHandler.StartChallenge)
        }
    }
    return r}
---
pkg/judge0/client.go
package judge0
import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)
type Client struct {
    baseURL string
    apiKey  string
    client  *http.Client
}
type Submission struct {
    SourceCode     string `json:"source_code"`
    LanguageID     int    `json:"language_id"`
    Stdin          string `json:"stdin,omitempty"`
    ExpectedOutput string `json:"expected_output,omitempty"`
}
type SubmissionResult struct {
    Token         string  `json:"token"`
    Stdout        *string `json:"stdout"`
    Stderr        *string `json:"stderr"`
    CompileOutput *string `json:"compile_output"`
    Message       *string `json:"message"`
    Status        Status  `json:"status"`
    Time          *string `json:"time"`
    Memory        *int    `json:"memory"`
}
type Status struct {
    ID          int    `json:"id"`
    Description string `json:"description"`
}
func NewClient(baseURL, apiKey string) *Client {
    return &Client{
        baseURL: baseURL,
        apiKey:  apiKey,
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}
func (c *Client) Submit(submission Submission) (*SubmissionResult, error) {
    url := fmt.Sprintf("%s/submissions?base64_encoded=false&wait=true", c.baseURL)
    jsonData, err := json.Marshal(submission)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal submission: %w", err)
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    req.Header.Set("Content-Type", "application/json")
    if c.apiKey != "" {
        req.Header.Set("X-RapidAPI-Key", c.apiKey)
    }
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
        return nil, fmt.Errorf("judge0 returned error: %s - %s", resp.Status, string(body))
    }
    var result SubmissionResult
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }
    return &result, nil
}
func (c *Client) GetSubmission(token string) (*SubmissionResult, error) {
    url := fmt.Sprintf("%s/submissions/%s?base64_encoded=false", c.baseURL, token)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    if c.apiKey != "" {
        req.Header.Set("X-RapidAPI-Key", c.apiKey)
    }
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("judge0 returned error: %s - %s", resp.Status, string(body))
    }
    var result SubmissionResult
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, fmt.Errorf("failed to unmarshal response: %w", err)
    }
    return &result, nil
}
func (c *Client) HealthCheck() error {
    url := fmt.Sprintf("%s/about", c.baseURL)
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }
    resp, err := c.client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send request: %w", err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("judge0 health check failed: %s", resp.Status)
    }
    return nil
}
---
Step 5: Environment Variables
.env.example
# Server
PORT=8080
ENVIRONMENT=development
LOG_LEVEL=info
# Database
DATABASE_URL=postgresql://username:password@localhost:5432/wizardcore?sslmode=disable
# Supabase
SUPABASE_URL=https://your-project.supabase.co
SUPABASE_JWT_SECRET=your-jwt-secret-from-supabase
# Judge0
JUDGE0_API_URL=http://localhost:2358
JUDGE0_API_KEY=
# Redis (for caching and WebSockets)
REDIS_URL=localhost:6379
# Frontend
FRONTEND_URL=http://localhost:3000
# Features
ENABLE_WEBSOCKET=true
ENABLE_RATE_LIMITING=true
---
Step 6: Next.js Environment Variables
Frontend .env.local
# Supabase
NEXT_PUBLIC_SUPABASE_URL=https://your-project.supabase.co
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-anon-key
# Backend API
NEXT_PUBLIC_API_URL=http://localhost:8080
# Judge0 (deprecated - now proxied through backend)
# NEXT_PUBLIC_JUDGE0_API_URL=http://localhost:2358
---
Step 7: Docker Compose Setup
docker-compose.yml
version: '3.8'
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: wizardcore
      POSTGRES_PASSWORD: wizardcore
      POSTGRES_DB: wizardcore
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
  judge0:
    image: judge0/judge0:latest
    environment:
      - REDIS_HOST=redis
      - REDIS_PORT=6379
      - POSTGRES_HOST=judge0-postgres
      - POSTGRES_PORT=5432
      - POSTGRES_DB=judge0
      - POSTGRES_USER=judge0
      - POSTGRES_PASSWORD=judge0
    ports:
      - "2358:2358"
    depends_on:
      - judge0-postgres
      - redis
  judge0-postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: judge0
      POSTGRES_PASSWORD: judge0
      POSTGRES_DB: judge0
    volumes:
      - judge0_postgres_data:/var/lib/postgresql/data
  backend:
    build: ./wizardcore-backend
    ports:
      - "8080:8080"
    environment:
      - PORT=8080
      - DATABASE_URL=postgresql://wizardcore:wizardcore@postgres:5432/wizardcore?sslmode=disable
      - REDIS_URL=redis:6379
      - JUDGE0_API_URL=http://judge0:2358
      - FRONTEND_URL=http://localhost:3000
    depends_on:
      - postgres
      - redis
      - judge0
    volumes:
      - ./wizardcore-backend:/app
volumes:
  postgres_data:
  redis_data:
  judge0_postgres_data:
---
Step 8: Implementation Checklist
Phase 1: Foundation (Week 1)
- [ ] Set up Go project structure
- [ ] Create all database migrations
- [ ] Implement database connection and migrations
- [ ] Set up configuration management
- [ ] Create models for all entities
- [ ] Implement authentication middleware (Supabase JWT validation)
- [ ] Create basic router setup
- [ ] Set up Docker Compose for local development
Phase 2: Core API Endpoints (Week 2-3)
- [ ] Implement user repository and service
- [ ] Create user endpoints (create, get, update, delete)
- [ ] Implement pathway repository and service
- [ ] Create pathway endpoints (list, get, enroll)
- [ ] Implement exercise repository and service
- [ ] Create exercise endpoints (get, list by module)
- [ ] Set up test database and write integration tests
Phase 3: Submission System (Week 4)
- [ ] Implement Judge0 client
- [ ] Create submission repository
- [ ] Implement submission service with Judge0 integration
- [ ] Create submission endpoints (create, get, list)
- [ ] Implement test case validation logic
- [ ] Add XP and progress update triggers
- [ ] Add draft save functionality
Phase 4: Progress & Achievements (Week 5)
- [ ] Implement progress tracking service
- [ ] Create progress endpoints
- [ ] Implement achievement system
- [ ] Create achievement unlock logic
- [ ] Implement daily activity tracking
- [ ] Create milestone system
- [ ] Build leaderboard update system (cron job)
Phase 5: Practice Arena (Week 6)
- [ ] Implement match repository and service
- [ ] Create matchmaking logic
- [ ] Implement WebSocket hub for real-time duels
- [ ] Create practice stats tracking
- [ ] Build speed run logic
- [ ] Implement random challenge generator
Phase 6: Frontend Integration (Week 7-8)
- [ ] Create auth callback route in Next.js
- [ ] Update all frontend components to use API
- [ ] Implement error handling and loading states
- [ ] Add toast notifications for user feedback
- [ ] Implement data fetching with SWR or React Query
- [ ] Add WebSocket connections for real-time features
- [ ] Update all forms to submit to backend
- [ ] Remove all mock data from frontend
Phase 7: Additional Features (Week 9-10)
- [ ] Implement notification system
- [ ] Create deadline management
- [ ] Build certificate generation
- [ ] Add data export functionality
- [ ] Implement 2FA system
- [ ] Add rate limiting
- [ ] Implement caching with Redis
- [ ] Add search functionality
Phase 8: Testing & Optimization (Week 11-12)
- [ ] Write unit tests for all services
- [ ] Write integration tests for all endpoints
- [ ] Perform load testing
- [ ] Optimize database queries
- [ ] Add database indexes
- [ ] Implement API response caching
- [ ] Add monitoring and logging
- [ ] Create API documentation
Phase 9: Deployment (Week 13)
- [ ] Set up production database (Supabase PostgreSQL)
- [ ] Deploy Go backend (Railway, Render, or AWS)
- [ ] Deploy Judge0 instance
- [ ] Configure environment variables
- [ ] Set up CI/CD pipeline
- [ ] Configure CORS properly
- [ ] Set up SSL certificates
- [ ] Deploy frontend to Vercel
Phase 10: Content Population (Week 14)
- [ ] Create seed scripts for pathways
- [ ] Add modules to each pathway
- [ ] Create exercises with test cases
- [ ] Add starter code for each exercise
- [ ] Configure achievements
- [ ] Test complete user journey
- [ ] Fix any bugs found in testing
---
Summary
This comprehensive plan provides:
1. 20 frontend files that need backend integration
2. Complete database schema with 11 migration files
3. Full Go backend structure with models, repositories, services, and handlers
4. Judge0 integration proxied through backend
5. Authentication flow with Supabase JWT validation
6. WebSocket support for real-time features
7. 14-week implementation timeline
The backend will be built in Go with:
- PostgreSQL for data storage
- Redis for caching and WebSocket state
- Judge0 for code execution
- Gin framework for HTTP routing
- Clean architecture with separation of concerns