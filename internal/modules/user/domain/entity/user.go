package entity

import (
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/user/domain/valueobject"
	"github.com/google/uuid"
)

// User represents the core user entity in the domain
type User struct {
	id          uuid.UUID
	email       valueobject.Email
	password    valueobject.Password
	lastLoginAt *time.Time
	name        string
	createdAt   time.Time
	updatedAt   time.Time
	isActive    bool
}

// NewUser creates a new user instance
func NewUser(email valueobject.Email, password valueobject.Password, name string) *User {
	now := time.Now()
	return &User{
		id:        uuid.New(),
		email:     email,
		password:  password,
		name:      name,
		isActive:  true,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructUser reconstructs a user from persistence
func ReconstructUser(
	id uuid.UUID,
	email valueobject.Email,
	password valueobject.Password,
	name string,
	isActive bool,
	createdAt time.Time,
	updatedAt time.Time,
	lastLoginAt *time.Time,
) *User {
	return &User{
		id:          id,
		email:       email,
		password:    password,
		name:        name,
		isActive:    isActive,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
		lastLoginAt: lastLoginAt,
	}
}

// ID returns the user's unique identifier
func (u *User) ID() uuid.UUID {
	return u.id
}

// Email returns the user's email
func (u *User) Email() valueobject.Email {
	return u.email
}

// Password returns the user's password
func (u *User) Password() valueobject.Password {
	return u.password
}

// Name returns the user's name
func (u *User) Name() string {
	return u.name
}

// IsActive returns whether the user is active
func (u *User) IsActive() bool {
	return u.isActive
}

// CreatedAt returns when the user was created
func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

// UpdatedAt returns when the user was last updated
func (u *User) UpdatedAt() time.Time {
	return u.updatedAt
}

// LastLoginAt returns when the user last logged in
func (u *User) LastLoginAt() *time.Time {
	return u.lastLoginAt
}

// UpdateName updates the user's name
func (u *User) UpdateName(name string) {
	u.name = name
	u.updatedAt = time.Now()
}

// UpdateEmail updates the user's email
func (u *User) UpdateEmail(email valueobject.Email) {
	u.email = email
	u.updatedAt = time.Now()
}

// UpdatePassword updates the user's password
func (u *User) UpdatePassword(password valueobject.Password) {
	u.password = password
	u.updatedAt = time.Now()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.isActive = false
	u.updatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.isActive = true
	u.updatedAt = time.Now()
}

// RecordLogin records a successful login
func (u *User) RecordLogin() {
	now := time.Now()
	u.lastLoginAt = &now
	u.updatedAt = now
}

// VerifyPassword checks if the provided password matches the user's password
func (u *User) VerifyPassword(plainPassword string) bool {
	return u.password.Compare(plainPassword)
}
