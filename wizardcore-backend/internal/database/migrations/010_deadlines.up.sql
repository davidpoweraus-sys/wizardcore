-- Create user deadlines table
CREATE TABLE IF NOT EXISTS user_deadlines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    deadline_type VARCHAR(50) NOT NULL, -- 'project', 'quiz', 'lab', 'assignment'
    exercise_id UUID REFERENCES exercises(id) ON DELETE SET NULL,
    pathway_id UUID REFERENCES pathways(id) ON DELETE SET NULL,
    module_id UUID REFERENCES modules(id) ON DELETE SET NULL,
    due_date TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    is_completed BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_deadlines_user_id ON user_deadlines(user_id);
CREATE INDEX idx_deadlines_due_date ON user_deadlines(due_date);
CREATE INDEX idx_deadlines_is_completed ON user_deadlines(is_completed);