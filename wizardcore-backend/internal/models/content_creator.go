package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ContentCreatorProfile represents a content creator's profile
type ContentCreatorProfile struct {
	ID                  uuid.UUID      `json:"id" db:"id"`
	UserID              uuid.UUID      `json:"user_id" db:"user_id"`
	Bio                 *string        `json:"bio,omitempty" db:"bio"`
	Specialization      pq.StringArray `json:"specialization" db:"specialization"`
	Website             *string        `json:"website,omitempty" db:"website"`
	GithubURL           *string        `json:"github_url,omitempty" db:"github_url"`
	LinkedinURL         *string        `json:"linkedin_url,omitempty" db:"linkedin_url"`
	TwitterURL          *string        `json:"twitter_url,omitempty" db:"twitter_url"`
	IsVerified          bool           `json:"is_verified" db:"is_verified"`
	VerificationDate    *time.Time     `json:"verification_date,omitempty" db:"verification_date"`
	TotalContentCreated int            `json:"total_content_created" db:"total_content_created"`
	TotalStudents       int            `json:"total_students" db:"total_students"`
	AverageRating       float64        `json:"average_rating" db:"average_rating"`
	CreatedAt           time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at" db:"updated_at"`
}

// ContentReview represents a review for content awaiting approval
type ContentReview struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ContentType   string     `json:"content_type" db:"content_type"` // 'pathway', 'module', 'exercise'
	ContentID     uuid.UUID  `json:"content_id" db:"content_id"`
	ReviewerID    *uuid.UUID `json:"reviewer_id,omitempty" db:"reviewer_id"`
	Status        string     `json:"status" db:"status"` // 'pending', 'approved', 'rejected', 'needs_revision'
	ReviewNotes   *string    `json:"review_notes,omitempty" db:"review_notes"`
	RevisionNotes *string    `json:"revision_notes,omitempty" db:"revision_notes"`
	ReviewedAt    *time.Time `json:"reviewed_at,omitempty" db:"reviewed_at"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ContentVersionHistory tracks changes to content over time
type ContentVersionHistory struct {
	ID          uuid.UUID              `json:"id" db:"id"`
	ContentType string                 `json:"content_type" db:"content_type"` // 'pathway', 'module', 'exercise'
	ContentID   uuid.UUID              `json:"content_id" db:"content_id"`
	Version     int                    `json:"version" db:"version"`
	Data        map[string]interface{} `json:"data" db:"data"`
	CreatedBy   *uuid.UUID             `json:"created_by,omitempty" db:"created_by"`
	ChangeNotes *string                `json:"change_notes,omitempty" db:"change_notes"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// ContentAnalytics tracks metrics for creator content
type ContentAnalytics struct {
	ID            uuid.UUID `json:"id" db:"id"`
	ContentType   string    `json:"content_type" db:"content_type"` // 'pathway', 'module', 'exercise'
	ContentID     uuid.UUID `json:"content_id" db:"content_id"`
	CreatorID     uuid.UUID `json:"creator_id" db:"creator_id"`
	Views         int       `json:"views" db:"views"`
	Enrollments   int       `json:"enrollments" db:"enrollments"`
	Completions   int       `json:"completions" db:"completions"`
	AverageRating float64   `json:"average_rating" db:"average_rating"`
	TotalRatings  int       `json:"total_ratings" db:"total_ratings"`
	Date          time.Time `json:"date" db:"date"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ContentRating represents a user's rating and review of content
type ContentRating struct {
	ID          uuid.UUID `json:"id" db:"id"`
	ContentType string    `json:"content_type" db:"content_type"` // 'pathway', 'module', 'exercise'
	ContentID   uuid.UUID `json:"content_id" db:"content_id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Rating      int       `json:"rating" db:"rating"` // 1-5
	Review      *string   `json:"review,omitempty" db:"review"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CreateContentCreatorProfileRequest is the request to create a creator profile
type CreateContentCreatorProfileRequest struct {
	Bio            *string  `json:"bio"`
	Specialization []string `json:"specialization"`
	Website        *string  `json:"website"`
	GithubURL      *string  `json:"github_url"`
	LinkedinURL    *string  `json:"linkedin_url"`
	TwitterURL     *string  `json:"twitter_url"`
}

// UpdateContentCreatorProfileRequest is the request to update a creator profile
type UpdateContentCreatorProfileRequest struct {
	Bio            *string  `json:"bio"`
	Specialization []string `json:"specialization"`
	Website        *string  `json:"website"`
	GithubURL      *string  `json:"github_url"`
	LinkedinURL    *string  `json:"linkedin_url"`
	TwitterURL     *string  `json:"twitter_url"`
}

// CreatorStats represents statistics for a content creator
type CreatorStats struct {
	TotalPathways    int     `json:"total_pathways"`
	TotalModules     int     `json:"total_modules"`
	TotalExercises   int     `json:"total_exercises"`
	TotalStudents    int     `json:"total_students"`
	AverageRating    float64 `json:"average_rating"`
	TotalRatings     int     `json:"total_ratings"`
	TotalViews       int     `json:"total_views"`
	TotalEnrollments int     `json:"total_enrollments"`
	TotalCompletions int     `json:"total_completions"`
	CompletionRate   float64 `json:"completion_rate"`
	PendingReviews   int     `json:"pending_reviews"`
	PublishedContent int     `json:"published_content"`
	DraftContent     int     `json:"draft_content"`
}

// CreatePathwayRequest is the request to create a new pathway
type CreatePathwayRequest struct {
	Title         string   `json:"title" validate:"required"`
	Subtitle      *string  `json:"subtitle"`
	Description   *string  `json:"description"`
	Level         string   `json:"level" validate:"required,oneof=Beginner Intermediate Advanced Expert"`
	DurationWeeks int      `json:"duration_weeks" validate:"required,min=1"`
	ColorGradient *string  `json:"color_gradient"`
	Icon          *string  `json:"icon"`
	Prerequisites []string `json:"prerequisites"`
	SortOrder     int      `json:"sort_order"`
	Status        string   `json:"status" validate:"oneof=draft published"`
}

// UpdatePathwayRequest is the request to update a pathway
type UpdatePathwayRequest struct {
	Title         *string  `json:"title"`
	Subtitle      *string  `json:"subtitle"`
	Description   *string  `json:"description"`
	Level         *string  `json:"level" validate:"omitempty,oneof=Beginner Intermediate Advanced Expert"`
	DurationWeeks *int     `json:"duration_weeks" validate:"omitempty,min=1"`
	ColorGradient *string  `json:"color_gradient"`
	Icon          *string  `json:"icon"`
	Prerequisites []string `json:"prerequisites"`
	SortOrder     *int     `json:"sort_order"`
	Status        *string  `json:"status" validate:"omitempty,oneof=draft published archived under_review"`
}

// CreateModuleRequest is the request to create a new module
type CreateModuleRequest struct {
	PathwayID      uuid.UUID `json:"pathway_id" validate:"required"`
	Title          string    `json:"title" validate:"required"`
	Description    *string   `json:"description"`
	SortOrder      int       `json:"sort_order" validate:"required"`
	EstimatedHours *int      `json:"estimated_hours" validate:"omitempty,min=1"`
	XPReward       int       `json:"xp_reward" validate:"min=0"`
	Status         string    `json:"status" validate:"oneof=draft published"`
}

// UpdateModuleRequest is the request to update a module
type UpdateModuleRequest struct {
	Title          *string `json:"title"`
	Description    *string `json:"description"`
	SortOrder      *int    `json:"sort_order"`
	EstimatedHours *int    `json:"estimated_hours" validate:"omitempty,min=1"`
	XPReward       *int    `json:"xp_reward" validate:"omitempty,min=0"`
	Status         *string `json:"status" validate:"omitempty,oneof=draft published archived under_review"`
}

// CreateExerciseRequest is the request to create a new exercise
type CreateExerciseRequest struct {
	ModuleID         uuid.UUID               `json:"module_id" validate:"required"`
	Title            string                  `json:"title" validate:"required"`
	Difficulty       string                  `json:"difficulty" validate:"required,oneof=BEGINNER INTERMEDIATE ADVANCED"`
	Points           int                     `json:"points" validate:"min=0"`
	TimeLimitMinutes *int                    `json:"time_limit_minutes" validate:"omitempty,min=1"`
	SortOrder        int                     `json:"sort_order" validate:"required"`
	Objectives       []string                `json:"objectives"`
	Content          *string                 `json:"content"` // Markdown
	Examples         map[string]interface{}  `json:"examples"`
	Description      *string                 `json:"description"`
	Constraints      []string                `json:"constraints"`
	Hints            []string                `json:"hints"`
	StarterCode      *string                 `json:"starter_code"`
	SolutionCode     *string                 `json:"solution_code"`
	LanguageID       int                     `json:"language_id" validate:"required"`
	Tags             []string                `json:"tags"`
	Status           string                  `json:"status" validate:"oneof=draft published"`
	TestCases        []CreateTestCaseRequest `json:"test_cases" validate:"required,min=1"`
}

// UpdateExerciseRequest is the request to update an exercise
type UpdateExerciseRequest struct {
	Title            *string                `json:"title"`
	Difficulty       *string                `json:"difficulty" validate:"omitempty,oneof=BEGINNER INTERMEDIATE ADVANCED"`
	Points           *int                   `json:"points" validate:"omitempty,min=0"`
	TimeLimitMinutes *int                   `json:"time_limit_minutes" validate:"omitempty,min=1"`
	SortOrder        *int                   `json:"sort_order"`
	Objectives       []string               `json:"objectives"`
	Content          *string                `json:"content"`
	Examples         map[string]interface{} `json:"examples"`
	Description      *string                `json:"description"`
	Constraints      []string               `json:"constraints"`
	Hints            []string               `json:"hints"`
	StarterCode      *string                `json:"starter_code"`
	SolutionCode     *string                `json:"solution_code"`
	LanguageID       *int                   `json:"language_id"`
	Tags             []string               `json:"tags"`
	Status           *string                `json:"status" validate:"omitempty,oneof=draft published archived under_review"`
}

// CreateTestCaseRequest is the request to create a test case
type CreateTestCaseRequest struct {
	Input          *string `json:"input"`
	ExpectedOutput string  `json:"expected_output" validate:"required"`
	IsHidden       bool    `json:"is_hidden"`
	Points         int     `json:"points" validate:"min=0"`
	SortOrder      int     `json:"sort_order" validate:"required"`
}

// SubmitContentForReviewRequest is the request to submit content for review
type SubmitContentForReviewRequest struct {
	ContentType   string    `json:"content_type" validate:"required,oneof=pathway module exercise"`
	ContentID     uuid.UUID `json:"content_id" validate:"required"`
	RevisionNotes *string   `json:"revision_notes"`
}

// ReviewContentRequest is the request to review content (admin only)
type ReviewContentRequest struct {
	ReviewID    uuid.UUID `json:"review_id" validate:"required"`
	Status      string    `json:"status" validate:"required,oneof=approved rejected needs_revision"`
	ReviewNotes *string   `json:"review_notes"`
}

// ContentWithCreator extends content models with creator information
type ContentWithCreator struct {
	Content      interface{}            `json:"content"`
	Creator      *User                  `json:"creator,omitempty"`
	Profile      *ContentCreatorProfile `json:"creator_profile,omitempty"`
	IsOwner      bool                   `json:"is_owner"`
	CanEdit      bool                   `json:"can_edit"`
	ReviewStatus *string                `json:"review_status,omitempty"`
}

// Export models for complete content export
type ExportTestCase struct {
	Input          *string `json:"input"`
	ExpectedOutput string  `json:"expected_output"`
	IsHidden       bool    `json:"is_hidden"`
	Points         int     `json:"points"`
	SortOrder      int     `json:"sort_order"`
}

type ExportExercise struct {
	Title            string                 `json:"title"`
	Difficulty       string                 `json:"difficulty"`
	Points           int                    `json:"points"`
	TimeLimitMinutes *int                   `json:"time_limit_minutes,omitempty"`
	SortOrder        int                    `json:"sort_order"`
	Objectives       []string               `json:"objectives"`
	Content          *string                `json:"content,omitempty"`
	Examples         map[string]interface{} `json:"examples,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Constraints      []string               `json:"constraints"`
	Hints            []string               `json:"hints"`
	StarterCode      *string                `json:"starter_code,omitempty"`
	SolutionCode     *string                `json:"solution_code,omitempty"`
	LanguageID       int                    `json:"language_id"`
	Tags             []string               `json:"tags"`
	TestCases        []ExportTestCase       `json:"test_cases"`
}

type ExportModule struct {
	Title          string           `json:"title"`
	Description    *string          `json:"description,omitempty"`
	SortOrder      int              `json:"sort_order"`
	EstimatedHours *int             `json:"estimated_hours,omitempty"`
	XPReward       int              `json:"xp_reward"`
	Exercises      []ExportExercise `json:"exercises"`
}

type ExportPathway struct {
	Title         string         `json:"title"`
	Subtitle      *string        `json:"subtitle,omitempty"`
	Description   *string        `json:"description,omitempty"`
	Level         string         `json:"level"`
	DurationWeeks int            `json:"duration_weeks"`
	ColorGradient *string        `json:"color_gradient,omitempty"`
	Icon          *string        `json:"icon,omitempty"`
	Prerequisites []string       `json:"prerequisites"`
	SortOrder     int            `json:"sort_order"`
	Modules       []ExportModule `json:"modules"`
}

type ImportPathwayRequest struct {
	Pathway ExportPathway `json:"pathway" validate:"required"`
	Status  string        `json:"status" validate:"oneof=draft published"`
}

type ExportResponse struct {
	Pathway  ExportPathway `json:"pathway"`
	Metadata struct {
		ExportedAt time.Time `json:"exported_at"`
		Version    string    `json:"version"`
		CreatorID  uuid.UUID `json:"creator_id"`
	} `json:"metadata"`
}
