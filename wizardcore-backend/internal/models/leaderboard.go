package models

import (
	"time"

	"github.com/google/uuid"
)

type LeaderboardEntry struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	UserID        uuid.UUID  `json:"user_id" db:"user_id"`
	Username      string     `json:"username" db:"username"`
	AvatarURL     *string    `json:"avatar_url,omitempty" db:"avatar_url"`
	Timeframe     string     `json:"timeframe" db:"timeframe"`
	PathwayID     *uuid.UUID `json:"pathway_id,omitempty" db:"pathway_id"`
	Rank          int        `json:"rank" db:"rank"`
	PreviousRank  *int       `json:"previous_rank,omitempty" db:"previous_rank"`
	XP            int        `json:"xp" db:"xp"`
	StreakDays    int        `json:"streak_days" db:"streak_days"`
	BadgeCount    int        `json:"badge_count" db:"badge_count"`
	CountryCode   *string    `json:"country_code,omitempty"`
	Trend         string     `json:"trend"` // up, down, same
	IsCurrentUser bool       `json:"is_current_user"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

type LeaderboardResponse struct {
	Leaderboard []LeaderboardEntry `json:"leaderboard"`
	Stats       LeaderboardStats   `json:"stats"`
	Pagination  Pagination         `json:"pagination"`
}

type LeaderboardStats struct {
	TotalLearners      int     `json:"total_learners"`
	CurrentUserRank    int     `json:"current_user_rank"`
	CurrentUserChange  int     `json:"current_user_change"`
	TopXP              int     `json:"top_xp"`
	TopUsername        string  `json:"top_username"`
	CountryCount       int     `json:"country_count"`
}

type Pagination struct {
	Total   int `json:"total"`
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}