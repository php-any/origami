package php

import (
	"strings"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewImplodeFunction() data.FuncStmt {
	return &ImplodeFunction{}
}

type ImplodeFunction struct{}

func (f *ImplodeFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	separatorValue, _ := ctx.GetIndexValue(0)
	arrayValue, _ := ctx.GetIndexValue(1)

	// 处理参数顺序：implode 可以接受 (separator, array) 或 (array, separator)
	var separator string
	var array *data.ArrayValue

	if separatorValue != nil {
		if _, ok := separatorValue.(*data.ArrayValue); ok {
			// 第一个参数是数组，第二个参数应该是分隔符
			array = separatorValue.(*data.ArrayValue)
			if arrayValue != nil {
				separator = arrayValue.AsString()
			}
		} else {
			// 第一个参数是分隔符
			separator = separatorValue.AsString()
			if arrayValue != nil {
				if arr, ok := arrayValue.(*data.ArrayValue); ok {
					array = arr
				}
			}
		}
	}

	if array == nil {
		return data.NewStringValue(""), nil
	}

	// 将数组元素转换为字符串并连接
	var parts []string
	for _, val := range array.Value {
		parts = append(parts, val.AsString())
	}

	return data.NewStringValue(strings.Join(parts, separator)), nil
}

func (f *ImplodeFunction) GetName() string {
	return "implode"
}

func (f *ImplodeFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "separator", 0, nil, nil),
		node.NewParameter(nil, "array", 1, nil, nil),
	}
}

func (f *ImplodeFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "separator", 0, data.NewBaseType("string")),
		node.NewVariable(nil, "array", 1, data.NewBaseType("array")),
	}
}
