package output

import "github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"

type EventRepository interface {
	Create(event *entity.Event) error
	GetCurrentByResCode(resCode string) (*entity.Event, error)
	UpdateCurrent(event *entity.Event) error
}
