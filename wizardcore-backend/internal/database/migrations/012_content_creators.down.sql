-- Migration: 012_content_creators.down.sql
-- Rollback content creator support

-- Drop triggers first
DROP TRIGGER IF EXISTS update_creator_stats_pathways ON pathways;
DROP FUNCTION IF EXISTS update_content_creator_stats();
DROP FUNCTION IF EXISTS get_content_creator_stats(UUID);

-- Drop tables in reverse order
DROP TABLE IF EXISTS content_version_history;
DROP TABLE IF EXISTS content_reviews;
DROP TABLE IF EXISTS content_creator_profiles;

-- Remove columns from exercises
ALTER TABLE exercises DROP COLUMN IF EXISTS created_by;
ALTER TABLE exercises DROP COLUMN IF EXISTS status;
ALTER TABLE exercises DROP COLUMN IF EXISTS version;
ALTER TABLE exercises DROP COLUMN IF EXISTS published_at;
ALTER TABLE exercises DROP COLUMN IF EXISTS requires_approval;

-- Remove columns from modules
ALTER TABLE modules DROP COLUMN IF EXISTS created_by;
ALTER TABLE modules DROP COLUMN IF EXISTS status;
ALTER TABLE modules DROP COLUMN IF EXISTS version;
ALTER TABLE modules DROP COLUMN IF EXISTS published_at;

-- Remove columns from pathways
ALTER TABLE pathways DROP COLUMN IF EXISTS created_by;
ALTER TABLE pathways DROP COLUMN IF EXISTS status;
ALTER TABLE pathways DROP COLUMN IF EXISTS version;
ALTER TABLE pathways DROP COLUMN IF EXISTS published_at;
ALTER TABLE pathways DROP COLUMN IF EXISTS review_notes;

-- Remove role column from users
ALTER TABLE users DROP COLUMN IF EXISTS role;

-- Drop indexes
DROP INDEX IF EXISTS idx_users_role;
DROP INDEX IF EXISTS idx_pathways_created_by;
DROP INDEX IF EXISTS idx_pathways_status;
DROP INDEX IF EXISTS idx_modules_created_by;
DROP INDEX IF EXISTS idx_modules_status;
DROP INDEX IF EXISTS idx_exercises_created_by;
DROP INDEX IF EXISTS idx_exercises_status;
DROP INDEX IF EXISTS idx_content_reviews_content;
DROP INDEX IF EXISTS idx_content_reviews_status;
DROP INDEX IF EXISTS idx_content_version_history_content;
DROP INDEX IF EXISTS idx_content_version_history_version;