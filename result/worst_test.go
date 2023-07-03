package result

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWorstState(t *testing.T) {

	assert.Equal(t, 3, WorstState(3))
	assert.Equal(t, 2, WorstState(2))
	assert.Equal(t, 1, WorstState(1))
	assert.Equal(t, 0, WorstState(0))

	assert.Equal(t, 2, WorstState(0, 1, 2, 3))
	assert.Equal(t, 3, WorstState(0, 1, 3))
	assert.Equal(t, 1, WorstState(1, 0, 0))
	assert.Equal(t, 0, WorstState(0, 0, 0))

	assert.Equal(t, 3, WorstState(-1))
	assert.Equal(t, 3, WorstState(4))
}

func BenchmarkWorstState(b *testing.B) {
	b.ReportAllocs()

	// Initialize slice for benchmarking
	states := make([]int, 0, 100)
	for i := 0; i < 100; i++ {
		states = append(states, i%4)
	}

	for i := 0; i < b.N; i++ {
		s := WorstState(states...)
		_ = s
	}
}
