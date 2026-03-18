package output

import "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"

// LocationRepository defines the persistence port for location data.
type LocationRepository interface {
	Create(location *entity.Location) error
	Update(location *entity.Location) error
	FindByCode(code string) (*entity.Location, error)
	FindTree(rootCode string) (*entity.Location, error)
	ExistsByCode(code string) (bool, error)
	GetAll() ([]entity.Location, error)
}
