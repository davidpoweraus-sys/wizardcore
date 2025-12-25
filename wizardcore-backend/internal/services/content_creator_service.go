package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

type ContentCreatorService struct {
	creatorRepo *repositories.ContentCreatorRepository
	userRepo    *repositories.UserRepository
}

func NewContentCreatorService(
	creatorRepo *repositories.ContentCreatorRepository,
	userRepo *repositories.UserRepository,
) *ContentCreatorService {
	return &ContentCreatorService{
		creatorRepo: creatorRepo,
		userRepo:    userRepo,
	}
}

// Profile Management

func (s *ContentCreatorService) CreateProfile(userID uuid.UUID, req *models.CreateContentCreatorProfileRequest) (*models.ContentCreatorProfile, error) {
	// Check if user exists
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if profile already exists
	existingProfile, err := s.creatorRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing profile: %w", err)
	}
	if existingProfile != nil {
		return nil, fmt.Errorf("creator profile already exists")
	}

	// Create profile
	profile := &models.ContentCreatorProfile{
		UserID:         userID,
		Bio:            req.Bio,
		Specialization: req.Specialization,
		Website:        req.Website,
		GithubURL:      req.GithubURL,
		LinkedinURL:    req.LinkedinURL,
		TwitterURL:     req.TwitterURL,
	}

	if err := s.creatorRepo.CreateProfile(profile); err != nil {
		return nil, fmt.Errorf("failed to create profile: %w", err)
	}

	// TODO: Update user role to 'content_creator'
	// This would require adding a method to UserRepository

	return profile, nil
}

func (s *ContentCreatorService) GetProfile(userID uuid.UUID) (*models.ContentCreatorProfile, error) {
	profile, err := s.creatorRepo.GetProfileByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if profile == nil {
		return nil, fmt.Errorf("creator profile not found")
	}
	return profile, nil
}

func (s *ContentCreatorService) UpdateProfile(userID uuid.UUID, req *models.UpdateContentCreatorProfileRequest) error {
	// Verify profile exists
	_, err := s.GetProfile(userID)
	if err != nil {
		return err
	}

	updates := make(map[string]interface{})
	if req.Bio != nil {
		updates["bio"] = req.Bio
	}
	if req.Specialization != nil {
		updates["specialization"] = req.Specialization
	}
	if req.Website != nil {
		updates["website"] = req.Website
	}
	if req.GithubURL != nil {
		updates["github_url"] = req.GithubURL
	}
	if req.LinkedinURL != nil {
		updates["linkedin_url"] = req.LinkedinURL
	}
	if req.TwitterURL != nil {
		updates["twitter_url"] = req.TwitterURL
	}

	return s.creatorRepo.UpdateProfile(userID, updates)
}

func (s *ContentCreatorService) GetStats(userID uuid.UUID) (*models.CreatorStats, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, err
	}

	stats, err := s.creatorRepo.GetCreatorStats(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get creator stats: %w", err)
	}

	return stats, nil
}

// Pathway Management

func (s *ContentCreatorService) CreatePathway(userID uuid.UUID, req *models.CreatePathwayRequest) (*models.Pathway, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	pathway := &models.Pathway{
		Title:         req.Title,
		Subtitle:      req.Subtitle,
		Description:   req.Description,
		Level:         req.Level,
		DurationWeeks: req.DurationWeeks,
		ColorGradient: req.ColorGradient,
		Icon:          req.Icon,
		SortOrder:     req.SortOrder,
		Prerequisites: req.Prerequisites,
	}

	if err := s.creatorRepo.CreatePathway(pathway, userID); err != nil {
		return nil, fmt.Errorf("failed to create pathway: %w", err)
	}

	return pathway, nil
}

func (s *ContentCreatorService) GetCreatorPathways(userID uuid.UUID, status string) ([]*models.Pathway, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	pathways, err := s.creatorRepo.GetCreatorPathways(userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get creator pathways: %w", err)
	}

	return pathways, nil
}

func (s *ContentCreatorService) UpdatePathway(pathwayID, userID uuid.UUID, req *models.UpdatePathwayRequest) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("pathway", pathwayID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this pathway")
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = req.Title
	}
	if req.Subtitle != nil {
		updates["subtitle"] = req.Subtitle
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.Level != nil {
		updates["level"] = req.Level
	}
	if req.DurationWeeks != nil {
		updates["duration_weeks"] = req.DurationWeeks
	}
	if req.ColorGradient != nil {
		updates["color_gradient"] = req.ColorGradient
	}
	if req.Icon != nil {
		updates["icon"] = req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = req.Status
	}

	return s.creatorRepo.UpdatePathway(pathwayID, userID, updates)
}

func (s *ContentCreatorService) DeletePathway(pathwayID, userID uuid.UUID) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("pathway", pathwayID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this pathway")
	}

	return s.creatorRepo.DeletePathway(pathwayID, userID)
}

// Module Management

func (s *ContentCreatorService) CreateModule(userID uuid.UUID, req *models.CreateModuleRequest) (*models.Module, error) {
	// Verify user owns the pathway
	isOwner, err := s.creatorRepo.IsContentOwner("pathway", req.PathwayID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own the parent pathway")
	}

	module := &models.Module{
		PathwayID:      req.PathwayID,
		Title:          req.Title,
		Description:    req.Description,
		SortOrder:      req.SortOrder,
		EstimatedHours: req.EstimatedHours,
		XPReward:       req.XPReward,
	}

	if err := s.creatorRepo.CreateModule(module, userID); err != nil {
		return nil, fmt.Errorf("failed to create module: %w", err)
	}

	return module, nil
}

func (s *ContentCreatorService) GetCreatorModules(userID uuid.UUID, pathwayID *uuid.UUID) ([]*models.Module, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	modules, err := s.creatorRepo.GetCreatorModules(userID, pathwayID)
	if err != nil {
		return nil, fmt.Errorf("failed to get creator modules: %w", err)
	}

	return modules, nil
}

func (s *ContentCreatorService) UpdateModule(moduleID, userID uuid.UUID, req *models.UpdateModuleRequest) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("module", moduleID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this module")
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = req.Title
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = req.SortOrder
	}
	if req.EstimatedHours != nil {
		updates["estimated_hours"] = req.EstimatedHours
	}
	if req.XPReward != nil {
		updates["xp_reward"] = req.XPReward
	}
	if req.Status != nil {
		updates["status"] = req.Status
	}

	return s.creatorRepo.UpdateModule(moduleID, userID, updates)
}

func (s *ContentCreatorService) DeleteModule(moduleID, userID uuid.UUID) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("module", moduleID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this module")
	}

	return s.creatorRepo.DeleteModule(moduleID, userID)
}

// Exercise Management

func (s *ContentCreatorService) CreateExercise(userID uuid.UUID, req *models.CreateExerciseRequest) (*models.Exercise, error) {
	// Verify user owns the module
	isOwner, err := s.creatorRepo.IsContentOwner("module", req.ModuleID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own the parent module")
	}

	exercise := &models.Exercise{
		ModuleID:         req.ModuleID,
		Title:            req.Title,
		Difficulty:       req.Difficulty,
		Points:           req.Points,
		TimeLimitMinutes: req.TimeLimitMinutes,
		SortOrder:        req.SortOrder,
		Objectives:       req.Objectives,
		Content:          req.Content,
		Examples:         req.Examples,
		Description:      req.Description,
		Constraints:      req.Constraints,
		Hints:            req.Hints,
		StarterCode:      req.StarterCode,
		SolutionCode:     req.SolutionCode,
		LanguageID:       req.LanguageID,
		Tags:             req.Tags,
	}

	if err := s.creatorRepo.CreateExercise(exercise, userID); err != nil {
		return nil, fmt.Errorf("failed to create exercise: %w", err)
	}

	// Create test cases
	for _, tcReq := range req.TestCases {
		testCase := &models.TestCase{
			ExerciseID:     exercise.ID,
			Input:          tcReq.Input,
			ExpectedOutput: tcReq.ExpectedOutput,
			IsHidden:       tcReq.IsHidden,
			Points:         tcReq.Points,
			SortOrder:      tcReq.SortOrder,
		}
		if err := s.creatorRepo.CreateTestCase(testCase); err != nil {
			return nil, fmt.Errorf("failed to create test case: %w", err)
		}
	}

	return exercise, nil
}

func (s *ContentCreatorService) GetCreatorExercises(userID uuid.UUID, moduleID *uuid.UUID) ([]*models.Exercise, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	exercises, err := s.creatorRepo.GetCreatorExercises(userID, moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get creator exercises: %w", err)
	}

	return exercises, nil
}

func (s *ContentCreatorService) GetExerciseWithTestCases(exerciseID, userID uuid.UUID) (*models.ExerciseWithTests, error) {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("exercise", exerciseID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own this exercise")
	}

	// Get exercise
	exercises, err := s.creatorRepo.GetCreatorExercises(userID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get exercise: %w", err)
	}

	var exercise *models.Exercise
	for _, e := range exercises {
		if e.ID == exerciseID {
			exercise = e
			break
		}
	}
	if exercise == nil {
		return nil, fmt.Errorf("exercise not found")
	}

	// Get test cases
	testCases, err := s.creatorRepo.GetTestCasesByExercise(exerciseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}

	// Convert pointer slice to value slice
	testCaseValues := make([]models.TestCase, len(testCases))
	for i, tc := range testCases {
		testCaseValues[i] = *tc
	}

	return &models.ExerciseWithTests{
		Exercise:  *exercise,
		TestCases: testCaseValues,
	}, nil
}

func (s *ContentCreatorService) UpdateExercise(exerciseID, userID uuid.UUID, req *models.UpdateExerciseRequest) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("exercise", exerciseID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this exercise")
	}

	updates := make(map[string]interface{})
	if req.Title != nil {
		updates["title"] = req.Title
	}
	if req.Difficulty != nil {
		updates["difficulty"] = req.Difficulty
	}
	if req.Points != nil {
		updates["points"] = req.Points
	}
	if req.TimeLimitMinutes != nil {
		updates["time_limit_minutes"] = req.TimeLimitMinutes
	}
	if req.SortOrder != nil {
		updates["sort_order"] = req.SortOrder
	}
	if req.Objectives != nil {
		updates["objectives"] = req.Objectives
	}
	if req.Content != nil {
		updates["content"] = req.Content
	}
	if req.Examples != nil {
		updates["examples"] = req.Examples
	}
	if req.Description != nil {
		updates["description"] = req.Description
	}
	if req.Constraints != nil {
		updates["constraints"] = req.Constraints
	}
	if req.Hints != nil {
		updates["hints"] = req.Hints
	}
	if req.StarterCode != nil {
		updates["starter_code"] = req.StarterCode
	}
	if req.SolutionCode != nil {
		updates["solution_code"] = req.SolutionCode
	}
	if req.LanguageID != nil {
		updates["language_id"] = req.LanguageID
	}
	if req.Tags != nil {
		updates["tags"] = req.Tags
	}
	if req.Status != nil {
		updates["status"] = req.Status
	}

	return s.creatorRepo.UpdateExercise(exerciseID, userID, updates)
}

func (s *ContentCreatorService) DeleteExercise(exerciseID, userID uuid.UUID) error {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("exercise", exerciseID, userID)
	if err != nil {
		return fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return fmt.Errorf("unauthorized: user does not own this exercise")
	}

	return s.creatorRepo.DeleteExercise(exerciseID, userID)
}

// Review Management

func (s *ContentCreatorService) SubmitForReview(userID uuid.UUID, req *models.SubmitContentForReviewRequest) (*models.ContentReview, error) {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner(req.ContentType, req.ContentID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own this content")
	}

	review := &models.ContentReview{
		ContentType:   req.ContentType,
		ContentID:     req.ContentID,
		RevisionNotes: req.RevisionNotes,
	}

	if err := s.creatorRepo.CreateReview(review); err != nil {
		return nil, fmt.Errorf("failed to submit for review: %w", err)
	}

	return review, nil
}

func (s *ContentCreatorService) GetCreatorReviews(userID uuid.UUID) ([]*models.ContentReview, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	reviews, err := s.creatorRepo.GetReviewsByCreator(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get reviews: %w", err)
	}

	return reviews, nil
}

// Admin only - Review approval
func (s *ContentCreatorService) ReviewContent(adminID uuid.UUID, req *models.ReviewContentRequest) error {
	// TODO: Verify admin role
	// This would require checking the user's role in the database

	return s.creatorRepo.UpdateReview(req.ReviewID, req.Status, *req.ReviewNotes, adminID)
}

// Export/Import Operations

func (s *ContentCreatorService) ExportPathway(pathwayID, userID uuid.UUID) (*models.ExportResponse, error) {
	// Verify ownership
	isOwner, err := s.creatorRepo.IsContentOwner("pathway", pathwayID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify ownership: %w", err)
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized: user does not own this pathway")
	}

	pathway, err := s.creatorRepo.ExportPathway(pathwayID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to export pathway: %w", err)
	}

	response := &models.ExportResponse{
		Pathway: *pathway,
	}
	response.Metadata.ExportedAt = time.Now()
	response.Metadata.Version = "1.0"
	response.Metadata.CreatorID = userID

	return response, nil
}

func (s *ContentCreatorService) ImportPathway(userID uuid.UUID, req *models.ImportPathwayRequest) (*models.Pathway, error) {
	// Verify user is a content creator
	_, err := s.GetProfile(userID)
	if err != nil {
		return nil, fmt.Errorf("user is not a content creator: %w", err)
	}

	// Validate the imported pathway
	if req.Pathway.Title == "" {
		return nil, fmt.Errorf("pathway title is required")
	}
	if len(req.Pathway.Modules) == 0 {
		return nil, fmt.Errorf("pathway must have at least one module")
	}

	// Check for duplicate module sort orders
	moduleSorts := make(map[int]bool)
	for _, module := range req.Pathway.Modules {
		if module.Title == "" {
			return nil, fmt.Errorf("module title is required")
		}
		if moduleSorts[module.SortOrder] {
			return nil, fmt.Errorf("duplicate module sort order: %d", module.SortOrder)
		}
		moduleSorts[module.SortOrder] = true

		// Check for duplicate exercise sort orders within module
		exerciseSorts := make(map[int]bool)
		for _, exercise := range module.Exercises {
			if exercise.Title == "" {
				return nil, fmt.Errorf("exercise title is required")
			}
			if exercise.LanguageID <= 0 {
				return nil, fmt.Errorf("exercise language_id is required")
			}
			if len(exercise.TestCases) == 0 {
				return nil, fmt.Errorf("exercise must have at least one test case")
			}
			if exerciseSorts[exercise.SortOrder] {
				return nil, fmt.Errorf("duplicate exercise sort order: %d in module '%s'", exercise.SortOrder, module.Title)
			}
			exerciseSorts[exercise.SortOrder] = true
		}
	}

	pathway, err := s.creatorRepo.ImportPathway(userID, &req.Pathway, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to import pathway: %w", err)
	}

	return pathway, nil
}
