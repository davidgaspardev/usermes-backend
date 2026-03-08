package output

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/google/uuid"
)

// ShiftRepository defines the persistence contract for shift instances.
type ShiftRepository interface {
	// Save persists a new shift instance.
	Save(shift *entity.Shift) error

	// FindByID retrieves a shift instance by its unique identifier.
	FindByID(id uuid.UUID) (*entity.Shift, error)
}
