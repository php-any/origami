package parser

import (
	"strings"
	"testing"
)

func TestDiagVendorCompileAuthStringCorrupt(t *testing.T) {
	in := `<?php
protected function compileAuth($guard = null) {
    $guard = is_null($guard) ? '()' : $guard;
    return "<?php if(auth()->guard{$guard}->check()): ?>";
}`
	got := convertAltPHPSyntax("vendor/fake.php", in)
	// 字符串字面量 "<?php if(auth()->guard{$guard}->check()): ?>" 中的冒号形式必须保留，
	// 不能被 convertControlKeywords 误转换成花括号。
	if strings.Contains(got, "guard{$guard}->check()): ?>") &&
		!strings.Contains(got, "guard{$guard}->check()) { ?>") {
		// 字符串内的 if(...): 保留，转换正确
		return
	}
	t.Fatalf("convertControlKeywords corrupted double-quoted return string:\n%s", got)
}
