package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewCountFunction() data.FuncStmt {
	return &CountFunction{}
}

type CountFunction struct{}

func (f *CountFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, _ := ctx.GetIndexValue(0)
	modeValue, _ := ctx.GetIndexValue(1)

	if value == nil {
		return data.NewIntValue(0), nil
	}

	// 检查是否为 NullValue
	if _, ok := value.(*data.NullValue); ok {
		return data.NewIntValue(0), nil
	}

	// 处理数组
	if arrayVal, ok := value.(*data.ArrayValue); ok {
		// 处理 mode 参数（COUNT_RECURSIVE）
		mode := 0
		if modeValue != nil {
			if _, ok := modeValue.(*data.NullValue); !ok {
				if modeInt, ok := modeValue.(data.AsInt); ok {
					if m, err := modeInt.AsInt(); err == nil {
						mode = m
					}
				}
			}
		}

		if mode == 1 {
			// COUNT_RECURSIVE: 递归计算多维数组
			return data.NewIntValue(f.countRecursive(arrayVal)), nil
		}

		// 普通模式：只计算顶层元素
		return data.NewIntValue(len(arrayVal.List)), nil
	}

	// 处理对象
	if objectVal, ok := value.(*data.ObjectValue); ok {
		properties := objectVal.GetProperties()
		return data.NewIntValue(len(properties)), nil
	}

	// 处理 ClassValue / ThisValue（支持 Countable 接口）
	if classVal, ok := value.(*data.ClassValue); ok {
		if method, has := classVal.GetMethod("count"); has {
			fnCtx := classVal.CreateContext(method.GetVariables())
			ret, ctl := method.Call(fnCtx)
			if ctl == nil {
				if iv, ok2 := ret.(data.AsInt); ok2 {
					if i, err := iv.AsInt(); err == nil {
						return data.NewIntValue(i), nil
					}
				}
			}
		}
		return data.NewIntValue(len(classVal.ObjectValue.GetProperties())), nil
	}
	if thisVal, ok := value.(*data.ThisValue); ok {
		if method, has := thisVal.GetMethod("count"); has {
			fnCtx := thisVal.CreateContext(method.GetVariables())
			ret, ctl := method.Call(fnCtx)
			if ctl == nil {
				if iv, ok2 := ret.(data.AsInt); ok2 {
					if i, err := iv.AsInt(); err == nil {
						return data.NewIntValue(i), nil
					}
				}
			}
		}
		return data.NewIntValue(1), nil
	}
	// 其他类型返回 1
	return data.NewIntValue(1), nil
}

// countRecursive 递归计算数组元素数量
func (f *CountFunction) countRecursive(arrayVal *data.ArrayValue) int {
	count := len(arrayVal.List)
	valueList := arrayVal.ToValueList()
	for _, val := range valueList {
		if nestedArray, ok := val.(*data.ArrayValue); ok {
			count += f.countRecursive(nestedArray)
		}
	}
	return count
}

func (f *CountFunction) GetName() string {
	return "count"
}

func (f *CountFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
		node.NewParameter(nil, "mode", 1, node.NewNullLiteral(nil), nil),
	}
}

func (f *CountFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
		node.NewVariable(nil, "mode", 1, data.NewBaseType("int")),
	}
}
