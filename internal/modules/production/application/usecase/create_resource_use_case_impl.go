package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
)

type createResourceUseCaseImpl struct {
	resourceRepository output.ResourceRepository
	eventRepository    output.EventRepository
}

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
	return nil
}
