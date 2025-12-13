package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/application/port/output"
)

// IAM configures all routes for the user module
type IamRoutes struct {
	userHandler *UserHandler
	middleware  *AuthMiddleware
}

// NewIamRoutes creates a new instance of UserRoutes
func NewIamRoutes(userService input.UserService, tokenGenerator output.TokenGenerator) *IamRoutes {
	return &IamRoutes{
		userHandler: NewUserHandler(userService),
		middleware:  NewAuthMiddleware(tokenGenerator),
	}
}

// SetupRoutes registers all user routes with the Fiber app
func (r *IamRoutes) SetupRoutes(app *fiber.App) {
	// Create a route group for user endpoints
	iam := app.Group("/v1/api/iam")

	// Public routes (no authentication required)
	iam.Post("/user/register", r.userHandler.Register)
	iam.Post("/user/login", r.userHandler.Login)

	// Protected routes (authentication required)
	iam.Get("/user/me", r.middleware.Authenticate, r.userHandler.GetMe)
	iam.Get("/user/:id", r.middleware.Authenticate, r.userHandler.GetUserByID)
	iam.Put("/user/:id", r.middleware.Authenticate, r.userHandler.UpdateUser)
	iam.Post("/user/:id/change-password", r.middleware.Authenticate, r.userHandler.ChangePassword)
	iam.Post("/user/:id/deactivate", r.middleware.Authenticate, r.userHandler.DeactivateUser)
	iam.Post("/user/:id/activate", r.middleware.Authenticate, r.userHandler.ActivateUser)
}
