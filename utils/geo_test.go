package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDistance(t *testing.T) {
	distance := CalculateDistance(2.0508377164686333, 45.32936363251881, 2.053290964628916, 45.32910968662657, "")
	assert.Equal(t, 274.23128658767484, distance)

}
