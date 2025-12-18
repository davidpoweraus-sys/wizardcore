package services

import (
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type ProgressService struct {
	progressRepo *repositories.ProgressRepository
	userRepo     *repositories.UserRepository
	pathwayRepo  *repositories.PathwayRepository
	exerciseRepo *repositories.ExerciseRepository
}

func NewProgressService(progressRepo *repositories.ProgressRepository, userRepo *repositories.UserRepository, pathwayRepo *repositories.PathwayRepository, exerciseRepo *repositories.ExerciseRepository) *ProgressService {
	return &ProgressService{
		progressRepo: progressRepo,
		userRepo:     userRepo,
		pathwayRepo:  pathwayRepo,
		exerciseRepo: exerciseRepo,
	}
}

func (s *ProgressService) GetUserProgress(userID uuid.UUID) (*models.ProgressResponse, error) {
	pathways, err := s.progressRepo.GetUserPathwayProgress(userID)
	if err != nil {
		return nil, err
	}
	totals, err := s.progressRepo.GetUserProgressTotals(userID)
	if err != nil {
		return nil, err
	}
	return &models.ProgressResponse{
		Pathways: pathways,
		Totals:   *totals,
	}, nil
}

func (s *ProgressService) GetMilestones(userID uuid.UUID) ([]models.Milestone, error) {
	return s.progressRepo.GetUserMilestones(userID)
}

func (s *ProgressService) GetWeeklyActivity(userID uuid.UUID) (*models.WeeklyActivity, error) {
	return s.progressRepo.GetWeeklyActivity(userID)
}

func (s *ProgressService) GetWeeklyHours(userID uuid.UUID) (*models.WeeklyHours, error) {
	return s.progressRepo.GetWeeklyHours(userID)
}

// RecordSubmissionActivity records activity from a submission (to be called by submission service)
func (s *ProgressService) RecordSubmissionActivity(userID uuid.UUID, exerciseID uuid.UUID, xpEarned int, timeSpentMinutes int) error {
	// TODO: implement logic to update user_module_progress, user_daily_activity, and milestones
	// For now, just record daily activity
	// We need to get the current date
	// This is a placeholder
	return nil
}