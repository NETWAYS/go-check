package convert

import (
	"testing"
)

type bytesIECTestCase struct {
	name     string
	input    BytesIEC
	expected string
}

func TestBytesIEC_HumanReadable(t *testing.T) {
	tests := []bytesIECTestCase{
		{name: "zero", input: BytesIEC(0), expected: "0B"},
		{name: "less than 1K", input: BytesIEC(999), expected: "999B"},
		{name: "999 KiB", input: BytesIEC(999 * 1024), expected: "999KiB"},
		{name: "999 MiB", input: BytesIEC(999 * 1024 * 1024), expected: "999MiB"},
		{name: "999 GiB", input: BytesIEC(999 * 1024 * 1024 * 1024), expected: "999GiB"},
		{name: "999 TiB", input: BytesIEC(999 * 1024 * 1024 * 1024 * 1024), expected: "999TiB"},
		{name: "4 PiB", input: BytesIEC(4 * 1024 * 1024 * 1024 * 1024 * 1024), expected: "4PiB"},
		{name: "4096 PiB", input: BytesIEC(4 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024), expected: "4096PiB"},
		{name: "1263 MiB", input: BytesIEC(1263 * 1024 * 1024), expected: "1263MiB"},
		{name: "100 MiB", input: BytesIEC(100 * 1024 * 1024), expected: "100MiB"},
		{name: "123.05 MiB", input: BytesIEC(129032519), expected: "123.05MiB"},
		{name: "14.67 GiB", input: BytesIEC(15756365824), expected: "14.67GiB"},
		{name: "1024 KiB", input: BytesIEC(1024 * 1024), expected: "1024KiB"},
		{name: "2 MiB", input: BytesIEC(2 * 1024 * 1024), expected: "2MiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.HumanReadable(); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
			if got := tt.input.String(); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}
