package entity

import (
	"testing"

	"github.com/google/uuid"
)

func TestNewRootLocation(t *testing.T) {
	loc := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})

	if loc == nil {
		t.Fatal("Expected location, got nil")
	}
	if loc.Code() != "PLANT01" {
		t.Errorf("Expected code PLANT01, got %s", loc.Code())
	}
	if loc.Name() != "Plant One" {
		t.Errorf("Expected name Plant One, got %s", loc.Name())
	}
	if loc.Kind() != LocationKindPlant {
		t.Errorf("Expected kind PLANT, got %s", loc.Kind())
	}
	if loc.ParentCode() != "" {
		t.Errorf("Expected empty parent code, got %s", loc.ParentCode())
	}
	if len(loc.Children()) != 0 {
		t.Errorf("Expected no children, got %d", len(loc.Children()))
	}
}

func TestNewLocation(t *testing.T) {
	root := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})
	child := NewLocation("AREA01", "Area One", LocationKindArea, root, uuid.UUID{})

	if child.Code() != "AREA01" {
		t.Errorf("Expected code AREA01, got %s", child.Code())
	}
	if child.Kind() != LocationKindArea {
		t.Errorf("Expected kind AREA, got %s", child.Kind())
	}
	if child.ParentCode() != "PLANT01" {
		t.Errorf("Expected parent code PLANT01, got %s", child.ParentCode())
	}
}

func TestNewLocation_WithShiftPattern(t *testing.T) {
	patternID := uuid.New()
	loc := NewRootLocation("PLANT01", "Plant One", patternID)

	if loc.ShiftPatternID() != patternID {
		t.Errorf("Expected shiftPatternID %v, got %v", patternID, loc.ShiftPatternID())
	}
}

func TestLocation_AddChild(t *testing.T) {
	root := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})
	child := NewLocation("AREA01", "Area One", LocationKindArea, root, uuid.UUID{})
	root.AddChild(child)

	if len(root.Children()) != 1 {
		t.Fatalf("Expected 1 child, got %d", len(root.Children()))
	}
	if root.Children()[0].Code() != "AREA01" {
		t.Errorf("Expected child code AREA01, got %s", root.Children()[0].Code())
	}
}

func TestLocation_FindByCode(t *testing.T) {
	root := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})
	area := NewLocation("AREA01", "Area One", LocationKindArea, root, uuid.UUID{})
	line := NewLocation("LINE01", "Line One", LocationKindLine, area, uuid.UUID{})
	root.AddChild(area)
	area.AddChild(line)

	t.Run("find root", func(t *testing.T) {
		found := root.FindByCode("PLANT01")
		if found == nil || found.Code() != "PLANT01" {
			t.Error("Expected to find root")
		}
	})

	t.Run("find direct child", func(t *testing.T) {
		found := root.FindByCode("AREA01")
		if found == nil || found.Code() != "AREA01" {
			t.Error("Expected to find AREA01")
		}
	})

	t.Run("find nested child", func(t *testing.T) {
		found := root.FindByCode("LINE01")
		if found == nil || found.Code() != "LINE01" {
			t.Error("Expected to find LINE01")
		}
	})

	t.Run("not found", func(t *testing.T) {
		found := root.FindByCode("UNKNOWN")
		if found != nil {
			t.Error("Expected nil for unknown code")
		}
	})
}

func TestLocation_ExistsByCode(t *testing.T) {
	root := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})
	area := NewLocation("AREA01", "Area One", LocationKindArea, root, uuid.UUID{})
	root.AddChild(area)

	if !root.ExistsByCode("PLANT01") {
		t.Error("Expected PLANT01 to exist")
	}
	if !root.ExistsByCode("AREA01") {
		t.Error("Expected AREA01 to exist")
	}
	if root.ExistsByCode("UNKNOWN") {
		t.Error("Expected UNKNOWN to not exist")
	}
}

func TestIsValidLocationKind(t *testing.T) {
	tests := []struct {
		kind  string
		valid bool
	}{
		{"PLANT", true},
		{"AREA", true},
		{"LINE", true},
		{"SECTION", true},
		{"INVALID", false},
		{"", false},
		{"plant", false},
	}

	for _, tt := range tests {
		got := IsValidLocationKind(tt.kind)
		if got != tt.valid {
			t.Errorf("IsValidLocationKind(%q) = %v, want %v", tt.kind, got, tt.valid)
		}
	}
}

func TestLocation_FindByCode_NoChildren(t *testing.T) {
	root := NewRootLocation("PLANT01", "Plant One", uuid.UUID{})
	found := root.FindByCode("AREA01")
	if found != nil {
		t.Error("Expected nil when no children exist")
	}
}
