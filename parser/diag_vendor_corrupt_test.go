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
	want := `): ?>`
	if !strings.Contains(got, want) {
		t.Fatalf("expected alt syntax ): ?> in return string, got:\n%s", got)
	}
	if strings.Contains(got, ") {") && strings.Contains(got, "return \"<?php if") {
		t.Fatalf("convertControlKeywords corrupted double-quoted return string:\n%s", got)
	}
}
