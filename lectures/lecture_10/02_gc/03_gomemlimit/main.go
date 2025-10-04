package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	const memoryLimit = 1000 << 20
	debug.SetMemoryLimit(memoryLimit)
	debug.SetGCPercent(-1)

	fmt.Println("Memory limit set to", memoryLimit>>20, "MB")

	const chunkSize = 100 << 20

	for i := 0; i < 100; i++ {
		buf := make([]byte, chunkSize)
		_ = buf

		time.Sleep(200 * time.Millisecond)

		mem := new(runtime.MemStats)
		runtime.ReadMemStats(mem)
		fmt.Printf("Allocated: %3d MB | Heap: %3d MB | NumGC: %d\n",
			chunkSize>>20,
			mem.HeapAlloc/(1024*1024),
			mem.NumGC,
		)
	}
}
