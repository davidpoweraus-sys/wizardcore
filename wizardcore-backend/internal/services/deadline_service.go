package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type DeadlineService struct {
	deadlineRepo *repositories.DeadlineRepository
}

func NewDeadlineService(deadlineRepo *repositories.DeadlineRepository) *DeadlineService {
	return &DeadlineService{deadlineRepo: deadlineRepo}
}

func (s *DeadlineService) GetUserDeadlines(userID uuid.UUID, limit, offset int) ([]models.Deadline, error) {
	if limit <= 0 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	return s.deadlineRepo.FindByUserID(userID, limit, offset)
}

func (s *DeadlineService) CreateDeadline(deadline *models.Deadline) error {
	if deadline.UserID == uuid.Nil {
		return fmt.Errorf("user_id is required")
	}
	if deadline.Title == "" {
		return fmt.Errorf("title is required")
	}
	if deadline.DueDate.IsZero() {
		return fmt.Errorf("due_date is required")
	}
	return s.deadlineRepo.Create(deadline)
}

func (s *DeadlineService) UpdateDeadline(deadline *models.Deadline) error {
	if deadline.ID == uuid.Nil {
		return fmt.Errorf("deadline ID is required")
	}
	return s.deadlineRepo.Update(deadline)
}

func (s *DeadlineService) DeleteDeadline(userID, deadlineID uuid.UUID) error {
	return s.deadlineRepo.Delete(userID, deadlineID)
}