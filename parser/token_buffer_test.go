package parser

import (
	"github.com/php-any/origami/data"
	"testing"
)

func TestReleaseTokenBufferPreservesSourceAndCanParseAgain(t *testing.T) {
	p := NewParser()
	program, ctl := p.ParseString("$value = 42;", "first.php")
	if ctl != nil {
		t.Fatal(ctl.AsString())
	}
	from := program.Statements[0].(interface{ GetFrom() data.From }).GetFrom()
	start, end := from.GetPosition()
	variables := p.GetVariables()
	p.ReleaseTokenBuffer()
	if len(p.tokens) != 0 || len(variables) != 1 || from.GetSource() != "first.php" {
		t.Fatal("compilation scratch was retained or AST metadata changed")
	}
	if _, ctl := p.ParseString("$other = 7;", "second.php"); ctl != nil {
		t.Fatal(ctl.AsString())
	}
	afterStart, afterEnd := from.GetPosition()
	if afterStart != start || afterEnd != end || from.GetSource() != "first.php" {
		t.Fatal("a later parse changed the released AST's source range")
	}
}
