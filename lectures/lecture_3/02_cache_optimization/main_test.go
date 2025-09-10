package main

import (
	"testing"
)

var sink int64

func BenchmarkMatrixIteration(b *testing.B) {
	const n = 4096

	b.Run("row iteration", func(b *testing.B) {
		b.ReportAllocs()
		matrix := [n][n]int64{}
		for i := range matrix {
			for j := range matrix {
				matrix[i][j] = 1
			}
		}
		b.ResetTimer()

		for b.Loop() {
			var sum int64

			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					sum += matrix[i][j]
				}
			}

			sink += sink
		}
	})

	b.Run("column iteration", func(b *testing.B) {
		b.ReportAllocs()
		matrix := [n][n]int64{}
		for i := range matrix {
			for j := range matrix {
				matrix[i][j] = 1
			}
		}
		b.ResetTimer()

		for b.Loop() {
			var sum int64

			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					sum += matrix[j][i]
				}
			}
		}
	})
}
