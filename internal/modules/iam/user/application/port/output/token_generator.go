package output

import (
	"time"

	"github.com/google/uuid"
)

// TokenGenerator defines the output port for JWT token generation and validation
type TokenGenerator interface {
	// GenerateToken creates a new JWT token for the given user ID
	GenerateToken(userID uuid.UUID, email string, expiresIn time.Duration) (string, error)

	// ValidateToken validates a JWT token and returns the user ID
	ValidateToken(token string) (uuid.UUID, error)

	// RefreshToken generates a new token from an existing valid token
	RefreshToken(token string, expiresIn time.Duration) (string, error)
}
