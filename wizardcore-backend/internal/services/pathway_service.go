package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type PathwayService struct {
	pathwayRepo *repositories.PathwayRepository
	userRepo    *repositories.UserRepository
}

func NewPathwayService(pathwayRepo *repositories.PathwayRepository, userRepo *repositories.UserRepository) *PathwayService {
	return &PathwayService{pathwayRepo: pathwayRepo, userRepo: userRepo}
}

func (s *PathwayService) GetAllPathways() ([]models.Pathway, error) {
	return s.pathwayRepo.FindAll()
}

func (s *PathwayService) GetAllPathwaysWithEnrollment(userID uuid.UUID) ([]models.PathwayWithEnrollment, error) {
	return s.pathwayRepo.FindAllWithEnrollment(userID)
}

func (s *PathwayService) GetPathwayByID(id uuid.UUID) (*models.Pathway, error) {
	return s.pathwayRepo.FindByID(id)
}

func (s *PathwayService) GetUserEnrollments(userID uuid.UUID) ([]models.UserPathwayEnrollment, error) {
	return s.pathwayRepo.FindEnrollmentsByUserID(userID)
}

func (s *PathwayService) GetPathwayProgress(userID, pathwayID uuid.UUID) (*models.PathwayProgress, error) {
	return s.pathwayRepo.GetPathwayProgress(userID, pathwayID)
}

func (s *PathwayService) EnrollUser(userID, pathwayID uuid.UUID) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	// Check if pathway exists
	pathway, err := s.pathwayRepo.FindByID(pathwayID)
	if err != nil {
		return fmt.Errorf("failed to find pathway: %w", err)
	}
	if pathway == nil {
		return fmt.Errorf("pathway not found")
	}

	// Check if already enrolled
	existing, err := s.pathwayRepo.FindEnrollment(userID, pathwayID)
	if err != nil {
		return fmt.Errorf("failed to check enrollment: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("user already enrolled in this pathway")
	}

	// Create enrollment with default values
	enrollment := &models.UserPathwayEnrollment{
		UserID:             userID,
		PathwayID:          pathwayID,
		ProgressPercentage: 0,
		CompletedModules:   0,
		XPEarned:           0,
		StreakDays:         0,
		LastActivityAt:     nil,
		CompletedAt:        nil,
	}
	return s.pathwayRepo.CreateEnrollment(enrollment)
}

func (s *PathwayService) UpdateEnrollmentProgress(userID, pathwayID uuid.UUID, progressPercentage, completedModules, xpEarned, streakDays int) error {
	enrollment, err := s.pathwayRepo.FindEnrollment(userID, pathwayID)
	if err != nil {
		return fmt.Errorf("failed to find enrollment: %w", err)
	}
	if enrollment == nil {
		return fmt.Errorf("enrollment not found")
	}

	enrollment.ProgressPercentage = progressPercentage
	enrollment.CompletedModules = completedModules
	enrollment.XPEarned = xpEarned
	enrollment.StreakDays = streakDays
	// Update last activity to now
	now := time.Now()
	enrollment.LastActivityAt = &now

	return s.pathwayRepo.UpdateEnrollment(enrollment)
}
