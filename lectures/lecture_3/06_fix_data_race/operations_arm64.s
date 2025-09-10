#include "go_asm.h"
#include "textflag.h"

// func SyncRead(addr *int64) (val int64)
TEXT ·SyncRead(SB),NOSPLIT,$0-16
	MOVD	ptr+0(FP), R0

	// Load-Acquire Register
	LDAR	(R0), R0
	MOVD	R0, ret+8(FP)
	RET

// func SyncWrite(addr *int64, val int64)
TEXT ·SyncWrite(SB), NOSPLIT, $0-16
	MOVD	ptr+0(FP), R0
	MOVD	val+8(FP), R1

	// Store-Release register
	STLR	R1, (R0)
	RET
