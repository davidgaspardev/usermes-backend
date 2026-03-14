package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// GetShiftPatternByIDUseCase defines the use case for retrieving a shift pattern by ID.
type GetShiftPatternByIDUseCase interface {
	Execute(ctx context.Context, id uuid.UUID) (*entity.ShiftPattern, error)
}
