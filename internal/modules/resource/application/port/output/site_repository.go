package output

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
)

type SiteRepository interface {
	Create(ctx context.Context, site *entity.Site) error
}
