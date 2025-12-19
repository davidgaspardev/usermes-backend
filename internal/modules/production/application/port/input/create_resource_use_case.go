package input

import "context"

type CreateResourceCommand struct {
	PlantCode    string
	ShiftID      *string
	Code         string
	ResourceType string
	Tags         []string
	StopFactor   int16
	WhoCreated   string
}

type CreateResourceUseCase interface {
	Execute(ctx context.Context, command *CreateResourceCommand) error
}
