package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/domain/entity"
)

// UserService defines the input port for user operations (use cases)
type UserService interface {
	// Register creates a new user account
	Register(ctx context.Context, email, password, username, name string) (*entity.User, error)

	// Login authenticates a user and returns a token
	Login(ctx context.Context, username, password string) (token string, user *entity.User, err error)

	// GetUserByID retrieves a user by their ID
	GetUserByID(ctx context.Context, userID uuid.UUID) (*entity.User, error)

	// GetUserByEmail retrieves a user by their email
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)

	// GetUserByUsername retrieves a user by their username
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)

	// UpdateUser updates user information
	UpdateUser(ctx context.Context, userID uuid.UUID, name string) (*entity.User, error)

	// ChangePassword changes a user's password
	ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error

	// DeactivateUser deactivates a user account
	DeactivateUser(ctx context.Context, userID uuid.UUID) error

	// ActivateUser activates a user account
	ActivateUser(ctx context.Context, userID uuid.UUID) error
}
