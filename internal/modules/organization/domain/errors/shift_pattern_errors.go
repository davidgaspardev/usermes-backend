package errors

import "errors"

// Shift pattern domain errors.
var (
	// ErrShiftPatternNotFound is returned when a shift pattern is not found.
	ErrShiftPatternNotFound = errors.New("shift pattern not found")

	// ErrShiftPatternAlreadyExists is returned when a shift pattern with the same name already exists.
	ErrShiftPatternAlreadyExists = errors.New("shift pattern already exists")

	// ErrShiftEntryIsOff is returned when trying to compute times for an off-day entry.
	ErrShiftEntryIsOff = errors.New("cannot compute times for an off-day entry")

	// ErrInvalidShiftPatternCommand is returned when the command to create a shift pattern is invalid.
	ErrInvalidShiftPatternCommand = errors.New("invalid shift pattern command: name and entries are required")
)
