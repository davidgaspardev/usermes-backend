package entity

import "time"

// EventType represents the kind of production event.
type EventType string

// Valid EventType values for production events.
const (
	ResourceCreated         EventType = "RESOURCE_CREATED"
	OperatorSignIn          EventType = "OPERATOR_SIGN_IN"
	OperatorSignOut         EventType = "OPERATOR_SIGN_OUT"
	OperatorSignInWithItem  EventType = "OPERATOR_SIGN_IN_WITH_ITEM"
	OperatorSignOutWithItem EventType = "OPERATOR_SIGN_OUT_WITH_ITEM"
	SwapOperator            EventType = "SWAP_OPERATOR"
	InsertStop              EventType = "INSERT_STOP"
	OperatorSetProduction   EventType = "OPERATOR_SET_PRODUTION"
	OperatorChangeItem      EventType = "OPERATOR_CHANGE_ITEM"
	ChangeShift             EventType = "CHANGE_SHIFT"
	OperatorRemoveItem      EventType = "OPERATOR_REMOVE_ITEM"
	OperatorAddItem         EventType = "OPERATOR_ADD_ITEM"
)

// EventStatus represents the production status recorded by an event.
type EventStatus string

// Valid EventStatus values for production events.
const (
	Production EventStatus = "PRODUCTION"
	Stop       EventStatus = "STOP"
)

// Event represents a production event tied to a resource and shift.
type Event struct {
	dateStart     time.Time
	updatedAt     time.Time
	createdAt     time.Time
	dateEnd       *time.Time
	shiftID       string
	resCode       string
	prodCode      string
	eventType     EventType
	userID        string
	statusCode    string
	status        EventStatus
	prodQuantity  uint32
	scrapQuantity uint32
}

// NewEvent creates a new production Event with the given parameters.
func NewEvent(
	eventType EventType,
	status EventStatus,
	statusCode string,
	prodCode string,
	resCode string,
	prodQuantity uint32,
	scrapQuantity uint32,
	dateStart time.Time,
	shiftID string,
	userID string,
) *Event {
	time := time.Now()
	return &Event{
		eventType:     eventType,
		status:        status,
		statusCode:    statusCode,
		prodCode:      prodCode,
		resCode:       resCode,
		prodQuantity:  prodQuantity,
		scrapQuantity: scrapQuantity,
		dateStart:     dateStart,
		dateEnd:       nil,
		shiftID:       shiftID,
		userID:        userID,
		createdAt:     time,
		updatedAt:     time,
	}
}

// NewEventChangeShift closes the current event at the shift boundary and opens a new one under the
// next shift. It carries over status, statusCode, prodCode, resCode, userID, and quantities from
// the previous event. The caller must persist both the updated prev and the returned event.
func NewEventChangeShift(prev *Event, newShiftID string, boundaryTime time.Time) *Event {
	prev.dateEnd = &boundaryTime
	prev.updatedAt = boundaryTime

	return NewEvent(
		ChangeShift,
		prev.status,
		prev.statusCode,
		prev.prodCode,
		prev.resCode,
		0,
		0,
		boundaryTime,
		newShiftID,
		prev.userID,
	)
}

// NewEventResourceCreated creates a resource-created event for the given resource and shift.
func NewEventResourceCreated(
	resCode string,
	shiftID string,
	userID string,
) *Event {
	return NewEvent(
		ResourceCreated,
		Stop,
		"",
		"",
		resCode,
		0,
		0,
		time.Now(),
		shiftID,
		userID,
	)
}
