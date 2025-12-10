package entity

import (
	"time"

	"github.com/google/uuid"
)

// Plant represents a manufacturing plant/facility in the domain
type Plant struct {
	id        uuid.UUID
	ownerID   uuid.UUID
	createdAt time.Time
	updatedAt time.Time
	latitude  float64
	longitude float64
	code      string
	name      string
	isActive  bool
}

// NewPlant creates a new plant instance
func NewPlant(code, name string, latitude, longitude float64, ownerID uuid.UUID) *Plant {
	now := time.Now()
	return &Plant{
		id:        uuid.New(),
		code:      code,
		name:      name,
		latitude:  latitude,
		longitude: longitude,
		ownerID:   ownerID,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructPlant reconstructs a plant from persistence
func ReconstructPlant(
	id uuid.UUID,
	code, name string,
	latitude, longitude float64,
	ownerID uuid.UUID,
	isActive bool,
	createdAt, updatedAt time.Time,
) *Plant {
	return &Plant{
		id:        id,
		code:      code,
		name:      name,
		latitude:  latitude,
		longitude: longitude,
		ownerID:   ownerID,
		isActive:  isActive,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID returns the plant ID
func (p *Plant) ID() uuid.UUID {
	return p.id
}

// Code returns the plant code
func (p *Plant) Code() string {
	return p.code
}

// Name returns the plant name
func (p *Plant) Name() string {
	return p.name
}

// Latitude returns the plant latitude
func (p *Plant) Latitude() float64 {
	return p.latitude
}

// Longitude returns the plant longitude
func (p *Plant) Longitude() float64 {
	return p.longitude
}

// OwnerID returns the owner user ID
func (p *Plant) OwnerID() uuid.UUID {
	return p.ownerID
}

// IsActive returns whether the plant is active
func (p *Plant) IsActive() bool {
	return p.isActive
}

// CreatedAt returns the creation timestamp
func (p *Plant) CreatedAt() time.Time {
	return p.createdAt
}

// UpdatedAt returns the last update timestamp
func (p *Plant) UpdatedAt() time.Time {
	return p.updatedAt
}

// UpdateName updates the plant name
func (p *Plant) UpdateName(name string) {
	p.name = name
	p.updatedAt = time.Now()
}

// UpdateLocation updates the plant coordinates
func (p *Plant) UpdateLocation(latitude, longitude float64) {
	p.latitude = latitude
	p.longitude = longitude
	p.updatedAt = time.Now()
}

// Update updates all mutable fields at once
func (p *Plant) Update(name string, latitude, longitude float64) {
	p.name = name
	p.latitude = latitude
	p.longitude = longitude
	p.updatedAt = time.Now()
}

// Deactivate deactivates the plant
func (p *Plant) Deactivate() {
	p.isActive = false
	p.updatedAt = time.Now()
}

// Activate activates the plant
func (p *Plant) Activate() {
	p.isActive = true
	p.updatedAt = time.Now()
}
