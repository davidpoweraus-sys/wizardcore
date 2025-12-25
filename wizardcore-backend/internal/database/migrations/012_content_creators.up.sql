-- Add role column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'student' 
CHECK (role IN ('student', 'content_creator', 'admin'));

-- Create index on role for faster queries
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

-- Create content creator profiles table
CREATE TABLE IF NOT EXISTS content_creator_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE REFERENCES users(id) ON DELETE CASCADE,
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
    average_rating DECIMAL(3,2) DEFAULT 0.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_creator_profiles_user_id ON content_creator_profiles(user_id);
CREATE INDEX IF NOT EXISTS idx_creator_profiles_verified ON content_creator_profiles(is_verified);

-- Add creator fields to pathways table
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'under_review'));
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS review_notes TEXT;

CREATE INDEX IF NOT EXISTS idx_pathways_created_by ON pathways(created_by);
CREATE INDEX IF NOT EXISTS idx_pathways_status ON pathways(status);

-- Add creator fields to modules table
ALTER TABLE modules ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE modules ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'under_review'));
ALTER TABLE modules ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE modules ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_modules_created_by ON modules(created_by);
CREATE INDEX IF NOT EXISTS idx_modules_status ON modules(status);

-- Add creator fields to exercises table
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' 
CHECK (status IN ('draft', 'published', 'archived', 'under_review'));
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN DEFAULT false;

CREATE INDEX IF NOT EXISTS idx_exercises_created_by ON exercises(created_by);
CREATE INDEX IF NOT EXISTS idx_exercises_status ON exercises(status);

-- Create content reviews table
CREATE TABLE IF NOT EXISTS content_reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    reviewer_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL CHECK (status IN ('pending', 'approved', 'rejected', 'needs_revision')),
    review_notes TEXT,
    revision_notes TEXT,
    reviewed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reviews_content_type_id ON content_reviews(content_type, content_id);
CREATE INDEX IF NOT EXISTS idx_reviews_reviewer ON content_reviews(reviewer_id);
CREATE INDEX IF NOT EXISTS idx_reviews_status ON content_reviews(status);

-- Create content version history table
CREATE TABLE IF NOT EXISTS content_version_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    version INTEGER NOT NULL,
    data JSONB NOT NULL,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    change_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_version_history_content ON content_version_history(content_type, content_id, version DESC);
CREATE INDEX IF NOT EXISTS idx_version_history_created_by ON content_version_history(created_by);

-- Create content analytics table
CREATE TABLE IF NOT EXISTS content_analytics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    creator_id UUID REFERENCES users(id) ON DELETE CASCADE,
    views INTEGER DEFAULT 0,
    enrollments INTEGER DEFAULT 0,
    completions INTEGER DEFAULT 0,
    average_rating DECIMAL(3,2) DEFAULT 0.0,
    total_ratings INTEGER DEFAULT 0,
    date DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_type, content_id, date)
);

CREATE INDEX IF NOT EXISTS idx_analytics_creator ON content_analytics(creator_id, date DESC);
CREATE INDEX IF NOT EXISTS idx_analytics_content ON content_analytics(content_type, content_id, date DESC);

-- Create content ratings table
CREATE TABLE IF NOT EXISTS content_ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_type, content_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_ratings_content ON content_ratings(content_type, content_id);
CREATE INDEX IF NOT EXISTS idx_ratings_user ON content_ratings(user_id);

-- Update updated_at timestamp function for new tables
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
CREATE TRIGGER update_content_creator_profiles_updated_at BEFORE UPDATE ON content_creator_profiles 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_reviews_updated_at BEFORE UPDATE ON content_reviews 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_analytics_updated_at BEFORE UPDATE ON content_analytics 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_content_ratings_updated_at BEFORE UPDATE ON content_ratings 
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
