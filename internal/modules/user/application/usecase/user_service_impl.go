package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepository output.UserRepository
	tokenGenerator output.TokenGenerator
	tokenDuration  time.Duration
}

// NewUserService creates a new instance of UserServiceImpl
func NewUserService(
	userRepository output.UserRepository,
	tokenGenerator output.TokenGenerator,
	tokenDuration time.Duration,
) input.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
		tokenGenerator: tokenGenerator,
		tokenDuration:  tokenDuration,
	}
}

// Register creates a new user account
func (s *UserServiceImpl) Register(ctx context.Context, emailStr, password, name string) (*entity.User, error) {
	// Validate name
	if err := validateName(name); err != nil {
		return nil, err
	}

	// Create email value object
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	// Check if email already exists
	exists, err := s.userRepository.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrEmailAlreadyExists
	}

	// Create password value object (this will hash the password)
	passwordVO, err := valueobject.NewPassword(password)
	if err != nil {
		return nil, err
	}

	// Create new user entity
	user := entity.NewUser(email, passwordVO, name)

	// Persist user
	if err := s.userRepository.Save(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates a user and returns a token
func (s *UserServiceImpl) Login(ctx context.Context, emailStr, password string) (string, *entity.User, error) {
	// Create email value object
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	// Find user by email
	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return "", nil, errors.ErrUserInactive
	}

	// Verify password
	if !user.VerifyPassword(password) {
		return "", nil, errors.ErrInvalidCredentials
	}

	// Record login
	user.RecordLogin()
	if err := s.userRepository.Update(ctx, user); err != nil {
		return "", nil, err
	}

	// Generate token
	token, err := s.tokenGenerator.GenerateToken(user.ID(), user.Email().Value(), s.tokenDuration)
	if err != nil {
		return "", nil, errors.ErrTokenGenerationFailed
	}

	return token, user, nil
}

// GetUserByID retrieves a user by their ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error) {
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// GetUserByEmail retrieves a user by their email
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, emailStr string) (*entity.User, error) {
	email, err := valueobject.NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepository.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// UpdateUser updates user information
func (s *UserServiceImpl) UpdateUser(ctx context.Context, userID uuid.UUID, name string) (*entity.User, error) {
	// Validate name
	if err := validateName(name); err != nil {
		return nil, err
	}

	// Find user
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}

	// Update name
	user.UpdateName(name)

	// Persist changes
	if err := s.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePassword changes a user's password
func (s *UserServiceImpl) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	// Find user
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Verify old password
	if !user.VerifyPassword(oldPassword) {
		return errors.ErrInvalidPassword
	}

	// Create new password value object
	newPasswordVO, err := valueobject.NewPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	user.UpdatePassword(newPasswordVO)

	// Persist changes
	if err := s.userRepository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// DeactivateUser deactivates a user account
func (s *UserServiceImpl) DeactivateUser(ctx context.Context, userID uuid.UUID) error {
	// Find user
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Deactivate user
	user.Deactivate()

	// Persist changes
	if err := s.userRepository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// ActivateUser activates a user account
func (s *UserServiceImpl) ActivateUser(ctx context.Context, userID uuid.UUID) error {
	// Find user
	user, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return errors.ErrUserNotFound
	}

	// Activate user
	user.Activate()

	// Persist changes
	if err := s.userRepository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

// validateName validates the user's name
func validateName(name string) error {
	if name == "" {
		return errors.ErrUserNameRequired
	}

	if len(name) < 2 {
		return errors.ErrUserNameTooShort
	}

	if len(name) > 100 {
		return errors.ErrUserNameTooLong
	}

	return nil
}
