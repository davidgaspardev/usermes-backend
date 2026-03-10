package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type createLocationRootUseCaseImpl struct {
	locationRepository     output.LocationRepository
	shiftPatternRepository output.ShiftPatternRepository
}

// NewCreateLocationRootUseCase creates a new CreateLocationRootUseCase instance.
func NewCreateLocationRootUseCase(
	locationRepository output.LocationRepository,
	shiftPatternRepository output.ShiftPatternRepository,
) input.CreateLocationRootUseCase {
	return &createLocationRootUseCaseImpl{
		locationRepository:     locationRepository,
		shiftPatternRepository: shiftPatternRepository,
	}
}

func (uc *createLocationRootUseCaseImpl) Execute(ctx context.Context, command input.CreateLocationRootCommand) (*entity.Location, error) {
	locationExists, err := uc.locationRepository.ExistsByCode(command.Code)
	if err != nil {
		return nil, err
	}

	if locationExists {
		return nil, errors.ErrLocationAlreadyExists
	}

	shiftPatternID, err := resolveShiftPatternID(command.ShiftPatternID, uc.shiftPatternRepository)
	if err != nil {
		return nil, err
	}

	location := entity.NewRootLocation(command.Code, command.Name, shiftPatternID)
	if err := uc.locationRepository.Create(location); err != nil {
		return nil, err
	}

	return location, nil
}

// resolveShiftPatternID returns the given ID when provided, otherwise fetches the default pattern.
func resolveShiftPatternID(id *uuid.UUID, repo output.ShiftPatternRepository) (uuid.UUID, error) {
	if id != nil {
		return *id, nil
	}
	pattern, err := repo.GetDefault()
	if err != nil {
		return uuid.UUID{}, err
	}
	return pattern.ID(), nil
}
