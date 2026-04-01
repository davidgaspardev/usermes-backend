package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type createShiftPatternUseCaseImpl struct {
	locationRepository     output.LocationRepository
	shiftPatternRepository output.ShiftPatternRepository
}

// NewCreateShiftPatternUseCase creates a new CreateShiftPatternUseCase instance.
func NewCreateShiftPatternUseCase(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
) input.CreateShiftPatternUseCase {
	return &createShiftPatternUseCaseImpl{
		locationRepository:     locationRepository,
		shiftPatternRepository: shiftPatternRepository,
	}
}

func (uc *createShiftPatternUseCaseImpl) Execute(ctx context.Context, command input.CreateShiftPatternCommand) (*entity.ShiftPattern, error) {
	if command.LocationCode == "" || command.Name == "" || command.CycleLength < 1 || len(command.Entries) == 0 {
		return nil, domainerrors.ErrInvalidShiftPatternCommand
	}

	location, err := uc.locationRepository.FindByCode(command.LocationCode)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, domainerrors.ErrLocationNotFound
	}

	exists, err := uc.shiftPatternRepository.ExistsByName(command.Name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domainerrors.ErrShiftPatternAlreadyExists
	}

	pattern := entity.NewShiftPattern(command.Name, command.RefStartDate, command.CycleLength)

	for _, e := range command.Entries {
		pattern.AddEntry(entity.NewShiftEntry(e.DayIndex, e.Name, e.StartTime, e.EndTime))
	}

	if err := uc.shiftPatternRepository.Save(pattern); err != nil {
		return nil, err
	}

	location.AssignShiftPattern(pattern.ID())

	if err := uc.locationRepository.Update(location); err != nil {
		return nil, err
	}

	return pattern, nil
}
