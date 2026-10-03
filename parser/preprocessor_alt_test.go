package parser

import (
	"strings"
	"testing"
)

func TestAltConvertRealFile(t *testing.T) {
	// Keep the compiled Blade constructs reproducible without a particular
	// machine's framework cache directory or generated filename.
	code := `<?php if(session('error')): ?>
<div><?php echo session('error'); ?></div>
<?php endif; ?>
<?php if ($__bag->has($__errorArgs[0])) : ?>
<span><?php echo $message; ?></span>
<?php endif; ?>`
	if !hasControlColon(code) {
		t.Fatal("did not recognize compiled Blade control syntax")
	}
	out := convertAltPHPSyntax("test.php", code)
	if strings.Contains(out, "if(session('error')):") {
		t.Fatal("conversion of if(...): failed on real file")
	}
	if strings.Contains(out, "if ($__bag->has($__errorArgs[0])) :") {
		t.Fatal("conversion of if (cond) : failed on real file")
	}
	if strings.Count(out, "}") != 2 {
		t.Fatal("did not close both alternate control blocks")
	}
}
