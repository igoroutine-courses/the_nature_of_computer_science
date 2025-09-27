package main

import (
	"fmt"
	"runtime"
	"sync"
)

// Observed outcomes:
// (1, 0): 556764 times
// (0, 1): 443227 times
// (1, 1): 9 times

//go:noescape
//go:noinline
func SyncRead(addr *int64) (val int64)

//go:noescape
//go:noinline
func SyncWrite(addr *int64, val int64)

func main() {
	var (
		wg      sync.WaitGroup
		results = make(map[[2]int]int)
	)

	const iterations = 1_000_000

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	for i := 0; i < iterations; i++ {
		var x, y int64
		var r1, r2 int64

		// P
		wg.Add(1)
		go func() {
			defer wg.Done()

			//runtime.LockOSThread()
			//defer runtime.UnlockOSThread()

			x = 1
			//SyncWrite(&x, 1)

			r1 = y
			//SyncWrite(&r1, SyncRead(&y))
		}()

		// Q
		wg.Add(1)
		go func() {
			defer wg.Done()

			//runtime.LockOSThread()
			//defer runtime.UnlockOSThread()

			y = 1
			//SyncWrite(&y, 1)
			//
			r2 = x
			//SyncWrite(&r2, SyncRead(&x))
		}()

		wg.Wait()

		results[[2]int{int(r1), int(r2)}]++
		//results[[2]int{int(SyncRead(&r1)), int(SyncRead(&r2))}]++
	}

	fmt.Println("Observed outcomes:")
	for k, v := range results {
		fmt.Printf("(%d, %d): %d times\n", k[0], k[1], v)
	}
}
