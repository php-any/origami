package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
	"github.com/php-any/origami/utils"
)

// 辅助函数：调用无参无返回值方法（void）
func callVoidMethod(obj *data.ClassValue, name string) data.Control {
	if m, ok := obj.GetMethod(name); ok {
		fnCtx := obj.CreateContext(m.GetVariables())
		_, ctl := m.Call(fnCtx)
		return ctl
	}
	return nil
}

// 辅助函数：调用返回 Value 的方法
func callValueMethod(obj *data.ClassValue, name string) (data.Value, data.Control) {
	if m, ok := obj.GetMethod(name); ok {
		fnCtx := obj.CreateContext(m.GetVariables())
		v, ctl := m.Call(fnCtx)
		if ctl != nil {
			return nil, ctl
		}
		if val, ok := v.(data.Value); ok {
			return val, nil
		}
		return data.NewNullValue(), nil
	}
	return data.NewNullValue(), nil
}

// 辅助函数：调用返回 bool 的方法
func callBoolMethod(obj *data.ClassValue, name string) (bool, data.Control) {
	v, ctl := callValueMethod(obj, name)
	if ctl != nil {
		return false, ctl
	}
	if b, ok := v.(data.AsBool); ok {
		vb, err := b.AsBool()
		if err != nil {
			return false, utils.NewThrow(err)
		}
		return vb, nil
	}
	return v != nil, nil
}

// 辅助函数：检查类是否实现了接口
func checkInterfaceStructure(source data.ClassStmt, target data.InterfaceStmt) bool {
	// 获取目标接口的所有方法
	targetMethods := target.GetMethods()

	// 检查源类是否实现了目标接口的所有方法
	for _, targetMethod := range targetMethods {
		methodName := targetMethod.GetName()
		sourceMethod, exists := source.GetMethod(methodName)

		if !exists {
			return false
		}

		// 检查方法签名是否匹配（参数数量）
		if len(sourceMethod.GetParams()) != len(targetMethod.GetParams()) {
			return false
		}
	}

	return true
}

// EndFunction 实现 end 函数
// 将数组的内部指针移动到最后一个元素，并返回该元素的值
type EndFunction struct{}

func NewEndFunction() data.FuncStmt {
	return &EndFunction{}
}

func (f *EndFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue, _ := ctx.GetIndexValue(0)

	if arrayValue == nil {
		return data.NewNullValue(), nil
	}

	// 使用类型 switch 处理不同类型
	switch val := arrayValue.(type) {
	case *data.ArrayValue:
		// 处理数组
		if len(val.Value) == 0 {
			return data.NewNullValue(), nil
		}
		// 返回最后一个元素
		return val.Value[len(val.Value)-1], nil

	case *data.ObjectValue:
		// 处理对象（关联数组）
		var lastValue data.Value
		var hasValue bool

		// 使用 RangeProperties 按插入顺序遍历，获取最后一个元素
		val.RangeProperties(func(key string, value data.Value) bool {
			lastValue = value
			hasValue = true
			return true // 继续遍历到最后一个
		})

		if !hasValue {
			return data.NewNullValue(), nil
		}
		return lastValue, nil

	case *data.ClassValue:
		// 处理 Iterator 对象
		// 检查是否实现了 Iterator 接口
		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(val.Class, targetInterface) {
				// 重置迭代器到开始位置
				if ctl := callVoidMethod(val, "rewind"); ctl != nil {
					return data.NewNullValue(), nil
				}

				var lastValue data.Value
				var hasValue bool

				// 遍历到最后一个元素
				for {
					// 检查当前位置是否有效
					valid, ctl := callBoolMethod(val, "valid")
					if ctl != nil {
						return data.NewNullValue(), nil
					}
					if !valid {
						break
					}

					// 获取当前元素
					currentVal, ctl := callValueMethod(val, "current")
					if ctl != nil {
						return data.NewNullValue(), nil
					}
					lastValue = currentVal
					hasValue = true

					// 移动到下一个位置
					if ctl := callVoidMethod(val, "next"); ctl != nil {
						return data.NewNullValue(), nil
					}
				}

				if !hasValue {
					return data.NewNullValue(), nil
				}
				return lastValue, nil
			}
		}
		// 不是 Iterator 接口，返回 null
		return data.NewNullValue(), nil

	default:
		// 不是数组类型，返回 null
		return data.NewNullValue(), nil
	}
}

func (f *EndFunction) GetName() string {
	return "end"
}

func (f *EndFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, data.Mixed{}),
	}
}

func (f *EndFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
	}
}
