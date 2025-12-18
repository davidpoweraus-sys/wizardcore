package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Pathway struct {
	ID             uuid.UUID      `json:"id" db:"id"`
	Title          string         `json:"title" db:"title"`
	Subtitle       *string        `json:"subtitle,omitempty" db:"subtitle"`
	Description    *string        `json:"description,omitempty" db:"description"`
	Level          string         `json:"level" db:"level"`
	DurationWeeks  int            `json:"duration_weeks" db:"duration_weeks"`
	StudentCount   int            `json:"student_count" db:"student_count"`
	Rating         float64        `json:"rating" db:"rating"`
	ModuleCount    int            `json:"module_count" db:"module_count"`
	ColorGradient  *string        `json:"color_gradient,omitempty" db:"color_gradient"`
	Icon           *string        `json:"icon,omitempty" db:"icon"`
	IsLocked       bool           `json:"is_locked" db:"is_locked"`
	SortOrder      int            `json:"sort_order" db:"sort_order"`
	Prerequisites  pq.StringArray `json:"prerequisites" db:"prerequisites"`
	CreatedAt      time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at" db:"updated_at"`
}

type PathwayWithEnrollment struct {
	Pathway
	IsEnrolled bool `json:"is_enrolled"`
	Progress   int  `json:"progress"`
}

type Module struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	PathwayID      uuid.UUID  `json:"pathway_id" db:"pathway_id"`
	Title          string     `json:"title" db:"title"`
	Description    *string    `json:"description,omitempty" db:"description"`
	SortOrder      int        `json:"sort_order" db:"sort_order"`
	EstimatedHours *int       `json:"estimated_hours,omitempty" db:"estimated_hours"`
	XPReward       int        `json:"xp_reward" db:"xp_reward"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at" db:"updated_at"`
}

type UserPathwayEnrollment struct {
	ID                uuid.UUID  `json:"id" db:"id"`
	UserID            uuid.UUID  `json:"user_id" db:"user_id"`
	PathwayID         uuid.UUID  `json:"pathway_id" db:"pathway_id"`
	ProgressPercentage int       `json:"progress_percentage" db:"progress_percentage"`
	CompletedModules  int        `json:"completed_modules" db:"completed_modules"`
	XPEarned          int        `json:"xp_earned" db:"xp_earned"`
	StreakDays        int        `json:"streak_days" db:"streak_days"`
	LastActivityAt    *time.Time `json:"last_activity_at,omitempty" db:"last_activity_at"`
	EnrolledAt        time.Time  `json:"enrolled_at" db:"enrolled_at"`
	CompletedAt       *time.Time `json:"completed_at,omitempty" db:"completed_at"`
}

type PathwayProgress struct {
	PathwayID         uuid.UUID  `json:"pathway_id"`
	Title             string     `json:"title"`
	Progress          int        `json:"progress_percentage"`
	CompletedModules  int        `json:"completed_modules"`
	TotalModules      int        `json:"total_modules"`
	XPEarned          int        `json:"xp_earned"`
	StreakDays        int        `json:"streak_days"`
	LastActivity      *time.Time `json:"last_activity,omitempty"`
}