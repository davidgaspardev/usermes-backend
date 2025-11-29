package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
	"github.com/google/uuid"
)

// Mock Repository
type MockUserRepository struct {
	users        map[uuid.UUID]*entity.User
	usersByEmail map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users:        make(map[uuid.UUID]*entity.User),
		usersByEmail: make(map[string]*entity.User),
	}
}

func (m *MockUserRepository) Save(ctx context.Context, user *entity.User) error {
	if _, exists := m.usersByEmail[user.Email().Value()]; exists {
		return errors.ErrEmailAlreadyExists
	}
	m.users[user.ID()] = user
	m.usersByEmail[user.Email().Value()] = user
	return nil
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if _, exists := m.users[user.ID()]; !exists {
		return errors.ErrUserNotFound
	}
	m.users[user.ID()] = user
	return nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, exists := m.users[id]
	if !exists {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	user, exists := m.usersByEmail[email.Value()]
	if !exists {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (m *MockUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	_, exists := m.usersByEmail[email.Value()]
	return exists, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	user, exists := m.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}
	delete(m.users, id)
	delete(m.usersByEmail, user.Email().Value())
	return nil
}

func (m *MockUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	users := make([]*entity.User, 0)
	for _, user := range m.users {
		users = append(users, user)
	}
	return users, nil
}

// Mock Token Generator
type MockTokenGenerator struct{}

func (m *MockTokenGenerator) GenerateToken(userID uuid.UUID, email string, expiresIn time.Duration) (string, error) {
	return "mock-token-" + userID.String(), nil
}

func (m *MockTokenGenerator) ValidateToken(token string) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (m *MockTokenGenerator) RefreshToken(token string, expiresIn time.Duration) (string, error) {
	return "refreshed-" + token, nil
}

func setupTestService() *UserServiceImpl {
	return &UserServiceImpl{
		userRepository: NewMockUserRepository(),
		tokenGenerator: &MockTokenGenerator{},
		tokenDuration:  24 * time.Hour,
	}
}

func TestUserService_Register(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	user, err := service.Register(ctx, "test@example.com", "Password123", "Test User")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	if user == nil {
		t.Fatal("User should not be nil")
	}

	if user.Name() != "Test User" {
		t.Errorf("Expected name 'Test User', got '%s'", user.Name())
	}

	if user.Email().Value() != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email().Value())
	}
}

func TestUserService_Register_DuplicateEmail(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	_, err := service.Register(ctx, "test@example.com", "Password123", "User 1")
	if err != nil {
		t.Fatalf("First registration failed: %v", err)
	}

	_, err = service.Register(ctx, "test@example.com", "Password456", "User 2")
	if err != errors.ErrEmailAlreadyExists {
		t.Errorf("Expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestUserService_Register_InvalidData(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	tests := []struct {
		name     string
		email    string
		password string
		userName string
		wantErr  bool
	}{
		{"invalid email", "invalid-email", "Password123", "User", true},
		{"invalid password", "test@example.com", "short", "User", true},
		{"empty name", "test@example.com", "Password123", "", true},
		{"short name", "test@example.com", "Password123", "A", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Register(ctx, tt.email, tt.password, tt.userName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_Login(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register user first
	_, err := service.Register(ctx, "test@example.com", "Password123", "Test User")
	if err != nil {
		t.Fatalf("Register failed: %v", err)
	}

	// Login
	token, user, err := service.Login(ctx, "test@example.com", "Password123")
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	if token == "" {
		t.Error("Token should not be empty")
	}

	if user == nil {
		t.Fatal("User should not be nil")
	}

	if user.LastLoginAt() == nil {
		t.Error("LastLoginAt should be set after login")
	}
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register user
	service.Register(ctx, "test@example.com", "Password123", "Test User")

	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"wrong password", "test@example.com", "WrongPassword"},
		{"wrong email", "wrong@example.com", "Password123"},
		{"nonexistent user", "nonexistent@example.com", "Password123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := service.Login(ctx, tt.email, tt.password)
			if err != errors.ErrInvalidCredentials {
				t.Errorf("Expected ErrInvalidCredentials, got %v", err)
			}
		})
	}
}

func TestUserService_GetUserByID(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register user
	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	// Get user by ID
	user, err := service.GetUserByID(ctx, registered.ID())
	if err != nil {
		t.Fatalf("GetUserByID failed: %v", err)
	}

	if user.ID() != registered.ID() {
		t.Error("Retrieved user ID doesn't match")
	}
}

func TestUserService_GetUserByID_NotFound(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	_, err := service.GetUserByID(ctx, uuid.New())
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_GetUserByEmail(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	user, err := service.GetUserByEmail(ctx, "test@example.com")
	if err != nil {
		t.Fatalf("GetUserByEmail failed: %v", err)
	}

	if user.ID() != registered.ID() {
		t.Error("Retrieved user ID doesn't match")
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register user
	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Old Name")

	// Update user
	updated, err := service.UpdateUser(ctx, registered.ID(), "New Name")
	if err != nil {
		t.Fatalf("UpdateUser failed: %v", err)
	}

	if updated.Name() != "New Name" {
		t.Errorf("Expected name 'New Name', got '%s'", updated.Name())
	}
}

func TestUserService_UpdateUser_InvalidName(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	tests := []struct {
		name    string
		newName string
		wantErr bool
	}{
		{"empty name", "", true},
		{"too short", "A", true},
		{"too long", string(make([]byte, 101)), true},
		{"valid name", "Valid Name", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.UpdateUser(ctx, registered.ID(), tt.newName)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserService_ChangePassword(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register user
	registered, _ := service.Register(ctx, "test@example.com", "OldPassword123", "Test User")

	// Change password
	err := service.ChangePassword(ctx, registered.ID(), "OldPassword123", "NewPassword456")
	if err != nil {
		t.Fatalf("ChangePassword failed: %v", err)
	}

	// Try login with new password
	_, _, err = service.Login(ctx, "test@example.com", "NewPassword456")
	if err != nil {
		t.Error("Should be able to login with new password")
	}

	// Try login with old password (should fail)
	_, _, err = service.Login(ctx, "test@example.com", "OldPassword123")
	if err == nil {
		t.Error("Should not be able to login with old password")
	}
}

func TestUserService_ChangePassword_WrongOldPassword(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	err := service.ChangePassword(ctx, registered.ID(), "WrongPassword", "NewPassword456")
	if err != errors.ErrInvalidPassword {
		t.Errorf("Expected ErrInvalidPassword, got %v", err)
	}
}

func TestUserService_ChangePassword_InvalidNewPassword(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	err := service.ChangePassword(ctx, registered.ID(), "Password123", "short")
	if err == nil {
		t.Error("Should fail with invalid new password")
	}
}

func TestUserService_DeactivateUser(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")

	err := service.DeactivateUser(ctx, registered.ID())
	if err != nil {
		t.Fatalf("DeactivateUser failed: %v", err)
	}

	// Try to login (should fail)
	_, _, err = service.Login(ctx, "test@example.com", "Password123")
	if err != errors.ErrUserInactive {
		t.Errorf("Expected ErrUserInactive, got %v", err)
	}
}

func TestUserService_DeactivateUser_NotFound(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	err := service.DeactivateUser(ctx, uuid.New())
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestUserService_ActivateUser(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")
	service.DeactivateUser(ctx, registered.ID())

	err := service.ActivateUser(ctx, registered.ID())
	if err != nil {
		t.Fatalf("ActivateUser failed: %v", err)
	}

	// Should be able to login now
	_, _, err = service.Login(ctx, "test@example.com", "Password123")
	if err != nil {
		t.Error("Should be able to login after reactivation")
	}
}

func TestUserService_ActivateUser_NotFound(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	err := service.ActivateUser(ctx, uuid.New())
	if err != errors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got %v", err)
	}
}

func TestNewUserService(t *testing.T) {
	repo := NewMockUserRepository()
	tokenGen := &MockTokenGenerator{}
	duration := 24 * time.Hour

	service := NewUserService(repo, tokenGen, duration)

	if service == nil {
		t.Fatal("Service should not be nil")
	}
}

func TestUserService_Register_InvalidEmail(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	_, err := service.Register(ctx, "invalid-email", "Password123", "Test User")
	if err == nil {
		t.Error("Should fail with invalid email")
	}
}

func TestUserService_Login_InactiveUser(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	// Register and deactivate user
	registered, _ := service.Register(ctx, "test@example.com", "Password123", "Test User")
	service.DeactivateUser(ctx, registered.ID())

	// Try to login
	_, _, err := service.Login(ctx, "test@example.com", "Password123")
	if err != errors.ErrUserInactive {
		t.Errorf("Expected ErrUserInactive, got %v", err)
	}
}

func TestUserService_GetUserByEmail_InvalidEmail(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	_, err := service.GetUserByEmail(ctx, "invalid-email")
	if err == nil {
		t.Error("Should fail with invalid email")
	}
}

func TestUserService_UpdateUser_NotFound(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	_, err := service.UpdateUser(ctx, uuid.New(), "New Name")
	if err == nil {
		t.Error("Should fail when user not found")
	}
}

func TestUserService_ChangePassword_UserNotFound(t *testing.T) {
	service := setupTestService()
	ctx := context.Background()

	err := service.ChangePassword(ctx, uuid.New(), "OldPass", "NewPass")
	if err == nil {
		t.Error("Should fail when user not found")
	}
}
