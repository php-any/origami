package php

import (
	"unicode/utf8"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MbDetectEncodingFunction 实现 mb_detect_encoding（简化版）。
// mb_detect_encoding(string $string, array|string|null $encodings = null, bool $strict = false): string|false
// 对合法 UTF-8 返回 "UTF-8"；否则在 strict 时返回 false，非 strict 仍尝试返回 "UTF-8"。
type MbDetectEncodingFunction struct{}

func NewMbDetectEncodingFunction() data.FuncStmt {
	return &MbDetectEncodingFunction{}
}

func (f *MbDetectEncodingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewBoolValue(false), nil
	}
	s := v.AsString()
	strict := false
	if sv, ok := ctx.GetIndexValue(2); ok && sv != nil {
		if b, ok := sv.(data.AsBool); ok {
			strict, _ = b.AsBool()
		}
	}
	if utf8.ValidString(s) {
		return data.NewStringValue("UTF-8"), nil
	}
	if strict {
		return data.NewBoolValue(false), nil
	}
	// 非 strict：仍给出一个编码猜测，便于调用方继续走 mb_* 路径
	return data.NewStringValue("UTF-8"), nil
}

func (f *MbDetectEncodingFunction) GetName() string { return "mb_detect_encoding" }

func (f *MbDetectEncodingFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, data.String{}),
		node.NewParameter(nil, "encodings", 1, node.NewNullLiteral(nil), nil),
		node.NewParameter(nil, "strict", 2, data.NewBoolValue(false), data.Bool{}),
	}
}

func (f *MbDetectEncodingFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.String{}),
		node.NewVariable(nil, "encodings", 1, nil),
		node.NewVariable(nil, "strict", 2, data.Bool{}),
	}
}
