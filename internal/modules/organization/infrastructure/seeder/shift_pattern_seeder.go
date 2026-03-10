package seeder

import (
	"fmt"
	"time"

	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/organization/domain/entity"
)

const defaultPatternName = "Default 3-Shift Rotation"

// defaultEntries defines the 3-day rotating cycle: Morning → Afternoon → Night.
var defaultEntries = []struct {
	name      string
	startTime string
	endTime   string
}{
	{"Morning", "06:00", "14:00"},
	{"Afternoon", "14:00", "22:00"},
	{"Night", "22:00", "06:00"},
}

// SeedDefaultShiftPatterns seeds the default shift pattern if it does not already exist.
// The default pattern is a 3-day rotating cycle anchored to 2024-01-01.
// It is safe to call on every boot.
func SeedDefaultShiftPatterns(repo output.ShiftPatternRepository) error {
	exists, err := repo.ExistsByName(defaultPatternName)
	if err != nil {
		return fmt.Errorf("shift pattern seeder: check existence: %w", err)
	}

	if exists {
		return nil
	}

	// Anchor to a fixed reference date for a deterministic cycle.
	refStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	pattern := entity.NewShiftPattern(defaultPatternName, refStart)

	for i, e := range defaultEntries {
		pattern.AddEntry(entity.NewShiftEntry(i, e.name, e.startTime, e.endTime))
	}

	if err := repo.Save(pattern); err != nil {
		return fmt.Errorf("shift pattern seeder: save: %w", err)
	}

	return nil
}
