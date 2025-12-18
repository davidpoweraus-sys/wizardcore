package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type ExerciseRepository struct {
	db *sql.DB
}

func NewExerciseRepository(db *sql.DB) *ExerciseRepository {
	return &ExerciseRepository{db: db}
}

func (r *ExerciseRepository) FindByID(id uuid.UUID) (*models.Exercise, error) {
	query := `
		SELECT id, module_id, title, difficulty, points, time_limit_minutes,
		       sort_order, objectives, content, examples, description,
		       constraints, hints, starter_code, solution_code, language_id,
		       tags, concurrent_solvers, total_submissions, total_completions,
		       average_completion_time, created_at, updated_at
		FROM exercises
		WHERE id = $1
	`
	var e models.Exercise
	var timeLimit, avgCompletionTime sql.NullInt64
	var content, description, starterCode, solutionCode sql.NullString
	var examplesBytes []byte
	var objectives, constraints, hints, tags pq.StringArray
	err := r.db.QueryRow(query, id).Scan(
		&e.ID,
		&e.ModuleID,
		&e.Title,
		&e.Difficulty,
		&e.Points,
		&timeLimit,
		&e.SortOrder,
		&objectives,
		&content,
		&examplesBytes,
		&description,
		&constraints,
		&hints,
		&starterCode,
		&solutionCode,
		&e.LanguageID,
		&tags,
		&e.ConcurrentSolvers,
		&e.TotalSubmissions,
		&e.TotalCompletions,
		&avgCompletionTime,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find exercise by ID: %w", err)
	}
	if timeLimit.Valid {
		val := int(timeLimit.Int64)
		e.TimeLimitMinutes = &val
	}
	if avgCompletionTime.Valid {
		val := int(avgCompletionTime.Int64)
		e.AvgCompletionTime = &val
	}
	if content.Valid {
		e.Content = &content.String
	}
	if description.Valid {
		e.Description = &description.String
	}
	if starterCode.Valid {
		e.StarterCode = &starterCode.String
	}
	if solutionCode.Valid {
		e.SolutionCode = &solutionCode.String
	}
	e.Objectives = objectives
	e.Constraints = constraints
	e.Hints = hints
	e.Tags = tags
	// Unmarshal examples JSONB
	if len(examplesBytes) > 0 {
		var examples map[string]interface{}
		if err := json.Unmarshal(examplesBytes, &examples); err != nil {
			return nil, fmt.Errorf("failed to unmarshal examples: %w", err)
		}
		e.Examples = examples
	}
	return &e, nil
}

func (r *ExerciseRepository) FindByModuleID(moduleID uuid.UUID) ([]models.Exercise, error) {
	query := `
		SELECT id, module_id, title, difficulty, points, time_limit_minutes,
		       sort_order, objectives, content, examples, description,
		       constraints, hints, starter_code, solution_code, language_id,
		       tags, concurrent_solvers, total_submissions, total_completions,
		       average_completion_time, created_at, updated_at
		FROM exercises
		WHERE module_id = $1
		ORDER BY sort_order, title
	`
	rows, err := r.db.Query(query, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to query exercises by module: %w", err)
	}
	defer rows.Close()

	var exercises []models.Exercise
	for rows.Next() {
		var e models.Exercise
		var timeLimit, avgCompletionTime sql.NullInt64
		var content, description, starterCode, solutionCode sql.NullString
		var examplesBytes []byte
		var objectives, constraints, hints, tags pq.StringArray
		err := rows.Scan(
			&e.ID,
			&e.ModuleID,
			&e.Title,
			&e.Difficulty,
			&e.Points,
			&timeLimit,
			&e.SortOrder,
			&objectives,
			&content,
			&examplesBytes,
			&description,
			&constraints,
			&hints,
			&starterCode,
			&solutionCode,
			&e.LanguageID,
			&tags,
			&e.ConcurrentSolvers,
			&e.TotalSubmissions,
			&e.TotalCompletions,
			&avgCompletionTime,
			&e.CreatedAt,
			&e.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan exercise: %w", err)
		}
		if timeLimit.Valid {
			val := int(timeLimit.Int64)
			e.TimeLimitMinutes = &val
		}
		if avgCompletionTime.Valid {
			val := int(avgCompletionTime.Int64)
			e.AvgCompletionTime = &val
		}
		if content.Valid {
			e.Content = &content.String
		}
		if description.Valid {
			e.Description = &description.String
		}
		if starterCode.Valid {
			e.StarterCode = &starterCode.String
		}
		if solutionCode.Valid {
			e.SolutionCode = &solutionCode.String
		}
		e.Objectives = objectives
		e.Constraints = constraints
		e.Hints = hints
		e.Tags = tags
		if len(examplesBytes) > 0 {
			var examples map[string]interface{}
			if err := json.Unmarshal(examplesBytes, &examples); err != nil {
				return nil, fmt.Errorf("failed to unmarshal examples: %w", err)
			}
			e.Examples = examples
		}
		exercises = append(exercises, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating exercises rows: %w", err)
	}
	return exercises, nil
}

func (r *ExerciseRepository) FindTestCases(exerciseID uuid.UUID) ([]models.TestCase, error) {
	query := `
		SELECT id, exercise_id, input, expected_output, is_hidden, points, sort_order, created_at
		FROM test_cases
		WHERE exercise_id = $1
		ORDER BY sort_order
	`
	rows, err := r.db.Query(query, exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query test cases: %w", err)
	}
	defer rows.Close()

	var testCases []models.TestCase
	for rows.Next() {
		var tc models.TestCase
		var input sql.NullString
		err := rows.Scan(
			&tc.ID,
			&tc.ExerciseID,
			&input,
			&tc.ExpectedOutput,
			&tc.IsHidden,
			&tc.Points,
			&tc.SortOrder,
			&tc.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test case: %w", err)
		}
		if input.Valid {
			tc.Input = &input.String
		}
		testCases = append(testCases, tc)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating test cases rows: %w", err)
	}
	return testCases, nil
}

func (r *ExerciseRepository) UpdateStats(exerciseID uuid.UUID, concurrentSolvers, totalSubmissions, totalCompletions int, avgCompletionTime *int) error {
	query := `
		UPDATE exercises
		SET concurrent_solvers = $2,
		    total_submissions = $3,
		    total_completions = $4,
		    average_completion_time = $5,
		    updated_at = $6
		WHERE id = $1
	`
	now := time.Now()
	_, err := r.db.Exec(query, exerciseID, concurrentSolvers, totalSubmissions, totalCompletions, avgCompletionTime, now)
	if err != nil {
		return fmt.Errorf("failed to update exercise stats: %w", err)
	}
	return nil
}

// GetRandomExercise returns a random exercise from the database
func (r *ExerciseRepository) GetRandomExercise() (*models.Exercise, error) {
	query := `
		SELECT id, module_id, title, difficulty, points, time_limit_minutes,
		       sort_order, objectives, content, examples, description,
		       constraints, hints, starter_code, solution_code, language_id,
		       tags, concurrent_solvers, total_submissions, total_completions,
		       average_completion_time, created_at, updated_at
		FROM exercises
		ORDER BY RANDOM()
		LIMIT 1
	`
	var e models.Exercise
	var timeLimit, avgCompletionTime sql.NullInt64
	var content, description, starterCode, solutionCode sql.NullString
	var examplesBytes []byte
	var objectives, constraints, hints, tags pq.StringArray
	err := r.db.QueryRow(query).Scan(
		&e.ID,
		&e.ModuleID,
		&e.Title,
		&e.Difficulty,
		&e.Points,
		&timeLimit,
		&e.SortOrder,
		&objectives,
		&content,
		&examplesBytes,
		&description,
		&constraints,
		&hints,
		&starterCode,
		&solutionCode,
		&e.LanguageID,
		&tags,
		&e.ConcurrentSolvers,
		&e.TotalSubmissions,
		&e.TotalCompletions,
		&avgCompletionTime,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get random exercise: %w", err)
	}
	if timeLimit.Valid {
		val := int(timeLimit.Int64)
		e.TimeLimitMinutes = &val
	}
	if avgCompletionTime.Valid {
		val := int(avgCompletionTime.Int64)
		e.AvgCompletionTime = &val
	}
	if content.Valid {
		e.Content = &content.String
	}
	if description.Valid {
		e.Description = &description.String
	}
	if starterCode.Valid {
		e.StarterCode = &starterCode.String
	}
	if solutionCode.Valid {
		e.SolutionCode = &solutionCode.String
	}
	e.Objectives = objectives
	e.Constraints = constraints
	e.Hints = hints
	e.Tags = tags
	// Unmarshal examples JSONB
	if len(examplesBytes) > 0 {
		var examples map[string]interface{}
		if err := json.Unmarshal(examplesBytes, &examples); err != nil {
			return nil, fmt.Errorf("failed to unmarshal examples: %w", err)
		}
		e.Examples = examples
	}
	return &e, nil
}