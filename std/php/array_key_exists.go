package php

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayKeyExistsFunction() data.FuncStmt {
	return &ArrayKeyExistsFunction{}
}

type ArrayKeyExistsFunction struct{}

func (f *ArrayKeyExistsFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	keyValue, _ := ctx.GetIndexValue(0)
	arrayValue, _ := ctx.GetIndexValue(1)

	if keyValue == nil || arrayValue == nil {
		return data.NewBoolValue(false), nil
	}

	keyStr := keyValue.AsString()

	// 检查数组
	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		// 对于数组，检查索引是否存在
		if keyInt, ok := keyValue.(data.AsInt); ok {
			if i, err := keyInt.AsInt(); err == nil {
				if i >= 0 && i < len(arrayVal.Value) {
					return data.NewBoolValue(true), nil
				}
			}
		}
		// 检查字符串键（关联数组）
		for i, val := range arrayVal.Value {
			// 这里简化处理，实际应该检查数组的键
			// 暂时返回 false
			_ = i
			_ = val
		}
		return data.NewBoolValue(false), nil
	}

	// 检查对象
	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		_, exists := objectVal.GetProperty(keyStr)
		return data.NewBoolValue(exists), nil
	}

	return data.NewBoolValue(false), nil
}

func (f *ArrayKeyExistsFunction) GetName() string {
	return "array_key_exists"
}

func (f *ArrayKeyExistsFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "key", 0, nil, nil),
		node.NewParameter(nil, "array", 1, nil, nil),
	}
}

func (f *ArrayKeyExistsFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "key", 0, data.NewBaseType("string|int")),
		node.NewVariable(nil, "array", 1, data.NewBaseType("array")),
	}
}
