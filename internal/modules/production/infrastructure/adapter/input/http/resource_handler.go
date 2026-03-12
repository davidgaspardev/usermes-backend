package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/infrastructure/dto"
)

// ResourceHandler handles HTTP requests for resource operations
type ResourceHandler struct {
	service input.ResourceService
}

// NewResourceHandler creates a new instance of ResourceHandler
func NewResourceHandler(service input.ResourceService) *ResourceHandler {
	return &ResourceHandler{
		service: service,
	}
}

// Create handles the creation of a new resource
func (h *ResourceHandler) Create(c *fiber.Ctx) error {
	locationCode := c.Params("location_code")
	if locationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Location code is required in URL",
		))
	}

	var req dto.CreateResourceRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	resource, err := h.service.Create(c.Context(), locationCode, req.Code, req.Type, req.StopFactor, req.Tags)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToResourceResponse(resource))
}

// GetByID handles retrieving a resource by ID
func (h *ResourceHandler) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_id",
			"Invalid resource ID format",
		))
	}

	resource, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToResourceResponse(resource))
}

// GetByCode handles retrieving a resource by code
func (h *ResourceHandler) GetByCode(c *fiber.Ctx) error {
	locationCode := c.Params("location_code")
	if locationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Location code is required in URL",
		))
	}

	code := c.Params("res_code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_code",
			"Resource code is required",
		))
	}

	resource, err := h.service.GetByCode(c.Context(), locationCode, code)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToResourceResponse(resource))
}

// GetAll handles retrieving all resources with pagination
func (h *ResourceHandler) GetAll(c *fiber.Ctx) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	resources, err := h.service.GetAll(c.Context(), limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToResourceListResponse(resources, limit, offset))
}

// GetByType handles retrieving resources by type
func (h *ResourceHandler) GetByType(c *fiber.Ctx) error {
	resourceType := c.Params("type")
	if resourceType == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_type",
			"Resource type is required",
		))
	}

	return h.getResourcesWithPagination(c, func(limit, offset int) ([]*entity.Resource, error) {
		return h.service.GetByType(c.Context(), resourceType, limit, offset)
	})
}

// getResourcesWithPagination is a helper function to handle pagination logic
func (h *ResourceHandler) getResourcesWithPagination(
	c *fiber.Ctx,
	fetchFunc func(limit, offset int) ([]*entity.Resource, error),
) error {
	limit := c.QueryInt("limit", 10)
	offset := c.QueryInt("offset", 0)

	resources, err := fetchFunc(limit, offset)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToResourceListResponse(resources, limit, offset))
}

// Update handles updating an existing resource
func (h *ResourceHandler) Update(c *fiber.Ctx) error {
	locationCode := c.Params("location_code")
	if locationCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Location code is required in URL",
		))
	}

	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid resource ID",
		))
	}

	var req dto.UpdateResourceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	resource, err := h.service.Update(c.Context(), id, locationCode, req.Code, req.Type, req.StopFactor, req.Tags)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToResourceResponse(resource))
}

// Delete handles deleting a resource
func (h *ResourceHandler) Delete(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_id",
			"Invalid resource ID format",
		))
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(dto.NewSuccessResponse(
		"Resource deleted successfully",
		nil,
	))
}

// handleError maps domain errors to HTTP responses
func (h *ResourceHandler) handleError(c *fiber.Ctx, err error) error {
	switch err {
	case errors.ErrResourceNotFound:
		return c.Status(fiber.StatusNotFound).JSON(dto.NewErrorResponse(
			"resource_not_found",
			"Resource not found",
		))
	case errors.ErrResourceAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(dto.NewErrorResponse(
			"resource_already_exists",
			"Resource already exists",
		))
	case errors.ErrCodeAlreadyExists:
		return c.Status(fiber.StatusConflict).JSON(dto.NewErrorResponse(
			"code_already_exists",
			"A resource with this code already exists",
		))
	case errors.ErrInvalidCode:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_code",
			"Invalid resource code: must be between 2 and 50 characters",
		))
	case errors.ErrInvalidType:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_type",
			"Invalid resource type: must be between 2 and 50 characters",
		))
	case errors.ErrInvalidStopFactor:
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_stop_factor",
			"Invalid stop factor: must be non-negative",
		))
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(dto.NewErrorResponse(
			"internal_error",
			"An internal error occurred",
		))
	}
}
