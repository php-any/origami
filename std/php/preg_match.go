package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/std/php/preg"
)

// PregMatchFunction 实现 preg_match 函数
type PregMatchFunction struct{}

func NewPregMatchFunction() data.FuncStmt {
	return &PregMatchFunction{}
}

func (f *PregMatchFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	patternValue, _ := ctx.GetIndexValue(0)
	subjectValue, _ := ctx.GetIndexValue(1)
	matchesValue, _ := ctx.GetIndexValue(2)
	// flagsValue, _ := ctx.GetIndexValue(3)
	// offsetValue, _ := ctx.GetIndexValue(4)

	if patternValue == nil || subjectValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()

	// 使用 preg.CompileAny 统一处理 PHP 风格的正则表达式（支持 lookahead/lookbehind）
	re, err := preg.CompileAny(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	// Find matches
	// preg_match finds the first match.
	loc := re.FindStringSubmatchIndex(subject)
	if loc == nil {
		// 无匹配时，清空 $matches
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = data.NewArrayValue([]data.Value{})
		}
		return data.NewIntValue(0), nil // No match
	}

	// 如果传入了第三个参数，填充匹配结果
	{
		matchStrs := []data.Value{}
		// loc contains [start, end, start, end...]
		for i := 0; i < len(loc); i += 2 {
			start, end := loc[i], loc[i+1]
			if start == -1 {
				matchStrs = append(matchStrs, data.NewStringValue("")) // Unmatched group?
			} else {
				matchStrs = append(matchStrs, data.NewStringValue(subject[start:end]))
			}
		}
		newMatches := data.NewArrayValue(matchStrs)
		// 通过 ZVal 引用写回（与 preg_match_all 保持一致）
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = newMatches
		} else if matchesValue != nil {
			if r, ok := matchesValue.(*data.ReferenceValue); ok {
				r.Ctx.SetVariableValue(r.Val, newMatches)
			}
		}
	}

	return data.NewIntValue(1), nil // Match found
}

func (f *PregMatchFunction) GetName() string {
	return "preg_match"
}

func (f *PregMatchFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "pattern", 0, nil, nil),
		node.NewParameter(nil, "subject", 1, nil, nil),
		node.NewParameterReference(nil, "matches", 2, data.NewBaseType("array")),
		node.NewParameter(nil, "flags", 3, node.NewIntLiteral(nil, "0"), nil),
		node.NewParameter(nil, "offset", 4, node.NewIntLiteral(nil, "0"), nil),
	}
}

func (f *PregMatchFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "pattern", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "subject", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "matches", 2, data.NewBaseType("array")), // Should be reference?
		node.NewVariable(nil, "flags", 3, data.NewBaseType("int")),
		node.NewVariable(nil, "offset", 4, data.NewBaseType("int")),
	}
}
