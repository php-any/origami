package parser

import (
	"testing"

	"github.com/php-any/origami/lexer"
)

func TestReproClassNamedArg(t *testing.T) {
	src := `<?php
class Finder {
    public function addComponent($name, $class) {}
}
class Foo {
    public function test() {
        $finder = new Finder();
        $name = "x";
        $class = "y";
        $finder->addComponent(name: $name, class: $class);
    }
}
`
	lx := lexer.NewLexer()
	toks := lx.Tokenize(src)
	for i, tok := range toks {
		t.Logf("tok[%d] type=%d literal=%q line=%d pos=%d", i, tok.Type(), tok.Literal(), tok.Line(), tok.Pos())
	}

	p := NewParser()
	_, ctl := p.ParseString(src, "repro.php")
	if ctl != nil {
		t.Fatalf("parse error: %v", ctl)
	}
}
