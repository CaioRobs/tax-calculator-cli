package services

import (
	"math"
	"testing"
)

func nearlyEqual(a, b, tol float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if math.IsInf(a, 1) && math.IsInf(b, 1) {
		return true
	}
	if math.IsInf(a, -1) && math.IsInf(b, -1) {
		return true
	}
	return math.Abs(a-b) <= tol
}

func TestRound2_CommonCases(t *testing.T) {
	tests := []struct {
		in   float64
		out  float64
		name string
	}{
		{in: 16.666666, out: 16.67, name: "ceil-like to two decimals"},
		{in: 16.664, out: 16.66, name: "down on < .005"},
		{in: 16.665, out: 16.67, name: "half away from zero positive"},
		{in: -16.665, out: -16.67, name: "half away from zero negative"},
		{in: 0.0, out: 0.0, name: "zero"},
	}
	for _, tt := range tests {
		got := Round2(tt.in)
		if !nearlyEqual(got, tt.out, 1e-9) {
			t.Fatalf("%s: Round2(%v) = %v, want %v", tt.name, tt.in, got, tt.out)
		}
	}
}

func TestRoundTo_VariousPlaces(t *testing.T) {
	tests := []struct {
		in     float64
		places int
		out    float64
		name   string
	}{
		{in: 1.2344, places: 3, out: 1.234, name: "3 places down"},
		{in: 1.2345, places: 3, out: 1.235, name: "3 places half away"},
		{in: 123.0, places: 0, out: 123.0, name: "0 places"},
		{in: -1.235, places: 2, out: -1.24, name: "negative 2 places"},
	}
	for _, tt := range tests {
		got := RoundTo(tt.in, tt.places)
		if !nearlyEqual(got, tt.out, 1e-9) {
			t.Fatalf("%s: RoundTo(%v,%d) = %v, want %v", tt.name, tt.in, tt.places, got, tt.out)
		}
	}
}

func TestRoundTo_SpecialValues(t *testing.T) {
	if got := RoundTo(math.NaN(), 2); !math.IsNaN(got) {
		t.Fatalf("NaN should propagate, got %v", got)
	}
	if got := RoundTo(math.Inf(1), 2); !math.IsInf(got, 1) {
		t.Fatalf("+Inf should propagate, got %v", got)
	}
	if got := RoundTo(math.Inf(-1), 2); !math.IsInf(got, -1) {
		t.Fatalf("-Inf should propagate, got %v", got)
	}
}
