package main

import (
	"fmt"
	"time"
	"unsafe"
)

//go:noescape
func estimatePgxArgSize(arg any) int {
	switch v := arg.(type) {
	case nil:
		return 0
	case int:
		return int(unsafe.Sizeof(v))
	case string:
		return len(v)
	case []int32:
		return len(v) * int(unsafe.Sizeof(0))
	case []byte:
		return len(v)
	case time.Time:
		return 8
	default:
		// code with reflection
		return len(fmt.Sprintf("%v", arg))
	}
}

// ./main.go:10:6: can only use //go:noescape with external func implementations

func main() {
	fmt.Println(estimatePgxArgSize(1))                                     // 8
	fmt.Println(estimatePgxArgSize("Hello, @igoroutine"))                  // 1
	fmt.Println(estimatePgxArgSize(map[string]any{"hello": "igoroutine"})) // 21
}
