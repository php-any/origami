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
	var valueList []data.Value

	if separatorValue != nil && isImplodeArray(separatorValue) {
		valueList = implodeValueList(separatorValue)
		if arrayValue != nil {
			separator = arrayValue.AsString()
		}
	} else {
		if separatorValue != nil {
			separator = separatorValue.AsString()
		}
		if arrayValue != nil {
			valueList = implodeValueList(arrayValue)
		}
	}

	if valueList == nil {
		return data.NewStringValue(""), nil
	}

	var parts []string
	for _, val := range valueList {
		part, ctl := node.ValueToDisplayString(ctx, val)
		if ctl != nil {
			return nil, ctl
		}
		parts = append(parts, part)
	}

	return data.NewStringValue(strings.Join(parts, separator)), nil
}

func isImplodeArray(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func implodeValueList(v data.Value) []data.Value {
	switch arr := v.(type) {
	case *data.ArrayValue:
		return arr.ToValueList()
	case *data.ObjectValue:
		var list []data.Value
		arr.RangeProperties(func(_ string, val data.Value) bool {
			list = append(list, val)
			return true
		})
		return list
	default:
		return nil
	}
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
