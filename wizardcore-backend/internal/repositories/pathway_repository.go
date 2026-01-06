package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type PathwayRepository struct {
	db *sql.DB
}

func NewPathwayRepository(db *sql.DB) *PathwayRepository {
	return &PathwayRepository{db: db}
}

func (r *PathwayRepository) FindAll() ([]models.Pathway, error) {
	query := `
		SELECT id, title, subtitle, description, level, duration_weeks,
		       student_count, rating, module_count, color_gradient, icon,
		       is_locked, sort_order, prerequisites, created_at, updated_at
		FROM pathways
		ORDER BY sort_order, title
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Pathway{}, fmt.Errorf("failed to query pathways: %w", err)
	}
	defer rows.Close()

	pathways := make([]models.Pathway, 0)
	for rows.Next() {
		var p models.Pathway
		var subtitle, description, colorGradient, icon sql.NullString
		var prerequisites pq.StringArray
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&subtitle,
			&description,
			&p.Level,
			&p.DurationWeeks,
			&p.StudentCount,
			&p.Rating,
			&p.ModuleCount,
			&colorGradient,
			&icon,
			&p.IsLocked,
			&p.SortOrder,
			&prerequisites,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return []models.Pathway{}, fmt.Errorf("failed to scan pathway: %w", err)
		}
		if subtitle.Valid {
			p.Subtitle = &subtitle.String
		}
		if description.Valid {
			p.Description = &description.String
		}
		if colorGradient.Valid {
			p.ColorGradient = &colorGradient.String
		}
		if icon.Valid {
			p.Icon = &icon.String
		}
		p.Prerequisites = prerequisites
		pathways = append(pathways, p)
	}
	if err := rows.Err(); err != nil {
		return []models.Pathway{}, fmt.Errorf("error iterating pathways rows: %w", err)
	}
	return pathways, nil
}

func (r *PathwayRepository) FindAllWithEnrollment(userID uuid.UUID) ([]models.PathwayWithEnrollment, error) {
	query := `
		SELECT
			p.id, p.title, p.subtitle, p.description, p.level, p.duration_weeks,
			p.student_count, p.rating, p.module_count, p.color_gradient, p.icon,
			p.is_locked, p.sort_order, p.prerequisites, p.created_at, p.updated_at,
			CASE WHEN upe.user_id IS NOT NULL THEN true ELSE false END as is_enrolled,
			COALESCE(upe.progress_percentage, 0) as progress
		FROM pathways p
		LEFT JOIN user_pathway_enrollments upe ON p.id = upe.pathway_id AND upe.user_id = $1
		ORDER BY p.sort_order, p.title
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.PathwayWithEnrollment{}, fmt.Errorf("failed to query pathways with enrollment: %w", err)
	}
	defer rows.Close()

	pathways := make([]models.PathwayWithEnrollment, 0)
	for rows.Next() {
		var p models.PathwayWithEnrollment
		var subtitle, description, colorGradient, icon sql.NullString
		var prerequisites pq.StringArray
		err := rows.Scan(
			&p.ID,
			&p.Title,
			&subtitle,
			&description,
			&p.Level,
			&p.DurationWeeks,
			&p.StudentCount,
			&p.Rating,
			&p.ModuleCount,
			&colorGradient,
			&icon,
			&p.IsLocked,
			&p.SortOrder,
			&prerequisites,
			&p.CreatedAt,
			&p.UpdatedAt,
			&p.IsEnrolled,
			&p.Progress,
		)
		if err != nil {
			return []models.PathwayWithEnrollment{}, fmt.Errorf("failed to scan pathway with enrollment: %w", err)
		}
		if subtitle.Valid {
			p.Subtitle = &subtitle.String
		}
		if description.Valid {
			p.Description = &description.String
		}
		if colorGradient.Valid {
			p.ColorGradient = &colorGradient.String
		}
		if icon.Valid {
			p.Icon = &icon.String
		}
		p.Prerequisites = prerequisites
		pathways = append(pathways, p)
	}
	if err := rows.Err(); err != nil {
		return []models.PathwayWithEnrollment{}, fmt.Errorf("error iterating pathways rows: %w", err)
	}
	return pathways, nil
}

func (r *PathwayRepository) FindByID(id uuid.UUID) (*models.Pathway, error) {
	query := `
		SELECT id, title, subtitle, description, level, duration_weeks,
		       student_count, rating, module_count, color_gradient, icon,
		       is_locked, sort_order, prerequisites, created_at, updated_at
		FROM pathways
		WHERE id = $1
	`
	var p models.Pathway
	var subtitle, description, colorGradient, icon sql.NullString
	var prerequisites pq.StringArray
	err := r.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Title,
		&subtitle,
		&description,
		&p.Level,
		&p.DurationWeeks,
		&p.StudentCount,
		&p.Rating,
		&p.ModuleCount,
		&colorGradient,
		&icon,
		&p.IsLocked,
		&p.SortOrder,
		&prerequisites,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find pathway by ID: %w", err)
	}
	if subtitle.Valid {
		p.Subtitle = &subtitle.String
	}
	if description.Valid {
		p.Description = &description.String
	}
	if colorGradient.Valid {
		p.ColorGradient = &colorGradient.String
	}
	if icon.Valid {
		p.Icon = &icon.String
	}
	p.Prerequisites = prerequisites
	return &p, nil
}

func (r *PathwayRepository) FindEnrollmentsByUserID(userID uuid.UUID) ([]models.UserPathwayEnrollment, error) {
	query := `
		SELECT id, user_id, pathway_id, progress_percentage, completed_modules,
		       xp_earned, streak_days, last_activity_at, enrolled_at, completed_at
		FROM user_pathway_enrollments
		WHERE user_id = $1
		ORDER BY enrolled_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.UserPathwayEnrollment{}, fmt.Errorf("failed to query enrollments: %w", err)
	}
	defer rows.Close()

	enrollments := make([]models.UserPathwayEnrollment, 0)
	for rows.Next() {
		var e models.UserPathwayEnrollment
		var lastActivityAt, completedAt sql.NullTime
		err := rows.Scan(
			&e.ID,
			&e.UserID,
			&e.PathwayID,
			&e.ProgressPercentage,
			&e.CompletedModules,
			&e.XPEarned,
			&e.StreakDays,
			&lastActivityAt,
			&e.EnrolledAt,
			&completedAt,
		)
		if err != nil {
			return []models.UserPathwayEnrollment{}, fmt.Errorf("failed to scan enrollment: %w", err)
		}
		if lastActivityAt.Valid {
			e.LastActivityAt = &lastActivityAt.Time
		}
		if completedAt.Valid {
			e.CompletedAt = &completedAt.Time
		}
		enrollments = append(enrollments, e)
	}
	if err := rows.Err(); err != nil {
		return []models.UserPathwayEnrollment{}, fmt.Errorf("error iterating enrollment rows: %w", err)
	}
	return enrollments, nil
}

func (r *PathwayRepository) FindEnrollment(userID, pathwayID uuid.UUID) (*models.UserPathwayEnrollment, error) {
	query := `
		SELECT id, user_id, pathway_id, progress_percentage, completed_modules,
		       xp_earned, streak_days, last_activity_at, enrolled_at, completed_at
		FROM user_pathway_enrollments
		WHERE user_id = $1 AND pathway_id = $2
	`
	var e models.UserPathwayEnrollment
	var lastActivityAt, completedAt sql.NullTime
	err := r.db.QueryRow(query, userID, pathwayID).Scan(
		&e.ID,
		&e.UserID,
		&e.PathwayID,
		&e.ProgressPercentage,
		&e.CompletedModules,
		&e.XPEarned,
		&e.StreakDays,
		&lastActivityAt,
		&e.EnrolledAt,
		&completedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find enrollment: %w", err)
	}
	if lastActivityAt.Valid {
		e.LastActivityAt = &lastActivityAt.Time
	}
	if completedAt.Valid {
		e.CompletedAt = &completedAt.Time
	}
	return &e, nil
}

func (r *PathwayRepository) CreateEnrollment(enrollment *models.UserPathwayEnrollment) error {
	query := `
		INSERT INTO user_pathway_enrollments (
			id, user_id, pathway_id, progress_percentage, completed_modules,
			xp_earned, streak_days, last_activity_at, enrolled_at, completed_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, enrolled_at
	`
	if enrollment.ID == uuid.Nil {
		enrollment.ID = uuid.New()
	}
	now := time.Now()
	if enrollment.EnrolledAt.IsZero() {
		enrollment.EnrolledAt = now
	}
	err := r.db.QueryRow(
		query,
		enrollment.ID,
		enrollment.UserID,
		enrollment.PathwayID,
		enrollment.ProgressPercentage,
		enrollment.CompletedModules,
		enrollment.XPEarned,
		enrollment.StreakDays,
		enrollment.LastActivityAt,
		enrollment.EnrolledAt,
		enrollment.CompletedAt,
	).Scan(&enrollment.ID, &enrollment.EnrolledAt)
	if err != nil {
		return fmt.Errorf("failed to create enrollment: %w", err)
	}
	return nil
}

func (r *PathwayRepository) UpdateEnrollment(enrollment *models.UserPathwayEnrollment) error {
	query := `
		UPDATE user_pathway_enrollments
		SET progress_percentage = $2,
		    completed_modules = $3,
		    xp_earned = $4,
		    streak_days = $5,
		    last_activity_at = $6,
		    completed_at = $7
		WHERE id = $1
		RETURNING id
	`
	err := r.db.QueryRow(
		query,
		enrollment.ID,
		enrollment.ProgressPercentage,
		enrollment.CompletedModules,
		enrollment.XPEarned,
		enrollment.StreakDays,
		enrollment.LastActivityAt,
		enrollment.CompletedAt,
	).Scan(&enrollment.ID)
	if err != nil {
		return fmt.Errorf("failed to update enrollment: %w", err)
	}
	return nil
}

func (r *PathwayRepository) GetPathwayProgress(userID, pathwayID uuid.UUID) (*models.PathwayProgress, error) {
	// This query joins enrollment with pathway and modules to compute progress.
	// For simplicity, we'll just return the enrollment progress.
	// A more sophisticated query could compute total modules etc.
	query := `
		SELECT
			p.id,
			p.title,
			COALESCE(upe.progress_percentage, 0) as progress_percentage,
			COALESCE(upe.completed_modules, 0) as completed_modules,
			(SELECT COUNT(*) FROM modules m WHERE m.pathway_id = p.id) as total_modules,
			COALESCE(upe.xp_earned, 0) as xp_earned,
			COALESCE(upe.streak_days, 0) as streak_days,
			upe.last_activity_at
		FROM pathways p
		LEFT JOIN user_pathway_enrollments upe ON p.id = upe.pathway_id AND upe.user_id = $1
		WHERE p.id = $2
	`
	var progress models.PathwayProgress
	var lastActivityAt sql.NullTime
	err := r.db.QueryRow(query, userID, pathwayID).Scan(
		&progress.PathwayID,
		&progress.Title,
		&progress.Progress,
		&progress.CompletedModules,
		&progress.TotalModules,
		&progress.XPEarned,
		&progress.StreakDays,
		&lastActivityAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pathway progress: %w", err)
	}
	if lastActivityAt.Valid {
		progress.LastActivity = &lastActivityAt.Time
	}
	return &progress, nil
}
