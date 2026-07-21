package parser

import "testing"

func TestParseParenthesizedSubtractionExpression(t *testing.T) {
	p := NewParser()
	_, ctl := p.ParseExpressionFromString("($i-1)")
	if ctl != nil {
		t.Fatalf("parse failed: %s", ctl.AsString())
	}
}

func TestParseParenthesizedNullCoalesceArrowFn(t *testing.T) {
	p := NewParser()
	_, ctl := p.ParseExpressionFromString("(null ?? fn ($scale) => $scale <= 1000)(500)")
	if ctl != nil {
		t.Fatalf("parse failed: %s", ctl.AsString())
	}
}
