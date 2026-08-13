package php

import (
	"strconv"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewFilterVarFunction() data.FuncStmt {
	return &FilterVarFunction{}
}

type FilterVarFunction struct{}

func (fn *FilterVarFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	filter, _ := ctx.GetIndexValue(1)

	// 获取 filter 为整数
	filterInt := 516 // FILTER_DEFAULT
	if fv, ok := filter.(data.AsInt); ok {
		if n, err := fv.AsInt(); err == nil {
			filterInt = n
		}
	}

	switch filterInt {
	case 257: // FILTER_VALIDATE_INT
		if v, ok := value.(*data.StringValue); ok {
			if _, err := strconv.Atoi(v.Value); err == nil {
				return value, nil
			}
		}
		if _, ok := value.(*data.IntValue); ok {
			return value, nil
		}
		return data.NewBoolValue(false), nil // 验证失败返回 false
	case 258: // FILTER_VALIDATE_BOOLEAN
		if _, ok := value.(*data.BoolValue); ok {
			return value, nil
		}
		return data.NewBoolValue(false), nil
	default:
		// 默认：返回原值
		return value, nil
	}
}

func (fn *FilterVarFunction) GetName() string {
	return "filter_var"
}

func (fn *FilterVarFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "filter", 1, data.NewIntValue(516), nil),
	}
}

func (fn *FilterVarFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.Mixed{}),
		node.NewVariable(nil, "filter", 1, data.NewBaseType("int")),
	}
}
