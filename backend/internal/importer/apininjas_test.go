package importer

import "testing"

func TestToFloat(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect float64
	}{
		{2.5, 2.5},
		{2, 2.0},
		{"3.5", 3.5},
		{"invalid", 0},
		{nil, 0},
	}

	for _, tc := range tests {
		got := toFloat(tc.input)
		if got != tc.expect {
			t.Errorf("toFloat(%v) = %v; want %v", tc.input, got, tc.expect)
		}
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input  interface{}
		expect string
	}{
		{"hello", "hello"},
		{123, "123"},
		{45.6, "45.6"},
		{nil, ""}, // Assuming we want empty string for nil
	}

	for _, tc := range tests {
		got := toString(tc.input)
		if got != tc.expect {
			t.Errorf("toString(%v) = '%v'; want '%v'", tc.input, got, tc.expect)
		}
	}
}
