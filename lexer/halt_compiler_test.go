package lexer

import (
	"strings"
	"testing"

	"github.com/php-any/origami/token"
)

func TestCompilerHaltCaseAndOffset(t *testing.T) {
	for _, name := range []string{"__halt_compiler", "__HALT_COMPILER", "__HaLt_CoMpIlEr"} {
		t.Run(name, func(t *testing.T) {
			stub := "#!/usr/bin/env php\n<?php /* 中文 */ " + name + " /* comment */ (\n) ;"
			tokens := NewLexer().TokenizeTemplate(stub + "\x00\xff<?php declare(strict_types=1);")
			found := false
			for _, tok := range tokens {
				if tok.Type() == token.HALT_COMPILER {
					found = true
				}
				if tok.End() > len(stub) {
					t.Fatalf("scanned archive data: %v", tok)
				}
			}
			if !found || tokens[len(tokens)-1].End() != len(stub) {
				t.Fatalf("missing halt or incorrect byte offset: %v", tokens)
			}
		})
	}
}

func TestCompilerHaltIdentifierBoundary(t *testing.T) {
	for _, name := range []string{"__HALT_COMPILER_extra", "__HALT_COMPILER2", "__HALT_COMPILER中文"} {
		for _, tok := range NewLexer().TokenizeTemplate("<?php " + name + "();") {
			if tok.Type() == token.HALT_COMPILER {
				t.Fatalf("identifier %q treated as halt", name)
			}
		}
	}
}

func TestCompilerHaltInTextDoesNotStop(t *testing.T) {
	source := "<?php /* __HALT_COMPILER(); */ $__text = '__HALT_COMPILER();'; echo 'after';"
	tokens := NewLexer().TokenizeTemplate(source)
	for _, tok := range tokens {
		if tok.Type() == token.HALT_COMPILER {
			t.Fatal("text treated as halt")
		}
	}
	if tokens[len(tokens)-1].End() != len(strings.TrimSpace(source)) {
		t.Fatal("stopped before end of source")
	}
}
