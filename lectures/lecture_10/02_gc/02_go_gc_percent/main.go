package main

import (
	"fmt"
	"runtime"
	"runtime/debug"
	"time"
)

func main() {
	debug.SetGCPercent(200) // -1, 0, 10, 100, 200

	const chunkSize = 20 << 20

	var allocated [][]byte

	for i := 0; i < 100; i++ {
		buf := make([]byte, chunkSize)
		allocated = append(allocated, buf)
		time.Sleep(200 * time.Millisecond)

		mem := new(runtime.MemStats)
		runtime.ReadMemStats(mem)

		fmt.Printf("Step: %3d | HeapAlloc: %4d MB | NumGC: %3d\n",
			i,
			mem.HeapAlloc/(1024*1024),
			mem.NumGC,
		)
	}
}
