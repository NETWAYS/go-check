package convert

import (
	"testing"
)

type bytesSITestCase struct {
	input    BytesSI
	expected string
}

func TestBytesSI_HumanReadable2(t *testing.T) {
	tests := []bytesSITestCase{
		{BytesSI(0), "0B"},
		{BytesSI(999), "999B"},
		{BytesSI(999 * 1000), "999KB"},
		{BytesSI(999 * 1000 * 1000), "999MB"},
		{BytesSI(999 * 1000 * 1000 * 1000), "999GB"},
		{BytesSI(999 * 1000 * 1000 * 1000 * 1000), "999TB"},
		{BytesSI(4 * 1000 * 1000 * 1000 * 1000 * 1000), "4PB"},
		{BytesSI(4 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000), "4000PB"},
		{BytesSI(4 * 1000 * 1000 * 1000 * 1000), "4TB"},
		{BytesSI(4 * 1000 * 1000 * 1000 * 1000 * 1000), "4PB"},
		{BytesSI(1263 * 1000 * 1000), "1263MB"},
		{BytesSI(123050 * 1000), "123.05MB"},
		{BytesSI(14670 * 1000 * 1000), "14.67GB"},
		{BytesSI(1000 * 1000), "1000KB"},
		{BytesSI(2 * 1000 * 1000), "2MB"},
		{BytesSI(3 * 1000 * 1000), "3MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.input.HumanReadable(); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
			if got := tt.input.String(); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
