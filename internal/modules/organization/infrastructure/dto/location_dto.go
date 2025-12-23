package dto

import "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"

type CreateLocationRootRequest struct {
	Code string `json:"code" validate:"required"`
	Name string `json:"name" validate:"required"`
}

type AddLocationRequest struct {
	Code       string `json:"code" validate:"required"`
	Name       string `json:"name" validate:"required"`
	ParentCode string `json:"parent_code" validate:"required"`
	RootCode   string `json:"root_code" validate:"required"`
}

type LocationResponse struct {
	Code     string             `json:"code"`
	Name     string             `json:"name"`
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
		Children: children,
	}
}
