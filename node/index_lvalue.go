package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

// Resolved only for complex dimensions. The retained slot is an address, not a
// PHP reference binding, so resolving a write does not introduce array aliases.
type resolvedDimensionOwner struct {
	slot            *data.ZVal
	overloadedClass string
}

func (r *resolvedDimensionOwner) GetValue(data.Context) (data.GetValue, data.Control) {
	return r.slot.ReadValue(), nil
}
func (r *resolvedDimensionOwner) GetZVal(data.Context) (*data.ZVal, data.Control) { return r.slot, nil }
func (r *resolvedDimensionOwner) GetIndex() int                                   { return -1 }
func (r *resolvedDimensionOwner) GetName() string                                 { return "dimension" }
func (r *resolvedDimensionOwner) GetType() data.Types                             { return nil }
func (r *resolvedDimensionOwner) SetValue(ctx data.Context, value data.Value) data.Control {
	if r.overloadedClass != "" {
		emitIndirectModificationNotice(nil, r.overloadedClass)
		return nil
	}
	prepared, ctl := r.slot.PrepareWrite(value, ctx)
	if ctl != nil {
		return ctl
	}
	data.CowAssign(r.slot, prepared)
	return nil
}

func resolveDimensionOwner(ctx data.Context, expression data.GetValue) (data.GetValue, data.Control) {
	return resolveDimensionOwnerMode(ctx, expression, true)
}
func resolveDimensionOwnerMode(ctx data.Context, expression data.GetValue, create bool) (data.GetValue, data.Control) {
	switch expr := expression.(type) {
	case *IndexExpression:
		resolved, ctl := resolveDimensionMode(ctx, expr, create)
		if ctl != nil {
			return nil, ctl
		}
		parent, ctl := resolved.Array.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		parent = unwrapIndexAssignTarget(parent)
		var object *data.ClassValue
		switch value := parent.(type) {
		case *data.ClassValue:
			object = value
		case *data.ThisValue:
			object = value.ClassValue
		}
		if object != nil && checkArrayAccess(ctx, object.Class) {
			key, ctl := resolved.Index.GetValue(ctx)
			if ctl != nil {
				return nil, ctl
			}
			value, ctl := callArrayAccessOffsetGetReference(ctx, object, key)
			if ctl != nil {
				return nil, ctl
			}
			if value, ok := value.(data.Value); ok {
				if slot, ctl := data.ReferenceSlot(value); slot != nil || ctl != nil {
					if ctl != nil {
						return nil, ctl
					}
					return &resolvedDimensionOwner{slot: slot}, nil
				}
				return &resolvedDimensionOwner{slot: data.NewZVal(value), overloadedClass: object.Class.GetName()}, nil
			}
			return &resolvedDimensionOwner{slot: data.NewZVal(data.NewNullValue())}, nil
		}
		if !create {
			if array, ok := parent.(*data.ArrayValue); ok {
				array = cowSeparateNestedArray(ctx, resolved.Array, array).(*data.ArrayValue)
				key, ctl := resolved.Index.GetValue(ctx)
				if ctl != nil {
					return nil, ctl
				}
				if slot, ok := array.LookupZValByStringKey(key.(data.Value).AsString()); ok {
					return &resolvedDimensionOwner{slot: slot}, nil
				}
			}
			return &resolvedDimensionOwner{slot: data.NewZVal(data.NewNullValue())}, nil
		}
		if _, null := parent.(*data.NullValue); null {
			array := data.NewArrayValueFromSlots(nil)
			if setter, ok := resolved.Array.(data.Variable); ok {
				if ctl := setter.SetValue(ctx, array); ctl != nil {
					return nil, ctl
				}
			}
		}
		slot, ctl := resolved.getOrCreateZValResolved(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &resolvedDimensionOwner{slot: slot}, nil
	case *CallObjectProperty:
		object, ctl := expr.Object.GetValue(ctx)
		if ctl != nil {
			return nil, ctl
		}
		if value, ok := object.(*data.ThisValue); ok {
			object = value.ClassValue
		}
		if value, ok := object.(*data.ClassValue); ok {
			if property, found := lookupObjectProperty(ctx, value, expr.Property); found {
				if readonly, ok := property.(*ClassProperty); ok && readonly.IsReadonly {
					current, ctl := property.GetValue(value)
					if ctl != nil {
						return nil, ctl
					}
					if _, object := current.(*data.ClassValue); object {
						return &literalGetValue{v: current}, nil
					}
					return nil, data.NewErrorThrowByName(expr.GetFrom(), fmt.Errorf("Cannot modify readonly property %s::$%s", readonly.DeclaringClass, expr.Property), "Error")
				}
			}
		}
		proxy := &CallObjectProperty{Node: expr.Node, Object: &literalGetValue{v: object}, Property: expr.Property}
		if value, ok := object.(*data.ClassValue); ok {
			if property, found := lookupObjectProperty(ctx, value, expr.Property); found && propertyAccessible(ctx, value, property, expr.Property) {
				if !value.ObjectValue.HasProperty(data.PropertyStorageName(property)) {
					if !create {
						return &literalGetValue{v: data.NewNullValue()}, nil
					}
					if _, ctl := proxy.AssignValue(ctx, data.NewArrayValueFromSlots(nil)); ctl != nil {
						return nil, ctl
					}
				}
				slot, ctl := property.GetZVal(value)
				if ctl != nil {
					return nil, ctl
				}
				return &resolvedDimensionOwner{slot: slot}, nil
			}
		}
		return proxy, nil
	case *CallObjectDynamicProperty:
		frozen, ctl := freezeLvalue(ctx, expr)
		if ctl != nil {
			return nil, ctl
		}
		return resolveDimensionOwnerMode(ctx, frozen, create)
	case *CallStaticProperty, *CallStaticPropertyLater, *CallStaticKeywordProperty, *CallSelfProperty:
		// Static dimensions must retain the property's actual storage slot.
		// A temporary value loses writes as soon as COW separates the array.
		frozen, ctl := freezeLvalue(ctx, expr)
		if ctl != nil {
			return nil, ctl
		}
		slot, ctl := frozen.(interface {
			GetZVal(data.Context) (*data.ZVal, data.Control)
		}).GetZVal(ctx)
		if ctl != nil {
			return nil, ctl
		}
		return &resolvedDimensionOwner{slot: slot}, nil
	case data.Variable:
		return expr, nil
	default:
		var value data.GetValue
		var ctl data.Control
		if call, ok := expression.(referenceCall); ok {
			value, ctl = call.GetReferenceValue(ctx)
		} else {
			value, ctl = expression.GetValue(ctx)
		}
		if ctl != nil {
			return nil, ctl
		}
		if value, ok := value.(data.Value); ok {
			if slot, ctl := data.ReferenceSlot(value); slot != nil || ctl != nil {
				if ctl != nil {
					return nil, ctl
				}
				return &resolvedDimensionOwner{slot: slot}, nil
			}
			return &resolvedDimensionOwner{slot: data.NewZVal(value)}, nil
		}
		return nil, data.NewTypeError(nil, fmt.Errorf("Cannot resolve dimension owner %T", value))
	}
}

func resolveDimension(ctx data.Context, expression *IndexExpression) (*IndexExpression, data.Control) {
	return resolveDimensionMode(ctx, expression, true)
}
func resolveDimensionMode(ctx data.Context, expression *IndexExpression, create bool) (*IndexExpression, data.Control) {
	owner, ctl := resolveDimensionOwnerMode(ctx, expression.Array, create)
	if ctl != nil {
		return nil, ctl
	}
	key, ctl := expression.Index.GetValue(ctx)
	if ctl != nil {
		return nil, ctl
	}
	return &IndexExpression{Node: expression.Node, Array: owner, Index: &literalGetValue{v: key}, Append: expression.Append}, nil
}

func complexDimension(expression *IndexExpression) bool {
	if _, simple := expression.Array.(*VariableExpression); !simple {
		return true
	}
	switch expression.Index.(type) {
	case *StringLiteral, *IntLiteral, *literalGetValue:
		return false
	}
	return true
}

// Only compound dimension assignments need an operand replacement. Cached AST
// nodes stay immutable; the copy belongs to this execution.
func rebindDimensionOperand(expression, old, replacement data.GetValue) (data.GetValue, bool) {
	switch operation := expression.(type) {
	case *BinaryAdd:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinarySub:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryMul:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryQuo:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryRem:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryDot:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryPow:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryBitAnd:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryBitOr:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryBitXor:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryShl:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	case *BinaryShr:
		if operation.Left == old {
			copy := *operation
			copy.Left = replacement
			return &copy, true
		}
	}
	return expression, false
}
