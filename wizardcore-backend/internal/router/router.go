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
	creatorRepo := repositories.NewContentCreatorRepository(db)
	// rbacRepo := repositories.NewRBACRepository(db, logger) // Not currently used
	activityRepo := repositories.NewActivityRepository(db, logger)
	preferencesRepo := repositories.NewPreferencesRepository(db)

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
	userService := services.NewUserService(userRepo, preferencesRepo)
	pathwayService := services.NewPathwayService(pathwayRepo, userRepo)
	exerciseService := services.NewExerciseService(exerciseRepo)
	practiceService := services.NewPracticeService(matchRepo, userRepo, exerciseRepo)
	progressService := services.NewProgressService(progressRepo, userRepo, pathwayRepo, exerciseRepo, activityRepo, logger)
	submissionService := services.NewSubmissionService(submissionRepo, exerciseRepo, userRepo, judge0Client, practiceService, progressService)
	achievementService := services.NewAchievementService(achievementRepo, userRepo)
	leaderboardService := services.NewLeaderboardService(leaderboardRepo, userRepo, redisClient)
	searchService := services.NewSearchService(searchRepo)
	creatorService := services.NewContentCreatorService(creatorRepo, userRepo)
	activityService := services.NewActivityService(activityRepo, progressRepo, logger)
	// rbacService := services.NewRBACService(rbacRepo, userRepo, logger) // Not currently used

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userService, logger)
	userHandler := handlers.NewUserHandler(userService, progressService, activityService, logger)
	pathwayHandler := handlers.NewPathwayHandler(pathwayService, userService, logger)
	exerciseHandler := handlers.NewExerciseHandler(exerciseService, logger)
	submissionHandler := handlers.NewSubmissionHandler(submissionService, logger)
	achievementHandler := handlers.NewAchievementHandler(achievementService, logger)
	leaderboardHandler := handlers.NewLeaderboardHandler(leaderboardService, logger)
	progressHandler := handlers.NewProgressHandler(progressService, logger)
	practiceHandler := handlers.NewPracticeHandler(practiceService, logger)
	searchHandler := handlers.NewSearchHandler(searchService, logger)
	websocketHandler := handlers.NewWebSocketHandler(hub)
	creatorHandler := handlers.NewContentCreatorHandler(creatorService, logger)

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

			// Content Creator routes (requires content_creator or admin role)
			creator := protected.Group("/content-creator")
			creator.Use(middleware.ContentCreatorMiddleware(db))
			{
				// Profile
				creator.POST("/profile", creatorHandler.CreateProfile)
				creator.GET("/profile", creatorHandler.GetProfile)
				creator.PUT("/profile", creatorHandler.UpdateProfile)
				creator.GET("/stats", creatorHandler.GetStats)

				// Pathways
				creator.POST("/pathways", creatorHandler.CreatePathway)
				creator.GET("/pathways", creatorHandler.GetPathways)
				creator.PUT("/pathways/:id", creatorHandler.UpdatePathway)
				creator.DELETE("/pathways/:id", creatorHandler.DeletePathway)

				// Modules
				creator.POST("/modules", creatorHandler.CreateModule)
				creator.GET("/modules", creatorHandler.GetModules)
				creator.PUT("/modules/:id", creatorHandler.UpdateModule)
				creator.DELETE("/modules/:id", creatorHandler.DeleteModule)

				// Exercises
				creator.POST("/exercises", creatorHandler.CreateExercise)
				creator.GET("/exercises", creatorHandler.GetExercises)
				creator.GET("/exercises/:id", creatorHandler.GetExercise)

				// Reviews
				creator.POST("/reviews", creatorHandler.SubmitForReview)
				creator.GET("/reviews", creatorHandler.GetReviews)
			}

			// Admin routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminMiddleware(db))
			{
				// Content review
				admin.POST("/reviews", creatorHandler.ReviewContent)
			}
		}
	}

	return r
}
