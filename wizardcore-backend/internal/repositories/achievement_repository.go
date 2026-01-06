package repositories

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type AchievementRepository struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) *AchievementRepository {
	return &AchievementRepository{db: db}
}

// FindAll returns all achievements
func (r *AchievementRepository) FindAll() ([]models.Achievement, error) {
	query := `
		SELECT id, title, description, icon, color_gradient, rarity, xp_reward, criteria_type, criteria_value, criteria_metadata, is_hidden, sort_order, created_at
		FROM achievements
		ORDER BY sort_order
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Achievement{}, err
	}
	defer rows.Close()

	achievements := make([]models.Achievement, 0)
	for rows.Next() {
		var a models.Achievement
		var criteriaValue sql.NullInt64
		var criteriaMetadata []byte
		err := rows.Scan(
			&a.ID,
			&a.Title,
			&a.Description,
			&a.Icon,
			&a.ColorGradient,
			&a.Rarity,
			&a.XPReward,
			&a.CriteriaType,
			&criteriaValue,
			&criteriaMetadata,
			&a.IsHidden,
			&a.SortOrder,
			&a.CreatedAt,
		)
		if err != nil {
			return []models.Achievement{}, err
		}
		if criteriaValue.Valid {
			val := int(criteriaValue.Int64)
			a.CriteriaValue = &val
		}
		if len(criteriaMetadata) > 0 {
			var metadata map[string]interface{}
			if err := json.Unmarshal(criteriaMetadata, &metadata); err == nil {
				a.CriteriaMetadata = metadata
			}
		}
		achievements = append(achievements, a)
	}
	return achievements, nil
}

// FindByID returns a single achievement by ID
func (r *AchievementRepository) FindByID(id uuid.UUID) (*models.Achievement, error) {
	query := `
		SELECT id, title, description, icon, color_gradient, rarity, xp_reward, criteria_type, criteria_value, criteria_metadata, is_hidden, sort_order, created_at
		FROM achievements
		WHERE id = $1
	`
	var a models.Achievement
	var criteriaValue sql.NullInt64
	var criteriaMetadata []byte
	err := r.db.QueryRow(query, id).Scan(
		&a.ID,
		&a.Title,
		&a.Description,
		&a.Icon,
		&a.ColorGradient,
		&a.Rarity,
		&a.XPReward,
		&a.CriteriaType,
		&criteriaValue,
		&criteriaMetadata,
		&a.IsHidden,
		&a.SortOrder,
		&a.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if criteriaValue.Valid {
		val := int(criteriaValue.Int64)
		a.CriteriaValue = &val
	}
	if len(criteriaMetadata) > 0 {
		var metadata map[string]interface{}
		if err := json.Unmarshal(criteriaMetadata, &metadata); err == nil {
			a.CriteriaMetadata = metadata
		}
	}
	return &a, nil
}

// GetUserAchievements returns achievements with progress for a specific user
func (r *AchievementRepository) GetUserAchievements(userID uuid.UUID) ([]models.AchievementWithProgress, error) {
	// Get all achievements
	allAchievements, err := r.FindAll()
	if err != nil {
		return []models.AchievementWithProgress{}, err
	}

	// Get user's earned achievements
	query := `
		SELECT achievement_id, progress, earned_at
		FROM user_achievements
		WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.AchievementWithProgress{}, err
	}
	defer rows.Close()

	earnedMap := make(map[uuid.UUID]struct {
		progress int
		earnedAt *time.Time
	})
	for rows.Next() {
		var achievementID uuid.UUID
		var progress int
		var earnedAt sql.NullTime
		err := rows.Scan(&achievementID, &progress, &earnedAt)
		if err != nil {
			return []models.AchievementWithProgress{}, err
		}
		var earnedTime *time.Time
		if earnedAt.Valid {
			earnedTime = &earnedAt.Time
		}
		earnedMap[achievementID] = struct {
			progress int
			earnedAt *time.Time
		}{progress, earnedTime}
	}

	// Combine
	result := make([]models.AchievementWithProgress, 0, len(allAchievements))
	for _, a := range allAchievements {
		awp := models.AchievementWithProgress{
			Achievement: a,
			Earned:      false,
			Progress:    0,
			EarnedDate:  nil,
		}
		if earned, ok := earnedMap[a.ID]; ok {
			awp.Earned = earned.earnedAt != nil
			awp.Progress = earned.progress
			awp.EarnedDate = earned.earnedAt
		}
		result = append(result, awp)
	}
	return result, nil
}

// UpsertUserAchievement creates or updates a user's achievement progress
func (r *AchievementRepository) UpsertUserAchievement(userID, achievementID uuid.UUID, progress int, earned bool) error {
	var earnedAt interface{}
	if earned {
		earnedAt = time.Now()
	} else {
		earnedAt = nil
	}
	query := `
		INSERT INTO user_achievements (id, user_id, achievement_id, progress, earned_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4)
		ON CONFLICT (user_id, achievement_id) DO UPDATE SET
			progress = EXCLUDED.progress,
			earned_at = EXCLUDED.earned_at
	`
	_, err := r.db.Exec(query, userID, achievementID, progress, earnedAt)
	return err
}

// CheckAndUnlock checks if a user meets criteria for any achievements and unlocks them
func (r *AchievementRepository) CheckAndUnlock(userID uuid.UUID, criteriaType string, criteriaValue int) ([]uuid.UUID, error) {
	// Find achievements with matching criteria type where user hasn't earned yet
	query := `
		SELECT a.id, a.criteria_value
		FROM achievements a
		LEFT JOIN user_achievements ua ON a.id = ua.achievement_id AND ua.user_id = $1 AND ua.earned_at IS NOT NULL
		WHERE a.criteria_type = $2 AND ua.achievement_id IS NULL
	`
	rows, err := r.db.Query(query, userID, criteriaType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var unlocked []uuid.UUID
	for rows.Next() {
		var achievementID uuid.UUID
		var requiredValue sql.NullInt64
		err := rows.Scan(&achievementID, &requiredValue)
		if err != nil {
			return nil, err
		}
		if !requiredValue.Valid || criteriaValue >= int(requiredValue.Int64) {
			// Unlock achievement
			err := r.UpsertUserAchievement(userID, achievementID, criteriaValue, true)
			if err != nil {
				return nil, err
			}
			unlocked = append(unlocked, achievementID)
		} else {
			// Update progress
			err := r.UpsertUserAchievement(userID, achievementID, criteriaValue, false)
			if err != nil {
				return nil, err
			}
		}
	}
	return unlocked, nil
}