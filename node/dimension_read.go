package node

import (
	"errors"

	"github.com/php-any/origami/data"
)

// ReadDimensionQuiet evaluates every receiver/key once without creating slots.
// probe uses isset/offsetExists semantics; fetch reads the final value only when
// required by empty or coalesce. Parent dimensions always need their values.
func ReadDimensionQuiet(ctx data.Context, index *IndexExpression, probe, fetch bool) (data.GetValue, bool, data.Control) {
	var container data.GetValue
	var ctl data.Control
	if parent, ok := index.Array.(*IndexExpression); ok {
		container, _, ctl = ReadDimensionQuiet(ctx, parent, probe, true)
	} else if _, dynamic := index.Array.(*CallObjectDynamicProperty); dynamic {
		frozen, control := freezeLvalue(ctx, index.Array)
		if control != nil {
			return nil, false, control
		}
		copy := *index
		copy.Array = frozen
		return ReadDimensionQuiet(ctx, &copy, probe, fetch)
	} else if property, ok := index.Array.(*CallObjectProperty); ok {
		object, control := property.Object.GetValue(ctx)
		if control != nil {
			return nil, false, control
		}
		frozen := &CallObjectProperty{Node: property.Node, Object: &literalGetValue{v: object}, Property: property.Property}
		if value, exists, control := probeObjectProperty(ctx, frozen); control != nil {
			return nil, false, control
		} else if !exists {
			// PHP：没有 __isset 时，isset/empty/?? 的属性下标会调用 __get 再检查键。
			// isset($obj->prop) 本身仍然不调用 __get。有 __isset 且返回 false 时也不调用。
			value, use, control := magicPropertyForDimension(ctx, frozen)
			if control != nil {
				return nil, false, control
			}
			if use {
				container = value
			}
			if container == nil {
				container = data.NewNullValue()
			}
		} else if value != nil {
			container = value
		} else {
			container, ctl = frozen.GetValue(ctx)
		}
	} else {
		container, ctl = index.Array.GetValue(ctx)
	}
	if ctl != nil {
		return nil, false, ctl
	}
	key, ctl := index.Index.GetValue(ctx)
	if ctl != nil {
		return nil, false, ctl
	}
	// Resolve and validate keys in one dispatch on the ordinary array path.
	if array, ok := container.(*data.ArrayValue); ok {
		var slot *data.ZVal
		switch key := key.(type) {
		case *data.StringValue:
			slot, _ = array.LookupZValByStringKey(key.Value)
		case *data.NullValue:
			slot, _ = array.LookupZValByStringKey("")
		case data.AsInt:
			integer, err := key.AsInt()
			if err != nil {
				return nil, false, data.NewErrorThrowByName(index.GetFrom(), err, "TypeError")
			}
			slot, _ = array.FindSlotByIntKey(integer)
		default:
			return nil, false, data.NewErrorThrowByName(index.GetFrom(), errors.New("Illegal offset type"), "TypeError")
		}
		if slot == nil {
			return data.NewNullValue(), false, nil
		}
		value := slot.ReadValue()
		if value == nil {
			value = data.NewNullValue()
		}
		return value, !probe || issetNonNullValue(value), nil
	}
	var object *data.ClassValue
	switch value := container.(type) {
	case *data.ClassValue:
		object = value
	case *data.ThisValue:
		object = value.ClassValue
	}
	if object != nil && checkArrayAccess(ctx, object.Class) {
		if probe {
			exists, ctl := callArrayAccessOffsetExists(ctx, object, key.(data.Value))
			if ctl != nil || !exists || !fetch {
				return data.NewNullValue(), exists, ctl
			}
		}
		value, ctl := callArrayAccessOffsetGet(ctx, object, key)
		return value, true, ctl
	}
	if object != nil {
		return nil, false, data.NewErrorThrowByName(index.GetFrom(), errors.New("Cannot use object as array"), "Error")
	}
	value, exists := readIndexNoWarn(container, key)
	if !exists {
		// Strings support quiet byte offsets too.
		if text, ok := container.(*data.StringValue); ok {
			if integer, ok := key.(data.AsInt); ok {
				offset, err := integer.AsInt()
				if offset < 0 {
					offset += len(text.Value)
				}
				if err == nil && offset >= 0 && offset < len(text.Value) {
					return data.NewStringValue(text.Value[offset : offset+1]), true, nil
				}
			}
		}
		return data.NewNullValue(), false, nil
	}
	if probe && !issetNonNullValue(value.(data.Value)) {
		return value, false, nil
	}
	return value, true, nil
}

// magicPropertyForDimension 在属性下标的 isset/empty/?? 中补上 PHP 的 __get。
// 仅当属性未声明、没有动态值、也没有 __isset 时读取；__isset 返回 false 时保持不调用 __get。
func magicPropertyForDimension(ctx data.Context, property *CallObjectProperty) (data.GetValue, bool, data.Control) {
	owner, ctl := property.Object.GetValue(ctx)
	if ctl != nil {
		return nil, false, ctl
	}
	var object *data.ClassValue
	switch value := owner.(type) {
	case *data.ClassValue:
		object = value
	case *data.ThisValue:
		object = value.ClassValue
	default:
		return nil, false, nil
	}
	if object == nil || object.ObjectValue == nil {
		return nil, false, nil
	}
	if _, declared := lookupObjectProperty(ctx, object, property.Property); declared {
		return nil, false, nil
	}
	if object.ObjectValue.HasProperty(property.Property) {
		return nil, false, nil
	}
	if _, found := object.GetMethod("__isset"); found {
		return nil, false, nil
	}
	if _, found := object.GetMethod("__get"); !found {
		return nil, false, nil
	}
	value, ctl := property.GetValue(ctx)
	if ctl != nil {
		return nil, false, ctl
	}
	return value, true, nil
}

// FreezeLvalue retains expression results rather than mutating the cached AST.
// Variable addresses are kept for a later conditional write.
func freezeLvalue(ctx data.Context, expression data.GetValue) (data.GetValue, data.Control) {
	switch value := expression.(type) {
	case *IndexExpression:
		owner, ctl := freezeLvalue(ctx, value.Array)
		if ctl != nil {
			return nil, ctl
		}
		key, ctl := value.Index.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &IndexExpression{Node: value.Node, Array: owner, Index: &literalGetValue{v: key}, Append: value.Append}, nil
	case *CallObjectDynamicProperty:
		object, ctl := value.Object.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		name, ctl := value.NameExpr.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &CallObjectProperty{Node: value.Node, Object: &literalGetValue{v: object}, Property: name.(data.Value).AsString()}, nil
	case *CallObjectProperty:
		object, ctl := value.Object.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &CallObjectProperty{Node: value.Node, Object: &literalGetValue{v: object}, Property: value.Property}, nil
	case *CallStaticKeywordProperty, *CallStaticPropertyLater, *CallSelfProperty:
		return value, nil
	case *CallStaticProperty:
		if _, class := value.Stmt.(data.ClassStmt); class {
			return value, nil
		}
		owner, ctl := value.Stmt.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &CallStaticProperty{Node: value.Node, Stmt: &literalGetValue{v: owner}, Property: value.Property}, nil
	case data.Variable:
		return value, nil
	default:
		result, ctl := expression.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &literalGetValue{v: result}, nil
	}
}

type NullCoalesceAssign struct {
	*Node
	Left, Right data.GetValue
}

func (n *NullCoalesceAssign) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	left, ctl := freezeLvalue(ctx, n.Left)
	if ctl != nil {
		return nil, ctl
	}
	value, ctl := (&NullCoalesceExpression{Node: n.Node, Left: left, Right: &literalGetValue{v: data.NewNullValue()}}).GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	if value != nil {
		if _, null := value.(*data.NullValue); !null {
			return value, nil
		}
	}
	return NewBinaryAssign(n.GetFrom(), left, n.Right).GetValue(ctx)
}
