package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/valueobject"
)

// UserRepository defines the output port for user persistence operations
type UserRepository interface {
	// Save persists a new user
	Save(ctx context.Context, user *entity.User) error

	// Update updates an existing user
	Update(ctx context.Context, user *entity.User) error

	// FindByID retrieves a user by their unique identifier
	FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)

	// FindByEmail retrieves a user by their email address
	FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error)

	// FindByUsername retrieves a user by their username
	FindByUsername(ctx context.Context, username valueobject.Username) (*entity.User, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error)

	// ExistsByUsername checks if a user with the given username exists
	ExistsByUsername(ctx context.Context, username valueobject.Username) (bool, error)

	// Delete removes a user from the repository
	Delete(ctx context.Context, id uuid.UUID) error

	// FindAll retrieves all users (with optional pagination)
	FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error)
}
