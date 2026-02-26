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
	valueList := arrayVal.ToValueList()
	for _, val := range valueList {
		if strict {
			// 严格模式：PHP === 语义，类型和值都必须相同（按值比较，非指针）
			if valueEqualStrict(needleValue, val) {
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

// valueEqualStrict 实现 PHP === 语义：类型相同且值相同（按值比较，非指针）
func valueEqualStrict(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch va := a.(type) {
	case *data.StringValue:
		if vb, ok := b.(*data.StringValue); ok {
			return va.Value == vb.Value
		}
	case *data.IntValue:
		if vb, ok := b.(*data.IntValue); ok {
			return va.Value == vb.Value
		}
	case *data.FloatValue:
		if vb, ok := b.(*data.FloatValue); ok {
			return va.Value == vb.Value
		}
	case *data.BoolValue:
		if vb, ok := b.(*data.BoolValue); ok {
			return va.Value == vb.Value
		}
	case *data.NullValue:
		_, ok := b.(*data.NullValue)
		return ok
	default:
		// 其他类型（如 ArrayValue、ObjectValue）在 PHP 中 === 为引用相等，这里保守按指针比较
		return a == b
	}
	return false
}
