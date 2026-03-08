package seeder

import (
	"fmt"

	"github.com/davidgaspardev/usermes-backend/internal/modules/production/application/port/output"
	"github.com/davidgaspardev/usermes-backend/internal/modules/production/domain/entity"
)

var defaultShiftTemplates = []struct {
	name      string
	startTime string
	endTime   string
}{
	{"Morning", "06:00", "14:00"},
	{"Afternoon", "14:00", "22:00"},
	{"Night", "22:00", "06:00"},
}

// SeedDefaultShiftTemplates seeds the default shift templates (Morning / Afternoon / Night)
// if they do not already exist. It is safe to call on every boot.
func SeedDefaultShiftTemplates(repo output.ShiftTemplateRepository) error {
	for _, t := range defaultShiftTemplates {
		exists, err := repo.ExistsByName(t.name)
		if err != nil {
			return fmt.Errorf("shift template seeder: check existence of %q: %w", t.name, err)
		}
		if exists {
			continue
		}

		template := entity.NewShiftTemplate(t.name, t.startTime, t.endTime)
		if err := repo.Save(template); err != nil {
			return fmt.Errorf("shift template seeder: save %q: %w", t.name, err)
		}
	}
	return nil
}
