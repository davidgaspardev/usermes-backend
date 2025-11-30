package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
)

func createTestUser(t *testing.T) *User {
	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	password, err := valueobject.NewPassword("Password123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	return NewUser(email, password, "Test User")
}

func TestNewUser(t *testing.T) {
	user := createTestUser(t)

	if user.ID() == uuid.Nil {
		t.Error("User ID should not be nil")
	}

	if user.Name() != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name())
	}

	if !user.IsActive() {
		t.Error("New user should be active")
	}

	if user.CreatedAt().IsZero() {
		t.Error("CreatedAt should be set")
	}

	if user.UpdatedAt().IsZero() {
		t.Error("UpdatedAt should be set")
	}

	if user.LastLoginAt() != nil {
		t.Error("LastLoginAt should be nil for new user")
	}
}

func TestUser_UpdateName(t *testing.T) {
	user := createTestUser(t)
	originalUpdatedAt := user.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	user.UpdateName("New Name")

	if user.Name() != "New Name" {
		t.Errorf("Expected name 'New Name', got '%s'", user.Name())
	}

	if !user.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestUser_UpdateEmail(t *testing.T) {
	user := createTestUser(t)
	newEmail, err := valueobject.NewEmail("newemail@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	user.UpdateEmail(newEmail)

	if user.Email().Value() != "newemail@example.com" {
		t.Errorf("Expected email 'newemail@example.com', got '%s'", user.Email().Value())
	}
}

func TestUser_UpdatePassword(t *testing.T) {
	user := createTestUser(t)
	newPassword, err := valueobject.NewPassword("NewPassword456")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	user.UpdatePassword(newPassword)

	if !user.VerifyPassword("NewPassword456") {
		t.Error("New password should be verified")
	}

	if user.VerifyPassword("Password123") {
		t.Error("Old password should not work")
	}
}

func TestUser_VerifyPassword(t *testing.T) {
	user := createTestUser(t)

	if !user.VerifyPassword("Password123") {
		t.Error("Correct password should be verified")
	}

	if user.VerifyPassword("WrongPassword") {
		t.Error("Wrong password should not be verified")
	}
}

func TestUser_Deactivate(t *testing.T) {
	user := createTestUser(t)

	user.Deactivate()

	if user.IsActive() {
		t.Error("User should be inactive after deactivation")
	}
}

func TestUser_Activate(t *testing.T) {
	user := createTestUser(t)
	user.Deactivate()

	user.Activate()

	if !user.IsActive() {
		t.Error("User should be active after activation")
	}
}

func TestUser_RecordLogin(t *testing.T) {
	user := createTestUser(t)

	if user.LastLoginAt() != nil {
		t.Error("New user should not have LastLoginAt")
	}

	user.RecordLogin()

	if user.LastLoginAt() == nil {
		t.Error("LastLoginAt should be set after login")
	}

	firstLogin := *user.LastLoginAt()
	time.Sleep(10 * time.Millisecond)

	user.RecordLogin()

	if !user.LastLoginAt().After(firstLogin) {
		t.Error("LastLoginAt should be updated on subsequent logins")
	}
}

func TestReconstructUser(t *testing.T) {
	id := uuid.New()
	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	password, err := valueobject.NewPassword("Password123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}
	now := time.Now()
	lastLogin := now.Add(-1 * time.Hour)

	user := ReconstructUser(
		id,
		email,
		password,
		"Test User",
		true,
		now,
		now,
		&lastLogin,
	)

	if user.ID() != id {
		t.Error("Reconstructed user should have correct ID")
	}

	if user.Name() != "Test User" {
		t.Error("Reconstructed user should have correct name")
	}

	if !user.IsActive() {
		t.Error("Reconstructed user should be active")
	}

	if user.LastLoginAt() == nil {
		t.Error("Reconstructed user should have LastLoginAt")
	}
}

func TestUser_AllGetters(t *testing.T) {
	user := createTestUser(t)

	// Test all getter methods
	if user.ID() == uuid.Nil {
		t.Error("ID should not be nil")
	}
	if user.Email().Value() != "test@example.com" {
		t.Error("Email getter failed")
	}
	if user.Password().Hash() == "" {
		t.Error("Password getter failed")
	}
	if user.Name() != "Test User" {
		t.Error("Name getter failed")
	}
	if !user.IsActive() {
		t.Error("IsActive getter failed")
	}
	if user.CreatedAt().IsZero() {
		t.Error("CreatedAt getter failed")
	}
	if user.UpdatedAt().IsZero() {
		t.Error("UpdatedAt getter failed")
	}
	if user.LastLoginAt() != nil {
		t.Error("LastLoginAt should be nil initially")
	}
}
