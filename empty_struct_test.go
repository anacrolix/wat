// +build !go1.9

package wat

import (
	"testing"

	"github.com/go-quicktest/qt"
)

// struct{} values always have the same address.
func testEmptyStructEquality(t *testing.T) {
	qt.Check(t, qt.IsTrue(struct{}{} == struct{}{}))
	qt.Check(t, qt.IsFalse(new(bool) == new(bool)))
	qt.Check(t, qt.IsTrue(new(struct{}) == new(struct{})))
	qt.Check(t, qt.IsTrue(&struct{}{} == &struct{}{}))
	var a, b struct{}
	qt.Check(t, qt.IsTrue(&a == &b))
}
