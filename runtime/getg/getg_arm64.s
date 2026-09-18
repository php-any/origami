#include "textflag.h"

TEXT ·G(SB), NOSPLIT, $0-8
	MOVD g, ret+0(FP)
	RET
