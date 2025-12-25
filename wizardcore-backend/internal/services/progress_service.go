package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
	"go.uber.org/zap"
)

type ProgressService struct {
	progressRepo *repositories.ProgressRepository
	userRepo     *repositories.UserRepository
	pathwayRepo  *repositories.PathwayRepository
	exerciseRepo *repositories.ExerciseRepository
	activityRepo *repositories.ActivityRepository
	logger       *zap.Logger
}

func NewProgressService(progressRepo *repositories.ProgressRepository, userRepo *repositories.UserRepository, pathwayRepo *repositories.PathwayRepository, exerciseRepo *repositories.ExerciseRepository, activityRepo *repositories.ActivityRepository, logger *zap.Logger) *ProgressService {
	return &ProgressService{
		progressRepo: progressRepo,
		userRepo:     userRepo,
		pathwayRepo:  pathwayRepo,
		exerciseRepo: exerciseRepo,
		activityRepo: activityRepo,
		logger:       logger,
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
	ctx := context.Background()

	// Get exercise details for activity title
	exercise, err := s.exerciseRepo.FindByID(exerciseID)
	if err != nil {
		s.logger.Error("Failed to get exercise details",
			zap.Error(err),
			zap.String("exercise_id", exerciseID.String()),
		)
		// Continue without exercise title
		exercise = &models.Exercise{Title: "Unknown Exercise"}
	}

	// Record daily activity
	today := time.Now().UTC().Truncate(24 * time.Hour)
	err = s.progressRepo.RecordDailyActivity(userID, today, 1, xpEarned, timeSpentMinutes, 1, true)
	if err != nil {
		s.logger.Error("Failed to record daily activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return fmt.Errorf("failed to record daily activity: %w", err)
	}

	// Create activity record
	err = s.activityRepo.CreateExerciseSubmissionActivity(ctx, userID, exerciseID, exercise.Title, xpEarned, timeSpentMinutes)
	if err != nil {
		s.logger.Error("Failed to create exercise submission activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("exercise_id", exerciseID.String()),
		)
		// Don't fail the whole operation if activity recording fails
		s.logger.Warn("Continuing despite activity recording failure")
	}

	// Update user's total XP
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		s.logger.Error("Failed to get user for XP update",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		// Continue without XP update
	} else {
		user.TotalXP += xpEarned
		err = s.userRepo.Update(user)
		if err != nil {
			s.logger.Error("Failed to update user XP",
				zap.Error(err),
				zap.String("user_id", userID.String()),
			)
			// Continue without XP update
		}
	}

	// TODO: Update user_module_progress if exercise is part of a module
	// TODO: Check and award milestones

	s.logger.Info("Submission activity recorded",
		zap.String("user_id", userID.String()),
		zap.String("exercise_id", exerciseID.String()),
		zap.Int("xp_earned", xpEarned),
		zap.Int("time_spent_minutes", timeSpentMinutes),
	)

	return nil
}
