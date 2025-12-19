# Content Creator Implementation Plan

## Overview
This document outlines the implementation plan for adding content creator functionality to the WizardCore cybersecurity learning platform. The system will allow authenticated users with appropriate permissions to create and manage courses, modules, and exercises.

## Current Architecture Analysis

### Existing Database Schema
1. **Users Table**: Basic user information, no role system
2. **Pathways Table**: Courses with basic metadata
3. **Modules Table**: Learning modules within pathways
4. **Exercises Table**: Code exercises with Judge0 integration
5. **Test Cases Table**: Test cases for exercises
6. **Submissions Table**: Student exercise submissions

### Existing Backend Services
- User authentication/authorization (Supabase)
- Exercise evaluation (Judge0 integration)
- Pathway/module management
- Progress tracking

## Implementation Requirements

### 1. User Role System
- Add `role` field to users table (student, content_creator, admin)
- Content creator profiles with verification system
- Role-based access control

### 2. Content Management System
- Create/update/delete pathways (courses)
- Create/update/delete modules within pathways
- Create/update/delete exercises within modules
- Content versioning and history
- Draft/published/archived status management

### 3. Exercise Creation with Judge0 Integration
- Code exercise creation interface
- Test case management (visible/hidden test cases)
- Starter code and solution templates
- Language selection (Python, JavaScript, C, etc.)
- Difficulty levels and point values

### 4. Review and Approval Workflow
- Content submission for admin review
- Review notes and feedback system
- Revision tracking
- Publishing controls

### 5. Analytics and Reporting
- Content creator dashboard
- Student engagement metrics
- Exercise completion statistics
- Revenue sharing (if applicable)

## Technical Implementation Plan

### Phase 1: Database Schema Updates
1. **Migration 012_content_creators.up.sql**:
   - Add `role` column to users table
   - Create `content_creator_profiles` table
   - Add creator metadata to pathways, modules, exercises
   - Add content status fields
   - Create content review and version history tables

### Phase 2: Backend Services
1. **ContentCreatorService**:
   - Profile management
   - Content creation/editing
   - Review workflow management
   - Statistics and analytics

2. **ContentCreatorHandler**:
   - REST API endpoints for all content operations
   - Authentication and authorization middleware
   - Input validation and error handling

### Phase 3: Frontend Components
1. **Content Creator Dashboard**:
   - Overview of created content
   - Statistics and metrics
   - Quick actions

2. **Course Creation Interface**:
   - Pathway creation form
   - Module management
   - Exercise builder with code editor

3. **Exercise Builder**:
   - Integrated code editor (Monaco)
   - Test case management
   - Judge0 language selection
   - Preview functionality

### Phase 4: Integration
1. **API Integration**:
   - Connect frontend to backend endpoints
   - Real-time validation
   - Error handling and user feedback

2. **Judge0 Integration**:
   - Exercise validation during creation
   - Test case verification
   - Performance optimization

## Database Schema Changes

### New Tables
```sql
-- Content creator profiles
CREATE TABLE content_creator_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id),
    bio TEXT,
    specialization TEXT[],
    website VARCHAR(255),
    github_url VARCHAR(255),
    linkedin_url VARCHAR(255),
    twitter_url VARCHAR(255),
    is_verified BOOLEAN DEFAULT false,
    verification_date TIMESTAMP,
    total_content_created INTEGER DEFAULT 0,
    total_students INTEGER DEFAULT 0,
    average_rating DECIMAL(3,2) DEFAULT 0.0
);

-- Content reviews
CREATE TABLE content_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    reviewer_id UUID REFERENCES users(id),
    status VARCHAR(50) CHECK (status IN ('pending', 'approved', 'rejected', 'needs_revision')),
    review_notes TEXT,
    revision_notes TEXT,
    reviewed_at TIMESTAMP
);

-- Content version history
CREATE TABLE content_version_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    version INTEGER NOT NULL,
    data JSONB NOT NULL,
    created_by UUID REFERENCES users(id),
    change_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Modified Tables
```sql
-- Add to users table
ALTER TABLE users ADD COLUMN role VARCHAR(50) DEFAULT 'student' 
CHECK (role IN ('student', 'content_creator', 'admin'));

-- Add to pathways table
ALTER TABLE pathways ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE pathways ADD COLUMN status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE pathways ADD COLUMN version INTEGER DEFAULT 1;
ALTER TABLE pathways ADD COLUMN published_at TIMESTAMP;
ALTER TABLE pathways ADD COLUMN review_notes TEXT;

-- Add to modules table
ALTER TABLE modules ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE modules ADD COLUMN status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE modules ADD COLUMN version INTEGER DEFAULT 1;
ALTER TABLE modules ADD COLUMN published_at TIMESTAMP;

-- Add to exercises table
ALTER TABLE exercises ADD COLUMN created_by UUID REFERENCES users(id);
ALTER TABLE exercises ADD COLUMN status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE exercises ADD COLUMN version INTEGER DEFAULT 1;
ALTER TABLE exercises ADD COLUMN published_at TIMESTAMP;
ALTER TABLE exercises ADD COLUMN requires_approval BOOLEAN DEFAULT false;
```

## API Endpoints

### Content Creator Profile
- `GET /api/content-creator/profile` - Get creator profile
- `PUT /api/content-creator/profile` - Update creator profile
- `GET /api/content-creator/stats` - Get creator statistics

### Pathway Management
- `POST /api/content-creator/pathways` - Create new pathway
- `PUT /api/content-creator/pathways/:id` - Update pathway
- `GET /api/content-creator/pathways` - List creator's pathways

### Module Management
- `POST /api/content-creator/pathways/:pathwayId/modules` - Create module
- `PUT /api/content-creator/modules/:id` - Update module
- `GET /api/content-creator/modules` - List creator's modules

### Exercise Management
- `POST /api/content-creator/modules/:moduleId/exercises` - Create exercise
- `PUT /api/content-creator/exercises/:id` - Update exercise
- `GET /api/content-creator/exercises` - List creator's exercises

### Review Workflow
- `POST /api/content-creator/reviews` - Submit content for review
- `GET /api/content-creator/reviews` - Get review status

## Frontend Components Structure

### Pages
1. `/creator/dashboard` - Content creator dashboard
2. `/creator/pathways` - Pathway management
3. `/creator/pathways/new` - Create new pathway
4. `/creator/pathways/:id/edit` - Edit pathway
5. `/creator/modules/new` - Create new module
6. `/creator/exercises/new` - Create new exercise

### Components
1. `CreatorDashboard.tsx` - Main dashboard component
2. `PathwayForm.tsx` - Pathway creation/editing form
3. `ModuleForm.tsx` - Module creation/editing form
4. `ExerciseBuilder.tsx` - Exercise creation with code editor
5. `TestCaseManager.tsx` - Test case management
6. `ContentStatusBadge.tsx` - Status indicator component
7. `ReviewWorkflow.tsx` - Review submission interface

## Judge0 Integration for Exercise Creation

### Features
1. **Language Support**: Python, JavaScript, C, C++, Java, Go, Rust
2. **Code Validation**: Real-time syntax checking
3. **Test Execution**: Run test cases during creation
4. **Performance Testing**: Time/memory limit validation

### Implementation
```typescript
interface ExerciseCreationData {
  title: string;
  description: string;
  difficulty: 'beginner' | 'intermediate' | 'advanced';
  language_id: number; // Judge0 language ID
  starter_code: string;
  solution_code: string;
  test_cases: TestCase[];
  time_limit_minutes: number;
  points: number;
}

interface TestCase {
  input: string;
  expected_output: string;
  is_hidden: boolean;
  points: number;
}
```

## Security Considerations

1. **Authentication**: Supabase JWT tokens
2. **Authorization**: Role-based access control
3. **Input Validation**: Sanitize all user inputs
4. **Rate Limiting**: Prevent abuse of creation endpoints
5. **Content Moderation**: Review workflow for sensitive content

## Testing Strategy

### Unit Tests
- Service layer logic
- Input validation
- Business rules

### Integration Tests
- API endpoints
- Database operations
- Judge0 integration

### End-to-End Tests
- Complete content creation workflow
- User interface interactions
- Cross-browser compatibility

## Deployment Plan

### Phase 1: Development
- Implement database migrations
- Create backend services
- Develop core frontend components

### Phase 2: Staging
- Deploy to staging environment
- User acceptance testing
- Performance testing

### Phase 3: Production
- Gradual rollout to beta users
- Monitor system performance
- Collect user feedback

## Success Metrics

1. **Adoption Rate**: Percentage of users becoming content creators
2. **Content Quality**: Average exercise completion rate
3. **Creator Satisfaction**: Creator feedback scores
4. **System Performance**: API response times under load
5. **Revenue Impact**: If monetization is implemented

## Timeline Estimate

- **Week 1-2**: Database schema and backend services
- **Week 3-4**: Frontend components and UI/UX
- **Week 5**: Integration and testing
- **Week 6**: Staging deployment and UAT
- **Week 7**: Production rollout

## Risks and Mitigations

1. **Technical Complexity**: Judge0 integration may be complex
   - Mitigation: Start with basic Python exercises, expand gradually

2. **Content Quality**: Low-quality user-generated content
   - Mitigation: Implement review workflow and quality guidelines

3. **Performance**: Large number of exercises may impact performance
   - Mitigation: Implement pagination, caching, and CDN for static content

4. **Security**: Malicious code in user-created exercises
   - Mitigation: Sandboxed execution, input validation, content moderation

## Next Steps

1. Review and approve this implementation plan
2. Set up development environment
3. Begin with Phase 1: Database schema updates
4. Schedule regular progress reviews
5. Establish quality assurance process

---

*Last Updated: December 19, 2025*  
*Version: 1.0*  
*Author: WizardCore Development Team*