#include "go_asm.h"
#include "textflag.h"

// func SyncRead(addr *int64) (val int64)
TEXT ·SyncRead(SB),NOSPLIT,$0-16
	MOVQ	ptr+0(FP), AX
	MOVQ    (AX), DX
	MOVQ	DX, ret+8(FP)
	RET

// func SyncWrite(addr *int64, val int64)
TEXT ·SyncWrite(SB), NOSPLIT, $0-16
	MOVQ	ptr+0(FP), BX
	MOVQ	val+8(FP), AX

	// XCHG (Exchange) always implies a full memory barrier and atomicity with memory operands
	XCHGQ	AX, 0(BX)
	RET
