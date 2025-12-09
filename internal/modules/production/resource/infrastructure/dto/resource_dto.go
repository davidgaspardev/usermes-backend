package dto

import (
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/resource/domain/entity"
)

// CreateResourceRequest represents the request body for creating a resource
type CreateResourceRequest struct {
	ShiftID    *string  `json:"shift_id,omitempty"`
	Code       string   `json:"code" validate:"required,min=2,max=50"`
	Type       string   `json:"type" validate:"required,min=2,max=50"`
	Tags       []string `json:"tags,omitempty"`
	StopFactor int16    `json:"stop_factor" validate:"gte=0"`
}

// UpdateResourceRequest represents the request body for updating a resource
type UpdateResourceRequest struct {
	ShiftID    *string  `json:"shift_id,omitempty"`
	Code       string   `json:"code" validate:"required,min=2,max=50"`
	Type       string   `json:"type" validate:"required,min=2,max=50"`
	Tags       []string `json:"tags,omitempty"`
	StopFactor int16    `json:"stop_factor" validate:"gte=0"`
}

// ResourceResponse represents the response body for resource data
type ResourceResponse struct {
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ShiftID    *string   `json:"shift_id,omitempty"`
	ID         string    `json:"id"`
	Code       string    `json:"code"`
	Type       string    `json:"type"`
	Tags       []string  `json:"tags,omitempty"`
	StopFactor int16     `json:"stop_factor"`
}

// ResourceListResponse represents the response body for a list of resources
type ResourceListResponse struct {
	Resources []*ResourceResponse `json:"resources"`
	Total     int                 `json:"total"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
}

// ErrorResponse represents the response body for errors
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message"`
}

// ToResourceResponse converts a Resource entity to ResourceResponse DTO
func ToResourceResponse(resource *entity.Resource) ResourceResponse {
	return ResourceResponse{
		ID:         resource.ID().String(),
		Code:       resource.Code(),
		ShiftID:    resource.ShiftID(),
		Type:       resource.Type(),
		StopFactor: resource.StopFactor(),
		Tags:       resource.Tags(),
		CreatedAt:  resource.CreatedAt(),
		UpdatedAt:  resource.UpdatedAt(),
	}
}

// ToResourceListResponse converts a slice of Resources to ResourceListResponse
func ToResourceListResponse(resources []*entity.Resource, limit, offset int) ResourceListResponse {
	responses := make([]*ResourceResponse, len(resources))
	for i, resource := range resources {
		r := ToResourceResponse(resource)
		responses[i] = &r
	}

	return ResourceListResponse{
		Resources: responses,
		Total:     len(responses),
		Limit:     limit,
		Offset:    offset,
	}
}

// NewErrorResponse creates a new ErrorResponse
func NewErrorResponse(errMsg, message string) ErrorResponse {
	return ErrorResponse{
		Error:   errMsg,
		Message: message,
	}
}

// NewSuccessResponse creates a new SuccessResponse
func NewSuccessResponse(message string, data interface{}) SuccessResponse {
	return SuccessResponse{
		Message: message,
		Data:    data,
	}
}
