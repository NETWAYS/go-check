package main

import (
	"os"
	"testing"

	"github.com/NETWAYS/go-check/testhelper"
)

func TestMyMain(t *testing.T) {
	actual := testhelper.RunMainTest(main)

	expected := `[WARNING] - Check2
\_ [OK] Check1
\_ [WARNING] Check2
\_ [WARNING] Check3
\_ [OK] Check4
|foo=23 bar=42 'foo2 bar'=46

would exit with code 1
`

	if actual != expected {
		t.Fatalf("expected %q\n, got %q", expected, actual)
	}
}

func TestMain(m *testing.M) {
	testhelper.EnableTestMode()
	os.Exit(m.Run())
}
