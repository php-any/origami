package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func TestStrictTypesPHPFile(t *testing.T) {
	for _, test := range []struct {
		source string
		fatal  bool
	}{
		{"<?php\ndeclare(strict_types=1);\n", false},
		{"<?php\r\n/* comment */\r\ndeclare(strict_types=1);\n", false},
		{"<?php\n; declare(strict_types=1);", true},
		{"\n<?php declare(strict_types=1);", true},
		{"<?= 1; ?><?php declare(strict_types=1);", true},
	} {
		t.Run(test.source, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "strict.php")
			if err := os.WriteFile(path, []byte(test.source), 0600); err != nil {
				t.Fatal(err)
			}
			program, ctl := NewParser().ParseFile(path)
			if test.fatal {
				if failure, ok := ctl.(*data.ThrowValue); !ok || !failure.PHPCompileFatal {
					t.Fatalf("expected compile fatal: %v", ctl)
				}
				return
			}
			if ctl != nil {
				t.Fatal(ctl.AsString())
			}
			if !program.StrictTypes {
				t.Fatal("strict mode lost in PHP file")
			}
		})
	}
}

func TestStrictTypesCompilationUnit(t *testing.T) {
	cases := []struct {
		source        string
		strict, fatal bool
	}{
		{"declare(strict_types=1);", true, false},
		{"declare(strict_types=0);", false, false},
		{"/* comment */ declare(strict_types=0x1);", true, false},
		{"declare(strict_types=0b1);", true, false},
		{"declare(ticks=1); declare(strict_types=1);", true, false},
		{"declare(strict_types=1); declare(strict_types=0);", false, false},
		{"declare(strict_types=2);", false, true},
		{"declare(strict_types=1.0);", false, true},
		{"declare(strict_types='1');", false, true},
		{"echo ''; declare(strict_types=1);", false, true},
		{"function f() { declare(strict_types=1); }", false, true},
		{"declare(strict_types=1) {}", false, true},
	}
	for _, test := range cases {
		t.Run(test.source, func(t *testing.T) {
			parser := NewParser()
			program, ctl := parser.ParseString(test.source, "strict.php")
			if test.fatal {
				failure, ok := ctl.(*data.ThrowValue)
				if !ok || !failure.PHPCompileFatal {
					t.Fatalf("expected compile fatal, got %v", ctl)
				}
				return
			}
			if ctl != nil {
				t.Fatal(ctl.AsString())
			}
			if program.StrictTypes != test.strict {
				t.Fatalf("strict = %v", program.StrictTypes)
			}
		})
	}
}

func TestStrictTypesDeclarationMetadata(t *testing.T) {
	parser := NewParser()
	program, ctl := parser.ParseString("declare(strict_types=1); function declared(): string { return 'ok'; }", "strict.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	declaration, ok := program.Statements[len(program.Statements)-1].(*node.FunctionStatement)
	if !ok || !declaration.StrictTypes {
		t.Fatalf("strict declaration metadata lost: %T", declaration)
	}
	// A parser can be reused for a different weak compilation unit.
	program, ctl = parser.ParseString("function weak(): string { return 'ok'; }", "weak.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	if program.StrictTypes {
		t.Fatal("strict mode leaked between files")
	}
}
