package input

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type CreateLocationRootCommand struct {
	Code string
	Name string
}

type CreateLocationRootUseCase interface {
	Execute(ctx context.Context, command CreateLocationRootCommand) (*entity.Location, error)
}
