package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type deleteShiftPatternUseCaseImpl struct {
	shiftPatternRepository output.ShiftPatternRepository
}

// NewDeleteShiftPatternUseCase creates a new DeleteShiftPatternUseCase instance.
func NewDeleteShiftPatternUseCase(repo output.ShiftPatternRepository) input.DeleteShiftPatternUseCase {
	return &deleteShiftPatternUseCaseImpl{shiftPatternRepository: repo}
}

func (uc *deleteShiftPatternUseCaseImpl) Execute(ctx context.Context, id uuid.UUID) error {
	pattern, err := uc.shiftPatternRepository.FindByID(id)
	if err != nil {
		return err
	}
	if pattern == nil {
		return domainerrors.ErrShiftPatternNotFound
	}

	return uc.shiftPatternRepository.Delete(id)
}
