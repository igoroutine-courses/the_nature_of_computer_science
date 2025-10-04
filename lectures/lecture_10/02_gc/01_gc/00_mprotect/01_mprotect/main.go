package main

import (
	"fmt"
	"syscall"
)

func main() {
	pageSize := syscall.Getpagesize()

	data, err := syscall.Mmap(
		-1,
		0,
		pageSize,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_ANON|syscall.MAP_PRIVATE,
	)

	if err != nil {
		panic(err)
	}
	defer syscall.Munmap(data)
	copy(data, "hello")

	err = syscall.Mprotect(data, syscall.PROT_READ)
	if err != nil {
		panic(err)
	}

	// Read ok
	fmt.Println(string(data[:5]))

	// Write - SIGSEGV
	// data[0] = 'H'

	// disable protection
	err = syscall.Mprotect(data, syscall.PROT_READ|syscall.PROT_WRITE)
	if err != nil {
		panic(err)
	}

	data[0] = 'H'
	fmt.Println(string(data[:5]))
}
