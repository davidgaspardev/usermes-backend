package entity

import (
	"time"

	"github.com/google/uuid"
)

// Shift is a concrete occurrence of a ShiftTemplate on a specific date/time range.
// Its ID is stored on the production Event to link the event to a shift period.
type Shift struct {
	startAt      time.Time
	endAt        time.Time
	createdAt    time.Time
	calendarCode string
	name         string
	templateID   uuid.UUID
	id           uuid.UUID
}

// NewShift creates a Shift instance from a ShiftTemplate with concrete start and end timestamps.
func NewShift(template *ShiftTemplate, calendarCode string, startAt, endAt time.Time) *Shift {
	return &Shift{
		id:           uuid.New(),
		templateID:   template.id,
		calendarCode: calendarCode,
		name:         template.name,
		startAt:      startAt,
		endAt:        endAt,
		createdAt:    time.Now(),
	}
}

// ID returns the shift instance's unique identifier.
func (s *Shift) ID() uuid.UUID {
	return s.id
}

// TemplateID returns the ID of the ShiftTemplate this instance was created from.
func (s *Shift) TemplateID() uuid.UUID {
	return s.templateID
}

// CalendarCode returns the code of the ShiftCalendar this shift belongs to.
func (s *Shift) CalendarCode() string {
	return s.calendarCode
}

// Name returns the shift's display name (inherited from its template).
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
