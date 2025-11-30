package usecase

import (
	"context"
	"strings"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/errors"
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
	code string,
	shiftID *string,
	resourceType string,
	stopFactor int16,
) (*entity.Resource, error) {
	// Validate code
	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	// Validate type
	if err := s.validateType(resourceType); err != nil {
		return nil, err
	}

	// Validate stop factor
	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	// Check if code already exists
	exists, err := s.repository.ExistsByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.ErrCodeAlreadyExists
	}

	// Create resource
	resource := entity.NewResource(code, shiftID, resourceType, stopFactor)

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
	code string,
	shiftID *string,
	resourceType string,
	stopFactor int16,
) (*entity.Resource, error) {
	// Validate code
	if err := s.validateCode(code); err != nil {
		return nil, err
	}

	// Validate type
	if err := s.validateType(resourceType); err != nil {
		return nil, err
	}

	// Validate stop factor
	if err := s.validateStopFactor(stopFactor); err != nil {
		return nil, err
	}

	// Find existing resource
	resource, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if code is being changed and if new code already exists
	if resource.Code() != code {
		exists, err := s.repository.ExistsByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.ErrCodeAlreadyExists
		}
	}

	// Update resource
	resource.Update(code, shiftID, resourceType, stopFactor)

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

// GetByCode retrieves a resource by code
func (s *ResourceServiceImpl) GetByCode(ctx context.Context, code string) (*entity.Resource, error) {
	return s.repository.FindByCode(ctx, code)
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

// validateType validates the resource type
func (s *ResourceServiceImpl) validateType(resourceType string) error {
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
