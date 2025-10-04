package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {}

// go build -gcflags='-m=4' .
// ./main.go:37:3:[1] estimatePgxArgSize stmt: return len(fmt.Sprintf("%v", ... argument...))
// ./main.go:21:25: parameter arg leaks to {heap} for estimatePgxArgSize with derefs=0:
// ./main.go:21:25:   flow: {storage for ... argument} ← arg:
// ./main.go:21:25:     from ... argument (slice-literal-element) at ./main.go:37:25
// ./main.go:21:25:   flow: {heap} ← {storage for ... argument}:
// ./main.go:21:25:     from ... argument (spill) at ./main.go:37:25
// ./main.go:21:25:     from fmt.Sprintf("%v", ... argument...) (call parameter) at ./main.go:37:25
// ./main.go:21:25: leaking param: arg
// ./main.go:37:25: ... argument does not escape

//go:noinline
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
