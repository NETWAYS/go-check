package convert

import (
	"fmt"
	"testing"
)

func TestParseBytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected uint64
	}{
		{name: "zero", expected: 0, input: "0B"},
		{name: "0", expected: 0, input: "0"},
		{name: "1024", expected: 1024, input: "1024"},
		{name: "2 MiB", expected: 2 * 1024 * 1024, input: "2MiB"},
		{name: "1MB", expected: 1000 * 1000, input: "1MB"},
		{name: "100MB", expected: 1000 * 1000 * 100, input: "100MB"},
		{name: " 23 GiB", expected: 23 * 1024 * 1024 * 1024, input: "  23   GiB  "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseBytes(tt.input)

			if err != nil {
				t.Fatalf("did not expect error, got %v", err)
			}

			if actual != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, actual)
			}
		})
	}
}

func TestParseBytes_WithError(t *testing.T) {
	tests := []struct {
		input string
	}{
		{
			input: "unittest",
		},
		{
			input: "1PHD",
		},
		{
			input: "",
		},
		{
			input: "1.2.3.4MB",
		},
		{
			input: "-1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := ParseBytes(tt.input)

			if err == nil {
				t.Fatalf("did expect error, got nil")
			}
		})
	}
}

func ExampleParseBytes() {
	b, err := ParseBytes("999KiB")

	if err != nil {
		panic("Could not parse string into number of bytes")
	}

	fmt.Println(b)
	// Output: 1022976
}

func TestBytesIEC_HumanReadable(t *testing.T) {
	tests := []struct {
		name     string
		input    uint64
		expected string
	}{
		{name: "zero", input: 0, expected: "0B"},
		{name: "less than 1K", input: 999, expected: "999B"},
		{name: "999 KiB", input: 999 * 1024, expected: "999KiB"},
		{name: "999 MiB", input: 999 * 1024 * 1024, expected: "999MiB"},
		{name: "999 GiB", input: 999 * 1024 * 1024 * 1024, expected: "999GiB"},
		{name: "999 TiB", input: 999 * 1024 * 1024 * 1024 * 1024, expected: "999TiB"},
		{name: "4 PiB", input: 4 * 1024 * 1024 * 1024 * 1024 * 1024, expected: "4PiB"},
		{name: "4096 PiB", input: 4 * 1024 * 1024 * 1024 * 1024 * 1024 * 1024, expected: "4096PiB"},
		{name: "1263 MiB", input: 1263 * 1024 * 1024, expected: "1263MiB"},
		{name: "100 MiB", input: 100 * 1024 * 1024, expected: "100MiB"},
		{name: "123.05 MiB", input: 129032519, expected: "123.05MiB"},
		{name: "14.67 GiB", input: 15756365824, expected: "14.67GiB"},
		{name: "1024 KiB", input: 1024 * 1024, expected: "1024KiB"},
		{name: "2 MiB", input: 2 * 1024 * 1024, expected: "2MiB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesIEC(tt.input); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func ExampleBytesIEC() {
	b := BytesIEC(999 * 1024)

	fmt.Println(b)
	// Output: 999KiB
}

func TestBytesSI_HumanReadable(t *testing.T) {
	tests := []struct {
		input    uint64
		expected string
	}{
		{input: 0, expected: "0B"},
		{input: 999, expected: "999B"},
		{input: 999 * 1000, expected: "999KB"},
		{input: 999 * 1000 * 1000, expected: "999MB"},
		{input: 999 * 1000 * 1000 * 1000, expected: "999GB"},
		{input: 999 * 1000 * 1000 * 1000 * 1000, expected: "999TB"},
		{input: 4 * 1000 * 1000 * 1000 * 1000 * 1000, expected: "4PB"},
		{input: 4 * 1000 * 1000 * 1000 * 1000 * 1000 * 1000, expected: "4000PB"},
		{input: 4 * 1000 * 1000 * 1000 * 1000, expected: "4TB"},
		{input: 4 * 1000 * 1000 * 1000 * 1000 * 1000, expected: "4PB"},
		{input: 1263 * 1000 * 1000, expected: "1263MB"},
		{input: 123050 * 1000, expected: "123.05MB"},
		{input: 14670 * 1000 * 1000, expected: "14.67GB"},
		{input: 1000 * 1000, expected: "1000KB"},
		{input: 2 * 1000 * 1000, expected: "2MB"},
		{input: 3 * 1000 * 1000, expected: "3MB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := BytesSI(tt.input); got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func ExampleBytesSI() {
	b := BytesSI(999 * 1000)

	fmt.Println(b)
	// Output: 999KB
}
