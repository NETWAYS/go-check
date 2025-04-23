package convert

import (
	"testing"
)

func TestParseBytes(t *testing.T) {
	b, err := ParseBytes("1024")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := b.(BytesIEC); !ok {
		t.Fatalf("expected type BytesIEC, got %T", b)
	}
	if b.Bytes() != 1024 {
		t.Fatalf("expected 1024 bytes, got %d", b.Bytes())
	}

	b, err = ParseBytes("1MB")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := b.(BytesSI); !ok {
		t.Fatalf("expected type BytesSI, got %T", b)
	}
	if b.Bytes() != 1000*1000 {
		t.Fatalf("expected 1000000 bytes, got %d", b.Bytes())
	}

	b, err = ParseBytes("1 MiB")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := b.(BytesIEC); !ok {
		t.Fatalf("expected type BytesIEC, got %T", b)
	}
	if b.Bytes() != 1024*1024 {
		t.Fatalf("expected 1048576 bytes, got %d", b.Bytes())
	}

	b, err = ParseBytes("100MB")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	b, err = ParseBytes("100MiB")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	b, err = ParseBytes("  23   GiB  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if _, ok := b.(BytesIEC); !ok {
		t.Fatalf("expected type BytesIEC, got %T", b)
	}
	if b.Bytes() != 23*1024*1024*1024 {
		t.Fatalf("expected 24742653952 bytes, got %d", b.Bytes())
	}

	b, err = ParseBytes("1.2.3.4MB")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if b != nil {
		t.Fatalf("expected nil, got %v", b)
	}

	b, err = ParseBytes("1PHD")
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if b != nil {
		t.Fatalf("expected nil, got %v", b)
	}
}
