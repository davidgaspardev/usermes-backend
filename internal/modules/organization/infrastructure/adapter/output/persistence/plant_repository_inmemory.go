package persistence

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// InMemoryPlantRepository implements PlantRepository using in-memory storage
type InMemoryPlantRepository struct {
	plants map[uuid.UUID]*entity.Plant
	mu     sync.RWMutex
}

// NewInMemoryPlantRepository creates a new instance of InMemoryPlantRepository
func NewInMemoryPlantRepository() output.PlantRepository {
	return &InMemoryPlantRepository{
		plants: make(map[uuid.UUID]*entity.Plant),
	}
}

// Save persists a new plant
func (r *InMemoryPlantRepository) Save(ctx context.Context, plant *entity.Plant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.plants[plant.ID()] = plant
	return nil
}

// Update updates an existing plant
func (r *InMemoryPlantRepository) Update(ctx context.Context, plant *entity.Plant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plants[plant.ID()]; !exists {
		return errors.New("plant not found")
	}

	r.plants[plant.ID()] = plant
	return nil
}

// FindByID retrieves a plant by its ID
func (r *InMemoryPlantRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Plant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plant, exists := r.plants[id]
	if !exists {
		return nil, errors.New("plant not found")
	}

	return plant, nil
}

// FindByCode retrieves a plant by its code
func (r *InMemoryPlantRepository) FindByCode(ctx context.Context, code string) (*entity.Plant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plant := range r.plants {
		if plant.Code() == code {
			return plant, nil
		}
	}

	return nil, errors.New("plant not found")
}

// FindByOwner retrieves all plants owned by a user
func (r *InMemoryPlantRepository) FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]*entity.Plant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plants := make([]*entity.Plant, 0)
	for _, plant := range r.plants {
		if plant.OwnerID() == ownerID {
			plants = append(plants, plant)
		}
	}

	return plants, nil
}

// FindAll retrieves all plants
func (r *InMemoryPlantRepository) FindAll(ctx context.Context) ([]*entity.Plant, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plants := make([]*entity.Plant, 0, len(r.plants))
	for _, plant := range r.plants {
		plants = append(plants, plant)
	}

	return plants, nil
}

// ExistsByCode checks if a plant with the given code exists
func (r *InMemoryPlantRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, plant := range r.plants {
		if plant.Code() == code {
			return true, nil
		}
	}

	return false, nil
}

// Delete removes a plant by its ID
func (r *InMemoryPlantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.plants[id]; !exists {
		return errors.New("plant not found")
	}

	delete(r.plants, id)
	return nil
}
