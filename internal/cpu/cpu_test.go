package cpu

import (
	"reflect"
	"testing"
)

func TestFormCpuInfo(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected CPUStats
	}{
		{
			name:     "Normal case",
			input:    []string{" 2.0 us", " 1.5 sy", " 95.0 id", " 1.5 wa"},
			expected: CPUStats{Usr: 2.0, Sys: 1.5, Idle: 95.0, Iowait: 1.5},
		},
		{
			name:     "Missing fields",
			input:    []string{" 2.0 us", " 98.0 id"},
			expected: CPUStats{Usr: 2.0, Sys: 0, Idle: 98.0, Iowait: 0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formCPUInfo(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("formCpuInfo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParamToFloat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float32
	}{
		{"Simple case", " 2.0 us", 2.0},
		{"Integer value", " 5 sy", 5.0},
		{"Trailing spaces", " 3.5 id ", 3.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := paramToFloat(tt.input)
			if result != tt.expected {
				t.Errorf("paramToFloat(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}
