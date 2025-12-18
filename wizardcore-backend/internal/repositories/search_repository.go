package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/yourusername/wizardcore-backend/internal/models"
)

type SearchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *SearchRepository {
	return &SearchRepository{db: db}
}

func (r *SearchRepository) Search(query string, limit, offset int) (*models.SearchResults, error) {
	if query == "" {
		return &models.SearchResults{
			Pathways:  []models.Pathway{},
			Exercises: []models.Exercise{},
			Users:     []models.User{},
		}, nil
	}

	// Search pathways
	pathways, err := r.searchPathways(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search pathways: %w", err)
	}

	// Search exercises
	exercises, err := r.searchExercises(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search exercises: %w", err)
	}

	// Search users
	users, err := r.searchUsers(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("search users: %w", err)
	}

	return &models.SearchResults{
		Pathways:  pathways,
		Exercises: exercises,
		Users:     users,
	}, nil
}

func (r *SearchRepository) searchPathways(query string, limit, offset int) ([]models.Pathway, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.Query(`
		SELECT id, title, subtitle, description, level, duration_weeks, student_count, rating, module_count, color_gradient, icon, is_locked, sort_order, prerequisites, created_at, updated_at
		FROM pathways
		WHERE LOWER(title) LIKE $1 OR LOWER(description) LIKE $2
		ORDER BY title
		LIMIT $3 OFFSET $4
	`, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pathways []models.Pathway
	for rows.Next() {
		var p models.Pathway
		err := rows.Scan(
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
			&p.Prerequisites,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		pathways = append(pathways, p)
	}
	return pathways, nil
}

func (r *SearchRepository) searchExercises(query string, limit, offset int) ([]models.Exercise, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.Query(`
		SELECT id, module_id, title, difficulty, points, time_limit_minutes, sort_order, objectives, content, examples, description, constraints, hints, starter_code, solution_code, language_id, tags, concurrent_solvers, total_submissions, total_completions, average_completion_time, created_at, updated_at
		FROM exercises
		WHERE LOWER(title) LIKE $1 OR LOWER(description) LIKE $2
		ORDER BY title
		LIMIT $3 OFFSET $4
	`, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		err := rows.Scan(
			&e.ID,
			&e.ModuleID,
			&e.Title,
			&e.Difficulty,
			&e.Points,
			&e.TimeLimitMinutes,
			&e.SortOrder,
			&e.Objectives,
			&e.Content,
			&e.Examples,
			&e.Description,
			&e.Constraints,
			&e.Hints,
			&e.StarterCode,
			&e.SolutionCode,
			&e.LanguageID,
			&e.Tags,
			&e.ConcurrentSolvers,
			&e.TotalSubmissions,
			&e.TotalCompletions,
			&e.AvgCompletionTime,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		exercises = append(exercises, e)
	}
	return exercises, nil
}

func (r *SearchRepository) searchUsers(query string, limit, offset int) ([]models.User, error) {
	searchTerm := "%" + strings.ToLower(query) + "%"
	rows, err := r.db.Query(`
		SELECT id, supabase_user_id, email, display_name, avatar_url, bio, location, website, github_username, twitter_username, total_xp, practice_score, global_rank, current_streak, longest_streak, last_activity_date, created_at, updated_at
		FROM users
		WHERE LOWER(email) LIKE $1 OR LOWER(display_name) LIKE $2
		ORDER BY display_name
		LIMIT $3 OFFSET $4
	`, searchTerm, searchTerm, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		err := rows.Scan(
			&u.ID,
			&u.SupabaseUserID,
			&u.Email,
			&u.DisplayName,
			&u.AvatarURL,
			&u.Bio,
			&u.Location,
			&u.Website,
			&u.GithubUsername,
			&u.TwitterUsername,
			&u.TotalXP,
			&u.PracticeScore,
			&u.GlobalRank,
			&u.CurrentStreak,
			&u.LongestStreak,
			&u.LastActivityDate,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}