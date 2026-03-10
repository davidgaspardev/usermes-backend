package persistence

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type memoryShiftPatternRepository struct {
	patterns []*entity.ShiftPattern
}

// NewMemoryShiftPatternRepository creates a new in-memory ShiftPatternRepository.
func NewMemoryShiftPatternRepository() output.ShiftPatternRepository {
	return &memoryShiftPatternRepository{
		patterns: make([]*entity.ShiftPattern, 0),
	}
}

func (r *memoryShiftPatternRepository) Save(pattern *entity.ShiftPattern) error {
	r.patterns = append(r.patterns, pattern)
	return nil
}

func (r *memoryShiftPatternRepository) FindByID(id uuid.UUID) (*entity.ShiftPattern, error) {
	for _, p := range r.patterns {
		if p.ID() == id {
			return p, nil
		}
	}
	return nil, nil
}

func (r *memoryShiftPatternRepository) GetDefault() (*entity.ShiftPattern, error) {
	if len(r.patterns) == 0 {
		return nil, domainerrors.ErrShiftPatternNotFound
	}
	return r.patterns[0], nil
}

func (r *memoryShiftPatternRepository) GetAll() ([]*entity.ShiftPattern, error) {
	result := make([]*entity.ShiftPattern, len(r.patterns))
	copy(result, r.patterns)
	return result, nil
}

func (r *memoryShiftPatternRepository) ExistsByName(name string) (bool, error) {
	for _, p := range r.patterns {
		if p.Name() == name {
			return true, nil
		}
	}
	return false, nil
}
