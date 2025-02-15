package main

import (
	"fmt"
	"unsafe"
)

func main() {
	value := 42
	pValue := &value

	uintPtrValue := int(uintptr(unsafe.Pointer(pValue)))
	fmt.Println("OLD: ", uintPtrValue)

	// NEW:  1152922878996868816
	newAddress := uintptr(uintPtrValue | (uintPtrValue ^ (uintPtrValue & (1 << 60))))
	fmt.Println("NEW: ", newAddress)

	fmt.Println(*(*int)(unsafe.Pointer(newAddress)))
}
