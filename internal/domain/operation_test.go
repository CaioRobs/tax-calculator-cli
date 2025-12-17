package domain

import (
	"encoding/json"
	"testing"
)

func TestOperation_JSON_Unmarshal_HappyPath(t *testing.T) {
	input := []byte(`{"operation":"buy","unit-cost":10.5,"quantity":3}`)
	var op Operation
	if err := json.Unmarshal(input, &op); err != nil {
		t.Fatalf("unexpected error on unmarshal: %v", err)
	}
	if op.Operation != "buy" {
		t.Fatalf("operation: expected buy, got %q", op.Operation)
	}
	if op.UnitCost != 10.5 {
		t.Fatalf("unit-cost: expected 10.5, got %v", op.UnitCost)
	}
	if op.Quantity != 3 {
		t.Fatalf("quantity: expected 3, got %d", op.Quantity)
	}
}

func TestOperation_JSON_Marshal_Tags(t *testing.T) {
	op := Operation{Operation: "sell", UnitCost: 20.0, Quantity: 5}
	b, err := json.Marshal(op)
	if err != nil {
		t.Fatalf("unexpected error on marshal: %v", err)
	}
	// Verifica presença e valores das chaves com os nomes dos tags
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unexpected error parsing marshaled json: %v", err)
	}
	if got, ok := m["operation"].(string); !ok || got != "sell" {
		t.Fatalf("operation key missing or wrong: %v", m["operation"])
	}
	if got, ok := m["unit-cost"].(float64); !ok || got != 20.0 {
		t.Fatalf("unit-cost key missing or wrong: %v", m["unit-cost"])
	}
	if got, ok := m["quantity"].(float64); !ok || got != 5.0 {
		t.Fatalf("quantity key missing or wrong: %v", m["quantity"])
	}
}

func TestOperation_JSON_Unmarshal_InvalidQuantityType(t *testing.T) {
	// quantity como número não inteiro deve falhar ao unmarshal para int
	input := []byte(`{"operation":"buy","unit-cost":10.0,"quantity":3.14}`)
	var op Operation
	if err := json.Unmarshal(input, &op); err == nil {
		t.Fatalf("expected error for non-integer quantity, got nil")
	}
}

func TestOperation_JSON_Unmarshal_MissingFields_Defaults(t *testing.T) {
	// campos ausentes devem virar zero-values em Go
	input := []byte(`{"operation":"buy"}`)
	var op Operation
	if err := json.Unmarshal(input, &op); err != nil {
		t.Fatalf("unexpected error on unmarshal: %v", err)
	}
	if op.Operation != "buy" {
		t.Fatalf("operation: expected buy, got %q", op.Operation)
	}
	if op.UnitCost != 0 {
		t.Fatalf("unit-cost: expected 0, got %v", op.UnitCost)
	}
	if op.Quantity != 0 {
		t.Fatalf("quantity: expected 0, got %d", op.Quantity)
	}
}
