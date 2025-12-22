package errors

import "errors"

var (
	ErrLocationAlreadyExists  = errors.New("location already exists")
	ErrLocationTreeNotFound   = errors.New("location tree not found")
	ErrLocationNotFound       = errors.New("location not found")
	ErrInvalidLocationCommand = errors.New("invalid location command")
)
