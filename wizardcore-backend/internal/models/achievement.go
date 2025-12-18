package models

import (
	"time"

	"github.com/google/uuid"
)

type Achievement struct {
	ID               uuid.UUID              `json:"id" db:"id"`
	Title            string                 `json:"title" db:"title"`
	Description      *string                `json:"description,omitempty" db:"description"`
	Icon             *string                `json:"icon,omitempty" db:"icon"`
	ColorGradient    *string                `json:"color_gradient,omitempty" db:"color_gradient"`
	Rarity           string                 `json:"rarity" db:"rarity"`
	XPReward         int                    `json:"xp_reward" db:"xp_reward"`
	CriteriaType     string                 `json:"criteria_type" db:"criteria_type"`
	CriteriaValue    *int                   `json:"criteria_value,omitempty" db:"criteria_value"`
	CriteriaMetadata map[string]interface{} `json:"criteria_metadata,omitempty" db:"criteria_metadata"`
	IsHidden         bool                   `json:"is_hidden" db:"is_hidden"`
	SortOrder        int                    `json:"sort_order" db:"sort_order"`
	CreatedAt        time.Time              `json:"created_at" db:"created_at"`
}

type UserAchievement struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	AchievementID uuid.UUID  `json:"achievement_id" db:"achievement_id"`
	Progress      int        `json:"progress" db:"progress"`
	EarnedAt      *time.Time `json:"earned_at,omitempty" db:"earned_at"`
}

type AchievementWithProgress struct {
	Achievement
	Earned     bool       `json:"earned"`
	Progress   int        `json:"progress"`
	EarnedDate *time.Time `json:"earned_date,omitempty"`
}