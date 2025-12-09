package site

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/infrastructure/dto"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SiteHandler struct {
	service input.SiteService
}

func NewSiteHandler(service input.SiteService) *SiteHandler {
	return &SiteHandler{
		service: service,
	}
}

func (h *SiteHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateSiteRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	userID, err := uuid.NewUUID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse(
			"internal_error",
			"Failed to generate UUID",
		))
	}

	site, err := h.service.Create(c.Context(), req.Code, req.Name, req.Latitude, req.Longitude, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse(
			"internal_error",
			"Failed to create site",
		))
	}

	response := dto.ToSiteResponse(site)
	return c.Status(fiber.StatusCreated).JSON(response)
}
