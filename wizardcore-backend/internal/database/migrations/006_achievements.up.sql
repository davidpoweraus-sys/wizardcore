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