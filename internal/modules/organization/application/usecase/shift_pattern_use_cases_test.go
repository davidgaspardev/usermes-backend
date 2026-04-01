package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	domainerrors "github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/adapter/output/persistence"
)

func mustParseHHMM(s string) time.Time {
	t, err := time.Parse("15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func defaultEntryCommands() []input.ShiftEntryCommand {
	return []input.ShiftEntryCommand{
		{DayIndex: 0, Name: "Morning", StartTime: mustParseHHMM("06:00"), EndTime: mustParseHHMM("14:00")},
		{DayIndex: 1, Name: "Afternoon", StartTime: mustParseHHMM("14:00"), EndTime: mustParseHHMM("22:00")},
		{DayIndex: 2, Name: "Night", StartTime: mustParseHHMM("22:00"), EndTime: mustParseHHMM("06:00")},
	}
}

// ---- Create ----------------------------------------------------------------

func TestCreateShiftPatternUseCase_Execute(t *testing.T) {
	ctx := context.Background()
	refStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("creates pattern and links to location", func(t *testing.T) {
		locationRepo := persistence.NewLocationRepositoryInMemory()
		patternRepo := persistence.NewMemoryShiftPatternRepository()

		loc, err := NewCreateLocationRootUseCase(locationRepo).Execute(ctx, input.CreateLocationRootCommand{
			Code: "PLANT01", Name: "Plant",
		})
		if err != nil || loc == nil {
			t.Fatalf("failed to seed location: %v", err)
		}

		cmd := input.CreateShiftPatternCommand{
			LocationCode: "PLANT01",
			Name:         "3-Shift",
			RefStartDate: refStart,
			CycleLength:  3,
			Entries:      defaultEntryCommands(),
		}

		pattern, err := NewCreateShiftPatternUseCase(locationRepo, patternRepo).Execute(ctx, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if pattern == nil {
			t.Fatal("expected pattern, got nil")
		}
		if pattern.Name() != "3-Shift" {
			t.Errorf("expected name 3-Shift, got %s", pattern.Name())
		}
		if pattern.CycleLength() != 3 {
			t.Errorf("expected cycleLength 3, got %d", pattern.CycleLength())
		}
		if len(pattern.Entries()) != 3 {
			t.Errorf("expected 3 entries, got %d", len(pattern.Entries()))
		}
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		locationRepo := persistence.NewLocationRepositoryInMemory()
		patternRepo := persistence.NewMemoryShiftPatternRepository()

		_, err := NewCreateShiftPatternUseCase(locationRepo, patternRepo).Execute(ctx, input.CreateShiftPatternCommand{})
		if err != domainerrors.ErrInvalidShiftPatternCommand {
			t.Errorf("expected ErrInvalidShiftPatternCommand, got %v", err)
		}
	})

	t.Run("returns error when location not found", func(t *testing.T) {
		locationRepo := persistence.NewLocationRepositoryInMemory()
		patternRepo := persistence.NewMemoryShiftPatternRepository()

		_, err := NewCreateShiftPatternUseCase(locationRepo, patternRepo).Execute(ctx, input.CreateShiftPatternCommand{
			LocationCode: "MISSING",
			Name:         "X",
			RefStartDate: refStart,
			CycleLength:  1,
			Entries:      []input.ShiftEntryCommand{{DayIndex: 0, Name: "S1", StartTime: mustParseHHMM("06:00"), EndTime: mustParseHHMM("14:00")}},
		})
		if err != domainerrors.ErrLocationNotFound {
			t.Errorf("expected ErrLocationNotFound, got %v", err)
		}
	})

	t.Run("returns error when pattern name already exists", func(t *testing.T) {
		locationRepo := persistence.NewLocationRepositoryInMemory()
		patternRepo := persistence.NewMemoryShiftPatternRepository()

		if _, err := NewCreateLocationRootUseCase(locationRepo).Execute(ctx, input.CreateLocationRootCommand{
			Code: "PLANT01", Name: "Plant",
		}); err != nil {
			t.Fatalf("seed: %v", err)
		}

		cmd := input.CreateShiftPatternCommand{
			LocationCode: "PLANT01",
			Name:         "Duplicate",
			RefStartDate: refStart,
			CycleLength:  1,
			Entries:      []input.ShiftEntryCommand{{DayIndex: 0, Name: "S1", StartTime: mustParseHHMM("06:00"), EndTime: mustParseHHMM("14:00")}},
		}

		if _, err := NewCreateShiftPatternUseCase(locationRepo, patternRepo).Execute(ctx, cmd); err != nil {
			t.Fatalf("first create: %v", err)
		}

		// second location to avoid ErrLocationNotFound
		if _, err := NewCreateLocationRootUseCase(locationRepo).Execute(ctx, input.CreateLocationRootCommand{
			Code: "PLANT02", Name: "Plant 2",
		}); err != nil {
			t.Fatalf("seed 2: %v", err)
		}
		cmd.LocationCode = "PLANT02"

		_, err := NewCreateShiftPatternUseCase(locationRepo, patternRepo).Execute(ctx, cmd)
		if err != domainerrors.ErrShiftPatternAlreadyExists {
			t.Errorf("expected ErrShiftPatternAlreadyExists, got %v", err)
		}
	})
}

// ---- Get by ID -------------------------------------------------------------

func TestGetShiftPatternByIDUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns pattern by id", func(t *testing.T) {
		repo := seededShiftPatternRepo(t)
		pattern, _ := repo.GetDefault()

		result, err := NewGetShiftPatternByIDUseCase(repo).Execute(ctx, pattern.ID())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil || result.ID() != pattern.ID() {
			t.Error("expected pattern with same ID")
		}
	})

	t.Run("returns error for unknown id", func(t *testing.T) {
		repo := persistence.NewMemoryShiftPatternRepository()

		_, err := NewGetShiftPatternByIDUseCase(repo).Execute(ctx, uuid.New())
		if err != domainerrors.ErrShiftPatternNotFound {
			t.Errorf("expected ErrShiftPatternNotFound, got %v", err)
		}
	})
}

// ---- Update ----------------------------------------------------------------

func TestUpdateShiftPatternUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("updates pattern successfully", func(t *testing.T) {
		repo := seededShiftPatternRepo(t)
		existing, _ := repo.GetDefault()

		newEntries := []input.ShiftEntryCommand{
			{DayIndex: 0, Name: "Day", StartTime: mustParseHHMM("08:00"), EndTime: mustParseHHMM("16:00")},
		}
		cmd := input.UpdateShiftPatternCommand{
			ID:           existing.ID(),
			Name:         "Updated",
			RefStartDate: existing.RefStartDate(),
			CycleLength:  1,
			Entries:      newEntries,
		}

		updated, err := NewUpdateShiftPatternUseCase(repo).Execute(ctx, cmd)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updated.Name() != "Updated" {
			t.Errorf("expected name Updated, got %s", updated.Name())
		}
		if updated.CycleLength() != 1 {
			t.Errorf("expected cycleLength 1, got %d", updated.CycleLength())
		}
		if len(updated.Entries()) != 1 {
			t.Errorf("expected 1 entry, got %d", len(updated.Entries()))
		}
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		repo := persistence.NewMemoryShiftPatternRepository()

		_, err := NewUpdateShiftPatternUseCase(repo).Execute(ctx, input.UpdateShiftPatternCommand{})
		if err != domainerrors.ErrInvalidShiftPatternCommand {
			t.Errorf("expected ErrInvalidShiftPatternCommand, got %v", err)
		}
	})

	t.Run("returns error for unknown id", func(t *testing.T) {
		repo := persistence.NewMemoryShiftPatternRepository()

		_, err := NewUpdateShiftPatternUseCase(repo).Execute(ctx, input.UpdateShiftPatternCommand{
			ID:           uuid.New(),
			Name:         "X",
			RefStartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			CycleLength:  1,
			Entries:      []input.ShiftEntryCommand{{DayIndex: 0, Name: "S", StartTime: mustParseHHMM("06:00"), EndTime: mustParseHHMM("14:00")}},
		})
		if err != domainerrors.ErrShiftPatternNotFound {
			t.Errorf("expected ErrShiftPatternNotFound, got %v", err)
		}
	})
}

// ---- Delete ----------------------------------------------------------------

func TestDeleteShiftPatternUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("deletes pattern successfully", func(t *testing.T) {
		repo := seededShiftPatternRepo(t)
		existing, _ := repo.GetDefault()

		err := NewDeleteShiftPatternUseCase(repo).Execute(ctx, existing.ID())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		result, _ := repo.FindByID(existing.ID())
		if result != nil {
			t.Error("expected pattern to be deleted")
		}
	})

	t.Run("returns error for unknown id", func(t *testing.T) {
		repo := persistence.NewMemoryShiftPatternRepository()

		err := NewDeleteShiftPatternUseCase(repo).Execute(ctx, uuid.New())
		if err != domainerrors.ErrShiftPatternNotFound {
			t.Errorf("expected ErrShiftPatternNotFound, got %v", err)
		}
	})
}
