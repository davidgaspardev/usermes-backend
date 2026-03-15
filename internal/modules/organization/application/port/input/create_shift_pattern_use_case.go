package input

import (
	"context"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// ShiftEntryCommand holds input data for a single shift entry.
// Days with no entries are implicitly off.
type ShiftEntryCommand struct {
	StartTime time.Time
	EndTime   time.Time
	Name      string
	DayIndex  int
}

// CreateShiftPatternCommand holds input data for creating a shift pattern.
// LocationCode is required: the pattern is created and immediately linked to that location.
// CycleLength must be >= 1 and > the highest DayIndex in Entries.
type CreateShiftPatternCommand struct {
	RefStartDate time.Time
	LocationCode string
	Name         string
	Entries      []ShiftEntryCommand
	CycleLength  int
}

// CreateShiftPatternUseCase defines the use case for creating a shift pattern.
type CreateShiftPatternUseCase interface {
	Execute(ctx context.Context, command CreateShiftPatternCommand) (*entity.ShiftPattern, error)
}
