package internal

import (
	"fmt"
	"time"
	"unsafe"
)

//go:noinline
func EstimatePgxArgSizeInternal(arg any) int {
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
