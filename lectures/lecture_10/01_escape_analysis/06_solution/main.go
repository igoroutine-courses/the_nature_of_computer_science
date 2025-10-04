package main

import (
	_ "escape/06_solution/internal" // !!!
	_ "unsafe"                      // !!!
)

//go:noescape
//go:linkname estimatePgxArgSize escape/06_solution/internal.EstimatePgxArgSizeInternal
func estimatePgxArgSize(arg any) int

// go build -gcflags='-m=10' main.go
// ./main.go:10:25: arg does not escape | escAnalyze
// .   DCLFUNC esc(no) main.main ABI:ABIInternal ABIRefs:{ABIInternal} InlinabilityChecked FUNC-func() tc(1) # main.go:12:6
// .   DCLFUNC-body
// .   .   BLOCK tc(1) # main.go:12:6

func main() {}
