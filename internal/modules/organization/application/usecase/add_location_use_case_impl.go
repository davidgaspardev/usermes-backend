package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type addLocationUseCaseImpl struct {
	locationRepository output.LocationRepository
}

func NewAddLocationUseCase(locationRepository output.LocationRepository) input.AddLocationUseCase {
	return &addLocationUseCaseImpl{
		locationRepository: locationRepository,
	}
}

func (uc *addLocationUseCaseImpl) Execute(ctx context.Context, command input.AddLocationCommand) (*entity.Location, error) {
	if !isCommandValid(&command) {
		return nil, errors.ErrInvalidLocationCommand
	}

	locationTree, err := uc.locationRepository.FindTree(command.RootCode)
	if err != nil {
		return nil, err
	}

	if locationTree == nil {
		return nil, errors.ErrLocationTreeNotFound
	}

	if locationAlreadyExists := locationTree.ExistsByLocation(command.Code); locationAlreadyExists {
		return nil, errors.ErrLocationAlreadyExists
	}

	parentLocation := locationTree.FindByCode(command.ParentCode)
	if parentLocation == nil {
		return nil, errors.ErrLocationNotFound
	}

	location := entity.NewLocation(
		command.Code,
		command.Name,
		parentLocation,
	)
	parentLocation.AddChild(location)

	if err := uc.locationRepository.Create(location); err != nil {
		return nil, err
	}

	return locationTree, nil
}

func isCommandValid(command *input.AddLocationCommand) bool {
	return command.Code != "" &&
		command.Name != "" &&
		command.ParentCode != "" &&
		command.RootCode != "" &&
		command.Code != command.RootCode &&
		command.Code != command.ParentCode
}
