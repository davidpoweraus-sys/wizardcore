-- Populate achievements table with actual badges
-- This includes all achievements from the hardcoded frontend data

-- First, let's check what we already have
-- We already have "First Steps" and "Python Master" from previous seed

-- Add more achievements
INSERT INTO achievements (title, description, icon, color_gradient, rarity, xp_reward, criteria_type, criteria_value, sort_order) VALUES
-- Speed Coder: Solve 10 exercises under 5 minutes each
('Speed Coder', 'Solve 10 exercises under 5 minutes each', 'zap', 'linear-gradient(135deg, #00f5d4 0%, #9b5de5 100%)', 'Rare', 250, 'fast_exercises', 10, 2),

-- Week Warrior: Maintain a 7-day learning streak
('Week Warrior', 'Maintain a 7-day learning streak', 'flame', 'linear-gradient(135deg, #ff6b6b 0%, #ff8e8e 100%)', 'Uncommon', 150, 'streak_days', 7, 3),

-- Bug Hunter: Find and report a security vulnerability
('Bug Hunter', 'Find and report a security vulnerability', 'bullseye', 'linear-gradient(135deg, #8338ec 0%, #3a86ff 100%)', 'Legendary', 1000, 'vulnerabilities_reported', 1, 7),

-- Assembly Guru: Write functional x64 shellcode
('Assembly Guru', 'Write functional x64 shellcode', 'shield', 'linear-gradient(135deg, #4361ee 0%, #7209b7 100%)', 'Epic', 500, 'assembly_exercises', 5, 5),

-- 100% Club: Achieve perfect score on 20 exercises
('100% Club', 'Achieve perfect score on 20 exercises', 'check', 'linear-gradient(135deg, #f72585 0%, #b5179e 100%)', 'Rare', 300, 'perfect_scores', 20, 4),

-- Rootkit Researcher: Complete Rootkit Development course
('Rootkit Researcher', 'Complete Rootkit Development course', 'lock', 'linear-gradient(135deg, #495057 0%, #212529 100%)', 'Mythic', 1500, 'pathway_completed', 1, 8),

-- Additional achievements for better gamification
-- Daily Learner: Complete exercises for 30 consecutive days
('Daily Learner', 'Complete exercises for 30 consecutive days', 'calendar', 'linear-gradient(135deg, #ff9e00 0%, #ffd166 100%)', 'Epic', 750, 'streak_days', 30, 6),

-- Code Marathon: Complete 100 total exercises
('Code Marathon', 'Complete 100 total exercises', 'infinity', 'linear-gradient(135deg, #06d6a0 0%, #118ab2 100%)', 'Epic', 800, 'exercises_completed', 100, 9),

-- Early Bird: Complete 5 exercises before 8 AM
('Early Bird', 'Complete 5 exercises before 8 AM', 'sun', 'linear-gradient(135deg, #ffd166 0%, #ef476f 100%)', 'Uncommon', 200, 'early_exercises', 5, 10),

-- Night Owl: Complete 10 exercises after 10 PM
('Night Owl', 'Complete 10 exercises after 10 PM', 'moon', 'linear-gradient(135deg, #073b4c 0%, #118ab2 100%)', 'Uncommon', 200, 'late_exercises', 10, 11),

-- Social Butterfly: Share 5 achievements
('Social Butterfly', 'Share 5 achievements', 'share', 'linear-gradient(135deg, #7209b7 0%, #f72585 100%)', 'Common', 100, 'achievements_shared', 5, 12),

-- Practice Master: Complete 50 practice matches
('Practice Master', 'Complete 50 practice matches', 'swords', 'linear-gradient(135deg, #ff595e 0%, #ffca3a 100%)', 'Rare', 400, 'practice_matches', 50, 13),

-- Quiz Whiz: Score 90% or higher on 10 quizzes
('Quiz Whiz', 'Score 90% or higher on 10 quizzes', 'brain', 'linear-gradient(135deg, #8ac926 0%, #1982c4 100%)', 'Uncommon', 250, 'high_score_quizzes', 10, 14),

-- Project Pioneer: Complete 5 capstone projects
('Project Pioneer', 'Complete 5 capstone projects', 'rocket', 'linear-gradient(135deg, #ff595e 0%, #6a4c93 100%)', 'Epic', 600, 'projects_completed', 5, 15),

-- Community Helper: Help 10 other learners in forums
('Community Helper', 'Help 10 other learners in forums', 'heart', 'linear-gradient(135deg, #ff595e 0%, #ffca3a 100%)', 'Rare', 350, 'forum_helps', 10, 16),

-- Feedback Provider: Submit 20 exercise feedback reports
('Feedback Provider', 'Submit 20 exercise feedback reports', 'message', 'linear-gradient(135deg, #1982c4 0%, #8ac926 100%)', 'Uncommon', 200, 'feedback_submitted', 20, 17),

-- Pathfinder: Complete 3 different pathways
('Pathfinder', 'Complete 3 different pathways', 'compass', 'linear-gradient(135deg, #6a4c93 0%, #1982c4 100%)', 'Epic', 700, 'pathways_completed', 3, 18),

-- XP Collector: Earn 10,000 total XP
('XP Collector', 'Earn 10,000 total XP', 'treasure', 'linear-gradient(135deg, #ffd166 0%, #06d6a0 100%)', 'Legendary', 1200, 'total_xp', 10000, 19),

-- All-Rounder: Earn at least one achievement from each category
('All-Rounder', 'Earn at least one achievement from each category', 'star', 'linear-gradient(135deg, #ff595e 0%, #ffca3a 0%, #8ac926 0%, #1982c4 0%, #6a4c93 100%)', 'Mythic', 2000, 'category_achievements', 5, 20)
;

-- Update existing achievements with better metadata
UPDATE achievements SET 
  icon = 'star',
  color_gradient = 'linear-gradient(135deg, #ffd166 0%, #f3722c 100%)',
  rarity = 'Common'
WHERE title = 'First Steps';

UPDATE achievements SET 
  icon = 'python',
  color_gradient = 'linear-gradient(135deg, #3776ab 0%, #ffd343 100%)',
  rarity = 'Epic'
WHERE title = 'Python Master';