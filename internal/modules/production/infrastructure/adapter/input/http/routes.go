package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
)

// ProductionRoutes configures all routes for the production module
type ProductionRoutes struct {
	resourceHandler *ResourceHandler
}

// NewProductionRoutes creates a new instance of ProductionRoutes
func NewProductionRoutes(resourceService input.ResourceService) *ProductionRoutes {
	return &ProductionRoutes{
		resourceHandler: NewResourceHandler(resourceService),
	}
}

// SetupRoutes registers all production routes with the Fiber app
func (r *ProductionRoutes) SetupRoutes(app *fiber.App) {
	resources := app.Group("/v1/api/production/plants/:plant_code")

	// Create a new resource in the specified plant
	resources.Post("/resources/", r.resourceHandler.Create)
	resources.Get("/resources/", r.resourceHandler.GetAll)
	resources.Get("/resources/:res_code", r.resourceHandler.GetByCode)
	resources.Get("/resources/:type", r.resourceHandler.GetByType)
	resources.Get("/resources/shift/:shiftId", r.resourceHandler.GetByShiftID)
	resources.Put("/resources/:id", r.resourceHandler.Update)
	resources.Delete("/resources/:id", r.resourceHandler.Delete)
}
