-- Create submissions table
CREATE TABLE IF NOT EXISTS submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    source_code TEXT NOT NULL,
    language_id INTEGER NOT NULL,
    
    -- Judge0 results
    judge0_token VARCHAR(255),
    status VARCHAR(50), -- 'pending', 'processing', 'accepted', 'wrong_answer', 'time_limit_exceeded', 'runtime_error', 'compilation_error'
    stdout TEXT,
    stderr TEXT,
    compile_output TEXT,
    execution_time DECIMAL(10,3), -- in seconds
    memory_used INTEGER, -- in KB
    
    -- Scoring
    test_cases_passed INTEGER DEFAULT 0,
    test_cases_total INTEGER DEFAULT 0,
    points_earned INTEGER DEFAULT 0,
    is_correct BOOLEAN DEFAULT false,
    
    -- Metadata
    submission_type VARCHAR(50) DEFAULT 'solution', -- 'draft', 'solution', 'practice'
    ip_address INET,
    user_agent TEXT,
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_submissions_user_id ON submissions(user_id);
CREATE INDEX idx_submissions_exercise_id ON submissions(exercise_id);
CREATE INDEX idx_submissions_status ON submissions(status);
CREATE INDEX idx_submissions_created_at ON submissions(created_at DESC);
-- Create submission test results
CREATE TABLE IF NOT EXISTS submission_test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    submission_id UUID REFERENCES submissions(id) ON DELETE CASCADE,
    test_case_id UUID REFERENCES test_cases(id) ON DELETE CASCADE,
    passed BOOLEAN NOT NULL,
    actual_output TEXT,
    execution_time DECIMAL(10,3),
    memory_used INTEGER,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_test_results_submission_id ON submission_test_results(submission_id);