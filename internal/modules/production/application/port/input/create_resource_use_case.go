package input

import "context"

// CreateResourceCommand holds the input data for creating a new resource.
type CreateResourceCommand struct {
	ShiftID      *string
	PlantCode    string
	Code         string
	ResourceType string
	WhoCreated   string
	Tags         []string
	StopFactor   int16
}

// CreateResourceUseCase is the input port for creating a new resource.
type CreateResourceUseCase interface {
	Execute(ctx context.Context, command *CreateResourceCommand) error
}
