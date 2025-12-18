package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type SubmissionRepository struct {
	db *sql.DB
}

func NewSubmissionRepository(db *sql.DB) *SubmissionRepository {
	return &SubmissionRepository{db: db}
}

func (r *SubmissionRepository) Create(submission *models.Submission) error {
	query := `
		INSERT INTO submissions (
			id, user_id, exercise_id, source_code, language_id,
			judge0_token, status, stdout, stderr, compile_output,
			execution_time, memory_used, test_cases_passed, test_cases_total,
			points_earned, is_correct, submission_type, ip_address, user_agent,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21)
		RETURNING created_at, updated_at
	`
	if submission.ID == uuid.Nil {
		submission.ID = uuid.New()
	}
	now := time.Now()
	err := r.db.QueryRow(
		query,
		submission.ID,
		submission.UserID,
		submission.ExerciseID,
		submission.SourceCode,
		submission.LanguageID,
		submission.Judge0Token,
		submission.Status,
		submission.Stdout,
		submission.Stderr,
		submission.CompileOutput,
		submission.ExecutionTime,
		submission.MemoryUsed,
		submission.TestCasesPassed,
		submission.TestCasesTotal,
		submission.PointsEarned,
		submission.IsCorrect,
		submission.SubmissionType,
		submission.IPAddress,
		submission.UserAgent,
		now,
		now,
	).Scan(&submission.CreatedAt, &submission.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create submission: %w", err)
	}
	return nil
}

func (r *SubmissionRepository) FindByID(id uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, source_code, language_id,
			judge0_token, status, stdout, stderr, compile_output,
			execution_time, memory_used, test_cases_passed, test_cases_total,
			points_earned, is_correct, submission_type, ip_address, user_agent,
			created_at, updated_at
		FROM submissions
		WHERE id = $1
	`
	submission := &models.Submission{}
	err := r.db.QueryRow(query, id).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ExerciseID,
		&submission.SourceCode,
		&submission.LanguageID,
		&submission.Judge0Token,
		&submission.Status,
		&submission.Stdout,
		&submission.Stderr,
		&submission.CompileOutput,
		&submission.ExecutionTime,
		&submission.MemoryUsed,
		&submission.TestCasesPassed,
		&submission.TestCasesTotal,
		&submission.PointsEarned,
		&submission.IsCorrect,
		&submission.SubmissionType,
		&submission.IPAddress,
		&submission.UserAgent,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find submission by ID: %w", err)
	}
	return submission, nil
}

func (r *SubmissionRepository) FindByExerciseIDAndUserID(exerciseID, userID uuid.UUID) ([]*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, source_code, language_id,
			judge0_token, status, stdout, stderr, compile_output,
			execution_time, memory_used, test_cases_passed, test_cases_total,
			points_earned, is_correct, submission_type, ip_address, user_agent,
			created_at, updated_at
		FROM submissions
		WHERE exercise_id = $1 AND user_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(query, exerciseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query submissions: %w", err)
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		submission := &models.Submission{}
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ExerciseID,
			&submission.SourceCode,
			&submission.LanguageID,
			&submission.Judge0Token,
			&submission.Status,
			&submission.Stdout,
			&submission.Stderr,
			&submission.CompileOutput,
			&submission.ExecutionTime,
			&submission.MemoryUsed,
			&submission.TestCasesPassed,
			&submission.TestCasesTotal,
			&submission.PointsEarned,
			&submission.IsCorrect,
			&submission.SubmissionType,
			&submission.IPAddress,
			&submission.UserAgent,
			&submission.CreatedAt,
			&submission.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan submission: %w", err)
		}
		submissions = append(submissions, submission)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return submissions, nil
}

func (r *SubmissionRepository) FindLatestByExerciseIDAndUserID(exerciseID, userID uuid.UUID) (*models.Submission, error) {
	query := `
		SELECT id, user_id, exercise_id, source_code, language_id,
			judge0_token, status, stdout, stderr, compile_output,
			execution_time, memory_used, test_cases_passed, test_cases_total,
			points_earned, is_correct, submission_type, ip_address, user_agent,
			created_at, updated_at
		FROM submissions
		WHERE exercise_id = $1 AND user_id = $2
		ORDER BY created_at DESC
		LIMIT 1
	`
	submission := &models.Submission{}
	err := r.db.QueryRow(query, exerciseID, userID).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ExerciseID,
		&submission.SourceCode,
		&submission.LanguageID,
		&submission.Judge0Token,
		&submission.Status,
		&submission.Stdout,
		&submission.Stderr,
		&submission.CompileOutput,
		&submission.ExecutionTime,
		&submission.MemoryUsed,
		&submission.TestCasesPassed,
		&submission.TestCasesTotal,
		&submission.PointsEarned,
		&submission.IsCorrect,
		&submission.SubmissionType,
		&submission.IPAddress,
		&submission.UserAgent,
		&submission.CreatedAt,
		&submission.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find latest submission: %w", err)
	}
	return submission, nil
}

func (r *SubmissionRepository) Update(submission *models.Submission) error {
	query := `
		UPDATE submissions
		SET
			judge0_token = $2,
			status = $3,
			stdout = $4,
			stderr = $5,
			compile_output = $6,
			execution_time = $7,
			memory_used = $8,
			test_cases_passed = $9,
			test_cases_total = $10,
			points_earned = $11,
			is_correct = $12,
			submission_type = $13,
			ip_address = $14,
			user_agent = $15,
			updated_at = $16
		WHERE id = $1
		RETURNING updated_at
	`
	now := time.Now()
	err := r.db.QueryRow(
		query,
		submission.ID,
		submission.Judge0Token,
		submission.Status,
		submission.Stdout,
		submission.Stderr,
		submission.CompileOutput,
		submission.ExecutionTime,
		submission.MemoryUsed,
		submission.TestCasesPassed,
		submission.TestCasesTotal,
		submission.PointsEarned,
		submission.IsCorrect,
		submission.SubmissionType,
		submission.IPAddress,
		submission.UserAgent,
		now,
	).Scan(&submission.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}
	return nil
}

func (r *SubmissionRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM submissions WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete submission: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("submission not found")
	}
	return nil
}