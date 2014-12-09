package curious

import "testing"

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
