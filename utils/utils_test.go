package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNumbers(t *testing.T) {
	assert.Equal(t, 10, len(GenerateNumbers(10)))
	assert.Empty(t, 0, GenerateNumbers(0))
	generate1 := GenerateNumbers(20)
	time.Sleep(2 * time.Nanosecond)
	generate2 := GenerateNumbers(20)
	assert.NotEqual(t, generate1, generate2)
}
