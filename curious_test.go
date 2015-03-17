package curious

import (
	"reflect"
	"testing"
	"unsafe"
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

