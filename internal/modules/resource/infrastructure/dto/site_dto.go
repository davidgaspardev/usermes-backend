package dto

import "github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"

type CreateSiteRequest struct {
	Code      string  `json:"code" validate:"required"`
	Name      string  `json:"name" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

type SiteResponse struct {
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	OwnerID   string  `json:"owner_id"`
}

func ToSiteResponse(site *entity.Site) *SiteResponse {
	return &SiteResponse{
		Code:      site.Code,
		Name:      site.Name,
		Latitude:  site.Latitude,
		Longitude: site.Longitude,
		OwnerID:   site.OwnerID.String(),
	}
}
