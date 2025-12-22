package entity

import (
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type Location struct {
	code      string
	name      string
	parent    *Location
	children  []*Location
	createdAt time.Time
	updatedAt time.Time
}

func NewLocation(code string, name string, parent *Location) *Location {
	now := time.Now()
	return &Location{
		code:      code,
		name:      name,
		parent:    parent,
		children:  nil,
		createdAt: now,
		updatedAt: now,
	}
}

func NewRootLocation(code string, name string) *Location {
	return NewLocation(code, name, nil)
}

func (l *Location) Code() string {
	return l.code
}

func (l *Location) Name() string {
	return l.name
}

func (l *Location) ParentCode() string {
	return l.parent.code
}

func (l *Location) Children() []*Location {
	return l.children
}

func (l *Location) FindByCode(code string) (*Location, error) {
	if l.children == nil {
		return nil, errors.ErrLocationNotFound
	}

	for _, child := range l.children {
		if child.code == code {
			return child, nil
		} else {
			locationFound, err := child.FindByCode(code)
			if err != nil {
				return nil, err
			}
			return locationFound, nil
		}
	}

	return nil, nil
}

func (l *Location) AddChild(location *Location) {
	if l.children == nil {
		l.children = make([]*Location, 1)
		l.children[0] = location
	} else {
		l.children = append(l.children, location)
	}
}
