package http

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/user/application/port/output"
	"github.com/gofiber/fiber/v2"
)

// UserRoutes configures all routes for the user module
type UserRoutes struct {
	handler    *UserHandler
	middleware *AuthMiddleware
}

// NewUserRoutes creates a new instance of UserRoutes
func NewUserRoutes(userService input.UserService, tokenGenerator output.TokenGenerator) *UserRoutes {
	return &UserRoutes{
		handler:    NewUserHandler(userService),
		middleware: NewAuthMiddleware(tokenGenerator),
	}
}

// SetupRoutes registers all user routes with the Fiber app
func (r *UserRoutes) SetupRoutes(app *fiber.App) {
	// Create a route group for user endpoints
	users := app.Group("/api/users")

	// Public routes (no authentication required)
	users.Post("/register", r.handler.Register)
	users.Post("/login", r.handler.Login)

	// Protected routes (authentication required)
	users.Get("/me", r.middleware.Authenticate, r.handler.GetMe)
	users.Get("/:id", r.middleware.Authenticate, r.handler.GetUserByID)
	users.Put("/:id", r.middleware.Authenticate, r.handler.UpdateUser)
	users.Post("/:id/change-password", r.middleware.Authenticate, r.handler.ChangePassword)
	users.Post("/:id/deactivate", r.middleware.Authenticate, r.handler.DeactivateUser)
	users.Post("/:id/activate", r.middleware.Authenticate, r.handler.ActivateUser)
}
