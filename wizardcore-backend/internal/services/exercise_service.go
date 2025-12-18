package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type ExerciseService struct {
	exerciseRepo *repositories.ExerciseRepository
}

func NewExerciseService(exerciseRepo *repositories.ExerciseRepository) *ExerciseService {
	return &ExerciseService{exerciseRepo: exerciseRepo}
}

func (s *ExerciseService) GetExerciseByID(id uuid.UUID) (*models.ExerciseWithTests, error) {
	exercise, err := s.exerciseRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercise: %w", err)
	}
	if exercise == nil {
		return nil, nil
	}

	testCases, err := s.exerciseRepo.FindTestCases(id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch test cases: %w", err)
	}

	return &models.ExerciseWithTests{
		Exercise:  *exercise,
		TestCases: testCases,
	}, nil
}

func (s *ExerciseService) GetExercisesByModuleID(moduleID uuid.UUID) ([]models.Exercise, error) {
	return s.exerciseRepo.FindByModuleID(moduleID)
}

func (s *ExerciseService) GetExerciseStats(exerciseID uuid.UUID) (*models.ExerciseStats, error) {
	exercise, err := s.exerciseRepo.FindByID(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch exercise: %w", err)
	}
	if exercise == nil {
		return nil, nil
	}

	// Calculate completion rate
	var completionRate int
	if exercise.TotalSubmissions > 0 {
		completionRate = (exercise.TotalCompletions * 100) / exercise.TotalSubmissions
	}

	return &models.ExerciseStats{
		ConcurrentSolvers: exercise.ConcurrentSolvers,
		TotalSubmissions:  exercise.TotalSubmissions,
		CompletionRate:    completionRate,
	}, nil
}

func (s *ExerciseService) UpdateExerciseStats(exerciseID uuid.UUID, concurrentSolvers, totalSubmissions, totalCompletions int, avgCompletionTime *int) error {
	return s.exerciseRepo.UpdateStats(exerciseID, concurrentSolvers, totalSubmissions, totalCompletions, avgCompletionTime)
}