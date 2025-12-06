package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/dto"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService input.UserService
}

// NewUserHandler creates a new instance of UserHandler
func NewUserHandler(userService input.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register handles user registration
// POST /api/users/register
func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	// Call use case
	user, err := h.userService.Register(c.Context(), req.Email, req.Password, req.Username, req.Name)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusCreated).JSON(dto.NewSuccessResponse(
		"User registered successfully",
		dto.ToUserResponse(user),
	))
}

// Login handles user login
// POST /api/users/login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	// Call use case
	token, user, err := h.userService.Login(c.Context(), req.Username, req.Password)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.ToLoginResponse(token, user))
}

// GetUserByID handles getting a user by ID
// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	// Parse user ID from path parameter
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid user ID format",
		))
	}

	// Call use case
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.ToUserResponse(user))
}

// GetMe handles getting the current authenticated user
// GET /api/users/me
func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	// Get user ID from context (set by authentication middleware)
	userID, ok := c.Locals("userID").(uuid.UUID)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			"User not authenticated",
		))
	}

	// Call use case
	user, err := h.userService.GetUserByID(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.ToUserResponse(user))
}

// UpdateUser handles updating user information
// PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	// Parse user ID from path parameter
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid user ID format",
		))
	}

	var req dto.UpdateUserRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	// Call use case
	user, err := h.userService.UpdateUser(c.Context(), userID, req.Name)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse(
		"User updated successfully",
		dto.ToUserResponse(user),
	))
}

// ChangePassword handles changing user password
// POST /api/users/:id/change-password
func (h *UserHandler) ChangePassword(c *fiber.Ctx) error {
	// Parse user ID from path parameter
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid user ID format",
		))
	}

	var req dto.ChangePasswordRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	// Call use case
	err = h.userService.ChangePassword(c.Context(), userID, req.OldPassword, req.NewPassword)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse(
		"Password changed successfully",
		nil,
	))
}

// DeactivateUser handles deactivating a user account
// POST /api/users/:id/deactivate
func (h *UserHandler) DeactivateUser(c *fiber.Ctx) error {
	// Parse user ID from path parameter
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid user ID format",
		))
	}

	// Call use case
	err = h.userService.DeactivateUser(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse(
		"User deactivated successfully",
		nil,
	))
}

// ActivateUser handles activating a user account
// POST /api/users/:id/activate
func (h *UserHandler) ActivateUser(c *fiber.Ctx) error {
	// Parse user ID from path parameter
	userIDStr := c.Params("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"bad_request",
			"Invalid user ID format",
		))
	}

	// Call use case
	err = h.userService.ActivateUser(c.Context(), userID)
	if err != nil {
		return h.handleError(c, err)
	}

	// Return response
	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse(
		"User activated successfully",
		nil,
	))
}

// handleError maps domain errors to HTTP responses
func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	switch err {
	case errors.ErrInvalidEmail, errors.ErrEmailTooLong, errors.ErrInvalidEmailFormat:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_email",
			err.Error(),
		))
	case errors.ErrPasswordRequired, errors.ErrPasswordTooShort, errors.ErrPasswordTooLong, errors.ErrPasswordWeak:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_password",
			err.Error(),
		))
	case errors.ErrUsernameRequired, errors.ErrUsernameTooShort, errors.ErrUsernameTooLong, errors.ErrInvalidUsernameFormat:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_username",
			err.Error(),
		))
	case errors.ErrEmailAlreadyExists, errors.ErrUserAlreadyExists, errors.ErrUsernameAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(dto.NewErrorResponse(
			"conflict",
			err.Error(),
		))
	case errors.ErrUserNotFound:
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse(
			"not_found",
			err.Error(),
		))
	case errors.ErrInvalidCredentials:
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			"Invalid username or password",
		))
	case errors.ErrUserInactive:
		return c.Status(fiber.StatusForbidden).JSON(dto.NewErrorResponse(
			"forbidden",
			err.Error(),
		))
	case errors.ErrUnauthorized, errors.ErrInvalidToken, errors.ErrTokenExpired:
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			err.Error(),
		))
	case errors.ErrUserNameRequired, errors.ErrUserNameTooShort, errors.ErrUserNameTooLong:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_name",
			err.Error(),
		))
	case errors.ErrInvalidPassword:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_password",
			err.Error(),
		))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse(
			"internal_error",
			"An internal error occurred",
		))
	}
}
