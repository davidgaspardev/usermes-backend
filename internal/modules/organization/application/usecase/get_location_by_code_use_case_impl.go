package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type getLocationByCodeUseCase struct {
	locationRepo output.LocationRepository
}

// NewGetLocationByCodeUseCase creates a new GetLocationByCodeUseCase instance.
func NewGetLocationByCodeUseCase(locationRepo output.LocationRepository) input.GetLocationByCodeUseCase {
	return &getLocationByCodeUseCase{locationRepo: locationRepo}
}

func (uc *getLocationByCodeUseCase) Execute(ctx context.Context, code string) (*entity.Location, error) {
	location, err := uc.locationRepo.FindTree(code)
	if err != nil {
		return nil, err
	}
	if location == nil {
		return nil, errors.ErrLocationTreeNotFound
	}
	return location, nil
}
