package repositories

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
)

type ProgressRepository struct {
	db *sql.DB
}

func NewProgressRepository(db *sql.DB) *ProgressRepository {
	return &ProgressRepository{db: db}
}

// GetUserPathwayProgress returns progress for all pathways a user is enrolled in
func (r *ProgressRepository) GetUserPathwayProgress(userID uuid.UUID) ([]models.PathwayProgress, error) {
	query := `
		SELECT 
			p.id,
			p.title,
			COALESCE(upe.progress_percentage, 0) as progress_percentage,
			COALESCE(upe.completed_modules, 0) as completed_modules,
			p.module_count as total_modules,
			COALESCE(upe.xp_earned, 0) as xp_earned,
			COALESCE(upe.streak_days, 0) as streak_days,
			upe.last_activity_at
		FROM pathways p
		LEFT JOIN user_pathway_enrollments upe ON p.id = upe.pathway_id AND upe.user_id = $1
		WHERE upe.user_id IS NOT NULL
		ORDER BY p.sort_order
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var progresses []models.PathwayProgress
	for rows.Next() {
		var pp models.PathwayProgress
		var lastActivity sql.NullTime
		err := rows.Scan(
			&pp.PathwayID,
			&pp.Title,
			&pp.Progress,
			&pp.CompletedModules,
			&pp.TotalModules,
			&pp.XPEarned,
			&pp.StreakDays,
			&lastActivity,
		)
		if err != nil {
			return nil, err
		}
		if lastActivity.Valid {
			pp.LastActivity = &lastActivity.Time
		}
		progresses = append(progresses, pp)
	}
	return progresses, nil
}

// GetUserProgressTotals returns aggregated totals for a user
func (r *ProgressRepository) GetUserProgressTotals(userID uuid.UUID) (*models.ProgressTotals, error) {
	query := `
		SELECT 
			COALESCE(SUM(upe.xp_earned), 0) as total_xp,
			COALESCE(SUM(CASE WHEN upe.last_activity_at >= NOW() - INTERVAL '7 days' THEN upe.xp_earned ELSE 0 END), 0) as xp_this_week,
			COALESCE(AVG(upe.progress_percentage), 0)::int as overall_progress,
			COALESCE(u.current_streak, 0) as current_streak,
			COALESCE(SUM(upe.completed_modules), 0) as modules_completed,
			COALESCE(SUM(p.module_count), 0) as modules_total
		FROM users u
		LEFT JOIN user_pathway_enrollments upe ON u.id = upe.user_id
		LEFT JOIN pathways p ON upe.pathway_id = p.id
		WHERE u.id = $1
		GROUP BY u.id, u.current_streak
	`
	var totals models.ProgressTotals
	err := r.db.QueryRow(query, userID).Scan(
		&totals.TotalXP,
		&totals.XPThisWeek,
		&totals.OverallProgress,
		&totals.CurrentStreak,
		&totals.ModulesCompleted,
		&totals.ModulesTotal,
	)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		// Return zero totals
		return &models.ProgressTotals{}, nil
	}
	return &totals, nil
}

// GetUserMilestones returns milestones achieved by a user
func (r *ProgressRepository) GetUserMilestones(userID uuid.UUID) ([]models.Milestone, error) {
	query := `
		SELECT id, user_id, title, description, milestone_type, xp_awarded, achieved_at
		FROM user_milestones
		WHERE user_id = $1
		ORDER BY achieved_at DESC
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []models.Milestone
	for rows.Next() {
		var m models.Milestone
		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.Title,
			&m.Description,
			&m.MilestoneType,
			&m.XPAwarded,
			&m.AchievedAt,
		)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}
	return milestones, nil
}

// GetWeeklyActivity returns weekly activity data for a user
func (r *ProgressRepository) GetWeeklyActivity(userID uuid.UUID) (*models.WeeklyActivity, error) {
	// Get last 7 days of activity
	query := `
		SELECT 
			TO_CHAR(activity_date, 'Dy') as day,
			COALESCE(SUM(exercises_completed), 0) as exercises_completed,
			COALESCE(SUM(time_spent_minutes), 0) as time_spent_minutes
		FROM user_daily_activity
		WHERE user_id = $1 AND activity_date >= CURRENT_DATE - INTERVAL '7 days'
		GROUP BY activity_date
		ORDER BY activity_date
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var weeklyData []models.DailyActivity
	totalExercises := 0
	totalMinutes := 0
	days := 0
	for rows.Next() {
		var day string
		var exercises, minutes int
		err := rows.Scan(&day, &exercises, &minutes)
		if err != nil {
			return nil, err
		}
		weeklyData = append(weeklyData, models.DailyActivity{
			Day:   day,
			Value: exercises,
			Hours: minutes / 60,
		})
		totalExercises += exercises
		totalMinutes += minutes
		days++
	}

	avgDailyTime := 0
	if days > 0 {
		avgDailyTime = totalMinutes / days
	}

	// Calculate completion rate (placeholder)
	completionRate := 0
	// Calculate streak (placeholder)
	currentStreak := 0
	// Trend percentage (placeholder)
	trendPercentage := 0

	return &models.WeeklyActivity{
		WeeklyData:      weeklyData,
		AvgDailyTime:    avgDailyTime,
		CompletionRate:  completionRate,
		CurrentStreak:   currentStreak,
		TrendPercentage: trendPercentage,
	}, nil
}

// GetWeeklyHours returns weekly hours spent
func (r *ProgressRepository) GetWeeklyHours(userID uuid.UUID) (*models.WeeklyHours, error) {
	// Get the current week's start (Monday) and end (Sunday)
	now := time.Now()
	weekStart := now.AddDate(0, 0, -int(now.Weekday())+1) // Monday
	weekEnd := weekStart.AddDate(0, 0, 6)                 // Sunday

	query := `
		SELECT 
			activity_date,
			COALESCE(SUM(time_spent_minutes), 0) as minutes
		FROM user_daily_activity
		WHERE user_id = $1 AND activity_date >= $2 AND activity_date <= $3
		GROUP BY activity_date
		ORDER BY activity_date
	`
	rows, err := r.db.Query(query, userID, weekStart, weekEnd)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	dailyHours := make([]int, 7)
	totalHours := 0
	for rows.Next() {
		var date time.Time
		var minutes int
		err := rows.Scan(&date, &minutes)
		if err != nil {
			return nil, err
		}
		dayIndex := int(date.Weekday()) - 1 // Monday = 0
		if dayIndex >= 0 && dayIndex < 7 {
			dailyHours[dayIndex] = minutes / 60
			totalHours += minutes / 60
		}
	}

	// Calculate trend (compare with previous week)
	prevWeekStart := weekStart.AddDate(0, 0, -7)
	prevWeekEnd := weekEnd.AddDate(0, 0, -7)
	var prevWeekHours int
	prevQuery := `SELECT COALESCE(SUM(time_spent_minutes), 0)/60 FROM user_daily_activity WHERE user_id = $1 AND activity_date >= $2 AND activity_date <= $3`
	err = r.db.QueryRow(prevQuery, userID, prevWeekStart, prevWeekEnd).Scan(&prevWeekHours)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	trendUp := totalHours > prevWeekHours
	changeHours := totalHours - prevWeekHours

	return &models.WeeklyHours{
		WeekStart:   weekStart,
		WeekEnd:     weekEnd,
		TotalHours:  totalHours,
		DailyHours:  dailyHours,
		TrendUp:     trendUp,
		ChangeHours: changeHours,
	}, nil
}

// RecordDailyActivity records or updates daily activity for a user
func (r *ProgressRepository) RecordDailyActivity(userID uuid.UUID, date time.Time, exercisesCompleted, xpEarned, timeSpentMinutes, submissionsCount int, streakMaintained bool) error {
	query := `
		INSERT INTO user_daily_activity (user_id, activity_date, exercises_completed, xp_earned, time_spent_minutes, submissions_count, streak_maintained)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (user_id, activity_date) DO UPDATE SET
			exercises_completed = EXCLUDED.exercises_completed,
			xp_earned = EXCLUDED.xp_earned,
			time_spent_minutes = EXCLUDED.time_spent_minutes,
			submissions_count = EXCLUDED.submissions_count,
			streak_maintained = EXCLUDED.streak_maintained
	`
	_, err := r.db.Exec(query, userID, date, exercisesCompleted, xpEarned, timeSpentMinutes, submissionsCount, streakMaintained)
	return err
}

// AddMilestone adds a new milestone for a user
func (r *ProgressRepository) AddMilestone(milestone *models.Milestone) error {
	query := `
		INSERT INTO user_milestones (id, user_id, title, description, milestone_type, xp_awarded, achieved_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query, milestone.ID, milestone.UserID, milestone.Title, milestone.Description, milestone.MilestoneType, milestone.XPAwarded, milestone.AchievedAt)
	return err
}
