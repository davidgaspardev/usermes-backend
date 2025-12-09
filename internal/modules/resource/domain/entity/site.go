package entity

import (
	"time"

	"github.com/google/uuid"
)

type Site struct {
	createAt  time.Time
	updateAt  time.Time
	Code      string
	Name      string
	Latitude  float64
	Longitude float64
	OwnerID   uuid.UUID
}

func NewSite(code, name string, latitude, longitude float64, ownerID uuid.UUID, createAt, updateAt time.Time) *Site {
	return &Site{
		createAt:  createAt,
		updateAt:  updateAt,
		Code:      code,
		Name:      name,
		Latitude:  latitude,
		Longitude: longitude,
		OwnerID:   ownerID,
	}
}
