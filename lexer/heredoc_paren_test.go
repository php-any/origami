package lexer

import (
	"testing"

	"github.com/php-any/origami/token"
)

func TestHeredocAfterParen(t *testing.T) {
	input := "(<<<HTML\ntest\nHTML)"
	tokens := NewLexer().Tokenize(input)
	for i, tok := range tokens {
		lit := tok.Literal()
		if len(lit) > 30 {
			lit = lit[:30] + "..."
		}
		t.Logf("%d: type=%v lit=%q", i, tok.Type(), lit)
	}
	found := false
	for _, tok := range tokens {
		if tok.Type() == token.HEREDOC {
			found = true
		}
	}
	if !found {
		t.Fatal("expected HEREDOC token after (")
	}
}

func TestHeredocBeforeCloseParen(t *testing.T) {
	input := "foo($base.<<<HTML\ntest\nHTML)"
	tokens := NewLexer().Tokenize(input)
	found := false
	for _, tok := range tokens {
		if tok.Type() == token.HEREDOC {
			found = true
		}
	}
	if !found {
		t.Fatal("expected HEREDOC token before closing paren in call")
	}
}

func TestHeredocAfterDot(t *testing.T) {
	input := "$base.<<<HTML\ntest\nHTML"
	tokens := NewLexer().Tokenize(input)
	found := false
	for _, tok := range tokens {
		if tok.Type() == token.HEREDOC {
			found = true
		}
	}
	if !found {
		t.Fatal("expected HEREDOC token after dot")
	}
}
