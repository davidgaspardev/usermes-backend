package persistence

import (
	"context"
	"sync"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
)

type MemorySiteRepository struct {
	sites map[string]*entity.Site
	mu    sync.RWMutex
}

func NewMemorySiteRepository() output.SiteRepository {
	return &MemorySiteRepository{
		sites: make(map[string]*entity.Site),
	}
}

func (r *MemorySiteRepository) Create(ctx context.Context, site *entity.Site) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sites[site.Code] = site
	return nil
}
