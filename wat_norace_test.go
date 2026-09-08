// +build !race

package wat

import (
	"sync"
	"testing"

	"github.com/bradfitz/iter"

	"github.com/go-quicktest/qt"
)

func TestSyncPoolZeroesItems(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return 42
		},
	}
	qt.Assert(t, qt.Equals(p.Get(), 42))
	p.Put(1)
	qt.Assert(t, qt.Equals(p.Get(), 1))
	qt.Assert(t, qt.Equals(p.Get(), 42))
	p.Put([]int{1, 2})
	qt.Assert(t, qt.DeepEquals[any](p.Get(), []int{1, 2}))
	for range iter.N(100) {
		p.Put(make([]byte, 100000))
	}
	got := 0
	for range iter.N(100) {
		if _, ok := p.Get().([]byte); ok {
			got++
		}
	}
	t.Logf("got %d back", got)
}
