package std

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ObjectFunction 实现 PHP 风格的 (object) 类型转换。
// 主要为了支持 Laravel helpers 中的 (object) $arguments 用法：
// - 如果传入的是数组（包括关联数组），则将其转换为类似 stdClass 的对象：
//   - 数组的键将作为对象属性名（数值键会被转换为字符串）
//
// - 如果传入的是已是对象（ObjectValue / ClassValue），则直接返回
// - 其它标量值会被包装到一个带有 "scalar" 属性的对象中（与 PHP 行为类似）
type ObjectFunction struct{}

func NewObjectFunction() data.FuncStmt { return &ObjectFunction{} }

func (f *ObjectFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil {
		return data.NewObjectValue(), nil
	}

	switch val := v.(type) {
	case *data.ObjectValue:
		// 已是对象，直接返回
		return val, nil
	case *data.ClassValue:
		// 类实例，按对象语义返回
		return val, nil
	case *data.ArrayValue:
		// 数组 -> 对象：数值键与字符串键都转换为属性名（数值键转为字符串）
		obj := data.NewObjectValue()
		for i, z := range val.List {
			key := fmt.Sprintf("%d", i)
			if z != nil && z.Name != "" {
				key = z.Name
			}

			if z == nil || z.Value == nil {
				obj.SetProperty(key, data.NewNullValue())
			} else {
				obj.SetProperty(key, z.Value)
			}
		}
		return obj, nil
	default:
		// 其它标量/值：包装为带有 scalar 属性的对象
		obj := data.NewObjectValue()
		obj.SetProperty("scalar", v)
		return obj, nil
	}
}

func (f *ObjectFunction) GetName() string { return "object" }

func (f *ObjectFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "value", 0, nil, nil),
	}
}

func (f *ObjectFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "value", 0, data.NewBaseType("mixed")),
	}
}
