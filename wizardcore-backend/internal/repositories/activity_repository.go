package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"go.uber.org/zap"
)

type ActivityRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewActivityRepository(db *sql.DB, logger *zap.Logger) *ActivityRepository {
	return &ActivityRepository{
		db:     db,
		logger: logger,
	}
}

// CreateActivity creates a new user activity
func (r *ActivityRepository) CreateActivity(ctx context.Context, activity *models.UserActivity) error {
	query := `
		INSERT INTO user_activities (
			id, user_id, activity_type, title, description, icon, color, metadata, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	// Convert metadata to JSONB
	var metadataJSON []byte
	if activity.Metadata != nil {
		var err error
		metadataJSON, err = json.Marshal(activity.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	_, err := r.db.ExecContext(ctx, query,
		activity.ID,
		activity.UserID,
		activity.ActivityType,
		activity.Title,
		activity.Description,
		activity.Icon,
		activity.Color,
		metadataJSON,
		activity.CreatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create activity: %w", err)
	}

	return nil
}

// GetUserActivities retrieves activities for a user with pagination
func (r *ActivityRepository) GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.UserActivity, error) {
	query := `
		SELECT 
			id, user_id, activity_type, title, description, icon, color, metadata, created_at
		FROM user_activities 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query user activities: %w", err)
	}
	defer rows.Close()

	var activities []models.UserActivity
	for rows.Next() {
		var activity models.UserActivity
		var description, icon, color sql.NullString
		var metadataJSON []byte

		err := rows.Scan(
			&activity.ID,
			&activity.UserID,
			&activity.ActivityType,
			&activity.Title,
			&description,
			&icon,
			&color,
			&metadataJSON,
			&activity.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan activity: %w", err)
		}

		// Handle nullable fields
		if description.Valid {
			activity.Description = &description.String
		}
		if icon.Valid {
			activity.Icon = &icon.String
		}
		if color.Valid {
			activity.Color = &color.String
		}

		// Parse metadata JSON
		if len(metadataJSON) > 0 {
			var metadata map[string]interface{}
			if err := json.Unmarshal(metadataJSON, &metadata); err != nil {
				return nil, fmt.Errorf("failed to unmarshal metadata: %w", err)
			}
			activity.Metadata = metadata
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating activities: %w", err)
	}

	return activities, nil
}

// CountUserActivities returns the total number of activities for a user
func (r *ActivityRepository) CountUserActivities(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM user_activities WHERE user_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user activities: %w", err)
	}

	return count, nil
}

// CreateExerciseSubmissionActivity creates an activity record for exercise submission
func (r *ActivityRepository) CreateExerciseSubmissionActivity(ctx context.Context, userID uuid.UUID, exerciseID uuid.UUID, exerciseTitle string, xpEarned int, timeSpentMinutes int) error {
	activity := &models.UserActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: "practice",
		Title:        fmt.Sprintf("Completed: %s", exerciseTitle),
		Description:  stringPtr(fmt.Sprintf("Earned %d XP in %d minutes", xpEarned, timeSpentMinutes)),
		Icon:         stringPtr("code"),
		Color:        stringPtr("text-blue-400"),
		Metadata: map[string]interface{}{
			"exercise_id":        exerciseID.String(),
			"xp_earned":          xpEarned,
			"time_spent_minutes": timeSpentMinutes,
			"activity_date":      time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
	}

	return r.CreateActivity(ctx, activity)
}

// CreateModuleCompletionActivity creates an activity record for module completion
func (r *ActivityRepository) CreateModuleCompletionActivity(ctx context.Context, userID uuid.UUID, moduleID uuid.UUID, moduleTitle string, xpEarned int) error {
	activity := &models.UserActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: "completion",
		Title:        fmt.Sprintf("Module Completed: %s", moduleTitle),
		Description:  stringPtr(fmt.Sprintf("Earned %d XP", xpEarned)),
		Icon:         stringPtr("book-open"),
		Color:        stringPtr("text-green-400"),
		Metadata: map[string]interface{}{
			"module_id":     moduleID.String(),
			"xp_earned":     xpEarned,
			"activity_date": time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
	}

	return r.CreateActivity(ctx, activity)
}

// CreatePathwayEnrollmentActivity creates an activity record for pathway enrollment
func (r *ActivityRepository) CreatePathwayEnrollmentActivity(ctx context.Context, userID uuid.UUID, pathwayID uuid.UUID, pathwayTitle string) error {
	activity := &models.UserActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: "completion",
		Title:        fmt.Sprintf("Enrolled in: %s", pathwayTitle),
		Description:  stringPtr("Started a new learning pathway"),
		Icon:         stringPtr("trophy"),
		Color:        stringPtr("text-yellow-400"),
		Metadata: map[string]interface{}{
			"pathway_id":    pathwayID.String(),
			"activity_date": time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
	}

	return r.CreateActivity(ctx, activity)
}

// CreateAchievementActivity creates an activity record for achievement unlock
func (r *ActivityRepository) CreateAchievementActivity(ctx context.Context, userID uuid.UUID, achievementID uuid.UUID, achievementTitle string, xpEarned int) error {
	activity := &models.UserActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: "achievement",
		Title:        fmt.Sprintf("Achievement Unlocked: %s", achievementTitle),
		Description:  stringPtr(fmt.Sprintf("Earned %d XP", xpEarned)),
		Icon:         stringPtr("trophy"),
		Color:        stringPtr("text-purple-400"),
		Metadata: map[string]interface{}{
			"achievement_id": achievementID.String(),
			"xp_earned":      xpEarned,
			"activity_date":  time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
	}

	return r.CreateActivity(ctx, activity)
}

// CreateStreakActivity creates an activity record for streak maintenance
func (r *ActivityRepository) CreateStreakActivity(ctx context.Context, userID uuid.UUID, streakDays int) error {
	activity := &models.UserActivity{
		ID:           uuid.New(),
		UserID:       userID,
		ActivityType: "streak",
		Title:        fmt.Sprintf("%d Day Streak!", streakDays),
		Description:  stringPtr("Keep up the great work!"),
		Icon:         stringPtr("clock"),
		Color:        stringPtr("text-neon-cyan"),
		Metadata: map[string]interface{}{
			"streak_days":   streakDays,
			"activity_date": time.Now().Format(time.RFC3339),
		},
		CreatedAt: time.Now(),
	}

	return r.CreateActivity(ctx, activity)
}

// Helper function to convert string to pointer
func stringPtr(s string) *string {
	return &s
}
