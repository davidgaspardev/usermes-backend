package output

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// ShiftPatternRepository defines the persistence contract for shift patterns.
type ShiftPatternRepository interface {
	// Save persists a new shift pattern.
	Save(pattern *entity.ShiftPattern) error

	// Update persists changes to an existing shift pattern.
	Update(pattern *entity.ShiftPattern) error

	// Delete removes a shift pattern by its unique identifier.
	Delete(id uuid.UUID) error

	// FindByID retrieves a shift pattern by its unique identifier.
	FindByID(id uuid.UUID) (*entity.ShiftPattern, error)

	// GetDefault returns the default shift pattern (the first seeded one).
	// Returns ErrShiftPatternNotFound when no patterns have been seeded.
	GetDefault() (*entity.ShiftPattern, error)

	// ExistsByName reports whether a shift pattern with the given name already exists.
	ExistsByName(name string) (bool, error)
}
