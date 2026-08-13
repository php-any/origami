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
	flagsValue, _ := ctx.GetIndexValue(3)
	offsetValue, _ := ctx.GetIndexValue(4)

	if patternValue == nil || subjectValue == nil {
		return data.NewBoolValue(false), nil
	}

	pattern := patternValue.AsString()
	subject := subjectValue.AsString()
	offset := 0
	if offsetValue != nil {
		if asInt, ok := offsetValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				offset = v
			}
		}
	}
	flags := 0
	if flagsValue != nil {
		if asInt, ok := flagsValue.(data.AsInt); ok {
			if v, err := asInt.AsInt(); err == nil {
				flags = v
			}
		}
	}

	// 使用 preg.CompileAny 统一处理 PHP 风格的正则表达式（支持 lookahead/lookbehind）
	re, err := preg.CompileAny(pattern)
	if err != nil {
		// PHP 行为: 发出 warning，返回 false；这里只返回 false
		return data.NewBoolValue(false), nil
	}

	anchored := preg.HasModifier(pattern, 'A')
	captures := preg.FindCaptures(re, subject, offset, anchored)
	if captures == nil {
		// 无匹配时，清空 $matches
		if z := ctx.GetIndexZVal(2); z != nil {
			z.Value = data.NewArrayValue([]data.Value{})
		}
		return data.NewIntValue(0), nil // No match
	}

	// 如果传入了第三个参数，填充匹配结果
	newMatches := preg.BuildMatchArray(captures, flags)
	if z := ctx.GetIndexZVal(2); z != nil {
		z.Value = newMatches
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
		node.NewParameterReference(nil, "matches", 2, nil, data.NewBaseType("array")),
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
