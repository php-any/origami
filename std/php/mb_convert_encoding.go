package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MbConvertEncodingFunction 实现 mb_convert_encoding 函数（简化版，仅支持 UTF-8）
type MbConvertEncodingFunction struct{}

func NewMbConvertEncodingFunction() data.FuncStmt { return &MbConvertEncodingFunction{} }

func (f *MbConvertEncodingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	// 简化实现：假设输入已经是 UTF-8，直接返回
	return data.NewStringValue(v.AsString()), nil
}

func (f *MbConvertEncodingFunction) GetName() string { return "mb_convert_encoding" }
func (f *MbConvertEncodingFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "string", 0, nil, nil),
		node.NewParameter(nil, "to_encoding", 1, node.NewStringLiteralByAst(nil, "UTF-8"), nil),
		node.NewParameter(nil, "from_encoding", 2, node.NewNullLiteral(nil), nil),
	}
}
func (f *MbConvertEncodingFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "string", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "to_encoding", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "from_encoding", 2, data.NewNullableType(data.NewBaseType("string"))),
	}
}
