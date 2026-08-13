package lexer

import (
	"testing"

	"github.com/php-any/origami/token"
)

func TestHeredocTokenType(t *testing.T) {
	input := `$a = <<<EOT
hello
EOT;
$x = 1 << 2;`
	tokens := NewLexer().Tokenize(input)
	var heredoc, shl int
	for _, tok := range tokens {
		switch tok.Type() {
		case token.HEREDOC:
			heredoc++
		case token.SHL:
			shl++
		case token.HEREDOC_START:
			t.Fatalf("不应单独出现 HEREDOC_START，完整 heredoc 应合并为 HEREDOC token")
		}
	}
	if heredoc != 1 {
		t.Fatalf("期望 1 个 HEREDOC token，实际 %d", heredoc)
	}
	if shl != 1 {
		t.Fatalf("期望 1 个 SHL(<<) token，实际 %d", shl)
	}
}

func TestNowdocTokenType(t *testing.T) {
	input := `$a = <<<'ND'
line
ND;`
	tokens := NewLexer().Tokenize(input)
	var nowdoc int
	for _, tok := range tokens {
		if tok.Type() == token.NOWDOC {
			nowdoc++
		}
	}
	if nowdoc != 1 {
		t.Fatalf("期望 1 个 NOWDOC token，实际 %d", nowdoc)
	}
}
