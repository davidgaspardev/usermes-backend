package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/domain/errors"
)

// PlantServiceImpl implements the PlantService interface
type PlantServiceImpl struct {
	plantRepository output.PlantRepository
}

// NewPlantService creates a new instance of PlantServiceImpl
func NewPlantService(plantRepository output.PlantRepository) input.PlantService {
	return &PlantServiceImpl{
		plantRepository: plantRepository,
	}
}

// Create creates a new plant
func (s *PlantServiceImpl) Create(
	ctx context.Context,
	code string,
	name string,
	latitude float64,
	longitude float64,
	ownerID uuid.UUID,
) (*entity.Plant, error) {
	// Validate code
	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	// Validate name
	if err := s.validateName(name); err != nil {
		return nil, err
	}

	// Validate coordinates
	if err := s.validateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}

	// Validate owner ID
	if ownerID == uuid.Nil {
		return nil, errors.ErrInvalidOwnerID
	}

	// Normalize code to uppercase
	code = strings.ToUpper(strings.TrimSpace(code))

	// Check if code already exists
	exists, err := s.plantRepository.ExistsByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrPlantCodeAlreadyExists
	}

	// Create plant entity
	plant := entity.NewPlant(code, name, latitude, longitude, ownerID)

	// Save plant
	if err := s.plantRepository.Save(ctx, plant); err != nil {
		return nil, err
	}

	return plant, nil
}

// GetByID retrieves a plant by its ID
func (s *PlantServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Plant, error) {
	plant, err := s.plantRepository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrPlantNotFound
	}

	return plant, nil
}

// GetByCode retrieves a plant by its code
func (s *PlantServiceImpl) GetByCode(ctx context.Context, code string) (*entity.Plant, error) {
	// Normalize code to uppercase
	code = strings.ToUpper(strings.TrimSpace(code))

	plant, err := s.plantRepository.FindByCode(ctx, code)
	if err != nil {
		return nil, errors.ErrPlantNotFound
	}

	return plant, nil
}

// GetByOwner retrieves all plants owned by a user
func (s *PlantServiceImpl) GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*entity.Plant, error) {
	if ownerID == uuid.Nil {
		return nil, errors.ErrInvalidOwnerID
	}

	plants, err := s.plantRepository.FindByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return plants, nil
}

// ListAll retrieves all active plants
func (s *PlantServiceImpl) ListAll(ctx context.Context) ([]*entity.Plant, error) {
	plants, err := s.plantRepository.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	// Filter only active plants
	activePlants := make([]*entity.Plant, 0)
	for _, plant := range plants {
		if plant.IsActive() {
			activePlants = append(activePlants, plant)
		}
	}

	return activePlants, nil
}

// Update updates plant information
func (s *PlantServiceImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	name string,
	latitude float64,
	longitude float64,
) (*entity.Plant, error) {
	// Validate name
	if err := s.validateName(name); err != nil {
		return nil, err
	}

	// Validate coordinates
	if err := s.validateCoordinates(latitude, longitude); err != nil {
		return nil, err
	}

	// Find plant
	plant, err := s.plantRepository.FindByID(ctx, id)
	if err != nil {
		return nil, errors.ErrPlantNotFound
	}

	// Update plant
	plant.Update(name, latitude, longitude)

	// Save changes
	if err := s.plantRepository.Update(ctx, plant); err != nil {
		return nil, err
	}

	return plant, nil
}

// Deactivate deactivates a plant
func (s *PlantServiceImpl) Deactivate(ctx context.Context, id uuid.UUID) error {
	// Find plant
	plant, err := s.plantRepository.FindByID(ctx, id)
	if err != nil {
		return errors.ErrPlantNotFound
	}

	// Deactivate plant
	plant.Deactivate()

	// Save changes
	if err := s.plantRepository.Update(ctx, plant); err != nil {
		return err
	}

	return nil
}

// Activate activates a plant
func (s *PlantServiceImpl) Activate(ctx context.Context, id uuid.UUID) error {
	// Find plant
	plant, err := s.plantRepository.FindByID(ctx, id)
	if err != nil {
		return errors.ErrPlantNotFound
	}

	// Activate plant
	plant.Activate()

	// Save changes
	if err := s.plantRepository.Update(ctx, plant); err != nil {
		return err
	}

	return nil
}

// Delete removes a plant
func (s *PlantServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if plant exists
	_, err := s.plantRepository.FindByID(ctx, id)
	if err != nil {
		return errors.ErrPlantNotFound
	}

	// Delete plant
	if err := s.plantRepository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// validateCode validates the plant code
func (s *PlantServiceImpl) validateCode(code string) error {
	code = strings.TrimSpace(code)

	if code == "" {
		return errors.ErrPlantCodeRequired
	}

	if len(code) < 2 {
		return errors.ErrPlantCodeTooShort
	}

	if len(code) > 20 {
		return errors.ErrPlantCodeTooLong
	}

	return nil
}

// validateName validates the plant name
func (s *PlantServiceImpl) validateName(name string) error {
	name = strings.TrimSpace(name)

	if name == "" {
		return errors.ErrPlantNameRequired
	}

	if len(name) < 2 {
		return errors.ErrPlantNameTooShort
	}

	if len(name) > 100 {
		return errors.ErrPlantNameTooLong
	}

	return nil
}

// validateCoordinates validates latitude and longitude
func (s *PlantServiceImpl) validateCoordinates(latitude, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.ErrInvalidLatitude
	}

	if longitude < -180 || longitude > 180 {
		return errors.ErrInvalidLongitude
	}

	return nil
}
