package lexer

import (
	"strings"
	"testing"

	"github.com/php-any/origami/token"
)

func TestLineCommentStopsAtPhpCloseTag(t *testing.T) {
	input := "<?php // note.?>\n<span>after</span>"
	tokens := NewLexer().TokenizeTemplate(input)
	foundHTML := false
	for _, tok := range tokens {
		if tok.Type() == token.HTML_TAG && strings.Contains(tok.Literal(), "<span>after</span>") {
			foundHTML = true
			break
		}
	}
	if !foundHTML {
		var dump []string
		for _, tok := range tokens {
			dump = append(dump, tok.Literal())
		}
		t.Fatalf("?> 应切回 HTML，token=%q", dump)
	}
}
