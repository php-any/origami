package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/php-any/origami/token"
)

func TestPharStubHaltOffset(t *testing.T) {
	for _, ending := range []string{"\n", "\r\n"} {
		for _, name := range []string{"__halt_compiler", "__HALT_COMPILER", "__HaLt_CoMpIlEr"} {
			stub := "#!/usr/bin/env php" + ending + "<?php /* 中文 */ " +
				"declare(strict_types=1); " + name + " /* comment */ (" + ending + ");"
			path := filepath.Join(t.TempDir(), "archive.phar")
			if err := os.WriteFile(path, []byte(stub+"\x00\xff<?php declare(strict_types=1);"), 0600); err != nil {
				t.Fatal(err)
			}
			p := NewParser()
			program, ctl := p.ParseFile(path)
			if ctl != nil {
				t.Fatalf("%s: %s", name, ctl.AsString())
			}
			if !program.StrictTypes || p.compilerHaltOffset != len(stub) {
				t.Fatalf("strict=%v offset=%d want=%d", program.StrictTypes, p.compilerHaltOffset, len(stub))
			}
		}
	}
}

func TestPharShebangNotEmitted(t *testing.T) {
	for _, tok := range tokenizeFileContent(NewParser().lexer, "archive.phar", []byte("#!/usr/bin/env php\n<?php echo 'ok';")) {
		if tok.Type() == token.HTML_TAG {
			t.Fatalf("shebang emitted as output: %q", tok.Literal())
		}
	}
}
