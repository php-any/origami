package parser

import "testing"

func TestParseParenthesizedSubtractionExpression(t *testing.T) {
	p := NewParser()
	_, ctl := p.ParseExpressionFromString("($i-1)")
	if ctl != nil {
		t.Fatalf("parse failed: %s", ctl.AsString())
	}
}
