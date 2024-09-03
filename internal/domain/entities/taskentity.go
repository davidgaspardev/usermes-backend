package entities

import "github.com/davidgaspardev/usermes-backend/internal/domain/valueobjects"

type TaskEntity struct {
	ID uint64
	valueobjects.CreateTask
}
