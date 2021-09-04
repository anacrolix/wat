package wat

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/bradfitz/iter"
)

type copyElem = uintptr

var sizeofCopyElem uintptr

func init() {
	var arbitraryCopyElem copyElem
	sizeofCopyElem = unsafe.Sizeof(&arbitraryCopyElem)
}

func BenchmarkCopy(b *testing.B) {
	for _, f := range []struct {
		name   string
		copyDo func(b *testing.B, count int)
	}{{
		name: "RegularAppend",
		copyDo: func(b *testing.B, size int) {
			var sl []copyElem
			for i := 0; i < size; i++ {
				sl = append(sl, uintptr(i))
			}
		},
	}, {
		name: "PreallocateAppend",
		copyDo: func(b *testing.B, size int) {
			sl := make([]copyElem, 0, size)
			for i := 0; i < size; i++ {
				sl = append(sl, uintptr(i))
			}
		},
	}, {
		name: "AssignAllocated",
		copyDo: func(b *testing.B, size int) {
			sl := make([]copyElem, size)
			for i := 0; i < size; i++ {
				sl[i] = uintptr(i)
			}
		},
	}} {
		b.Run(f.name, func(b *testing.B) {
			for _, count := range []int{0x10, 0x100, 0x1000, 0x10000, 0x100000} {
				b.Run(fmt.Sprintf("%v", count), func(b *testing.B) {
					b.SetBytes(int64(sizeofCopyElem * uintptr(count)))
					for range iter.N(b.N) {
						f.copyDo(b, count)
					}
				})
			}
		})
	}
}
