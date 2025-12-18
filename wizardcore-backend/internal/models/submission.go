package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID               uuid.UUID  `json:"id" db:"id"`
	UserID           uuid.UUID  `json:"user_id" db:"user_id"`
	ExerciseID       uuid.UUID  `json:"exercise_id" db:"exercise_id"`
	SourceCode       string     `json:"source_code" db:"source_code"`
	LanguageID       int        `json:"language_id" db:"language_id"`
	Judge0Token      *string    `json:"judge0_token,omitempty" db:"judge0_token"`
	Status           string     `json:"status" db:"status"`
	Stdout           *string    `json:"stdout,omitempty" db:"stdout"`
	Stderr           *string    `json:"stderr,omitempty" db:"stderr"`
	CompileOutput    *string    `json:"compile_output,omitempty" db:"compile_output"`
	ExecutionTime    *float64   `json:"execution_time,omitempty" db:"execution_time"`
	MemoryUsed       *int       `json:"memory_used,omitempty" db:"memory_used"`
	TestCasesPassed  int        `json:"test_cases_passed" db:"test_cases_passed"`
	TestCasesTotal   int        `json:"test_cases_total" db:"test_cases_total"`
	PointsEarned     int        `json:"points_earned" db:"points_earned"`
	IsCorrect        bool       `json:"is_correct" db:"is_correct"`
	SubmissionType   string     `json:"submission_type" db:"submission_type"`
	IPAddress        *string    `json:"ip_address,omitempty" db:"ip_address"`
	UserAgent        *string    `json:"user_agent,omitempty" db:"user_agent"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
}

type SubmissionTestResult struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	SubmissionID  uuid.UUID  `json:"submission_id" db:"submission_id"`
	TestCaseID    uuid.UUID  `json:"test_case_id" db:"test_case_id"`
	Passed        bool       `json:"passed" db:"passed"`
	ActualOutput  *string    `json:"actual_output,omitempty" db:"actual_output"`
	ExecutionTime *float64   `json:"execution_time,omitempty" db:"execution_time"`
	MemoryUsed    *int       `json:"memory_used,omitempty" db:"memory_used"`
	ErrorMessage  *string    `json:"error_message,omitempty" db:"error_message"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
}

type CreateSubmissionRequest struct {
	ExerciseID uuid.UUID  `json:"exercise_id" validate:"required"`
	SourceCode string     `json:"source_code" validate:"required"`
	LanguageID int        `json:"language_id" validate:"required"`
	MatchID    *uuid.UUID `json:"match_id,omitempty"`
}

type SubmissionResponse struct {
	Submission
	TestResults []SubmissionTestResult `json:"test_results,omitempty"`
}