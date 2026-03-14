package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

// shiftEntryRecord is the storage representation of a shift entry.
// startTime and endTime are stored as milliseconds since midnight.
type shiftEntryRecord struct {
	Name      string
	StartTime uint64 // ms since midnight
	EndTime   uint64 // ms since midnight
	DayIndex  int
	IsOff     bool
}

// shiftPatternRecord is the storage representation of a shift pattern.
// All time fields are stored as Unix milliseconds.
type shiftPatternRecord struct {
	ID           uuid.UUID
	Name         string
	RefStartDate uint64 // Unix ms
	CycleLength  int
	Entries      []shiftEntryRecord
	CreatedAt    uint64 // Unix ms
}

func timeOfDayToMS(t time.Time) uint64 {
	return uint64(t.Hour())*3600000 + uint64(t.Minute())*60000
}

func msToTimeOfDay(ms uint64) time.Time {
	h := int(ms / 3600000)
	m := int((ms % 3600000) / 60000)
	return time.Date(0, 1, 1, h, m, 0, 0, time.UTC)
}

func toShiftPatternRecord(p *entity.ShiftPattern) shiftPatternRecord {
	entries := make([]shiftEntryRecord, len(p.Entries()))
	for i, e := range p.Entries() {
		entries[i] = shiftEntryRecord{
			Name:      e.Name(),
			StartTime: timeOfDayToMS(e.StartTime()),
			EndTime:   timeOfDayToMS(e.EndTime()),
			DayIndex:  e.DayIndex(),
			IsOff:     e.IsOff(),
		}
	}
	return shiftPatternRecord{
		ID:           p.ID(),
		Name:         p.Name(),
		RefStartDate: uint64(p.RefStartDate().UnixMilli()),
		CycleLength:  p.CycleLength(),
		Entries:      entries,
		CreatedAt:    uint64(p.CreatedAt().UnixMilli()),
	}
}

func fromShiftPatternRecord(rec shiftPatternRecord) *entity.ShiftPattern {
	refStart := time.UnixMilli(int64(rec.RefStartDate)).UTC()
	createdAt := time.UnixMilli(int64(rec.CreatedAt)).UTC()

	p := entity.ReconstructShiftPattern(rec.ID, rec.Name, refStart, rec.CycleLength, createdAt)
	for _, e := range rec.Entries {
		if e.IsOff {
			p.AddEntry(entity.NewDayOffEntry(e.DayIndex))
		} else {
			p.AddEntry(entity.NewShiftEntry(e.DayIndex, e.Name, msToTimeOfDay(e.StartTime), msToTimeOfDay(e.EndTime)))
		}
	}
	return p
}

type memoryShiftPatternRepository struct {
	patterns []shiftPatternRecord
}

// NewMemoryShiftPatternRepository creates a new in-memory ShiftPatternRepository.
func NewMemoryShiftPatternRepository() output.ShiftPatternRepository {
	return &memoryShiftPatternRepository{
		patterns: make([]shiftPatternRecord, 0),
	}
}

func (r *memoryShiftPatternRepository) Save(pattern *entity.ShiftPattern) error {
	r.patterns = append(r.patterns, toShiftPatternRecord(pattern))
	return nil
}

func (r *memoryShiftPatternRepository) Update(pattern *entity.ShiftPattern) error {
	for i, p := range r.patterns {
		if p.ID == pattern.ID() {
			r.patterns[i] = toShiftPatternRecord(pattern)
			return nil
		}
	}
	return domainerrors.ErrShiftPatternNotFound
}

func (r *memoryShiftPatternRepository) Delete(id uuid.UUID) error {
	for i, p := range r.patterns {
		if p.ID == id {
			r.patterns = append(r.patterns[:i], r.patterns[i+1:]...)
			return nil
		}
	}
	return domainerrors.ErrShiftPatternNotFound
}

func (r *memoryShiftPatternRepository) FindByID(id uuid.UUID) (*entity.ShiftPattern, error) {
	for _, p := range r.patterns {
		if p.ID == id {
			return fromShiftPatternRecord(p), nil
		}
	}
	return nil, nil
}

func (r *memoryShiftPatternRepository) GetDefault() (*entity.ShiftPattern, error) {
	if len(r.patterns) == 0 {
		return nil, domainerrors.ErrShiftPatternNotFound
	}
	return fromShiftPatternRecord(r.patterns[0]), nil
}

func (r *memoryShiftPatternRepository) ExistsByName(name string) (bool, error) {
	for _, p := range r.patterns {
		if p.Name == name {
			return true, nil
		}
	}
	return false, nil
}
