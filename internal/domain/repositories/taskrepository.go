package repostories

import "github.com/davidgaspardev/usermes-backend/internal/domain/entities"

type TaskEntityRef = *entities.TaskEntity

type TaskRepository interface {
	FindById(id uint64) (TaskEntityRef, error)
}
