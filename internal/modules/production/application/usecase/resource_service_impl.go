package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
)

// ResourceServiceImpl implements the ResourceService interface
type ResourceServiceImpl struct {
	repository output.ResourceRepository
}

// NewResourceService creates a new instance of ResourceServiceImpl
func NewResourceService(repository output.ResourceRepository) input.ResourceService {
	return &ResourceServiceImpl{
		repository: repository,
	}
}

// Create creates a new resource
func (s *ResourceServiceImpl) Create(
	ctx context.Context,
	plantCode string,
	code string,
	shiftID *string,
	resourceType string,
	stopFactor int16,
	tags []string,
) (*entity.Resource, error) {
	// Validate plant code
	plantCode = strings.ToUpper(strings.TrimSpace(plantCode))
	if plantCode == "" {
		return nil, errors.ErrInvalidCode
	}

	// Validate code
	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	// Validate resource type
	if err := s.validateResourceType(resourceType); err != nil {
		return nil, err
	}

	// Validate stop factor
	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	// Normalize code
	code = strings.ToUpper(strings.TrimSpace(code))

	// Check if code already exists in this plant
	existing, err := s.repository.FindByCode(ctx, code)
	if err == nil && existing != nil && existing.PlantCode() == plantCode {
		return nil, errors.ErrResourceCodeAlreadyExists
	}

	// Create resource
	resource := entity.NewResource(plantCode, code, shiftID, resourceType, stopFactor, tags)

	// Save resource
	if err := s.repository.Save(ctx, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

// Update updates an existing resource
func (s *ResourceServiceImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	plantCode string,
	code string,
	shiftID *string,
	resourceType string,
	stopFactor int16,
	tags []string,
) (*entity.Resource, error) {
	// Validate plant code
	plantCode = strings.ToUpper(strings.TrimSpace(plantCode))
	if plantCode == "" {
		return nil, errors.ErrInvalidCode
	}

	// Validate code
	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	// Validate resource type
	if err := s.validateResourceType(resourceType); err != nil {
		return nil, err
	}

	// Validate stop factor
	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	// Normalize code
	code = strings.ToUpper(strings.TrimSpace(code))

	// Find existing resource
	resource, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update resource (including plant code if changed)
	if resource.PlantCode() != plantCode {
		resource.UpdatePlantCode(plantCode)
	}
	resource.Update(code, shiftID, resourceType, stopFactor, tags)

	// Save changes
	if err := s.repository.Update(ctx, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

// Delete deletes a resource by ID
func (s *ResourceServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if resource exists
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Delete resource
	return s.repository.Delete(ctx, id)
}

// GetByID retrieves a resource by ID
func (s *ResourceServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error) {
	return s.repository.FindByID(ctx, id)
}

// GetByCode retrieves a resource by code within a plant
func (s *ResourceServiceImpl) GetByCode(ctx context.Context, plantCode, code string) (*entity.Resource, error) {
	plantCode = strings.ToUpper(strings.TrimSpace(plantCode))
	code = strings.ToUpper(strings.TrimSpace(code))

	resource, err := s.repository.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	// Verify resource belongs to the specified plant
	if resource.PlantCode() != plantCode {
		return nil, errors.ErrResourceNotFound
	}

	return resource, nil
}

// GetByPlant retrieves all resources for a specific plant
func (s *ResourceServiceImpl) GetByPlant(ctx context.Context, plantCode string, limit, offset int) ([]*entity.Resource, error) {
	plantCode = strings.ToUpper(strings.TrimSpace(plantCode))

	allResources, err := s.repository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	// Filter by plant code
	plantResources := make([]*entity.Resource, 0)
	for _, resource := range allResources {
		if resource.PlantCode() == plantCode {
			plantResources = append(plantResources, resource)
		}
	}

	return plantResources, nil
}

// GetAll retrieves all resources with pagination
func (s *ResourceServiceImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.Resource, error) {
	return s.repository.FindAll(ctx, limit, offset)
}

// GetByType retrieves all resources of a specific type
func (s *ResourceServiceImpl) GetByType(ctx context.Context, resourceType string, limit, offset int) ([]*entity.Resource, error) {
	return s.repository.FindByType(ctx, resourceType, limit, offset)
}

// GetByShiftID retrieves all resources assigned to a specific shift
func (s *ResourceServiceImpl) GetByShiftID(ctx context.Context, shiftID string, limit, offset int) ([]*entity.Resource, error) {
	return s.repository.FindByShiftID(ctx, shiftID, limit, offset)
}

// validateCode validates the resource code
func (s *ResourceServiceImpl) validateCode(code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return errors.ErrInvalidCode
	}
	if len(code) < 2 || len(code) > 50 {
		return errors.ErrInvalidCode
	}
	return nil
}

// validateResourceType validates the resource type
func (s *ResourceServiceImpl) validateResourceType(resourceType string) error {
	resourceType = strings.TrimSpace(resourceType)
	if resourceType == "" {
		return errors.ErrInvalidType
	}
	if len(resourceType) < 2 || len(resourceType) > 50 {
		return errors.ErrInvalidType
	}
	return nil
}

// validateStopFactor validates the stop factor
func (s *ResourceServiceImpl) validateStopFactor(stopFactor int16) error {
	if stopFactor < 0 {
		return errors.ErrInvalidStopFactor
	}
	return nil
}
