package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
)

// OrganizationRoutes registers HTTP routes for the organization module.
type OrganizationRoutes struct {
	locationHandler *LocationHandler
}

// NewOrganizationRoutes creates a new OrganizationRoutes with the given repositories.
func NewOrganizationRoutes(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
) *OrganizationRoutes {
	return &OrganizationRoutes{
		locationHandler: NewLocationHandler(locationRepository, shiftPatternRepository),
	}
}

// SetupRoutes registers all organization routes on the given Fiber app.
func (o *OrganizationRoutes) SetupRoutes(app *fiber.App) {
	locationRoutes := app.Group("/v1/api/organization/locations")

	locationRoutes.Post("/", o.locationHandler.Create)
	locationRoutes.Post("/add", o.locationHandler.Add)
	locationRoutes.Get("/", o.locationHandler.GetAll)
	locationRoutes.Get("/:location_code", o.locationHandler.GetByCode)
}
