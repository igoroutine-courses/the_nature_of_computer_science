package main

import "testing"

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
