package output

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
)

// ShiftRepository defines the persistence contract for shift instances.
type ShiftRepository interface {
	// Save persists a new shift instance.
	Save(shift *entity.Shift) error

	// FindByID retrieves a shift instance by its unique identifier.
	FindByID(id uuid.UUID) (*entity.Shift, error)
}
