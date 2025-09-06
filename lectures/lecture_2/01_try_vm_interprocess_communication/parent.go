package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"
	"unsafe"
)

const childPath = "open_lessons/vmemory/try_vm_interprocess_communication/child.go"

func getValidAddress() *int {
	s := make([]int, 1<<20)
	return &s[0]
}

func runChild(
	ctx context.Context,
	logger *slog.Logger,
	address *int,
) {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	childFilePath := filepath.Join(wd, childPath)

	cmd := exec.CommandContext(ctx, "go", "run", childFilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	// ok for demo, not valid
	cmd.Args = append(cmd.Args, fmt.Sprint(uintptr(unsafe.Pointer(address))))

	err = cmd.Start()

	if err != nil {
		logger.Error("child process init error", slog.Any("error", err))
		panic(err)
	}
}

func main() {
	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	address := getValidAddress()
	runChild(ctx, logger, address)

	*address = 42
	for {
		v := *address

		if v != 42 {
			fmt.Println("Replaced")
			break
		}

		fmt.Printf("Wait for replacement %v\n", address)
		time.Sleep(time.Second * 5)
	}
}
