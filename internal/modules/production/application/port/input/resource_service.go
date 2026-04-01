package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
)

// ResourceService defines the input port for resource use cases
type ResourceService interface {
	// Create creates a new resource
	Create(
		ctx context.Context,
		locationCode string,
		code string,
		resourceType string,
		stopFactor int16,
		tags []string,
	) (*entity.Resource, error)

	// Update updates an existing resource
	Update(
		ctx context.Context,
		id uuid.UUID,
		locationCode string,
		code string,
		resourceType string,
		stopFactor int16,
		tags []string,
	) (*entity.Resource, error)

	// Delete deletes a resource by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// GetByID retrieves a resource by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error)

	// GetByCode retrieves a resource by code within a location
	GetByCode(ctx context.Context, locationCode, code string) (*entity.Resource, error)

	// GetByLocation retrieves all resources for a specific location
	GetByLocation(ctx context.Context, locationCode string, limit, offset int) ([]*entity.Resource, error)

	// GetAll retrieves all resources with pagination
	GetAll(ctx context.Context, limit, offset int) ([]*entity.Resource, error)

	// GetByType retrieves all resources of a specific type
	GetByType(ctx context.Context, resourceType string, limit, offset int) ([]*entity.Resource, error)
}
