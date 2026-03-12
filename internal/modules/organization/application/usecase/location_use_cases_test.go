package usecase

import (
	"context"
	"testing"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/input"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/adapter/output/persistence"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/infrastructure/seeder"
)

// seededShiftPatternRepo returns an in-memory shift pattern repo pre-loaded with the default pattern.
func seededShiftPatternRepo(t *testing.T) output.ShiftPatternRepository {
	t.Helper()
	repo := persistence.NewMemoryShiftPatternRepository()
	if err := seeder.SeedDefaultShiftPatterns(repo); err != nil {
		t.Fatalf("failed to seed shift patterns: %v", err)
	}
	return repo
}

func TestCreateLocationRootUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("creates root location successfully", func(t *testing.T) {
		repo := persistence.NewLocationRepositoryInMemory()
		uc := NewCreateLocationRootUseCase(repo)

		loc, err := uc.Execute(ctx, input.CreateLocationRootCommand{
			Code: "PLANT01",
			Name: "Plant One",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if loc == nil {
			t.Fatal("Expected location, got nil")
		}
		if loc.Code() != "PLANT01" {
			t.Errorf("Expected code PLANT01, got %s", loc.Code())
		}
		// Plants have no shift pattern
		if loc.ShiftPatternID().String() != "00000000-0000-0000-0000-000000000000" {
			t.Error("Expected zero ShiftPatternID for plant location")
		}
	})

	t.Run("returns error when location already exists", func(t *testing.T) {
		repo := persistence.NewLocationRepositoryInMemory()
		uc := NewCreateLocationRootUseCase(repo)
		cmd := input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"}

		_, err := uc.Execute(ctx, cmd)
		if err != nil {
			t.Fatalf("First creation should succeed, got %v", err)
		}

		_, err = uc.Execute(ctx, cmd)
		if err != errors.ErrLocationAlreadyExists {
			t.Errorf("Expected ErrLocationAlreadyExists, got %v", err)
		}
	})
}

func TestAddLocationUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	setup := func() (input.AddLocationUseCase, input.CreateLocationRootUseCase) {
		repo := persistence.NewLocationRepositoryInMemory()
		return NewAddLocationUseCase(repo, seededShiftPatternRepo(t)), NewCreateLocationRootUseCase(repo)
	}

	t.Run("adds child location successfully", func(t *testing.T) {
		addUC, createUC := setup()

		_, err := createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		tree, err := addUC.Execute(ctx, input.AddLocationCommand{
			Code:       "AREA01",
			Name:       "Area One",
			Kind:       "AREA",
			ParentCode: "PLANT01",
			RootCode:   "PLANT01",
		})

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if tree == nil {
			t.Fatal("Expected tree, got nil")
		}
		if !tree.ExistsByCode("AREA01") {
			t.Error("Expected AREA01 to exist in tree")
		}
	})

	t.Run("returns error when root tree not found", func(t *testing.T) {
		addUC, _ := setup()

		_, err := addUC.Execute(ctx, input.AddLocationCommand{
			Code:       "AREA01",
			Name:       "Area One",
			Kind:       "AREA",
			ParentCode: "PLANT01",
			RootCode:   "PLANT01",
		})

		if err != errors.ErrLocationTreeNotFound {
			t.Errorf("Expected ErrLocationTreeNotFound, got %v", err)
		}
	})

	t.Run("returns error when location already exists", func(t *testing.T) {
		addUC, createUC := setup()

		_, err := createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		cmd := input.AddLocationCommand{
			Code: "AREA01", Name: "Area One", Kind: "AREA",
			ParentCode: "PLANT01", RootCode: "PLANT01",
		}
		_, err = addUC.Execute(ctx, cmd)
		if err != nil {
			t.Fatalf("First add should succeed, got %v", err)
		}

		_, err = addUC.Execute(ctx, cmd)
		if err != errors.ErrLocationAlreadyExists {
			t.Errorf("Expected ErrLocationAlreadyExists, got %v", err)
		}
	})

	t.Run("returns error when parent not found", func(t *testing.T) {
		addUC, createUC := setup()

		_, err := createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		_, err = addUC.Execute(ctx, input.AddLocationCommand{
			Code:       "AREA01",
			Name:       "Area One",
			Kind:       "AREA",
			ParentCode: "NONEXISTENT",
			RootCode:   "PLANT01",
		})

		if err != errors.ErrLocationNotFound {
			t.Errorf("Expected ErrLocationNotFound, got %v", err)
		}
	})

	t.Run("returns error for invalid command", func(t *testing.T) {
		addUC, _ := setup()

		_, err := addUC.Execute(ctx, input.AddLocationCommand{})
		if err != errors.ErrInvalidLocationCommand {
			t.Errorf("Expected ErrInvalidLocationCommand, got %v", err)
		}
	})

	t.Run("returns error when kind is PLANT", func(t *testing.T) {
		addUC, createUC := setup()

		_, err := createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		_, err = addUC.Execute(ctx, input.AddLocationCommand{
			Code:       "PLANT02",
			Name:       "Plant Two",
			Kind:       "PLANT",
			ParentCode: "PLANT01",
			RootCode:   "PLANT01",
		})

		if err != errors.ErrInvalidLocationCommand {
			t.Errorf("Expected ErrInvalidLocationCommand, got %v", err)
		}
	})
}

func TestGetAllLocationsUseCase_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("returns empty list when no locations", func(t *testing.T) {
		repo := persistence.NewLocationRepositoryInMemory()
		uc := NewGetAllLocationsUseCase(repo)

		locations, err := uc.Execute(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(locations) != 0 {
			t.Errorf("Expected empty list, got %d", len(locations))
		}
	})

	t.Run("returns all root locations", func(t *testing.T) {
		repo := persistence.NewLocationRepositoryInMemory()
		createUC := NewCreateLocationRootUseCase(repo)
		getAllUC := NewGetAllLocationsUseCase(repo)

		_, err := createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT01", Name: "Plant One"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}
		_, err = createUC.Execute(ctx, input.CreateLocationRootCommand{Code: "PLANT02", Name: "Plant Two"})
		if err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		locations, err := getAllUC.Execute(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(locations) != 2 {
			t.Errorf("Expected 2 locations, got %d", len(locations))
		}
	})
}
