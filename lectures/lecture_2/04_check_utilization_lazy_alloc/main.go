package main

import (
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

// vmmap
func main() {
	fmt.Println("PID:", os.Getpid())

	debug.SetGCPercent(-1)
	debug.SetMemoryLimit(-1)

	b := make([][]int, 0)

	for i := range 100_000 {
		other := make([]int, 1<<28)
		for j := range other {
			other[j] = 123
		}

		b = append(b, other)

		time.Sleep(time.Millisecond * 300)
		fmt.Println(i)
	}
}
