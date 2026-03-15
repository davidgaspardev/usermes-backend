package entity

import (
	"time"

	"github.com/google/uuid"

	organizationerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

// ShiftEntry is a single shift slot within a ShiftPattern cycle.
// Multiple entries can share the same DayIndex, representing concurrent shifts on that day.
type ShiftEntry struct {
	startTime time.Time
	endTime   time.Time
	name      string
	dayIndex  int
	isOff     bool
}

// NewShiftEntry creates a working ShiftEntry at the given day position in the cycle.
func NewShiftEntry(dayIndex int, name string, startTime, endTime time.Time) *ShiftEntry {
	return &ShiftEntry{
		dayIndex:  dayIndex,
		name:      name,
		startTime: startTime,
		endTime:   endTime,
		isOff:     false,
	}
}

// NewDayOffEntry creates a rest-day ShiftEntry at the given day position in the cycle.
func NewDayOffEntry(dayIndex int) *ShiftEntry {
	return &ShiftEntry{
		dayIndex: dayIndex,
		name:     "Off",
		isOff:    true,
	}
}

// DayIndex returns the 0-based day position of this entry within its pattern cycle.
func (e *ShiftEntry) DayIndex() int {
	return e.dayIndex
}

// Name returns the shift entry's display name (e.g. "Shift 1", "Night", "Off").
func (e *ShiftEntry) Name() string {
	return e.name
}

// StartTime returns the start time-of-day (zero value for off days).
func (e *ShiftEntry) StartTime() time.Time {
	return e.startTime
}

// EndTime returns the end time-of-day (zero value for off days).
func (e *ShiftEntry) EndTime() time.Time {
	return e.endTime
}

// IsOff reports whether this entry is a rest/off day with no production shift.
func (e *ShiftEntry) IsOff() bool {
	return e.isOff
}

// TimesForDate computes the concrete startAt and endAt timestamps for this entry
// on the given calendar date. Overnight shifts (endTime < startTime) automatically
// advance endAt to the following day.
// Returns ErrShiftEntryIsOff if called on an off-day entry.
func (e *ShiftEntry) TimesForDate(date time.Time) (startAt, endAt time.Time, err error) {
	if e.isOff {
		return time.Time{}, time.Time{}, organizationerrors.ErrShiftEntryIsOff
	}

	loc := date.Location()
	startAt = time.Date(date.Year(), date.Month(), date.Day(), e.startTime.Hour(), e.startTime.Minute(), 0, 0, loc)
	endAt = time.Date(date.Year(), date.Month(), date.Day(), e.endTime.Hour(), e.endTime.Minute(), 0, 0, loc)

	// overnight shift: end must be strictly after start
	if !endAt.After(startAt) {
		endAt = endAt.AddDate(0, 0, 1)
	}

	return startAt, endAt, nil
}

// ShiftPattern is a rotating schedule anchored to a reference date.
// CycleLength defines how many days are in one full rotation.
// Each day in the cycle (0-based index) can hold multiple ShiftEntry items,
// representing concurrent shifts (e.g. Shift 1, Shift 2, Shift 3 on the same day).
// Use EntriesForDate to resolve which shifts apply to any given calendar date.
type ShiftPattern struct {
	refStartDate time.Time
	createdAt    time.Time
	name         string
	entries      []*ShiftEntry
	id           uuid.UUID
	cycleLength  int // explicit cycle length in days
}

// NewShiftPattern creates a new ShiftPattern with the given name, reference start date,
// and explicit cycle length.
func NewShiftPattern(name string, refStartDate time.Time, cycleLength int) *ShiftPattern {
	return &ShiftPattern{
		id:           uuid.New(),
		name:         name,
		refStartDate: refStartDate,
		cycleLength:  cycleLength,
		entries:      nil,
		createdAt:    time.Now(),
	}
}

// ReconstructShiftPattern rebuilds a ShiftPattern from stored values without
// generating a new ID or resetting timestamps. Used by the persistence layer.
func ReconstructShiftPattern(id uuid.UUID, name string, refStartDate time.Time, cycleLength int, createdAt time.Time) *ShiftPattern {
	return &ShiftPattern{
		id:           id,
		name:         name,
		refStartDate: refStartDate,
		cycleLength:  cycleLength,
		entries:      nil,
		createdAt:    createdAt,
	}
}

// ID returns the shift pattern's unique identifier.
func (p *ShiftPattern) ID() uuid.UUID {
	return p.id
}

// Name returns the shift pattern's display name.
func (p *ShiftPattern) Name() string {
	return p.name
}

// RefStartDate returns the anchor date from which the cycle is calculated.
func (p *ShiftPattern) RefStartDate() time.Time {
	return p.refStartDate
}

// CycleLength returns the number of days in one full rotation cycle.
func (p *ShiftPattern) CycleLength() int {
	return p.cycleLength
}

// Entries returns all shift entries across the entire cycle.
func (p *ShiftPattern) Entries() []*ShiftEntry {
	return p.entries
}

// AddEntry appends a ShiftEntry to the pattern.
func (p *ShiftPattern) AddEntry(entry *ShiftEntry) {
	p.entries = append(p.entries, entry)
}

// EntriesForDate returns all ShiftEntry items that apply to the given calendar date.
// The day offset within the cycle is computed as (days since refStartDate) modulo CycleLength.
// Returns nil if the pattern has no entries or cycleLength is zero.
func (p *ShiftPattern) EntriesForDate(date time.Time) []*ShiftEntry {
	if len(p.entries) == 0 || p.cycleLength == 0 {
		return nil
	}

	ref := time.Date(p.refStartDate.Year(), p.refStartDate.Month(), p.refStartDate.Day(), 0, 0, 0, 0, time.UTC)
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	daysSince := int(d.Sub(ref).Hours() / 24)
	dayOffset := daysSince % p.cycleLength
	if dayOffset < 0 {
		dayOffset += p.cycleLength
	}

	var result []*ShiftEntry
	for _, e := range p.entries {
		if e.dayIndex == dayOffset {
			result = append(result, e)
		}
	}
	return result
}

// CreatedAt returns the shift pattern creation timestamp.
func (p *ShiftPattern) CreatedAt() time.Time {
	return p.createdAt
}

// SetName updates the shift pattern name.
func (p *ShiftPattern) SetName(name string) {
	p.name = name
}

// SetCycleLength updates the cycle length.
func (p *ShiftPattern) SetCycleLength(cycleLength int) {
	p.cycleLength = cycleLength
}

// SetEntries replaces the full set of shift entries.
func (p *ShiftPattern) SetEntries(entries []*ShiftEntry) {
	p.entries = entries
}
