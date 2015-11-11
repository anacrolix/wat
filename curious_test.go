package curious

import (
	"math/rand"
	"mime"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"testing"
	"unsafe"

	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/require"
)

// Looks like if we append endlessly, we're given new backing arrays.
func BenchmarkEndlessAppend(t *testing.B) {
	for range make([]struct{}, t.N) {
		sl := make([]int, 0x10000)
		for i := 0; i < 0x10000; i++ {
			sl = append(sl, make([]int, 0x10000)...)
			sl = sl[len(sl)/2:]
		}
	}
}

// Is a zero value slice the same as a nil slice of the same type?
func TestSliceZeroValue(t *testing.T) {
	sl := []byte{}
	p := (*reflect.SliceHeader)(unsafe.Pointer(&sl))
	t.Log(*p)
	var sl1 []byte
	p = (*reflect.SliceHeader)(unsafe.Pointer(&sl1))
	t.Log(*p)
}

// Range takes it's own slice reference. i is reset every time around the
// loop.
func TestRangeAlterIndex(t *testing.T) {
	b := []byte{'a', 'b', 'c'}
	var c []byte
	for i, j := range b {
		b = nil
		t.Log(i, j)
		i--
		t.Log(i)
		c = append(c, 'd'+byte(i))[1:]
	}
	t.Log(c)
}

type realSlimShady struct{}

func (realSlimShady) StandUp()    {}
func (realSlimShady) Cross8Mile() {}

type slimShady interface {
	StandUp()
	// Only the real Slim Shady can cross 8 Mile
	Cross8Mile()
}

// Posers can stand up
type poser interface {
	StandUp()
}

// Eminem looks like a poser
type eminem struct {
	poser
}

func TestAssertComposedType(t *testing.T) {
	// Eminem implements poser, but he's the real Shady
	var realShady poser = eminem{realSlimShady{}}
	// He should be able to cross 8 Mile.
	_, ok := realShady.(slimShady)
	if ok {
		// If he can, Go has changed.
		t.FailNow()
	}
}

func TestMakesliceNegative(t *testing.T) {
	var l int64
	l = -1
	defer func() {
		r := recover()
		if !strings.Contains(r.(error).Error(), "len out of range") {
			t.FailNow()
		}
	}()
	_ = make([]byte, l)
}

// See if named returned values are set even if return values are all
// specified in the return statement. It seems that they are.
func TestDirectReturnSetsNamedValues(t *testing.T) {
	var intercepted bool
	f := func() (named bool) {
		defer func() {
			// Saved the named return value.
			intercepted = named
		}()
		// Bypass the named return value.
		return true
	}
	f()
	t.Log(intercepted)
	if intercepted != true {
		t.FailNow()
	}
}

var constantSlice = []string{"a", "b", "c", "d", "e"}

func sliceIndex(i int) string {
	return []string{"a", "b", "c", "d", "e"}[i]
}

func arrayIndex(i int) string {
	return [...]string{"a", "b", "c", "d", "e"}[i]
}

func BenchmarkConstantIndex(b *testing.B) {
	for range iter.N(b.N) {
		constantIndex(rand.Intn(5))
	}
}

func constantIndex(i int) string {
	return constantSlice[i]
}

func BenchmarkSliceIndex(b *testing.B) {
	for range iter.N(b.N) {
		sliceIndex(rand.Intn(5))
	}
}

func BenchmarkArrayIndex(b *testing.B) {
	for range iter.N(b.N) {
		arrayIndex(rand.Intn(5))
	}
}

func TestChannelReference(t *testing.T) {
	a := struct {
		b chan struct{}
	}{make(chan struct{})}
	c := struct {
		d <-chan struct{}
	}{a.b}
	e := c
	select {
	case <-e.d:
		t.FailNow()
	default:
	}
	go close(a.b)
	<-e.d
}

func TestExtensionMimeTypes(t *testing.T) {
	t.Log(mime.TypeByExtension("/some/path.mp4"))
	t.Log(mime.TypeByExtension(".avi"))
	t.Log(mime.TypeByExtension("/some/path.mp4"))
	t.Log(mime.TypeByExtension(".ogv"))
}

func TestSliceStringsAllocation(t *testing.T) {
	b := make([]byte, 1000000000)
	for i := range b {
		b[i] = byte(i)
	}
	s := string(b)
	runtime.GC()
	var ss []string
	for i := range iter.N(1000) {
		ss = append(ss, s[i:])
	}
	// time.Sleep(30 * time.Second)
}

func touchedBytes(n int) []byte {
	ret := make([]byte, n)
	for i := range ret {
		ret[i] = byte(i)
	}
	return ret
}

const oneGB = 1000000000

func TestSlicedStringTrimmed(t *testing.T) {
	s := string([]byte(string(touchedBytes(oneGB))[:1000]))
	runtime.GC()
	// time.Sleep(time.Second * 30)
	s += ""
}

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

func TestFmtF(t *testing.T) {
	t.Logf("%+q", '\xcf')
	t.Logf("%#q", '\xcf')
	t.Logf("%q", '\xcf')
}
