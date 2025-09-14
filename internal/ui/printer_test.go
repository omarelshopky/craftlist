package ui

import (
	"testing"
)

func TestHumanizeNumber(t *testing.T) {
	p := NewPrinter()

	tests := []struct {
		name     string
		input    int
		expected string
	}{
		{
			name:     "small number",
			input:    123,
			expected: "123",
		},
		{
			name:     "zero",
			input:    0,
			expected: "0",
		},
		{
			name:     "negative number",
			input:    -456,
			expected: "-456",
		},
		{
			name:     "large number",
			input:    123456789,
			expected: "123,456,789",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := p.humanizeNumber(tc.input)
			if result != tc.expected {
				t.Errorf("humanizeNumber(%d) = %s; expected %s", tc.input, result, tc.expected)
			}
		})
	}
}
