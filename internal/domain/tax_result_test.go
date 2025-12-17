package domain

import (
	"encoding/json"
	"math"
	"testing"
)

func TestTaxResult_JSON_Unmarshal_HappyPath(t *testing.T) {
	input := []byte(`{"tax":123.45}`)
	var tr TaxResult
	if err := json.Unmarshal(input, &tr); err != nil {
		t.Fatalf("unexpected error on unmarshal: %v", err)
	}
	if tr.Tax != 123.45 {
		t.Fatalf("tax: expected 123.45, got %v", tr.Tax)
	}
}

func TestTaxResult_JSON_Marshal_Tags(t *testing.T) {
	tr := TaxResult{Tax: 20.0}
	b, err := json.Marshal(tr)
	if err != nil {
		t.Fatalf("unexpected error on marshal: %v", err)
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatalf("unexpected error parsing marshaled json: %v", err)
	}
	if got, ok := m["tax"].(float64); !ok || got != 20.0 {
		t.Fatalf("tax key missing or wrong: %v", m["tax"])
	}
}

func TestTaxResult_JSON_Unmarshal_MissingField_Default(t *testing.T) {
	input := []byte(`{}`)
	var tr TaxResult
	if err := json.Unmarshal(input, &tr); err != nil {
		t.Fatalf("unexpected error on unmarshal: %v", err)
	}
	if tr.Tax != 0 {
		t.Fatalf("tax: expected 0, got %v", tr.Tax)
	}
}

func TestTaxResult_JSON_Marshal_NaN_Error(t *testing.T) {
	tr := TaxResult{Tax: math.NaN()}
	if _, err := json.Marshal(tr); err == nil {
		t.Fatalf("expected error when marshaling NaN, got nil")
	}
}
