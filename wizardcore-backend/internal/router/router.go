package router

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/wizardcore-backend/internal/config"
	"github.com/yourusername/wizardcore-backend/internal/handlers"
	"github.com/yourusername/wizardcore-backend/internal/middleware"
	"github.com/yourusername/wizardcore-backend/internal/repositories"
	"github.com/yourusername/wizardcore-backend/internal/services"
	"github.com/yourusername/wizardcore-backend/internal/websocket"
	"github.com/yourusername/wizardcore-backend/pkg/judge0"
	"github.com/yourusername/wizardcore-backend/pkg/redis"

	"go.uber.org/zap"
)

func Setup(db *sql.DB, cfg *config.Config, logger *zap.Logger, hub *websocket.Hub) *gin.Engine {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// CORS
	r.Use(middleware.CORSMiddleware(cfg.CORSAllowedOrigins))

	// Rate limiting (global)
	r.Use(middleware.RateLimitMiddleware(cfg.RateLimitRPS, cfg.RateLimitBurst))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	pathwayRepo := repositories.NewPathwayRepository(db)
	exerciseRepo := repositories.NewExerciseRepository(db)
	submissionRepo := repositories.NewSubmissionRepository(db)
	achievementRepo := repositories.NewAchievementRepository(db)
	progressRepo := repositories.NewProgressRepository(db)
	leaderboardRepo := repositories.NewLeaderboardRepository(db)
	matchRepo := repositories.NewMatchRepository(db)
	searchRepo := repositories.NewSearchRepository(db)

	// Initialize Judge0 client
	judge0Client := judge0.NewClient(cfg.Judge0APIURL, cfg.Judge0APIKey)

	// Initialize Redis client (optional)
	var redisClient *redis.Client
	if cfg.RedisURL != "" && cfg.RedisURL != "localhost:6379" {
		var err error
		redisClient, err = redis.NewClient(cfg.RedisURL)
		if err != nil {
			logger.Warn("Failed to connect to Redis, caching disabled", zap.Error(err))
		} else {
			logger.Info("Redis connected successfully")
		}
	}

	// Initialize services
	userService := services.NewUserService(userRepo)
	pathwayService := services.NewPathwayService(pathwayRepo, userRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	practiceService := services.NewPracticeService(matchRepo, userRepo, exerciseRepo)
	submissionService := services.NewSubmissionService(submissionRepo, exerciseRepo, userRepo, judge0Client, practiceService)
	achievementService := services.NewAchievementService(achievementRepo, userRepo)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo, userRepo, redisClient)
	progressService := services.NewProgressService(progressRepo, userRepo, pathwayRepo, exerciseRepo)
	searchService := services.NewSearchService(searchRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, logger)
	userHandler := handlers.NewUserHandler(userService, progressService, logger)
	pathwayHandler := handlers.NewPathwayHandler(pathwayService, userService, logger)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService, logger)
	submissionHandler := handlers.NewSubmissionHandler(submissionService, logger)
	achievementHandler := handlers.NewAchievementHandler(achievementService, logger)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService, logger)
	progressHandler := handlers.NewProgressHandler(progressService, logger)
	practiceHandler := handlers.NewPracticeHandler(practiceService, logger)
	searchHandler := handlers.NewSearchHandler(searchService, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub)

	// API routes
	api := r.Group("/api/v1")
	{
		// Public routes
		api.POST("/users", authHandler.CreateUser)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.SupabaseJWTSecret))
		{
			// User routes
			protected.GET("/users/me", userHandler.GetCurrentUser)
			protected.PUT("/users/me/profile", userHandler.UpdateProfile)
			protected.GET("/users/me/preferences", userHandler.GetPreferences)
			protected.PUT("/users/me/preferences", userHandler.UpdatePreferences)
			protected.GET("/users/me/stats", userHandler.GetStats)
			protected.GET("/users/me/activities", userHandler.GetActivities)
			protected.GET("/users/me/nav-counts", userHandler.GetNavCounts)
			protected.DELETE("/users/me", userHandler.DeleteAccount)
			protected.GET("/users/me/export", userHandler.ExportData)

			// Pathway routes
			protected.GET("/pathways", pathwayHandler.GetAllPathways)
			protected.GET("/pathways/:id", pathwayHandler.GetPathway)
			protected.POST("/pathways/:id/enroll", pathwayHandler.EnrollPathway)
			protected.GET("/users/me/pathways", pathwayHandler.GetUserPathways)
			protected.GET("/users/me/deadlines", pathwayHandler.GetDeadlines)

			// Exercise routes
			protected.GET("/exercises", exerciseHandler.GetExercisesByModule)
			protected.GET("/exercises/:id", exerciseHandler.GetExercise)
			protected.GET("/exercises/:id/stats", exerciseHandler.GetExerciseStats)

			// Submission routes
			protected.POST("/submissions", submissionHandler.CreateSubmission)
			protected.GET("/submissions/latest/:exercise_id", submissionHandler.GetLatestSubmission)
			protected.POST("/submissions/save-draft/:exercise_id", submissionHandler.SaveDraft)
			protected.GET("/submissions/:id", submissionHandler.GetSubmission)

			// Achievement routes
			protected.GET("/users/me/achievements", achievementHandler.GetUserAchievements)

			// Leaderboard routes
			protected.GET("/leaderboard", leaderboardHandler.GetLeaderboard)

			// Progress routes
			protected.GET("/users/me/progress", progressHandler.GetUserProgress)
			protected.GET("/users/me/milestones", progressHandler.GetMilestones)
			protected.GET("/users/me/activity/weekly", progressHandler.GetWeeklyActivity)
			protected.GET("/users/me/activity/weekly-hours", progressHandler.GetWeeklyHours)

			// Practice routes
			protected.GET("/practice/challenges", practiceHandler.GetChallenges)
			protected.GET("/practice/areas", practiceHandler.GetAreas)
			protected.GET("/users/me/practice/stats", practiceHandler.GetStats)
			protected.GET("/users/me/matches", practiceHandler.GetRecentMatches)
			protected.POST("/practice/challenges/:type/start", practiceHandler.StartChallenge)

			// Search route
			protected.GET("/search", searchHandler.Search)

			// WebSocket route
			protected.GET("/ws", websocketHandler.ServeWebSocket)
		}
	}

	return r
}