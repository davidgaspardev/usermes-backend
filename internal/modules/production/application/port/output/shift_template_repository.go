package output

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/google/uuid"
)

// ShiftTemplateRepository defines the persistence contract for shift templates.
type ShiftTemplateRepository interface {
	// Save persists a new shift template.
	Save(template *entity.ShiftTemplate) error

	// FindByID retrieves a shift template by its unique identifier.
	FindByID(id uuid.UUID) (*entity.ShiftTemplate, error)

	// GetAll returns all shift templates.
	GetAll() ([]*entity.ShiftTemplate, error)

	// ExistsByName reports whether a shift template with the given name already exists.
	ExistsByName(name string) (bool, error)
}
