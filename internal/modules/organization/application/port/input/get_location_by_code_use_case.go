package input

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// GetLocationByCodeUseCase defines the use case for retrieving a location tree by its root code.
type GetLocationByCodeUseCase interface {
	Execute(ctx context.Context, code string) (*entity.Location, error)
}
