# WIZARDCORE COMPREHENSIVE TODO LIST

*Last Updated: December 27, 2025 - Complete Audit Update*
*Status: In Progress - Critical Frontend Mock Data Identified*

## LEGEND
- [ ] = Not started
- [x] = Completed
- [-] = In progress
- [!] = Blocked/needs attention

---

## 🚨 CRITICAL PRIORITY (Core Product Non-Functional)

### 1. LEARNING ENVIRONMENT - ENTIRE CORE PRODUCT IS MOCKED
**File**: `components/learning/LearningEnvironment.tsx`
**Status**: 🔴 CRITICAL - The entire learning experience is fake!

- [ ] **Backend API**: Ensure `/exercises/:id` endpoint returns complete exercise data
  - [ ] Include: title, module, objectives, content, examples, description, constraints, test cases, hints, starter code
  - [ ] Add language_id, difficulty, points, time_limit fields
  - [ ] Return concurrent_solvers count (real-time or cached)
- [ ] **Backend API**: Implement `/submissions/:exercise_id/latest` to get user's saved code
- [ ] **Backend API**: Complete `/submissions` POST endpoint with Judge0 integration
  - [ ] Execute code via Judge0
  - [ ] Run against all test cases
  - [ ] Award XP on success
  - [ ] Record activity and progress
- [ ] **Frontend**: Replace hardcoded lesson object (lines 26-61)
- [ ] **Frontend**: Replace hardcoded exercise object (lines 63-89)
- [ ] **Frontend**: Replace hardcoded starter code (lines 10-24)
- [ ] **Frontend**: Implement real handleSubmit() function (currently just console.log at line 113-116)
- [ ] **Frontend**: Fetch exercise data on component mount
- [ ] **Frontend**: Auto-save user code every 30 seconds
- [ ] **Testing**: End-to-end test of submission flow with Judge0

**Impact**: Users cannot actually complete exercises or learn!

---

### 2. PATHWAY ENROLLMENT SYSTEM - USERS CAN'T ENROLL
**File**: `components/pathways/PathwaySelection.tsx`
**Status**: 🔴 CRITICAL - Users cannot join courses!

- [ ] **Backend API**: Implement `POST /pathways/:id/enroll`
  - [ ] Create user_pathway_enrollments record
  - [ ] Return enrollment status
  - [ ] Send welcome notification
  - [ ] Record activity
- [ ] **Backend API**: Update `GET /pathways` to include enrollment status
  - [ ] Add `is_enrolled` boolean to response
  - [ ] Add `progress` percentage if enrolled
- [ ] **Frontend**: Remove hardcoded pathways array (lines 6-97)
- [ ] **Frontend**: Fetch pathways from API
- [ ] **Frontend**: Connect "Enroll Now" button to POST endpoint (line 166-168)
- [ ] **Frontend**: Connect "Take Assessment" button (line 189)
- [ ] **Frontend**: Add enrollment confirmation modal
- [ ] **Frontend**: Check prerequisites before enrollment
- [ ] **Integration**: Redirect to first module after enrollment

**Impact**: Users see courses but cannot enroll in them!

---

### 3. ACHIEVEMENTS SYSTEM - GAMIFICATION BROKEN
**Files**: `components/achievements/AchievementsDisplay.tsx`
**Status**: 🔴 HIGH - Motivation/engagement system non-functional

- [ ] **Database**: Populate `achievements` table with actual badges
  - [ ] Define criteria for each achievement
  - [ ] Set XP rewards, rarity levels, icons
- [ ] **Backend API**: Implement `GET /users/me/achievements`
  - [ ] Return all achievements with earned status
  - [ ] Include progress percentage for unearned badges
  - [ ] Include earned_date for completed badges
  - [ ] Calculate stats (earned_count, global_rank, total_xp)
- [ ] **Backend Service**: Create achievement awarding logic
  - [ ] Check criteria after each user action
  - [ ] Award badge when criteria met
  - [ ] Send notification on unlock
  - [ ] Record activity
- [ ] **Frontend**: Remove hardcoded badges array (lines 5-86)
- [ ] **Frontend**: Remove hardcoded leaderboard (lines 88-97)
- [ ] **Frontend**: Fetch achievements from API
- [ ] **Frontend**: Add achievement unlock animation
- [ ] **Frontend**: Add badge sharing feature

**Impact**: Users don't get rewarded for accomplishments!

---

### 4. LEADERBOARD SYSTEM - COMPETITION BROKEN
**File**: `components/leaderboard/LeaderboardTable.tsx`
**Status**: 🔴 HIGH - Social/competitive features non-functional

- [ ] **Backend API**: Implement `GET /leaderboard`
  - [ ] Support timeframe filtering (all, month, week)
  - [ ] Support pathway filtering
  - [ ] Add pagination (limit, offset)
  - [ ] Highlight current user (is_current_user: true)
  - [ ] Include rank trend (up/down/same)
  - [ ] Add country codes for users
- [ ] **Backend**: Complete leaderboard_repository.go TODOs
  - [ ] Set country code (line 85)
  - [ ] Calculate rank trends (line 85)
  - [ ] Mark current user (line 85)
  - [ ] Compute user rank and change (line 125)
- [ ] **Backend Service**: Complete leaderboard_service.go TODOs
  - [ ] Set trend and is_current_user (line 84)
- [ ] **Frontend**: Remove hardcoded leaderboardData (lines 5-16)
- [ ] **Frontend**: Remove hardcoded timeframes (lines 18-24)
- [ ] **Frontend**: Remove hardcoded stats (lines 35, 47, 63, 76)
- [ ] **Frontend**: Fetch leaderboard from API
- [ ] **Frontend**: Implement timeframe filters with API calls
- [ ] **Frontend**: Add pagination controls
- [ ] **Frontend**: Add "View Profile" modal

**Impact**: Users can't compete or see rankings!

---

## 🔴 HIGH PRIORITY (Core Functionality)

### 5. PRACTICE ARENA - ENTIRE PRACTICE SYSTEM MOCKED
**File**: `components/practice/PracticeArena.tsx`
**Status**: 🔴 HIGH - Practice features don't work

- [ ] **Backend API**: Implement `GET /practice/challenges`
  - [ ] Return challenge types with metadata
- [ ] **Backend API**: Implement `GET /practice/areas`
  - [ ] Return practice areas by language/topic
  - [ ] Include exercise counts and completion stats
- [ ] **Backend API**: Implement `GET /users/me/practice/stats`
  - [ ] Return practice_score, duels_won, win_rate, avg_time
  - [ ] Return live_opponents count
- [ ] **Backend API**: Implement `GET /users/me/matches?limit=5`
  - [ ] Return recent match history
- [ ] **Backend API**: Implement `POST /practice/challenges/:type/start`
  - [ ] Return exercise ID for challenge
  - [ ] Start timer/tracking
- [ ] **Backend API**: Implement matchmaking system for duels
  - [ ] WebSocket or polling for live matches
- [ ] **Frontend**: Remove hardcoded challengeTypes (lines 6-43)
- [ ] **Frontend**: Remove hardcoded practiceAreas (lines 45-52)
- [ ] **Frontend**: Remove hardcoded recentMatches (lines 54-59)
- [ ] **Frontend**: Remove hardcoded stats (lines 78, 94, 109, 124)
- [ ] **Frontend**: Fetch practice data from API
- [ ] **Frontend**: Implement challenge start logic
- [ ] **Frontend**: Implement real-time duel system

**Impact**: Practice mode doesn't work at all!

---

### 6. IMPLEMENT USER ACTIVITY TRACKING SYSTEM
- [x] **Database**: Create `user_activities` table
- [x] **Model**: Create `Activity` model in `internal/models/activity.go`
- [x] **Repository**: Create `activity_repository.go` with CRUD operations
- [x] **Service**: Create `activity_service.go` with business logic
- [x] **Handler**: Update `user_handler.go` `GetActivities()` to return real data
- [-] **Triggers**: Add activity recording for:
  - [x] Exercise submissions
  - [ ] Module completions  
  - [ ] Pathway enrollments
  - [ ] Achievement unlocks
- [x] **Frontend**: `RecentActivity.tsx` already fetches from API ✓

---

### 7. COMPLETE PROGRESS SERVICE IMPLEMENTATION
- [x] **Fix `RecordSubmissionActivity()`** in `progress_service.go`:
  - [x] Calculate and record time spent on exercise
  - [x] Update `user_module_progress` table
  - [x] Update `user_daily_activity` table
  - [ ] Check and award milestones (TODO at line 119-120)
- [ ] **Add progress aggregation methods**:
  - [ ] Daily study time calculation
  - [ ] Weekly progress summaries
  - [ ] Streak tracking logic
- [x] **Integration**: Connect progress service to submission handler

---

### 8. FIX STUDY TIME TRACKING IN USER STATS
- [x] **Database**: Ensure `user_daily_activity` table has `study_time_minutes` field
- [x] **Query**: Add aggregation query in `progress_repository.go` to sum study time
- [x] **Service**: Update `GetUserProgress()` to include study time calculation
- [x] **Handler**: Update `GetStats()` in `user_handler.go` to use real study time
- [x] **Frontend**: `QuickStats.tsx` already fetches from API ✓

---

### 9. COMPLETE NOTIFICATION SYSTEM
**File**: `wizardcore-backend/internal/handlers/notification_handler.go`
**Status**: 🔴 All notification handlers are stubs!

- [ ] **Backend**: Fix GetNotifications() (currently returns empty array at line 51)
  - [ ] Convert Supabase user ID to internal user ID (TODO at line 50)
  - [ ] Query notifications table
  - [ ] Return notifications with unread count
- [ ] **Backend**: Fix MarkAsRead() (currently just logs at line 72)
  - [ ] Convert Supabase user ID to internal user ID (TODO at line 68)
  - [ ] Update notification read status
- [ ] **Backend**: Fix DeleteNotification() (currently just logs at line 93)
  - [ ] Convert Supabase user ID to internal user ID (TODO at line 91)
  - [ ] Delete notification from database
- [ ] **Backend**: Create notification service
  - [ ] Create notifications on events
  - [ ] Support multiple notification types
- [ ] **Frontend**: Update Header.tsx notification bell
  - [ ] Fetch notifications from API
  - [ ] Show unread count badge
  - [ ] Display notification panel on click
  - [ ] Mark as read on click

**Impact**: Users don't receive any notifications!

---

### 10. IMPLEMENT NAVIGATION BADGE COUNTS
**File**: `wizardcore-backend/internal/handlers/user_handler.go`
**Status**: 🟡 MEDIUM - Nav badges show no counts

- [ ] **Backend**: Implement GetNavCounts() (line 301 - currently returns null)
  - [ ] Count unread notifications
  - [ ] Count new achievements
  - [ ] Count practice challenges available
  - [ ] Return as JSON object
- [ ] **Frontend**: Update Sidebar.tsx to display badge counts
- [ ] **Frontend**: Update Header.tsx to show notification count

---

### 11. IMPLEMENT SEARCH FUNCTIONALITY
**File**: `components/dashboard/Header.tsx`
**Status**: 🟡 MEDIUM - Search bar does nothing

- [ ] **Backend API**: Implement `GET /search?q={query}&type={type}`
  - [ ] Search courses/pathways
  - [ ] Search modules
  - [ ] Search exercises
  - [ ] Search help articles
  - [ ] Return results with highlights
- [ ] **Frontend**: Connect search input to API (lines 39-43)
- [ ] **Frontend**: Add search results dropdown
- [ ] **Frontend**: Add keyboard navigation (arrow keys, enter)
- [ ] **Frontend**: Add recent searches
- [ ] **UX**: Debounce search input (300ms)

---

## 🟡 MEDIUM PRIORITY (User Experience)

### 12. PROGRESS TRACKER PAGE - DUPLICATE/CONFLICTING DATA
**File**: `components/progress/ProgressTracker.tsx`
**Status**: 🟡 MEDIUM - May be duplicate of existing progress API

- [ ] **Assess**: Determine if this duplicates PathwayProgressList.tsx functionality
- [ ] **Backend**: If needed, ensure `/users/me/progress` returns all data
  - [ ] Pathway progress with modules completed
  - [ ] Overall stats (total XP, streak, modules)
- [ ] **Backend**: If needed, implement `/users/me/milestones`
  - [ ] Return completed milestones with dates and XP
- [ ] **Backend**: If needed, implement `/users/me/activity/weekly-hours`
  - [ ] Return daily hour counts for the week
- [ ] **Frontend**: Remove hardcoded pathways (lines 5-66)
- [ ] **Frontend**: Remove hardcoded milestones (lines 68-74)
- [ ] **Frontend**: Remove hardcoded weekly hours (lines 235-242)
- [ ] **Frontend**: Fetch progress data from API
- [ ] **Decision**: Merge with existing progress components or keep separate?

---

### 13. IMPLEMENT USER PREFERENCES SYSTEM
- [x] **Database**: `user_preferences` table already exists ✓
- [x] **Model**: `UserPreferences` model exists ✓
- [x] **Repository**: `preferences_repository.go` created ✓
- [x] **Handler**: `GetPreferences()` and `UpdatePreferences()` implemented ✓
- [ ] **Frontend**: Create preferences UI in settings panel
  - [ ] Theme selector
  - [ ] Notification preferences
  - [ ] Learning preferences
  - [ ] Email frequency
- [ ] **Integration**: Apply theme preferences globally

---

### 14. ADD DEADLINE MANAGEMENT SYSTEM
**File**: `wizardcore-backend/internal/handlers/pathway_handler.go`
**Status**: 🟡 MEDIUM - Deadlines stubbed

- [ ] **Database**: `user_deadlines` table already exists ✓
- [ ] **Model**: Complete `Deadline` model
- [ ] **Repository**: Create `deadline_repository.go`
- [ ] **Handler**: Implement `GetDeadlines()` in pathway_handler.go (TODO at line 125)
  - [ ] Query deadlines for user
  - [ ] Filter by upcoming/overdue
  - [ ] Sort by due date
- [ ] **Service**: Add deadline reminder logic
  - [ ] Send notifications before due date
  - [ ] Mark overdue deadlines
- [ ] **Frontend**: Update dashboard page.tsx deadlines section (lines 38-51)
  - [ ] Remove hardcoded deadlines
  - [ ] Fetch from `/users/me/deadlines`
  - [ ] Show countdown timers
  - [ ] Add deadline creation UI

---

### 15. COMPLETE CONTENT CREATOR EXAMPLES SYSTEM
**File**: `app/creator/exercises/new/NewExerciseContent.tsx`
**Status**: 🟡 LOW - Preview feature missing

- [ ] **Database**: Update `exercises` table `examples` field JSON schema
- [ ] **Validation**: Add JSON schema validation for examples
- [ ] **Frontend**: Implement preview modal (TODO at line 81)
  - [ ] Create modal component
  - [ ] Render exercise preview
  - [ ] Show test cases
  - [ ] Show hints
- [ ] **Frontend**: Update `ExerciseBuilder.tsx` to handle examples properly
- [ ] **API**: Validate examples on exercise create/update
- [ ] **Cleanup**: Remove .bak files (page.tsx.bak with duplicate TODO at line 84)

---

### 16. ADD ADMIN PERMISSION CHECKS
**File**: `wizardcore-backend/internal/handlers/content_creator_handler.go`
**Status**: 🟡 MEDIUM - Security issue!

- [ ] **Backend**: Implement admin check (TODO at line 704)
  - [ ] Check user role before allowing action
  - [ ] Return 403 Forbidden if not admin
- [ ] **Service**: Update content_creator_service.go
  - [ ] Implement role updates (TODO at line 63)
  - [ ] Verify admin role (TODO at line 543)
- [ ] **Testing**: Test role-based access control

---

### 17. IMPLEMENT PATHWAY CARD NAVIGATION
**File**: `components/dashboard/PathwayCard.tsx`
**Status**: 🟡 LOW - Buttons don't work

- [ ] **Frontend**: Add onClick to "Continue" button (line 50-52)
  - [ ] Navigate to `/dashboard/learning?pathway={id}&resume=true`
  - [ ] Track button click analytics
- [ ] **Frontend**: Add onClick to "Details" button (line 53-55)
  - [ ] Navigate to `/dashboard/pathways/{id}`
- [ ] **Frontend**: Pass pathway_id as prop
- [ ] **Integration**: Implement resume logic (fetch last incomplete exercise)

---

## 🔵 LOW PRIORITY (Advanced Features)

### 18. IMPLEMENT FULL RBAC PERMISSION SYSTEM
- [ ] **Data**: Populate `permissions` table with actual permissions
- [ ] **Assignments**: Create default role-permission assignments
- [ ] **Middleware**: Implement permission checking middleware
- [ ] **Handler**: Update RBAC handlers to use repository
- [ ] **Service**: Complete RBAC service TODOs
  - [ ] Add checks for dangerous permissions (TODO at line 246)
  - [ ] Determine which roles grant permissions (TODO at line 503)
- [ ] **Audit**: Implement permission usage logging
- [ ] **UI**: Create admin panel for RBAC management

---

### 19. ADD ANALYTICS AND REPORTING SYSTEM
- [ ] **Database**: Create analytics views/tables
- [ ] **Queries**: Add aggregation queries for:
  - [ ] User engagement metrics
  - [ ] Course completion rates
  - [ ] Popular content analysis
  - [ ] Retention metrics
- [ ] **API**: Create analytics endpoints
- [ ] **Frontend**: Create admin analytics dashboard

---

### 20. COMPLETE SUPABASE AUTH INTEGRATION
- [ ] **Schema**: Run Supabase Auth migrations
- [ ] **Sync**: Implement user synchronization between Supabase and local users
- [ ] **Webhooks**: Set up Supabase webhooks for user events
- [ ] **Testing**: Test full auth flow end-to-end

---

### 21. FIX TEST INFRASTRUCTURE AND MIGRATIONS
- [ ] **Migrations**: Create proper migration files
- [ ] **Test DB**: Fix `testdb.go` to run migrations
- [ ] **Fixtures**: Create test data fixtures
- [ ] **Integration Tests**: Add tests for all handlers
- [ ] **E2E Tests**: Add end-to-end test suite

---

## 🧹 CLEANUP TASKS

### 22. REMOVE UNUSED CODE
- [ ] **Judge0 Service**: Remove unused functions in `lib/judge0/service.ts`
  - [ ] Remove batchSubmit() (lines 75-104) - never called
  - [ ] Remove getSubmission() (lines 109-123) - never called
  - [ ] Remove healthCheck() (lines 128-138) - never called
  - [ ] OR implement these functions if needed
- [ ] **Backup Files**: Delete .bak files
  - [ ] `app/creator/exercises/new/page.tsx.bak`
  - [ ] `app/creator/modules/new/page.tsx.bak`
- [ ] **Console Logs**: Remove debug console.log statements
  - [ ] `app/creator/exercises/new/NewExerciseContent.tsx` line 80 (preview log)
  - [ ] Keep error logging console.error statements

---

### 23. REFACTOR JUDGE0 INTEGRATION
**Status**: 🟡 MEDIUM - Security/architecture issue

- [ ] **Decision**: Move Judge0 calls to backend or keep in frontend?
- [ ] **Backend**: If moving, create `/code/execute` endpoint
  - [ ] Proxy to Judge0
  - [ ] Hide API credentials
  - [ ] Add rate limiting
  - [ ] Cache results
- [ ] **Frontend**: Update components to use backend endpoint
- [ ] **Security**: Never expose Judge0 API key to frontend

---

## ADDITIONAL TASKS IDENTIFIED

### Database Schema Updates:
- [ ] Add missing foreign key constraints
- [ ] Add indexes for performance
- [ ] Create database views for common queries
- [ ] Add database functions for complex operations

### API Enhancements:
- [ ] Add pagination to list endpoints (TODO at user_handler.go line 242)
- [ ] Add filtering and sorting to list endpoints
- [ ] Add rate limiting
- [ ] Add request validation middleware
- [ ] Add comprehensive error handling

### Frontend Improvements:
- [x] Add loading states for async operations (QuickStats, ProgressChart, RecentActivity ✓)
- [ ] Add error boundaries
- [ ] Implement proper form validation
- [ ] Add offline support indicators
- [ ] Improve accessibility (ARIA labels, keyboard navigation)

### Performance Optimizations:
- [ ] Add database query optimization
- [ ] Implement caching layer (Redis)
- [ ] Add CDN for static assets
- [ ] Implement code splitting
- [ ] Add performance monitoring

### Security Enhancements:
- [ ] Add input sanitization
- [ ] Implement CSRF protection
- [ ] Add security headers
- [ ] Implement audit logging
- [ ] Add security scanning to CI/CD

### DevOps & Deployment:
- [ ] Create comprehensive CI/CD pipeline
- [ ] Add health checks
- [ ] Implement monitoring and alerting
- [ ] Create backup and recovery procedures
- [ ] Add environment-specific configurations

---

## PROGRESS TRACKING

### Completed Tasks:
- [x] Created `user_activities` table migration
- [x] Created `Activity` model (`UserActivity`)
- [x] Created `activity_repository.go` with CRUD operations
- [x] Created `activity_service.go` with business logic
- [x] Updated `user_handler.go` `GetActivities()` to return real data
- [x] Fixed `RecordSubmissionActivity()` in `progress_service.go`
- [x] Added study time aggregation query in `progress_repository.go`
- [x] Updated `GetStats()` to use real study time calculation
- [x] Connected progress service to submission handler
- [x] Added activity recording to submission service
- [x] Tested activity recording integration
- [x] Created `preferences_repository.go`
- [x] Updated user service with preferences methods
- [x] Implemented `GetPreferences()` and `UpdatePreferences()` handlers
- [x] Updated router with preferences repository
- [x] Frontend components using API: QuickStats, ProgressChart, RecentActivity, PathwayProgressList

### In Progress:
- [-] Frontend mock data removal (major effort required)

### Blocked/Needs Attention:
- [!] **CRITICAL**: Learning environment is completely non-functional
- [!] **CRITICAL**: Users cannot enroll in pathways
- [!] **HIGH**: Achievements system not working
- [!] **HIGH**: Leaderboard system not working
- [!] **HIGH**: Practice arena not working

---

## IMMEDIATE NEXT STEPS (CRITICAL PATH)

### Week 1: Core Learning System
1. [ ] **Learning Environment API** (3-4 days)
   - [ ] Complete `/exercises/:id` endpoint with all fields
   - [ ] Implement `/submissions` POST with Judge0 integration
   - [ ] Implement `/submissions/:exercise_id/latest` GET
   - [ ] Test Judge0 integration end-to-end

2. [ ] **Learning Environment Frontend** (2-3 days)
   - [ ] Replace all hardcoded data with API calls
   - [ ] Implement real submission flow
   - [ ] Add auto-save functionality
   - [ ] Test complete learning flow

### Week 2: Pathway Enrollment
3. [ ] **Pathway Enrollment Backend** (2 days)
   - [ ] Implement `POST /pathways/:id/enroll`
   - [ ] Update `GET /pathways` with enrollment status
   - [ ] Add activity recording

4. [ ] **Pathway Enrollment Frontend** (1 day)
   - [ ] Connect enrollment buttons
   - [ ] Add confirmation modal
   - [ ] Update UI after enrollment

### Week 3: Gamification Systems
5. [ ] **Achievements System** (3 days)
   - [ ] Populate achievements table
   - [ ] Implement `/users/me/achievements` endpoint
   - [ ] Create achievement awarding logic
   - [ ] Update frontend

6. [ ] **Leaderboard System** (2 days)
   - [ ] Implement `/leaderboard` endpoint
   - [ ] Complete repository TODOs
   - [ ] Update frontend

### Week 4: Practice & Polish
7. [ ] **Practice Arena** (3 days)
   - [ ] Implement practice endpoints
   - [ ] Update frontend

8. [ ] **Notifications & Navigation** (2 days)
   - [ ] Complete notification handlers
   - [ ] Implement GetNavCounts
   - [ ] Update frontend

---

## ESTIMATED TIMELINE

### Critical Priority (Must Have):
- **Learning Environment**: 1 week
- **Pathway Enrollment**: 3-4 days
- **Achievements**: 3 days
- **Leaderboard**: 2 days
- **Total**: ~3 weeks

### High Priority (Should Have):
- **Practice Arena**: 3 days
- **Notifications**: 2 days
- **Search**: 2 days
- **Total**: ~1 week

### Medium Priority (Nice to Have):
- **Progress Tracker**: 2 days
- **Deadlines**: 2 days
- **Preferences UI**: 1 day
- **Total**: ~1 week

### Low Priority (Future):
- **RBAC**: 1 week
- **Analytics**: 1 week
- **Testing**: 1 week
- **Total**: ~3 weeks

**Grand Total**: 8-9 weeks for complete implementation

---

## KEY FILES TO MODIFY

### Backend (Go):
- `wizardcore-backend/internal/handlers/exercise_handler.go` - Learning environment
- `wizardcore-backend/internal/handlers/submission_handler.go` - Code submission
- `wizardcore-backend/internal/handlers/pathway_handler.go` - Enrollment, deadlines
- `wizardcore-backend/internal/handlers/user_handler.go` - Nav counts
- `wizardcore-backend/internal/handlers/notification_handler.go` - All 3 functions
- `wizardcore-backend/internal/handlers/achievement_handler.go` - Achievements
- `wizardcore-backend/internal/handlers/leaderboard_handler.go` - Leaderboard
- `wizardcore-backend/internal/handlers/practice_handler.go` - Practice arena
- `wizardcore-backend/internal/services/progress_service.go` - Milestones
- `wizardcore-backend/internal/services/achievement_service.go` - Awarding logic
- `wizardcore-backend/internal/repositories/leaderboard_repository.go` - Trends, highlighting

### Frontend (TypeScript/React):
- `components/learning/LearningEnvironment.tsx` - Core learning
- `components/pathways/PathwaySelection.tsx` - Enrollment
- `components/achievements/AchievementsDisplay.tsx` - Achievements
- `components/leaderboard/LeaderboardTable.tsx` - Leaderboard
- `components/practice/PracticeArena.tsx` - Practice
- `components/progress/ProgressTracker.tsx` - Progress
- `components/dashboard/Header.tsx` - Search, notifications
- `components/dashboard/PathwayCard.tsx` - Navigation
- `app/dashboard/page.tsx` - Deadlines

---

## NOTES

### Current Status Assessment (Updated):
- **Backend**: ~70% complete
  - ✅ Core CRUD works
  - ✅ Some endpoints functional (stats, activities, progress)
  - ❌ Many endpoints are stubs (notifications, nav counts, deadlines)
  - ❌ Missing critical endpoints (achievements, leaderboard, practice)
  
- **Frontend**: ~40% functional
  - ✅ 4 components use real API (QuickStats, ProgressChart, RecentActivity, PathwayProgressList)
  - ❌ 8+ major components use 100% mock data
  - ❌ Core learning environment is fake
  - ❌ Enrollment doesn't work
  
- **Database**: ~90% ready
  - ✅ Schema is comprehensive
  - ✅ Tables exist for all features
  - ❌ Many tables empty or underutilized
  - ❌ Achievements table needs data

### Risk Assessment:
- 🔴 **CRITICAL**: Users believe the app works but it's mostly a facade
- 🔴 **CRITICAL**: Cannot complete exercises (core product broken)
- 🔴 **CRITICAL**: Cannot enroll in courses (conversion funnel broken)
- 🟡 **HIGH**: No gamification (engagement will suffer)
- 🟡 **HIGH**: No competition (viral growth limited)

### Success Criteria:
- [ ] User can sign up, enroll in pathway, complete exercise, earn achievement
- [ ] User sees real progress, real stats, real leaderboard position
- [ ] User receives notifications for achievements and deadlines
- [ ] User can practice and compete in duels
- [ ] All mock data removed from frontend

---

## HOW TO UPDATE THIS LIST

1. When starting a task, change `[ ]` to `[-]`
2. When completing a task, change `[-]` or `[ ]` to `[x]`
3. When blocked, change to `[!]` and add note in "Blocked/Needs Attention"
4. Update the "Last Updated" date at the top
5. Move completed tasks to "Completed Tasks" section
6. Add new tasks as identified from code audits

### Example:
```
- [-] Task in progress
- [x] Task completed
- [!] Task blocked (needs API design clarification)
```

---

*This document reflects the complete state of mock data, TODOs, and stubs found in comprehensive codebase audit conducted December 27, 2025.*
