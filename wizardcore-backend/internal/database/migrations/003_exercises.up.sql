-- Create exercises table
CREATE TABLE IF NOT EXISTS exercises (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    module_id UUID REFERENCES modules(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    difficulty VARCHAR(50) NOT NULL, -- 'BEGINNER', 'INTERMEDIATE', 'ADVANCED'
    points INTEGER DEFAULT 0,
    time_limit_minutes INTEGER,
    sort_order INTEGER NOT NULL,
    
    -- Lesson content
    objectives TEXT[],
    content TEXT, -- Markdown
    examples JSONB,
    
    -- Exercise details
    description TEXT,
    constraints TEXT[],
    hints TEXT[],
    starter_code TEXT,
    solution_code TEXT,
    
    -- Language
    language_id INTEGER NOT NULL, -- Judge0 language ID
    
    -- Metadata
    tags TEXT[],
    concurrent_solvers INTEGER DEFAULT 0,
    total_submissions INTEGER DEFAULT 0,
    total_completions INTEGER DEFAULT 0,
    average_completion_time INTEGER, -- in seconds
    
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_exercises_module_id ON exercises(module_id);
CREATE INDEX idx_exercises_difficulty ON exercises(difficulty);
-- Create test cases table
CREATE TABLE IF NOT EXISTS test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id UUID REFERENCES exercises(id) ON DELETE CASCADE,
    input TEXT,
    expected_output TEXT,
    is_hidden BOOLEAN DEFAULT false,
    points INTEGER DEFAULT 0,
    sort_order INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_test_cases_exercise_id ON test_cases(exercise_id);