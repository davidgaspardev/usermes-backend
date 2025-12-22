package http

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/gofiber/fiber/v2"
)

type OrganizationRoutes struct {
	locationHandler *LocationHandler
}

func NewOrganizationRoutes(
	locationHandler output.LocationRepository,
) *OrganizationRoutes {
	return &OrganizationRoutes{
		locationHandler: NewLocationHandler(locationHandler),
	}
}

func (o *OrganizationRoutes) SetupRoutes(app *fiber.App) {
	locationRoutes := app.Group("/v1/api/organization/locations")

	locationRoutes.Post("/", o.locationHandler.Create)
}
