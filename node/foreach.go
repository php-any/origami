package node

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/utils"
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
			return false, utils.NewThrow(err)
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
			var shouldBreak bool
			var shouldReturn bool

			// 使用 RangeProperties 保证遍历顺序与插入顺序一致
			array.RangeProperties(func(i string, element data.Value) bool {
				ctx.SetVariableValue(u.Value, element)
				if u.Key != nil {
					ctx.SetVariableValue(u.Key, data.NewStringValue(i))
				}

				for _, statement := range u.Body {
					v, c = statement.GetValue(ctx)
					if c != nil {
						if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
							shouldBreak = true
							return false // 停止遍历
						}
						if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
							return true // 继续下一次迭代
						}
						shouldReturn = true
						return false // 停止遍历
					}
				}
				return true
			})

			if shouldBreak {
				return nil, nil
			}
			if shouldReturn {
				return nil, c
			}
			return v, nil
		}
		// 不再 fallthrough 到数组
	case *data.ArrayValue:
		var v data.GetValue
		var c data.Control

		// 遍历数组
		valueList := array.ToValueList()
		for i, element := range valueList {
			// 设置值变量
			acl := u.Value.SetValue(ctx, element)
			if acl != nil {
				return nil, acl
			}
			// 如果有键变量，设置键变量
			if u.Key != nil {
				keyValue := data.NewIntValue(i)
				ctx.SetVariableValue(u.Key, keyValue)
			}

			// 执行循环体
			for bodyIndex, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					switch ctrl := c.(type) {
					case data.BreakControl:
						if ctrl.IsBreak() {
							return nil, nil
						}
					case data.ContinueControl:
						if ctrl.IsContinue() {
							goto nextArrayElement
						}
					case data.YieldValueControl:
						// yield：包装成 ForeachArrayYieldControl，保存迭代状态
						return nil, NewForeachArrayYieldControl(u, valueList, i, bodyIndex+1, ctrl)
					}
					// return/throw 直接返回
					return nil, c
				}
			}
		nextArrayElement:
		}
		return v, nil
	case data.Iterator:
		// 直接实现 Iterator 接口的值
		var v data.GetValue
		var c data.Control

		// rewind
		_, ctl := array.Rewind(ctx)
		if ctl != nil {
			return nil, ctl
		}

		for {
			// valid
			validV, ctl := array.Valid(ctx)
			if ctl != nil {
				return nil, ctl
			}
			// 将 valid 结果转换为 bool
			if validBool, ok := validV.(data.AsBool); ok {
				valid, err := validBool.AsBool()
				if err != nil {
					return nil, utils.NewThrow(err)
				}
				if !valid {
					break
				}
			} else {
				// 如果无法转换为 bool，则检查是否为非空值
				if validV == nil {
					break
				}
			}

			// current
			valV, ctl := array.Current(ctx)
			if ctl != nil {
				return nil, ctl
			}
			ctx.SetVariableValue(u.Value, valV)

			// key
			var keyV data.Value
			if u.Key != nil {
				kv, kctl := array.Key(ctx)
				if kctl != nil {
					return nil, kctl
				}
				keyV = kv
				ctx.SetVariableValue(u.Key, keyV)
			}

			shouldSkipNext := false
			for _, statement := range u.Body {
				v, c = statement.GetValue(ctx)
				if c != nil {
					if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
						return nil, nil
					}
					if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
						// 推进迭代器进入下一次
						if ctl := array.Next(ctx); ctl != nil {
							return nil, ctl
						}
						shouldSkipNext = true
						break
					}
					return nil, checkThrowControlFrom(statement, c)
				}
			}

			// next（如果 continue 已经调用过则跳过）
			if !shouldSkipNext {
				if ctl := array.Next(ctx); ctl != nil {
					return nil, ctl
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

// ForeachArrayYieldControl 表示在 foreach 数组循环体内遇到 yield 时的暂停状态
// 实现 data.YieldControl，使得 FuncYieldStackState.Next() 可以正确恢复迭代
type ForeachArrayYieldControl struct {
	*ForeachStatement
	ValueList  []data.Value
	ArrayIndex int // 当前数组元素的索引（下次恢复后还需继续遍历的位置）
	BodyIndex  int // 当前 body 中的位置（下次从哪里继续）
	Value      data.YieldValueControl
}

func NewForeachArrayYieldControl(stmt *ForeachStatement, valueList []data.Value, arrayIndex int, bodyIndex int, v data.YieldValueControl) *ForeachArrayYieldControl {
	return &ForeachArrayYieldControl{
		ForeachStatement: stmt,
		ValueList:        valueList,
		ArrayIndex:       arrayIndex,
		BodyIndex:        bodyIndex,
		Value:            v,
	}
}

func (f *ForeachArrayYieldControl) GetYieldKey() data.Value {
	return f.Value.GetYieldKey()
}

func (f *ForeachArrayYieldControl) GetYieldValue() data.Value {
	return f.Value.GetYieldValue()
}

func (f *ForeachArrayYieldControl) AsString() string {
	return "foreach array yield"
}

func (f *ForeachArrayYieldControl) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	// 推进到下一个 yield
	ctl := f.Next(ctx)
	if ctl != nil {
		// ctl == f 表示找到了下一个 yield；ctl == error 表示发生了错误
		return nil, ctl
	}
	// 迭代完毕，返回 nil
	return nil, nil
}

// Next 从上次暂停的位置继续执行，直到找到下一个 yield 或迭代完毕
// 返回的 Control 是：nil（结束）、self（下一个 yield）或 error
func (f *ForeachArrayYieldControl) Next(ctx data.Context) data.Control {
	// 先继续当前数组元素的剩余 body
	startBodyIndex := f.BodyIndex
	f.BodyIndex = 0

	for i := f.ArrayIndex; i < len(f.ValueList); i++ {
		element := f.ValueList[i]
		// 设置变量
		if acl := f.ForeachStatement.Value.SetValue(ctx, element); acl != nil {
			return acl
		}
		if f.ForeachStatement.Key != nil {
			ctx.SetVariableValue(f.ForeachStatement.Key, data.NewIntValue(i))
		}

		startBody := startBodyIndex
		startBodyIndex = 0 // 后续元素从头开始

		for bodyIndex := startBody; bodyIndex < len(f.Body); bodyIndex++ {
			statement := f.Body[bodyIndex]
			_, c := statement.GetValue(ctx)
			if c != nil {
				if ctrl, ok := c.(data.BreakControl); ok && ctrl.IsBreak() {
					return nil
				}
				if ctrl, ok := c.(data.ContinueControl); ok && ctrl.IsContinue() {
					break // 跳到下一个数组元素
				}
				if ctrl, ok := c.(data.YieldValueControl); ok {
					f.Value = ctrl
					f.ArrayIndex = i
					f.BodyIndex = bodyIndex + 1
					return f
				}
				return c
			}
		}
	}
	return nil
}

func (f *ForeachArrayYieldControl) Rewind(ctx data.Context) (data.Value, data.Control) {
	return data.NewNullValue(), nil
}

func (f *ForeachArrayYieldControl) Valid(ctx data.Context) (data.Value, data.Control) {
	return data.NewBoolValue(true), nil
}

func (f *ForeachArrayYieldControl) Current(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldValue(), nil
}

func (f *ForeachArrayYieldControl) Key(ctx data.Context) (data.Value, data.Control) {
	return f.Value.GetYieldKey(), nil
}

func (f *ForeachArrayYieldControl) CreateStackState(ctx data.Context, fn data.FuncStmt, originalBody []data.GetValue, bodyIndex int) data.Generator {
	// 构建 newBody：将 bodyIndex 位置（foreach语句）替换为 f 自身
	newBody := make([]data.GetValue, len(originalBody))
	copy(newBody, originalBody)
	newBody[bodyIndex] = f
	// 获取当前 yield 的值
	currentKey := f.Value.GetYieldKey()
	currentValue := f.Value.GetYieldValue()
	return NewFuncYieldStackState(ctx, fn, newBody, bodyIndex, currentKey, currentValue)
}

type ForeachValueTarget struct {
	V []data.Variable
}

func (f *ForeachValueTarget) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return f, nil
}

func (f *ForeachValueTarget) GetIndex() int {
	return 0
}

func (f *ForeachValueTarget) GetName() string {
	return ""
}

func (f *ForeachValueTarget) GetType() data.Types {
	return nil
}

func (f *ForeachValueTarget) SetValue(ctx data.Context, value data.Value) data.Control {
	switch d := value.(type) {
	case *data.ArrayValue:
		for i, val := range d.List {
			f.V[i].SetValue(ctx, val.Value)
		}
	case *data.ObjectValue:
		i := 0
		for _, val := range d.GetProperties() {
			f.V[i].SetValue(ctx, val)
			i++
		}
	}

	return nil
}
