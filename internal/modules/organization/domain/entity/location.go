package entity

import (
	"time"
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
	if l.parent == nil {
		return ""
	}
	return l.parent.code
}

func (l *Location) Children() []*Location {
	return l.children
}

func (l *Location) FindByCode(code string) *Location {
	if l.code == code {
		return l
	}

	if l.children == nil {
		return nil
	}

	for _, child := range l.children {
		if child.code == code {
			return child
		}
		locationFound := child.FindByCode(code)
		if locationFound != nil {
			return locationFound
		}
	}

	return nil
}

func (l *Location) ExistsByCode(code string) bool {
	location := l.FindByCode(code)
	return location != nil
}

func (l *Location) AddChild(location *Location) {
	l.children = append(l.children, location)
}
