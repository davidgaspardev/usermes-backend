package dto

import (
	"errors"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// ShiftEntryRequest is the request body for a single shift entry.
// Days with no entries are implicitly off — there is no need to declare off days.
type ShiftEntryRequest struct {
	Name      string `json:"name"`
	StartTime string `json:"start_time"` // ISO 8601 "HH:MM:SS"
	EndTime   string `json:"end_time"`   // ISO 8601 "HH:MM:SS"
	DayIndex  int    `json:"day_index"`  // 0-based day position within the cycle
}

// CreateShiftPatternRequest is the request body for creating a shift pattern.
type CreateShiftPatternRequest struct {
	Name         string              `json:"name"`
	RefStartDate string              `json:"ref_start_date"`
	Entries      []ShiftEntryRequest `json:"entries"`
	CycleLength  int                 `json:"cycle_length"`
}

// Validate validates the create shift pattern request.
func (r *CreateShiftPatternRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.RefStartDate == "" {
		return errors.New("ref_start_date is required")
	}
	if r.CycleLength < 1 {
		return errors.New("cycle_length must be >= 1")
	}
	if len(r.Entries) == 0 {
		return errors.New("entries are required")
	}
	for _, e := range r.Entries {
		if e.DayIndex < 0 || e.DayIndex >= r.CycleLength {
			return errors.New("each entry day_index must be in [0, cycle_length)")
		}
		if e.Name == "" || e.StartTime == "" || e.EndTime == "" {
			return errors.New("each entry must have name, start_time, and end_time")
		}
	}
	return nil
}

// ParseRefStartDate parses the ref_start_date field into a time.Time value.
func (r *CreateShiftPatternRequest) ParseRefStartDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.RefStartDate)
}

// UpdateShiftPatternRequest is the request body for updating a shift pattern.
type UpdateShiftPatternRequest struct {
	Name         string              `json:"name"`
	RefStartDate string              `json:"ref_start_date"`
	Entries      []ShiftEntryRequest `json:"entries"`
	CycleLength  int                 `json:"cycle_length"`
}

// Validate validates the update shift pattern request.
func (r *UpdateShiftPatternRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.RefStartDate == "" {
		return errors.New("ref_start_date is required")
	}
	if r.CycleLength < 1 {
		return errors.New("cycle_length must be >= 1")
	}
	if len(r.Entries) == 0 {
		return errors.New("entries are required")
	}
	for _, e := range r.Entries {
		if e.DayIndex < 0 || e.DayIndex >= r.CycleLength {
			return errors.New("each entry day_index must be in [0, cycle_length)")
		}
		if e.Name == "" || e.StartTime == "" || e.EndTime == "" {
			return errors.New("each entry must have name, start_time, and end_time")
		}
	}
	return nil
}

// ParseRefStartDate parses the ref_start_date field into a time.Time value.
func (r *UpdateShiftPatternRequest) ParseRefStartDate() (time.Time, error) {
	return time.Parse("2006-01-02", r.RefStartDate)
}

// ShiftEntryResponse is the JSON response for a single shift entry.
type ShiftEntryResponse struct {
	Name      string `json:"name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	DayIndex  int    `json:"day_index"`
}

// ShiftPatternResponse is the JSON response for a shift pattern.
type ShiftPatternResponse struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	RefStartDate string               `json:"ref_start_date"`
	CreatedAt    string               `json:"created_at"`
	Entries      []ShiftEntryResponse `json:"entries"`
	CycleLength  int                  `json:"cycle_length"`
}

// ToShiftPatternResponse converts a ShiftPattern entity to a ShiftPatternResponse DTO.
func ToShiftPatternResponse(p *entity.ShiftPattern) ShiftPatternResponse {
	entries := make([]ShiftEntryResponse, len(p.Entries()))
	for i, e := range p.Entries() {
		entries[i] = ShiftEntryResponse{
			DayIndex:  e.DayIndex(),
			Name:      e.Name(),
			StartTime: e.StartTime().Format("15:04:05"),
			EndTime:   e.EndTime().Format("15:04:05"),
		}
	}

	return ShiftPatternResponse{
		ID:           p.ID().String(),
		Name:         p.Name(),
		RefStartDate: p.RefStartDate().Format("2006-01-02"),
		CycleLength:  p.CycleLength(),
		Entries:      entries,
		CreatedAt:    p.CreatedAt().Format(time.RFC3339),
	}
}
