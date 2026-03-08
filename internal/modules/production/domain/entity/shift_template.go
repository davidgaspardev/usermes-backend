package entity

import (
	"time"

	"github.com/google/uuid"
)

// ShiftTemplate defines a named recurring time window pattern (e.g. Morning 06:00–14:00).
// It is used as the blueprint to create concrete Shift instances for production events.
type ShiftTemplate struct {
	createdAt time.Time
	name      string
	startTime string // "HH:MM" format
	endTime   string // "HH:MM" format
	id        uuid.UUID
}

// NewShiftTemplate creates a new ShiftTemplate with the given name and time window.
func NewShiftTemplate(name, startTime, endTime string) *ShiftTemplate {
	return &ShiftTemplate{
		id:        uuid.New(),
		name:      name,
		startTime: startTime,
		endTime:   endTime,
		createdAt: time.Now(),
	}
}

// ID returns the shift template's unique identifier.
func (s *ShiftTemplate) ID() uuid.UUID {
	return s.id
}

// Name returns the shift template's display name.
func (s *ShiftTemplate) Name() string {
	return s.name
}

// StartTime returns the shift template's start time in "HH:MM" format.
func (s *ShiftTemplate) StartTime() string {
	return s.startTime
}

// EndTime returns the shift template's end time in "HH:MM" format.
func (s *ShiftTemplate) EndTime() string {
	return s.endTime
}

// CreatedAt returns the shift template creation timestamp.
func (s *ShiftTemplate) CreatedAt() time.Time {
	return s.createdAt
}
