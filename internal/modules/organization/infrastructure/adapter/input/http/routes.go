package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
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
	locationRoutes.Post("/add", o.locationHandler.Add)
	locationRoutes.Get("/", o.locationHandler.GetAll)
}
