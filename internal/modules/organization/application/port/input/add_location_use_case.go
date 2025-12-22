package input

import (
	"context"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

type AddLocationCommand struct {
	Code       string
	Name       string
	ParentCode string
	RootCode   string
}

type AddLocationUseCase interface {
	Execute(ctx *context.Context, command AddLocationCommand) (*entity.Location, error)
}
