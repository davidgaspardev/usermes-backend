package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/domain/errors"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/plant/infrastructure/adapter/output/persistence"
)

func TestPlantService_Create(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful creation", func(t *testing.T) {
		code := "SP01"
		name := "São Paulo Plant"
		latitude := -23.5505
		longitude := -46.6333
		ownerID := uuid.New()

		plant, err := service.Create(ctx, code, name, latitude, longitude, ownerID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if plant.Code() != "SP01" { // Should be uppercase
			t.Errorf("Expected code SP01, got %s", plant.Code())
		}

		if plant.Name() != name {
			t.Errorf("Expected name %s, got %s", name, plant.Name())
		}

		if plant.Latitude() != latitude {
			t.Errorf("Expected latitude %f, got %f", latitude, plant.Latitude())
		}

		if plant.Longitude() != longitude {
			t.Errorf("Expected longitude %f, got %f", longitude, plant.Longitude())
		}

		if plant.OwnerID() != ownerID {
			t.Errorf("Expected ownerID %s, got %s", ownerID, plant.OwnerID())
		}

		if !plant.IsActive() {
			t.Error("Expected plant to be active")
		}
	})

	t.Run("code normalization to uppercase", func(t *testing.T) {
		code := "rj01"
		name := "Rio Plant"
		ownerID := uuid.New()

		plant, err := service.Create(ctx, code, name, -22.9068, -43.1729, ownerID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if plant.Code() != "RJ01" {
			t.Errorf("Expected code to be normalized to RJ01, got %s", plant.Code())
		}
	})

	t.Run("duplicate code", func(t *testing.T) {
		code := "DUP01"
		name := "Duplicate Plant"
		ownerID := uuid.New()

		_, err := service.Create(ctx, code, name, 0, 0, ownerID)
		if err != nil {
			t.Fatalf("Expected no error on first creation, got %v", err)
		}

		_, err = service.Create(ctx, code, name, 0, 0, ownerID)
		if err != errors.ErrPlantCodeAlreadyExists {
			t.Errorf("Expected ErrPlantCodeAlreadyExists, got %v", err)
		}
	})

	t.Run("invalid code - empty", func(t *testing.T) {
		_, err := service.Create(ctx, "", "Plant", 0, 0, uuid.New())
		if err != errors.ErrPlantCodeRequired {
			t.Errorf("Expected ErrPlantCodeRequired, got %v", err)
		}
	})

	t.Run("invalid code - too short", func(t *testing.T) {
		_, err := service.Create(ctx, "A", "Plant", 0, 0, uuid.New())
		if err != errors.ErrPlantCodeTooShort {
			t.Errorf("Expected ErrPlantCodeTooShort, got %v", err)
		}
	})

	t.Run("invalid code - too long", func(t *testing.T) {
		longCode := "VERYLONGCODEEXCEEDING20CHARS"
		_, err := service.Create(ctx, longCode, "Plant", 0, 0, uuid.New())
		if err != errors.ErrPlantCodeTooLong {
			t.Errorf("Expected ErrPlantCodeTooLong, got %v", err)
		}
	})

	t.Run("invalid name - empty", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST01", "", 0, 0, uuid.New())
		if err != errors.ErrPlantNameRequired {
			t.Errorf("Expected ErrPlantNameRequired, got %v", err)
		}
	})

	t.Run("invalid name - too short", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST02", "A", 0, 0, uuid.New())
		if err != errors.ErrPlantNameTooShort {
			t.Errorf("Expected ErrPlantNameTooShort, got %v", err)
		}
	})

	t.Run("invalid name - too long", func(t *testing.T) {
		longName := "This is a very long plant name that exceeds the maximum allowed length of one hundred characters for sure"
		_, err := service.Create(ctx, "TEST03", longName, 0, 0, uuid.New())
		if err != errors.ErrPlantNameTooLong {
			t.Errorf("Expected ErrPlantNameTooLong, got %v", err)
		}
	})

	t.Run("invalid latitude - below range", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST04", "Plant", -91, 0, uuid.New())
		if err != errors.ErrInvalidLatitude {
			t.Errorf("Expected ErrInvalidLatitude, got %v", err)
		}
	})

	t.Run("invalid latitude - above range", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST05", "Plant", 91, 0, uuid.New())
		if err != errors.ErrInvalidLatitude {
			t.Errorf("Expected ErrInvalidLatitude, got %v", err)
		}
	})

	t.Run("invalid longitude - below range", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST06", "Plant", 0, -181, uuid.New())
		if err != errors.ErrInvalidLongitude {
			t.Errorf("Expected ErrInvalidLongitude, got %v", err)
		}
	})

	t.Run("invalid longitude - above range", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST07", "Plant", 0, 181, uuid.New())
		if err != errors.ErrInvalidLongitude {
			t.Errorf("Expected ErrInvalidLongitude, got %v", err)
		}
	})

	t.Run("invalid owner ID - nil", func(t *testing.T) {
		_, err := service.Create(ctx, "TEST08", "Plant", 0, 0, uuid.Nil)
		if err != errors.ErrInvalidOwnerID {
			t.Errorf("Expected ErrInvalidOwnerID, got %v", err)
		}
	})
}

func TestPlantService_GetByID(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		plant, err := service.Create(ctx, "GET01", "Get Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		retrieved, err := service.GetByID(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if retrieved.ID() != plant.ID() {
			t.Errorf("Expected ID %s, got %s", plant.ID(), retrieved.ID())
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		_, err := service.GetByID(ctx, uuid.New())
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})
}

func TestPlantService_GetByCode(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		code := "CODE01"
		plant, err := service.Create(ctx, code, "Code Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		retrieved, err := service.GetByCode(ctx, code)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if retrieved.ID() != plant.ID() {
			t.Errorf("Expected ID %s, got %s", plant.ID(), retrieved.ID())
		}
	})

	t.Run("case insensitive search", func(t *testing.T) {
		code := "LOWER01"
		_, err := service.Create(ctx, code, "Lower Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		retrieved, err := service.GetByCode(ctx, "lower01")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if retrieved.Code() != code {
			t.Errorf("Expected code %s, got %s", code, retrieved.Code())
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		_, err := service.GetByCode(ctx, "NOTEXIST")
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})
}

func TestPlantService_GetByOwner(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful retrieval", func(t *testing.T) {
		ownerID := uuid.New()

		_, err := service.Create(ctx, "OWN01", "Owner Plant 1", 0, 0, ownerID)
		if err != nil {
			t.Fatalf("Failed to create plant 1: %v", err)
		}

		_, err = service.Create(ctx, "OWN02", "Owner Plant 2", 0, 0, ownerID)
		if err != nil {
			t.Fatalf("Failed to create plant 2: %v", err)
		}

		// Create plant with different owner
		_, err = service.Create(ctx, "OWN03", "Other Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant 3: %v", err)
		}

		plants, err := service.GetByOwner(ctx, ownerID)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(plants) != 2 {
			t.Errorf("Expected 2 plants, got %d", len(plants))
		}
	})

	t.Run("no plants found", func(t *testing.T) {
		plants, err := service.GetByOwner(ctx, uuid.New())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(plants) != 0 {
			t.Errorf("Expected 0 plants, got %d", len(plants))
		}
	})

	t.Run("invalid owner ID", func(t *testing.T) {
		_, err := service.GetByOwner(ctx, uuid.Nil)
		if err != errors.ErrInvalidOwnerID {
			t.Errorf("Expected ErrInvalidOwnerID, got %v", err)
		}
	})
}

func TestPlantService_ListAll(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("list active plants only", func(t *testing.T) {
		plant1, err := service.Create(ctx, "LIST01", "List Plant 1", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant 1: %v", err)
		}

		plant2, err := service.Create(ctx, "LIST02", "List Plant 2", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant 2: %v", err)
		}

		// Deactivate one plant
		err = service.Deactivate(ctx, plant1.ID())
		if err != nil {
			t.Fatalf("Failed to deactivate plant: %v", err)
		}

		plants, err := service.ListAll(ctx)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Should only return active plant
		activeCount := 0
		for _, p := range plants {
			if p.IsActive() {
				activeCount++
			}
			if p.ID() == plant2.ID() && !p.IsActive() {
				t.Error("Expected plant2 to be active")
			}
		}

		if activeCount == 0 {
			t.Error("Expected at least one active plant")
		}
	})
}

func TestPlantService_Update(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful update", func(t *testing.T) {
		plant, err := service.Create(ctx, "UPD01", "Old Name", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		newName := "New Name"
		newLatitude := -23.5505
		newLongitude := -46.6333

		updated, err := service.Update(ctx, plant.ID(), newName, newLatitude, newLongitude)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if updated.Name() != newName {
			t.Errorf("Expected name %s, got %s", newName, updated.Name())
		}

		if updated.Latitude() != newLatitude {
			t.Errorf("Expected latitude %f, got %f", newLatitude, updated.Latitude())
		}

		if updated.Longitude() != newLongitude {
			t.Errorf("Expected longitude %f, got %f", newLongitude, updated.Longitude())
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		_, err := service.Update(ctx, uuid.New(), "Name", 0, 0)
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})

	t.Run("invalid name", func(t *testing.T) {
		plant, err := service.Create(ctx, "UPD02", "Valid Name", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		_, err = service.Update(ctx, plant.ID(), "", 0, 0)
		if err != errors.ErrPlantNameRequired {
			t.Errorf("Expected ErrPlantNameRequired, got %v", err)
		}
	})

	t.Run("invalid coordinates", func(t *testing.T) {
		plant, err := service.Create(ctx, "UPD03", "Valid Name", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		_, err = service.Update(ctx, plant.ID(), "Name", 91, 0)
		if err != errors.ErrInvalidLatitude {
			t.Errorf("Expected ErrInvalidLatitude, got %v", err)
		}
	})
}

func TestPlantService_Deactivate(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful deactivation", func(t *testing.T) {
		plant, err := service.Create(ctx, "DEACT01", "Deactivate Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		err = service.Deactivate(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		retrieved, err := service.GetByID(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Failed to retrieve plant: %v", err)
		}

		if retrieved.IsActive() {
			t.Error("Expected plant to be inactive")
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		err := service.Deactivate(ctx, uuid.New())
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})
}

func TestPlantService_Activate(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful activation", func(t *testing.T) {
		plant, err := service.Create(ctx, "ACT01", "Activate Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		// Deactivate first
		err = service.Deactivate(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Failed to deactivate plant: %v", err)
		}

		// Now activate
		err = service.Activate(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		retrieved, err := service.GetByID(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Failed to retrieve plant: %v", err)
		}

		if !retrieved.IsActive() {
			t.Error("Expected plant to be active")
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		err := service.Activate(ctx, uuid.New())
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})
}

func TestPlantService_Delete(t *testing.T) {
	repository := persistence.NewInMemoryPlantRepository()
	service := NewPlantService(repository)
	ctx := context.Background()

	t.Run("successful deletion", func(t *testing.T) {
		plant, err := service.Create(ctx, "DEL01", "Delete Plant", 0, 0, uuid.New())
		if err != nil {
			t.Fatalf("Failed to create plant: %v", err)
		}

		err = service.Delete(ctx, plant.ID())
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		_, err = service.GetByID(ctx, plant.ID())
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected plant to be deleted, got %v", err)
		}
	})

	t.Run("plant not found", func(t *testing.T) {
		err := service.Delete(ctx, uuid.New())
		if err != errors.ErrPlantNotFound {
			t.Errorf("Expected ErrPlantNotFound, got %v", err)
		}
	})
}
