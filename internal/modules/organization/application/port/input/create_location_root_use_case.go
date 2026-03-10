package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// CreateLocationRootCommand holds the input data for creating a root location.
// ShiftPatternID is optional; when nil the default seeded pattern is used.
type CreateLocationRootCommand struct {
	ShiftPatternID *uuid.UUID
	Code           string
	Name           string
}

// CreateLocationRootUseCase defines the use case for creating a root location.
type CreateLocationRootUseCase interface {
	Execute(ctx context.Context, command CreateLocationRootCommand) (*entity.Location, error)
}
