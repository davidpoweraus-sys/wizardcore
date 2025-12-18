package models

import (
	"time"

	"github.com/google/uuid"
)

// Progress represents user progress on a module
type Progress struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	UserID             uuid.UUID  `json:"user_id" db:"user_id"`
	ModuleID           uuid.UUID  `json:"module_id" db:"module_id"`
	PathwayID          uuid.UUID  `json:"pathway_id" db:"pathway_id"`
	ProgressPercentage int        `json:"progress_percentage" db:"progress_percentage"`
	CompletedExercises int        `json:"completed_exercises" db:"completed_exercises"`
	TotalExercises     int        `json:"total_exercises" db:"total_exercises"`
	XPEarned           int        `json:"xp_earned" db:"xp_earned"`
	TimeSpentMinutes   int        `json:"time_spent_minutes" db:"time_spent_minutes"`
	StartedAt          time.Time  `json:"started_at" db:"started_at"`
	CompletedAt        *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	LastActivityAt     *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
}

// ProgressResponse is the API response for GET /api/v1/users/me/progress
type ProgressResponse struct {
	Pathways []PathwayProgress `json:"pathways"`
	Totals   ProgressTotals    `json:"totals"`
}

// ProgressTotals aggregates overall progress
type ProgressTotals struct {
	TotalXP          int `json:"total_xp"`
	XPThisWeek       int `json:"xp_this_week"`
	OverallProgress  int `json:"overall_progress"`
	CurrentStreak    int `json:"current_streak"`
	ModulesCompleted int `json:"modules_completed"`
	ModulesTotal     int `json:"modules_total"`
}

// Milestone represents a user milestone
type Milestone struct {
	ID            uuid.UUID `json:"id" db:"id"`
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	Title         string    `json:"title" db:"title"`
	Description   *string   `json:"description,omitempty" db:"description"`
	MilestoneType string    `json:"milestone_type" db:"milestone_type"`
	XPAwarded     int       `json:"xp_awarded" db:"xp_awarded"`
	AchievedAt    time.Time `json:"achieved_at" db:"achieved_at"`
}

// WeeklyActivity represents weekly activity data
type WeeklyActivity struct {
	WeeklyData        []DailyActivity `json:"weekly_data"`
	AvgDailyTime      int             `json:"avg_daily_time_minutes"`
	CompletionRate    int             `json:"completion_rate"`
	CurrentStreak     int             `json:"current_streak"`
	TrendPercentage   int             `json:"trend_percentage"`
}

// DailyActivity represents a single day's activity
type DailyActivity struct {
	Day   string `json:"day"`
	Value int    `json:"value"`
	Hours int    `json:"hours"`
}

// WeeklyHours represents weekly hours spent
type WeeklyHours struct {
	WeekStart   time.Time `json:"week_start"`
	WeekEnd     time.Time `json:"week_end"`
	TotalHours  int       `json:"total_hours"`
	DailyHours  []int     `json:"daily_hours"`
	TrendUp     bool      `json:"trend_up"`
	ChangeHours int       `json:"change_hours"`
}