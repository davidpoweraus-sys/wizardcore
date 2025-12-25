package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type PreferencesRepository struct {
	db *sql.DB
}

func NewPreferencesRepository(db *sql.DB) *PreferencesRepository {
	return &PreferencesRepository{db: db}
}

// GetUserPreferences retrieves preferences for a user
func (r *PreferencesRepository) GetUserPreferences(ctx context.Context, userID uuid.UUID) (*models.UserPreferences, error) {
	query := `
		SELECT 
			user_id, theme, language, email_notifications, push_notifications,
			public_profile, show_progress, auto_save, sound_effects, two_factor_enabled,
			created_at, updated_at
		FROM user_preferences
		WHERE user_id = $1
	`

	var preferences models.UserPreferences
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&preferences.UserID,
		&preferences.Theme,
		&preferences.Language,
		&preferences.EmailNotifications,
		&preferences.PushNotifications,
		&preferences.PublicProfile,
		&preferences.ShowProgress,
		&preferences.AutoSave,
		&preferences.SoundEffects,
		&preferences.TwoFactorEnabled,
		&preferences.CreatedAt,
		&preferences.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		// Return default preferences
		return r.createDefaultPreferences(ctx, userID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user preferences: %w", err)
	}

	return &preferences, nil
}

// UpdateUserPreferences updates preferences for a user
func (r *PreferencesRepository) UpdateUserPreferences(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// Start building the query
	query := "UPDATE user_preferences SET "
	params := []interface{}{}
	paramCount := 1

	// Add each field to update
	for field, value := range updates {
		query += fmt.Sprintf("%s = $%d, ", field, paramCount)
		params = append(params, value)
		paramCount++
	}

	// Add updated_at
	query += fmt.Sprintf("updated_at = $%d ", paramCount)
	params = append(params, time.Now())
	paramCount++

	// Add WHERE clause
	query += fmt.Sprintf("WHERE user_id = $%d", paramCount)
	params = append(params, userID)

	// Execute the update
	result, err := r.db.ExecContext(ctx, query, params...)
	if err != nil {
		return fmt.Errorf("failed to update user preferences: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// No existing preferences, create with defaults and updates
		return r.createPreferencesWithUpdates(ctx, userID, updates)
	}

	return nil
}

// createDefaultPreferences creates default preferences for a user
func (r *PreferencesRepository) createDefaultPreferences(ctx context.Context, userID uuid.UUID) (*models.UserPreferences, error) {
	query := `
		INSERT INTO user_preferences (
			user_id, theme, language, email_notifications, push_notifications,
			public_profile, show_progress, auto_save, sound_effects, two_factor_enabled,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING created_at, updated_at
	`

	now := time.Now()
	defaultPrefs := &models.UserPreferences{
		UserID:             userID,
		Theme:              "dark",
		Language:           "en",
		EmailNotifications: true,
		PushNotifications:  false,
		PublicProfile:      true,
		ShowProgress:       true,
		AutoSave:           true,
		SoundEffects:       true,
		TwoFactorEnabled:   false,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err := r.db.QueryRowContext(ctx, query,
		defaultPrefs.UserID,
		defaultPrefs.Theme,
		defaultPrefs.Language,
		defaultPrefs.EmailNotifications,
		defaultPrefs.PushNotifications,
		defaultPrefs.PublicProfile,
		defaultPrefs.ShowProgress,
		defaultPrefs.AutoSave,
		defaultPrefs.SoundEffects,
		defaultPrefs.TwoFactorEnabled,
		defaultPrefs.CreatedAt,
		defaultPrefs.UpdatedAt,
	).Scan(&defaultPrefs.CreatedAt, &defaultPrefs.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create default preferences: %w", err)
	}

	return defaultPrefs, nil
}

// createPreferencesWithUpdates creates preferences with custom updates
func (r *PreferencesRepository) createPreferencesWithUpdates(ctx context.Context, userID uuid.UUID, updates map[string]interface{}) error {
	// Start with defaults
	defaultPrefs := models.UserPreferences{
		UserID:             userID,
		Theme:              "dark",
		Language:           "en",
		EmailNotifications: true,
		PushNotifications:  false,
		PublicProfile:      true,
		ShowProgress:       true,
		AutoSave:           true,
		SoundEffects:       true,
		TwoFactorEnabled:   false,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Apply updates
	if theme, ok := updates["theme"].(string); ok {
		defaultPrefs.Theme = theme
	}
	if language, ok := updates["language"].(string); ok {
		defaultPrefs.Language = language
	}
	if emailNotifications, ok := updates["email_notifications"].(bool); ok {
		defaultPrefs.EmailNotifications = emailNotifications
	}
	if pushNotifications, ok := updates["push_notifications"].(bool); ok {
		defaultPrefs.PushNotifications = pushNotifications
	}
	if publicProfile, ok := updates["public_profile"].(bool); ok {
		defaultPrefs.PublicProfile = publicProfile
	}
	if showProgress, ok := updates["show_progress"].(bool); ok {
		defaultPrefs.ShowProgress = showProgress
	}
	if autoSave, ok := updates["auto_save"].(bool); ok {
		defaultPrefs.AutoSave = autoSave
	}
	if soundEffects, ok := updates["sound_effects"].(bool); ok {
		defaultPrefs.SoundEffects = soundEffects
	}
	if twoFactorEnabled, ok := updates["two_factor_enabled"].(bool); ok {
		defaultPrefs.TwoFactorEnabled = twoFactorEnabled
	}

	query := `
		INSERT INTO user_preferences (
			user_id, theme, language, email_notifications, push_notifications,
			public_profile, show_progress, auto_save, sound_effects, two_factor_enabled,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err := r.db.ExecContext(ctx, query,
		defaultPrefs.UserID,
		defaultPrefs.Theme,
		defaultPrefs.Language,
		defaultPrefs.EmailNotifications,
		defaultPrefs.PushNotifications,
		defaultPrefs.PublicProfile,
		defaultPrefs.ShowProgress,
		defaultPrefs.AutoSave,
		defaultPrefs.SoundEffects,
		defaultPrefs.TwoFactorEnabled,
		defaultPrefs.CreatedAt,
		defaultPrefs.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create preferences with updates: %w", err)
	}

	return nil
}
