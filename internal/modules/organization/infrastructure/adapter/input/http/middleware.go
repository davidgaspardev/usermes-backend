package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/dto"
)

// LocationMiddleware provides middleware for location-related validations.
type LocationMiddleware struct {
	locationRepository output.LocationRepository
}

// NewLocationMiddleware creates a new LocationMiddleware instance.
func NewLocationMiddleware(locationRepository output.LocationRepository) *LocationMiddleware {
	return &LocationMiddleware{locationRepository: locationRepository}
}

// VerifyLocationCode checks that the :location_code route param refers to an existing location.
func (m *LocationMiddleware) VerifyLocationCode(c *fiber.Ctx) error {
	code := c.Params("location_code")

	exists, err := m.locationRepository.ExistsByCode(code)
	if err != nil {
		return err
	}

	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse(
			"not_found",
			"Location not found",
		))
	}

	return c.Next()
}
