package models

import (
	"time"

	"github.com/google/uuid"
)

type Deadline struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	UserID       uuid.UUID  `json:"user_id" db:"user_id"`
	Title        string     `json:"title" db:"title"`
	Description  *string    `json:"description,omitempty" db:"description"`
	DeadlineType string     `json:"deadline_type" db:"deadline_type"`
	ExerciseID   *uuid.UUID `json:"exercise_id,omitempty" db:"exercise_id"`
	PathwayID    *uuid.UUID `json:"pathway_id,omitempty" db:"pathway_id"`
	ModuleID     *uuid.UUID `json:"module_id,omitempty" db:"module_id"`
	DueDate      time.Time  `json:"due_date" db:"due_date"`
	CompletedAt  *time.Time `json:"completed_at,omitempty" db:"completed_at"`
	IsCompleted  bool       `json:"is_completed" db:"is_completed"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
}