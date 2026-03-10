package errors

import "errors"

var (
	// ErrResourceNotFound is returned when a resource is not found
	ErrResourceNotFound = errors.New("resource not found")

	// ErrResourceAlreadyExists is returned when trying to create a resource that already exists
	ErrResourceAlreadyExists = errors.New("resource already exists")

	// ErrCodeAlreadyExists is returned when a resource with the same code already exists
	ErrCodeAlreadyExists = errors.New("resource with this code already exists")

	// ErrResourceCodeAlreadyExists is returned when a resource with the same code already exists in the same plant
	ErrResourceCodeAlreadyExists = errors.New("resource with this code already exists in this plant")

	// ErrInvalidCode is returned when the resource code is invalid
	ErrInvalidCode = errors.New("invalid resource code")

	// ErrInvalidType is returned when the resource type is invalid
	ErrInvalidType = errors.New("invalid resource type")

	// ErrInvalidStopFactor is returned when the stop factor is invalid
	ErrInvalidStopFactor = errors.New("invalid stop factor: must be non-negative")

	// ErrShiftNotFound is returned when a shift instance is not found
	ErrShiftNotFound = errors.New("shift not found")
)
