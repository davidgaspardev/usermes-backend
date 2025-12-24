package entity

import (
	"time"
)

type LocationKind string

const (
	LocationKindPlant   LocationKind = "PLANT"
	LocationKindArea    LocationKind = "AREA"
	LocationKindLine    LocationKind = "LINE"
	LocationKindSection LocationKind = "SECTION"
)

type Location struct {
	code      string
	name      string
	kind      LocationKind
	parent    *Location
	children  []*Location
	createdAt time.Time
	updatedAt time.Time
}

func NewLocation(code string, name string, kind LocationKind, parent *Location) *Location {
	now := time.Now()
	return &Location{
		code:      code,
		name:      name,
		kind:      kind,
		parent:    parent,
		children:  nil,
		createdAt: now,
		updatedAt: now,
	}
}

func NewRootLocation(code string, name string) *Location {
	return NewLocation(code, name, LocationKindPlant, nil)
}

func (l *Location) Code() string {
	return l.code
}

func (l *Location) Name() string {
	return l.name
}

func (l *Location) Kind() LocationKind {
	return l.kind
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

func IsValidLocationKind(kind string) bool {
	switch LocationKind(kind) {
	case LocationKindPlant, LocationKindArea, LocationKindLine, LocationKindSection:
		return true
	default:
		return false
	}
}

func (l *Location) AddChild(location *Location) {
	l.children = append(l.children, location)
}
