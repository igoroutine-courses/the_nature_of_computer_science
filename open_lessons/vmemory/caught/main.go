package main

import (
	"fmt"
	"runtime/debug"
	"unsafe"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("caught")
		}
	}()

	debug.SetPanicOnFault(true)

	a := (*int)(unsafe.Pointer(uintptr(1234)))
	fmt.Println(*a)
}
