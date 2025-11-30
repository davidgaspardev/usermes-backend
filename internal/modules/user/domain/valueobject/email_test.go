package valueobject

import (
	"testing"
)

func TestNewEmail_Valid(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"valid email", "test@example.com"},
		{"valid email with subdomain", "user@mail.example.com"},
		{"valid email with plus", "user+tag@example.com"},
		{"valid email with dots", "first.last@example.com"},
		{"valid email with numbers", "user123@example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			email, err := NewEmail(tt.email)
			if err != nil {
				t.Errorf("NewEmail(%s) returned error: %v", tt.email, err)
			}
			if email.Value() == "" {
				t.Errorf("NewEmail(%s) returned empty value", tt.email)
			}
		})
	}
}

func TestNewEmail_Invalid(t *testing.T) {
	tests := []struct {
		name  string
		email string
	}{
		{"empty email", ""},
		{"missing @", "testexample.com"},
		{"missing domain", "test@"},
		{"missing local part", "@example.com"},
		{"invalid format", "test@@example.com"},
		{"spaces", "test @example.com"},
		{"no TLD", "test@example"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.email)
			if err == nil {
				t.Errorf("NewEmail(%s) should return error but didn't", tt.email)
			}
		})
	}
}

func TestEmail_Normalization(t *testing.T) {
	email1, err := NewEmail("Test@Example.COM")
	if err != nil {
		t.Fatalf("Failed to create email1: %v", err)
	}
	email2, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email2: %v", err)
	}

	if email1.Value() != email2.Value() {
		t.Errorf("Email normalization failed: %s != %s", email1.Value(), email2.Value())
	}

	if email1.Value() != "test@example.com" {
		t.Errorf("Expected normalized email 'test@example.com', got '%s'", email1.Value())
	}
}

func TestEmail_Equals(t *testing.T) {
	email1, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email1: %v", err)
	}
	email2, err := NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("Failed to create email2: %v", err)
	}
	email3, err := NewEmail("other@example.com")
	if err != nil {
		t.Fatalf("Failed to create email3: %v", err)
	}

	if !email1.Equals(email2) {
		t.Error("Same emails should be equal")
	}

	if email1.Equals(email3) {
		t.Error("Different emails should not be equal")
	}
}

func TestEmail_String(t *testing.T) {
	emailStr := "test@example.com"
	email, err := NewEmail(emailStr)
	if err != nil {
		t.Fatalf("Failed to create email: %v", err)
	}

	if email.String() != emailStr {
		t.Errorf("Expected String() to return '%s', got '%s'", emailStr, email.String())
	}
}

func TestEmail_TooLong(t *testing.T) {
	// Email longer than 255 characters
	longEmail := string(make([]byte, 250)) + "@test.com"
	_, err := NewEmail(longEmail)
	if err == nil {
		t.Error("Should fail for email longer than 255 characters")
	}
}

func TestEmail_Trimming(t *testing.T) {
	email, err := NewEmail("  test@example.com  ")
	if err != nil {
		t.Errorf("Should trim spaces: %v", err)
	}
	if email.Value() != "test@example.com" {
		t.Errorf("Expected 'test@example.com', got '%s'", email.Value())
	}
}
