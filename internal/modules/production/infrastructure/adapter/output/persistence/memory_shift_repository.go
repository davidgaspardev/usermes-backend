package persistence

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
)

type memoryShiftRepository struct {
	shifts []*entity.Shift
}

// NewMemoryShiftRepository creates a new in-memory ShiftRepository.
func NewMemoryShiftRepository() output.ShiftRepository {
	return &memoryShiftRepository{
		shifts: make([]*entity.Shift, 0),
	}
}

func (r *memoryShiftRepository) Save(shift *entity.Shift) error {
	r.shifts = append(r.shifts, shift)
	return nil
}

func (r *memoryShiftRepository) FindByID(id uuid.UUID) (*entity.Shift, error) {
	for _, s := range r.shifts {
		if s.ID() == id {
			return s, nil
		}
	}
	return nil, nil
}
