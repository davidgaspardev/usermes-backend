package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/domain/entity"
)

// ResourceRepository defines the output port for resource persistence operations
type ResourceRepository interface {
	// Save persists a new resource
	Save(ctx context.Context, resource *entity.Resource) error

	// Update updates an existing resource
	Update(ctx context.Context, resource *entity.Resource) error

	// Delete removes a resource by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// FindByID retrieves a resource by ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error)

	// FindByCode retrieves a resource by code
	FindByCode(ctx context.Context, code string) (*entity.Resource, error)

	// ExistsByCode checks if a resource with the given code exists
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// FindAll retrieves all resources with pagination
	// limit: maximum number of results to return (0 = no limit)
	// offset: number of results to skip
	FindAll(ctx context.Context, limit, offset int) ([]*entity.Resource, error)

	// FindByType retrieves all resources of a specific type
	FindByType(ctx context.Context, resourceType string, limit, offset int) ([]*entity.Resource, error)

	// FindByShiftID retrieves all resources assigned to a specific shift
	FindByShiftID(ctx context.Context, shiftID string, limit, offset int) ([]*entity.Resource, error)
}
