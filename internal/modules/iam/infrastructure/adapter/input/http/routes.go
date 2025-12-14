package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/application/port/output"
)

// IAMRoutes configures all routes for the IAM module
type IAMRoutes struct {
	userHandler *UserHandler
	middleware  *AuthMiddleware
}

// NewIAMRoutes creates a new instance of IAMRoutes
func NewIAMRoutes(userService input.UserService, tokenGenerator output.TokenGenerator) *IAMRoutes {
	return &IAMRoutes{
		userHandler: NewUserHandler(userService),
		middleware:  NewAuthMiddleware(tokenGenerator),
	}
}

// SetupRoutes registers all IAM routes with the Fiber app
func (r *IAMRoutes) SetupRoutes(app *fiber.App) {
	iam := app.Group("/v1/api/iam")

	// Public routes (no authentication required)
	iam.Post("/users/register", r.userHandler.Register)
	iam.Post("/users/login", r.userHandler.Login)

	// Protected routes (authentication required)
	iam.Get("/users/me", r.middleware.Authenticate, r.userHandler.GetMe)
	iam.Get("/users/:id", r.middleware.Authenticate, r.userHandler.GetUserByID)
	iam.Put("/users/:id", r.middleware.Authenticate, r.userHandler.UpdateUser)
	iam.Post("/users/:id/change-password", r.middleware.Authenticate, r.userHandler.ChangePassword)
	iam.Post("/users/:id/deactivate", r.middleware.Authenticate, r.userHandler.DeactivateUser)
	iam.Post("/users/:id/activate", r.middleware.Authenticate, r.userHandler.ActivateUser)
}
