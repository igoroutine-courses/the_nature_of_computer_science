package main

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

//go:noinline
func estimatePgxArgSizeWithDefault(arg any) int {
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

//go:noinline
func estimatePgxArgSizeWithoutDefault(arg any) int {
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
	}

	panic("unexpected type")
}

func BenchmarkEstimatePgxArgSize(b *testing.B) {
	const sliceSize = 100

	b.Run("with default", func(b *testing.B) {
		b.ReportAllocs()
		s := make([]int32, sliceSize)
		for i := range s {
			s[i] = 27
		}

		b.ResetTimer()

		for b.Loop() {
			estimatePgxArgSizeWithDefault(s)
		}
	})

	b.Run("without default", func(b *testing.B) {
		b.ReportAllocs()
		s := make([]int32, sliceSize)
		for i := range s {
			s[i] = 27
		}

		b.ResetTimer()

		for b.Loop() {
			estimatePgxArgSizeWithoutDefault(s)
		}
	})
}
