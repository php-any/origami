package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayCombineFunction 实现 array_combine 函数
// array_combine(keys, values): 以 keys 数组的值作为键、values 数组的值作为值，创建关联数组
type ArrayCombineFunction struct{}

func NewArrayCombineFunction() data.FuncStmt {
	return &ArrayCombineFunction{}
}

func (f *ArrayCombineFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	keysVal, _ := ctx.GetIndexValue(0)
	valuesVal, _ := ctx.GetIndexValue(1)

	if keysVal == nil || valuesVal == nil {
		return data.NewBoolValue(false), nil
	}

	// 获取 keys 列表
	var keys []data.Value
	switch kv := keysVal.(type) {
	case *data.ArrayValue:
		keys = kv.ToValueList()
	case *data.ObjectValue:
		kv.RangeProperties(func(_ string, v data.Value) bool {
			keys = append(keys, v)
			return true
		})
	default:
		return data.NewBoolValue(false), nil
	}

	// 获取 values 列表
	var values []data.Value
	switch vv := valuesVal.(type) {
	case *data.ArrayValue:
		values = vv.ToValueList()
	case *data.ObjectValue:
		vv.RangeProperties(func(_ string, v data.Value) bool {
			values = append(values, v)
			return true
		})
	default:
		return data.NewBoolValue(false), nil
	}

	if len(keys) != len(values) {
		return data.NewBoolValue(false), nil
	}

	// 检查所有 key 是否为连续整数 0,1,2,...,n-1
	// 如果是，返回 ArrayValue（PHP 行为：整数键的 array_combine 返回索引数组）
	allSequential := true
	for i, key := range keys {
		if iv, ok := key.(*data.IntValue); ok {
			if iv.Value != i {
				allSequential = false
				break
			}
		} else {
			allSequential = false
			break
		}
	}

	if allSequential {
		return data.NewArrayValue(values), nil
	}

	result := data.NewObjectValue()
	for i, key := range keys {
		result.SetProperty(key.AsString(), values[i])
	}
	return result, nil
}

func (f *ArrayCombineFunction) GetName() string {
	return "array_combine"
}

func (f *ArrayCombineFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "keys", 0, nil, nil),
		node.NewParameter(nil, "values", 1, nil, nil),
	}
}

func (f *ArrayCombineFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "keys", 0, data.Mixed{}),
		node.NewVariable(nil, "values", 1, data.Mixed{}),
	}
}
