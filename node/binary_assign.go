package node

import (
	"errors"
	"fmt"

	"github.com/php-any/origami/data"
)

type BinaryAssign struct {
	*Node `pp:"-"`
	Left  data.GetValue
	Right data.GetValue
}

func NewBinaryAssign(from data.From, left, right data.GetValue) BinaryExpression {
	switch l := left.(type) {
	case *VariableList:
		return &BinaryAssignVariableList{
			Node:  NewNode(from),
			Left:  l,
			Right: right,
		}
	case data.Variable:
		return &BinaryAssignVariable{
			Node:  NewNode(from),
			Left:  l,
			Right: right,
		}
	}

	return &BinaryAssign{
		Node:  NewNode(from),
		Left:  left,
		Right: right,
	}
}

func (b *BinaryAssign) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	if rv == nil {
		rv = data.NewNullValue()
	}
	if v, ok := rv.(data.Value); ok {
		switch l := b.Left.(type) {
		case *CallObjectProperty:
			temp, acl := l.Object.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			switch object := temp.(type) {
			case *data.ClassValue: // 需要检查属性类型
				property, ok := object.GetPropertyStmt(l.Property)
				if ok {
					if property.GetType() != nil && !property.GetType().Is(v) {
						return nil, data.NewErrorThrow(b.GetFrom(), fmt.Errorf("%s 属性 %s 因为类型不一致无法赋值", TryGetCallClassName(object), l.Property))
					}
				}
				return v, object.SetProperty(l.Property, v)
			case data.SetProperty:
				return v, object.SetProperty(l.Property, v)
			default:
				return nil, data.NewErrorThrow(b.GetFrom(), errors.New("object is not set property"))
			}
		case *IndexExpression:
			// 统一走 IndexExpression 自身的 SetValue 逻辑，以支持任意嵌套：
			// $namespace['commands'][1]['sub'] = 'foo';
			if ctl := l.SetValue(ctx, v); ctl != nil {
				return nil, ctl
			}
			return v, nil
		case *CallStaticProperty:
			return v, l.SetProperty(ctx, l.Property, v)
		case *CallSelfProperty:
			return v, l.SetProperty(ctx, l.Property, v)
		case *BinaryNeStrict: // !==

			return data.NewBoolValue(false), nil
		case *BinaryLand: // &&
			lv, acl := l.Left.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			rv, acl := l.Right.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}

			if v, ok := lv.(data.AsBool); ok {
				bv, err := v.AsBool()
				if err != nil {
					return nil, data.NewErrorThrow(b.from, err)
				}
				if !bv {
					return data.NewBoolValue(false), nil
				}
			} else {
				return data.NewBoolValue(false), nil
			}

			if v, ok := rv.(data.AsBool); ok {
				bv, err := v.AsBool()
				if err != nil {
					return nil, data.NewErrorThrow(b.from, err)
				}
				if !bv {
					return data.NewBoolValue(false), nil
				}
			} else {
				return data.NewBoolValue(false), nil
			}

			return data.NewBoolValue(true), nil
		case *BinaryEqStrict, *BinaryEq, *BinaryNe:
			// 处理其他比较运算符的相似情况
			// 由于不同的比较类型，我们需要通过反射或类型断言来提取 Left 和 Right
			// 为简化，这里只给出友好的错误信息
			return nil, data.NewErrorThrow(b.from, fmt.Errorf("赋值表达式的左侧不能是比较表达式，请使用括号: (%T)", l))
		case *CallStaticPropertyLater:
			return v, l.SetProperty(ctx, l.property, v)
		case *CallStaticKeywordProperty:
			return v, l.SetProperty(ctx, l.Property, v)
		case *Array:
			if rv, ok := rv.(*data.ArrayValue); ok {
				valueList := rv.ToValueList()
				for i, value := range valueList {
					if i < len(l.V) {
						set := l.V[i]
						if set, ok := set.(data.Variable); ok {
							set.SetValue(ctx, value)
						}
					}
				}
			} else if cv, ok := v.(*data.ClassValue); ok {
				// 支持实现了 ArrayAccess 的对象解构（如 Collection）
				if method, exists := cv.GetMethod("offsetGet"); exists {
					for i, set := range l.V {
						fnCtx := cv.CreateContext(method.GetVariables())
						if len(method.GetVariables()) > 0 {
							fnCtx.SetVariableValue(method.GetVariables()[0], data.NewIntValue(i))
						}
						ret, ctl := method.Call(fnCtx)
						if ctl != nil {
							return nil, ctl
						}
						var val data.Value = data.NewNullValue()
						if rv2, ok2 := ret.(data.Value); ok2 {
							val = rv2
						}
						if s, ok2 := set.(data.Variable); ok2 {
							if ctl := s.SetValue(ctx, val); ctl != nil {
								return nil, ctl
							}
						}
					}
				}
			}
			return data.NewBoolValue(true), nil
		// 其它 data.Variable 类型（包括各类超全局变量节点），统一走其自身的 SetValue 逻辑；
		// 若某个超全局是只读的，应在对应节点的 SetValue 中返回错误或忽略写入。
		default:
			return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO 赋值表达式遇到未支持的类型: %T", l))
		}
	}

	return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO BinaryAssign rv=%T left=%T", rv, b.Left))
}

type BinaryAssignVariable struct {
	*Node `pp:"-"`
	Left  data.Variable
	Right data.GetValue
}

func (b *BinaryAssignVariable) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	if v, ok := rv.(data.Value); ok {
		return v, b.Left.SetValue(ctx, v)
	}

	if rv == nil {
		v := data.NewNullValue()
		return v, b.Left.SetValue(ctx, v)
	}

	return nil, data.NewErrorThrow(b.from, fmt.Errorf("TODO BinaryAssign rv=%T left=%T", rv, b.Left))
}

type BinaryAssignVariableList struct {
	*Node `pp:"-"`
	Left  *VariableList
	Right data.GetValue
}

func (b *BinaryAssignVariableList) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}

	// 检查是否是ReturnControl
	if returnControl, ok := rv.(data.ReturnControl); ok {
		// 从ReturnControl中提取ArrayValue
		arrayValue := returnControl.ReturnValue()
		// 使用VariableList的SetValue方法来处理多变量赋值
		ctl := b.Left.SetValue(ctx, arrayValue)
		if ctl != nil {
			return nil, ctl
		}
		return arrayValue, nil
	}

	// 处理普通值
	if v, ok := rv.(data.Value); ok {
		// 如果值是实现了 ArrayAccess 的对象，通过 offsetGet 逐个取值
		if cv, ok2 := v.(*data.ClassValue); ok2 {
			if method, exists := cv.GetMethod("offsetGet"); exists {
				for i, lv := range b.Left.Vars {
					fnCtx := cv.CreateContext(method.GetVariables())
					if len(method.GetVariables()) > 0 {
						fnCtx.SetVariableValue(method.GetVariables()[0], data.NewIntValue(i))
					}
					ret, ctl := method.Call(fnCtx)
					if ctl != nil {
						return nil, ctl
					}
					var val data.Value = data.NewNullValue()
					if rv, ok3 := ret.(data.Value); ok3 {
						val = rv
					}
					if ctl := lv.SetValue(ctx, val); ctl != nil {
						return nil, ctl
					}
				}
				return v, nil
			}
		}
		// 使用VariableList的SetValue方法来处理多变量赋值
		ctl := b.Left.SetValue(ctx, v)
		if ctl != nil {
			return nil, ctl
		}
		return v, nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("多变量赋值失败"))
}
