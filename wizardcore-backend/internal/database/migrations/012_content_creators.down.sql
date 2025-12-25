-- Drop triggers
DROP TRIGGER IF EXISTS update_content_creator_profiles_updated_at ON content_creator_profiles;
DROP TRIGGER IF EXISTS update_content_reviews_updated_at ON content_reviews;
DROP TRIGGER IF EXISTS update_content_analytics_updated_at ON content_analytics;
DROP TRIGGER IF EXISTS update_content_ratings_updated_at ON content_ratings;

-- Drop tables
DROP TABLE IF EXISTS content_ratings;
DROP TABLE IF EXISTS content_analytics;
DROP TABLE IF EXISTS content_version_history;
DROP TABLE IF EXISTS content_reviews;
DROP TABLE IF EXISTS content_creator_profiles;

-- Remove columns from exercises
ALTER TABLE exercises DROP COLUMN IF EXISTS requires_approval;
ALTER TABLE exercises DROP COLUMN IF EXISTS published_at;
ALTER TABLE exercises DROP COLUMN IF EXISTS version;
ALTER TABLE exercises DROP COLUMN IF EXISTS status;
ALTER TABLE exercises DROP COLUMN IF EXISTS created_by;

-- Remove columns from modules
ALTER TABLE modules DROP COLUMN IF EXISTS published_at;
ALTER TABLE modules DROP COLUMN IF EXISTS version;
ALTER TABLE modules DROP COLUMN IF EXISTS status;
ALTER TABLE modules DROP COLUMN IF EXISTS created_by;

-- Remove columns from pathways
ALTER TABLE pathways DROP COLUMN IF EXISTS review_notes;
ALTER TABLE pathways DROP COLUMN IF EXISTS published_at;
ALTER TABLE pathways DROP COLUMN IF EXISTS version;
ALTER TABLE pathways DROP COLUMN IF EXISTS status;
ALTER TABLE pathways DROP COLUMN IF EXISTS created_by;

-- Remove role column from users
ALTER TABLE users DROP COLUMN IF EXISTS role;
