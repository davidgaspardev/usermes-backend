package valueobject

import (
	"regexp"
	"strings"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
)

// Email represents a validated email address
type Email struct {
	value string
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a new Email value object with validation
func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, errors.ErrInvalidEmail
	}

	if len(email) > 255 {
		return Email{}, errors.ErrEmailTooLong
	}

	if !emailRegex.MatchString(email) {
		return Email{}, errors.ErrInvalidEmailFormat
	}

	return Email{value: email}, nil
}

// Value returns the email string value
func (e Email) Value() string {
	return e.value
}

// String implements the Stringer interface
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
