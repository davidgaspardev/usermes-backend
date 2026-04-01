package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// CreateLocationRootCommand holds the input data for creating a root location.
// Plants do not carry a shift pattern — only sub-locations do.
type CreateLocationRootCommand struct {
	Code      string
	Name      string
	CreatedBy uuid.UUID
}

// CreateLocationRootUseCase defines the use case for creating a root location.
type CreateLocationRootUseCase interface {
	Execute(ctx context.Context, command CreateLocationRootCommand) (*entity.Location, error)
}
