package errors

import "errors"

// Domain errors for User module
var (
	// Email errors
	ErrInvalidEmail       = errors.New("invalid email address")
	ErrEmailTooLong       = errors.New("email address is too long")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrEmailAlreadyExists = errors.New("email already exists")

	// Password errors
	ErrPasswordRequired   = errors.New("password is required")
	ErrPasswordTooShort   = errors.New("password must be at least 8 characters long")
	ErrPasswordTooLong    = errors.New("password is too long")
	ErrPasswordWeak       = errors.New("password must contain at least one letter and one number")
	ErrPasswordHashFailed = errors.New("failed to hash password")
	ErrInvalidPassword    = errors.New("invalid password")

	// Username errors
	ErrUsernameRequired      = errors.New("username is required")
	ErrUsernameTooShort      = errors.New("username must be at least 3 characters long")
	ErrUsernameTooLong       = errors.New("username must be at most 30 characters long")
	ErrInvalidUsernameFormat = errors.New("username can only contain letters, numbers, underscores, and hyphens")
	ErrUsernameAlreadyExists = errors.New("username already exists")

	// User errors
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserInactive       = errors.New("user account is inactive")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNameRequired   = errors.New("user name is required")
	ErrUserNameTooShort   = errors.New("user name is too short")
	ErrUserNameTooLong    = errors.New("user name is too long")

	// Authentication errors
	ErrUnauthorized          = errors.New("unauthorized")
	ErrInvalidToken          = errors.New("invalid token")
	ErrTokenExpired          = errors.New("token expired")
	ErrTokenGenerationFailed = errors.New("failed to generate token")
)
