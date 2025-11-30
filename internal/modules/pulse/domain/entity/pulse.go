package entity

import "time"

type Pulse struct {
	date      time.Time
	createdAt time.Time
	updatedAt time.Time
	id        string
	ownerID   string
	factor    int16
}

func NewPulse(id, ownerID string, factor int16, date time.Time) *Pulse {
	now := time.Now()

	return &Pulse{
		id:        id,
		ownerID:   ownerID,
		factor:    factor,
		date:      date,
		createdAt: now,
		updatedAt: now,
	}
}
