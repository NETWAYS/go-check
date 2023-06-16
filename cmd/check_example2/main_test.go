package main

import (
	"os"
	"testing"

	"github.com/NETWAYS/go-check/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestMyMain(t *testing.T) {
	stdout := testhelper.RunMainTest(main)

	resultString := `WARNING - states: warning=1 ok=1
\_ [OK] Check1
\_ [WARNING] Check2
|foo=23 bar=42 'foo2 bar'=46

would exit with code 1
`

	assert.Equal(t, resultString, stdout)
}

func TestMain(m *testing.M) {
	testhelper.EnableTestMode()
	os.Exit(m.Run())
}
