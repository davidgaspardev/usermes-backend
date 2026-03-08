package usecase

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type getAllLocationsUseCase struct {
	locationRepo output.LocationRepository
}

// NewGetAllLocationsUseCase creates a new GetAllLocationsUseCase instance.
func NewGetAllLocationsUseCase(locationRepo output.LocationRepository) input.GetAllLocationsUseCase {
	return &getAllLocationsUseCase{
		locationRepo: locationRepo,
	}
}

func (uc *getAllLocationsUseCase) Execute(ctx context.Context) ([]entity.Location, error) {
	return uc.locationRepo.GetAll()
}
