package output

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// ShiftPatternRepository defines the persistence contract for shift patterns.
type ShiftPatternRepository interface {
	// Save persists a new shift pattern.
	Save(pattern *entity.ShiftPattern) error

	// FindByID retrieves a shift pattern by its unique identifier.
	FindByID(id uuid.UUID) (*entity.ShiftPattern, error)

	// GetAll returns all shift patterns.
	GetAll() ([]*entity.ShiftPattern, error)

	// ExistsByName reports whether a shift pattern with the given name already exists.
	ExistsByName(name string) (bool, error)
}
