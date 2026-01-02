package dto

import (
	"errors"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type CreateLocationRootRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

func (r *CreateLocationRootRequest) Validate() error {
	if r.Code == "" || r.Name == "" {
		return errors.New("code and name are required")
	}
	return nil
}

type AddLocationRequest struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	Kind       string `json:"kind" validate:"required"`
	ParentCode string `json:"parent_code" validate:"required"`
	RootCode   string `json:"root_code" validate:"required"`
}

type LocationResponse struct {
	Code     string             `json:"code"`
	Name     string             `json:"name"`
	Kind     string             `json:"kind"`
	Children []LocationResponse `json:"children"`
}

func ToLocationResponse(location *entity.Location) LocationResponse {
	children := make([]LocationResponse, len(location.Children()))
	for i, child := range location.Children() {
		children[i] = ToLocationResponse(child)
	}

	return LocationResponse{
		Code:     location.Code(),
		Name:     location.Name(),
		Kind:     string(location.Kind()),
		Children: children,
	}
}

type AllLocationsResponse struct {
	Locations []LocationResponse `json:"locations"`
}

func ToAllLocationsResponse(locations []entity.Location) AllLocationsResponse {
	locationsResponse := make([]LocationResponse, len(locations))
	for i, location := range locations {
		locationsResponse[i] = ToLocationResponse(&location)
	}
	return AllLocationsResponse{
		Locations: locationsResponse,
	}
}
