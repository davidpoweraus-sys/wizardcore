package services

import (
	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type AchievementService struct {
	achievementRepo *repositories.AchievementRepository
	userRepo        *repositories.UserRepository
}

func NewAchievementService(achievementRepo *repositories.AchievementRepository, userRepo *repositories.UserRepository) *AchievementService {
	return &AchievementService{
		achievementRepo: achievementRepo,
		userRepo:        userRepo,
	}
}

func (s *AchievementService) GetUserAchievements(userID uuid.UUID) ([]models.AchievementWithProgress, error) {
	return s.achievementRepo.GetUserAchievements(userID)
}

// UnlockAchievements checks and unlocks achievements based on criteria
func (s *AchievementService) UnlockAchievements(userID uuid.UUID, criteriaType string, criteriaValue int) ([]uuid.UUID, error) {
	return s.achievementRepo.CheckAndUnlock(userID, criteriaType, criteriaValue)
}

// RecordAchievementProgress updates progress for a specific achievement type
func (s *AchievementService) RecordAchievementProgress(userID uuid.UUID, criteriaType string, criteriaValue int) error {
	_, err := s.achievementRepo.CheckAndUnlock(userID, criteriaType, criteriaValue)
	return err
}