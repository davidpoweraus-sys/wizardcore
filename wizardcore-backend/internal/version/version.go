package version

import "os"

// Version constants
var (
	// Version is the current version of the application
	Version = "1.0.0"
	
	// BuildTime should be set during compilation via -ldflags
	BuildTime = getBuildTime()
	
	// GitCommit should be set during compilation via -ldflags
	GitCommit = getGitCommit()
	
	// Environment should be set during compilation or via config
	Environment = getEnvironment()
)

// getBuildTime returns build time from environment or default
func getBuildTime() string {
	if bt := os.Getenv("BUILD_TIMESTAMP"); bt != "" {
		return bt
	}
	return "2026-01-06T06:10:00Z"
}

// getGitCommit returns git commit from environment or default
func getGitCommit() string {
	if gc := os.Getenv("GIT_COMMIT"); gc != "" {
		return gc
	}
	return "unknown"
}

// getEnvironment returns environment from environment or default
func getEnvironment() string {
	if env := os.Getenv("ENVIRONMENT"); env != "" {
		return env
	}
	return "production"
}

// Info holds version information
type Info struct {
	Version     string `json:"version"`
	BuildTime   string `json:"build_time"`
	GitCommit   string `json:"git_commit"`
	Environment string `json:"environment"`
}

// GetInfo returns the current version information
func GetInfo() Info {
	return Info{
		Version:     Version,
		BuildTime:   BuildTime,
		GitCommit:   GitCommit,
		Environment: Environment,
	}
}