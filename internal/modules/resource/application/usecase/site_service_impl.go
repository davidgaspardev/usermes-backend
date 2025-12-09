package usecase

import (
	"context"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
	"github.com/google/uuid"
)

type SiteService struct {
	siteRepo output.SiteRepository
}

func NewSiteService(
	siteRepo output.SiteRepository,
) input.SiteService {
	return &SiteService{}
}

func (s *SiteService) Create(
	ctx context.Context,
	code, name string,
	latitude, longitude float64,
	ownerID uuid.UUID,
) (*entity.Site, error) {
	var createAt, updateAt = time.Now(), time.Now()
	site := entity.NewSite(code, name, latitude, longitude, ownerID, createAt, updateAt)
	if err := s.siteRepo.Create(ctx, site); err != nil {
		return nil, err
	}
	return site, nil
}
