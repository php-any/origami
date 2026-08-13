package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayUniqueFunction 实现 array_unique 函数
// array_unique(array $array, int $flags = SORT_STRING): array
// 移除数组中重复的值，保留第一次出现的键
type ArrayUniqueFunction struct{}

func NewArrayUniqueFunction() data.FuncStmt {
	return &ArrayUniqueFunction{}
}

func (f *ArrayUniqueFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 获取第一个参数：数组
	arrayValue, _ := ctx.GetIndexValue(0)
	if arrayValue == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 获取第二个参数：flags（可选，默认为 SORT_STRING = 2）
	flagsValue, _ := ctx.GetIndexValue(1)
	flags := 2 // SORT_STRING 默认值
	if flagsValue != nil {
		if _, ok := flagsValue.(*data.NullValue); !ok {
			if intVal, ok := flagsValue.(data.AsInt); ok {
				if f, err := intVal.AsInt(); err == nil {
					flags = f
				}
			}
		}
	}

	// 处理数组
	if arrayVal, ok := arrayValue.(*data.ArrayValue); ok {
		return f.processArray(arrayVal.ToValueList(), flags), nil
	}

	// 处理对象（关联数组）
	if objectVal, ok := arrayValue.(*data.ObjectValue); ok {
		return f.processObject(objectVal, flags), nil
	}

	// 不是数组类型，返回空数组
	return data.NewArrayValue([]data.Value{}), nil
}

// processArray 处理 ArrayValue
func (f *ArrayUniqueFunction) processArray(values []data.Value, flags int) data.GetValue {
	// 使用 map 来跟踪已见过的值
	seen := make(map[string]bool)
	result := make([]data.Value, 0)

	for _, val := range values {
		key := f.getValueKey(val, flags)
		if !seen[key] {
			seen[key] = true
			result = append(result, val)
		}
	}

	return data.NewArrayValue(result)
}

// processObject 处理 ObjectValue（关联数组）
func (f *ArrayUniqueFunction) processObject(objectVal *data.ObjectValue, flags int) data.GetValue {
	// 使用 map 来跟踪已见过的值
	seen := make(map[string]bool)
	result := data.NewObjectValue()

	properties := objectVal.GetProperties()
	for key, val := range properties {
		valueKey := f.getValueKey(val, flags)
		if !seen[valueKey] {
			seen[valueKey] = true
			result.SetProperty(key, val)
		}
	}

	return result
}

// getValueKey 根据 flags 获取值的唯一键
func (f *ArrayUniqueFunction) getValueKey(val data.Value, flags int) string {
	switch flags {
	case 1: // SORT_NUMERIC
		return f.getNumericKey(val)
	case 2: // SORT_STRING
		return val.AsString()
	case 5: // SORT_LOCALE_STRING
		// 简化实现：使用字符串比较
		return val.AsString()
	default: // SORT_REGULAR (0) 或其他
		return f.getRegularKey(val)
	}
}

// getNumericKey 获取数值键
func (f *ArrayUniqueFunction) getNumericKey(val data.Value) string {
	if intVal, ok := val.(*data.IntValue); ok {
		if i, err := intVal.AsInt(); err == nil {
			return fmt.Sprintf("i:%d", i)
		}
	}
	if floatVal, ok := val.(*data.FloatValue); ok {
		if f, err := floatVal.AsFloat(); err == nil {
			return fmt.Sprintf("f:%g", f)
		}
	}
	// 如果不是数值类型，转换为字符串
	return "s:" + val.AsString()
}

// getRegularKey 获取常规键（SORT_REGULAR）
func (f *ArrayUniqueFunction) getRegularKey(val data.Value) string {
	// SORT_REGULAR 模式：尝试数值比较，否则使用字符串
	if intVal, ok := val.(*data.IntValue); ok {
		if i, err := intVal.AsInt(); err == nil {
			return fmt.Sprintf("i:%d", i)
		}
	}
	if floatVal, ok := val.(*data.FloatValue); ok {
		if f, err := floatVal.AsFloat(); err == nil {
			return fmt.Sprintf("f:%g", f)
		}
	}
	// 其他类型使用字符串表示
	return "s:" + val.AsString()
}

func (f *ArrayUniqueFunction) GetName() string {
	return "array_unique"
}

func (f *ArrayUniqueFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(2), nil),
	}
}

func (f *ArrayUniqueFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "flags", 1, data.NewBaseType("int")),
	}
}
