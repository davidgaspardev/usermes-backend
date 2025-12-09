package valueobject

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/errors"
)

// Password represents a hashed password
type Password struct {
	hash string
}

const (
	minPasswordLength = 8
	maxPasswordLength = 72 // bcrypt limitation
	bcryptCost        = 12
)

// NewPassword creates a new Password value object from a plain text password
func NewPassword(plainPassword string) (Password, error) {
	plainPassword = strings.TrimSpace(plainPassword)

	if err := validatePassword(plainPassword); err != nil {
		return Password{}, err
	}

	hash, err := hashPassword(plainPassword)
	if err != nil {
		return Password{}, errors.ErrPasswordHashFailed
	}

	return Password{hash: hash}, nil
}

// NewPasswordFromHash creates a Password value object from an existing hash
// This is used when reconstructing a user from the database
func NewPasswordFromHash(hash string) Password {
	return Password{hash: hash}
}

// Hash returns the hashed password
func (p Password) Hash() string {
	return p.hash
}

// Compare checks if the provided plain password matches the hashed password
func (p Password) Compare(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(plainPassword))
	return err == nil
}

// validatePassword validates the plain password against business rules
func validatePassword(plainPassword string) error {
	if plainPassword == "" {
		return errors.ErrPasswordRequired
	}

	if len(plainPassword) < minPasswordLength {
		return errors.ErrPasswordTooShort
	}

	if len(plainPassword) > maxPasswordLength {
		return errors.ErrPasswordTooLong
	}

	// Check for at least one letter and one number
	hasLetter := false
	hasNumber := false

	for _, char := range plainPassword {
		switch {
		case char >= 'a' && char <= 'z', char >= 'A' && char <= 'Z':
			hasLetter = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}

		if hasLetter && hasNumber {
			break
		}
	}

	if !hasLetter || !hasNumber {
		return errors.ErrPasswordWeak
	}

	return nil
}

// hashPassword hashes a plain text password using bcrypt
func hashPassword(plainPassword string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// GenerateRandomPassword generates a secure random password
func GenerateRandomPassword(length int) (string, error) {
	if length < minPasswordLength {
		length = minPasswordLength
	}
	if length > maxPasswordLength {
		length = maxPasswordLength
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes)[:length], nil
}
