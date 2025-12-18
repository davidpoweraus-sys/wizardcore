package models

import (
	"time"

	"github.com/google/uuid"
)

// ChallengeType represents a type of practice challenge
type ChallengeType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// PracticeArea represents a practice area (language/topic) with completion stats
type PracticeArea struct {
	Name          string `json:"name"`
	ExerciseCount int    `json:"exercise_count"`
	CompletedCount int   `json:"completed_count"`
	ColorGradient string `json:"color_gradient"`
}

// PracticeMatch represents a practice match (duel, speed run, etc.)
type PracticeMatch struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	MatchType        string     `json:"match_type" db:"match_type"`
	Status           string     `json:"status" db:"status"`
	ExerciseID       uuid.UUID  `json:"exercise_id" db:"exercise_id"`
	TimeLimitMinutes *int       `json:"time_limit_minutes,omitempty" db:"time_limit_minutes"`
	StartedAt        *time.Time `json:"started_at,omitempty" db:"started_at"`
	EndedAt          *time.Time `json:"ended_at,omitempty" db:"ended_at"`
	CreatedAt        *time.Time `json:"created_at,omitempty" db:"created_at"`
}

// MatchParticipant represents a participant in a practice match
type MatchParticipant struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	MatchID      uuid.UUID  `json:"match_id" db:"match_id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	SubmissionID *uuid.UUID `json:"submission_id,omitempty" db:"submission_id"`
	Score        int        `json:"score" db:"score"`
	Rank         *int       `json:"rank,omitempty" db:"rank"`
	Result       string     `json:"result" db:"result"`
	XPEarned     int        `json:"xp_earned" db:"xp_earned"`
	JoinedAt     *time.Time `json:"joined_at,omitempty" db:"joined_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty" db:"finished_at"`
}

// UserPracticeStats represents a user's practice statistics
type UserPracticeStats struct {
	UserID                    uuid.UUID `json:"user_id" db:"user_id"`
	DuelsTotal                int       `json:"duels_total" db:"duels_total"`
	DuelsWon                  int       `json:"duels_won" db:"duels_won"`
	DuelsLost                 int       `json:"duels_lost" db:"duels_lost"`
	DuelsDraw                 int       `json:"duels_draw" db:"duels_draw"`
	SpeedRunsCompleted        int       `json:"speed_runs_completed" db:"speed_runs_completed"`
	BestSpeedRunTime          *int      `json:"best_speed_run_time,omitempty" db:"best_speed_run_time"`
	RandomChallengesCompleted int       `json:"random_challenges_completed" db:"random_challenges_completed"`
	TotalPracticeXP           int       `json:"total_practice_xp" db:"total_practice_xp"`
	PracticeScore             int       `json:"practice_score" db:"practice_score"`
	PracticeRank              *int      `json:"practice_rank,omitempty" db:"practice_rank"`
	AvgCompletionTime         *int      `json:"avg_completion_time,omitempty" db:"avg_completion_time"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`
}