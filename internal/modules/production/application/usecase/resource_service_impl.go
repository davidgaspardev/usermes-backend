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
	locationCode string,
	code string,
	resourceType string,
	stopFactor int16,
	tags []string,
) (*entity.Resource, error) {
	locationCode = strings.ToUpper(strings.TrimSpace(locationCode))
	if locationCode == "" {
		return nil, errors.ErrInvalidCode
	}

	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	if err := s.validateResourceType(resourceType); err != nil {
		return nil, err
	}

	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	existing, err := s.repository.FindByCode(ctx, code)
	if err == nil && existing != nil && existing.LocationCode() == locationCode {
		return nil, errors.ErrResourceCodeAlreadyExists
	}

	resource := entity.NewResource(locationCode, code, resourceType, stopFactor, tags)

	if err := s.repository.Save(ctx, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

// Update updates an existing resource
func (s *ResourceServiceImpl) Update(
	ctx context.Context,
	id uuid.UUID,
	locationCode string,
	code string,
	resourceType string,
	stopFactor int16,
	tags []string,
) (*entity.Resource, error) {
	locationCode = strings.ToUpper(strings.TrimSpace(locationCode))
	if locationCode == "" {
		return nil, errors.ErrInvalidCode
	}

	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	if err := s.validateResourceType(resourceType); err != nil {
		return nil, err
	}

	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	code = strings.ToUpper(strings.TrimSpace(code))

	resource, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if resource.LocationCode() != locationCode {
		resource.UpdateLocationCode(locationCode)
	}
	resource.Update(code, resourceType, stopFactor, tags)

	if err := s.repository.Update(ctx, resource); err != nil {
		return nil, err
	}

	return resource, nil
}

// Delete deletes a resource by ID
func (s *ResourceServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repository.Delete(ctx, id)
}

// GetByID retrieves a resource by ID
func (s *ResourceServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error) {
	return s.repository.FindByID(ctx, id)
}

// GetByCode retrieves a resource by code within a location
func (s *ResourceServiceImpl) GetByCode(ctx context.Context, locationCode, code string) (*entity.Resource, error) {
	locationCode = strings.ToUpper(strings.TrimSpace(locationCode))
	code = strings.ToUpper(strings.TrimSpace(code))

	resource, err := s.repository.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if resource.LocationCode() != locationCode {
		return nil, errors.ErrResourceNotFound
	}

	return resource, nil
}

// GetByLocation retrieves all resources for a specific location
func (s *ResourceServiceImpl) GetByLocation(ctx context.Context, locationCode string, limit, offset int) ([]*entity.Resource, error) {
	locationCode = strings.ToUpper(strings.TrimSpace(locationCode))

	allResources, err := s.repository.FindAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Resource, 0)
	for _, resource := range allResources {
		if resource.LocationCode() == locationCode {
			result = append(result, resource)
		}
	}

	return result, nil
}

// GetAll retrieves all resources with pagination
func (s *ResourceServiceImpl) GetAll(ctx context.Context, limit, offset int) ([]*entity.Resource, error) {
	return s.repository.FindAll(ctx, limit, offset)
}

// GetByType retrieves all resources of a specific type
func (s *ResourceServiceImpl) GetByType(ctx context.Context, resourceType string, limit, offset int) ([]*entity.Resource, error) {
	return s.repository.FindByType(ctx, resourceType, limit, offset)
}

func (s *ResourceServiceImpl) validateCode(code string) error {
	code = strings.TrimSpace(code)
	if code == "" || len(code) < 2 || len(code) > 50 {
		return errors.ErrInvalidCode
	}
	return nil
}

func (s *ResourceServiceImpl) validateResourceType(resourceType string) error {
	resourceType = strings.TrimSpace(resourceType)
	if resourceType == "" || len(resourceType) < 2 || len(resourceType) > 50 {
		return errors.ErrInvalidType
	}
	return nil
}

func (s *ResourceServiceImpl) validateStopFactor(stopFactor int16) error {
	if stopFactor < 0 {
		return errors.ErrInvalidStopFactor
	}
	return nil
}
