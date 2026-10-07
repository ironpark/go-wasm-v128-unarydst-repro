#include "textflag.h"

TEXT ·store16(SB), NOSPLIT, $0-16
	I64Load dst+0(FP)
	I32WrapI64
	I64Load src+8(FP)
	I32WrapI64
	V128Load $0
	V128Store $16
	RET

TEXT ·lane3(SB), NOSPLIT, $0-12
	Get SP
	I64Load p+0(FP)
	I32WrapI64
	V128Load $0
	I32x4ExtractLane $3
	I32Store ret+8(FP)
	RET

TEXT ·load16(SB), NOSPLIT, $0-16
	I64Load dst+0(FP)
	I32WrapI64
	I64Load src+8(FP)
	I32WrapI64
	V128Load $16
	V128Store $0
	RET
