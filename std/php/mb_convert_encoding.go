package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// MbConvertEncodingFunction 实现 mb_convert_encoding。
// PHP：string|array $string — 传入 array 时递归转换每个元素并返回 array。
type MbConvertEncodingFunction struct{}

func NewMbConvertEncodingFunction() data.FuncStmt { return &MbConvertEncodingFunction{} }

func (f *MbConvertEncodingFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewStringValue(""), nil
	}
	return mbConvertEncodingValue(v), nil
}

func mbConvertEncodingValue(v data.Value) data.Value {
	switch t := v.(type) {
	case *data.ArrayValue:
		list := make([]*data.ZVal, len(t.List))
		for i, z := range t.List {
			if z == nil {
				list[i] = data.NewZVal(data.NewNullValue())
				continue
			}
			nz := data.NewZVal(mbConvertEncodingValue(z.Value))
			nz.Name = z.Name
			list[i] = nz
		}
		return &data.ArrayValue{List: list}
	default:
		// 简化：假定已是目标编码（常见 UTF-8）
		return data.NewStringValue(v.AsString())
	}
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
		node.NewVariable(nil, "string", 0, nil),
		node.NewVariable(nil, "to_encoding", 1, data.NewBaseType("string")),
		node.NewVariable(nil, "from_encoding", 2, data.NewNullableType(data.NewBaseType("string"))),
	}
}
