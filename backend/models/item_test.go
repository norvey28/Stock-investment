package models

import (
	"encoding/json"
	"testing"
)

func TestMoneyUnmarshal_StringWithDollar(t *testing.T) {
	var m Money
	if err := json.Unmarshal([]byte(`"$42.00"`), &m); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if float64(m) != 42.0 {
		t.Fatalf("expected 42.0, got %v", m)
	}
}

func TestMoneyUnmarshal_Number(t *testing.T) {
	var m Money
	if err := json.Unmarshal([]byte(`42.5`), &m); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if float64(m) != 42.5 {
		t.Fatalf("expected 42.5, got %v", m)
	}
}
