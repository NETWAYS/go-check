package main

import (
	"os"
	"strings"
	"testing"

	"github.com/NETWAYS/go-check/testhelper"
)

func TestMyMain(t *testing.T) {
	actual := testhelper.RunMainTest(main, "--help")
	expected := `would exit with code 3`

	if !strings.Contains(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	actual = testhelper.RunMainTest(main, "--warning", "20")
	expected = "[OK] - value is 10"

	if !strings.Contains(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	actual = testhelper.RunMainTest(main, "--warning", "10", "--value", "11")
	expected = "[WARNING] - value is 11"

	if !strings.Contains(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}

func TestMain(m *testing.M) {
	testhelper.EnableTestMode()
	os.Exit(m.Run())
}
