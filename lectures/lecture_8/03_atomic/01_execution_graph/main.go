package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Observed outcomes:
// (1, 0): 590166 times
// (0, 1): 389059 times
// (0, 0): 20773 times
// (1, 1): 2 times

func main() {
	var (
		wg      sync.WaitGroup
		results = make(map[[2]int]int)
	)

	const iterations = 1_000_000

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for i := 0; i < iterations; i++ {
		var x, y int32
		var r1, r2 int32

		// P
		wg.Add(1)
		go func() {
			defer wg.Done()

			runtime.LockOSThread()
			defer runtime.UnlockOSThread()

			x = 1
			r1 = y
		}()

		// Q
		wg.Add(1)
		go func() {
			defer wg.Done()

			runtime.LockOSThread()
			defer runtime.UnlockOSThread()

			y = 1
			r2 = x
		}()

		wg.Wait()

		results[[2]int{int(r1), int(r2)}]++
	}

	fmt.Println("Observed outcomes:")
	for k, v := range results {
		fmt.Printf("(%d, %d): %d times\n", k[0], k[1], v)
	}
}
