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