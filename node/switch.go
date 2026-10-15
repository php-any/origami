package node

import (
	"github.com/php-any/origami/data"
	"strconv"
	"strings"
)

// A nil CaseValue denotes default and retains its source position.
type SwitchCase struct {
	*Node
	CaseValue  data.GetValue
	Statements []data.GetValue
}

func (s *SwitchCase) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	var value data.GetValue
	for _, statement := range s.Statements {
		var ctl data.Control
		value, ctl = statement.GetValue(ctx)
		if ctl != nil {
			return value, ctl
		}
	}
	return value, nil
}

type SwitchStatement struct {
	*Node
	Condition   data.GetValue
	Cases       []SwitchCase
	DefaultCase []data.GetValue
}

func NewSwitchStatement(from data.From, condition data.GetValue, cases []SwitchCase, defaultCase []data.GetValue) *SwitchStatement {
	// Legacy callers supplied default separately, after all numbered cases.
	if defaultCase != nil {
		cases = append(cases, SwitchCase{Statements: defaultCase})
	}
	return &SwitchStatement{Node: NewNode(from), Condition: condition, Cases: cases, DefaultCase: defaultCase}
}

// Evaluate labels until the first match, then execute branches in source order
// without evaluating later labels. Only use default if no case matched.
func (s *SwitchStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	condition, ctl := s.Condition.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	start, fallback := -1, -1
	for i := range s.Cases {
		branch := &s.Cases[i]
		if branch.CaseValue == nil {
			fallback = i
			continue
		}
		value, ctl := branch.CaseValue.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if s.isMatch(ctx, condition, value) {
			start = i
			break
		}
	}
	if start < 0 {
		start = fallback
	}
	if start < 0 {
		return data.NewNullValue(), nil
	}
	var value data.GetValue
	for i := start; i < len(s.Cases); i++ {
		value, ctl = s.Cases[i].GetValue(ctx)
		if ctl != nil {
			if bc, ok := ctl.(data.BreakControl); ok && bc.IsBreak() {
				return value, consumeBreakLevel(bc)
			}
			if cc, ok := ctl.(data.ContinueControl); ok && cc.IsContinue() {
				return value, consumeContinueLevel(cc)
			}
			return value, ctl
		}
	}
	return value, nil
}

// Boolean conversion precedes numeric comparison. Null only equals the empty
// string among strings; numbers compare numerically without truncating floats.
func (s *SwitchStatement) isMatch(ctx data.Context, left, right data.GetValue) bool {
	left, right = unwrapValue(left), unwrapValue(right)
	_, lb := left.(*data.BoolValue)
	_, rb := right.(*data.BoolValue)
	if lb || rb {
		return switchTruthy(left) == switchTruthy(right)
	}
	if isPHPNull(left) || isPHPNull(right) {
		other := right
		if isPHPNull(right) {
			other = left
		}
		if str, ok := other.(*data.StringValue); ok {
			return str.Value == ""
		}
		return !switchTruthy(other)
	}
	li, lf, ln, lint := switchNumber(left)
	ri, rf, rn, rint := switchNumber(right)
	if ln && rn {
		if lint && rint {
			return li == ri
		}
		return lf == rf
	}
	if ls, ok := left.(*data.StringValue); ok {
		if rs, ok := right.(*data.StringValue); ok {
			return ls.Value == rs.Value
		}
		if rn {
			return ls.Value == right.(data.AsString).AsString()
		}
	}
	if rs, ok := right.(*data.StringValue); ok && ln {
		return left.(data.AsString).AsString() == rs.Value
	}
	if lts, ok := getDateTimeTimestamp(ctx, left); ok {
		if rts, ok := getDateTimeTimestamp(ctx, right); ok {
			return lts == rts
		}
	}
	return isStrictEqual(left, right)
}

func switchTruthy(value data.GetValue) bool {
	if value == nil {
		return false
	}
	if value, ok := value.(data.AsBool); ok {
		truth, err := value.AsBool()
		return err == nil && truth
	}
	return true
}

// Keep integer precision, and accept PHP decimal numeric strings with ASCII
// whitespace. Reject hex, NaN, infinity and Go digit separators.
func switchNumber(value data.GetValue) (integer int64, decimal float64, numeric, integral bool) {
	switch v := value.(type) {
	case *data.IntValue:
		return int64(v.Value), float64(v.Value), true, true
	case *data.FloatValue:
		return 0, v.Value, true, false
	case *data.StringValue:
		text := strings.Trim(v.Value, " \t\n\r\v\f")
		if text == "" {
			return
		}
		for i := 0; i < len(text); i++ {
			c := text[i]
			if (c < '0' || c > '9') && c != '+' && c != '-' && c != '.' && c != 'e' && c != 'E' {
				return
			}
		}
		if n, err := strconv.ParseInt(text, 10, 64); err == nil {
			return n, float64(n), true, true
		}
		if n, err := strconv.ParseFloat(text, 64); err == nil {
			return 0, n, true, false
		}
	}
	return
}
