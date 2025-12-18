package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	SupabaseUserID  uuid.UUID  `json:"supabase_user_id" db:"supabase_user_id"`
	Email           string     `json:"email" db:"email" validate:"required,email"`
	DisplayName     *string    `json:"display_name,omitempty" db:"display_name"`
	AvatarURL       *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Bio             *string    `json:"bio,omitempty" db:"bio"`
	Location        *string    `json:"location,omitempty" db:"location"`
	Website         *string    `json:"website,omitempty" db:"website"`
	GithubUsername  *string    `json:"github_username,omitempty" db:"github_username"`
	TwitterUsername *string    `json:"twitter_username,omitempty" db:"twitter_username"`
	TotalXP         int        `json:"total_xp" db:"total_xp"`
	PracticeScore   int        `json:"practice_score" db:"practice_score"`
	GlobalRank      *int       `json:"global_rank,omitempty" db:"global_rank"`
	CurrentStreak   int        `json:"current_streak" db:"current_streak"`
	LongestStreak   int        `json:"longest_streak" db:"longest_streak"`
	LastActivityDate *time.Time `json:"last_activity_date,omitempty" db:"last_activity_date"`
	CreatedAt       time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at" db:"updated_at"`
}

type UserPreferences struct {
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Theme              string    `json:"theme" db:"theme"`
	Language           string    `json:"language" db:"language"`
	EmailNotifications bool      `json:"email_notifications" db:"email_notifications"`
	PushNotifications  bool      `json:"push_notifications" db:"push_notifications"`
	PublicProfile      bool      `json:"public_profile" db:"public_profile"`
	ShowProgress       bool      `json:"show_progress" db:"show_progress"`
	AutoSave           bool      `json:"auto_save" db:"auto_save"`
	SoundEffects       bool      `json:"sound_effects" db:"sound_effects"`
	TwoFactorEnabled   bool      `json:"two_factor_enabled" db:"two_factor_enabled"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type UserActivity struct {
	ID           uuid.UUID              `json:"id" db:"id"`
	UserID       uuid.UUID              `json:"user_id" db:"user_id"`
	ActivityType string                 `json:"activity_type" db:"activity_type"`
	Title        string                 `json:"title" db:"title"`
	Description  *string                `json:"description,omitempty" db:"description"`
	Icon         *string                `json:"icon,omitempty" db:"icon"`
	Color        *string                `json:"color,omitempty" db:"color"`
	Metadata     map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt    time.Time              `json:"created_at" db:"created_at"`
}

type UserStats struct {
	ActiveCourses          int    `json:"active_courses"`
	ActiveCoursesChange    string `json:"active_courses_change"`
	CompletionRate         int    `json:"completion_rate"`
	CompletionRateChange   string `json:"completion_rate_change"`
	StudyTimeHours         int    `json:"study_time_hours"`
	StudyTimeWeek          int    `json:"study_time_week"`
	XPTotal                int    `json:"xp_total"`
	XPToday                int    `json:"xp_today"`
}

type CreateUserRequest struct {
	SupabaseUserID uuid.UUID `json:"supabase_user_id" validate:"required"`
	Email          string    `json:"email" validate:"required,email"`
	DisplayName    string    `json:"display_name"`
}

type UpdateUserProfileRequest struct {
	DisplayName     *string `json:"display_name"`
	Bio             *string `json:"bio"`
	Location        *string `json:"location"`
	Website         *string `json:"website"`
	GithubUsername  *string `json:"github_username"`
	TwitterUsername *string `json:"twitter_username"`
}

type UpdatePreferencesRequest struct {
	Theme              *string `json:"theme"`
	Language           *string `json:"language"`
	EmailNotifications *bool   `json:"email_notifications"`
	PushNotifications  *bool   `json:"push_notifications"`
	PublicProfile      *bool   `json:"public_profile"`
	ShowProgress       *bool   `json:"show_progress"`
	AutoSave           *bool   `json:"auto_save"`
	SoundEffects       *bool   `json:"sound_effects"`
}