package persistence

import (
	"context"
	"sync"
	"testing"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
)

func TestNewMemoryResourceRepository(t *testing.T) {
	repo := NewMemoryResourceRepository()
	if repo == nil {
		t.Fatal("Expected repository to be created")
	}
}

func TestMemoryResourceRepository_Save(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("save new resource", func(t *testing.T) {
		resource := entity.NewResource("SP01", "RES001", nil, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify it was saved
		found, err := repo.FindByID(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Expected no error finding resource, got %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find saved resource")
		}
		if found.ID() != resource.ID() {
			t.Errorf("Expected ID %s, got %s", resource.ID(), found.ID())
		}
	})

	t.Run("save resource with shiftID", func(t *testing.T) {
		shiftID := "shift123"
		resource := entity.NewResource("SP01", "RES002", &shiftID, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		found, err := repo.FindByID(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found.ShiftID() == nil {
			t.Fatal("Expected shiftID to be set")
		}
		if *found.ShiftID() != shiftID {
			t.Errorf("Expected shiftID %s, got %s", shiftID, *found.ShiftID())
		}
	})

	t.Run("save duplicate code", func(t *testing.T) {
		code := "RES003"
		resource1 := entity.NewResource("SP01", code, nil, "machine", 5, nil)

		err := repo.Save(ctx, resource1)
		if err != nil {
			t.Fatalf("Expected no error on first save, got %v", err)
		}

		resource2 := entity.NewResource("SP01", code, nil, "tool", 3, nil)

		err = repo.Save(ctx, resource2)
		if err == nil {
			t.Fatal("Expected error when saving duplicate code")
		}
	})
}

func TestMemoryResourceRepository_Update(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("update existing resource", func(t *testing.T) {
		resource := entity.NewResource("SP01", "RES100", nil, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}

		// Update the resource
		newType := "tool"
		resource.UpdateType(newType)

		err = repo.Update(ctx, resource)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify update
		found, err := repo.FindByID(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Failed to find resource: %v", err)
		}
		if found.Type() != newType {
			t.Errorf("Expected type %s, got %s", newType, found.Type())
		}
	})

	t.Run("update non-existent resource", func(t *testing.T) {
		resource := entity.NewResource("SP01", "RES101", nil, "machine", 5, nil)

		// Try to update without saving first
		err := repo.Update(ctx, resource)
		if err == nil {
			t.Fatal("Expected error when updating non-existent resource")
		}
	})
}

func TestMemoryResourceRepository_Delete(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("delete existing resource", func(t *testing.T) {
		resource := entity.NewResource("SP01", "RES200", nil, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}

		err = repo.Delete(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Verify deletion - should return error
		_, err = repo.FindByID(ctx, resource.ID())
		if err == nil {
			t.Error("Expected error for deleted resource")
		}
	})

	t.Run("delete non-existent resource", func(t *testing.T) {
		nonExistentID := uuid.New()
		err := repo.Delete(ctx, nonExistentID)
		if err == nil {
			t.Fatal("Expected error when deleting non-existent resource")
		}
	})
}

func TestMemoryResourceRepository_FindByID(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("find existing resource", func(t *testing.T) {
		resource := entity.NewResource("SP01", "RES300", nil, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}

		found, err := repo.FindByID(ctx, resource.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find resource")
		}
		if found.ID() != resource.ID() {
			t.Errorf("Expected ID %s, got %s", resource.ID(), found.ID())
		}
	})

	t.Run("find non-existent resource", func(t *testing.T) {
		nonExistentID := uuid.New()
		_, err := repo.FindByID(ctx, nonExistentID)
		if err == nil {
			t.Fatal("Expected error for non-existent resource")
		}
	})
}

func TestMemoryResourceRepository_FindByCode(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("find existing resource by code", func(t *testing.T) {
		code := "RES400"
		resource := entity.NewResource("SP01", code, nil, "machine", 5, nil)

		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}

		found, err := repo.FindByCode(ctx, code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if found == nil {
			t.Fatal("Expected to find resource")
		}
		if found.Code() != code {
			t.Errorf("Expected code %s, got %s", code, found.Code())
		}
	})

	t.Run("find non-existent code", func(t *testing.T) {
		_, err := repo.FindByCode(ctx, "NONEXISTENT")
		if err == nil {
			t.Fatal("Expected error for non-existent code")
		}
	})
}

func TestMemoryResourceRepository_ExistsByCode(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	code := "RES500"
	resource := entity.NewResource("SP01", code, nil, "machine", 5, nil)

	err := repo.Save(ctx, resource)
	if err != nil {
		t.Fatalf("Failed to save resource: %v", err)
	}

	t.Run("existing code", func(t *testing.T) {
		exists, err := repo.ExistsByCode(ctx, code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !exists {
			t.Error("Expected code to exist")
		}
	})

	t.Run("non-existent code", func(t *testing.T) {
		exists, err := repo.ExistsByCode(ctx, "NONEXISTENT")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if exists {
			t.Error("Expected code to not exist")
		}
	})
}

func TestMemoryResourceRepository_FindAll(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	// Create multiple resources
	for i := 0; i < 5; i++ {
		code := "RES600" + string(rune('A'+i))
		resource := entity.NewResource("SP01", code, nil, "machine", int16(i), nil) // nolint:gosec
		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}
	}

	t.Run("find all with no pagination", func(t *testing.T) {
		resources, err := repo.FindAll(ctx, 100, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) < 5 {
			t.Errorf("Expected at least 5 resources, got %d", len(resources))
		}
	})

	t.Run("find all with limit", func(t *testing.T) {
		resources, err := repo.FindAll(ctx, 2, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
	})

	t.Run("find all with zero limit", func(t *testing.T) {
		resources, err := repo.FindAll(ctx, 0, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) < 5 {
			t.Errorf("Expected at least 5 resources, got %d", len(resources))
		}
	})
}

func TestMemoryResourceRepository_FindByType(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	// Create resources with different types
	types := []string{"machine", "machine", "tool", "tool", "equipment"}
	for i, resType := range types {
		code := "RES700" + string(rune('A'+i))
		resource := entity.NewResource("SP01", code, nil, resType, int16(i), nil) // nolint:gosec
		err := repo.Save(ctx, resource)
		if err != nil {
			t.Fatalf("Failed to save resource: %v", err)
		}
	}

	t.Run("find by existing type", func(t *testing.T) {
		resources, err := repo.FindByType(ctx, "machine", 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
		for _, r := range resources {
			if r.Type() != "machine" {
				t.Errorf("Expected type machine, got %s", r.Type())
			}
		}
	})

	t.Run("find by non-existent type", func(t *testing.T) {
		resources, err := repo.FindByType(ctx, "nonexistent", 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 0 {
			t.Errorf("Expected 0 resources, got %d", len(resources))
		}
	})

	t.Run("find by type with pagination", func(t *testing.T) {
		resources, err := repo.FindByType(ctx, "tool", 1, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 1 {
			t.Errorf("Expected 1 resource, got %d", len(resources))
		}
	})
}

func TestMemoryResourceRepository_FindByShiftID(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	// Create resources with different shift IDs
	shift1 := "shift001"
	shift2 := "shift002"

	resource1 := entity.NewResource("SP01", "RES800", &shift1, "machine", 1, nil)
	resource2 := entity.NewResource("SP01", "RES801", &shift1, "tool", 2, nil)
	resource3 := entity.NewResource("SP01", "RES802", &shift2, "machine", 3, nil)
	resource4 := entity.NewResource("SP01", "RES803", nil, "tool", 4, nil)

	if err := repo.Save(ctx, resource1); err != nil {
		t.Fatalf("Failed to save resource1: %v", err)
	}
	if err := repo.Save(ctx, resource2); err != nil {
		t.Fatalf("Failed to save resource2: %v", err)
	}
	if err := repo.Save(ctx, resource3); err != nil {
		t.Fatalf("Failed to save resource3: %v", err)
	}
	if err := repo.Save(ctx, resource4); err != nil {
		t.Fatalf("Failed to save resource4: %v", err)
	}

	t.Run("find by existing shift ID", func(t *testing.T) {
		resources, err := repo.FindByShiftID(ctx, shift1, 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(resources))
		}
		for _, r := range resources {
			if r.ShiftID() == nil || *r.ShiftID() != shift1 {
				t.Errorf("Expected shiftID %s, got %v", shift1, r.ShiftID())
			}
		}
	})

	t.Run("find by non-existent shift ID", func(t *testing.T) {
		resources, err := repo.FindByShiftID(ctx, "nonexistent", 10, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 0 {
			t.Errorf("Expected 0 resources, got %d", len(resources))
		}
	})

	t.Run("find by shift ID with pagination", func(t *testing.T) {
		resources, err := repo.FindByShiftID(ctx, shift1, 1, 0)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(resources) != 1 {
			t.Errorf("Expected 1 resource, got %d", len(resources))
		}
	})
}

func TestMemoryResourceRepository_ThreadSafety(t *testing.T) {
	repo := NewMemoryResourceRepository()
	ctx := context.Background()

	t.Run("concurrent saves", func(t *testing.T) {
		var wg sync.WaitGroup
		errors := make(chan error, 10)

		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				code := "CONCURRENT" + string(rune('A'+idx))
				resource := entity.NewResource("SP01", code, nil, "machine", int16(idx), nil) // nolint:gosec
				if err := repo.Save(ctx, resource); err != nil {
					errors <- err
				}
			}(i)
		}

		wg.Wait()
		close(errors)

		// Check for errors
		for err := range errors {
			t.Errorf("Concurrent save error: %v", err)
		}

		// Verify all were saved
		resources, err := repo.FindAll(ctx, 100, 0)
		if err != nil {
			t.Fatalf("Failed to find all: %v", err)
		}

		// Count resources with CONCURRENT prefix
		count := 0
		for _, r := range resources {
			if len(r.Code()) >= 10 && r.Code()[:10] == "CONCURRENT" {
				count++
			}
		}

		if count != 10 {
			t.Errorf("Expected 10 concurrent resources, got %d", count)
		}
	})

	t.Run("concurrent reads", func(t *testing.T) {
		// Save a resource first
		resource := entity.NewResource("SP01", "READTEST", nil, "machine", 5, nil)
		if err := repo.Save(ctx, resource); err != nil {
			t.Fatalf("Failed to save test resource: %v", err)
		}

		var wg sync.WaitGroup
		errors := make(chan error, 20)

		// Concurrent reads
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, err := repo.FindByCode(ctx, "READTEST")
				if err != nil {
					errors <- err
				}
			}()
		}

		wg.Wait()
		close(errors)

		for err := range errors {
			t.Errorf("Concurrent read error: %v", err)
		}
	})
}
