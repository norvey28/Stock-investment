package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestMoneyMarshal(t *testing.T) {
	m := Money(123.45)
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}
	if string(b) != "123.45" && string(b) != "123.45000000000001" {
		t.Fatalf("unexpected marshal output: %s", string(b))
	}
}

func TestMoneyValueAndScan(t *testing.T) {
	m := Money(77.77)
	v, err := m.Value()
	if err != nil {
		t.Fatalf("value error: %v", err)
	}
	if _, ok := v.(float64); !ok {
		t.Fatalf("expected float64 value, got %T", v)
	}

	// Scan from float64
	var m2 Money
	if err := m2.Scan(v); err != nil {
		t.Fatalf("scan error: %v", err)
	}
	if float64(m2) != 77.77 {
		t.Fatalf("expected 77.77, got %v", m2)
	}

	// Scan from []byte
	var m3 Money
	if err := m3.Scan([]byte("88.88")); err != nil {
		t.Fatalf("scan error bytes: %v", err)
	}
	if float64(m3) != 88.88 {
		t.Fatalf("expected 88.88, got %v", m3)
	}

	// Scan from string
	var m4 Money
	if err := m4.Scan("99.99"); err != nil {
		t.Fatalf("scan error string: %v", err)
	}
	if float64(m4) != 99.99 {
		t.Fatalf("expected 99.99, got %v", m4)
	}
}

func TestItemJSONRoundtrip(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	item := Item{
		ID:         uuid.New(),
		Ticker:     "BOX",
		TargetFrom: Money(42.0),
		TargetTo:   Money(36.0),
		Company:    "BOX",
		Action:     "initiated by",
		Brokerage:  "",
		RatingFrom: "Buy",
		RatingTo:   "Neutral",
		Time:       now,
		CreatedAt:  now,
	}

	b, err := json.Marshal(item)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var got Item
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	// Comparar valores relevantes
	if got.Ticker != item.Ticker || got.Company != item.Company || got.Action != item.Action {
		t.Fatalf("mismatch after roundtrip: got=%+v want=%+v", got, item)
	}
	if float64(got.TargetFrom) != float64(item.TargetFrom) {
		t.Fatalf("money mismatch: got=%v want=%v", got.TargetFrom, item.TargetFrom)
	}
}
