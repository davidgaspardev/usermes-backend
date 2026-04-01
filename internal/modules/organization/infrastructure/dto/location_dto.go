package dto

import (
	"errors"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

// CreateLocationRootRequest is the request body for creating a root location.
type CreateLocationRootRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

// Validate validates the create root location request.
func (r *CreateLocationRootRequest) Validate() error {
	if r.Code == "" || r.Name == "" {
		return errors.New("code and name are required")
	}
	return nil
}

// AddLocationRequest is the request body for adding a child location.
type AddLocationRequest struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Kind       string `json:"kind" validate:"required"`
	ParentCode string `json:"parent_code" validate:"required"`
	RootCode   string `json:"root_code" validate:"required"`
}

// Validate validates the add location request.
func (r *AddLocationRequest) Validate() error {
	if r.Code == "" || r.Name == "" || r.Kind == "" || r.ParentCode == "" || r.RootCode == "" {
		return errors.New("code, name, kind, parent_code, and root_code are required")
	}
	return nil
}

// LocationResponse is the JSON response for a single location node.
type LocationResponse struct {
	Code           string             `json:"code"`
	Name           string             `json:"name"`
	Kind           string             `json:"kind"`
	CreatedBy      string             `json:"created_by,omitempty"`
	ShiftPatternID string             `json:"shift_pattern_id,omitempty"`
	Children       []LocationResponse `json:"children,omitempty"`
}

// ToLocationResponse converts a Location entity to a LocationResponse DTO.
func ToLocationResponse(location *entity.Location) LocationResponse {
	children := make([]LocationResponse, len(location.Children()))
	for i, child := range location.Children() {
		children[i] = ToLocationResponse(child)
	}

	resp := LocationResponse{
		Code:     location.Code(),
		Name:     location.Name(),
		Kind:     string(location.Kind()),
		Children: children,
	}

	if id := location.CreatedBy(); id != (uuid.UUID{}) {
		resp.CreatedBy = id.String()
	}

	if id := location.ShiftPatternID(); id != (uuid.UUID{}) {
		resp.ShiftPatternID = id.String()
	}

	return resp
}

// AllLocationsResponse is the JSON response for a list of location trees.
type AllLocationsResponse struct {
	Locations []LocationResponse `json:"locations"`
}

// ToAllLocationsResponse converts a slice of Location entities to an AllLocationsResponse DTO.
func ToAllLocationsResponse(locations []entity.Location) AllLocationsResponse {
	locationsResponse := make([]LocationResponse, len(locations))
	for i, location := range locations {
		locationsResponse[i] = ToLocationResponse(&location)
	}
	return AllLocationsResponse{
		Locations: locationsResponse,
	}
}
