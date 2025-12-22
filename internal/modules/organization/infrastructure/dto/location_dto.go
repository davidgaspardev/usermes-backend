package dto

import "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"

type CreateLocationRootRequest struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type LocationRootResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func ToLocationRootResponse(locationRoot *entity.Location) LocationRootResponse {
	return LocationRootResponse{
		Code: locationRoot.Code(),
		Name: locationRoot.Name(),
	}
}
