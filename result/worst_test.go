package result_test

import (
	"fmt"
	"github.com/NETWAYS/go-check"
	"github.com/NETWAYS/go-check/result"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorstState(t *testing.T) {
	assert.Equal(t, 3, result.WorstState(3))
	assert.Equal(t, 2, result.WorstState(2))
	assert.Equal(t, 1, result.WorstState(1))
	assert.Equal(t, 0, result.WorstState(0))

	assert.Equal(t, 2, result.WorstState(0, 1, 2, 3))
	assert.Equal(t, 3, result.WorstState(0, 1, 3))
	assert.Equal(t, 1, result.WorstState(1, 0, 0))
	assert.Equal(t, 0, result.WorstState(0, 0, 0))

	assert.Equal(t, 3, result.WorstState(-1))
	assert.Equal(t, 3, result.WorstState(4))
}

func ExampleWorstState() {
	state := result.WorstState(check.Unknown, check.Critical, check.OK)
	fmt.Println(state)
	// Output:
	// 2
}
