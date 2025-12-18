package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	// Validate required fields
	if user.SupabaseUserID == uuid.Nil {
		return fmt.Errorf("supabase_user_id is required")
	}
	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	// Check if user already exists with same supabase_user_id or email
	existing, err := s.userRepo.FindBySupabaseUserID(user.SupabaseUserID)
	if err != nil {
		return fmt.Errorf("failed to check existing user: %w", err)
	}
	if existing != nil {
		return fmt.Errorf("user with supabase_user_id %s already exists", user.SupabaseUserID)
	}
	// TODO: check email uniqueness (optional)

	// Set default values
	if user.TotalXP < 0 {
		user.TotalXP = 0
	}
	if user.PracticeScore < 0 {
		user.PracticeScore = 0
	}
	if user.CurrentStreak < 0 {
		user.CurrentStreak = 0
	}
	if user.LongestStreak < 0 {
		user.LongestStreak = 0
	}

	// Create user
	return s.userRepo.Create(user)
}

func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) GetUserBySupabaseUserID(supabaseUserID uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindBySupabaseUserID(supabaseUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user by Supabase user ID: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}

func (s *UserService) UpdateUser(user *models.User) error {
	// Ensure user exists
	existing, err := s.userRepo.FindByID(user.ID)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("user not found")
	}
	// Perform update
	return s.userRepo.Update(user)
}

func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.Delete(id)
}