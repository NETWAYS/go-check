package convert

import (
	"testing"
)

func TestBytesSI_HumanReadable(t *testing.T) {
	if BytesSI(0).HumanReadable() != "0B" {
		t.Fatalf("expected '0B', got %s", BytesSI(0).HumanReadable())
	}
	if BytesSI(999).HumanReadable() != "999B" {
		t.Fatalf("expected '999B', got %s", BytesSI(999).HumanReadable())
	}
	if BytesSI(999*1000).HumanReadable() != "999KB" {
		t.Fatalf("expected '999KB', got %s", BytesSI(999*1000).HumanReadable())
	}
	if BytesSI(999*1000*1000).HumanReadable() != "999MB" {
		t.Fatalf("expected '999MB', got %s", BytesSI(999*1000*1000).HumanReadable())
	}
	if BytesSI(999*1000*1000*1000).HumanReadable() != "999GB" {
		t.Fatalf("expected '999GB', got %s", BytesSI(999*1000*1000*1000).HumanReadable())
	}
	if BytesSI(999*1000*1000*1000*1000).HumanReadable() != "999TB" {
		t.Fatalf("expected '999TB', got %s", BytesSI(999*1000*1000*1000*1000).HumanReadable())
	}
	if BytesSI(4*1000*1000*1000*1000*1000).HumanReadable() != "4PB" {
		t.Fatalf("expected '4PB', got %s", BytesSI(4*1000*1000*1000*1000*1000).HumanReadable())
	}
	if BytesSI(4*1000*1000*1000*1000*1000*1000).HumanReadable() != "4000PB" {
		t.Fatalf("expected '4000PB', got %s", BytesSI(4*1000*1000*1000*1000*1000*1000).HumanReadable())
	}
	if BytesSI(4*1000*1000*1000*1000).HumanReadable() != "4TB" {
		t.Fatalf("expected '4TB', got %s", BytesSI(4*1000*1000*1000*1000).HumanReadable())
	}
	if BytesSI(4*1000*1000*1000*1000*1000).HumanReadable() != "4PB" {
		t.Fatalf("expected '4PB', got %s", BytesSI(4*1000*1000*1000*1000*1000).HumanReadable())
	}
	if BytesSI(1263*1000*1000).HumanReadable() != "1263MB" {
		t.Fatalf("expected '1263MB', got %s", BytesSI(1263*1000*1000).HumanReadable())
	}
	if BytesSI(123050*1000).HumanReadable() != "123.05MB" {
		t.Fatalf("expected '123.05MB', got %s", BytesSI(123050*1000).HumanReadable())
	}
	if BytesSI(14670*1000*1000).HumanReadable() != "14.67GB" {
		t.Fatalf("expected '14.67GB', got %s", BytesSI(14670*1000*1000).HumanReadable())
	}
	if BytesSI(1000*1000).HumanReadable() != "1000KB" {
		t.Fatalf("expected '1000KB', got %s", BytesSI(1000*1000).HumanReadable())
	}
	if BytesSI(2*1000*1000).HumanReadable() != "2MB" {
		t.Fatalf("expected '2MB', got %s", BytesSI(2*1000*1000).HumanReadable())
	}
	if BytesSI(3*1000*1000).HumanReadable() != "3MB" {
		t.Fatalf("expected '3MB', got %s", BytesSI(3*1000*1000).HumanReadable())
	}
}
