package node

import (
	"errors"
	"fmt"
	"strconv"

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
			if ctl := assignListTargets(b.from, ctx, l.V, v); ctl != nil {
				return nil, ctl
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

// assignListTargets 把右值按位置分发给解构目标列表，是 [$a, $b] = / list($a, $b) = 的公共实现。
//
// 右值可以是数组，也可以是实现了 ArrayAccess 的对象（如 Collection，按 0..n-1 取 offsetGet）。
func assignListTargets(from data.From, ctx data.Context, targets []data.GetValue, value data.Value) data.Control {
	var source *data.ArrayValue
	switch src := value.(type) {
	case *data.ArrayValue:
		source = src

	case *data.ClassValue:
		// 支持实现了 ArrayAccess 的对象解构（如 Collection）
		if method, exists := src.GetMethod("offsetGet"); exists {
			for i, set := range targets {
				if _, skip := set.(*data.NullValue); skip {
					continue
				}
				fnCtx := implicitMethodFrame(ctx, src, method)
				if len(method.GetVariables()) > 0 {
					fnCtx.SetVariableValue(method.GetVariables()[0], data.NewIntValue(i))
				}
				ret, ctl := method.Call(fnCtx)
				if ctl != nil {
					return ctl
				}
				var val data.Value = data.NewNullValue()
				if rv2, ok2 := ret.(data.Value); ok2 {
					val = rv2
				}
				if ctl := assignListTarget(from, ctx, set, val); ctl != nil {
					return ctl
				}
			}
			return nil
		}
	}
	for i, target := range targets {
		if _, skip := target.(*data.NullValue); skip {
			continue
		}
		var item data.Value = data.NewNullValue()
		if source != nil {
			if slot, _ := source.FindSlotByIntKey(i); slot != nil {
				item = slot.ReadValue()
			} else if ctl := emitUndefinedArrayKeyWarning(ctx, from, strconv.Itoa(i), true); ctl != nil {
				return ctl
			}
		}
		if ctl := assignListTarget(from, ctx, target, item); ctl != nil {
			return ctl
		}
	}
	return nil
}

// assignListTarget 把单个解构目标写成 value。
//
// PHP 的 list() / [$a, $b] 左侧可以是任意可写左值：变量、超全局、下标（$o->arr[0]、$a[0][]）、
// 对象属性（$o->p）、静态属性，以及嵌套解构。这里逐类分派；此前只处理 data.Variable，
// 其余类型被**静默跳过**（既不赋值也不报错），因此 [$child->prefixes[0], ...] = [...] 这类
// 合法写法在整个解释器里无声失效。
func assignListTarget(from data.From, ctx data.Context, target data.GetValue, value data.Value) data.Control {
	// data.Variable 覆盖普通变量、超全局、类属性、可变变量等既有路径，行为保持不变。
	if setter, ok := target.(data.Variable); ok {
		return setter.SetValue(ctx, value)
	}
	switch t := target.(type) {
	case *data.NullValue:
		// 解析器用 NullValue 占位 list 解构中的空槽位：[$a, , $b] = [...]。
		// 该位置的值被丢弃，不做写入。
		return nil
	case *IndexExpression:
		// 下标左值自带写路径（含复杂下标的 resolveDimension/COW 处理）。
		// *IndexExpression 缺少 GetIndex/GetName/GetType，不满足 data.Variable，必须单独分派。
		return t.SetValue(ctx, value)
	case *CallObjectDynamicProperty:
		_, ctl := t.AssignValue(ctx, value)
		return ctl
	case *CallStaticProperty:
		return t.SetProperty(ctx, t.Property, value)
	case *CallStaticPropertyLater:
		return t.SetProperty(ctx, t.property, value)
	case *CallStaticKeywordProperty:
		return t.SetProperty(ctx, t.Property, value)
	case *CallSelfProperty:
		return t.SetProperty(ctx, t.Property, value)
	case *VarVar:
		return t.SetValue(ctx, value)
	case *Array:
		// PHP 支持嵌套解构：[$a, [$b, $c]] = [1, [2, 3]]
		return assignListTargets(from, ctx, t.V, value)
	}
	return data.NewErrorThrow(from, fmt.Errorf("解构赋值的目标不可写: %T", target))
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
