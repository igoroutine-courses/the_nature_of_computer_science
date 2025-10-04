package main

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

//go:noinline
func estimatePgxArgSize(arg any) int {
	switch v := arg.(type) {
	case nil:
		return 0
	case int:
		return int(unsafe.Sizeof(v))
	case string:
		return len(v)
	case []int32:
		return len(v) * int(unsafe.Sizeof(0))
	case []byte:
		return len(v)
	case time.Time:
		return 8
	default:
		// code with reflection
		return len(fmt.Sprintf("%v", arg))
	}
}

func BenchmarkEstimatePgxArgSize(b *testing.B) {
	b.ReportAllocs()

	s := make([]int32, 100)
	for i := range s {
		s[i] = 27
	}

	b.ResetTimer()

	for b.Loop() {
		estimatePgxArgSize(s)
	}
}
