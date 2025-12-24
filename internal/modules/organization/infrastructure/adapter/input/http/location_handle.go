package http

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/usecase"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/dto"
	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	locationRepository output.LocationRepository
}

func NewLocationHandler(locationRepository output.LocationRepository) *LocationHandler {
	return &LocationHandler{
		locationRepository: locationRepository,
	}
}

func (h *LocationHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateLocationRootRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
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

func (h *LocationHandler) Add(c *fiber.Ctx) error {
	var req dto.AddLocationRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	command := input.AddLocationCommand{
		Code:       req.Code,
		Name:       req.Name,
		Kind:       req.Kind,
		ParentCode: req.ParentCode,
		RootCode:   req.RootCode,
	}

	location, err := usecase.NewAddLocationUseCase(h.locationRepository).Execute(c.Context(), command)
	if err != nil {
		return err
	}

	response := dto.ToLocationResponse(location)
	return c.Status(fiber.StatusCreated).JSON(response)
}
