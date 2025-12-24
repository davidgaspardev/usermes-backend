package persistence

import (
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type LocationRecord struct {
	Code       string
	Name       string
	Kind       string
	ParentCode string
}

type locationRepositoryInMemory struct {
	locations []LocationRecord
}

func NewLocationRepositoryInMemory() output.LocationRepository {
	return &locationRepositoryInMemory{
		locations: make([]LocationRecord, 0),
	}
}

func (r *locationRepositoryInMemory) Create(location *entity.Location) error {
	r.locations = append(r.locations, LocationRecord{
		Code:       location.Code(),
		Name:       location.Name(),
		Kind:       string(location.Kind()),
		ParentCode: location.ParentCode(),
	})
	return nil
}

func (r *locationRepositoryInMemory) FindTree(rootCode string) (*entity.Location, error) {
	for _, location := range r.locations {
		if location.Code == rootCode && location.ParentCode == "" {
			return r.buildLocationTree(rootCode), nil
		}
	}

	return nil, nil
}

func (r *locationRepositoryInMemory) buildLocationTree(rootCode string) *entity.Location {
	var root *entity.Location

	// Encontra o location raiz para este parentCode
	for _, location := range r.locations {
		if location.Code == rootCode {
			root = entity.NewLocation(
				location.Code,
				location.Name,
				entity.LocationKind(location.Kind),
				nil,
			)
			break
		}
	}

	if root == nil {
		return nil
	}

	// Constrói a árvore de filhos recursivamente
	r.buildChildren(root)

	return root
}

func (r *locationRepositoryInMemory) buildChildren(parent *entity.Location) {
	for _, location := range r.locations {
		if location.ParentCode == parent.Code() {
			childCopy := entity.NewLocation(
				location.Code,
				location.Name,
				entity.LocationKind(location.Kind),
				parent,
			)
			r.buildChildren(childCopy)
			parent.AddChild(childCopy)
		}
	}
}

func (r *locationRepositoryInMemory) ExistsByCode(code string) (bool, error) {
	for _, location := range r.locations {
		if location.Code == code {
			return true, nil
		}
	}

	return false, nil
}

func (r *locationRepositoryInMemory) GetAll() ([]entity.Location, error) {
	var locations = []entity.Location{}

	for _, location := range r.locations {
		if location.ParentCode == "" {
			locationTree, err := r.FindTree(location.Code)
			if err != nil {
				return nil, err
			}

			locations = append(locations, *locationTree)
		}
	}

	return locations, nil
}
