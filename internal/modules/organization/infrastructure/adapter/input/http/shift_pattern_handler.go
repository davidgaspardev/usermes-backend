package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/usecase"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/dto"
)

// ShiftPatternHandler handles HTTP requests for shift pattern operations.
type ShiftPatternHandler struct {
	locationRepository     output.LocationRepository
	shiftPatternRepository output.ShiftPatternRepository
}

// NewShiftPatternHandler creates a new ShiftPatternHandler.
func NewShiftPatternHandler(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
) *ShiftPatternHandler {
	return &ShiftPatternHandler{
		locationRepository:     locationRepository,
		shiftPatternRepository: shiftPatternRepository,
	}
}

// Create handles POST /locations/:location_code/shift-patterns.
// Creates a shift pattern and immediately links it to the given location.
func (h *ShiftPatternHandler) Create(c *fiber.Ctx) error {
	locationCode := c.Params("location_code")

	var req dto.CreateShiftPatternRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	refStartDate, err := req.ParseRefStartDate()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			"ref_start_date must be in YYYY-MM-DD format",
		))
	}

	entries, err := parseEntryCommands(req.Entries)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	command := input.CreateShiftPatternCommand{
		LocationCode: locationCode,
		Name:         req.Name,
		RefStartDate: refStartDate,
		CycleLength:  req.CycleLength,
		Entries:      entries,
	}

	pattern, err := usecase.NewCreateShiftPatternUseCase(h.locationRepository, h.shiftPatternRepository).Execute(c.Context(), command)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ToShiftPatternResponse(pattern))
}

// GetByID handles GET /shift-patterns/:id.
func (h *ShiftPatternHandler) GetByID(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_id",
			"Invalid shift pattern ID",
		))
	}

	pattern, err := usecase.NewGetShiftPatternByIDUseCase(h.shiftPatternRepository).Execute(c.Context(), id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToShiftPatternResponse(pattern))
}

// Update handles PUT /shift-patterns/:id.
func (h *ShiftPatternHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_id",
			"Invalid shift pattern ID",
		))
	}

	var req dto.UpdateShiftPatternRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_request",
			"Invalid request body",
		))
	}

	if err := req.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	entries, err := parseEntryCommands(req.Entries)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"validation_error",
			err.Error(),
		))
	}

	command := input.UpdateShiftPatternCommand{
		ID:          id,
		Name:        req.Name,
		CycleLength: req.CycleLength,
		Entries:     entries,
	}

	pattern, err := usecase.NewUpdateShiftPatternUseCase(h.shiftPatternRepository).Execute(c.Context(), command)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(dto.ToShiftPatternResponse(pattern))
}

// Delete handles DELETE /shift-patterns/:id.
func (h *ShiftPatternHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(dto.NewErrorResponse(
			"invalid_id",
			"Invalid shift pattern ID",
		))
	}

	if err := usecase.NewDeleteShiftPatternUseCase(h.shiftPatternRepository).Execute(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func parseID(c *fiber.Ctx) (uuid.UUID, error) {
	return uuid.Parse(c.Params("id"))
}

// parseEntryCommands converts DTO entry requests into use-case commands,
// parsing "HH:MM" strings into time.Time values at the HTTP boundary.
func parseEntryCommands(entries []dto.ShiftEntryRequest) ([]input.ShiftEntryCommand, error) {
	cmds := make([]input.ShiftEntryCommand, len(entries))
	for i, e := range entries {
		startTime, err := time.Parse("15:04", e.StartTime)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "each entry start_time must be in HH:MM format")
		}
		endTime, err := time.Parse("15:04", e.EndTime)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "each entry end_time must be in HH:MM format")
		}
		cmds[i] = input.ShiftEntryCommand{
			DayIndex:  e.DayIndex,
			Name:      e.Name,
			StartTime: startTime,
			EndTime:   endTime,
		}
	}
	return cmds, nil
}
