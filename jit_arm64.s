//go:build arm64

#include "funcdata.h"
#include "textflag.h"

// See the comments on EmitGoEntryPreamble for what this function is supposed to do.
TEXT ·entrypoint(SB), NOSPLIT|NOFRAME, $0-48
	MOVD fn+0(FP), R27
	MOVD a+8(FP), R0
	MOVD b+16(FP), R1
	JMP  (R27)

TEXT ·afterGoFunctionCallEntrypoint(SB), NOSPLIT|NOFRAME, $0-32
	MOVD RSP, R27    // Move SP to R27 (temporary register) since SP cannot be stored directly in str instructions.
	MOVD R27, 24(R0) // Store R27 into [RO, #ExecutionContextOffsets.OriginalFramePointer]
	MOVD R30, 32(R0) // Store R30 into [R0, #ExecutionContextOffsets.GoReturnAddress]

	// Load the new stack pointer (which sits somewhere in Go-allocated stack) into SP.
	MOVD R19, RSP
	JMP  (R20)
