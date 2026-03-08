package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
)

type createResourceUseCaseImpl struct {
	resourceRepository output.ResourceRepository
	eventRepository    output.EventRepository
}

// NewCreateResourceUseCase creates a new CreateResourceUseCase implementation.
func NewCreateResourceUseCase(
	resourceRepository output.ResourceRepository,
	eventRepository output.EventRepository,
) input.CreateResourceUseCase {
	return &createResourceUseCaseImpl{
		resourceRepository: resourceRepository,
		eventRepository:    eventRepository,
	}
}

func (u *createResourceUseCaseImpl) Execute(ctx context.Context, command *input.CreateResourceCommand) error {
	if exists, err := u.resourceRepository.ExistsByCode(ctx, command.Code); err != nil {
		return err
	} else if exists {
		return errors.ErrResourceAlreadyExists
	}

	resource := entity.NewResource(
		command.PlantCode,
		command.Code,
		command.ShiftID,
		command.ResourceType,
		command.StopFactor,
		command.Tags,
	)
	event := entity.NewEventResourceCreated(
		resource.Code(),
		*command.ShiftID,
		command.WhoCreated,
	)

	if err := u.resourceRepository.Save(ctx, resource); err != nil {
		return err
	}

	if err := u.eventRepository.Create(event); err != nil {
		if deleteErr := u.resourceRepository.Delete(ctx, resource.ID()); deleteErr != nil {
			return deleteErr
		}
		return err
	}

	return nil
}
