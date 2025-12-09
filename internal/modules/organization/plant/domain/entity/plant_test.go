package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewPlant(t *testing.T) {
	code := "SP01"
	name := "São Paulo Plant"
	latitude := -23.5505
	longitude := -46.6333
	ownerID := uuid.New()

	plant := NewPlant(code, name, latitude, longitude, ownerID)

	if plant == nil {
		t.Fatal("Expected plant to be created, got nil")
	}

	if plant.ID() == uuid.Nil {
		t.Error("Expected plant ID to be generated")
	}

	if plant.Code() != code {
		t.Errorf("Expected code %s, got %s", code, plant.Code())
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
		t.Error("Expected new plant to be active")
	}

	if plant.CreatedAt().IsZero() {
		t.Error("Expected createdAt to be set")
	}

	if plant.UpdatedAt().IsZero() {
		t.Error("Expected updatedAt to be set")
	}

	if plant.CreatedAt().After(time.Now()) {
		t.Error("Expected createdAt to be in the past")
	}
}

func TestReconstructPlant(t *testing.T) {
	id := uuid.New()
	code := "RJ01"
	name := "Rio de Janeiro Plant"
	latitude := -22.9068
	longitude := -43.1729
	ownerID := uuid.New()
	isActive := true
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	plant := ReconstructPlant(id, code, name, latitude, longitude, ownerID, isActive, createdAt, updatedAt)

	if plant == nil {
		t.Fatal("Expected plant to be reconstructed, got nil")
	}

	if plant.ID() != id {
		t.Errorf("Expected ID %s, got %s", id, plant.ID())
	}

	if plant.Code() != code {
		t.Errorf("Expected code %s, got %s", code, plant.Code())
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

	if plant.IsActive() != isActive {
		t.Errorf("Expected isActive %t, got %t", isActive, plant.IsActive())
	}

	if !plant.CreatedAt().Equal(createdAt) {
		t.Errorf("Expected createdAt %v, got %v", createdAt, plant.CreatedAt())
	}

	if !plant.UpdatedAt().Equal(updatedAt) {
		t.Errorf("Expected updatedAt %v, got %v", updatedAt, plant.UpdatedAt())
	}
}

func TestPlant_UpdateName(t *testing.T) {
	plant := NewPlant("SP01", "Old Name", -23.5505, -46.6333, uuid.New())
	oldUpdatedAt := plant.UpdatedAt()

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	newName := "New Name"
	plant.UpdateName(newName)

	if plant.Name() != newName {
		t.Errorf("Expected name %s, got %s", newName, plant.Name())
	}

	if !plant.UpdatedAt().After(oldUpdatedAt) {
		t.Error("Expected updatedAt to be updated")
	}
}

func TestPlant_UpdateLocation(t *testing.T) {
	plant := NewPlant("SP01", "São Paulo Plant", -23.5505, -46.6333, uuid.New())
	oldUpdatedAt := plant.UpdatedAt()

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	newLatitude := -22.9068
	newLongitude := -43.1729
	plant.UpdateLocation(newLatitude, newLongitude)

	if plant.Latitude() != newLatitude {
		t.Errorf("Expected latitude %f, got %f", newLatitude, plant.Latitude())
	}

	if plant.Longitude() != newLongitude {
		t.Errorf("Expected longitude %f, got %f", newLongitude, plant.Longitude())
	}

	if !plant.UpdatedAt().After(oldUpdatedAt) {
		t.Error("Expected updatedAt to be updated")
	}
}

func TestPlant_Update(t *testing.T) {
	plant := NewPlant("SP01", "Old Name", -23.5505, -46.6333, uuid.New())
	oldUpdatedAt := plant.UpdatedAt()

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	newName := "New Name"
	newLatitude := -22.9068
	newLongitude := -43.1729
	plant.Update(newName, newLatitude, newLongitude)

	if plant.Name() != newName {
		t.Errorf("Expected name %s, got %s", newName, plant.Name())
	}

	if plant.Latitude() != newLatitude {
		t.Errorf("Expected latitude %f, got %f", newLatitude, plant.Latitude())
	}

	if plant.Longitude() != newLongitude {
		t.Errorf("Expected longitude %f, got %f", newLongitude, plant.Longitude())
	}

	if !plant.UpdatedAt().After(oldUpdatedAt) {
		t.Error("Expected updatedAt to be updated")
	}
}

func TestPlant_Deactivate(t *testing.T) {
	plant := NewPlant("SP01", "São Paulo Plant", -23.5505, -46.6333, uuid.New())

	if !plant.IsActive() {
		t.Fatal("Expected plant to be active initially")
	}

	oldUpdatedAt := plant.UpdatedAt()

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	plant.Deactivate()

	if plant.IsActive() {
		t.Error("Expected plant to be inactive after deactivation")
	}

	if !plant.UpdatedAt().After(oldUpdatedAt) {
		t.Error("Expected updatedAt to be updated")
	}
}

func TestPlant_Activate(t *testing.T) {
	plant := NewPlant("SP01", "São Paulo Plant", -23.5505, -46.6333, uuid.New())
	plant.Deactivate()

	if plant.IsActive() {
		t.Fatal("Expected plant to be inactive")
	}

	oldUpdatedAt := plant.UpdatedAt()

	// Wait a bit to ensure time difference
	time.Sleep(10 * time.Millisecond)

	plant.Activate()

	if !plant.IsActive() {
		t.Error("Expected plant to be active after activation")
	}

	if !plant.UpdatedAt().After(oldUpdatedAt) {
		t.Error("Expected updatedAt to be updated")
	}
}

func TestPlant_Getters(t *testing.T) {
	id := uuid.New()
	code := "SP01"
	name := "São Paulo Plant"
	latitude := -23.5505
	longitude := -46.6333
	ownerID := uuid.New()
	isActive := true
	createdAt := time.Now().Add(-24 * time.Hour)
	updatedAt := time.Now()

	plant := ReconstructPlant(id, code, name, latitude, longitude, ownerID, isActive, createdAt, updatedAt)

	// Test all getters
	t.Run("ID getter", func(t *testing.T) {
		if plant.ID() != id {
			t.Errorf("Expected ID %s, got %s", id, plant.ID())
		}
	})

	t.Run("Code getter", func(t *testing.T) {
		if plant.Code() != code {
			t.Errorf("Expected code %s, got %s", code, plant.Code())
		}
	})

	t.Run("Name getter", func(t *testing.T) {
		if plant.Name() != name {
			t.Errorf("Expected name %s, got %s", name, plant.Name())
		}
	})

	t.Run("Latitude getter", func(t *testing.T) {
		if plant.Latitude() != latitude {
			t.Errorf("Expected latitude %f, got %f", latitude, plant.Latitude())
		}
	})

	t.Run("Longitude getter", func(t *testing.T) {
		if plant.Longitude() != longitude {
			t.Errorf("Expected longitude %f, got %f", longitude, plant.Longitude())
		}
	})

	t.Run("OwnerID getter", func(t *testing.T) {
		if plant.OwnerID() != ownerID {
			t.Errorf("Expected ownerID %s, got %s", ownerID, plant.OwnerID())
		}
	})

	t.Run("IsActive getter", func(t *testing.T) {
		if plant.IsActive() != isActive {
			t.Errorf("Expected isActive %t, got %t", isActive, plant.IsActive())
		}
	})

	t.Run("CreatedAt getter", func(t *testing.T) {
		if !plant.CreatedAt().Equal(createdAt) {
			t.Errorf("Expected createdAt %v, got %v", createdAt, plant.CreatedAt())
		}
	})

	t.Run("UpdatedAt getter", func(t *testing.T) {
		if !plant.UpdatedAt().Equal(updatedAt) {
			t.Errorf("Expected updatedAt %v, got %v", updatedAt, plant.UpdatedAt())
		}
	})
}
