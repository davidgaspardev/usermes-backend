package valueobject

import (
	"testing"
)

func TestNewPassword_Valid(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"valid password", "Password123"},
		{"minimum length", "Pass1234"},
		{"with special chars", "Pass@123!"},
		{"long password", "ThisIsAVeryLongPassword123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pwd, err := NewPassword(tt.password)
			if err != nil {
				t.Errorf("NewPassword(%s) returned error: %v", tt.password, err)
			}
			if pwd.Hash() == "" {
				t.Error("Password hash should not be empty")
			}
			if pwd.Hash() == tt.password {
				t.Error("Password should be hashed, not plain text")
			}
		})
	}
}

func TestNewPassword_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"empty password", ""},
		{"too short", "Pass1"},
		{"no number", "Password"},
		{"no letter", "12345678"},
		{"only spaces", "        "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.password)
			if err == nil {
				t.Errorf("NewPassword(%s) should return error but didn't", tt.password)
			}
		})
	}
}

func TestPassword_Compare(t *testing.T) {
	plainPassword := "Password123"
	pwd, err := NewPassword(plainPassword)
	if err != nil {
		t.Fatalf("Failed to create password: %v", err)
	}

	// Test correct password
	if !pwd.Compare(plainPassword) {
		t.Error("Compare should return true for correct password")
	}

	// Test wrong password
	if pwd.Compare("WrongPassword123") {
		t.Error("Compare should return false for wrong password")
	}

	// Test empty password
	if pwd.Compare("") {
		t.Error("Compare should return false for empty password")
	}
}

func TestNewPasswordFromHash(t *testing.T) {
	hash := "$2a$12$abcdefghijklmnopqrstuvwxyz123456789"
	pwd := NewPasswordFromHash(hash)

	if pwd.Hash() != hash {
		t.Errorf("Expected hash '%s', got '%s'", hash, pwd.Hash())
	}
}

func TestPassword_HashIsDifferentEachTime(t *testing.T) {
	plainPassword := "Password123"

	pwd1, _ := NewPassword(plainPassword)
	pwd2, _ := NewPassword(plainPassword)

	if pwd1.Hash() == pwd2.Hash() {
		t.Error("Same password should produce different hashes (due to salt)")
	}

	// But both should validate against the original password
	if !pwd1.Compare(plainPassword) || !pwd2.Compare(plainPassword) {
		t.Error("Both hashes should validate against original password")
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	pwd, err := GenerateRandomPassword(16)
	if err != nil {
		t.Errorf("GenerateRandomPassword returned error: %v", err)
	}

	if len(pwd) != 16 {
		t.Errorf("Expected password length 16, got %d", len(pwd))
	}

	// Generate another and ensure they're different
	pwd2, _ := GenerateRandomPassword(16)
	if pwd == pwd2 {
		t.Error("Random passwords should be different")
	}
}

func TestPassword_TooLong(t *testing.T) {
	// Password longer than 72 characters (bcrypt limit)
	longPassword := string(make([]byte, 80)) + "A1"
	_, err := NewPassword(longPassword)
	if err == nil {
		t.Error("Should fail for password longer than 72 characters")
	}
}

func TestPassword_Trimming(t *testing.T) {
	pwd, err := NewPassword("  Password123  ")
	if err != nil {
		t.Errorf("Should trim spaces: %v", err)
	}
	// Should compare without spaces
	if !pwd.Compare("Password123") {
		t.Error("Should compare trimmed password")
	}
}

func TestGenerateRandomPassword_MinLength(t *testing.T) {
	// Test with length below minimum
	pwd, err := GenerateRandomPassword(5)
	if err != nil {
		t.Errorf("GenerateRandomPassword returned error: %v", err)
	}
	// Should be adjusted to minimum length
	if len(pwd) < minPasswordLength {
		t.Errorf("Generated password should be at least %d characters", minPasswordLength)
	}
}

func TestGenerateRandomPassword_MaxLength(t *testing.T) {
	// Test with length above maximum
	pwd, err := GenerateRandomPassword(100)
	if err != nil {
		t.Errorf("GenerateRandomPassword returned error: %v", err)
	}
	// Should be adjusted to maximum length
	if len(pwd) > maxPasswordLength {
		t.Errorf("Generated password should be at most %d characters", maxPasswordLength)
	}
}

func TestPassword_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"with spaces at end", "Password123   ", false},
		{"exactly 8 chars", "Pass1234", false},
		{"exactly 72 chars", string(make([]byte, 64)) + "Pass1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
