package main

import (
	"fmt"
	"sync"
)

//go:noescape
//go:noinline
func SyncRead(addr *int64) (val int64)

//go:noescape
//go:noinline
func SyncWrite(addr *int64, val int64)

// (1, 1): 6 times
// (1, 0): 9901067 times
// (0, 1): 98927 times

func main() {
	wg := new(sync.WaitGroup)
	results := make(map[[2]int]int)

	for i := 0; i < 10_000_000; i++ {
		var (
			x  int64 = 0
			y  int64 = 0
			r1 int64 = 0
			r2 int64 = 0
		)

		wg.Go(func() {
			SyncWrite(&x, 1)             // x = 1
			SyncWrite(&r1, SyncRead(&y)) // r1 = y
		})

		wg.Go(func() {
			SyncWrite(&y, 1)             // y = 1
			SyncWrite(&r2, SyncRead(&x)) // r2 = x
		})

		wg.Wait()
		results[[2]int{int(SyncRead(&r1)), int(SyncRead(&r2))}]++
	}

	for k, v := range results {
		fmt.Printf("(%d, %d): %d times\n", k[0], k[1], v)
	}
}
