package status

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestString(t *testing.T) {
	assert.Equal(t, "OK", String(OK))
}
