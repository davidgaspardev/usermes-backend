package input

import (
	"context"

	"github.com/google/uuid"
)

// DeleteShiftPatternUseCase defines the use case for deleting a shift pattern.
type DeleteShiftPatternUseCase interface {
	Execute(ctx context.Context, id uuid.UUID) error
}
