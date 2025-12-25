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

type ActivityService struct {
	activityRepo *repositories.ActivityRepository
	progressRepo *repositories.ProgressRepository
	logger       *zap.Logger
}

func NewActivityService(activityRepo *repositories.ActivityRepository, progressRepo *repositories.ProgressRepository, logger *zap.Logger) *ActivityService {
	return &ActivityService{
		activityRepo: activityRepo,
		progressRepo: progressRepo,
		logger:       logger,
	}
}

// GetUserActivities retrieves activities for a user with pagination
func (s *ActivityService) GetUserActivities(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.UserActivity, error) {
	activities, err := s.activityRepo.GetUserActivities(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("Failed to get user activities",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		return nil, fmt.Errorf("failed to get user activities: %w", err)
	}

	return activities, nil
}

// RecordExerciseSubmission records an activity for exercise submission and updates progress
func (s *ActivityService) RecordExerciseSubmission(ctx context.Context, userID uuid.UUID, exerciseID uuid.UUID, exerciseTitle string, xpEarned int, timeSpentMinutes int) error {
	// Create activity record
	err := s.activityRepo.CreateExerciseSubmissionActivity(ctx, userID, exerciseID, exerciseTitle, xpEarned, timeSpentMinutes)
	if err != nil {
		s.logger.Error("Failed to create exercise submission activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("exercise_id", exerciseID.String()),
		)
		return fmt.Errorf("failed to create exercise submission activity: %w", err)
	}

	// Update daily activity record
	today := time.Now().UTC().Truncate(24 * time.Hour)
	err = s.progressRepo.RecordDailyActivity(userID, today, 1, xpEarned, timeSpentMinutes, 1, true)
	if err != nil {
		s.logger.Error("Failed to record daily activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
		// Don't fail the whole operation if daily activity recording fails
		s.logger.Warn("Continuing despite daily activity recording failure")
	}

	s.logger.Info("Exercise submission activity recorded",
		zap.String("user_id", userID.String()),
		zap.String("exercise_id", exerciseID.String()),
		zap.Int("xp_earned", xpEarned),
		zap.Int("time_spent_minutes", timeSpentMinutes),
	)

	return nil
}

// RecordModuleCompletion records an activity for module completion
func (s *ActivityService) RecordModuleCompletion(ctx context.Context, userID uuid.UUID, moduleID uuid.UUID, moduleTitle string, xpEarned int) error {
	err := s.activityRepo.CreateModuleCompletionActivity(ctx, userID, moduleID, moduleTitle, xpEarned)
	if err != nil {
		s.logger.Error("Failed to create module completion activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("module_id", moduleID.String()),
		)
		return fmt.Errorf("failed to create module completion activity: %w", err)
	}

	s.logger.Info("Module completion activity recorded",
		zap.String("user_id", userID.String()),
		zap.String("module_id", moduleID.String()),
		zap.Int("xp_earned", xpEarned),
	)

	return nil
}

// RecordPathwayEnrollment records an activity for pathway enrollment
func (s *ActivityService) RecordPathwayEnrollment(ctx context.Context, userID uuid.UUID, pathwayID uuid.UUID, pathwayTitle string) error {
	err := s.activityRepo.CreatePathwayEnrollmentActivity(ctx, userID, pathwayID, pathwayTitle)
	if err != nil {
		s.logger.Error("Failed to create pathway enrollment activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("pathway_id", pathwayID.String()),
		)
		return fmt.Errorf("failed to create pathway enrollment activity: %w", err)
	}

	s.logger.Info("Pathway enrollment activity recorded",
		zap.String("user_id", userID.String()),
		zap.String("pathway_id", pathwayID.String()),
	)

	return nil
}

// RecordAchievementUnlock records an activity for achievement unlock
func (s *ActivityService) RecordAchievementUnlock(ctx context.Context, userID uuid.UUID, achievementID uuid.UUID, achievementTitle string, xpEarned int) error {
	err := s.activityRepo.CreateAchievementActivity(ctx, userID, achievementID, achievementTitle, xpEarned)
	if err != nil {
		s.logger.Error("Failed to create achievement activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.String("achievement_id", achievementID.String()),
		)
		return fmt.Errorf("failed to create achievement activity: %w", err)
	}

	s.logger.Info("Achievement activity recorded",
		zap.String("user_id", userID.String()),
		zap.String("achievement_id", achievementID.String()),
		zap.Int("xp_earned", xpEarned),
	)

	return nil
}

// RecordStreakMaintenance records an activity for streak maintenance
func (s *ActivityService) RecordStreakMaintenance(ctx context.Context, userID uuid.UUID, streakDays int) error {
	err := s.activityRepo.CreateStreakActivity(ctx, userID, streakDays)
	if err != nil {
		s.logger.Error("Failed to create streak activity",
			zap.Error(err),
			zap.String("user_id", userID.String()),
			zap.Int("streak_days", streakDays),
		)
		return fmt.Errorf("failed to create streak activity: %w", err)
	}

	s.logger.Info("Streak activity recorded",
		zap.String("user_id", userID.String()),
		zap.Int("streak_days", streakDays),
	)

	return nil
}

// GetRecentActivities returns recent activities formatted for frontend
func (s *ActivityService) GetRecentActivities(ctx context.Context, userID uuid.UUID, limit int) ([]models.UserActivity, error) {
	activities, err := s.activityRepo.GetUserActivities(ctx, userID, limit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent activities: %w", err)
	}

	return activities, nil
}
