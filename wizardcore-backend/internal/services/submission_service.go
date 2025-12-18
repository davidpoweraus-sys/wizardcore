package services

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
	"github.com/yourusername/wizardcore-backend/pkg/judge0"
)

type SubmissionService struct {
	submissionRepo  *repositories.SubmissionRepository
	exerciseRepo    *repositories.ExerciseRepository
	userRepo        *repositories.UserRepository
	judge0Client    *judge0.Client
	practiceService *PracticeService
}

func NewSubmissionService(submissionRepo *repositories.SubmissionRepository, exerciseRepo *repositories.ExerciseRepository, userRepo *repositories.UserRepository, judge0Client *judge0.Client, practiceService *PracticeService) *SubmissionService {
	return &SubmissionService{
		submissionRepo:  submissionRepo,
		exerciseRepo:    exerciseRepo,
		userRepo:        userRepo,
		judge0Client:    judge0Client,
		practiceService: practiceService,
	}
}

func (s *SubmissionService) CreateSubmission(submission *models.Submission) error {
	return s.CreateSubmissionWithMatch(submission, nil)
}

func (s *SubmissionService) CreateSubmissionWithMatch(submission *models.Submission, matchID *uuid.UUID) error {
	// Fetch exercise to get test cases
	exercise, err := s.exerciseRepo.FindByID(submission.ExerciseID)
	if err != nil {
		return fmt.Errorf("failed to fetch exercise: %w", err)
	}
	if exercise == nil {
		return fmt.Errorf("exercise not found")
	}

	// Fetch test cases
	testCases, err := s.exerciseRepo.FindTestCases(submission.ExerciseID)
	if err != nil {
		return fmt.Errorf("failed to fetch test cases: %w", err)
	}

	// Prepare submission fields
	submission.Status = "pending"
	submission.TestCasesTotal = len(testCases)
	submission.TestCasesPassed = 0
	submission.PointsEarned = 0
	submission.IsCorrect = false

	// Create submission record (pending)
	if err := s.submissionRepo.Create(submission); err != nil {
		return fmt.Errorf("failed to create submission record: %w", err)
	}

	// Evaluate each test case
	var totalPoints int
	for _, tc := range testCases {
		// Skip hidden test cases for now (they should be evaluated but not shown)
		if tc.IsHidden {
			continue
		}
		// Prepare Judge0 submission
		judge0Submission := judge0.Submission{
			SourceCode:     submission.SourceCode,
			LanguageID:     submission.LanguageID,
			Stdin:          "",
			ExpectedOutput: "",
		}
		if tc.Input != nil {
			judge0Submission.Stdin = *tc.Input
		}
		if tc.ExpectedOutput != "" {
			judge0Submission.ExpectedOutput = tc.ExpectedOutput
		}

		// Submit to Judge0
		result, err := s.judge0Client.Submit(judge0Submission)
		if err != nil {
			submission.Status = "judge0_error"
			errMsg := err.Error()
			submission.Stderr = &errMsg
			// Update submission with error
			s.submissionRepo.Update(submission)
			return fmt.Errorf("judge0 submission failed: %w", err)
		}

		// Determine if test passed
		passed := false
		if result.Status.ID == 3 { // Accepted
			if result.Stdout != nil && strings.TrimSpace(*result.Stdout) == strings.TrimSpace(tc.ExpectedOutput) {
				passed = true
			}
		}

		// Update submission with result (simplified)
		if passed {
			submission.TestCasesPassed++
			totalPoints += tc.Points
		}

		// Store test result (optional, we'd need a repository for submission_test_results)
	}

	// Update submission status and points
	if submission.TestCasesPassed == submission.TestCasesTotal {
		submission.Status = "accepted"
		submission.IsCorrect = true
	} else {
		submission.Status = "wrong_answer"
	}
	submission.PointsEarned = totalPoints

	// Update submission in DB
	if err := s.submissionRepo.Update(submission); err != nil {
		return fmt.Errorf("failed to update submission: %w", err)
	}

	// Update exercise stats
	s.exerciseRepo.UpdateStats(
		submission.ExerciseID,
		exercise.ConcurrentSolvers,
		exercise.TotalSubmissions+1,
		exercise.TotalCompletions,
		exercise.AvgCompletionTime,
	)

	// Update user XP (simplified)
	user, err := s.userRepo.FindByID(submission.UserID)
	if err != nil {
		// Log error but continue
		fmt.Printf("failed to fetch user for XP update: %v\n", err)
	} else if user != nil {
		user.TotalXP += submission.PointsEarned
		if err := s.userRepo.Update(user); err != nil {
			fmt.Printf("failed to update user XP: %v\n", err)
		}
	}

	// If matchID is provided, record match result
	if matchID != nil && s.practiceService != nil {
		result := "loss"
		if submission.IsCorrect {
			result = "win"
		}
		err = s.practiceService.RecordMatchResult(*matchID, submission.UserID, submission.PointsEarned, result, submission.PointsEarned, &submission.ID)
		if err != nil {
			// Log error but don't fail the submission
			fmt.Printf("failed to record match result: %v\n", err)
		}
	}

	return nil
}

func (s *SubmissionService) GetSubmissionByID(id uuid.UUID) (*models.Submission, error) {
	return s.submissionRepo.FindByID(id)
}

func (s *SubmissionService) GetLatestSubmission(exerciseID, userID uuid.UUID) (*models.Submission, error) {
	return s.submissionRepo.FindLatestByExerciseIDAndUserID(exerciseID, userID)
}

func (s *SubmissionService) SaveDraft(submission *models.Submission) error {
	submission.SubmissionType = "draft"
	submission.Status = "draft"
	return s.submissionRepo.Create(submission)
}