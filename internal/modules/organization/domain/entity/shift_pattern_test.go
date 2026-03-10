package entity

import (
	"testing"
	"time"
)

func TestNewShiftEntry(t *testing.T) {
	e := NewShiftEntry(0, "Morning", "06:00", "14:00")

	if e.DayIndex() != 0 {
		t.Errorf("expected dayIndex 0, got %d", e.DayIndex())
	}
	if e.Name() != "Morning" {
		t.Errorf("expected name Morning, got %s", e.Name())
	}
	if e.StartTime() != "06:00" {
		t.Errorf("expected startTime 06:00, got %s", e.StartTime())
	}
	if e.EndTime() != "14:00" {
		t.Errorf("expected endTime 14:00, got %s", e.EndTime())
	}
	if e.IsOff() {
		t.Error("expected IsOff false")
	}
}

func TestNewDayOffEntry(t *testing.T) {
	e := NewDayOffEntry(6)

	if e.DayIndex() != 6 {
		t.Errorf("expected dayIndex 6, got %d", e.DayIndex())
	}
	if e.Name() != "Off" {
		t.Errorf("expected name Off, got %s", e.Name())
	}
	if !e.IsOff() {
		t.Error("expected IsOff true")
	}
	if e.StartTime() != "" {
		t.Errorf("expected empty startTime, got %s", e.StartTime())
	}
	if e.EndTime() != "" {
		t.Errorf("expected empty endTime, got %s", e.EndTime())
	}
}

func TestShiftEntry_TimesForDate_Normal(t *testing.T) {
	e := NewShiftEntry(0, "Morning", "06:00", "14:00")
	date := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

	startAt, endAt, err := e.TimesForDate(date)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startAt.Hour() != 6 || startAt.Minute() != 0 {
		t.Errorf("expected startAt 06:00, got %v", startAt)
	}
	if endAt.Hour() != 14 || endAt.Minute() != 0 {
		t.Errorf("expected endAt 14:00, got %v", endAt)
	}
	if endAt.Day() != startAt.Day() {
		t.Error("expected same day for non-overnight shift")
	}
}

func TestShiftEntry_TimesForDate_Overnight(t *testing.T) {
	e := NewShiftEntry(2, "Night", "22:00", "06:00")
	date := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

	startAt, endAt, err := e.TimesForDate(date)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if startAt.Hour() != 22 {
		t.Errorf("expected startAt hour 22, got %d", startAt.Hour())
	}
	if endAt.Hour() != 6 {
		t.Errorf("expected endAt hour 6, got %d", endAt.Hour())
	}
	if endAt.Day() != startAt.Day()+1 {
		t.Errorf("expected endAt next day for overnight shift")
	}
}

func TestShiftEntry_TimesForDate_OffDay(t *testing.T) {
	e := NewDayOffEntry(6)
	date := time.Date(2026, 3, 8, 0, 0, 0, 0, time.UTC)

	_, _, err := e.TimesForDate(date)

	if err == nil {
		t.Error("expected error for off-day entry, got nil")
	}
}

func TestNewShiftPattern(t *testing.T) {
	refStart := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	p := NewShiftPattern("Weekly", refStart)

	if p == nil {
		t.Fatal("expected pattern, got nil")
	}
	if p.Name() != "Weekly" {
		t.Errorf("expected name Weekly, got %s", p.Name())
	}
	if !p.RefStartDate().Equal(refStart) {
		t.Errorf("expected refStartDate %v, got %v", refStart, p.RefStartDate())
	}
	if p.PeriodDays() != 0 {
		t.Errorf("expected 0 entries, got %d", p.PeriodDays())
	}
	if p.ID() == (p.ID()) && p.CreatedAt().IsZero() {
		t.Error("expected non-zero createdAt")
	}
}

func TestShiftPattern_AddEntry(t *testing.T) {
	p := NewShiftPattern("3-shift", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))
	p.AddEntry(NewShiftEntry(0, "Morning", "06:00", "14:00"))
	p.AddEntry(NewShiftEntry(1, "Afternoon", "14:00", "22:00"))
	p.AddEntry(NewShiftEntry(2, "Night", "22:00", "06:00"))

	if p.PeriodDays() != 3 {
		t.Errorf("expected 3 entries, got %d", p.PeriodDays())
	}
	if len(p.Entries()) != 3 {
		t.Errorf("expected 3 entries slice, got %d", len(p.Entries()))
	}
}

func TestShiftPattern_EntryForDate_NoEntries(t *testing.T) {
	p := NewShiftPattern("empty", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC))

	if p.EntryForDate(time.Now()) != nil {
		t.Error("expected nil for pattern with no entries")
	}
}

func TestShiftPattern_EntryForDate_Rotation(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	p := NewShiftPattern("3-shift", ref)
	p.AddEntry(NewShiftEntry(0, "Morning", "06:00", "14:00"))
	p.AddEntry(NewShiftEntry(1, "Afternoon", "14:00", "22:00"))
	p.AddEntry(NewShiftEntry(2, "Night", "22:00", "06:00"))

	tests := []struct {
		date         time.Time
		expectedName string
	}{
		{time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), "Morning"},   // day 0
		{time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), "Afternoon"}, // day 1
		{time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), "Night"},     // day 2
		{time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC), "Morning"},   // day 3 → wraps
		{time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC), "Morning"},   // day 6 → wraps
	}

	for _, tt := range tests {
		entry := p.EntryForDate(tt.date)
		if entry == nil {
			t.Fatalf("expected entry for %v, got nil", tt.date)
		}
		if entry.Name() != tt.expectedName {
			t.Errorf("date %v: expected %s, got %s", tt.date, tt.expectedName, entry.Name())
		}
	}
}

func TestShiftPattern_EntryForDate_BeforeRef(t *testing.T) {
	ref := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC) // day 0 = Afternoon
	p := NewShiftPattern("3-shift", ref)
	p.AddEntry(NewShiftEntry(0, "Morning", "06:00", "14:00"))
	p.AddEntry(NewShiftEntry(1, "Afternoon", "14:00", "22:00"))
	p.AddEntry(NewShiftEntry(2, "Night", "22:00", "06:00"))

	// 1 day before ref → index = 3 - (1 % 3) = 2 → Night
	date := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	entry := p.EntryForDate(date)
	if entry == nil {
		t.Fatal("expected entry, got nil")
	}
	if entry.Name() != "Night" {
		t.Errorf("expected Night for date before ref, got %s", entry.Name())
	}
}
