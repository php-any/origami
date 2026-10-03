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
	case *VariableExpression:
		// 解析阶段模式识别：对常见整数赋值模式发出 VarFastAssign 节点。
		// 统一使用单一具体类型，保证接口调用位点单态（monomorphic）；
		// 预提取变量索引与字面量值，运行时快速路径无节点类型断言。
		switch r := right.(type) {
		case *VariableExpression:
			return &VarFastAssign{
				Node:   NewNode(from),
				Dst:    l,
				DstIdx: l.Index,
				LhsIdx: r.Index, // src 变量索引存入 LhsIdx
				Slow:   right,
				op:     vfaOpCopy,
			}
		case *BinaryMul:
			lhsIdx, lhsLit := preExtract(r.Left)
			rhsIdx, rhsLit := preExtract(r.Right)
			return &VarFastAssign{
				Node:   NewNode(from),
				Dst:    l,
				DstIdx: l.Index,
				LhsIdx: lhsIdx,
				RhsIdx: rhsIdx,
				LhsLit: lhsLit,
				RhsLit: rhsLit,
				Slow:   right,
				op:     vfaOpMul,
			}
		case *BinaryAdd:
			lhsIdx, lhsLit := preExtract(r.Left)
			rhsIdx, rhsLit := preExtract(r.Right)
			return &VarFastAssign{
				Node:   NewNode(from),
				Dst:    l,
				DstIdx: l.Index,
				LhsIdx: lhsIdx,
				RhsIdx: rhsIdx,
				LhsLit: lhsLit,
				RhsLit: rhsLit,
				Slow:   right,
				op:     vfaOpAdd,
			}
		case *IntLiteral:
			if iv, ok := r.V.(*data.IntValue); ok {
				return &VarFastAssign{
					Node:   NewNode(from),
					Dst:    l,
					DstIdx: l.Index,
					LhsIdx: -1,
					LhsLit: iv.Value,
					Slow:   right,
					op:     vfaOpCopy,
				}
			}
		default:
			return &BinaryAssignVariable{Node: NewNode(from), Left: l, Right: right}
		}
	case *CallObjectProperty, *CallObjectDynamicProperty:
		return &BinaryAssign{Node: NewNode(from), Left: left, Right: right}
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
	if index, ok := b.Left.(*IndexExpression); ok {
		_, compound := rebindDimensionOperand(b.Right, index, index)
		if complexDimension(index) || compound {
			resolved, ctl := resolveDimension(ctx, index)
			if ctl != nil {
				return nil, ctl
			}
			right, _ := rebindDimensionOperand(b.Right, index, resolved)
			value, ctl := right.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
			if value == nil {
				value = data.NewNullValue()
			}
			return value, resolved.setValueResolved(ctx, value.(data.Value))
		}
	}
	rv, rCtl := b.Right.GetValue(ctx)
	if rCtl != nil {
		return nil, rCtl
	}
	if rv == nil {
		rv = data.NewNullValue()
	}
	if v, ok := rv.(data.Value); ok {
		switch l := b.Left.(type) {
		case *CallObjectDynamicProperty:
			return l.AssignValue(ctx, v)
		case *CallObjectProperty:
			return l.AssignValue(ctx, v)
		case *IndexExpression:
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
			// [$a, $b] = $arr：支持 ArrayValue 与 ObjectValue（关联数组，如 Collection::toArray()）
			var valueList []data.Value
			switch src := rv.(type) {
			case *data.ArrayValue:
				valueList = src.ToValueList()

			case *data.ClassValue:
				// 支持实现了 ArrayAccess 的对象解构（如 Collection）
				if method, exists := src.GetMethod("offsetGet"); exists {
					for i, set := range l.V {
						fnCtx := implicitMethodFrame(ctx, src, method)
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
					return v, nil
				}
			}
			for i, target := range l.V {
				if setter, ok := target.(data.Variable); ok {
					var value data.Value = data.NewNullValue()
					if i < len(valueList) && valueList[i] != nil {
						value = valueList[i]
					}
					if ctl := setter.SetValue(ctx, value); ctl != nil {
						return nil, ctl
					}
				}
			}
			return v, nil
		case *VarVar:
			// $$var = ... ：委托给 VarVar 自身的 SetValue
			if ctl := l.SetValue(ctx, v); ctl != nil {
				return nil, ctl
			}
			return v, nil
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
		// 标量按值写入局部变量槽，避免与 AST 字面量等共享 *IntValue 指针，
		// 以便 for 循环 VarStmtIncr 可安全原地自增。
		// 仅适用于 *VariableExpression；对象属性等 Variable 无有效索引。
		if ve, ok := b.Left.(*VariableExpression); ok {
			if zv := ctx.GetIndexZVal(ve.Index); zv != nil && data.IsScalarAssignFast(v) {
				if zv.Guard() != nil {
					prepared, ctl := zv.PrepareWrite(v, ctx)
					if ctl != nil {
						return nil, ctl
					}
					v = prepared
				}
				data.AssignScalarToZVal(zv, v)
				return zv.ReadValue(), nil
			}
		}
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
					fnCtx := implicitMethodFrame(ctx, cv, method)
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
