package convert

import (
	"testing"
)

func TestBytesIEC_HumanReadable(t *testing.T) {
	if BytesIEC(0).HumanReadable() != "0B" {
		t.Fatalf("expected '0B', got %s", BytesIEC(0).HumanReadable())
	}
	if BytesIEC(999).HumanReadable() != "999B" {
		t.Fatalf("expected '999B', got %s", BytesIEC(999).HumanReadable())
	}
	if BytesIEC(999*1024).HumanReadable() != "999KiB" {
		t.Fatalf("expected '999KiB', got %s", BytesIEC(999*1024).HumanReadable())
	}
	if BytesIEC(999*1024*1024).HumanReadable() != "999MiB" {
		t.Fatalf("expected '999MiB', got %s", BytesIEC(999*1024*1024).HumanReadable())
	}
	if BytesIEC(999*1024*1024*1024).HumanReadable() != "999GiB" {
		t.Fatalf("expected '999GiB', got %s", BytesIEC(999*1024*1024*1024).HumanReadable())
	}
	if BytesIEC(999*1024*1024*1024*1024).HumanReadable() != "999TiB" {
		t.Fatalf("expected '999TiB', got %s", BytesIEC(999*1024*1024*1024*1024).HumanReadable())
	}
	if BytesIEC(4*1024*1024*1024*1024*1024).HumanReadable() != "4PiB" {
		t.Fatalf("expected '4PiB', got %s", BytesIEC(4*1024*1024*1024*1024*1024).HumanReadable())
	}
	if BytesIEC(4*1024*1024*1024*1024*1024*1024).HumanReadable() != "4096PiB" {
		t.Fatalf("expected '4096PiB', got %s", BytesIEC(4*1024*1024*1024*1024*1024*1024).HumanReadable())
	}
	if BytesIEC(1263*1024*1024).HumanReadable() != "1263MiB" {
		t.Fatalf("expected '1263MiB', got %s", BytesIEC(1263*1024*1024).HumanReadable())
	}
	if BytesIEC(100*1024*1024).HumanReadable() != "100MiB" {
		t.Fatalf("expected '100MiB', got %s", BytesIEC(100*1024*1024).HumanReadable())
	}
	if BytesIEC(129032519).HumanReadable() != "123.05MiB" {
		t.Fatalf("expected '123.05MiB', got %s", BytesIEC(129032519).HumanReadable())
	}
	if BytesIEC(15756365824).HumanReadable() != "14.67GiB" {
		t.Fatalf("expected '14.67GiB', got %s", BytesIEC(15756365824).HumanReadable())
	}
	if BytesIEC(1024*1024).HumanReadable() != "1024KiB" {
		t.Fatalf("expected '1024KiB', got %s", BytesIEC(1024*1024).HumanReadable())
	}
	if BytesIEC(2*1024*1024).HumanReadable() != "2MiB" {
		t.Fatalf("expected '2MiB', got %s", BytesIEC(2*1024*1024).HumanReadable())
	}
}
