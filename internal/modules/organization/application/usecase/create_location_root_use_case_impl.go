package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type createLocationRootUseCaseImpl struct {
	locationRepository output.LocationRepository
}

// NewCreateLocationRootUseCase creates a new CreateLocationRootUseCase instance.
func NewCreateLocationRootUseCase(
	locationRepository output.LocationRepository,
) input.CreateLocationRootUseCase {
	return &createLocationRootUseCaseImpl{
		locationRepository: locationRepository,
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

	location := entity.NewRootLocation(
		command.Code,
		command.Name,
	)
	if err := uc.locationRepository.Create(location); err != nil {
		return nil, err
	}

	return location, nil
}
