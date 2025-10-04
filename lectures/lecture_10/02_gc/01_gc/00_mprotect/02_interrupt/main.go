package main

import (
	"fmt"
	"sync"
	"syscall"
	"time"
)

/*
#include <signal.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

void sigbus_handler(int sig) {
    write(2, "Caught SIGBUS in C handler\n", 31);
}

void setup_sigbus() {
    struct sigaction sa;
    memset(&sa, 0, sizeof(sa));
    sa.sa_handler = sigbus_handler;
    sigaction(SIGBUS, &sa, NULL);
}
*/
import "C"

func main() {
	C.setup_sigbus()

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

	wg := new(sync.WaitGroup)
	stop := make(chan struct{})

	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("[goroutine %d] recovered from panic: %v\n", i, r)
				}
			}()

			for {
				select {
				case <-stop:
					return
				default:
					s := string(data[:5])
					fmt.Printf("[goroutine %d] read: %s\n", i, s)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()
	}

	time.Sleep(1 * time.Second)

	fmt.Println("\n[main] disabling memory access...\n")
	err = syscall.Mprotect(data, syscall.PROT_NONE)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}
