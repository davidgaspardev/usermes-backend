package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type updateShiftPatternUseCaseImpl struct {
	shiftPatternRepository output.ShiftPatternRepository
}

// NewUpdateShiftPatternUseCase creates a new UpdateShiftPatternUseCase instance.
func NewUpdateShiftPatternUseCase(repo output.ShiftPatternRepository) input.UpdateShiftPatternUseCase {
	return &updateShiftPatternUseCaseImpl{shiftPatternRepository: repo}
}

func (uc *updateShiftPatternUseCaseImpl) Execute(ctx context.Context, command input.UpdateShiftPatternCommand) (*entity.ShiftPattern, error) {
	if command.Name == "" || command.CycleLength < 1 || len(command.Entries) == 0 || command.RefStartDate.IsZero() {
		return nil, domainerrors.ErrInvalidShiftPatternCommand
	}

	pattern, err := uc.shiftPatternRepository.FindByID(command.ID)
	if err != nil {
		return nil, err
	}
	if pattern == nil {
		return nil, domainerrors.ErrShiftPatternNotFound
	}

	pattern.SetName(command.Name)
	pattern.SetRefStartDate(command.RefStartDate)
	pattern.SetCycleLength(command.CycleLength)

	entries := make([]*entity.ShiftEntry, len(command.Entries))
	for i, e := range command.Entries {
		entries[i] = entity.NewShiftEntry(e.DayIndex, e.Name, e.StartTime, e.EndTime)
	}
	pattern.SetEntries(entries)

	if err := uc.shiftPatternRepository.Update(pattern); err != nil {
		return nil, err
	}

	return pattern, nil
}
