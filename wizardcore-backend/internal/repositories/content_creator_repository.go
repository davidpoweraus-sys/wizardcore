package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type ContentCreatorRepository struct {
	db *sql.DB
}

func NewContentCreatorRepository(db *sql.DB) *ContentCreatorRepository {
	return &ContentCreatorRepository{db: db}
}

// Profile Operations

func (r *ContentCreatorRepository) CreateProfile(profile *models.ContentCreatorProfile) error {
	query := `
		INSERT INTO content_creator_profiles (
			user_id, bio, specialization, website, github_url, linkedin_url, twitter_url
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at, is_verified, verification_date, 
		          total_content_created, total_students, average_rating
	`
	return r.db.QueryRow(
		query,
		profile.UserID,
		profile.Bio,
		pq.Array(profile.Specialization),
		profile.Website,
		profile.GithubURL,
		profile.LinkedinURL,
		profile.TwitterURL,
	).Scan(
		&profile.ID,
		&profile.CreatedAt,
		&profile.UpdatedAt,
		&profile.IsVerified,
		&profile.VerificationDate,
		&profile.TotalContentCreated,
		&profile.TotalStudents,
		&profile.AverageRating,
	)
}

func (r *ContentCreatorRepository) GetProfileByUserID(userID uuid.UUID) (*models.ContentCreatorProfile, error) {
	query := `
		SELECT id, user_id, bio, specialization, website, github_url, linkedin_url, 
		       twitter_url, is_verified, verification_date, total_content_created, 
		       total_students, average_rating, created_at, updated_at
		FROM content_creator_profiles
		WHERE user_id = $1
	`
	profile := &models.ContentCreatorProfile{}
	err := r.db.QueryRow(query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Bio,
		pq.Array(&profile.Specialization),
		&profile.Website,
		&profile.GithubURL,
		&profile.LinkedinURL,
		&profile.TwitterURL,
		&profile.IsVerified,
		&profile.VerificationDate,
		&profile.TotalContentCreated,
		&profile.TotalStudents,
		&profile.AverageRating,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (r *ContentCreatorRepository) UpdateProfile(userID uuid.UUID, updates map[string]interface{}) error {
	query := `
		UPDATE content_creator_profiles 
		SET bio = COALESCE($2, bio),
		    specialization = COALESCE($3, specialization),
		    website = COALESCE($4, website),
		    github_url = COALESCE($5, github_url),
		    linkedin_url = COALESCE($6, linkedin_url),
		    twitter_url = COALESCE($7, twitter_url),
		    updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1
	`
	_, err := r.db.Exec(
		query,
		userID,
		updates["bio"],
		pq.Array(updates["specialization"]),
		updates["website"],
		updates["github_url"],
		updates["linkedin_url"],
		updates["twitter_url"],
	)
	return err
}

// Pathway Operations

func (r *ContentCreatorRepository) CreatePathway(pathway *models.Pathway, creatorID uuid.UUID) error {
	query := `
		INSERT INTO pathways (
			title, subtitle, description, level, duration_weeks, color_gradient, 
			icon, sort_order, prerequisites, created_by, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at, student_count, rating, module_count, 
		          is_locked, version, published_at
	`
	return r.db.QueryRow(
		query,
		pathway.Title,
		pathway.Subtitle,
		pathway.Description,
		pathway.Level,
		pathway.DurationWeeks,
		pathway.ColorGradient,
		pathway.Icon,
		pathway.SortOrder,
		pq.Array(pathway.Prerequisites),
		creatorID,
		"draft", // Default status for new pathways
	).Scan(
		&pathway.ID,
		&pathway.CreatedAt,
		&pathway.UpdatedAt,
		&pathway.StudentCount,
		&pathway.Rating,
		&pathway.ModuleCount,
		&pathway.IsLocked,
	)
}

func (r *ContentCreatorRepository) GetCreatorPathways(creatorID uuid.UUID, status string) ([]*models.Pathway, error) {
	query := `
		SELECT id, title, subtitle, description, level, duration_weeks, student_count, 
		       rating, module_count, color_gradient, icon, is_locked, sort_order, 
		       prerequisites, created_at, updated_at
		FROM pathways
		WHERE created_by = $1 AND ($2 = '' OR status = $2)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, creatorID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	pathways := []*models.Pathway{}
	for rows.Next() {
		p := &models.Pathway{}
		if err := rows.Scan(
			&p.ID,
			&p.Title,
			&p.Subtitle,
			&p.Description,
			&p.Level,
			&p.DurationWeeks,
			&p.StudentCount,
			&p.Rating,
			&p.ModuleCount,
			&p.ColorGradient,
			&p.Icon,
			&p.IsLocked,
			&p.SortOrder,
			pq.Array(&p.Prerequisites),
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		pathways = append(pathways, p)
	}
	return pathways, rows.Err()
}

func (r *ContentCreatorRepository) UpdatePathway(pathwayID, creatorID uuid.UUID, updates map[string]interface{}) error {
	query := `
		UPDATE pathways 
		SET title = COALESCE($3, title),
		    subtitle = COALESCE($4, subtitle),
		    description = COALESCE($5, description),
		    level = COALESCE($6, level),
		    duration_weeks = COALESCE($7, duration_weeks),
		    color_gradient = COALESCE($8, color_gradient),
		    icon = COALESCE($9, icon),
		    sort_order = COALESCE($10, sort_order),
		    status = COALESCE($11, status),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND created_by = $2
	`
	result, err := r.db.Exec(
		query,
		pathwayID,
		creatorID,
		updates["title"],
		updates["subtitle"],
		updates["description"],
		updates["level"],
		updates["duration_weeks"],
		updates["color_gradient"],
		updates["icon"],
		updates["sort_order"],
		updates["status"],
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pathway not found or unauthorized")
	}
	return nil
}

func (r *ContentCreatorRepository) DeletePathway(pathwayID, creatorID uuid.UUID) error {
	query := `DELETE FROM pathways WHERE id = $1 AND created_by = $2`
	result, err := r.db.Exec(query, pathwayID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("pathway not found or unauthorized")
	}
	return nil
}

// Module Operations

func (r *ContentCreatorRepository) CreateModule(module *models.Module, creatorID uuid.UUID) error {
	query := `
		INSERT INTO modules (
			pathway_id, title, description, sort_order, estimated_hours, xp_reward, 
			created_by, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at, version, published_at
	`
	return r.db.QueryRow(
		query,
		module.PathwayID,
		module.Title,
		module.Description,
		module.SortOrder,
		module.EstimatedHours,
		module.XPReward,
		creatorID,
		"draft",
	).Scan(
		&module.ID,
		&module.CreatedAt,
		&module.UpdatedAt,
	)
}

func (r *ContentCreatorRepository) GetCreatorModules(creatorID uuid.UUID, pathwayID *uuid.UUID) ([]*models.Module, error) {
	query := `
		SELECT id, pathway_id, title, description, sort_order, estimated_hours, 
		       xp_reward, created_at, updated_at
		FROM modules
		WHERE created_by = $1 AND ($2::uuid IS NULL OR pathway_id = $2)
		ORDER BY pathway_id, sort_order
	`
	rows, err := r.db.Query(query, creatorID, pathwayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	modules := []*models.Module{}
	for rows.Next() {
		m := &models.Module{}
		if err := rows.Scan(
			&m.ID,
			&m.PathwayID,
			&m.Title,
			&m.Description,
			&m.SortOrder,
			&m.EstimatedHours,
			&m.XPReward,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		modules = append(modules, m)
	}
	return modules, rows.Err()
}

func (r *ContentCreatorRepository) UpdateModule(moduleID, creatorID uuid.UUID, updates map[string]interface{}) error {
	query := `
		UPDATE modules 
		SET title = COALESCE($3, title),
		    description = COALESCE($4, description),
		    sort_order = COALESCE($5, sort_order),
		    estimated_hours = COALESCE($6, estimated_hours),
		    xp_reward = COALESCE($7, xp_reward),
		    status = COALESCE($8, status),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND created_by = $2
	`
	result, err := r.db.Exec(
		query,
		moduleID,
		creatorID,
		updates["title"],
		updates["description"],
		updates["sort_order"],
		updates["estimated_hours"],
		updates["xp_reward"],
		updates["status"],
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("module not found or unauthorized")
	}
	return nil
}

func (r *ContentCreatorRepository) DeleteModule(moduleID, creatorID uuid.UUID) error {
	query := `DELETE FROM modules WHERE id = $1 AND created_by = $2`
	result, err := r.db.Exec(query, moduleID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("module not found or unauthorized")
	}
	return nil
}

// Exercise Operations

func (r *ContentCreatorRepository) CreateExercise(exercise *models.Exercise, creatorID uuid.UUID) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	examplesJSON, err := json.Marshal(exercise.Examples)
	if err != nil {
		return fmt.Errorf("failed to marshal examples: %w", err)
	}

	query := `
		INSERT INTO exercises (
			module_id, title, difficulty, points, time_limit_minutes, sort_order,
			objectives, content, examples, description, constraints, hints,
			starter_code, solution_code, language_id, tags, created_by, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at, concurrent_solvers, total_submissions, 
		          total_completions, average_completion_time
	`
	err = tx.QueryRow(
		query,
		exercise.ModuleID,
		exercise.Title,
		exercise.Difficulty,
		exercise.Points,
		exercise.TimeLimitMinutes,
		exercise.SortOrder,
		pq.Array(exercise.Objectives),
		exercise.Content,
		examplesJSON,
		exercise.Description,
		pq.Array(exercise.Constraints),
		pq.Array(exercise.Hints),
		exercise.StarterCode,
		exercise.SolutionCode,
		exercise.LanguageID,
		pq.Array(exercise.Tags),
		creatorID,
		"draft",
	).Scan(
		&exercise.ID,
		&exercise.CreatedAt,
		&exercise.UpdatedAt,
		&exercise.ConcurrentSolvers,
		&exercise.TotalSubmissions,
		&exercise.TotalCompletions,
		&exercise.AvgCompletionTime,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ContentCreatorRepository) GetCreatorExercises(creatorID uuid.UUID, moduleID *uuid.UUID) ([]*models.Exercise, error) {
	query := `
		SELECT id, module_id, title, difficulty, points, time_limit_minutes, sort_order,
		       objectives, content, examples, description, constraints, hints,
		       starter_code, solution_code, language_id, tags, concurrent_solvers,
		       total_submissions, total_completions, average_completion_time,
		       created_at, updated_at
		FROM exercises
		WHERE created_by = $1 AND ($2::uuid IS NULL OR module_id = $2)
		ORDER BY module_id, sort_order
	`
	rows, err := r.db.Query(query, creatorID, moduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	exercises := []*models.Exercise{}
	for rows.Next() {
		e := &models.Exercise{}
		var examplesJSON []byte
		if err := rows.Scan(
			&e.ID,
			&e.ModuleID,
			&e.Title,
			&e.Difficulty,
			&e.Points,
			&e.TimeLimitMinutes,
			&e.SortOrder,
			pq.Array(&e.Objectives),
			&e.Content,
			&examplesJSON,
			&e.Description,
			pq.Array(&e.Constraints),
			pq.Array(&e.Hints),
			&e.StarterCode,
			&e.SolutionCode,
			&e.LanguageID,
			pq.Array(&e.Tags),
			&e.ConcurrentSolvers,
			&e.TotalSubmissions,
			&e.TotalCompletions,
			&e.AvgCompletionTime,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if len(examplesJSON) > 0 {
			if err := json.Unmarshal(examplesJSON, &e.Examples); err != nil {
				return nil, fmt.Errorf("failed to unmarshal examples: %w", err)
			}
		}

		exercises = append(exercises, e)
	}
	return exercises, rows.Err()
}

func (r *ContentCreatorRepository) UpdateExercise(exerciseID, creatorID uuid.UUID, updates map[string]interface{}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Build dynamic update query
	query := `
		UPDATE exercises 
		SET title = COALESCE($3, title),
		    difficulty = COALESCE($4, difficulty),
		    points = COALESCE($5, points),
		    time_limit_minutes = COALESCE($6, time_limit_minutes),
		    sort_order = COALESCE($7, sort_order),
		    objectives = COALESCE($8, objectives),
		    content = COALESCE($9, content),
		    examples = COALESCE($10, examples),
		    description = COALESCE($11, description),
		    constraints = COALESCE($12, constraints),
		    hints = COALESCE($13, hints),
		    starter_code = COALESCE($14, starter_code),
		    solution_code = COALESCE($15, solution_code),
		    language_id = COALESCE($16, language_id),
		    tags = COALESCE($17, tags),
		    status = COALESCE($18, status),
		    version = version + 1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND created_by = $2
	`

	// Prepare parameters
	examplesJSON := []byte("{}")
	if updates["examples"] != nil {
		examplesJSON, err = json.Marshal(updates["examples"])
		if err != nil {
			return fmt.Errorf("failed to marshal examples: %w", err)
		}
	}

	result, err := tx.Exec(
		query,
		exerciseID,
		creatorID,
		updates["title"],
		updates["difficulty"],
		updates["points"],
		updates["time_limit_minutes"],
		updates["sort_order"],
		pq.Array(updates["objectives"]),
		updates["content"],
		examplesJSON,
		updates["description"],
		pq.Array(updates["constraints"]),
		pq.Array(updates["hints"]),
		updates["starter_code"],
		updates["solution_code"],
		updates["language_id"],
		pq.Array(updates["tags"]),
		updates["status"],
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("exercise not found or unauthorized")
	}

	return tx.Commit()
}

func (r *ContentCreatorRepository) DeleteExercise(exerciseID, creatorID uuid.UUID) error {
	query := `DELETE FROM exercises WHERE id = $1 AND created_by = $2`
	result, err := r.db.Exec(query, exerciseID, creatorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("exercise not found or unauthorized")
	}
	return nil
}

// Test Case Operations

func (r *ContentCreatorRepository) CreateTestCase(testCase *models.TestCase) error {
	query := `
		INSERT INTO test_cases (
			exercise_id, input, expected_output, is_hidden, points, sort_order
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		testCase.ExerciseID,
		testCase.Input,
		testCase.ExpectedOutput,
		testCase.IsHidden,
		testCase.Points,
		testCase.SortOrder,
	).Scan(&testCase.ID, &testCase.CreatedAt)
}

func (r *ContentCreatorRepository) GetTestCasesByExercise(exerciseID uuid.UUID) ([]*models.TestCase, error) {
	query := `
		SELECT id, exercise_id, input, expected_output, is_hidden, points, sort_order, created_at
		FROM test_cases
		WHERE exercise_id = $1
		ORDER BY sort_order
	`
	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	testCases := []*models.TestCase{}
	for rows.Next() {
		tc := &models.TestCase{}
		if err := rows.Scan(
			&tc.ID,
			&tc.ExerciseID,
			&tc.Input,
			&tc.ExpectedOutput,
			&tc.IsHidden,
			&tc.Points,
			&tc.SortOrder,
			&tc.CreatedAt,
		); err != nil {
			return nil, err
		}
		testCases = append(testCases, tc)
	}
	return testCases, rows.Err()
}

// Statistics

func (r *ContentCreatorRepository) GetCreatorStats(creatorID uuid.UUID) (*models.CreatorStats, error) {
	stats := &models.CreatorStats{}

	// Get content counts
	countQuery := `
		SELECT 
			COUNT(DISTINCT CASE WHEN p.created_by = $1 THEN p.id END) as total_pathways,
			COUNT(DISTINCT CASE WHEN m.created_by = $1 THEN m.id END) as total_modules,
			COUNT(DISTINCT CASE WHEN e.created_by = $1 THEN e.id END) as total_exercises,
			COUNT(DISTINCT CASE WHEN p.created_by = $1 AND p.status = 'published' THEN p.id END) as published_pathways,
			COUNT(DISTINCT CASE WHEN p.created_by = $1 AND p.status = 'draft' THEN p.id END) as draft_pathways
		FROM pathways p
		FULL OUTER JOIN modules m ON m.pathway_id = p.id
		FULL OUTER JOIN exercises e ON e.module_id = m.id
	`
	var totalPathways, totalModules, totalExercises, publishedCount, draftCount int
	err := r.db.QueryRow(countQuery, creatorID).Scan(
		&totalPathways,
		&totalModules,
		&totalExercises,
		&publishedCount,
		&draftCount,
	)
	if err != nil {
		return nil, err
	}

	stats.TotalPathways = totalPathways
	stats.TotalModules = totalModules
	stats.TotalExercises = totalExercises
	stats.PublishedContent = publishedCount
	stats.DraftContent = draftCount

	// Get analytics data
	analyticsQuery := `
		SELECT 
			COALESCE(SUM(views), 0) as total_views,
			COALESCE(SUM(enrollments), 0) as total_enrollments,
			COALESCE(SUM(completions), 0) as total_completions,
			COALESCE(AVG(average_rating), 0) as avg_rating,
			COALESCE(SUM(total_ratings), 0) as total_ratings
		FROM content_analytics
		WHERE creator_id = $1
	`
	err = r.db.QueryRow(analyticsQuery, creatorID).Scan(
		&stats.TotalViews,
		&stats.TotalEnrollments,
		&stats.TotalCompletions,
		&stats.AverageRating,
		&stats.TotalRatings,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Calculate completion rate
	if stats.TotalEnrollments > 0 {
		stats.CompletionRate = float64(stats.TotalCompletions) / float64(stats.TotalEnrollments) * 100
	}

	// Get pending reviews count
	reviewQuery := `
		SELECT COUNT(*) FROM content_reviews
		WHERE content_id IN (
			SELECT id FROM pathways WHERE created_by = $1
			UNION SELECT id FROM modules WHERE created_by = $1
			UNION SELECT id FROM exercises WHERE created_by = $1
		) AND status = 'pending'
	`
	err = r.db.QueryRow(reviewQuery, creatorID).Scan(&stats.PendingReviews)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return stats, nil
}

// Content Reviews

func (r *ContentCreatorRepository) CreateReview(review *models.ContentReview) error {
	query := `
		INSERT INTO content_reviews (
			content_type, content_id, status, revision_notes
		) VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(
		query,
		review.ContentType,
		review.ContentID,
		"pending",
		review.RevisionNotes,
	).Scan(&review.ID, &review.CreatedAt, &review.UpdatedAt)
}

func (r *ContentCreatorRepository) GetReviewsByCreator(creatorID uuid.UUID) ([]*models.ContentReview, error) {
	query := `
		SELECT cr.id, cr.content_type, cr.content_id, cr.reviewer_id, cr.status,
		       cr.review_notes, cr.revision_notes, cr.reviewed_at, cr.created_at, cr.updated_at
		FROM content_reviews cr
		LEFT JOIN pathways p ON cr.content_type = 'pathway' AND cr.content_id = p.id
		LEFT JOIN modules m ON cr.content_type = 'module' AND cr.content_id = m.id
		LEFT JOIN exercises e ON cr.content_type = 'exercise' AND cr.content_id = e.id
		WHERE p.created_by = $1 OR m.created_by = $1 OR e.created_by = $1
		ORDER BY cr.created_at DESC
	`
	rows, err := r.db.Query(query, creatorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []*models.ContentReview{}
	for rows.Next() {
		r := &models.ContentReview{}
		if err := rows.Scan(
			&r.ID,
			&r.ContentType,
			&r.ContentID,
			&r.ReviewerID,
			&r.Status,
			&r.ReviewNotes,
			&r.RevisionNotes,
			&r.ReviewedAt,
			&r.CreatedAt,
			&r.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, r)
	}
	return reviews, rows.Err()
}

func (r *ContentCreatorRepository) UpdateReview(reviewID uuid.UUID, status, reviewNotes string, reviewerID uuid.UUID) error {
	query := `
		UPDATE content_reviews 
		SET status = $2, 
		    review_notes = $3, 
		    reviewer_id = $4, 
		    reviewed_at = CURRENT_TIMESTAMP,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	result, err := r.db.Exec(query, reviewID, status, reviewNotes, reviewerID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("review not found")
	}
	return nil
}

// Version History

func (r *ContentCreatorRepository) CreateVersionHistory(history *models.ContentVersionHistory) error {
	dataJSON, err := json.Marshal(history.Data)
	if err != nil {
		return fmt.Errorf("failed to marshal version data: %w", err)
	}

	query := `
		INSERT INTO content_version_history (
			content_type, content_id, version, data, created_by, change_notes
		) VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`
	return r.db.QueryRow(
		query,
		history.ContentType,
		history.ContentID,
		history.Version,
		dataJSON,
		history.CreatedBy,
		history.ChangeNotes,
	).Scan(&history.ID, &history.CreatedAt)
}

// Helper to check ownership
func (r *ContentCreatorRepository) IsContentOwner(contentType string, contentID, userID uuid.UUID) (bool, error) {
	var query string
	switch contentType {
	case "pathway":
		query = "SELECT EXISTS(SELECT 1 FROM pathways WHERE id = $1 AND created_by = $2)"
	case "module":
		query = "SELECT EXISTS(SELECT 1 FROM modules WHERE id = $1 AND created_by = $2)"
	case "exercise":
		query = "SELECT EXISTS(SELECT 1 FROM exercises WHERE id = $1 AND created_by = $2)"
	default:
		return false, fmt.Errorf("invalid content type")
	}

	var exists bool
	err := r.db.QueryRow(query, contentID, userID).Scan(&exists)
	return exists, err
}

// Export/Import Operations

func (r *ContentCreatorRepository) ExportPathway(pathwayID, creatorID uuid.UUID) (*models.ExportPathway, error) {
	// Verify ownership
	isOwner, err := r.IsContentOwner("pathway", pathwayID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own this pathway")
	}

	// Get pathway
	var pathway models.ExportPathway
	query := `
		SELECT title, subtitle, description, level, duration_weeks, 
		       color_gradient, icon, prerequisites, sort_order
		FROM pathways 
		WHERE id = $1
	`
	err = r.db.QueryRow(query, pathwayID).Scan(
		&pathway.Title,
		&pathway.Subtitle,
		&pathway.Description,
		&pathway.Level,
		&pathway.DurationWeeks,
		&pathway.ColorGradient,
		&pathway.Icon,
		pq.Array(&pathway.Prerequisites),
		&pathway.SortOrder,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get pathway: %w", err)
	}

	// Get modules for this pathway
	modulesQuery := `
		SELECT id, title, description, sort_order, estimated_hours, xp_reward
		FROM modules 
		WHERE pathway_id = $1 AND created_by = $2
		ORDER BY sort_order
	`
	moduleRows, err := r.db.Query(modulesQuery, pathwayID, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules: %w", err)
	}
	defer moduleRows.Close()

	for moduleRows.Next() {
		var module models.ExportModule
		var moduleID uuid.UUID
		err := moduleRows.Scan(
			&moduleID,
			&module.Title,
			&module.Description,
			&module.SortOrder,
			&module.EstimatedHours,
			&module.XPReward,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan module: %w", err)
		}

		// Get exercises for this module
		exercisesQuery := `
			SELECT title, difficulty, points, time_limit_minutes, sort_order,
			       objectives, content, examples, description, constraints,
			       hints, starter_code, solution_code, language_id, tags
			FROM exercises 
			WHERE module_id = $1 AND created_by = $2
			ORDER BY sort_order
		`
		exerciseRows, err := r.db.Query(exercisesQuery, moduleID, creatorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get exercises: %w", err)
		}
		defer exerciseRows.Close()

		for exerciseRows.Next() {
			var exercise models.ExportExercise
			var exerciseID uuid.UUID
			var objectives, constraints, hints, tags pq.StringArray
			var examplesJSON []byte
			var content, description, starterCode, solutionCode *string

			err := exerciseRows.Scan(
				&exercise.Title,
				&exercise.Difficulty,
				&exercise.Points,
				&exercise.TimeLimitMinutes,
				&exercise.SortOrder,
				&objectives,
				&content,
				&examplesJSON,
				&description,
				&constraints,
				&hints,
				&starterCode,
				&solutionCode,
				&exercise.LanguageID,
				&tags,
			)
			if err != nil {
				return nil, fmt.Errorf("failed to scan exercise: %w", err)
			}

			exercise.Objectives = objectives
			exercise.Constraints = constraints
			exercise.Hints = hints
			exercise.Tags = tags
			exercise.Content = content
			exercise.Description = description
			exercise.StarterCode = starterCode
			exercise.SolutionCode = solutionCode

			// Parse examples JSON
			if len(examplesJSON) > 0 {
				if err := json.Unmarshal(examplesJSON, &exercise.Examples); err != nil {
					return nil, fmt.Errorf("failed to unmarshal examples: %w", err)
				}
			}

			// Get test cases for this exercise
			testCasesQuery := `
				SELECT input, expected_output, is_hidden, points, sort_order
				FROM test_cases 
				WHERE exercise_id = $1
				ORDER BY sort_order
			`
			testCaseRows, err := r.db.Query(testCasesQuery, exerciseID)
			if err != nil {
				return nil, fmt.Errorf("failed to get test cases: %w", err)
			}
			defer testCaseRows.Close()

			for testCaseRows.Next() {
				var testCase models.ExportTestCase
				err := testCaseRows.Scan(
					&testCase.Input,
					&testCase.ExpectedOutput,
					&testCase.IsHidden,
					&testCase.Points,
					&testCase.SortOrder,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to scan test case: %w", err)
				}
				exercise.TestCases = append(exercise.TestCases, testCase)
			}

			module.Exercises = append(module.Exercises, exercise)
		}
		exerciseRows.Close()

		pathway.Modules = append(pathway.Modules, module)
	}
	moduleRows.Close()

	return &pathway, nil
}

func (r *ContentCreatorRepository) ImportPathway(creatorID uuid.UUID, pathway *models.ExportPathway, status string) (*models.Pathway, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Create pathway
	pathwayQuery := `
		INSERT INTO pathways (
			title, subtitle, description, level, duration_weeks, 
			color_gradient, icon, prerequisites, sort_order, created_by, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`
	var newPathway models.Pathway
	err = tx.QueryRow(
		pathwayQuery,
		pathway.Title,
		pathway.Subtitle,
		pathway.Description,
		pathway.Level,
		pathway.DurationWeeks,
		pathway.ColorGradient,
		pathway.Icon,
		pq.Array(pathway.Prerequisites),
		pathway.SortOrder,
		creatorID,
		status,
	).Scan(&newPathway.ID, &newPathway.CreatedAt, &newPathway.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create pathway: %w", err)
	}

	// Create modules
	for _, module := range pathway.Modules {
		moduleQuery := `
			INSERT INTO modules (
				pathway_id, title, description, sort_order, 
				estimated_hours, xp_reward, created_by, status
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id, created_at, updated_at
		`
		var moduleID uuid.UUID
		var createdAt, updatedAt time.Time
		err = tx.QueryRow(
			moduleQuery,
			newPathway.ID,
			module.Title,
			module.Description,
			module.SortOrder,
			module.EstimatedHours,
			module.XPReward,
			creatorID,
			status,
		).Scan(&moduleID, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to create module: %w", err)
		}

		// Create exercises
		for _, exercise := range module.Exercises {
			examplesJSON, err := json.Marshal(exercise.Examples)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal examples: %w", err)
			}

			exerciseQuery := `
				INSERT INTO exercises (
					module_id, title, difficulty, points, time_limit_minutes, sort_order,
					objectives, content, examples, description, constraints, hints,
					starter_code, solution_code, language_id, tags, created_by, status
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
				RETURNING id, created_at, updated_at
			`
			var exerciseID uuid.UUID
			var exCreatedAt, exUpdatedAt time.Time
			err = tx.QueryRow(
				exerciseQuery,
				moduleID,
				exercise.Title,
				exercise.Difficulty,
				exercise.Points,
				exercise.TimeLimitMinutes,
				exercise.SortOrder,
				pq.Array(exercise.Objectives),
				exercise.Content,
				examplesJSON,
				exercise.Description,
				pq.Array(exercise.Constraints),
				pq.Array(exercise.Hints),
				exercise.StarterCode,
				exercise.SolutionCode,
				exercise.LanguageID,
				pq.Array(exercise.Tags),
				creatorID,
				status,
			).Scan(&exerciseID, &exCreatedAt, &exUpdatedAt)
			if err != nil {
				return nil, fmt.Errorf("failed to create exercise: %w", err)
			}

			// Create test cases
			for _, testCase := range exercise.TestCases {
				testCaseQuery := `
					INSERT INTO test_cases (
						exercise_id, input, expected_output, is_hidden, points, sort_order
					) VALUES ($1, $2, $3, $4, $5, $6)
				`
				_, err = tx.Exec(
					testCaseQuery,
					exerciseID,
					testCase.Input,
					testCase.ExpectedOutput,
					testCase.IsHidden,
					testCase.Points,
					testCase.SortOrder,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to create test case: %w", err)
				}
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &newPathway, nil
}
