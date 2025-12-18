package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
	"github.com/yourusername/wizardcore-backend/pkg/redis"
)

type LeaderboardService struct {
	leaderboardRepo *repositories.LeaderboardRepository
	userRepo        *repositories.UserRepository
	redisClient     *redis.Client
}

func NewLeaderboardService(leaderboardRepo *repositories.LeaderboardRepository, userRepo *repositories.UserRepository, redisClient *redis.Client) *LeaderboardService {
	return &LeaderboardService{
		leaderboardRepo: leaderboardRepo,
		userRepo:        userRepo,
		redisClient:     redisClient,
	}
}

func (s *LeaderboardService) GetLeaderboard(timeframe string, pathwayID *string, page, perPage int) (*models.LeaderboardResponse, error) {
	// Validate timeframe
	if timeframe == "" {
		timeframe = "all"
	}
	if !isValidTimeframe(timeframe) {
		timeframe = "all"
	}

	// Convert pathwayID string to UUID pointer
	var pathwayUUID *uuid.UUID
	if pathwayID != nil && *pathwayID != "" {
		parsed, err := uuid.Parse(*pathwayID)
		if err == nil {
			pathwayUUID = &parsed
		}
	}

	// Calculate offset
	offset := (page - 1) * perPage
	if offset < 0 {
		offset = 0
	}

	// Try cache if Redis client is available
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("leaderboard:%s:%v:%d:%d", timeframe, pathwayUUID, page, perPage)
		ctx := context.Background()
		cached, err := s.redisClient.Get(ctx, cacheKey)
		if err == nil && cached != "" {
			var response models.LeaderboardResponse
			if err := json.Unmarshal([]byte(cached), &response); err == nil {
				return &response, nil
			}
		}
	}

	// Fetch entries
	entries, err := s.leaderboardRepo.GetLeaderboard(timeframe, pathwayUUID, perPage, offset)
	if err != nil {
		return nil, err
	}

	// Fetch stats
	stats, err := s.leaderboardRepo.GetLeaderboardStats(timeframe, pathwayUUID)
	if err != nil {
		return nil, err
	}

	// Get total count for pagination
	total, err := s.leaderboardRepo.GetTotalCount(timeframe, pathwayUUID)
	if err != nil {
		return nil, err
	}

	// TODO: set trend and is_current_user for each entry
	// For now, we'll leave them empty

	response := &models.LeaderboardResponse{
		Leaderboard: entries,
		Stats:       *stats,
		Pagination: models.Pagination{
			Total:   total,
			Page:    page,
			PerPage: perPage,
		},
	}

	// Cache the response
	if s.redisClient != nil {
		cacheKey := fmt.Sprintf("leaderboard:%s:%v:%d:%d", timeframe, pathwayUUID, page, perPage)
		ctx := context.Background()
		data, err := json.Marshal(response)
		if err == nil {
			_ = s.redisClient.Set(ctx, cacheKey, data, 30*time.Second) // Cache for 30 seconds
		}
	}

	return response, nil
}

func isValidTimeframe(timeframe string) bool {
	switch timeframe {
	case "all", "month", "week":
		return true
	default:
		return false
	}
}

// UpdateLeaderboard triggers a leaderboard recalculation (to be called by cron)
func (s *LeaderboardService) UpdateLeaderboard(timeframe string) error {
	return s.leaderboardRepo.UpdateLeaderboard(timeframe)
}