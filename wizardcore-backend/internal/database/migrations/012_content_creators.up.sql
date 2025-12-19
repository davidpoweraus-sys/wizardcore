-- Add content creator support
-- Migration: 012_content_creators.up.sql

-- Add role column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) DEFAULT 'student' CHECK (role IN ('student', 'content_creator', 'admin'));

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

-- Add content creator metadata to pathways
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;
ALTER TABLE pathways ADD COLUMN IF NOT EXISTS review_notes TEXT;

-- Add content creator metadata to modules
ALTER TABLE modules ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE modules ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE modules ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE modules ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;

-- Add content creator metadata to exercises
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS created_by UUID REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS status VARCHAR(50) DEFAULT 'published' CHECK (status IN ('draft', 'published', 'archived', 'review'));
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS published_at TIMESTAMP;
ALTER TABLE exercises ADD COLUMN IF NOT EXISTS requires_approval BOOLEAN DEFAULT false;

-- Create content reviews table for admin approval
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
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(content_type, content_id)
);

-- Create content version history table
CREATE TABLE IF NOT EXISTS content_version_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('pathway', 'module', 'exercise')),
    content_id UUID NOT NULL,
    version INTEGER NOT NULL,
    data JSONB NOT NULL, -- Full content data at this version
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    change_notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_pathways_created_by ON pathways(created_by);
CREATE INDEX idx_pathways_status ON pathways(status);
CREATE INDEX idx_modules_created_by ON modules(created_by);
CREATE INDEX idx_modules_status ON modules(status);
CREATE INDEX idx_exercises_created_by ON exercises(created_by);
CREATE INDEX idx_exercises_status ON exercises(status);
CREATE INDEX idx_content_reviews_content ON content_reviews(content_type, content_id);
CREATE INDEX idx_content_reviews_status ON content_reviews(status);
CREATE INDEX idx_content_version_history_content ON content_version_history(content_type, content_id);
CREATE INDEX idx_content_version_history_version ON content_version_history(content_type, content_id, version DESC);

-- Add trigger to update content_creator_profiles stats
CREATE OR REPLACE FUNCTION update_content_creator_stats()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_TABLE_NAME = 'pathways' AND NEW.created_by IS NOT NULL THEN
        UPDATE content_creator_profiles 
        SET total_content_created = total_content_created + 1,
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = NEW.created_by;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_creator_stats_pathways
AFTER INSERT ON pathways
FOR EACH ROW
EXECUTE FUNCTION update_content_creator_stats();

-- Create function to get content creator dashboard stats
CREATE OR REPLACE FUNCTION get_content_creator_stats(creator_id UUID)
RETURNS TABLE(
    total_pathways INTEGER,
    total_modules INTEGER,
    total_exercises INTEGER,
    total_students BIGINT,
    average_rating DECIMAL(3,2),
    pending_reviews INTEGER
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        (SELECT COUNT(*) FROM pathways WHERE created_by = creator_id AND status = 'published')::INTEGER,
        (SELECT COUNT(*) FROM modules WHERE created_by = creator_id AND status = 'published')::INTEGER,
        (SELECT COUNT(*) FROM exercises WHERE created_by = creator_id AND status = 'published')::INTEGER,
        (SELECT COUNT(DISTINCT user_id) FROM user_pathway_enrollments WHERE pathway_id IN 
            (SELECT id FROM pathways WHERE created_by = creator_id))::BIGINT,
        COALESCE((SELECT AVG(rating) FROM pathways WHERE created_by = creator_id AND rating > 0), 0.0),
        (SELECT COUNT(*) FROM content_reviews WHERE content_id IN 
            (SELECT id FROM pathways WHERE created_by = creator_id UNION
             SELECT id FROM modules WHERE created_by = creator_id UNION
             SELECT id FROM exercises WHERE created_by = creator_id)
            AND status = 'pending')::INTEGER;
END;
$$ LANGUAGE plpgsql;