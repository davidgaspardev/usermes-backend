package input

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/resource/domain/entity"
)

type SiteService interface {
	Create(
		ctx context.Context,
		code, name string,
		latitude, longitude float64,
		ownerID uuid.UUID,
	) (*entity.Site, error)
}
