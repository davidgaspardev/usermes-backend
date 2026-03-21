package persistence

import (
	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

// LocationRecord is the in-memory representation of a persisted location.
// CreatedAt and UpdatedAt are stored as Unix milliseconds.
type LocationRecord struct {
	Code           string
	Name           string
	Kind           string
	ParentCode     string
	CreatedBy      string
	ShiftPatternID uuid.UUID
	CreatedAt      uint64 // Unix ms
	UpdatedAt      uint64 // Unix ms
}

type locationRepositoryInMemory struct {
	locations []LocationRecord
}

// NewLocationRepositoryInMemory creates a new in-memory LocationRepository.
func NewLocationRepositoryInMemory() output.LocationRepository {
	return &locationRepositoryInMemory{
		locations: make([]LocationRecord, 0),
	}
}

func (r *locationRepositoryInMemory) Create(location *entity.Location) error {
	r.locations = append(r.locations, LocationRecord{
		Code:           location.Code(),
		Name:           location.Name(),
		Kind:           string(location.Kind()),
		ParentCode:     location.ParentCode(),
		CreatedBy:      location.CreatedBy(),
		ShiftPatternID: location.ShiftPatternID(),
		CreatedAt:      uint64(location.CreatedAt().UnixMilli()), //nolint:gosec
		UpdatedAt:      uint64(location.UpdatedAt().UnixMilli()), //nolint:gosec
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
				location.ShiftPatternID,
				location.CreatedBy,
			)
			break
		}
	}

	if root == nil {
		return nil
	}

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
				location.ShiftPatternID,
				location.CreatedBy,
			)
			r.buildChildren(childCopy)
			parent.AddChild(childCopy)
		}
	}
}

func (r *locationRepositoryInMemory) Update(location *entity.Location) error {
	for i, loc := range r.locations {
		if loc.Code == location.Code() {
			r.locations[i].ShiftPatternID = location.ShiftPatternID()
			r.locations[i].Name = location.Name()
			r.locations[i].UpdatedAt = uint64(location.UpdatedAt().UnixMilli()) //nolint:gosec
			return nil
		}
	}
	return domainerrors.ErrLocationNotFound
}

func (r *locationRepositoryInMemory) FindByCode(code string) (*entity.Location, error) {
	for _, loc := range r.locations {
		if loc.Code == code {
			return entity.NewLocation(
				loc.Code,
				loc.Name,
				entity.LocationKind(loc.Kind),
				nil,
				loc.ShiftPatternID,
				loc.CreatedBy,
			), nil
		}
	}
	return nil, nil
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
