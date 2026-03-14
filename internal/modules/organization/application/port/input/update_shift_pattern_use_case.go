package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// UpdateShiftPatternCommand holds input data for updating a shift pattern.
type UpdateShiftPatternCommand struct {
	ID          uuid.UUID
	Name        string
	CycleLength int
	Entries     []ShiftEntryCommand
}

// UpdateShiftPatternUseCase defines the use case for updating a shift pattern.
type UpdateShiftPatternUseCase interface {
	Execute(ctx context.Context, command UpdateShiftPatternCommand) (*entity.ShiftPattern, error)
}
