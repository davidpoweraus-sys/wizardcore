package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Exercise struct {
	ID                     uuid.UUID              `json:"id" db:"id"`
	ModuleID               uuid.UUID              `json:"module_id" db:"module_id"`
	Title                  string                 `json:"title" db:"title"`
	Difficulty             string                 `json:"difficulty" db:"difficulty"`
	Points                 int                    `json:"points" db:"points"`
	TimeLimitMinutes       *int                   `json:"time_limit_minutes,omitempty" db:"time_limit_minutes"`
	SortOrder              int                    `json:"sort_order" db:"sort_order"`
	Objectives             pq.StringArray         `json:"objectives" db:"objectives"`
	Content                *string                `json:"content,omitempty" db:"content"`
	Examples               map[string]interface{} `json:"examples,omitempty" db:"examples"`
	Description            *string                `json:"description,omitempty" db:"description"`
	Constraints            pq.StringArray         `json:"constraints" db:"constraints"`
	Hints                  pq.StringArray         `json:"hints" db:"hints"`
	StarterCode            *string                `json:"starter_code,omitempty" db:"starter_code"`
	SolutionCode           *string                `json:"solution_code,omitempty" db:"solution_code"`
	LanguageID             int                    `json:"language_id" db:"language_id"`
	Tags                   pq.StringArray         `json:"tags" db:"tags"`
	ConcurrentSolvers      int                    `json:"concurrent_solvers" db:"concurrent_solvers"`
	TotalSubmissions       int                    `json:"total_submissions" db:"total_submissions"`
	TotalCompletions       int                    `json:"total_completions" db:"total_completions"`
	AvgCompletionTime      *int                   `json:"average_completion_time,omitempty" db:"average_completion_time"`
	CreatedAt              time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time              `json:"updated_at" db:"updated_at"`
}

type TestCase struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ExerciseID     uuid.UUID `json:"exercise_id" db:"exercise_id"`
	Input          *string   `json:"input,omitempty" db:"input"`
	ExpectedOutput string    `json:"expected_output" db:"expected_output"`
	IsHidden       bool      `json:"is_hidden" db:"is_hidden"`
	Points         int       `json:"points" db:"points"`
	SortOrder      int       `json:"sort_order" db:"sort_order"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

type ExerciseWithTests struct {
	Exercise
	TestCases []TestCase `json:"test_cases"`
}

type ExerciseStats struct {
	ConcurrentSolvers int `json:"concurrent_solvers"`
	TotalSubmissions  int `json:"total_submissions"`
	CompletionRate    int `json:"completion_rate"`
}