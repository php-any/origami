package array

import (
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// RangeFunction 实现 range(start, end, step=1): array
//
// PHP 语义要点：
//   - range("1","3") → [1,2,3]（数字字符串走数值路径）
//   - range("a","c") → ["a","b","c"]（非数字单字符走字符路径）
type RangeFunction struct{}

func NewRangeFunction() data.FuncStmt {
	return &RangeFunction{}
}

func (f *RangeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	startVal, _ := ctx.GetIndexValue(0)
	endVal, _ := ctx.GetIndexValue(1)
	stepVal, _ := ctx.GetIndexValue(2)

	if startVal == nil || endVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 字符串单字符且非数字：range('a','c')
	if isSingleCharString(startVal) && isSingleCharString(endVal) &&
		!isNumericString(startVal) && !isNumericString(endVal) {
		start := int([]rune(startVal.AsString())[0])
		end := int([]rune(endVal.AsString())[0])
		step := 1
		if stepVal != nil {
			if iv, ok := stepVal.(data.AsInt); ok {
				if v, err := iv.AsInt(); err == nil && v != 0 {
					step = v
				}
			}
		}
		return data.NewArrayValue(rangeIntValues(start, end, step, true)), nil
	}

	startF, startIsFloat := asFloat(startVal)
	endF, endIsFloat := asFloat(endVal)
	stepF := 1.0
	stepIsFloat := false
	if stepVal != nil {
		if f64, isF, ok := asFloat3(stepVal); ok {
			if f64 != 0 {
				stepF = f64
				stepIsFloat = isF
			}
		}
	}

	useFloat := startIsFloat || endIsFloat || stepIsFloat
	if useFloat {
		return data.NewArrayValue(rangeFloatValues(startF, endF, stepF)), nil
	}

	startI := int(startF)
	endI := int(endF)
	stepI := int(stepF)
	if stepI == 0 {
		stepI = 1
	}
	return data.NewArrayValue(rangeIntValues(startI, endI, stepI, false)), nil
}

func isSingleCharString(v data.Value) bool {
	_, ok := v.(*data.StringValue)
	if !ok {
		return false
	}
	r := []rune(v.AsString())
	return len(r) == 1
}

func isNumericString(v data.Value) bool {
	sv, ok := v.(*data.StringValue)
	if !ok {
		return false
	}
	s := strings.TrimSpace(sv.AsString())
	if s == "" {
		return false
	}
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func asFloat3(v data.Value) (float64, bool, bool) {
	switch t := v.(type) {
	case *data.FloatValue:
		return t.Value, true, true
	case *data.IntValue:
		return float64(t.Value), false, true
	case *data.StringValue:
		s := strings.TrimSpace(t.AsString())
		if s == "" {
			return 0, false, false
		}
		// 纯整数字符串保持 int 路径（range("10","12") → int[]）
		if !strings.ContainsAny(s, ".eE") {
			if i, err := strconv.ParseInt(s, 10, 64); err == nil {
				return float64(i), false, true
			}
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return f, true, true
		}
		return 0, false, false
	case data.AsFloat:
		if f, err := t.AsFloat(); err == nil {
			return f, true, true
		}
	case data.AsInt:
		if i, err := t.AsInt(); err == nil {
			return float64(i), false, true
		}
	}
	return 0, false, false
}

func asFloat(v data.Value) (float64, bool) {
	f, isF, ok := asFloat3(v)
	if !ok {
		return 0, false
	}
	return f, isF
}

func rangeIntValues(start, end, step int, asChar bool) []data.Value {
	if step == 0 {
		step = 1
	}
	if start <= end && step < 0 {
		step = -step
	}
	if start > end && step > 0 {
		step = -step
	}

	out := make([]data.Value, 0)
	if start <= end {
		for i := start; i <= end; i += step {
			if asChar {
				out = append(out, data.NewStringValue(string(rune(i))))
			} else {
				out = append(out, data.NewIntValue(i))
			}
			if step == 0 {
				break
			}
		}
	} else {
		for i := start; i >= end; i += step {
			if asChar {
				out = append(out, data.NewStringValue(string(rune(i))))
			} else {
				out = append(out, data.NewIntValue(i))
			}
			if step == 0 {
				break
			}
		}
	}
	return out
}

func rangeFloatValues(start, end, step float64) []data.Value {
	if step == 0 {
		step = 1
	}
	if start <= end && step < 0 {
		step = -step
	}
	if start > end && step > 0 {
		step = -step
	}

	out := make([]data.Value, 0)
	const eps = 1e-12
	if start <= end {
		for x := start; x <= end+eps; x += step {
			out = append(out, data.NewFloatValue(x))
			if step == 0 {
				break
			}
		}
	} else {
		for x := start; x >= end-eps; x += step {
			out = append(out, data.NewFloatValue(x))
			if step == 0 {
				break
			}
		}
	}
	return out
}

func (f *RangeFunction) GetName() string { return "range" }

func (f *RangeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "start", 0, nil, nil),
		node.NewParameter(nil, "end", 1, nil, nil),
		node.NewParameter(nil, "step", 2, node.NewIntLiteral(nil, "1"), nil),
	}
}

func (f *RangeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "start", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "end", 1, data.NewBaseType("mixed")),
		node.NewVariable(nil, "step", 2, data.NewBaseType("mixed")),
	}
}
