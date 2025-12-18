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