#include "textflag.h"

// G 返回当前 goroutine 的 runtime.g。
// 必须用 get_tls + g(r) 两步；MOVQ (TLS) 在 windows 上只是 TLS 基址。
// 本包不能叫 runtime，否则汇编符号 runtime·tls_g 会链到错误的包。
TEXT ·G(SB), NOSPLIT, $0-8
	MOVQ TLS, AX
	MOVQ 0(AX)(TLS*1), AX
	MOVQ AX, ret+0(FP)
	RET
