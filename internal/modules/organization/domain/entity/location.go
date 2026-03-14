package entity

import (
	"time"

	"github.com/google/uuid"
)

// LocationKind represents the type of a location in the organization hierarchy.
type LocationKind string

// Valid LocationKind values for the organization hierarchy.
const (
	LocationKindPlant   LocationKind = "PLANT"
	LocationKindArea    LocationKind = "AREA"
	LocationKindLine    LocationKind = "LINE"
	LocationKindSection LocationKind = "SECTION"
)

// Location represents a node in the organization location tree.
type Location struct {
	createdAt      time.Time
	updatedAt      time.Time
	parent         *Location
	code           string
	name           string
	kind           LocationKind
	children       []*Location
	shiftPatternID uuid.UUID
}

// NewLocation creates a new Location with the given attributes.
func NewLocation(code string, name string, kind LocationKind, parent *Location, shiftPatternID uuid.UUID) *Location {
	now := time.Now()
	return &Location{
		code:           code,
		name:           name,
		kind:           kind,
		parent:         parent,
		children:       nil,
		shiftPatternID: shiftPatternID,
		createdAt:      now,
		updatedAt:      now,
	}
}

// NewRootLocation creates a new root Location (plant level) with no parent.
// Plants do not operate shifts, so no shift pattern is assigned.
func NewRootLocation(code string, name string) *Location {
	return NewLocation(code, name, LocationKindPlant, nil, uuid.UUID{})
}

// Code returns the location's unique code.
func (l *Location) Code() string {
	return l.code
}

// Name returns the location's display name.
func (l *Location) Name() string {
	return l.name
}

// Kind returns the location's kind (plant, area, line, section).
func (l *Location) Kind() LocationKind {
	return l.kind
}

// ParentCode returns the code of the parent location, or empty string if root.
func (l *Location) ParentCode() string {
	if l.parent == nil {
		return ""
	}
	return l.parent.code
}

// Children returns the direct child locations.
func (l *Location) Children() []*Location {
	return l.children
}

// FindByCode searches for a location by code within this subtree.
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

// ExistsByCode reports whether a location with the given code exists in this subtree.
func (l *Location) ExistsByCode(code string) bool {
	location := l.FindByCode(code)
	return location != nil
}

// IsValidLocationKind reports whether the given string is a valid LocationKind value.
func IsValidLocationKind(kind string) bool {
	switch LocationKind(kind) {
	case LocationKindPlant, LocationKindArea, LocationKindLine, LocationKindSection:
		return true
	default:
		return false
	}
}

// ShiftPatternID returns the shift pattern assigned to this location.
func (l *Location) ShiftPatternID() uuid.UUID {
	return l.shiftPatternID
}

// AddChild appends a child location to this location's children list.
func (l *Location) AddChild(location *Location) {
	l.children = append(l.children, location)
}

// CreatedAt returns the location creation timestamp.
func (l *Location) CreatedAt() time.Time {
	return l.createdAt
}

// UpdatedAt returns the last update timestamp.
func (l *Location) UpdatedAt() time.Time {
	return l.updatedAt
}

// AssignShiftPattern sets the shift pattern for this location.
func (l *Location) AssignShiftPattern(shiftPatternID uuid.UUID) {
	l.shiftPatternID = shiftPatternID
	l.updatedAt = time.Now()
}
