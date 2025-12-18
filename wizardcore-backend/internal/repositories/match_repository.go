package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

// CreateMatch creates a new practice match
func (r *MatchRepository) CreateMatch(match *models.PracticeMatch) error {
	query := `
		INSERT INTO practice_matches (id, match_type, status, exercise_id, time_limit_minutes, started_at, ended_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query, match.ID, match.MatchType, match.Status, match.ExerciseID, match.TimeLimitMinutes, match.StartedAt, match.EndedAt, match.CreatedAt)
	return err
}

// GetMatchByID retrieves a match by ID
func (r *MatchRepository) GetMatchByID(id uuid.UUID) (*models.PracticeMatch, error) {
	query := `
		SELECT id, match_type, status, exercise_id, time_limit_minutes, started_at, ended_at, created_at
		FROM practice_matches
		WHERE id = $1
	`
	var match models.PracticeMatch
	err := r.db.QueryRow(query, id).Scan(
		&match.ID,
		&match.MatchType,
		&match.Status,
		&match.ExerciseID,
		&match.TimeLimitMinutes,
		&match.StartedAt,
		&match.EndedAt,
		&match.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

// UpdateMatch updates a match
func (r *MatchRepository) UpdateMatch(match *models.PracticeMatch) error {
	query := `
		UPDATE practice_matches
		SET match_type = $2, status = $3, exercise_id = $4, time_limit_minutes = $5, started_at = $6, ended_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(query, match.ID, match.MatchType, match.Status, match.ExerciseID, match.TimeLimitMinutes, match.StartedAt, match.EndedAt)
	return err
}

// AddParticipant adds a participant to a match
func (r *MatchRepository) AddParticipant(participant *models.MatchParticipant) error {
	query := `
		INSERT INTO match_participants (id, match_id, user_id, submission_id, score, rank, result, xp_earned, joined_at, finished_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(query, participant.ID, participant.MatchID, participant.UserID, participant.SubmissionID, participant.Score, participant.Rank, participant.Result, participant.XPEarned, participant.JoinedAt, participant.FinishedAt)
	return err
}

// GetParticipantsByMatchID retrieves all participants for a match
func (r *MatchRepository) GetParticipantsByMatchID(matchID uuid.UUID) ([]models.MatchParticipant, error) {
	query := `
		SELECT id, match_id, user_id, submission_id, score, rank, result, xp_earned, joined_at, finished_at
		FROM match_participants
		WHERE match_id = $1
		ORDER BY rank NULLS LAST
	`
	rows, err := r.db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.MatchParticipant
	for rows.Next() {
		var p models.MatchParticipant
		err := rows.Scan(
			&p.ID,
			&p.MatchID,
			&p.UserID,
			&p.SubmissionID,
			&p.Score,
			&p.Rank,
			&p.Result,
			&p.XPEarned,
			&p.JoinedAt,
			&p.FinishedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	return participants, nil
}

// GetParticipantByMatchAndUser retrieves a participant by match ID and user ID
func (r *MatchRepository) GetParticipantByMatchAndUser(matchID uuid.UUID, userID uuid.UUID) (*models.MatchParticipant, error) {
	query := `
		SELECT id, match_id, user_id, submission_id, score, rank, result, xp_earned, joined_at, finished_at
		FROM match_participants
		WHERE match_id = $1 AND user_id = $2
	`
	var participant models.MatchParticipant
	err := r.db.QueryRow(query, matchID, userID).Scan(
		&participant.ID,
		&participant.MatchID,
		&participant.UserID,
		&participant.SubmissionID,
		&participant.Score,
		&participant.Rank,
		&participant.Result,
		&participant.XPEarned,
		&participant.JoinedAt,
		&participant.FinishedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &participant, nil
}

// UpdateParticipant updates a participant's result, score, xp, submission_id, and finished_at
func (r *MatchRepository) UpdateParticipant(participant *models.MatchParticipant) error {
	query := `
		UPDATE match_participants
		SET submission_id = $2, score = $3, rank = $4, result = $5, xp_earned = $6, finished_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(query,
		participant.ID,
		participant.SubmissionID,
		participant.Score,
		participant.Rank,
		participant.Result,
		participant.XPEarned,
		participant.FinishedAt,
	)
	return err
}

// GetUserPracticeStats retrieves practice stats for a user
func (r *MatchRepository) GetUserPracticeStats(userID uuid.UUID) (*models.UserPracticeStats, error) {
	query := `
		SELECT duels_total, duels_won, duels_lost, duels_draw, speed_runs_completed, best_speed_run_time, random_challenges_completed, total_practice_xp, practice_score, practice_rank, avg_completion_time
		FROM user_practice_stats
		WHERE user_id = $1
	`
	var stats models.UserPracticeStats
	err := r.db.QueryRow(query, userID).Scan(
		&stats.DuelsTotal,
		&stats.DuelsWon,
		&stats.DuelsLost,
		&stats.DuelsDraw,
		&stats.SpeedRunsCompleted,
		&stats.BestSpeedRunTime,
		&stats.RandomChallengesCompleted,
		&stats.TotalPracticeXP,
		&stats.PracticeScore,
		&stats.PracticeRank,
		&stats.AvgCompletionTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return zero stats
			return &models.UserPracticeStats{}, nil
		}
		return nil, err
	}
	return &stats, nil
}

// UpdateUserPracticeStats updates or inserts practice stats for a user
func (r *MatchRepository) UpdateUserPracticeStats(stats *models.UserPracticeStats) error {
	query := `
		INSERT INTO user_practice_stats (user_id, duels_total, duels_won, duels_lost, duels_draw, speed_runs_completed, best_speed_run_time, random_challenges_completed, total_practice_xp, practice_score, practice_rank, avg_completion_time, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		ON CONFLICT (user_id) DO UPDATE SET
			duels_total = EXCLUDED.duels_total,
			duels_won = EXCLUDED.duels_won,
			duels_lost = EXCLUDED.duels_lost,
			duels_draw = EXCLUDED.duels_draw,
			speed_runs_completed = EXCLUDED.speed_runs_completed,
			best_speed_run_time = EXCLUDED.best_speed_run_time,
			random_challenges_completed = EXCLUDED.random_challenges_completed,
			total_practice_xp = EXCLUDED.total_practice_xp,
			practice_score = EXCLUDED.practice_score,
			practice_rank = EXCLUDED.practice_rank,
			avg_completion_time = EXCLUDED.avg_completion_time,
			updated_at = EXCLUDED.updated_at
	`
	_, err := r.db.Exec(query,
		stats.UserID,
		stats.DuelsTotal,
		stats.DuelsWon,
		stats.DuelsLost,
		stats.DuelsDraw,
		stats.SpeedRunsCompleted,
		stats.BestSpeedRunTime,
		stats.RandomChallengesCompleted,
		stats.TotalPracticeXP,
		stats.PracticeScore,
		stats.PracticeRank,
		stats.AvgCompletionTime,
		time.Now(),
	)
	return err
}

// FindPendingDuelWithOneParticipant finds a pending duel match that has exactly one participant
func (r *MatchRepository) FindPendingDuelWithOneParticipant() (*models.PracticeMatch, error) {
	query := `
		SELECT pm.id, pm.match_type, pm.status, pm.exercise_id, pm.time_limit_minutes, pm.started_at, pm.ended_at, pm.created_at
		FROM practice_matches pm
		JOIN match_participants mp ON pm.id = mp.match_id
		WHERE pm.match_type = 'duel'
		  AND pm.status = 'pending'
		GROUP BY pm.id
		HAVING COUNT(mp.user_id) = 1
		LIMIT 1
	`
	var match models.PracticeMatch
	err := r.db.QueryRow(query).Scan(
		&match.ID,
		&match.MatchType,
		&match.Status,
		&match.ExerciseID,
		&match.TimeLimitMinutes,
		&match.StartedAt,
		&match.EndedAt,
		&match.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &match, nil
}

// GetRecentMatches returns recent matches for a user
func (r *MatchRepository) GetRecentMatches(userID uuid.UUID, limit int) ([]models.PracticeMatch, error) {
	query := `
		SELECT pm.id, pm.match_type, pm.status, pm.exercise_id, pm.time_limit_minutes, pm.started_at, pm.ended_at, pm.created_at
		FROM practice_matches pm
		JOIN match_participants mp ON pm.id = mp.match_id
		WHERE mp.user_id = $1
		ORDER BY pm.created_at DESC
		LIMIT $2
	`
	rows, err := r.db.Query(query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []models.PracticeMatch
	for rows.Next() {
		var match models.PracticeMatch
		err := rows.Scan(
			&match.ID,
			&match.MatchType,
			&match.Status,
			&match.ExerciseID,
			&match.TimeLimitMinutes,
			&match.StartedAt,
			&match.EndedAt,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		matches = append(matches, match)
	}
	return matches, nil
}