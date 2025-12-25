package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port               int
	DatabaseURL        string
	SupabaseURL        string
	SupabaseJWTSecret  string
	Judge0APIURL       string
	Judge0APIKey       string
	RedisURL           string
	CORSAllowedOrigins []string
	Environment        string
	LogLevel           string
	RateLimitRPS       float64
	RateLimitBurst     int
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8080"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	rps, err := strconv.ParseFloat(getEnv("RATE_LIMIT_RPS", "10"), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_RPS: %w", err)
	}
	burst, err := strconv.Atoi(getEnv("RATE_LIMIT_BURST", "30"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_BURST: %w", err)
	}

	// Get DATABASE_URL or construct from individual components
	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		// Construct DATABASE_URL from individual components
		host := getEnv("DATABASE_HOST", "localhost")
		portDB := getEnv("DATABASE_PORT", "5432")
		user := getEnv("DATABASE_USER", "wizardcore")
		password := getEnv("DATABASE_PASSWORD", "wizardcore_password")
		name := getEnv("DATABASE_NAME", "wizardcore")
		sslMode := getEnv("DATABASE_SSLMODE", "disable")

		databaseURL = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
			user, password, host, portDB, name, sslMode)
	}

	cfg := &Config{
		Port:               port,
		DatabaseURL:        databaseURL,
		SupabaseURL:        getEnv("SUPABASE_URL", ""),
		SupabaseJWTSecret:  getEnv("SUPABASE_JWT_SECRET", ""),
		Judge0APIURL:       getEnv("JUDGE0_API_URL", "http://localhost:2358"),
		Judge0APIKey:       getEnv("JUDGE0_API_KEY", ""),
		RedisURL:           getEnv("REDIS_URL", "localhost:6379"),
		CORSAllowedOrigins: []string{getEnv("FRONTEND_URL", "http://localhost:3000")},
		Environment:        getEnv("ENVIRONMENT", "development"),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		RateLimitRPS:       rps,
		RateLimitBurst:     burst,
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.SupabaseJWTSecret == "" {
		return nil, fmt.Errorf("SUPABASE_JWT_SECRET is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
