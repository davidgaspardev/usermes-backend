package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	organizationerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

// ShiftEntry is a single slot in a ShiftPattern sequence.
// It represents one day's shift assignment within the rotating cycle.
type ShiftEntry struct {
	name      string
	startTime string // "HH:MM", empty when IsOff
	endTime   string // "HH:MM", empty when IsOff
	dayIndex  int    // 0-based position in the sequence
	isOff     bool   // true for rest/off days
}

// NewShiftEntry creates a working ShiftEntry at the given position in the sequence.
func NewShiftEntry(dayIndex int, name, startTime, endTime string) *ShiftEntry {
	return &ShiftEntry{
		dayIndex:  dayIndex,
		name:      name,
		startTime: startTime,
		endTime:   endTime,
		isOff:     false,
	}
}

// NewDayOffEntry creates a rest-day ShiftEntry at the given position in the sequence.
func NewDayOffEntry(dayIndex int) *ShiftEntry {
	return &ShiftEntry{
		dayIndex: dayIndex,
		name:     "Off",
		isOff:    true,
	}
}

// DayIndex returns the 0-based position of this entry in its pattern sequence.
func (e *ShiftEntry) DayIndex() int {
	return e.dayIndex
}

// Name returns the shift entry's display name (e.g. "Morning", "Night", "Off").
func (e *ShiftEntry) Name() string {
	return e.name
}

// StartTime returns the start time in "HH:MM" format (empty for off days).
func (e *ShiftEntry) StartTime() string {
	return e.startTime
}

// EndTime returns the end time in "HH:MM" format (empty for off days).
func (e *ShiftEntry) EndTime() string {
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

	startH, startM, err := parseHHMM(e.startTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid startTime %q: %w", e.startTime, err)
	}

	endH, endM, err := parseHHMM(e.endTime)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid endTime %q: %w", e.endTime, err)
	}

	loc := date.Location()
	startAt = time.Date(date.Year(), date.Month(), date.Day(), startH, startM, 0, 0, loc)
	endAt = time.Date(date.Year(), date.Month(), date.Day(), endH, endM, 0, 0, loc)

	// overnight shift: end must be strictly after start
	if !endAt.After(startAt) {
		endAt = endAt.AddDate(0, 0, 1)
	}

	return startAt, endAt, nil
}

// parseHHMM parses a "HH:MM" string into hour and minute integers.
func parseHHMM(s string) (hour, min int, err error) {
	_, err = fmt.Sscanf(s, "%d:%d", &hour, &min)
	return
}

// ShiftPattern is a rotating sequence of ShiftEntry items anchored to a reference date.
// The period length is derived from the number of entries (7 = weekly, 15 = quinzenal, 30 = monthly, etc.).
// Use EntryForDate to resolve which ShiftEntry applies to any given calendar date.
type ShiftPattern struct {
	refStartDate time.Time
	createdAt    time.Time
	name         string
	entries      []*ShiftEntry
	id           uuid.UUID
}

// NewShiftPattern creates a new ShiftPattern anchored to the given reference start date.
func NewShiftPattern(name string, refStartDate time.Time) *ShiftPattern {
	return &ShiftPattern{
		id:           uuid.New(),
		name:         name,
		refStartDate: refStartDate,
		entries:      nil,
		createdAt:    time.Now(),
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

// PeriodDays returns the length of the rotation cycle (number of entries).
func (p *ShiftPattern) PeriodDays() int {
	return len(p.entries)
}

// Entries returns the ordered sequence of shift entries in this pattern.
func (p *ShiftPattern) Entries() []*ShiftEntry {
	return p.entries
}

// AddEntry appends a ShiftEntry to the sequence.
func (p *ShiftPattern) AddEntry(entry *ShiftEntry) {
	p.entries = append(p.entries, entry)
}

// EntryForDate returns the ShiftEntry that applies to the given calendar date,
// computed as (days since refStartDate) modulo PeriodDays.
// Returns nil if the pattern has no entries.
func (p *ShiftPattern) EntryForDate(date time.Time) *ShiftEntry {
	if len(p.entries) == 0 {
		return nil
	}

	ref := time.Date(p.refStartDate.Year(), p.refStartDate.Month(), p.refStartDate.Day(), 0, 0, 0, 0, time.UTC)
	d := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)

	daysSince := int(d.Sub(ref).Hours() / 24)

	if daysSince < 0 {
		// date is before the reference: wrap the cycle backwards
		mod := (-daysSince) % len(p.entries)
		if mod == 0 {
			daysSince = 0
		} else {
			daysSince = len(p.entries) - mod
		}
	}

	return p.entries[daysSince%len(p.entries)]
}

// CreatedAt returns the shift pattern creation timestamp.
func (p *ShiftPattern) CreatedAt() time.Time {
	return p.createdAt
}
