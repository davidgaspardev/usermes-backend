package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/domain/entity"
)

// PlantService defines the input port for plant operations (use cases)
type PlantService interface {
	// Create creates a new plant
	Create(
		ctx context.Context,
		code string,
		name string,
		latitude float64,
		longitude float64,
		ownerID uuid.UUID,
	) (*entity.Plant, error)

	// GetByID retrieves a plant by its ID
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Plant, error)

	// GetByCode retrieves a plant by its code
	GetByCode(ctx context.Context, code string) (*entity.Plant, error)

	// GetByOwner retrieves all plants owned by a user
	GetByOwner(ctx context.Context, ownerID uuid.UUID) ([]*entity.Plant, error)

	// ListAll retrieves all active plants
	ListAll(ctx context.Context) ([]*entity.Plant, error)

	// Update updates plant information
	Update(
		ctx context.Context,
		id uuid.UUID,
		name string,
		latitude float64,
		longitude float64,
	) (*entity.Plant, error)

	// Deactivate deactivates a plant
	Deactivate(ctx context.Context, id uuid.UUID) error

	// Activate activates a plant
	Activate(ctx context.Context, id uuid.UUID) error

	// Delete removes a plant
	Delete(ctx context.Context, id uuid.UUID) error
}
