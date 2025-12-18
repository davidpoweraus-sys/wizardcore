package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
)

func intPtr(i int) *int {
	return &i
}

type PracticeService struct {
	matchRepo   *repositories.MatchRepository
	userRepo    *repositories.UserRepository
	exerciseRepo *repositories.ExerciseRepository
}

func NewPracticeService(matchRepo *repositories.MatchRepository, userRepo *repositories.UserRepository, exerciseRepo *repositories.ExerciseRepository) *PracticeService {
	return &PracticeService{
		matchRepo:   matchRepo,
		userRepo:    userRepo,
		exerciseRepo: exerciseRepo,
	}
}

// GetChallenges returns available challenge types
func (s *PracticeService) GetChallenges() ([]models.ChallengeType, error) {
	// Hardcoded for now
	return []models.ChallengeType{
		{ID: "duel", Name: "1v1 Duel", Description: "Compete against another learner in real-time", Icon: "⚔️"},
		{ID: "speed_run", Name: "Speed Run", Description: "Complete exercises as fast as possible", Icon: "⏱️"},
		{ID: "random", Name: "Random Challenge", Description: "Get a random exercise to solve", Icon: "🎲"},
		{ID: "endurance", Name: "Endurance", Description: "Solve as many exercises as you can in a time limit", Icon: "🏋️"},
	}, nil
}

// GetAreas returns practice areas with completion stats
func (s *PracticeService) GetAreas(userID uuid.UUID) ([]models.PracticeArea, error) {
	// Hardcoded for now
	return []models.PracticeArea{
		{Name: "Python", ExerciseCount: 42, CompletedCount: 28, ColorGradient: "from-green-400 to-cyan-400"},
		{Name: "C & Assembly", ExerciseCount: 18, CompletedCount: 5, ColorGradient: "from-blue-400 to-indigo-400"},
		{Name: "JavaScript", ExerciseCount: 25, CompletedCount: 12, ColorGradient: "from-yellow-400 to-orange-400"},
		{Name: "SQL", ExerciseCount: 15, CompletedCount: 8, ColorGradient: "from-purple-400 to-pink-400"},
		{Name: "Reverse Engineering", ExerciseCount: 10, CompletedCount: 2, ColorGradient: "from-red-400 to-rose-400"},
		{Name: "Rootkit Development", ExerciseCount: 8, CompletedCount: 0, ColorGradient: "from-gray-400 to-black"},
	}, nil
}

// GetStats returns practice stats for a user
func (s *PracticeService) GetStats(userID uuid.UUID) (*models.UserPracticeStats, error) {
	return s.matchRepo.GetUserPracticeStats(userID)
}

// GetRecentMatches returns recent matches for a user
func (s *PracticeService) GetRecentMatches(userID uuid.UUID, limit int) ([]models.PracticeMatch, error) {
	return s.matchRepo.GetRecentMatches(userID, limit)
}

// StartChallenge initiates a new challenge
func (s *PracticeService) StartChallenge(userID uuid.UUID, challengeType string) (*models.PracticeMatch, error) {
	var exerciseID uuid.UUID
	var timeLimit *int

	// Determine exercise based on challenge type
	switch challengeType {
	case "duel":
		// Try to find an existing pending duel with one participant
		existingMatch, err := s.matchRepo.FindPendingDuelWithOneParticipant()
		if err != nil {
			return nil, err
		}
		if existingMatch != nil {
			// Join existing match
			match := existingMatch
			exerciseID = match.ExerciseID
			timeLimit = match.TimeLimitMinutes
			// Add participant
			now := time.Now()
			participant := &models.MatchParticipant{
				ID:           uuid.New(),
				MatchID:      match.ID,
				UserID:       userID,
				SubmissionID: nil,
				Score:        0,
				Rank:         nil,
				Result:       "",
				XPEarned:     0,
				JoinedAt:     &now,
				FinishedAt:   nil,
			}
			err = s.matchRepo.AddParticipant(participant)
			if err != nil {
				return nil, err
			}
			// Update match status to active (both participants joined)
			match.Status = "active"
			match.StartedAt = &now
			err = s.matchRepo.UpdateMatch(match)
			if err != nil {
				return nil, err
			}
			return match, nil
		}
		// No existing match, create a new one with random exercise
		exercise, err := s.exerciseRepo.GetRandomExercise()
		if err != nil {
			return nil, err
		}
		if exercise == nil {
			return nil, fmt.Errorf("no exercises available")
		}
		exerciseID = exercise.ID
		timeLimit = intPtr(10) // 10 minutes for duel
	case "random", "speed_run", "endurance":
		// Get a random exercise
		exercise, err := s.exerciseRepo.GetRandomExercise()
		if err != nil {
			return nil, err
		}
		if exercise == nil {
			return nil, fmt.Errorf("no exercises available")
		}
		exerciseID = exercise.ID
		// Set time limit based on challenge type
		if challengeType == "speed_run" {
			timeLimit = intPtr(5) // 5 minutes for speed run
		} else if challengeType == "endurance" {
			timeLimit = intPtr(30) // 30 minutes for endurance
		} else {
			timeLimit = intPtr(10) // 10 minutes for random
		}
	default:
		return nil, fmt.Errorf("unknown challenge type: %s", challengeType)
	}

	now := time.Now()
	var startedAt *time.Time
	var status string
	if challengeType == "duel" {
		// Duel starts when both participants join
		startedAt = nil
		status = "pending"
	} else {
		// Solo challenges start immediately
		startedAt = &now
		status = "active"
	}
	match := &models.PracticeMatch{
		ID:               uuid.New(),
		MatchType:        challengeType,
		Status:           status,
		ExerciseID:       exerciseID,
		TimeLimitMinutes: timeLimit,
		StartedAt:        startedAt,
		EndedAt:          nil,
		CreatedAt:        &now,
	}
	err := s.matchRepo.CreateMatch(match)
	if err != nil {
		return nil, err
	}
	// Add participant
	participant := &models.MatchParticipant{
		ID:           uuid.New(),
		MatchID:      match.ID,
		UserID:       userID,
		SubmissionID: nil,
		Score:        0,
		Rank:         nil,
		Result:       "",
		XPEarned:     0,
		JoinedAt:     &now,
		FinishedAt:   nil,
	}
	err = s.matchRepo.AddParticipant(participant)
	if err != nil {
		return nil, err
	}
	return match, nil
}

// RecordMatchResult records the result of a completed match
func (s *PracticeService) RecordMatchResult(matchID uuid.UUID, userID uuid.UUID, score int, result string, xpEarned int, submissionID *uuid.UUID) error {
	// Fetch match
	match, err := s.matchRepo.GetMatchByID(matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return fmt.Errorf("match not found")
	}

	// Fetch participant
	participant, err := s.matchRepo.GetParticipantByMatchAndUser(matchID, userID)
	if err != nil {
		return err
	}
	if participant == nil {
		return fmt.Errorf("participant not found")
	}

	now := time.Now()
	participant.SubmissionID = submissionID
	participant.Score = score
	participant.Result = result
	participant.XPEarned = xpEarned
	participant.FinishedAt = &now

	// Update participant
	err = s.matchRepo.UpdateParticipant(participant)
	if err != nil {
		return err
	}

	// Update match status if all participants have finished
	participants, err := s.matchRepo.GetParticipantsByMatchID(matchID)
	if err != nil {
		return err
	}
	allFinished := true
	for _, p := range participants {
		if p.FinishedAt == nil {
			allFinished = false
			break
		}
	}
	if allFinished {
		match.Status = "completed"
		match.EndedAt = &now
		err = s.matchRepo.UpdateMatch(match)
		if err != nil {
			return err
		}
	}

	// Update user practice stats
	stats, err := s.matchRepo.GetUserPracticeStats(userID)
	if err != nil {
		return err
	}
	// Ensure stats.UserID is set
	stats.UserID = userID

	// Update based on match type
	switch match.MatchType {
	case "duel":
		stats.DuelsTotal++
		switch result {
		case "win":
			stats.DuelsWon++
		case "loss":
			stats.DuelsLost++
		case "draw":
			stats.DuelsDraw++
		}
	case "speed_run":
		stats.SpeedRunsCompleted++
		// Compute completion time in seconds
		var startTime time.Time
		if match.StartedAt != nil {
			startTime = *match.StartedAt
		} else if match.CreatedAt != nil {
			startTime = *match.CreatedAt
		} else {
			startTime = now // fallback
		}
		completionTime := int(now.Sub(startTime).Seconds())
		if stats.BestSpeedRunTime == nil || completionTime < *stats.BestSpeedRunTime {
			stats.BestSpeedRunTime = &completionTime
		}
		// Update average completion time (simplified: just set to latest)
		// For proper average we'd need to store total time and count, but we'll skip for now
		stats.AvgCompletionTime = &completionTime
	case "random":
		stats.RandomChallengesCompleted++
	case "endurance":
		// Endurance stats not defined; could add later
	}
	stats.TotalPracticeXP += xpEarned
	stats.PracticeScore += score

	// Save updated stats
	err = s.matchRepo.UpdateUserPracticeStats(stats)
	if err != nil {
		return err
	}

	return nil
}