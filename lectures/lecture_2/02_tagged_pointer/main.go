package main

import (
	"fmt"
	"unsafe"
)

func main() {
	value := 42
	pValue := &value

	// ok for example (uintptr)
	address := int(uintptr(unsafe.Pointer(pValue)))
	fmt.Println("Old address: ", address)
	fmt.Println("Old address value: ", *(*int)(unsafe.Pointer(uintptr(address))))

	newAddress := address ^ (1 << 60)
	//newAddress := address ^ (1 << 57)
	//newAddress := address ^ (1 << 55) // !!!
	fmt.Println("New address: ", address)
	fmt.Println("New address value: ", *(*int)(unsafe.Pointer(uintptr(newAddress))))
}
