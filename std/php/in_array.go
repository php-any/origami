package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewInArrayFunction() data.FuncStmt {
	return &InArrayFunction{}
}

type InArrayFunction struct{}

func (f *InArrayFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	needleValue, _ := ctx.GetIndexValue(0)
	haystackValue, _ := ctx.GetIndexValue(1)
	strictValue, _ := ctx.GetIndexValue(2)

	if haystackValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 检查是否为数组
	arrayVal, ok := haystackValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	if needleValue == nil {
		return data.NewBoolValue(false), nil
	}

	// 处理严格模式
	strict := false
	if strictValue != nil {
		if _, ok := strictValue.(*data.NullValue); !ok {
			if strictBool, ok := strictValue.(data.AsBool); ok {
				if s, err := strictBool.AsBool(); err == nil {
					strict = s
				}
			}
		}
	}

	// 在数组中查找
	for _, val := range arrayVal.Value {
		if strict {
			// 严格模式：类型和值都必须相同
			if needleValue == val {
				return data.NewBoolValue(true), nil
			}
		} else {
			// 非严格模式：比较字符串值
			if needleValue.AsString() == val.AsString() {
				return data.NewBoolValue(true), nil
			}
		}
	}

	return data.NewBoolValue(false), nil
}

func (f *InArrayFunction) GetName() string {
	return "in_array"
}

func (f *InArrayFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "needle", 0, nil, nil),
		node.NewParameter(nil, "haystack", 1, nil, nil),
		node.NewParameter(nil, "strict", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *InArrayFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "needle", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "haystack", 1, data.NewBaseType("array")),
		node.NewVariable(nil, "strict", 2, data.NewBaseType("bool")),
	}
}
