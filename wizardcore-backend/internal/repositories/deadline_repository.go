package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type DeadlineRepository struct {
	db *sql.DB
}

func NewDeadlineRepository(db *sql.DB) *DeadlineRepository {
	return &DeadlineRepository{db: db}
}

func (r *DeadlineRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Deadline, error) {
	query := `
		SELECT id, user_id, title, description, deadline_type,
			exercise_id, pathway_id, module_id, due_date,
			completed_at, is_completed, created_at
		FROM user_deadlines
		WHERE user_id = $1
		ORDER BY due_date ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query deadlines: %w", err)
	}
	defer rows.Close()

	var deadlines []models.Deadline
	for rows.Next() {
		var d models.Deadline
		err := rows.Scan(
			&d.ID,
			&d.UserID,
			&d.Title,
			&d.Description,
			&d.DeadlineType,
			&d.ExerciseID,
			&d.PathwayID,
			&d.ModuleID,
			&d.DueDate,
			&d.CompletedAt,
			&d.IsCompleted,
			&d.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan deadline: %w", err)
		}
		deadlines = append(deadlines, d)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return deadlines, nil
}

func (r *DeadlineRepository) Create(deadline *models.Deadline) error {
	query := `
		INSERT INTO user_deadlines (
			id, user_id, title, description, deadline_type,
			exercise_id, pathway_id, module_id, due_date,
			completed_at, is_completed, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at
	`
	if deadline.ID == uuid.Nil {
		deadline.ID = uuid.New()
	}
	now := time.Now()
	err := r.db.QueryRow(
		query,
		deadline.ID,
		deadline.UserID,
		deadline.Title,
		deadline.Description,
		deadline.DeadlineType,
		deadline.ExerciseID,
		deadline.PathwayID,
		deadline.ModuleID,
		deadline.DueDate,
		deadline.CompletedAt,
		deadline.IsCompleted,
		now,
	).Scan(&deadline.ID, &deadline.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create deadline: %w", err)
	}
	return nil
}

func (r *DeadlineRepository) Update(deadline *models.Deadline) error {
	query := `
		UPDATE user_deadlines
		SET
			title = $2,
			description = $3,
			deadline_type = $4,
			exercise_id = $5,
			pathway_id = $6,
			module_id = $7,
			due_date = $8,
			completed_at = $9,
			is_completed = $10
		WHERE id = $1 AND user_id = $11
		RETURNING id
	`
	_, err := r.db.Exec(
		query,
		deadline.ID,
		deadline.Title,
		deadline.Description,
		deadline.DeadlineType,
		deadline.ExerciseID,
		deadline.PathwayID,
		deadline.ModuleID,
		deadline.DueDate,
		deadline.CompletedAt,
		deadline.IsCompleted,
		deadline.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update deadline: %w", err)
	}
	return nil
}

func (r *DeadlineRepository) Delete(userID, deadlineID uuid.UUID) error {
	query := `DELETE FROM user_deadlines WHERE user_id = $1 AND id = $2`
	result, err := r.db.Exec(query, userID, deadlineID)
	if err != nil {
		return fmt.Errorf("failed to delete deadline: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("deadline not found")
	}
	return nil
}