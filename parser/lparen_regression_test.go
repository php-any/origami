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

// Laravel Vite::mapWithKeys 使用 (expr) => value 作数组键，不能被误判为 lambda。
func TestParseParenthesizedExprAsArrayKey(t *testing.T) {
	p := NewParser()
	src := `[(1 ? $v : $k) => $value === true ? $key : $value]`
	_, ctl := p.ParseExpressionFromString(src)
	if ctl != nil {
		t.Fatalf("parse failed: %s", ctl.AsString())
	}
}

func TestParseOrigamiLambdaStillWorks(t *testing.T) {
	p := NewParser()
	_, ctl := p.ParseExpressionFromString("($a, $b) => $a + $b")
	if ctl != nil {
		t.Fatalf("parse failed: %s", ctl.AsString())
	}
}
