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
	case data.Variable:
		return &BinaryAssignVariable{
			Node:  NewNode(from),
			Left:  l,
			Right: right,
		}
	case *VariableList:
		return &BinaryAssignVariableList{
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
			// 索引赋值 $this->where[key] = value
			idxExpr := l
			arrayVal, acl := idxExpr.Array.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}
			if _, ok := arrayVal.(*data.NullValue); ok {
				if ie, ok := idxExpr.Array.(*IndexExpression); ok {
					// 多级访问，自动创建空数组
					arrayVal = data.NewObjectValue()
					_, acl = NewBinaryAssign(b.from, ie, arrayVal).GetValue(ctx)
					if acl != nil {
						return nil, acl
					}
				}
			}

			indexVal, acl := idxExpr.Index.GetValue(ctx)
			if acl != nil {
				return nil, acl
			}

			switch arr := arrayVal.(type) {
			case *data.ArrayValue:
				// 数组索引赋值
				i := 0
				if iv, ok := indexVal.(data.AsInt); ok {
					var err error
					i, err = iv.AsInt()
					if err != nil {
						return nil, data.NewErrorThrow(b.from, err)
					}
				} else if iv, ok := indexVal.(data.AsString); ok {
					objectVal := data.NewObjectValue()
					valueList := arr.ToValueList()
					for i2, value := range valueList {
						objectVal.SetProperty(fmt.Sprintf("%d", i2), value)
					}
					objectVal.SetProperty(iv.AsString(), v)
					// 重新赋值
					_, acl = NewBinaryAssign(b.from, idxExpr.Array, objectVal).GetValue(ctx)
					if acl != nil {
						return nil, acl
					}

				} else {
					return nil, data.NewErrorThrow(b.from, errors.New("数组索引不是整数类型"))
				}
				if i < 0 {
					return nil, data.NewErrorThrow(b.from, errors.New("数组索引不能为负数"))
				}
				if i >= len(arr.List) {
					// 自动扩容，填充 null
					for j := len(arr.List); j <= i; j++ {
						arr.List = append(arr.List, data.NewZVal(data.NewNullValue()))
					}
				}
				arr.List[i] = data.NewZVal(v)
				return v, nil
			case data.SetProperty:
				// 对象属性赋值
				if iv, ok := indexVal.(data.AsString); ok {
					arr.SetProperty(iv.AsString(), v)
					return v, nil
				} else {
					return nil, data.NewErrorThrow(b.from, errors.New("对象属性索引不是字符串类型"))
				}
			default:
				return nil, data.NewErrorThrow(b.from, errors.New("索引赋值仅支持数组或对象"))
			}
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
		case *CallStaticKeywordProperty:
			return v, l.SetProperty(ctx, l.Property, v)
		case *Array:
			if rv, ok := rv.(*data.ArrayValue); ok {
				valueList := rv.ToValueList()
				for i, value := range valueList {
					set := l.V[i]
					if set, ok := set.(data.Variable); ok {
						set.SetValue(ctx, value)
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

	return nil, data.NewErrorThrow(b.from, errors.New("TODO BinaryAssign"))
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

	return nil, data.NewErrorThrow(b.from, errors.New("TODO BinaryAssign"))
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
		// 使用VariableList的SetValue方法来处理多变量赋值
		ctl := b.Left.SetValue(ctx, v)
		if ctl != nil {
			return nil, ctl
		}
		return v, nil
	}

	return nil, data.NewErrorThrow(b.from, errors.New("多变量赋值失败"))
}
