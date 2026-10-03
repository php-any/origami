package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func argumentReferenceSlot(ctx data.Context, argument data.GetValue) (*data.ZVal, data.Control) {
	if index, ok := argument.(*IndexExpression); ok {
		return index.GetOrCreateZVal(ctx)
	}
	if property, ok := argument.(interface {
		GetZVal(data.Context) (*data.ZVal, data.Control)
	}); ok {
		return property.GetZVal(ctx)
	}
	switch argument := argument.(type) {
	case data.Variable:
		return ctx.GetIndexZVal(argument.GetIndex()), nil
	default:
		var value data.GetValue
		var ctl data.Control
		if call, ok := argument.(referenceCall); ok {
			value, ctl = call.GetReferenceValue(ctx)
		} else {
			value, ctl = argument.GetValue(ctx)
		}
		if ctl != nil {
			return nil, ctl
		}
		if value, ok := value.(data.Value); ok {
			if slot, ctl := data.ReferenceSlot(value); slot != nil || ctl != nil {
				return slot, ctl
			}
		}
	}
	return nil, data.NewErrorThrowByName(nil, data.NewError(nil, "Only variables can be passed by reference", nil), "Error")
}

func bindVariadicReferences(frame, caller data.Context, parameter *ParametersReference, arguments []data.GetValue) data.Control {
	array := data.NewArrayValue(nil).(*data.ArrayValue)
	bind := func(key data.Value, slot *data.ZVal) data.Control {
		if parameter.Type != data.TypeInvalid {
			_, accepted, ctl := data.PrepareDeclaredValueInContext(parameter.Type, data.NewZValValue(slot), frame)
			if ctl != nil {
				return ctl
			}
			if !accepted {
				return data.NewTypeError(parameter.from, fmt.Errorf("invalid variadic argument type"))
			}
		}
		if !array.BindReference(key, slot) {
			return data.NewTypeError(parameter.from, fmt.Errorf("invalid variadic argument key"))
		}
		return nil
	}
	for _, argument := range arguments {
		var key data.Value = data.NewIntValue(array.NextAppendIntKey())
		if named, ok := argument.(*NamedArgument); ok {
			key = data.NewStringValue(named.Name)
			argument = named.Value
		}
		if spread, ok := argument.(*SpreadArgument); ok {
			value, ctl := spread.Expr.GetValue(caller)
			if ctl != nil {
				return ctl
			}
			source, ok := value.(*data.ArrayValue)
			if !ok {
				return data.NewTypeError(parameter.from, fmt.Errorf("Cannot unpack non-array by reference"))
			}
			source = cowSeparateNestedArray(caller, spread.Expr, source).(*data.ArrayValue)
			for i, slot := range source.Range() {
				key := slot.PHPArrayKey(i)
				if _, integer := key.(*data.IntValue); integer {
					key = data.NewIntValue(array.NextAppendIntKey())
				}
				if ctl := bind(key, slot); ctl != nil {
					return ctl
				}
			}
			continue
		}
		slot, ctl := argumentReferenceSlot(caller, argument)
		if ctl != nil {
			return ctl
		}
		if ctl := bind(key, slot); ctl != nil {
			return ctl
		}
	}
	return frame.SetVariableValue(parameter, array)
}
