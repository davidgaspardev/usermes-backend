package input

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// UpdateShiftPatternCommand holds input data for updating a shift pattern.
type UpdateShiftPatternCommand struct {
	RefStartDate time.Time
	Name         string
	Entries      []ShiftEntryCommand
	CycleLength  int
	ID           uuid.UUID
}

// UpdateShiftPatternUseCase defines the use case for updating a shift pattern.
type UpdateShiftPatternUseCase interface {
	Execute(ctx context.Context, command UpdateShiftPatternCommand) (*entity.ShiftPattern, error)
}
