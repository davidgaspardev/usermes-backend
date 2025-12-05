package entity

import (
	"time"

	"github.com/google/uuid"
)

// Resource represents the core resource entity in the domain
type Resource struct {
	createdAt    time.Time
	updatedAt    time.Time
	shiftID      *string
	code         string
	resourceType string
	stopFactor   int16
	id           uuid.UUID
}

// NewResource creates a new resource instance
func NewResource(code string, shiftID *string, resourceType string, stopFactor int16) *Resource {
	now := time.Now()
	return &Resource{
		id:           uuid.New(),
		code:         code,
		shiftID:      shiftID,
		resourceType: resourceType,
		stopFactor:   stopFactor,
		createdAt:    now,
		updatedAt:    now,
	}
}

// ReconstructResource reconstructs a resource from persistence
func ReconstructResource(
	id uuid.UUID,
	code string,
	shiftID *string,
	resourceType string,
	stopFactor int16,
	createdAt time.Time,
	updatedAt time.Time,
) *Resource {
	return &Resource{
		id:           id,
		code:         code,
		shiftID:      shiftID,
		resourceType: resourceType,
		stopFactor:   stopFactor,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ID returns the resource ID
func (r *Resource) ID() uuid.UUID {
	return r.id
}

// Code returns the resource code
func (r *Resource) Code() string {
	return r.code
}

// ShiftID returns the shift ID (can be nil)
func (r *Resource) ShiftID() *string {
	return r.shiftID
}

// Type returns the resource type
func (r *Resource) Type() string {
	return r.resourceType
}

// StopFactor returns the stop factor
func (r *Resource) StopFactor() int16 {
	return r.stopFactor
}

// CreatedAt returns the creation timestamp
func (r *Resource) CreatedAt() time.Time {
	return r.createdAt
}

// UpdatedAt returns the last update timestamp
func (r *Resource) UpdatedAt() time.Time {
	return r.updatedAt
}

// UpdateCode updates the resource code
func (r *Resource) UpdateCode(code string) {
	r.code = code
	r.updatedAt = time.Now()
}

// UpdateShiftID updates the shift ID
func (r *Resource) UpdateShiftID(shiftID *string) {
	r.shiftID = shiftID
	r.updatedAt = time.Now()
}

// UpdateType updates the resource type
func (r *Resource) UpdateType(resourceType string) {
	r.resourceType = resourceType
	r.updatedAt = time.Now()
}

// UpdateStopFactor updates the stop factor
func (r *Resource) UpdateStopFactor(stopFactor int16) {
	r.stopFactor = stopFactor
	r.updatedAt = time.Now()
}

// Update updates all mutable fields at once
func (r *Resource) Update(code string, shiftID *string, resourceType string, stopFactor int16) {
	r.code = code
	r.shiftID = shiftID
	r.resourceType = resourceType
	r.stopFactor = stopFactor
	r.updatedAt = time.Now()
}
