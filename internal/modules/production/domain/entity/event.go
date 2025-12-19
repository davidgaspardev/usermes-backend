package entity

import "time"

type EventType string

const (
	ResourceCreated         EventType = "RESOURCE_CREATED"
	OperatorSignIn          EventType = "OPERATOR_SIGN_IN"
	OperatorSignOut         EventType = "OPERATOR_SIGN_OUT"
	OperatorSignInWithItem  EventType = "OPERATOR_SIGN_IN_WITH_ITEM"
	OperatorSignOutWithItem EventType = "OPERATOR_SIGN_OUT_WITH_ITEM"
)

type EventStatus string

const (
	Production EventStatus = "PRODUCTION"
	Stop       EventStatus = "STOP"
)

type Event struct {
	eventType     EventType
	status        EventStatus
	statusCode    string
	prodCode      string
	resCode       string
	prodQuantity  uint32
	scrapQuantity uint32
	dateStart     time.Time
	dateEnd       *time.Time
	shiftID       string
	userID        string
	createdAt     time.Time
	updatedAt     time.Time
}

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
