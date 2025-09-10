package main

import (
	"runtime"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func BenchmarkCacheContention(b *testing.B) {
	const iterations = 100
	workers := runtime.GOMAXPROCS(0)

	b.Run("no sync", func(b *testing.B) {
		b.ReportAllocs()

		g := new(errgroup.Group)
		g.SetLimit(workers)

		var value int

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				g.Go(func() error {
					for range iterations {
						value++
					}

					return nil
				})
			}
		})

		err := g.Wait()
		require.NoError(b, err)
	})

	b.Run("sync", func(b *testing.B) {
		b.ReportAllocs()

		g := new(errgroup.Group)
		g.SetLimit(workers)
		value := atomic.Int64{}

		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				g.Go(func() error {
					for range iterations {
						value.Add(1)
					}

					return nil
				})
			}
		})

		err := g.Wait()
		require.NoError(b, err)
	})
}
