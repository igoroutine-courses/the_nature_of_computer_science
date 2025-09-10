package main

import (
	"fmt"
	"sync"
)

// (1, 0): 9896564 times
// (0, 1): 100218 times
// (0, 0): 3217 times
// (1, 1): 1 times

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
			x = 1 // in cache
			r1 = y
		})

		wg.Go(func() {
			y = 1
			r2 = x
		})

		wg.Wait()
		results[[2]int{int(r1), int(r2)}]++

		// expected
		// 1 1
		// 0 1
		// 1 0
	}

	for k, v := range results {
		fmt.Printf("(%d, %d): %d times\n", k[0], k[1], v)
	}
}
