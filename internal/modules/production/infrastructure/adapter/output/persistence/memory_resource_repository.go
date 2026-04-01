package persistence

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
)

// MemoryResourceRepository is an in-memory implementation of ResourceRepository
type MemoryResourceRepository struct {
	resources       map[uuid.UUID]*entity.Resource
	resourcesByCode map[string]*entity.Resource
	mu              sync.RWMutex
}

// NewMemoryResourceRepository creates a new instance of MemoryResourceRepository
func NewMemoryResourceRepository() output.ResourceRepository {
	return &MemoryResourceRepository{
		resources:       make(map[uuid.UUID]*entity.Resource),
		resourcesByCode: make(map[string]*entity.Resource),
	}
}

// Save persists a new resource
func (r *MemoryResourceRepository) Save(ctx context.Context, resource *entity.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if resource with same ID already exists
	if _, exists := r.resources[resource.ID()]; exists {
		return errors.ErrResourceAlreadyExists
	}

	// Check if resource with same code already exists
	if _, exists := r.resourcesByCode[resource.Code()]; exists {
		return errors.ErrCodeAlreadyExists
	}

	// Store resource
	r.resources[resource.ID()] = resource
	r.resourcesByCode[resource.Code()] = resource

	return nil
}

// Update updates an existing resource
func (r *MemoryResourceRepository) Update(ctx context.Context, resource *entity.Resource) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if resource exists
	existing, exists := r.resources[resource.ID()]
	if !exists {
		return errors.ErrResourceNotFound
	}

	// If code changed, update code index
	if existing.Code() != resource.Code() {
		// Remove old code entry
		delete(r.resourcesByCode, existing.Code())

		// Check if new code already exists
		if _, exists := r.resourcesByCode[resource.Code()]; exists {
			// Restore old code entry before returning error
			r.resourcesByCode[existing.Code()] = existing
			return errors.ErrCodeAlreadyExists
		}

		// Add new code entry
		r.resourcesByCode[resource.Code()] = resource
	}

	// Update resource
	r.resources[resource.ID()] = resource

	return nil
}

// Delete removes a resource by ID
func (r *MemoryResourceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if resource exists
	resource, exists := r.resources[id]
	if !exists {
		return errors.ErrResourceNotFound
	}

	// Remove from both maps
	delete(r.resources, id)
	delete(r.resourcesByCode, resource.Code())

	return nil
}

// FindByID retrieves a resource by ID
func (r *MemoryResourceRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resources[id]
	if !exists {
		return nil, errors.ErrResourceNotFound
	}

	return resource, nil
}

// FindByCode retrieves a resource by code
func (r *MemoryResourceRepository) FindByCode(ctx context.Context, code string) (*entity.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	resource, exists := r.resourcesByCode[code]
	if !exists {
		return nil, errors.ErrResourceNotFound
	}

	return resource, nil
}

// ExistsByCode checks if a resource with the given code exists
func (r *MemoryResourceRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.resourcesByCode[code]
	return exists, nil
}

// FindAll retrieves all resources with pagination
func (r *MemoryResourceRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	resources := make([]*entity.Resource, 0, len(r.resources))
	for _, resource := range r.resources {
		resources = append(resources, resource)
	}

	// Apply pagination
	if offset >= len(resources) {
		return []*entity.Resource{}, nil
	}

	end := offset + limit
	if limit == 0 || end > len(resources) {
		end = len(resources)
	}

	return resources[offset:end], nil
}

// FindByType retrieves all resources of a specific type
func (r *MemoryResourceRepository) FindByType(ctx context.Context, resourceType string, limit, offset int) ([]*entity.Resource, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Filter resources by type
	filtered := make([]*entity.Resource, 0)
	for _, resource := range r.resources {
		if resource.Type() == resourceType {
			filtered = append(filtered, resource)
		}
	}

	// Apply pagination
	if offset >= len(filtered) {
		return []*entity.Resource{}, nil
	}

	end := offset + limit
	if limit == 0 || end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], nil
}
