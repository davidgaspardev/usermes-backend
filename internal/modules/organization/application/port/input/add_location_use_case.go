package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// AddLocationCommand holds the input data for adding a location to an existing tree.
// ShiftPatternID is optional; when nil the default seeded pattern is used.
type AddLocationCommand struct {
	ShiftPatternID *uuid.UUID
	Code           string
	Name           string
	Kind           string
	ParentCode     string
	RootCode       string
}

// AddLocationUseCase defines the use case for adding a child location.
type AddLocationUseCase interface {
	Execute(ctx context.Context, command AddLocationCommand) (*entity.Location, error)
}
