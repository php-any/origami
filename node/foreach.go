package node

import (
	"fmt"

	"github.com/php-any/origami/data"
)

// 工具：调用无参无返回值方法（void）
func callVoidMethod(obj *data.ClassValue, name string) data.Control {
	if m, ok := obj.GetMethod(name); ok {
		fnCtx := obj.CreateContext(m.GetVariables())
		_, ctl := m.Call(fnCtx)
		return ctl
	}
	return nil
}

// 工具：调用返回 Value 的方法
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

// 工具：调用返回 bool 的方法
func callBoolMethod(obj *data.ClassValue, name string) (bool, data.Control) {
	v, ctl := callValueMethod(obj, name)
	if ctl != nil {
		return false, ctl
	}
	if b, ok := v.(data.AsBool); ok {
		vb, err := b.AsBool()
		if err != nil {
			return false, data.NewErrorThrow(nil, err)
		}
		return vb, nil
	}
	return v != nil, nil
}

// ForeachStatement 表示foreach语句
type ForeachStatement struct {
	*Node `pp:"-"`
	Array data.GetValue // 要遍历的数组
	Key   data.Variable // 键变量名（可选）
	Value data.Variable // 值变量名
	Body  []data.GetValue
}

func (u *ForeachStatement) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 获取数组值
	arrayValue, ctl := u.Array.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}

	// 检查数组值是否为数组类型
	switch array := arrayValue.(type) {
	case *data.ClassValue:
		// 类实例：若实现语言层 Iterator 接口，则按五方法迭代
		var v data.GetValue
		var c data.Control

		if targetInterface, ok := ctx.GetVM().GetInterface("Iterator"); ok {
			if checkInterfaceStructure(array.Class, targetInterface) {
				// rewind
				if ctl := callVoidMethod(array, "rewind"); ctl != nil {
					return nil, ctl
				}

				for {
					// valid
					valid, ctl := callBoolMethod(array, "valid")
					if ctl != nil {
						return nil, ctl
					}
					if !valid {
						break
					}

					// current
					valV, ctl := callValueMethod(array, "current")
					if ctl != nil {
						return nil, ctl
					}

					// key
					var keyV data.Value
					if u.Key != nil {
						kv, kctl := callValueMethod(array, "key")
						if kctl != nil {
							return nil, kctl
						}
						keyV = kv
					}

					ctx.SetVariableValue(u.Value, valV)
					if u.Key != nil {
						ctx.SetVariableValue(u.Key, keyV)
					}

					for _, statement := range u.Body {
						v, c = statement.GetValue(ctx)
						if c != nil {
							if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
								return nil, nil
							}
							if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
								// 推进迭代器进入下一次
								if ctl := callVoidMethod(array, "next"); ctl != nil {
									return nil, ctl
								}
								break
							}
							return nil, checkThrowControlFrom(statement, c)
						}
					}

					// next
					if ctl := callVoidMethod(array, "next"); ctl != nil {
						return nil, ctl
					}
				}

				return v, nil
			}
		}

		// 非 Iterator 类实例则按对象属性遍历
		{
			var v data.GetValue
			var c data.Control

			for i, element := range array.GetProperties() {
				ctx.SetVariableValue(u.Value, element)
				if u.Key != nil {
					ctx.SetVariableValue(u.Key, data.NewStringValue(i))
				}

				for _, statement := range u.Body {
					v, c = statement.GetValue(ctx)
					if c != nil {
						if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
							return nil, nil
						}
						if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
							break
						}
						return nil, c
					}
				}
			}
			return v, nil
		}
		// 不再 fallthrough 到数组
	case *data.ArrayValue:
		var v data.GetValue
		var c data.Control

		// 遍历数组
		for i, element := range array.Value {
			// 设置值变量
			ctx.SetVariableValue(u.Value, element)

			// 如果有键变量，设置键变量
			if u.Key != nil {
				keyValue := data.NewIntValue(i)
				ctx.SetVariableValue(u.Key, keyValue)
			}

			// 执行循环体
			for _, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					// break 跳出循环
					if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
						return nil, nil
					}
					// continue 跳到下一次迭代
					if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
						continue
					}
					// return/throw 直接返回
					return nil, c
				}
			}
		}
		return v, nil
	case *data.ObjectValue:
		var v data.GetValue
		var c data.Control

		// 遍历数组
		for i, element := range array.GetProperties() {
			// 设置值变量
			ctx.SetVariableValue(u.Value, element)

			// 如果有键变量，设置键变量
			if u.Key != nil {
				keyValue := data.NewStringValue(i)
				ctx.SetVariableValue(u.Key, keyValue)
			}

			// 执行循环体
			for _, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					// break 跳出循环
					if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
						return nil, nil
					}
					// continue 跳到下一次迭代
					if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
						continue
					}
					// return/throw 直接返回
					return nil, c
				}
			}
		}
		return v, nil
	case *data.NullValue:
		return nil, nil
	}

	return nil, data.NewErrorThrow(u.from, fmt.Errorf("foreach 只能遍历数组、对象或实现 Iterator 的值"))
}

// NewForeachStatement 创建一个新的foreach语句
func NewForeachStatement(token *TokenFrom, array data.GetValue, key data.Variable, value data.Variable, body []data.GetValue) *ForeachStatement {
	return &ForeachStatement{
		Node:  NewNode(token),
		Array: array,
		Key:   key,
		Value: value,
		Body:  body,
	}
}
