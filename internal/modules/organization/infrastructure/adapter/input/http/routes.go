package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
)

// OrganizationRoutes registers HTTP routes for the organization module.
type OrganizationRoutes struct {
	locationHandler     *LocationHandler
	shiftPatternHandler *ShiftPatternHandler
	authMiddleware      fiber.Handler
}

// NewOrganizationRoutes creates a new OrganizationRoutes with the given repositories.
func NewOrganizationRoutes(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
	authMiddleware fiber.Handler,
) *OrganizationRoutes {
	return &OrganizationRoutes{
		locationHandler:     NewLocationHandler(locationRepository, shiftPatternRepository),
		shiftPatternHandler: NewShiftPatternHandler(locationRepository, shiftPatternRepository),
		authMiddleware:      authMiddleware,
	}
}

// SetupRoutes registers all organization routes on the given Fiber app.
func (o *OrganizationRoutes) SetupRoutes(app *fiber.App) {
	locationRoutes := app.Group("/v1/api/organization/locations", o.authMiddleware)

	locationRoutes.Post("/", o.locationHandler.Create)
	locationRoutes.Post("/add", o.locationHandler.Add)
	locationRoutes.Get("/", o.locationHandler.GetAll)
	locationRoutes.Get("/:location_code", o.locationHandler.GetByCode)
	locationRoutes.Post("/:location_code/shift-patterns", o.shiftPatternHandler.Create)

	shiftPatternRoutes := app.Group("/v1/api/organization/shift-patterns", o.authMiddleware)

	shiftPatternRoutes.Get("/:id", o.shiftPatternHandler.GetByID)
	shiftPatternRoutes.Put("/:id", o.shiftPatternHandler.Update)
	shiftPatternRoutes.Delete("/:id", o.shiftPatternHandler.Delete)
}
