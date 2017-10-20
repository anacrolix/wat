package wat

import (
	"fmt"
	"math"
	"math/rand"
	"mime"
	"net"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	_ "github.com/anacrolix/envpprof"

	"github.com/bradfitz/iter"
	"github.com/stretchr/testify/assert"
)

// Looks like if we append endlessly, we're given new backing arrays.
func BenchmarkEndlessAppend(t *testing.B) {
	for range make([]struct{}, t.N) {
		sl := make([]int, 0x10000)
		for i := 0; i < 0x100; i++ {
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

func touchedBytes(n int) []byte {
	ret := make([]byte, n)
	for i := range ret {
		ret[i] = byte(i)
	}
	return ret
}

const oneGB = 1000000000

func TestFmtF(t *testing.T) {
	t.Logf("%+q", '\xcf')
	t.Logf("%#q", '\xcf')
	t.Logf("%q", '\xcf')
}

func TestFmtDFloat(t *testing.T) {
	var f float64 = 42.123
	assert.EqualValues(t, "42", fmt.Sprintf("%d", int(f)))
}

func TestInt64Wrap(t *testing.T) {
	a := int64(1)
	a += math.MaxInt64
	assert.True(t, a < 0)
}

func TestReflectCustomTypes(t *testing.T) {
	type A []byte
	assert.Equal(t, reflect.Slice, reflect.TypeOf(A{}).Kind())
}

type A [1]byte

func (me A) Bytes() []byte { return me[:] }

type B A

func TestTypedefExposedMethods(t *testing.T) {
	b := B{}
	A(b).Bytes()
	// b.Bytes()
}

// Ensure that appending a slice to a nil slice produces a new backing array.
func TestAppendNilBytesNewBacking(t *testing.T) {
	a := []byte{1, 2, 3}
	b := append([]byte(nil), a...)
	assert.EqualValues(t, a, b)
	b[1] = 4
	assert.NotEqual(t, a, b)
	t.Log(a)
	t.Log(b)
}

const backingArraySliceSrcLen = 100000

func BenchmarkNewBackingArrayNil(b *testing.B) {
	for range iter.N(b.N) {
		_ = append([]byte(nil), make([]byte, backingArraySliceSrcLen)...)
	}
}

func BenchmarkNewBackingArrayWithCap(b *testing.B) {
	for range iter.N(b.N) {
		_ = append(make([]byte, 0, backingArraySliceSrcLen), make([]byte, backingArraySliceSrcLen)...)
	}
}

func BenchmarkNewBackingArrayCopy(b *testing.B) {
	for range iter.N(b.N) {
		copy(make([]byte, backingArraySliceSrcLen), make([]byte, backingArraySliceSrcLen))
	}
}

func TestFormatMap(t *testing.T) {
	t.Logf("%#v", "hiya")
	t.Logf("%v", map[string]int{"hello": 42})
	t.Logf("%+v", struct {
		A string
		B int
	}{"{hello world}", 42})
}

func TestResolveBadAddress(t *testing.T) {
	_, err := net.ResolveUDPAddr("udp", "0.131.255.145:33085")
	t.Log(err)
}

type funcEqualityReceiver struct {
}

func (me *funcEqualityReceiver) Method() {}

func TestFuncEquality(t *testing.T) {
	a := func() {}
	b := func() {}
	// assert.NotEqual(t, a, b)
	assert.NotEqual(t, reflect.ValueOf(a).Pointer(), reflect.ValueOf(b).Pointer())
	// var objA funcEqualityReceiver
	// var objB funcEqualityReceiver
	// assert.NotEqual(t, objA.Method, objB.Method)
}

func TestReturnTuple(t *testing.T) {
	f := func() (int, error) {
		return 42, nil
	}
	func(...interface{}) {}(f())
}

func TestSliceLoopVariableArray(t *testing.T) {
	var b [][]byte
	c := map[[1]byte]int{{1}: 1, {2}: 2}
	for a := range c {
		b = append(b, a[:])
	}
	t.Logf("%q", b)
	b = nil
	for _, a := range [][1]byte{{1}, {2}} {
		b = append(b, a[:])
	}
	t.Logf("%q", b)
	d := [1]byte{1}
	e := d[:]
	t.Logf("%q", e)
	d[0] = 2
	t.Logf("%q", e)
}

func TestQueryEscapeNul(t *testing.T) {
	assert.EqualValues(t, "P%00%8E", url.QueryEscape("\x50\x00\x8e"))
}

func TestEmptyStructEquality(t *testing.T) {
	testEmptyStructEquality(t)
}

func TestDeferRecover(t *testing.T) {
	f := func() (ret string) {
		defer func() { recover() }()
		ret = "default"
		panic("fuck")
	}
	assert.Equal(t, "default", f())
	g := func() (ret string) {
		defer recover()
		ret = "default"
		panic("fuck")
	}
	assert.Panics(t, func() { g() })
}
