package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
)

// --- stub resource repository ---

type stubResourceRepo struct {
	existsErr error
	saveErr   error
	deleteErr error
	exists    bool
}

func (r *stubResourceRepo) ExistsByCode(_ context.Context, _ string) (bool, error) {
	return r.exists, r.existsErr
}
func (r *stubResourceRepo) Save(_ context.Context, _ *entity.Resource) error { return r.saveErr }
func (r *stubResourceRepo) Delete(_ context.Context, _ uuid.UUID) error      { return r.deleteErr }
func (r *stubResourceRepo) Update(_ context.Context, _ *entity.Resource) error {
	return nil
}
func (r *stubResourceRepo) FindByID(_ context.Context, _ uuid.UUID) (*entity.Resource, error) {
	return nil, nil
}
func (r *stubResourceRepo) FindByCode(_ context.Context, _ string) (*entity.Resource, error) {
	return nil, nil
}
func (r *stubResourceRepo) FindAll(_ context.Context, _, _ int) ([]*entity.Resource, error) {
	return nil, nil
}
func (r *stubResourceRepo) FindByType(_ context.Context, _ string, _, _ int) ([]*entity.Resource, error) {
	return nil, nil
}
func (r *stubResourceRepo) FindByShiftID(_ context.Context, _ string, _, _ int) ([]*entity.Resource, error) {
	return nil, nil
}

// --- stub event repository ---

type stubEventRepo struct {
	createErr error
}

func (r *stubEventRepo) Create(_ *entity.Event) error                        { return r.createErr }
func (r *stubEventRepo) GetCurrentByResCode(_ string) (*entity.Event, error) { return nil, nil }
func (r *stubEventRepo) UpdateCurrent(_ *entity.Event) error                 { return nil }

// --- helpers ---

func validCommand() *input.CreateResourceCommand {
	shiftID := "shift-001"
	return &input.CreateResourceCommand{
		PlantCode:    "SP01",
		Code:         "MACHINE-001",
		ShiftID:      &shiftID,
		ResourceType: "CNC",
		StopFactor:   10,
		WhoCreated:   "user-001",
	}
}

// --- tests ---

func TestCreateResourceUseCase_Execute_Success(t *testing.T) {
	uc := NewCreateResourceUseCase(&stubResourceRepo{}, &stubEventRepo{})

	err := uc.Execute(context.Background(), validCommand())

	if err != nil {
		t.Errorf("expected nil error, got %v", err)
	}
}

func TestCreateResourceUseCase_Execute_AlreadyExists(t *testing.T) {
	uc := NewCreateResourceUseCase(&stubResourceRepo{exists: true}, &stubEventRepo{})

	err := uc.Execute(context.Background(), validCommand())

	if !errors.Is(err, domainerrors.ErrResourceAlreadyExists) {
		t.Errorf("expected ErrResourceAlreadyExists, got %v", err)
	}
}

func TestCreateResourceUseCase_Execute_ExistsCheckError(t *testing.T) {
	repoErr := errors.New("db error")
	uc := NewCreateResourceUseCase(&stubResourceRepo{existsErr: repoErr}, &stubEventRepo{})

	err := uc.Execute(context.Background(), validCommand())

	if !errors.Is(err, repoErr) {
		t.Errorf("expected db error, got %v", err)
	}
}

func TestCreateResourceUseCase_Execute_SaveError(t *testing.T) {
	saveErr := errors.New("save failed")
	uc := NewCreateResourceUseCase(&stubResourceRepo{saveErr: saveErr}, &stubEventRepo{})

	err := uc.Execute(context.Background(), validCommand())

	if !errors.Is(err, saveErr) {
		t.Errorf("expected save error, got %v", err)
	}
}

func TestCreateResourceUseCase_Execute_EventCreateError_Rollback(t *testing.T) {
	eventErr := errors.New("event create failed")
	uc := NewCreateResourceUseCase(&stubResourceRepo{}, &stubEventRepo{createErr: eventErr})

	err := uc.Execute(context.Background(), validCommand())

	if !errors.Is(err, eventErr) {
		t.Errorf("expected event error, got %v", err)
	}
}

func TestCreateResourceUseCase_Execute_EventCreateError_DeleteError(t *testing.T) {
	eventErr := errors.New("event create failed")
	deleteErr := errors.New("delete failed")
	uc := NewCreateResourceUseCase(
		&stubResourceRepo{deleteErr: deleteErr},
		&stubEventRepo{createErr: eventErr},
	)

	err := uc.Execute(context.Background(), validCommand())

	if !errors.Is(err, deleteErr) {
		t.Errorf("expected delete error during rollback, got %v", err)
	}
}
