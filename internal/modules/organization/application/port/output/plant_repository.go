package output

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// PlantRepository defines the output port for plant persistence operations
type PlantRepository interface {
	// Save persists a new plant
	Save(ctx context.Context, plant *entity.Plant) error

	// Update updates an existing plant
	Update(ctx context.Context, plant *entity.Plant) error

	// FindByID retrieves a plant by its ID
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Plant, error)

	// FindByCode retrieves a plant by its code
	FindByCode(ctx context.Context, code string) (*entity.Plant, error)

	// FindByOwner retrieves all plants owned by a user
	FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]*entity.Plant, error)

	// FindAll retrieves all plants
	FindAll(ctx context.Context) ([]*entity.Plant, error)

	// ExistsByCode checks if a plant with the given code exists
	ExistsByCode(ctx context.Context, code string) (bool, error)

	// Delete removes a plant by its ID
	Delete(ctx context.Context, id uuid.UUID) error
}
