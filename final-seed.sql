-- Seed data for WizardCore - Final fixed version

-- Insert sample pathways (with empty prerequisites array)
INSERT INTO pathways (id, title, subtitle, description, level, duration_weeks, student_count, rating, module_count, color_gradient, icon, is_locked, sort_order, prerequisites, created_at, updated_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Python for Offensive Security', 'Master Python for hacking', 'Learn Python from scratch with a focus on offensive security techniques.', 'beginner', 8, 1500, 4.8, 5, 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', 'python', false, 1, '{}', NOW(), NOW()),
    ('22222222-2222-2222-2222-222222222222', 'C & Assembly: The Exploit Developer''s Core', 'Low-level exploit development', 'Dive deep into C and Assembly to understand memory corruption and write exploits.', 'advanced', 12, 800, 4.9, 6, 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', 'c', false, 2, '{}', NOW(), NOW()),
    ('33333333-3333-3333-3333-333333333333', 'JavaScript & Browser Exploitation', 'Hack the browser', 'Explore JavaScript vulnerabilities and browser exploitation techniques.', 'intermediate', 10, 1200, 4.7, 4, 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', 'javascript', false, 3, '{}', NOW(), NOW());

-- Insert sample modules for first pathway
INSERT INTO modules (id, pathway_id, title, description, sort_order, estimated_hours, xp_reward, created_at, updated_at)
VALUES
    ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'Python Basics', 'Learn Python syntax and basic programming concepts.', 1, 10, 100, NOW(), NOW()),
    ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'Networking with Python', 'Use Python for network scanning and packet manipulation.', 2, 15, 150, NOW(), NOW()),
    ('cccccccc-cccc-cccc-cccc-cccccccccccc', '11111111-1111-1111-1111-111111111111', 'Web Scraping & Automation', 'Automate tasks and scrape websites for reconnaissance.', 3, 12, 120, NOW(), NOW());

-- Insert sample exercises
INSERT INTO exercises (id, module_id, title, difficulty, points, time_limit_minutes, sort_order, objectives, content, examples, description, constraints, hints, starter_code, solution_code, language_id, tags, concurrent_solvers, total_submissions, total_completions, average_completion_time, created_at, updated_at)
VALUES
    ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Hello World', 'easy', 10, 5, 1, '{"Print Hello World"}', 'Write a Python program that prints "Hello, World!" to the console.', '{"example": "print(''Hello, World!'')"}', 'Basic output exercise', '{}', '{"Use the print function"}', 'print("Hello, World!")', 'print("Hello, World!")', 71, '{"python", "beginner"}', 0, 0, 0, NULL, NOW(), NOW()),
    ('eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Port Scanner', 'medium', 50, 30, 1, '{"Scan ports 1-1024", "Use sockets"}', 'Create a simple port scanner that checks which ports are open on localhost.', '{"example": "import socket"}', 'Network programming exercise', '{"Use only standard library"}', '{"Look up socket module"}', 'import socket\n\ndef scan_port(port):\n    # TODO: implement\n    pass', 'import socket\n\ndef scan_port(port):\n    try:\n        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)\n        sock.settimeout(1)\n        result = sock.connect_ex(("127.0.0.1", port))\n        sock.close()\n        return result == 0\n    except Exception:\n        return False', 71, '{"python", "networking", "security"}', 0, 0, 0, NULL, NOW(), NOW());

-- Insert sample achievements (with correct column names)
INSERT INTO achievements (id, title, description, icon, xp_reward, criteria_type, criteria_value, created_at, updated_at)
VALUES
    ('ffffffff-ffff-ffff-ffff-ffffffffffff', 'First Steps', 'Complete your first exercise', 'trophy', 100, 'exercises_completed', 1, NOW(), NOW()),
    ('gggggggg-gggg-gggg-gggg-gggggggggggg', 'Python Master', 'Complete all Python pathway exercises', 'python', 500, 'pathway_completed', 1, NOW(), NOW());

-- Insert a sample user (for testing) - using proper UUID format
INSERT INTO users (id, supabase_user_id, email, display_name, avatar_url, bio, location, website, github_username, twitter_username, total_xp, practice_score, global_rank, current_streak, longest_streak, last_activity_date, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', 'test@example.com', 'Test User', 'https://avatar.com/test', 'I love hacking!', 'London', 'https://example.com', 'testuser', 'testuser', 0, 0, NULL, 0, 0, NULL, NOW(), NOW());

-- Insert user preferences
INSERT INTO user_preferences (user_id, theme, language, email_notifications, push_notifications, public_profile, show_progress, auto_save, sound_effects, two_factor_enabled, created_at, updated_at)
VALUES
    ('00000000-0000-0000-0000-000000000001', 'dark', 'en', true, true, true, true, true, true, false, NOW(), NOW());