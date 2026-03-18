package persistence

import (
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

// shiftEntryRecord is the storage representation of a shift entry.
// StartTime is stored as Unix milliseconds; Duration is the shift length in ms.
type shiftEntryRecord struct {
	Name      string
	StartTime uint64 // Unix ms
	Duration  uint64 // shift length in ms
	DayIndex  int
	IsOff     bool
}

// shiftPatternRecord is the storage representation of a shift pattern.
// All time fields are stored as Unix milliseconds.
type shiftPatternRecord struct {
	Name         string
	Entries      []shiftEntryRecord
	RefStartDate uint64
	CycleLength  int
	CreatedAt    uint64
	ID           uuid.UUID
}

func timeOfDayToMS(t time.Time) uint64 {
	return uint64(t.Hour())*3_600_000 + uint64(t.Minute())*60_000 + uint64(t.Second())*1_000 //nolint:gosec
}

func msToTimeOfDay(ms uint64) time.Time {
	h := int(ms / 3_600_000)             //nolint:gosec
	m := int((ms % 3_600_000) / 60_000)  //nolint:gosec
	s := int((ms % 60_000) / 1_000)      //nolint:gosec
	return time.Date(0, 1, 1, h, m, s, 0, time.UTC)
}

// durationMS computes the duration between two time-of-day values in milliseconds.
// Overnight shifts (end < start) are handled correctly by wrapping around 24 h.
// A full-day shift (start == end) is treated as exactly 24 h.
func durationMS(start, end time.Time) uint64 {
	const dayMS = int64(24 * 3_600_000)
	startMS := int64(timeOfDayToMS(start)) //nolint:gosec
	endMS := int64(timeOfDayToMS(end))     //nolint:gosec
	dur := (endMS - startMS + dayMS) % dayMS
	if dur == 0 {
		dur = dayMS
	}
	return uint64(dur) //nolint:gosec
}

// endFromDuration reconstructs the end datetime from a start datetime and duration in ms.
// Overnight shifts advance the date by one day.
func endFromDuration(start time.Time, dur uint64) time.Time {
	return start.Add(time.Duration(dur) * time.Millisecond) //nolint:gosec
}

func toShiftPatternRecord(p *entity.ShiftPattern) shiftPatternRecord {
	entries := make([]shiftEntryRecord, len(p.Entries()))
	for i, e := range p.Entries() {
		entries[i] = shiftEntryRecord{
			Name:      e.Name(),
			StartTime: uint64(e.StartTime().UnixMilli()), //nolint:gosec
			Duration:  durationMS(e.StartTime(), e.EndTime()),
			DayIndex:  e.DayIndex(),
			IsOff:     e.IsOff(),
		}
	}
	return shiftPatternRecord{
		ID:           p.ID(),
		Name:         p.Name(),
		RefStartDate: uint64(p.RefStartDate().UnixMilli()), //nolint:gosec
		CycleLength:  p.CycleLength(),
		Entries:      entries,
		CreatedAt:    uint64(p.CreatedAt().UnixMilli()), //nolint:gosec
	}
}

func fromShiftPatternRecord(rec shiftPatternRecord) *entity.ShiftPattern {
	refStart := time.UnixMilli(int64(rec.RefStartDate)).UTC() //nolint:gosec
	createdAt := time.UnixMilli(int64(rec.CreatedAt)).UTC()   //nolint:gosec

	p := entity.ReconstructShiftPattern(rec.ID, rec.Name, refStart, rec.CycleLength, createdAt)
	for _, e := range rec.Entries {
		if e.IsOff {
			p.AddEntry(entity.NewDayOffEntry(e.DayIndex))
		} else {
			start := time.UnixMilli(int64(e.StartTime)).UTC() //nolint:gosec
			end := endFromDuration(start, e.Duration)
			p.AddEntry(entity.NewShiftEntry(e.DayIndex, e.Name, start, end))
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
