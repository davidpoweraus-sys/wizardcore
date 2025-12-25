# Content Creator System - Quick Start Guide

## 🚀 Backend is Complete!

The entire backend infrastructure for the content creator system has been implemented and is ready to use.

## ✅ What's Working

### Database Schema ✓
- All 5 new tables created
- Existing tables modified with creator fields
- Indexes optimized
- Triggers for updated_at columns

### Backend Services ✓
- Profile management
- Pathway/Module/Exercise CRUD
- Review workflow
- Statistics and analytics
- Role-based authorization

### API Endpoints ✓
- 20+ endpoints ready to use
- Full CRUD for all content types
- Protected by role middleware

## 🧪 Testing the Backend

### Step 1: Start the Backend

```bash
cd wizardcore-backend
go run cmd/api/main.go
```

The migration will automatically run on startup.

### Step 2: Create a Test User with Creator Role

First, you'll need to update a user's role to `content_creator`:

```sql
-- Connect to your database
psql $DATABASE_URL

-- Update an existing user to be a content creator
UPDATE users 
SET role = 'content_creator' 
WHERE email = 'your-test-email@example.com';

-- Verify the change
SELECT id, email, role FROM users WHERE email = 'your-test-email@example.com';
```

### Step 3: Test API Endpoints

Get your JWT token from Supabase (login through your frontend or use Supabase dashboard).

```bash
# Set your token
export TOKEN="your-jwt-token-here"

# Test 1: Create Creator Profile
curl -X POST http://localhost:8080/api/v1/content-creator/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "bio": "Expert in exploit development and reverse engineering",
    "specialization": ["exploit-dev", "reverse-engineering", "binary-analysis"],
    "website": "https://example.com",
    "github_url": "https://github.com/username"
  }'

# Test 2: Get Creator Profile
curl http://localhost:8080/api/v1/content-creator/profile \
  -H "Authorization: Bearer $TOKEN"

# Test 3: Create a Pathway
curl -X POST http://localhost:8080/api/v1/content-creator/pathways \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Advanced Binary Exploitation",
    "subtitle": "Master modern exploitation techniques",
    "description": "Learn to exploit buffer overflows, heap corruption, and bypass modern protections",
    "level": "Advanced",
    "duration_weeks": 12,
    "color_gradient": "from-red-500 to-purple-600",
    "icon": "🔓",
    "sort_order": 1,
    "status": "draft"
  }'

# Test 4: Get Creator's Pathways
curl http://localhost:8080/api/v1/content-creator/pathways \
  -H "Authorization: Bearer $TOKEN"

# Test 5: Create a Module (use pathway_id from step 3)
curl -X POST http://localhost:8080/api/v1/content-creator/modules \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "pathway_id": "PATHWAY-UUID-HERE",
    "title": "Stack Buffer Overflows",
    "description": "Understanding and exploiting stack-based buffer overflows",
    "sort_order": 1,
    "estimated_hours": 8,
    "xp_reward": 500,
    "status": "draft"
  }'

# Test 6: Create an Exercise (use module_id from step 5)
curl -X POST http://localhost:8080/api/v1/content-creator/exercises \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "module_id": "MODULE-UUID-HERE",
    "title": "Basic Stack Overflow",
    "difficulty": "BEGINNER",
    "points": 100,
    "time_limit_minutes": 30,
    "sort_order": 1,
    "objectives": [
      "Understand stack memory layout",
      "Identify buffer overflow vulnerability",
      "Overwrite return address"
    ],
    "content": "# Stack Buffer Overflow\n\nIn this exercise, you will exploit a simple stack buffer overflow...",
    "description": "Write a Python script to exploit a vulnerable C program",
    "constraints": [
      "Must use Python 3",
      "No external libraries beyond standard library"
    ],
    "hints": [
      "Look for strcpy() usage",
      "Calculate the offset to the return address"
    ],
    "starter_code": "#!/usr/bin/env python3\n# Your exploit code here\n",
    "solution_code": "#!/usr/bin/env python3\nimport struct\npayload = b\"A\" * 64 + struct.pack(\"<Q\", 0xdeadbeef)\nprint(payload.decode(\"latin-1\"))",
    "language_id": 71,
    "tags": ["buffer-overflow", "stack", "exploitation"],
    "status": "draft",
    "test_cases": [
      {
        "input": "",
        "expected_output": "Exploit successful!",
        "is_hidden": false,
        "points": 100,
        "sort_order": 1
      }
    ]
  }'

# Test 7: Get Creator Statistics
curl http://localhost:8080/api/v1/content-creator/stats \
  -H "Authorization: Bearer $TOKEN"

# Test 8: Submit Content for Review
curl -X POST http://localhost:8080/api/v1/content-creator/reviews \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content_type": "pathway",
    "content_id": "PATHWAY-UUID-HERE",
    "revision_notes": "First submission, ready for review"
  }'
```

## 📊 Database Verification

```sql
-- Check if migration ran
SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 5;

-- View creator profiles
SELECT * FROM content_creator_profiles;

-- View content by creator
SELECT p.title, p.status, u.email as creator
FROM pathways p
JOIN users u ON p.created_by = u.id
WHERE u.role = 'content_creator';

-- View pending reviews
SELECT cr.*, 
       CASE 
         WHEN cr.content_type = 'pathway' THEN (SELECT title FROM pathways WHERE id = cr.content_id)
         WHEN cr.content_type = 'module' THEN (SELECT title FROM modules WHERE id = cr.content_id)
         WHEN cr.content_type = 'exercise' THEN (SELECT title FROM exercises WHERE id = cr.content_id)
       END as content_title
FROM content_reviews cr
WHERE cr.status = 'pending';
```

## 🎯 Frontend Development Starting Points

### 1. Create Content Creator Dashboard Page

```typescript
// app/creator/dashboard/page.tsx
import { api } from '@/lib/api'

export default async function CreatorDashboard() {
  const stats = await api.get('/content-creator/stats')
  
  return (
    <div>
      <h1>Creator Dashboard</h1>
      <div className="grid grid-cols-4 gap-4">
        <StatCard title="Total Pathways" value={stats.total_pathways} />
        <StatCard title="Total Students" value={stats.total_students} />
        <StatCard title="Avg Rating" value={stats.average_rating} />
        <StatCard title="Pending Reviews" value={stats.pending_reviews} />
      </div>
    </div>
  )
}
```

### 2. Create Pathway Form Component

```typescript
// components/creator/PathwayForm.tsx
'use client'

import { useState } from 'react'
import { api } from '@/lib/api'

export function PathwayForm() {
  const [formData, setFormData] = useState({
    title: '',
    subtitle: '',
    description: '',
    level: 'Beginner',
    duration_weeks: 4,
    status: 'draft'
  })

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    try {
      const pathway = await api.post('/content-creator/pathways', formData)
      console.log('Created pathway:', pathway)
      // Redirect or show success message
    } catch (error) {
      console.error('Failed to create pathway:', error)
    }
  }

  return <form onSubmit={handleSubmit}>{/* Form fields */}</form>
}
```

### 3. Create Exercise Builder with Monaco

```typescript
// components/creator/ExerciseBuilder.tsx
'use client'

import Editor from '@monaco-editor/react'
import { useState } from 'react'

export function ExerciseBuilder() {
  const [code, setCode] = useState('# Starter code here')
  const [testCases, setTestCases] = useState([])

  return (
    <div>
      <Editor
        height="400px"
        language="python"
        value={code}
        onChange={(value) => setCode(value || '')}
        theme="vs-dark"
      />
      {/* Test case manager */}
    </div>
  )
}
```

## 🔧 Environment Variables

Make sure these are set in your `.env.local`:

```env
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
NEXT_PUBLIC_SUPABASE_URL=your-supabase-url
NEXT_PUBLIC_SUPABASE_ANON_KEY=your-supabase-anon-key
```

## 📚 API Documentation

Full API documentation is available in:
- `Documentation/content-creator-backend-implementation-summary.md`

## 🐛 Troubleshooting

### Migration doesn't run
```bash
# Manually run migration
cd wizardcore-backend
go run cmd/api/main.go
# Check logs for migration status
```

### 403 Forbidden on creator routes
```sql
-- Make sure user has correct role
UPDATE users SET role = 'content_creator' WHERE email = 'your-email@example.com';
```

### Can't access content created by another user
This is expected! Ownership validation ensures creators can only modify their own content.

## 🎉 Ready to Build!

You now have a fully functional backend for:
- ✅ Content creator profiles
- ✅ Pathway, module, and exercise management
- ✅ Review workflow
- ✅ Analytics and statistics
- ✅ Role-based authorization

Start building the frontend components and connect them to these endpoints!

---

**Next Steps**:
1. Install Monaco Editor: `npm install @monaco-editor/react`
2. Create creator dashboard page
3. Build pathway/module/exercise forms
4. Implement exercise builder with code editor
5. Add review workflow UI

**Questions?** Check the implementation summary or examine the backend code directly!
