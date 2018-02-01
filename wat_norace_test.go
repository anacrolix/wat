// +build !race

package wat

import (
	"sync"
	"testing"

	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/require"
)

func TestSyncPoolZeroesItems(t *testing.T) {
	p := sync.Pool{
		New: func() interface{} {
			return 42
		},
	}
	require.EqualValues(t, 42, p.Get())
	p.Put(1)
	require.EqualValues(t, 1, p.Get())
	require.EqualValues(t, 42, p.Get())
	p.Put([]int{1, 2})
	require.EqualValues(t, []int{1, 2}, p.Get())
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
