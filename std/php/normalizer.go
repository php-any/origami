package php

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// 实现 normalizer_is_normalized 与 normalizer_normalize 两个函数，
// 以模拟 intl 扩展的最常用行为，主要供 Symfony polyfill/string 使用。

// NormalizerIsNormalizedFunction 对应 normalizer_is_normalized()
//
// 签名：
//
//	normalizer_is_normalized(?string $string, ?int $form = Normalizer::FORM_C): bool
//
// 这里简单实现为：
//   - 若参数为 null，则视为空串，返回 true
//   - 若不是有效 UTF-8，则返回 false
//   - 否则一律返回 true
type NormalizerIsNormalizedFunction struct{}

func NewNormalizerIsNormalizedFunction() data.FuncStmt {
	return &NormalizerIsNormalizedFunction{}
}

func (f *NormalizerIsNormalizedFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	// formVal, _ := ctx.GetIndexValue(1) // 当前实现忽略具体 form

	if strVal == nil {
		// null 视为空串
		return data.NewBoolValue(true), nil
	}
	s := strVal.AsString()

	if !utf8.ValidString(s) {
		return data.NewBoolValue(false), nil
	}
	return data.NewBoolValue(true), nil
}

func (f *NormalizerIsNormalizedFunction) GetName() string {
	return "normalizer_is_normalized"
}

func (f *NormalizerIsNormalizedFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "form", 1, node.NewIntLiteral(nil, "4"), nil), // 默认 FORM_C
	}
}

func (f *NormalizerIsNormalizedFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewNullableType(data.NewBaseType("string"))),
		node.NewVariable(nil, "form", 1, data.NewBaseType("int")),
	}
}

// NormalizerNormalizeFunction 对应 normalizer_normalize()
//
// 签名：
//
//	normalizer_normalize(?string $string, ?int $form = Normalizer::FORM_C): string|false
//
// 这里的近似行为：
//   - null 视为空串，返回空串
//   - 若不是有效 UTF-8，则返回 false
//   - 否则原样返回（不做真实 Unicode 归一化）
type NormalizerNormalizeFunction struct{}

func NewNormalizerNormalizeFunction() data.FuncStmt {
	return &NormalizerNormalizeFunction{}
}

func (f *NormalizerNormalizeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	strVal, _ := ctx.GetIndexValue(0)
	// formVal, _ := ctx.GetIndexValue(1) // 当前实现忽略具体 form

	if strVal == nil {
		return data.NewStringValue(""), nil
	}
	s := strVal.AsString()
	if !utf8.ValidString(s) {
		return data.NewBoolValue(false), nil
	}
	return data.NewStringValue(s), nil
}

func (f *NormalizerNormalizeFunction) GetName() string {
	return "normalizer_normalize"
}

func (f *NormalizerNormalizeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "form", 1, node.NewIntLiteral(nil, "4"), nil), // 默认 FORM_C
	}
}

func (f *NormalizerNormalizeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewNullableType(data.NewBaseType("string"))),
		node.NewVariable(nil, "form", 1, data.NewBaseType("int")),
	}
}
