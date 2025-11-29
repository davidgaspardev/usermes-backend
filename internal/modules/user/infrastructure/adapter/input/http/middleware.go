package http

import (
	"strings"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/infrastructure/dto"
	"github.com/gofiber/fiber/v2"
)

// AuthMiddleware provides authentication middleware for protected routes
type AuthMiddleware struct {
	tokenGenerator output.TokenGenerator
}

// NewAuthMiddleware creates a new instance of AuthMiddleware
func NewAuthMiddleware(tokenGenerator output.TokenGenerator) *AuthMiddleware {
	return &AuthMiddleware{
		tokenGenerator: tokenGenerator,
	}
}

// Authenticate is a middleware that validates JWT tokens
func (m *AuthMiddleware) Authenticate(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			"Authorization header is required",
		))
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			"Invalid authorization header format. Use: Bearer <token>",
		))
	}

	token := parts[1]

	// Validate token
	userID, err := m.tokenGenerator.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(dto.NewErrorResponse(
			"unauthorized",
			"Invalid or expired token",
		))
	}

	// Store user ID in context for use in handlers
	c.Locals("userID", userID)

	// Continue to next handler
	return c.Next()
}

// OptionalAuthenticate is a middleware that validates JWT tokens but doesn't require them
func (m *AuthMiddleware) OptionalAuthenticate(c *fiber.Ctx) error {
	// Get authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		// No token provided, continue without authentication
		return c.Next()
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		// Invalid format, continue without authentication
		return c.Next()
	}

	token := parts[1]

	// Validate token
	userID, err := m.tokenGenerator.ValidateToken(token)
	if err != nil {
		// Invalid token, continue without authentication
		return c.Next()
	}

	// Store user ID in context for use in handlers
	c.Locals("userID", userID)

	return c.Next()
}
