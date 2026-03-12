package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/infrastructure/adapter/output/persistence"
)

func TestNewResourceService(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)

	if service == nil {
		t.Fatal("Expected service to be created")
	}
}

func TestResourceService_Create(t *testing.T) { //nolint:gocyclo
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		code := "RES001"
		resourceType := "machine"
		stopFactor := int16(5)

		resource, err := service.Create(ctx, "PLANT-01", code, resourceType, stopFactor, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if resource.Code() != code {
			t.Errorf("Expected code %s, got %s", code, resource.Code())
		}
		if resource.Type() != resourceType {
			t.Errorf("Expected type %s, got %s", resourceType, resource.Type())
		}
		if resource.StopFactor() != stopFactor {
			t.Errorf("Expected stopFactor %d, got %d", stopFactor, resource.StopFactor())
		}
	})

	t.Run("invalid code - empty", func(t *testing.T) {
		_, err := service.Create(ctx, "PLANT-01", "", "machine", 5, nil)
		if err == nil {
			t.Fatal("Expected error for empty code")
		}
		if err != errors.ErrInvalidCode {
			t.Errorf("Expected ErrInvalidCode, got %v", err)
		}
	})

	t.Run("invalid code - too long", func(t *testing.T) {
		longCode := "THISISAVERYLONGCODETHATEXCEEDSFIFTYCHARACTERSLIMIT123456789"
		_, err := service.Create(ctx, "PLANT-01", longCode, "machine", 5, nil)
		if err == nil {
			t.Fatal("Expected error for code too long")
		}
		if err != errors.ErrInvalidCode {
			t.Errorf("Expected ErrInvalidCode, got %v", err)
		}
	})

	t.Run("invalid type - empty", func(t *testing.T) {
		_, err := service.Create(ctx, "PLANT-01", "RES003", "", 5, nil)
		if err == nil {
			t.Fatal("Expected error for empty type")
		}
		if err != errors.ErrInvalidType {
			t.Errorf("Expected ErrInvalidType, got %v", err)
		}
	})

	t.Run("invalid type - too long", func(t *testing.T) {
		longType := "THISISAVERYLONGTYPETHATEXCEEDSFIFTYCHARACTERSLIMIT1234567890"
		_, err := service.Create(ctx, "PLANT-01", "RES004", longType, 5, nil)
		if err == nil {
			t.Fatal("Expected error for type too long")
		}
		if err != errors.ErrInvalidType {
			t.Errorf("Expected ErrInvalidType, got %v", err)
		}
	})

	t.Run("invalid stopFactor - negative", func(t *testing.T) {
		_, err := service.Create(ctx, "PLANT-01", "RES005", "machine", -1, nil)
		if err == nil {
			t.Fatal("Expected error for negative stopFactor")
		}
		if err != errors.ErrInvalidStopFactor {
			t.Errorf("Expected ErrInvalidStopFactor, got %v", err)
		}
	})

	t.Run("duplicate code", func(t *testing.T) {
		code := "RES006"
		_, err := service.Create(ctx, "PLANT-01", code, "machine", 5, nil)
		if err != nil {
			t.Fatalf("Expected no error on first create, got %v", err)
		}

		_, err = service.Create(ctx, "PLANT-01", code, "tool", 3, nil)
		if err == nil {
			t.Fatal("Expected error for duplicate code")
		}
		if err != errors.ErrResourceCodeAlreadyExists {
			t.Errorf("Expected ErrResourceCodeAlreadyExists, got %v", err)
		}
	})
}

func TestResourceService_Update(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	initialResource, err := service.Create(ctx, "PLANT-01", "RES100", "machine", 5, nil)
	if err != nil {
		t.Fatalf("Failed to create initial resource: %v", err)
	}

	t.Run("successful update", func(t *testing.T) {
		newCode := "RES101"
		newType := "tool"
		newStopFactor := int16(10)

		updated, err := service.Update(ctx, initialResource.ID(), "PLANT-01", newCode, newType, newStopFactor, nil)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updated.Code() != newCode {
			t.Errorf("Expected code %s, got %s", newCode, updated.Code())
		}
		if updated.Type() != newType {
			t.Errorf("Expected type %s, got %s", newType, updated.Type())
		}
		if updated.StopFactor() != newStopFactor {
			t.Errorf("Expected stopFactor %d, got %d", newStopFactor, updated.StopFactor())
		}
	})

	t.Run("update with invalid code", func(t *testing.T) {
		_, err := service.Update(ctx, initialResource.ID(), "PLANT-01", "", "machine", 5, nil)
		if err == nil {
			t.Fatal("Expected error for invalid code")
		}
		if err != errors.ErrInvalidCode {
			t.Errorf("Expected ErrInvalidCode, got %v", err)
		}
	})

	t.Run("update with invalid type", func(t *testing.T) {
		_, err := service.Update(ctx, initialResource.ID(), "PLANT-01", "RES103", "", 5, nil)
		if err == nil {
			t.Fatal("Expected error for invalid type")
		}
		if err != errors.ErrInvalidType {
			t.Errorf("Expected ErrInvalidType, got %v", err)
		}
	})

	t.Run("update with invalid stopFactor", func(t *testing.T) {
		_, err := service.Update(ctx, initialResource.ID(), "PLANT-01", "RES104", "machine", -5, nil)
		if err == nil {
			t.Fatal("Expected error for invalid stopFactor")
		}
		if err != errors.ErrInvalidStopFactor {
			t.Errorf("Expected ErrInvalidStopFactor, got %v", err)
		}
	})
}

func TestResourceService_Delete(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	t.Run("successful delete", func(t *testing.T) {
		resource, err := service.Create(ctx, "PLANT-01", "RES200", "machine", 5, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		err = service.Delete(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = service.GetByID(ctx, resource.ID())
		if err == nil {
			t.Error("Expected error for deleted resource")
		}
	})

	t.Run("delete non-existent resource", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := service.Delete(ctx, nonExistentID)
		if err == nil {
			t.Fatal("Expected error for non-existent resource")
		}
	})
}

func TestResourceService_GetByID(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	t.Run("get existing resource", func(t *testing.T) {
		created, err := service.Create(ctx, "PLANT-01", "RES300", "machine", 5, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		found, err := service.GetByID(ctx, created.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find resource")
		}
		if found.ID() != created.ID() {
			t.Errorf("Expected ID %s, got %s", created.ID(), found.ID())
		}
	})

	t.Run("get non-existent resource", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := service.GetByID(ctx, nonExistentID)
		if err == nil {
			t.Fatal("Expected error for non-existent resource")
		}
	})
}

func TestResourceService_GetByCode(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	t.Run("get existing resource by code", func(t *testing.T) {
		code := "RES400"
		created, err := service.Create(ctx, "PLANT-01", code, "machine", 5, nil)
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}

		found, err := service.GetByCode(ctx, "PLANT-01", code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find resource")
		}
		if found.Code() != code {
			t.Errorf("Expected code %s, got %s", code, found.Code())
		}
		if found.ID() != created.ID() {
			t.Error("Expected same resource")
		}
	})

	t.Run("get non-existent resource by code", func(t *testing.T) {
		_, err := service.GetByCode(ctx, "PLANT-01", "NONEXISTENT")
		if err == nil {
			t.Fatal("Expected error for non-existent resource")
		}
	})
}

func TestResourceService_GetAll(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		code := "RES500" + string(rune('A'+i))
		_, err := service.Create(ctx, "PLANT-01", code, "machine", int16(i), nil) // nolint:gosec
		if err != nil {
			t.Fatalf("Failed to create resource: %v", err)
		}
	}

	t.Run("get all with pagination", func(t *testing.T) {
		resources, err := service.GetAll(ctx, 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) < 5 {
			t.Errorf("Expected at least 5 resources, got %d", len(resources))
		}
	})

	t.Run("get with limit", func(t *testing.T) {
		resources, err := service.GetAll(ctx, 2, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
	})
}

func TestResourceService_GetByType(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	if _, err := service.Create(ctx, "PLANT-01", "RES600", "machine", 1, nil); err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}
	if _, err := service.Create(ctx, "PLANT-01", "RES601", "machine", 2, nil); err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}
	if _, err := service.Create(ctx, "PLANT-01", "RES602", "tool", 3, nil); err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	t.Run("get by existing type", func(t *testing.T) {
		resources, err := service.GetByType(ctx, "machine", 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
	})

	t.Run("get by non-existent type", func(t *testing.T) {
		resources, err := service.GetByType(ctx, "nonexistent", 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 0 {
			t.Errorf("Expected 0 resources, got %d", len(resources))
		}
	})
}

func TestResourceService_UpdatedAt(t *testing.T) {
	repo := persistence.NewMemoryResourceRepository()
	service := NewResourceService(repo)
	ctx := context.Background()

	resource, err := service.Create(ctx, "PLANT-01", "RES800", "machine", 5, nil)
	if err != nil {
		t.Fatalf("Failed to create resource: %v", err)
	}

	createdAt := resource.CreatedAt()
	updatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)

	updated, err := service.Update(ctx, resource.ID(), "PLANT-01", "RES801", "machine", 5, nil)
	if err != nil {
		t.Fatalf("Failed to update resource: %v", err)
	}

	if !updated.CreatedAt().Equal(createdAt) {
		t.Error("CreatedAt should not change on update")
	}

	if !updated.UpdatedAt().After(updatedAt) {
		t.Error("UpdatedAt should be more recent after update")
	}
}
