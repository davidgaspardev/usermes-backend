package input

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// GetAllLocationsUseCase defines the use case for retrieving all locations.
type GetAllLocationsUseCase interface {
	Execute(ctx context.Context) ([]entity.Location, error)
}
