package persistence

import (
	"context"
	"sync"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/entity"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/iam/user/domain/valueobject"
)

// MemoryUserRepository is an in-memory implementation of UserRepository
type MemoryUserRepository struct {
	users           map[uuid.UUID]*entity.User
	usersByEmail    map[string]*entity.User
	usersByUsername map[string]*entity.User
	mu              sync.RWMutex
}

// NewMemoryUserRepository creates a new instance of MemoryUserRepository
func NewMemoryUserRepository() output.UserRepository {
	return &MemoryUserRepository{
		users:           make(map[uuid.UUID]*entity.User),
		usersByEmail:    make(map[string]*entity.User),
		usersByUsername: make(map[string]*entity.User),
	}
}

// Save persists a new user
func (r *MemoryUserRepository) Save(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user already exists
	if _, exists := r.users[user.ID()]; exists {
		return errors.ErrUserAlreadyExists
	}

	// Check if email already exists
	if _, exists := r.usersByEmail[user.Email().Value()]; exists {
		return errors.ErrEmailAlreadyExists
	}

	// Check if username already exists
	if _, exists := r.usersByUsername[user.Username().Value()]; exists {
		return errors.ErrUsernameAlreadyExists
	}

	// Store user
	r.users[user.ID()] = user
	r.usersByEmail[user.Email().Value()] = user
	r.usersByUsername[user.Username().Value()] = user

	return nil
}

// Update updates an existing user
func (r *MemoryUserRepository) Update(ctx context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if user exists
	_, exists := r.users[user.ID()]
	if !exists {
		return errors.ErrUserNotFound
	}

	// Find and remove old email and username from index
	var oldEmail, oldUsername string
	for email, u := range r.usersByEmail {
		if u.ID() == user.ID() {
			oldEmail = email
			break
		}
	}
	for username, u := range r.usersByUsername {
		if u.ID() == user.ID() {
			oldUsername = username
			break
		}
	}

	// If email changed, check for conflicts
	if oldEmail != "" && oldEmail != user.Email().Value() {
		if _, exists := r.usersByEmail[user.Email().Value()]; exists {
			return errors.ErrEmailAlreadyExists
		}
		delete(r.usersByEmail, oldEmail)
	}

	// If username changed, check for conflicts
	if oldUsername != "" && oldUsername != user.Username().Value() {
		if _, exists := r.usersByUsername[user.Username().Value()]; exists {
			return errors.ErrUsernameAlreadyExists
		}
		delete(r.usersByUsername, oldUsername)
	}

	// Update all maps
	r.users[user.ID()] = user
	r.usersByEmail[user.Email().Value()] = user
	r.usersByUsername[user.Username().Value()] = user

	return nil
}

// FindByID retrieves a user by their unique identifier
func (r *MemoryUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// FindByEmail retrieves a user by their email address
func (r *MemoryUserRepository) FindByEmail(ctx context.Context, email valueobject.Email) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByEmail[email.Value()]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// FindByUsername retrieves a user by their username
func (r *MemoryUserRepository) FindByUsername(ctx context.Context, username valueobject.Username) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.usersByUsername[username.Value()]
	if !exists {
		return nil, errors.ErrUserNotFound
	}

	return user, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *MemoryUserRepository) ExistsByEmail(ctx context.Context, email valueobject.Email) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usersByEmail[email.Value()]
	return exists, nil
}

// ExistsByUsername checks if a user with the given username exists
func (r *MemoryUserRepository) ExistsByUsername(ctx context.Context, username valueobject.Username) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.usersByUsername[username.Value()]
	return exists, nil
}

// Delete removes a user from the repository
func (r *MemoryUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, exists := r.users[id]
	if !exists {
		return errors.ErrUserNotFound
	}

	// Remove from all maps
	delete(r.users, id)
	delete(r.usersByEmail, user.Email().Value())
	delete(r.usersByUsername, user.Username().Value())

	return nil
}

// FindAll retrieves all users (with optional pagination)
func (r *MemoryUserRepository) FindAll(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Convert map to slice
	users := make([]*entity.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Apply pagination
	start := offset
	if start > len(users) {
		start = len(users)
	}

	end := start + limit
	if limit <= 0 || end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}
