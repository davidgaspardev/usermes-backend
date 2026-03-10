package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/usecase"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/dto"
)

// LocationHandler handles HTTP requests for location operations.
type LocationHandler struct {
	locationRepository     output.LocationRepository
	shiftPatternRepository output.ShiftPatternRepository
}

// NewLocationHandler creates a new LocationHandler with the given repositories.
func NewLocationHandler(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
) *LocationHandler {
	return &LocationHandler{
		locationRepository:     locationRepository,
		shiftPatternRepository: shiftPatternRepository,
	}
}

// Create handles POST requests to create a root location.
func (h *LocationHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateLocationRootRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	command := input.CreateLocationRootCommand{
		Code: req.Code,
		Name: req.Name,
	}

	locationRoot, err := usecase.NewCreateLocationRootUseCase(h.locationRepository).Execute(c.Context(), command)
	if err != nil {
		return err
	}

	response := dto.ToLocationResponse(locationRoot)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// Add handles POST requests to add a child location to an existing tree.
func (h *LocationHandler) Add(c *fiber.Ctx) error {
	var req dto.AddLocationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	// Validate request
	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	command := input.AddLocationCommand{
		Code:       req.Code,
		Name:       req.Name,
		Kind:       req.Kind,
		ParentCode: req.ParentCode,
		RootCode:   req.RootCode,
	}

	location, err := usecase.NewAddLocationUseCase(h.locationRepository, h.shiftPatternRepository).Execute(c.Context(), command)
	if err != nil {
		return err
	}

	response := dto.ToLocationResponse(location)
	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetAll handles GET requests to retrieve all locations.
func (h *LocationHandler) GetAll(c *fiber.Ctx) error {
	locations, err := usecase.NewGetAllLocationsUseCase(h.locationRepository).Execute(c.Context())
	if err != nil {
		return err
	}

	response := dto.ToAllLocationsResponse(locations)
	return c.Status(fiber.StatusOK).JSON(response)
}
