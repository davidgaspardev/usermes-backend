package persistence

import (
	"context"
	"testing"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
	"github.com/google/uuid"
)

func createTestUser(t *testing.T) *entity.User {
	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	password, err := valueobject.NewPassword("Password123")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	return entity.NewUser(email, password, "Test User")
}

func TestNewMemoryUserRepository(t *testing.T) {
	repo := NewMemoryUserRepository()

	if repo == nil {
		t.Fatal("Repository should not be nil")
	}
}

func TestMemoryUserRepository_Save(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()
	user := createTestUser(t)

	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify user was saved
	found, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.ID() != user.ID() {
		t.Errorf("Expected user ID %s, got %s", user.ID(), found.ID())
	}
}

func TestMemoryUserRepository_Save_DuplicateEmail(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user1 := createTestUser(t)
	err := repo.Save(ctx, user1)
	if err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	// Try to save another user with same email
	email, err := valueobject.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	password, err := valueobject.NewPassword("Password456")
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}
	user2 := entity.NewUser(email, password, "Another User")

	err = repo.Save(ctx, user2)
	if err != errors.ErrEmailAlreadyExists {
		t.Errorf("Expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestMemoryUserRepository_Save_DuplicateID(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Try to save same user again
	err = repo.Save(ctx, user)
	if err != errors.ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got %v", err)
	}
}

func TestMemoryUserRepository_Update(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	// Update user name
	user.UpdateName("Updated Name")

	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Verify update
	found, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("Failed to find user: %v", err)
	}
	if found.Name() != "Updated Name" {
		t.Errorf("Expected name 'Updated Name', got '%s'", found.Name())
	}
}

func TestMemoryUserRepository_Update_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)

	err := repo.Update(ctx, user)
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_Update_EmailChange(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	oldEmail := user.Email()

	// Update email
	newEmail, err := valueobject.NewEmail("newemail@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	user.UpdateEmail(newEmail)

	err = repo.Update(ctx, user)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Old email should not find the user
	_, err = repo.FindByEmail(ctx, oldEmail)
	if err != errors.ErrUserNotFound {
		t.Error("Old email should not find user after update")
	}

	// New email should find the user
	found, err := repo.FindByEmail(ctx, newEmail)
	if err != nil {
		t.Fatalf("FindByEmail with new email failed: %v", err)
	}
	if found.ID() != user.ID() {
		t.Error("Found user ID doesn't match")
	}
}

func TestMemoryUserRepository_FindByID(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	found, err := repo.FindByID(ctx, user.ID())
	if err != nil {
		t.Fatalf("FindByID failed: %v", err)
	}

	if found.ID() != user.ID() {
		t.Error("Found user ID doesn't match")
	}

	if found.Email().Value() != user.Email().Value() {
		t.Error("Found user email doesn't match")
	}
}

func TestMemoryUserRepository_FindByID_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	_, err := repo.FindByID(ctx, uuid.New())
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_FindByEmail(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	found, err := repo.FindByEmail(ctx, user.Email())
	if err != nil {
		t.Fatalf("FindByEmail failed: %v", err)
	}

	if found.ID() != user.ID() {
		t.Error("Found user ID doesn't match")
	}
}

func TestMemoryUserRepository_FindByEmail_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	email, err := valueobject.NewEmail("nonexistent@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	_, err = repo.FindByEmail(ctx, email)
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_ExistsByEmail(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	exists, err := repo.ExistsByEmail(ctx, user.Email())
	if err != nil {
		t.Fatalf("ExistsByEmail failed: %v", err)
	}

	if !exists {
		t.Error("User should exist")
	}
}

func TestMemoryUserRepository_ExistsByEmail_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	email, err := valueobject.NewEmail("nonexistent@example.com")
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}
	exists, err := repo.ExistsByEmail(ctx, email)
	if err != nil {
		t.Fatalf("ExistsByEmail failed: %v", err)
	}

	if exists {
		t.Error("User should not exist")
	}
}

func TestMemoryUserRepository_Delete(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	user := createTestUser(t)
	err := repo.Save(ctx, user)
	if err != nil {
		t.Fatalf("Failed to save user: %v", err)
	}

	err = repo.Delete(ctx, user.ID())
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify deletion
	_, err = repo.FindByID(ctx, user.ID())
	if err != errors.ErrUserNotFound {
		t.Error("User should not be found after deletion")
	}

	// Email should also not find user
	_, err = repo.FindByEmail(ctx, user.Email())
	if err != errors.ErrUserNotFound {
		t.Error("User should not be found by email after deletion")
	}
}

func TestMemoryUserRepository_Delete_NotFound(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	err := repo.Delete(ctx, uuid.New())
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestMemoryUserRepository_FindAll(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Add multiple users
	for i := 0; i < 5; i++ {
		email, err := valueobject.NewEmail("user" + string(rune('0'+i)) + "@example.com")
		if err != nil {
			t.Fatalf("Failed to create email: %v", err)
		}
		password, err := valueobject.NewPassword("Password123")
		if err != nil {
			t.Fatalf("Failed to create password: %v", err)
		}
		user := entity.NewUser(email, password, "User "+string(rune('0'+i)))
		err = repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Failed to save user: %v", err)
		}
	}

	users, err := repo.FindAll(ctx, 10, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}
}

func TestMemoryUserRepository_FindAll_Pagination(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Add multiple users
	for i := 0; i < 10; i++ {
		email, err := valueobject.NewEmail("user" + string(rune('0'+i)) + "@example.com")
		if err != nil {
			t.Fatalf("Failed to create email: %v", err)
		}
		password, err := valueobject.NewPassword("Password123")
		if err != nil {
			t.Fatalf("Failed to create password: %v", err)
		}
		user := entity.NewUser(email, password, "User "+string(rune('0'+i)))
		err = repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Failed to save user: %v", err)
		}
	}

	// Get first page
	users, err := repo.FindAll(ctx, 5, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users in first page, got %d", len(users))
	}

	// Get second page
	users, err = repo.FindAll(ctx, 5, 5)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users in second page, got %d", len(users))
	}

	// Get beyond available users
	users, err = repo.FindAll(ctx, 5, 15)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users beyond available, got %d", len(users))
	}
}

func TestMemoryUserRepository_FindAll_NoLimit(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Add multiple users
	for i := 0; i < 5; i++ {
		email, err := valueobject.NewEmail("user" + string(rune('0'+i)) + "@example.com")
		if err != nil {
			t.Fatalf("Failed to create email: %v", err)
		}
		password, err := valueobject.NewPassword("Password123")
		if err != nil {
			t.Fatalf("Failed to create password: %v", err)
		}
		user := entity.NewUser(email, password, "User "+string(rune('0'+i)))
		err = repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Failed to save user: %v", err)
		}
	}

	// Get all users (limit = 0)
	users, err := repo.FindAll(ctx, 0, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 5 {
		t.Errorf("Expected 5 users, got %d", len(users))
	}
}

func TestMemoryUserRepository_FindAll_Empty(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	users, err := repo.FindAll(ctx, 10, 0)
	if err != nil {
		t.Fatalf("FindAll failed: %v", err)
	}

	if len(users) != 0 {
		t.Errorf("Expected 0 users, got %d", len(users))
	}
}

func TestMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewMemoryUserRepository()
	ctx := context.Background()

	// Test concurrent saves
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(index int) {
			email, err := valueobject.NewEmail("user" + string(rune('0'+index)) + "@example.com")
			if err != nil {
				t.Errorf("Failed to create email: %v", err)
				done <- true
				return
			}
			password, err := valueobject.NewPassword("Password123")
			if err != nil {
				t.Errorf("Failed to create password: %v", err)
				done <- true
				return
			}
			user := entity.NewUser(email, password, "User "+string(rune('0'+index)))
			err = repo.Save(ctx, user)
			if err != nil {
				t.Errorf("Failed to save user: %v", err)
				done <- true
				return
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	users, err := repo.FindAll(ctx, 100, 0)
	if err != nil {
		t.Fatalf("Failed to find all users: %v", err)
	}
	if len(users) != 10 {
		t.Errorf("Expected 10 users after concurrent saves, got %d", len(users))
	}
}
