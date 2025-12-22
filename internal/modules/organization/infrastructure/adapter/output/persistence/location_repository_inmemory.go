package persistence

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type locationRepositoryInMemory struct {
	locations []*entity.Location
}

func NewLocationRepositoryInMemory() output.LocationRepository {
	return &locationRepositoryInMemory{
		locations: make([]*entity.Location, 0),
	}
}

func (r *locationRepositoryInMemory) Create(location *entity.Location) error {
	r.locations = append(r.locations, location)
	return nil
}

func (r *locationRepositoryInMemory) FindTree(rootCode string) (*entity.Location, error) {
	for _, location := range r.locations {
		if location.Code() == rootCode && location.ParentCode() == "" {
			return location, nil
		}
	}

	return nil, nil
}

func (r *locationRepositoryInMemory) FindByCode(code string) (*entity.Location, error) {
	for _, location := range r.locations {
		if location.Code() == code {
			return location, nil
		}
	}

	return nil, nil
}
