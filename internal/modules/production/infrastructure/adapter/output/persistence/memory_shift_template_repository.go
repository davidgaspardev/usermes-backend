package persistence

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/google/uuid"
)

type memoryShiftTemplateRepository struct {
	templates []*entity.ShiftTemplate
}

// NewMemoryShiftTemplateRepository creates a new in-memory ShiftTemplateRepository.
func NewMemoryShiftTemplateRepository() output.ShiftTemplateRepository {
	return &memoryShiftTemplateRepository{
		templates: make([]*entity.ShiftTemplate, 0),
	}
}

func (r *memoryShiftTemplateRepository) Save(template *entity.ShiftTemplate) error {
	r.templates = append(r.templates, template)
	return nil
}

func (r *memoryShiftTemplateRepository) FindByID(id uuid.UUID) (*entity.ShiftTemplate, error) {
	for _, t := range r.templates {
		if t.ID() == id {
			return t, nil
		}
	}
	return nil, nil
}

func (r *memoryShiftTemplateRepository) GetAll() ([]*entity.ShiftTemplate, error) {
	result := make([]*entity.ShiftTemplate, len(r.templates))
	copy(result, r.templates)
	return result, nil
}

func (r *memoryShiftTemplateRepository) ExistsByName(name string) (bool, error) {
	for _, t := range r.templates {
		if t.Name() == name {
			return true, nil
		}
	}
	return false, nil
}
