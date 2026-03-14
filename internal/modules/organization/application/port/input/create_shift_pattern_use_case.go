package input

import (
	"context"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// ShiftEntryCommand holds input data for a single shift entry.
// Days with no entries are implicitly off.
type ShiftEntryCommand struct {
	Name      string
	StartTime time.Time // time-of-day
	EndTime   time.Time // time-of-day
	DayIndex  int       // 0-based day position within the cycle
}

// CreateShiftPatternCommand holds input data for creating a shift pattern.
// LocationCode is required: the pattern is created and immediately linked to that location.
// CycleLength must be >= 1 and > the highest DayIndex in Entries.
type CreateShiftPatternCommand struct {
	LocationCode string
	Name         string
	RefStartDate time.Time
	CycleLength  int
	Entries      []ShiftEntryCommand
}

// CreateShiftPatternUseCase defines the use case for creating a shift pattern.
type CreateShiftPatternUseCase interface {
	Execute(ctx context.Context, command CreateShiftPatternCommand) (*entity.ShiftPattern, error)
}
