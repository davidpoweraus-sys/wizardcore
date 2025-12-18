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