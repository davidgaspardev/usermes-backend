package dto

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
)

func createTestUser(t *testing.T) *entity.User {
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(t, err)

	password, err := valueobject.NewPassword("Password123")
	require.NoError(t, err)

	user := entity.NewUser(email, password, "John Doe")

	return user
}

func TestToUserResponse(t *testing.T) {
	user := createTestUser(t)

	response := ToUserResponse(user)

	assert.Equal(t, user.ID().String(), response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "John Doe", response.Name)
	assert.True(t, response.IsActive)
	assert.Equal(t, user.CreatedAt(), response.CreatedAt)
	assert.Equal(t, user.UpdatedAt(), response.UpdatedAt)
	assert.Equal(t, user.LastLoginAt(), response.LastLoginAt)
}

func TestToUserResponse_WithLastLogin(t *testing.T) {
	user := createTestUser(t)

	// Record a login
	user.RecordLogin()

	response := ToUserResponse(user)

	assert.Equal(t, user.ID().String(), response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "John Doe", response.Name)
	assert.True(t, response.IsActive)
	assert.NotNil(t, response.LastLoginAt)
	assert.Equal(t, user.LastLoginAt(), response.LastLoginAt)
}

func TestToUserResponse_DeactivatedUser(t *testing.T) {
	user := createTestUser(t)
	user.Deactivate()

	response := ToUserResponse(user)

	assert.Equal(t, user.ID().String(), response.ID)
	assert.False(t, response.IsActive)
}

func TestToLoginResponse(t *testing.T) {
	user := createTestUser(t)
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token" // nolint:gosec

	response := ToLoginResponse(token, user)

	assert.Equal(t, token, response.Token)
	assert.Equal(t, user.ID().String(), response.User.ID)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "John Doe", response.User.Name)
	assert.True(t, response.User.IsActive)
}

func TestToLoginResponse_WithEmptyToken(t *testing.T) {
	user := createTestUser(t)
	token := ""

	response := ToLoginResponse(token, user)

	assert.Equal(t, "", response.Token)
	assert.Equal(t, user.ID().String(), response.User.ID)
}

func TestNewErrorResponse(t *testing.T) {
	t.Run("with both error and message", func(t *testing.T) {
		response := NewErrorResponse("validation_error", "Invalid email format")

		assert.Equal(t, "validation_error", response.Error)
		assert.Equal(t, "Invalid email format", response.Message)
	})

	t.Run("with empty message", func(t *testing.T) {
		response := NewErrorResponse("internal_error", "")

		assert.Equal(t, "internal_error", response.Error)
		assert.Equal(t, "", response.Message)
	})

	t.Run("with empty error", func(t *testing.T) {
		response := NewErrorResponse("", "Something went wrong")

		assert.Equal(t, "", response.Error)
		assert.Equal(t, "Something went wrong", response.Message)
	})
}

func TestNewSuccessResponse(t *testing.T) {
	t.Run("with data", func(t *testing.T) {
		data := map[string]string{"user_id": "123", "status": "active"}
		response := NewSuccessResponse("User created successfully", data)

		assert.Equal(t, "User created successfully", response.Message)
		assert.Equal(t, data, response.Data)
	})

	t.Run("with nil data", func(t *testing.T) {
		response := NewSuccessResponse("Operation completed", nil)

		assert.Equal(t, "Operation completed", response.Message)
		assert.Nil(t, response.Data)
	})

	t.Run("with empty message", func(t *testing.T) {
		response := NewSuccessResponse("", "some data")

		assert.Equal(t, "", response.Message)
		assert.Equal(t, "some data", response.Data)
	})
}

func TestRegisterRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "Password123",
			Name:     "John Doe",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty email", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "",
			Password: "Password123",
			Name:     "John Doe",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "email", validationErr.Field)
		assert.Equal(t, "email is required", validationErr.Message)
	})

	t.Run("empty password", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "",
			Name:     "John Doe",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "password", validationErr.Field)
		assert.Equal(t, "password is required", validationErr.Message)
	})

	t.Run("empty name", func(t *testing.T) {
		req := RegisterRequest{
			Email:    "test@example.com",
			Password: "Password123",
			Name:     "",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "name", validationErr.Field)
		assert.Equal(t, "name is required", validationErr.Message)
	})
}

func TestLoginRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "Password123",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty email", func(t *testing.T) {
		req := LoginRequest{
			Email:    "",
			Password: "Password123",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "email", validationErr.Field)
		assert.Equal(t, "email is required", validationErr.Message)
	})

	t.Run("empty password", func(t *testing.T) {
		req := LoginRequest{
			Email:    "test@example.com",
			Password: "",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "password", validationErr.Field)
		assert.Equal(t, "password is required", validationErr.Message)
	})
}

func TestUpdateUserRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := UpdateUserRequest{
			Name: "Updated Name",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty name", func(t *testing.T) {
		req := UpdateUserRequest{
			Name: "",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "name", validationErr.Field)
		assert.Equal(t, "name is required", validationErr.Message)
	})
}

func TestChangePasswordRequest_Validate(t *testing.T) {
	t.Run("valid request", func(t *testing.T) {
		req := ChangePasswordRequest{
			OldPassword: "OldPassword123",
			NewPassword: "NewPassword456",
		}

		err := req.Validate()
		assert.NoError(t, err)
	})

	t.Run("empty old password", func(t *testing.T) {
		req := ChangePasswordRequest{
			OldPassword: "",
			NewPassword: "NewPassword456",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "old_password", validationErr.Field)
		assert.Equal(t, "old password is required", validationErr.Message)
	})

	t.Run("empty new password", func(t *testing.T) {
		req := ChangePasswordRequest{
			OldPassword: "OldPassword123",
			NewPassword: "",
		}

		err := req.Validate()
		assert.Error(t, err)

		validationErr, ok := err.(*ValidationError)
		require.True(t, ok)
		assert.Equal(t, "new_password", validationErr.Field)
		assert.Equal(t, "new password is required", validationErr.Message)
	})
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}

	assert.Equal(t, "invalid email format", err.Error())
}

func TestRegisterRequest_JSONTags(t *testing.T) {
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "Password123",
		Name:     "John Doe",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Password123", req.Password)
	assert.Equal(t, "John Doe", req.Name)
}

func TestLoginRequest_JSONTags(t *testing.T) {
	req := LoginRequest{
		Email:    "test@example.com",
		Password: "Password123",
	}

	assert.Equal(t, "test@example.com", req.Email)
	assert.Equal(t, "Password123", req.Password)
}

func TestUpdateUserRequest_JSONTags(t *testing.T) {
	req := UpdateUserRequest{
		Name: "Updated Name",
	}

	assert.Equal(t, "Updated Name", req.Name)
}

func TestChangePasswordRequest_JSONTags(t *testing.T) {
	req := ChangePasswordRequest{
		OldPassword: "OldPassword123",
		NewPassword: "NewPassword456",
	}

	assert.Equal(t, "OldPassword123", req.OldPassword)
	assert.Equal(t, "NewPassword456", req.NewPassword)
}

func TestUserResponse_JSONTags(t *testing.T) {
	id := uuid.New()
	now := time.Now()
	lastLogin := time.Now().Add(-time.Hour)

	response := UserResponse{
		CreatedAt:   now,
		UpdatedAt:   now,
		LastLoginAt: &lastLogin,
		ID:          id.String(),
		Email:       "test@example.com",
		Name:        "John Doe",
		IsActive:    true,
	}

	assert.Equal(t, now, response.CreatedAt)
	assert.Equal(t, now, response.UpdatedAt)
	assert.Equal(t, &lastLogin, response.LastLoginAt)
	assert.Equal(t, id.String(), response.ID)
	assert.Equal(t, "test@example.com", response.Email)
	assert.Equal(t, "John Doe", response.Name)
	assert.True(t, response.IsActive)
}

func TestUserResponse_WithNilLastLoginAt(t *testing.T) {
	response := UserResponse{
		LastLoginAt: nil,
		IsActive:    false,
	}

	assert.Nil(t, response.LastLoginAt)
	assert.False(t, response.IsActive)
}

func TestLoginResponse_Structure(t *testing.T) {
	user := createTestUser(t)
	token := "test.jwt.token" // nolint:gosec

	response := LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}

	assert.Equal(t, token, response.Token)
	assert.Equal(t, user.ID().String(), response.User.ID)
	assert.Equal(t, "test@example.com", response.User.Email)
	assert.Equal(t, "John Doe", response.User.Name)
}

func TestUserIDParam_Structure(t *testing.T) {
	userID := uuid.New()
	param := UserIDParam{
		UserID: userID,
	}

	assert.Equal(t, userID, param.UserID)
}

// Benchmark tests for performance-critical conversion functions
func BenchmarkToUserResponse(b *testing.B) {
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(b, err)
	password, err := valueobject.NewPassword("Password123")
	require.NoError(b, err)
	user := entity.NewUser(email, password, "John Doe")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToUserResponse(user)
	}
}

func BenchmarkToLoginResponse(b *testing.B) {
	email, err := valueobject.NewEmail("test@example.com")
	require.NoError(b, err)
	password, err := valueobject.NewPassword("Password123")
	require.NoError(b, err)
	user := entity.NewUser(email, password, "John Doe")
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.test.token" // nolint:gosec

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToLoginResponse(token, user)
	}
}

func BenchmarkValidateRegisterRequest(b *testing.B) {
	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "Password123",
		Name:     "John Doe",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req.Validate() // nolint:errcheck
	}
}
