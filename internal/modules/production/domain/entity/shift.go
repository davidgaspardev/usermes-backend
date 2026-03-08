package entity

import (
	"time"

	"github.com/google/uuid"
)

// Shift is a concrete occurrence of a ShiftPattern on a specific date/time range.
// Its ID is stored on the production Event to link the event to a shift period.
type Shift struct {
	startAt   time.Time
	endAt     time.Time
	createdAt time.Time
	name      string
	patternID uuid.UUID
	id        uuid.UUID
}

// NewShift creates a Shift instance from a ShiftPattern ID with concrete start and end timestamps.
func NewShift(patternID uuid.UUID, name string, startAt, endAt time.Time) *Shift {
	return &Shift{
		id:        uuid.New(),
		patternID: patternID,
		name:      name,
		startAt:   startAt,
		endAt:     endAt,
		createdAt: time.Now(),
	}
}

// ID returns the shift instance's unique identifier.
func (s *Shift) ID() uuid.UUID {
	return s.id
}

// PatternID returns the ID of the ShiftPattern this instance was created from.
func (s *Shift) PatternID() uuid.UUID {
	return s.patternID
}

// Name returns the shift's display name (inherited from its pattern entry).
func (s *Shift) Name() string {
	return s.name
}

// StartAt returns the concrete start timestamp of this shift occurrence.
func (s *Shift) StartAt() time.Time {
	return s.startAt
}

// EndAt returns the concrete end timestamp of this shift occurrence.
func (s *Shift) EndAt() time.Time {
	return s.endAt
}

// CreatedAt returns when this shift instance was created.
func (s *Shift) CreatedAt() time.Time {
	return s.createdAt
}
