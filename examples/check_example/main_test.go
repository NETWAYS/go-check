package main

import (
	"os"
	"regexp"
	"testing"

	"github.com/NETWAYS/go-check/testhelper"
	"github.com/stretchr/testify/assert"
)

func TestMyMain(t *testing.T) {
	stdout := testhelper.RunMainTest(main, "--help")
	assert.Regexp(t, regexp.MustCompile(`would exit with code 3`), stdout)

	stdout = testhelper.RunMainTest(main, "--warning", "20")
	assert.Regexp(t, regexp.MustCompile(`^\[OK\] - value is 10\nwould exit with code 0\n$`), stdout)

	stdout = testhelper.RunMainTest(main, "--warning", "10", "--value", "11")
	assert.Regexp(t, regexp.MustCompile(`^\[WARNING\] - value is 11\nwould exit with code 1\n$`), stdout)
}

func TestMain(m *testing.M) {
	testhelper.EnableTestMode()
	os.Exit(m.Run())
}
