package wat

import (
	"testing"

	"github.com/go-quicktest/qt"
)

// An interface value holding a typed nil pointer is not itself nil.
// testify's assert.Nil incorrectly reported it as nil; qt.IsNil does not.
func TestNilInterfaceValue(t *testing.T) {
	var a interface{}
	qt.Check(t, qt.IsNil(a))
	qt.Check(t, qt.IsTrue(a == nil))
	a = (*int)(nil)
	qt.Check(t, qt.IsTrue(a != nil))
	qt.Check(t, qt.IsNotNil(a))
}
