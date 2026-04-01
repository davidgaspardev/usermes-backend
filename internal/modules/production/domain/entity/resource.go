package entity

import (
	"time"

	"github.com/google/uuid"
)

// Resource represents the core resource entity in the domain
type Resource struct {
	createdAt    time.Time
	updatedAt    time.Time
	locationCode string
	code         string
	resourceType string
	tags         []string
	stopFactor   int16
	id           uuid.UUID
}

// copyStringSlice creates a defensive copy of a string slice
func copyStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}

// NewResource creates a new resource instance
func NewResource(locationCode, code, resourceType string, stopFactor int16, tags []string) *Resource {
	now := time.Now()
	return &Resource{
		id:           uuid.New(),
		locationCode: locationCode,
		code:         code,
		resourceType: resourceType,
		stopFactor:   stopFactor,
		tags:         copyStringSlice(tags),
		createdAt:    now,
		updatedAt:    now,
	}
}

// ReconstructResource reconstructs a resource from persistence
func ReconstructResource(
	id uuid.UUID,
	locationCode, code string,
	resourceType string,
	stopFactor int16,
	tags []string,
	createdAt time.Time,
	updatedAt time.Time,
) *Resource {
	return &Resource{
		id:           id,
		locationCode: locationCode,
		code:         code,
		resourceType: resourceType,
		stopFactor:   stopFactor,
		tags:         copyStringSlice(tags),
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ID returns the resource ID
func (r *Resource) ID() uuid.UUID {
	return r.id
}

// LocationCode returns the location code
func (r *Resource) LocationCode() string {
	return r.locationCode
}

// Code returns the resource code
func (r *Resource) Code() string {
	return r.code
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

// Tags returns the resource tags
func (r *Resource) Tags() []string {
	return copyStringSlice(r.tags)
}

// UpdateCode updates the resource code
func (r *Resource) UpdateCode(code string) {
	r.code = code
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

// UpdateTags updates the resource tags
func (r *Resource) UpdateTags(tags []string) {
	r.tags = copyStringSlice(tags)
	r.updatedAt = time.Now()
}

// UpdateLocationCode updates the location code
func (r *Resource) UpdateLocationCode(locationCode string) {
	r.locationCode = locationCode
	r.updatedAt = time.Now()
}

// Update updates all mutable fields at once
func (r *Resource) Update(code, resourceType string, stopFactor int16, tags []string) {
	r.code = code
	r.resourceType = resourceType
	r.stopFactor = stopFactor
	r.tags = copyStringSlice(tags)
	r.updatedAt = time.Now()
}
