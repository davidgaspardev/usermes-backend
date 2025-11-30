package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewResource(t *testing.T) {
	code := "RES001"
	shiftID := "SHIFT123"
	resourceType := "MACHINE"
	stopFactor := int16(5)

	resource := NewResource(code, &shiftID, resourceType, stopFactor)

	if resource.ID() == uuid.Nil {
		t.Error("Resource ID should not be nil")
	}

	if resource.Code() != code {
		t.Errorf("Expected code '%s', got '%s'", code, resource.Code())
	}

	if resource.ShiftID() == nil {
		t.Error("ShiftID should not be nil")
	} else if *resource.ShiftID() != shiftID {
		t.Errorf("Expected shiftID '%s', got '%s'", shiftID, *resource.ShiftID())
	}

	if resource.Type() != resourceType {
		t.Errorf("Expected type '%s', got '%s'", resourceType, resource.Type())
	}

	if resource.StopFactor() != stopFactor {
		t.Errorf("Expected stopFactor '%d', got '%d'", stopFactor, resource.StopFactor())
	}

	if resource.CreatedAt().IsZero() {
		t.Error("CreatedAt should be set")
	}

	if resource.UpdatedAt().IsZero() {
		t.Error("UpdatedAt should be set")
	}
}

func TestNewResource_WithNilShiftID(t *testing.T) {
	code := "RES002"
	resourceType := "OPERATOR"
	stopFactor := int16(0)

	resource := NewResource(code, nil, resourceType, stopFactor)

	if resource.ShiftID() != nil {
		t.Error("ShiftID should be nil")
	}
}

func TestResource_UpdateCode(t *testing.T) {
	resource := NewResource("RES001", nil, "MACHINE", 5)
	originalUpdatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	newCode := "RES002"
	resource.UpdateCode(newCode)

	if resource.Code() != newCode {
		t.Errorf("Expected code '%s', got '%s'", newCode, resource.Code())
	}

	if !resource.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestResource_UpdateShiftID(t *testing.T) {
	resource := NewResource("RES001", nil, "MACHINE", 5)
	originalUpdatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	newShiftID := "SHIFT456"
	resource.UpdateShiftID(&newShiftID)

	if resource.ShiftID() == nil {
		t.Error("ShiftID should not be nil after update")
	} else if *resource.ShiftID() != newShiftID {
		t.Errorf("Expected shiftID '%s', got '%s'", newShiftID, *resource.ShiftID())
	}

	if !resource.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestResource_UpdateShiftID_ToNil(t *testing.T) {
	shiftID := "SHIFT123"
	resource := NewResource("RES001", &shiftID, "MACHINE", 5)

	resource.UpdateShiftID(nil)

	if resource.ShiftID() != nil {
		t.Error("ShiftID should be nil after update")
	}
}

func TestResource_UpdateType(t *testing.T) {
	resource := NewResource("RES001", nil, "MACHINE", 5)
	originalUpdatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	newType := "OPERATOR"
	resource.UpdateType(newType)

	if resource.Type() != newType {
		t.Errorf("Expected type '%s', got '%s'", newType, resource.Type())
	}

	if !resource.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestResource_UpdateStopFactor(t *testing.T) {
	resource := NewResource("RES001", nil, "MACHINE", 5)
	originalUpdatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	newStopFactor := int16(10)
	resource.UpdateStopFactor(newStopFactor)

	if resource.StopFactor() != newStopFactor {
		t.Errorf("Expected stopFactor '%d', got '%d'", newStopFactor, resource.StopFactor())
	}

	if !resource.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestResource_Update(t *testing.T) {
	resource := NewResource("RES001", nil, "MACHINE", 5)
	originalUpdatedAt := resource.UpdatedAt()

	time.Sleep(10 * time.Millisecond)
	newCode := "RES999"
	newShiftID := "SHIFT999"
	newType := "TOOL"
	newStopFactor := int16(15)

	resource.Update(newCode, &newShiftID, newType, newStopFactor)

	if resource.Code() != newCode {
		t.Errorf("Expected code '%s', got '%s'", newCode, resource.Code())
	}

	if resource.ShiftID() == nil || *resource.ShiftID() != newShiftID {
		t.Error("ShiftID should be updated")
	}

	if resource.Type() != newType {
		t.Errorf("Expected type '%s', got '%s'", newType, resource.Type())
	}

	if resource.StopFactor() != newStopFactor {
		t.Errorf("Expected stopFactor '%d', got '%d'", newStopFactor, resource.StopFactor())
	}

	if !resource.UpdatedAt().After(originalUpdatedAt) {
		t.Error("UpdatedAt should be updated")
	}
}

func TestReconstructResource(t *testing.T) {
	id := uuid.New()
	code := "RES001"
	shiftID := "SHIFT123"
	resourceType := "MACHINE"
	stopFactor := int16(5)
	now := time.Now()

	resource := ReconstructResource(
		id,
		code,
		&shiftID,
		resourceType,
		stopFactor,
		now,
		now,
	)

	if resource.ID() != id {
		t.Error("Reconstructed resource should have correct ID")
	}

	if resource.Code() != code {
		t.Error("Reconstructed resource should have correct code")
	}

	if resource.ShiftID() == nil || *resource.ShiftID() != shiftID {
		t.Error("Reconstructed resource should have correct shiftID")
	}

	if resource.Type() != resourceType {
		t.Error("Reconstructed resource should have correct type")
	}

	if resource.StopFactor() != stopFactor {
		t.Error("Reconstructed resource should have correct stopFactor")
	}

	if !resource.CreatedAt().Equal(now) {
		t.Error("Reconstructed resource should have correct createdAt")
	}

	if !resource.UpdatedAt().Equal(now) {
		t.Error("Reconstructed resource should have correct updatedAt")
	}
}

func TestResource_AllGetters(t *testing.T) {
	shiftID := "SHIFT123"
	resource := NewResource("RES001", &shiftID, "MACHINE", 5)

	// Test all getter methods
	if resource.ID() == uuid.Nil {
		t.Error("ID should not be nil")
	}
	if resource.Code() != "RES001" {
		t.Error("Code getter failed")
	}
	if resource.ShiftID() == nil || *resource.ShiftID() != shiftID {
		t.Error("ShiftID getter failed")
	}
	if resource.Type() != "MACHINE" {
		t.Error("Type getter failed")
	}
	if resource.StopFactor() != 5 {
		t.Error("StopFactor getter failed")
	}
	if resource.CreatedAt().IsZero() {
		t.Error("CreatedAt getter failed")
	}
	if resource.UpdatedAt().IsZero() {
		t.Error("UpdatedAt getter failed")
	}
}
