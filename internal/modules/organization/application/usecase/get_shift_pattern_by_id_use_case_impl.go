package usecase

import (
	"context"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
)

type getShiftPatternByIDUseCaseImpl struct {
	shiftPatternRepository output.ShiftPatternRepository
}

// NewGetShiftPatternByIDUseCase creates a new GetShiftPatternByIDUseCase instance.
func NewGetShiftPatternByIDUseCase(repo output.ShiftPatternRepository) input.GetShiftPatternByIDUseCase {
	return &getShiftPatternByIDUseCaseImpl{shiftPatternRepository: repo}
}

func (uc *getShiftPatternByIDUseCaseImpl) Execute(ctx context.Context, id uuid.UUID) (*entity.ShiftPattern, error) {
	pattern, err := uc.shiftPatternRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if pattern == nil {
		return nil, domainerrors.ErrShiftPatternNotFound
	}
	return pattern, nil
}
