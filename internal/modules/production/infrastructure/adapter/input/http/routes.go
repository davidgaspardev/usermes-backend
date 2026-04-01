package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
)

// ProductionRoutes configures all routes for the production module
type ProductionRoutes struct {
	resourceHandler        *ResourceHandler
	authMiddleware         fiber.Handler
	locationCodeMiddleware fiber.Handler
}

// NewProductionRoutes creates a new instance of ProductionRoutes
func NewProductionRoutes(resourceService input.ResourceService, authMiddleware fiber.Handler, locationCodeMiddleware fiber.Handler) *ProductionRoutes {
	return &ProductionRoutes{
		resourceHandler:        NewResourceHandler(resourceService),
		authMiddleware:         authMiddleware,
		locationCodeMiddleware: locationCodeMiddleware,
	}
}

// SetupRoutes registers all production routes with the Fiber app
func (r *ProductionRoutes) SetupRoutes(app *fiber.App) {
	resources := app.Group("/v1/api/production/locations/:location_code", r.authMiddleware, r.locationCodeMiddleware)

	resources.Post("/resources/", r.resourceHandler.Create)
	resources.Get("/resources/", r.resourceHandler.GetAll)
	resources.Get("/resources/:id", r.resourceHandler.GetByID)
	resources.Get("/resources/code/:res_code", r.resourceHandler.GetByCode)
	resources.Get("/resources/type/:type", r.resourceHandler.GetByType)
	resources.Put("/resources/:id", r.resourceHandler.Update)
	resources.Delete("/resources/:id", r.resourceHandler.Delete)
}
