package dto

import (
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// CreatePlantRequest represents the request to create a new plant
type CreatePlantRequest struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UpdatePlantRequest represents the request to update a plant
type UpdatePlantRequest struct {
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// PlantResponse represents a plant in the response
type PlantResponse struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"owner_id"`
	IsActive  bool      `json:"is_active"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse represents a success response with data
type SuccessResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// NewErrorResponse creates a new error response
func NewErrorResponse(errType, message string) ErrorResponse {
	return ErrorResponse{
		Error:   errType,
		Message: message,
	}
}

// NewSuccessResponse creates a new success response
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}

// ToPlantResponse converts a plant entity to a response DTO
func ToPlantResponse(plant *entity.Plant) PlantResponse {
	return PlantResponse{
		ID:        plant.ID(),
		Code:      plant.Code(),
		Name:      plant.Name(),
		Latitude:  plant.Latitude(),
		Longitude: plant.Longitude(),
		OwnerID:   plant.OwnerID(),
		IsActive:  plant.IsActive(),
		CreatedAt: plant.CreatedAt(),
		UpdatedAt: plant.UpdatedAt(),
	}
}

// ToPlantResponseList converts a list of plant entities to response DTOs
func ToPlantResponseList(plants []*entity.Plant) []PlantResponse {
	responses := make([]PlantResponse, len(plants))
	for i, plant := range plants {
		responses[i] = ToPlantResponse(plant)
	}
	return responses
}

// Validate validates the CreatePlantRequest
func (r *CreatePlantRequest) Validate() error {
	if r.Code == "" {
		return NewValidationError("code is required")
	}
	if r.Name == "" {
		return NewValidationError("name is required")
	}
	return nil
}

// Validate validates the UpdatePlantRequest
func (r *UpdatePlantRequest) Validate() error {
	if r.Name == "" {
		return NewValidationError("name is required")
	}
	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	message string
}

func (e ValidationError) Error() string {
	return e.message
}

// NewValidationError creates a new validation error
func NewValidationError(message string) error {
	return ValidationError{message: message}
}
