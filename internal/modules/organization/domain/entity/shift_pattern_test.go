package entity

import (
	"testing"
	"time"
)

func mustParseTime(s string) time.Time {
	t, err := time.Parse("15:04", s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNewShiftEntry(t *testing.T) {
	e := NewShiftEntry(0, "Morning", mustParseTime("06:00"), mustParseTime("14:00"))

	if e.DayIndex() != 0 {
		t.Errorf("expected dayIndex 0, got %d", e.DayIndex())
	}
	if e.Name() != "Morning" {
		t.Errorf("expected name Morning, got %s", e.Name())
	}
	if e.StartTime().Hour() != 6 || e.StartTime().Minute() != 0 {
		t.Errorf("expected startTime 06:00, got %v", e.StartTime())
	}
	if e.EndTime().Hour() != 14 || e.EndTime().Minute() != 0 {
		t.Errorf("expected endTime 14:00, got %v", e.EndTime())
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
	if !e.StartTime().IsZero() {
		t.Errorf("expected zero startTime, got %v", e.StartTime())
	}
	if !e.EndTime().IsZero() {
		t.Errorf("expected zero endTime, got %v", e.EndTime())
	}
}

func TestShiftEntry_TimesForDate_Normal(t *testing.T) {
	e := NewShiftEntry(0, "Morning", mustParseTime("06:00"), mustParseTime("14:00"))
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
	e := NewShiftEntry(2, "Night", mustParseTime("22:00"), mustParseTime("06:00"))
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
	p := NewShiftPattern("Weekly", refStart, 7)

	if p == nil {
		t.Fatal("expected pattern, got nil")
	}
	if p.Name() != "Weekly" {
		t.Errorf("expected name Weekly, got %s", p.Name())
	}
	if !p.RefStartDate().Equal(refStart) {
		t.Errorf("expected refStartDate %v, got %v", refStart, p.RefStartDate())
	}
	if p.CycleLength() != 7 {
		t.Errorf("expected cycleLength 7, got %d", p.CycleLength())
	}
	if len(p.Entries()) != 0 {
		t.Errorf("expected 0 entries, got %d", len(p.Entries()))
	}
	if p.ID() == (p.ID()) && p.CreatedAt().IsZero() {
		t.Error("expected non-zero createdAt")
	}
}

func TestShiftPattern_AddEntry(t *testing.T) {
	p := NewShiftPattern("3-shift", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 3)
	p.AddEntry(NewShiftEntry(0, "Morning", mustParseTime("06:00"), mustParseTime("14:00")))
	p.AddEntry(NewShiftEntry(1, "Afternoon", mustParseTime("14:00"), mustParseTime("22:00")))
	p.AddEntry(NewShiftEntry(2, "Night", mustParseTime("22:00"), mustParseTime("06:00")))

	if p.CycleLength() != 3 {
		t.Errorf("expected cycleLength 3, got %d", p.CycleLength())
	}
	if len(p.Entries()) != 3 {
		t.Errorf("expected 3 entries slice, got %d", len(p.Entries()))
	}
}

func TestShiftPattern_EntriesForDate_NoEntries(t *testing.T) {
	p := NewShiftPattern("empty", time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 3)

	if p.EntriesForDate(time.Now()) != nil {
		t.Error("expected nil for pattern with no entries")
	}
}

func TestShiftPattern_EntriesForDate_SinglePerDay(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday
	p := NewShiftPattern("3-shift", ref, 3)
	p.AddEntry(NewShiftEntry(0, "Morning", mustParseTime("06:00"), mustParseTime("14:00")))
	p.AddEntry(NewShiftEntry(1, "Afternoon", mustParseTime("14:00"), mustParseTime("22:00")))
	p.AddEntry(NewShiftEntry(2, "Night", mustParseTime("22:00"), mustParseTime("06:00")))

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
		entries := p.EntriesForDate(tt.date)
		if len(entries) != 1 {
			t.Fatalf("date %v: expected 1 entry, got %d", tt.date, len(entries))
		}
		if entries[0].Name() != tt.expectedName {
			t.Errorf("date %v: expected %s, got %s", tt.date, tt.expectedName, entries[0].Name())
		}
	}
}

func TestShiftPattern_EntriesForDate_MultiplePerDay(t *testing.T) {
	ref := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC) // Monday = day 0
	p := NewShiftPattern("Weekly", ref, 7)

	// Weekdays (0-4): 3 shifts each
	for day := 0; day < 5; day++ {
		p.AddEntry(NewShiftEntry(day, "Shift 1", mustParseTime("06:00"), mustParseTime("14:00")))
		p.AddEntry(NewShiftEntry(day, "Shift 2", mustParseTime("14:00"), mustParseTime("22:00")))
		p.AddEntry(NewShiftEntry(day, "Shift 3", mustParseTime("22:00"), mustParseTime("06:00")))
	}
	// Weekend (5-6): 2 shifts each
	for day := 5; day < 7; day++ {
		p.AddEntry(NewShiftEntry(day, "Shift 1", mustParseTime("06:00"), mustParseTime("14:00")))
		p.AddEntry(NewShiftEntry(day, "Shift 2", mustParseTime("14:00"), mustParseTime("22:00")))
	}

	monday := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	if entries := p.EntriesForDate(monday); len(entries) != 3 {
		t.Errorf("expected 3 entries for Monday, got %d", len(entries))
	}

	saturday := time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC)
	if entries := p.EntriesForDate(saturday); len(entries) != 2 {
		t.Errorf("expected 2 entries for Saturday, got %d", len(entries))
	}

	sunday := time.Date(2024, 1, 7, 0, 0, 0, 0, time.UTC)
	if entries := p.EntriesForDate(sunday); len(entries) != 2 {
		t.Errorf("expected 2 entries for Sunday, got %d", len(entries))
	}

	nextMonday := time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC)
	if entries := p.EntriesForDate(nextMonday); len(entries) != 3 {
		t.Errorf("expected 3 entries for next Monday, got %d", len(entries))
	}
}

func TestShiftPattern_EntriesForDate_BeforeRef(t *testing.T) {
	ref := time.Date(2024, 1, 4, 0, 0, 0, 0, time.UTC)
	p := NewShiftPattern("3-shift", ref, 3)
	p.AddEntry(NewShiftEntry(0, "Morning", mustParseTime("06:00"), mustParseTime("14:00")))
	p.AddEntry(NewShiftEntry(1, "Afternoon", mustParseTime("14:00"), mustParseTime("22:00")))
	p.AddEntry(NewShiftEntry(2, "Night", mustParseTime("22:00"), mustParseTime("06:00")))

	// 1 day before ref → offset = 3 - (1 % 3) = 2 → Night
	date := time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC)
	entries := p.EntriesForDate(date)
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Name() != "Night" {
		t.Errorf("expected Night for date before ref, got %s", entries[0].Name())
	}
}
