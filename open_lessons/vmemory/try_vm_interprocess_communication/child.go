package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"unsafe"
)

func main() {
	address, err := strconv.ParseInt(os.Args[1], 10, 64)

	if err != nil {
		panic(err)
	}

	p := (*int)(unsafe.Pointer(uintptr(address)))
	for {
		*p = 111

		fmt.Printf("Try to replace %v\n", p)
		time.Sleep(time.Second * 5)
	}
}
