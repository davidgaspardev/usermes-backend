package dto

import (
	"errors"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/domain/valueobjects"
)

type TaskDTO struct {
	Name      string   `json:"name"`
	ResCode   string   `json:"resCode"`
	CronCode  string   `json:"cronCode"`
	DateStart string   `json:"dateStart"`
	DateEnd   string   `json:"dateEnd"`
	Deps      []uint64 `json:"deps"`
}

func (task *TaskDTO) Validate() error {
	if task.Name != "" && task.ResCode != "" && task.CronCode != "" && task.DateStart != "" && task.DateEnd != "" {
		return nil
	}

	return errors.New("invalid request body")
}

func (task *TaskDTO) GetValueObject() (valueobjects.CreateTask, error) {
	dateStart, err := time.Parse(time.RFC3339, task.DateStart)
	if err != nil {
		return valueobjects.CreateTask{}, err
	}

	dateEnd, err := time.Parse(time.RFC3339, task.DateEnd)
	if err != nil {
		return valueobjects.CreateTask{}, err
	}

	var createTask = valueobjects.CreateTask{
		Name:      task.Name,
		ResCode:   task.ResCode,
		CronCode:  task.CronCode,
		DateStart: dateStart,
		DateEnd:   dateEnd,
		Deps:      task.Deps,
	}
	return createTask, nil
}
