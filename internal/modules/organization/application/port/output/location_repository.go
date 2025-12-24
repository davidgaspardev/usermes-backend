package output

import "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"

type LocationRepository interface {
	Create(location *entity.Location) error
	FindTree(rootCode string) (*entity.Location, error)
	ExistsByCode(code string) (bool, error)
	GetAll() ([]entity.Location, error)
}
