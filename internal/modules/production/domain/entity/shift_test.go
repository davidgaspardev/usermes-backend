package entity

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNewShift(t *testing.T) {
	patternID := uuid.New()
	startAt := time.Date(2026, 3, 8, 22, 0, 0, 0, time.UTC)
	endAt := time.Date(2026, 3, 9, 6, 0, 0, 0, time.UTC)

	s := NewShift(patternID, "Night", startAt, endAt)

	if s == nil {
		t.Fatal("expected shift, got nil")
	}
	if s.ID() == (uuid.UUID{}) {
		t.Error("expected non-zero ID")
	}
	if s.PatternID() != patternID {
		t.Errorf("expected patternID %v, got %v", patternID, s.PatternID())
	}
	if s.Name() != "Night" {
		t.Errorf("expected name Night, got %s", s.Name())
	}
	if !s.StartAt().Equal(startAt) {
		t.Errorf("expected startAt %v, got %v", startAt, s.StartAt())
	}
	if !s.EndAt().Equal(endAt) {
		t.Errorf("expected endAt %v, got %v", endAt, s.EndAt())
	}
	if s.CreatedAt().IsZero() {
		t.Error("expected non-zero createdAt")
	}
}
