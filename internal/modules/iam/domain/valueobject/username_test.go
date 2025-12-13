package valueobject

import (
	"testing"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/domain/errors"
)

func TestNewUsername_Valid(t *testing.T) {
	tests := []struct {
		name     string
		username string
		expected string
	}{
		{
			name:     "valid username",
			username: "john_doe",
			expected: "john_doe",
		},
		{
			name:     "minimum length",
			username: "abc",
			expected: "abc",
		},
		{
			name:     "with numbers",
			username: "user123",
			expected: "user123",
		},
		{
			name:     "with hyphen",
			username: "john-doe",
			expected: "john-doe",
		},
		{
			name:     "with underscore",
			username: "john_doe_123",
			expected: "john_doe_123",
		},
		{
			name:     "maximum length",
			username: "abcdefghijklmnopqrstuvwxyz1234",
			expected: "abcdefghijklmnopqrstuvwxyz1234",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username, err := NewUsername(tt.username)
			if err != nil {
				t.Fatalf("Expected no error, got: %v", err)
			}

			if username.Value() != tt.expected {
				t.Errorf("Expected username '%s', got '%s'", tt.expected, username.Value())
			}
		})
	}
}

func TestNewUsername_Invalid(t *testing.T) {
	tests := []struct {
		expectedErr error
		name        string
		username    string
	}{
		{
			name:        "empty username",
			username:    "",
			expectedErr: errors.ErrUsernameRequired,
		},
		{
			name:        "too short",
			username:    "ab",
			expectedErr: errors.ErrUsernameTooShort,
		},
		{
			name:        "too long",
			username:    "abcdefghijklmnopqrstuvwxyz12345",
			expectedErr: errors.ErrUsernameTooLong,
		},
		{
			name:        "with spaces",
			username:    "john doe",
			expectedErr: errors.ErrInvalidUsernameFormat,
		},
		{
			name:        "with special characters",
			username:    "john@doe",
			expectedErr: errors.ErrInvalidUsernameFormat,
		},
		{
			name:        "with dots",
			username:    "john.doe",
			expectedErr: errors.ErrInvalidUsernameFormat,
		},
		{
			name:        "only spaces",
			username:    "   ",
			expectedErr: errors.ErrUsernameRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUsername(tt.username)
			if err == nil {
				t.Fatal("Expected error, got nil")
			}

			if err != tt.expectedErr {
				t.Errorf("Expected error '%v', got '%v'", tt.expectedErr, err)
			}
		})
	}
}

func TestUsername_Trimming(t *testing.T) {
	username, err := NewUsername("  john_doe  ")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if username.Value() != "john_doe" {
		t.Errorf("Expected username 'john_doe', got '%s'", username.Value())
	}
}

func TestUsername_String(t *testing.T) {
	username, err := NewUsername("john_doe")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if username.String() != "john_doe" {
		t.Errorf("Expected string 'john_doe', got '%s'", username.String())
	}
}

func TestUsername_Equals(t *testing.T) {
	username1, err := NewUsername("john_doe")
	if err != nil {
		t.Fatalf("Failed to create username: %v", err)
	}
	username2, err := NewUsername("john_doe")
	if err != nil {
		t.Fatalf("Failed to create username: %v", err)
	}
	username3, err := NewUsername("jane_doe")
	if err != nil {
		t.Fatalf("Failed to create username: %v", err)
	}

	if !username1.Equals(username2) {
		t.Error("Expected usernames to be equal")
	}

	if username1.Equals(username3) {
		t.Error("Expected usernames to be different")
	}
}
