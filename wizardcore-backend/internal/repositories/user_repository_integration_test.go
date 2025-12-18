package repositories

import (
	"testing"

	"github.com/google/uuid"
	"github.com/yourusername/wizardcore-backend/internal/models"
	"github.com/yourusername/wizardcore-backend/internal/testutils"
)

func TestUserRepository_Create(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	displayName := "Test User"
	user := &models.User{
		SupabaseUserID:   uuid.New(),
		Email:            "test@example.com",
		DisplayName:      &displayName,
		AvatarURL:        nil,
		Bio:              nil,
		Location:         nil,
		Website:          nil,
		GithubUsername:   nil,
		TwitterUsername:  nil,
		TotalXP:          0,
		PracticeScore:    0,
		GlobalRank:       nil,
		CurrentStreak:    0,
		LongestStreak:    0,
		LastActivityDate: nil,
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if user.ID == uuid.Nil {
		t.Error("Expected user ID to be set")
	}
	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
	if user.UpdatedAt.IsZero() {
		t.Error("Expected UpdatedAt to be set")
	}

	// Fetch by ID
	fetched, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("Expected to find user by ID")
	}
	if fetched.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, fetched.Email)
	}
	if fetched.SupabaseUserID != user.SupabaseUserID {
		t.Errorf("Expected supabase_user_id %s, got %s", user.SupabaseUserID, fetched.SupabaseUserID)
	}
}

func TestUserRepository_FindBySupabaseUserID(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	supabaseID := uuid.New()
	displayName := "Find User"
	user := &models.User{
		SupabaseUserID:   supabaseID,
		Email:            "find@example.com",
		DisplayName:      &displayName,
		AvatarURL:        nil,
		Bio:              nil,
		Location:         nil,
		Website:          nil,
		GithubUsername:   nil,
		TwitterUsername:  nil,
		TotalXP:          0,
		PracticeScore:    0,
		GlobalRank:       nil,
		CurrentStreak:    0,
		LongestStreak:    0,
		LastActivityDate: nil,
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	fetched, err := repo.FindBySupabaseUserID(supabaseID)
	if err != nil {
		t.Fatalf("FindBySupabaseUserID failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("Expected to find user by Supabase user ID")
	}
	if fetched.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, fetched.Email)
	}

	// Non-existent ID
	nonExistent := uuid.New()
	fetched, err = repo.FindBySupabaseUserID(nonExistent)
	if err != nil {
		t.Fatalf("FindBySupabaseUserID with non-existent ID should not error: %v", err)
	}
	if fetched != nil {
		t.Error("Expected nil for non-existent Supabase user ID")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	originalName := "Original Name"
	user := &models.User{
		SupabaseUserID:   uuid.New(),
		Email:            "update@example.com",
		DisplayName:      &originalName,
		AvatarURL:        nil,
		Bio:              nil,
		Location:         nil,
		Website:          nil,
		GithubUsername:   nil,
		TwitterUsername:  nil,
		TotalXP:          0,
		PracticeScore:    0,
		GlobalRank:       nil,
		CurrentStreak:    0,
		LongestStreak:    0,
		LastActivityDate: nil,
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Update fields
	updatedName := "Updated Name"
	bio := "A bio"
	location := "Earth"
	user.DisplayName = &updatedName
	user.Bio = &bio
	user.Location = &location
	user.TotalXP = 100
	user.CurrentStreak = 5

	err = repo.Update(user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Fetch and verify
	fetched, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}
	if fetched.DisplayName == nil || *fetched.DisplayName != "Updated Name" {
		t.Errorf("Expected DisplayName 'Updated Name', got %v", fetched.DisplayName)
	}
	if fetched.Bio == nil || *fetched.Bio != "A bio" {
		t.Errorf("Expected Bio 'A bio', got %v", fetched.Bio)
	}
	if fetched.TotalXP != 100 {
		t.Errorf("Expected TotalXP 100, got %d", fetched.TotalXP)
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := testutils.SetupTestDB(t)
	repo := NewUserRepository(db)

	displayName := "Delete Me"
	user := &models.User{
		SupabaseUserID:   uuid.New(),
		Email:            "delete@example.com",
		DisplayName:      &displayName,
		AvatarURL:        nil,
		Bio:              nil,
		Location:         nil,
		Website:          nil,
		GithubUsername:   nil,
		TwitterUsername:  nil,
		TotalXP:          0,
		PracticeScore:    0,
		GlobalRank:       nil,
		CurrentStreak:    0,
		LongestStreak:    0,
		LastActivityDate: nil,
	}

	err := repo.Create(user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Delete
	err = repo.Delete(user.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	fetched, err := repo.FindByID(user.ID)
	if err != nil {
		t.Fatalf("FindByID after delete should not error: %v", err)
	}
	if fetched != nil {
		t.Error("Expected user to be deleted")
	}
}