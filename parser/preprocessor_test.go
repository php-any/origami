package parser

import (
	"testing"
)

func TestConvertAltPHPSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "if colon to brace",
			input:    "<?php if($condition): ?>\nok\n<?php endif; ?>",
			expected: "<?php if ($condition) { ?>\nok\n<?php } ?>",
		},
		{
			name:     "if else colon to brace",
			input:    "<?php if($a): ?>a<?php else: ?>b<?php endif; ?>",
			expected: "<?php if ($a) { ?>a<?php } else { ?>b<?php } ?>",
		},
		{
			name:     "elseif colon",
			input:    "<?php if($a): ?>a<?php elseif($b): ?>b<?php endif; ?>",
			expected: "<?php if ($a) { ?>a<?php } elseif ($b) { ?>b<?php } ?>",
		},
		{
			name:     "foreach colon",
			input:    "<?php foreach($items as $item): ?>\n<?php echo $item; ?>\n<?php endforeach; ?>",
			expected: "<?php foreach ($items as $item) { ?>\n<?php echo $item; ?>\n<?php } ?>",
		},
		{
			name:     "endif in string skipped",
			input:    "<?php return '<?php endif; ?>'; ?>",
			expected: "<?php return '<?php endif; ?>'; ?>",
		},
		{
			name:     "else colon in string skipped",
			input:    "<?php $x = 'else:'; echo $x; ?>",
			expected: "<?php $x = 'else:'; echo $x; ?>",
		},
		{
			name:     "regular php no change",
			input:    "<?php $a = $b ? $c : $d; ?>",
			expected: "<?php $a = $b ? $c : $d; ?>",
		},
		{
			name:     "elvis operator preserved",
			input:    "<?php $x = $a ?: $b; ?>",
			expected: "<?php $x = $a ?: $b; ?>",
		},
		{
			name:     "nested if colon",
			input:    "<?php if($a): ?><?php if($b): ?>x<?php endif; ?><?php endif; ?>",
			expected: "<?php if ($a) { ?><?php if ($b) { ?>x<?php } ?><?php } ?>",
		},
		{
			name:     "if with nested parens",
			input:    "<?php if(Route::has('login')): ?>\nok\n<?php endif; ?>",
			expected: "<?php if (Route::has('login')) { ?>\nok\n<?php } ?>",
		},
		{
			name:     "if with method chain in parens",
			input:    "<?php if(auth()->guard()->check()): ?>auth<?php else: ?>guest<?php endif; ?>",
			expected: "<?php if (auth()->guard()->check()) { ?>auth<?php } else { ?>guest<?php } ?>",
		},
		{
			name:     "endif in block comment skipped",
			input:    "<?php /* endif; */ echo \"hi\"; ?>",
			expected: "<?php /* endif; */ echo \"hi\"; ?>",
		},
		{
			name:     "else line comment skipped",
			input:    "<?php\n// else: not converted\n$a = 1; ?>",
			expected: "<?php\n// else: not converted\n$a = 1; ?>",
		},
		{
			name:     "blade endauth directive",
			input:    "<?php if(auth()->guard()->check()): ?>auth<?php endif; ?>@endauth",
			expected: "<?php if (auth()->guard()->check()) { ?>auth<?php } ?><?php } ?>",
		},
		{
			name:     "while colon",
			input:    "<?php while($c): ?><?php echo $c; ?><?php endwhile; ?>",
			expected: "<?php while ($c) { ?><?php echo $c; ?><?php } ?>",
		},
		{
			name:     "for colon",
			input:    "<?php for($i=0;$i<10;$i++): ?><?php echo $i; ?><?php endfor; ?>",
			expected: "<?php for ($i=0;$i<10;$i++) { ?><?php echo $i; ?><?php } ?>",
		},
		{
			name:     "switch colon",
			input:    "<?php switch($a): ?><?php case 1: ?>x<?php endswitch; ?>",
			expected: "<?php switch ($a) { ?><?php case 1: ?>x<?php } ?>",
		},
		// ---- Blade 编译产物（混合风格）场景 ----
		// Laravel Blade 编译器会把 @if/@elseif 编译成带冒号的开始标记，
		// 但把 @else/@endif 编译成花括号（<?php } else { ?> / <?php } ?>），
		// 因此产物中只有开始标记（if(...):）而没有 endif; 结束标记。
		// 快速检查必须能识别这种"只有开始标记"的情况，否则 if(...): 不会转换。
		{
			name:     "blade mixed only if colon no endif",
			input:    "<?php if($x): ?><p>A</p><?php } ?>",
			expected: "<?php if ($x) { ?><p>A</p><?php } ?>",
		},
		{
			name:     "blade mixed if elseif else no end markers",
			input:    "<?php if($a): ?>a<?php elseif($b): ?>b<?php } else { ?>c<?php } ?>",
			expected: "<?php if ($a) { ?>a<?php } elseif ($b) { ?>b<?php } else { ?>c<?php } ?>",
		},
		{
			name:     "blade mixed nested if with brace blocks",
			input:    "<?php if(Route::has('login')): ?><?php if (auth()->guard()->check()) { ?>D<?php } else { ?>L<?php } ?><?php } ?>",
			expected: "<?php if (Route::has('login')) { ?><?php if (auth()->guard()->check()) { ?>D<?php } else { ?>L<?php } ?><?php } ?>",
		},
		{
			name:     "blade mixed with nested parens in condition",
			input:    "<?php if(file_exists(public_path('a')) || file_exists(public_path('b'))): ?>v<?php } ?>",
			expected: "<?php if (file_exists(public_path('a')) || file_exists(public_path('b'))) { ?>v<?php } ?>",
		},
		{
			name:     "blade foreach colon only",
			input:    "<?php foreach($items as $i): ?><li>X</li><?php } ?>",
			expected: "<?php foreach ($items as $i) { ?><li>X</li><?php } ?>",
		},
		{
			name:     "regular brace block no change",
			input:    "<?php if ($x) { echo 'a'; } else { echo 'b'; } ?>",
			expected: "<?php if ($x) { echo 'a'; } else { echo 'b'; } ?>",
		},
		{
			name:     "ternary operator not misdetected",
			input:    "<?php $x = $a ? $b : $c; echo $x; ?>",
			expected: "<?php $x = $a ? $b : $c; echo $x; ?>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertAltPHPSyntax("test.php", tt.input)
			if got != tt.expected {
				t.Errorf("got:  %q\nwant: %q", got, tt.expected)
			}
		})
	}
}
