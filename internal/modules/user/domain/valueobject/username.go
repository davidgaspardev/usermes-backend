package valueobject

import (
	"regexp"
	"strings"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
)

// Username represents a validated username value object
type Username struct {
	value string
}

var (
	// usernameRegex validates username format (alphanumeric, underscore, hyphen)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
)

const (
	minUsernameLength = 3
	maxUsernameLength = 30
)

// NewUsername creates a new Username value object with validation
func NewUsername(username string) (Username, error) {
	// Trim whitespace
	username = strings.TrimSpace(username)

	// Check if username is empty
	if username == "" {
		return Username{}, errors.ErrUsernameRequired
	}

	// Check minimum length
	if len(username) < minUsernameLength {
		return Username{}, errors.ErrUsernameTooShort
	}

	// Check maximum length
	if len(username) > maxUsernameLength {
		return Username{}, errors.ErrUsernameTooLong
	}

	// Check format
	if !usernameRegex.MatchString(username) {
		return Username{}, errors.ErrInvalidUsernameFormat
	}

	return Username{value: username}, nil
}

// Value returns the username string value
func (u Username) Value() string {
	return u.value
}

// String implements the Stringer interface
func (u Username) String() string {
	return u.value
}

// Equals checks if two usernames are equal
func (u Username) Equals(other Username) bool {
	return u.value == other.value
}
