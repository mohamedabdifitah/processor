package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveObjectID(t *testing.T) {
	id := RemoveObjectID("ObjectID(\"655230599abddc91ae792158\")")
	assert.Equal(t, id, "655230599abddc91ae792158")

}
