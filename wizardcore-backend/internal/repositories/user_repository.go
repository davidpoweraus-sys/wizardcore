package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
		INSERT INTO users (
			id, supabase_user_id, email, display_name, avatar_url, bio,
			location, website, github_username, twitter_username,
			total_xp, practice_score, global_rank, current_streak,
			longest_streak, last_activity_date, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
		RETURNING id, created_at, updated_at
	`
	// Generate a new UUID if not set
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}
	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.ID,
		user.SupabaseUserID,
		user.Email,
		user.DisplayName,
		user.AvatarURL,
		user.Bio,
		user.Location,
		user.Website,
		user.GithubUsername,
		user.TwitterUsername,
		user.TotalXP,
		user.PracticeScore,
		user.GlobalRank,
		user.CurrentStreak,
		user.LongestStreak,
		user.LastActivityDate,
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, supabase_user_id, email, display_name, avatar_url, bio,
			location, website, github_username, twitter_username,
			total_xp, practice_score, global_rank, current_streak,
			longest_streak, last_activity_date, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.SupabaseUserID,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
		&user.Bio,
		&user.Location,
		&user.Website,
		&user.GithubUsername,
		&user.TwitterUsername,
		&user.TotalXP,
		&user.PracticeScore,
		&user.GlobalRank,
		&user.CurrentStreak,
		&user.LongestStreak,
		&user.LastActivityDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by ID: %w", err)
	}
	return user, nil
}

func (r *UserRepository) FindBySupabaseUserID(supabaseUserID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, supabase_user_id, email, display_name, avatar_url, bio,
			location, website, github_username, twitter_username,
			total_xp, practice_score, global_rank, current_streak,
			longest_streak, last_activity_date, created_at, updated_at
		FROM users
		WHERE supabase_user_id = $1
	`
	user := &models.User{}
	err := r.db.QueryRow(query, supabaseUserID).Scan(
		&user.ID,
		&user.SupabaseUserID,
		&user.Email,
		&user.DisplayName,
		&user.AvatarURL,
		&user.Bio,
		&user.Location,
		&user.Website,
		&user.GithubUsername,
		&user.TwitterUsername,
		&user.TotalXP,
		&user.PracticeScore,
		&user.GlobalRank,
		&user.CurrentStreak,
		&user.LongestStreak,
		&user.LastActivityDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find user by Supabase user ID: %w", err)
	}
	return user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	query := `
		UPDATE users
		SET
			email = $2,
			display_name = $3,
			avatar_url = $4,
			bio = $5,
			location = $6,
			website = $7,
			github_username = $8,
			twitter_username = $9,
			total_xp = $10,
			practice_score = $11,
			global_rank = $12,
			current_streak = $13,
			longest_streak = $14,
			last_activity_date = $15,
			updated_at = $16
		WHERE id = $1
		RETURNING updated_at
	`
	now := time.Now()
	err := r.db.QueryRow(
		query,
		user.ID,
		user.Email,
		user.DisplayName,
		user.AvatarURL,
		user.Bio,
		user.Location,
		user.Website,
		user.GithubUsername,
		user.TwitterUsername,
		user.TotalXP,
		user.PracticeScore,
		user.GlobalRank,
		user.CurrentStreak,
		user.LongestStreak,
		user.LastActivityDate,
		now,
	).Scan(&user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}