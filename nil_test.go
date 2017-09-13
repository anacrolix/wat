package wat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTestifyNilValue(t *testing.T) {
	var a interface{}
	assert.Nil(t, a)
	assert.True(t, a == nil)
	a = (*int)(nil)
	assert.True(t, a != nil)
	// Testify incorrectly says a is nil.
	assert.Nil(t, a)
}
