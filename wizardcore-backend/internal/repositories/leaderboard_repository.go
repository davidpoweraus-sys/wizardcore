package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type LeaderboardRepository struct {
	db *sql.DB
}

func NewLeaderboardRepository(db *sql.DB) *LeaderboardRepository {
	return &LeaderboardRepository{db: db}
}

// GetLeaderboard retrieves leaderboard entries for a given timeframe and optional pathway
func (r *LeaderboardRepository) GetLeaderboard(timeframe string, pathwayID *uuid.UUID, limit, offset int) ([]models.LeaderboardEntry, error) {
	query := `
		SELECT 
			le.id,
			le.user_id,
			u.display_name as username,
			u.avatar_url,
			le.timeframe,
			le.pathway_id,
			le.rank,
			le.previous_rank,
			le.xp,
			le.streak_days,
			le.badge_count,
			le.updated_at
		FROM leaderboard_entries le
		JOIN users u ON le.user_id = u.id
		WHERE le.timeframe = $1
		AND ($2::uuid IS NULL OR le.pathway_id = $2)
		ORDER BY le.rank
		LIMIT $3 OFFSET $4
	`
	var rows *sql.Rows
	var err error
	if pathwayID == nil {
		rows, err = r.db.Query(query, timeframe, nil, limit, offset)
	} else {
		rows, err = r.db.Query(query, timeframe, pathwayID, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.LeaderboardEntry
	for rows.Next() {
		var le models.LeaderboardEntry
		var pathwayID sql.NullString
		var previousRank sql.NullInt64
		err := rows.Scan(
			&le.ID,
			&le.UserID,
			&le.Username,
			&le.AvatarURL,
			&le.Timeframe,
			&pathwayID,
			&le.Rank,
			&previousRank,
			&le.XP,
			&le.StreakDays,
			&le.BadgeCount,
			&le.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		if pathwayID.Valid {
			id, err := uuid.Parse(pathwayID.String)
			if err == nil {
				le.PathwayID = &id
			}
		}
		if previousRank.Valid {
			rank := int(previousRank.Int64)
			le.PreviousRank = &rank
		}
		// TODO: set country code, trend, is_current_user (requires additional joins)
		entries = append(entries, le)
	}
	return entries, nil
}

// GetLeaderboardStats returns statistics about the leaderboard
func (r *LeaderboardRepository) GetLeaderboardStats(timeframe string, pathwayID *uuid.UUID) (*models.LeaderboardStats, error) {
	query := `
		SELECT 
			COUNT(DISTINCT le.user_id) as total_learners,
			MAX(le.xp) as top_xp,
			u.display_name as top_username
		FROM leaderboard_entries le
		JOIN users u ON le.user_id = u.id
		WHERE le.timeframe = $1
		AND ($2::uuid IS NULL OR le.pathway_id = $2)
		GROUP BY u.display_name
		ORDER BY le.xp DESC
		LIMIT 1
	`
	var stats models.LeaderboardStats
	var topUsername sql.NullString
	var pathwayParam interface{}
	if pathwayID == nil {
		pathwayParam = nil
	} else {
		pathwayParam = pathwayID
	}
	err := r.db.QueryRow(query, timeframe, pathwayParam).Scan(
		&stats.TotalLearners,
		&stats.TopXP,
		&topUsername,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if topUsername.Valid {
		stats.TopUsername = topUsername.String
	}
	// TODO: compute current user rank and change
	stats.CurrentUserRank = 0
	stats.CurrentUserChange = 0
	stats.CountryCount = 0 // would need country data
	return &stats, nil
}

// UpdateLeaderboard recalculates leaderboard entries for a given timeframe
func (r *LeaderboardRepository) UpdateLeaderboard(timeframe string) error {
	// This is a complex operation that should be run as a cron job
	// For now, we'll just log that it's called
	return nil
}

// GetTotalCount returns total number of entries for pagination
func (r *LeaderboardRepository) GetTotalCount(timeframe string, pathwayID *uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*)
		FROM leaderboard_entries
		WHERE timeframe = $1
		AND ($2::uuid IS NULL OR pathway_id = $2)
	`
	var count int
	var pathwayParam interface{}
	if pathwayID == nil {
		pathwayParam = nil
	} else {
		pathwayParam = pathwayID
	}
	err := r.db.QueryRow(query, timeframe, pathwayParam).Scan(&count)
	return count, err
}