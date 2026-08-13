package array

import (
	"fmt"
	"os"

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

	// PHP 8.1: deprecation when null is passed as key
	if _, isNull := keyValue.(*data.NullValue); isNull {
		fmt.Fprintln(os.Stderr, "Deprecated: Using null as the key parameter for array_key_exists() is deprecated, use an empty string instead")
	}

	keyStr := keyValue.AsString()

	// 检查数组
	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		if keyInt, ok := keyValue.(data.AsInt); ok {
			if i, err := keyInt.AsInt(); err == nil {
				// 必须按 PHP 整数键查找，不能用 List 下标（命名键会占槽位）
				if z, _ := arrayVal.FindSlotByIntKey(i); z != nil {
					return data.NewBoolValue(true), nil
				}
				return data.NewBoolValue(false), nil
			}
		}
		if _, ok := arrayVal.LookupZValByStringKey(keyStr); ok {
			return data.NewBoolValue(true), nil
		}
		return data.NewBoolValue(false), nil
	}

	// 关联数组在运行时可能以 ObjectValue 表示；键存在且值为 null 时仍应返回 true（对齐 PHP）
	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		return data.NewBoolValue(objectVal.HasProperty(keyStr)), nil
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
