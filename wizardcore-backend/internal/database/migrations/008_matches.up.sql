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