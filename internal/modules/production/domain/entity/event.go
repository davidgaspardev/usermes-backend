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
	Stop       EventStatus = "Stop"
)

type Event struct {
	Type          EventType
	Status        EventStatus
	StopCode      string
	ProdCode      string
	ResCode       string
	ProdQuantity  uint32
	ScrapQuantity uint32
	DateStart     time.Time
	DateEnd       *time.Time
	ShiftID       string
	UserID        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
