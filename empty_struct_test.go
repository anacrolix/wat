// +build !go1.9

package wat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// struct{} values always have the same address.
func testEmptyStructEquality(t *testing.T) {
	assert.True(t, struct{}{} == struct{}{})
	assert.False(t, new(bool) == new(bool))
	assert.True(t, new(struct{}) == new(struct{}))
	assert.True(t, &struct{}{} == &struct{}{})
	var a, b struct{}
	assert.True(t, &a == &b)
}
