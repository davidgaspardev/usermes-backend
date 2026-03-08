package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/domain/entity"
)

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"` //nolint:gosec // intentional: request DTO receives user-provided password
	Username string `json:"username"`
	Name     string `json:"name"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` //nolint:gosec // intentional: request DTO receives user-provided password
}

// UpdateUserRequest represents the request body for updating user information
type UpdateUserRequest struct {
	Name string `json:"name"`
}

// ChangePasswordRequest represents the request body for changing password
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// UserResponse represents the response body for user data
type UserResponse struct {
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	ID          string     `json:"id"`
	Email       string     `json:"email"`
	Username    string     `json:"username"`
	Name        string     `json:"name"`
	IsActive    bool       `json:"is_active"`
}

// LoginResponse represents the response body for successful login
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// ErrorResponse represents the response body for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ToUserResponse converts a User entity to UserResponse DTO
func ToUserResponse(user *entity.User) UserResponse {
	return UserResponse{
		ID:          user.ID().String(),
		Email:       user.Email().Value(),
		Username:    user.Username().Value(),
		Name:        user.Name(),
		IsActive:    user.IsActive(),
		CreatedAt:   user.CreatedAt(),
		UpdatedAt:   user.UpdatedAt(),
		LastLoginAt: user.LastLoginAt(),
	}
}

// ToLoginResponse converts a token and user entity to LoginResponse DTO
func ToLoginResponse(token string, user *entity.User) LoginResponse {
	return LoginResponse{
		Token: token,
		User:  ToUserResponse(user),
	}
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(errMsg, message string) ErrorResponse {
	return ErrorResponse{
		Error:   errMsg,
		Message: message,
	}
}

// NewSuccessResponse creates a new SuccessResponse
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}

// ValidateRegisterRequest validates the register request
func (r *RegisterRequest) Validate() error {
	if r.Email == "" {
		return &ValidationError{Field: "email", Message: "email is required"}
	}
	if r.Password == "" {
		return &ValidationError{Field: "password", Message: "password is required"}
	}
	if r.Username == "" {
		return &ValidationError{Field: "username", Message: "username is required"}
	}
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	return nil
}

// ValidateLoginRequest validates the login request
func (r *LoginRequest) Validate() error {
	if r.Username == "" {
		return &ValidationError{Field: "username", Message: "username is required"}
	}
	if r.Password == "" {
		return &ValidationError{Field: "password", Message: "password is required"}
	}
	return nil
}

// ValidateUpdateUserRequest validates the update user request
func (r *UpdateUserRequest) Validate() error {
	if r.Name == "" {
		return &ValidationError{Field: "name", Message: "name is required"}
	}
	return nil
}

// ValidateChangePasswordRequest validates the change password request
func (r *ChangePasswordRequest) Validate() error {
	if r.OldPassword == "" {
		return &ValidationError{Field: "old_password", Message: "old password is required"}
	}
	if r.NewPassword == "" {
		return &ValidationError{Field: "new_password", Message: "new password is required"}
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	return e.Message
}

// UserIDParam represents the user ID path parameter
type UserIDParam struct {
	UserID uuid.UUID
}
