package input

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type GetAllLocationsUseCase interface {
	Execute(ctx context.Context) ([]entity.Location, error)
}
