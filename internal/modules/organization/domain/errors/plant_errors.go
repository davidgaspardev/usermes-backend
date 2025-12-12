package errors

import "errors"

var (
	// ErrPlantNotFound is returned when a plant is not found
	ErrPlantNotFound = errors.New("plant not found")

	// ErrPlantCodeRequired is returned when plant code is empty
	ErrPlantCodeRequired = errors.New("plant code is required")

	// ErrPlantCodeTooShort is returned when plant code is too short
	ErrPlantCodeTooShort = errors.New("plant code must be at least 2 characters")

	// ErrPlantCodeTooLong is returned when plant code is too long
	ErrPlantCodeTooLong = errors.New("plant code must not exceed 20 characters")

	// ErrPlantCodeAlreadyExists is returned when plant code already exists
	ErrPlantCodeAlreadyExists = errors.New("plant code already exists")

	// ErrPlantNameRequired is returned when plant name is empty
	ErrPlantNameRequired = errors.New("plant name is required")

	// ErrPlantNameTooShort is returned when plant name is too short
	ErrPlantNameTooShort = errors.New("plant name must be at least 2 characters")

	// ErrPlantNameTooLong is returned when plant name is too long
	ErrPlantNameTooLong = errors.New("plant name must not exceed 100 characters")

	// ErrInvalidLatitude is returned when latitude is out of range
	ErrInvalidLatitude = errors.New("latitude must be between -90 and 90")

	// ErrInvalidLongitude is returned when longitude is out of range
	ErrInvalidLongitude = errors.New("longitude must be between -180 and 180")

	// ErrInvalidOwnerID is returned when owner ID is invalid
	ErrInvalidOwnerID = errors.New("owner ID is required")

	// ErrPlantInactive is returned when trying to operate on inactive plant
	ErrPlantInactive = errors.New("plant is inactive")

	// ErrUnauthorizedPlantAccess is returned when user doesn't have access to plant
	ErrUnauthorizedPlantAccess = errors.New("unauthorized access to plant")
)
