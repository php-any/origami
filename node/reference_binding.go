package node

import (
	"fmt"
	"github.com/php-any/origami/data"
)

func hasReferenceParameters(params []data.GetValue) bool {
	for _, parameter := range params {
		if p, ok := parameter.(interface{ IsReferenceParameter() bool }); ok && p.IsReferenceParameter() {
			return true
		}
	}
	return false
}

// This cold binder keeps expressions and unpacked buckets until the target
// parameter is known. Eager flattening loses references and evaluates lvalues
// twice. Ordinary positional calls retain their allocation-free fast binder.
func bindReferenceCall(frame, caller data.Context, params, args []data.GetValue, receiver data.Context) data.Control {
	byName := make(map[string]int, len(params))
	variadic := -1
	for i, p := range params {
		if name, ok := p.(data.GetName); ok {
			byName[name.GetName()] = i
		}
		if rest, ok := p.(interface{ IsVariadicParameter() bool }); ok && rest.IsVariadicParameter() {
			variadic = i
		}
	}
	bound := make([]bool, len(params))
	values := make([]data.Value, len(params))
	extra := make([]data.Value, 0)
	var rest *data.ArrayValue
	if variadic >= 0 {
		rest = data.NewArrayValue(nil).(*data.ArrayValue)
	}
	position, named := 0, false
	seenNames := make(map[string]bool)
	bind := func(name string, hasName bool, raw data.GetValue, slot *data.ZVal, value data.Value) data.Control {
		index := position
		if hasName {
			named = true
			if seenNames[name] {
				return data.NewErrorThrowByName(nil, fmt.Errorf("Named parameter $%s overwrites previous argument", name), "Error")
			}
			seenNames[name] = true
			var found bool
			index, found = byName[name]
			if !found {
				index = variadic
			}
			if index < 0 {
				return data.NewErrorThrowByName(nil, fmt.Errorf("Unknown named parameter $%s", name), "Error")
			}
		} else {
			if named {
				return data.NewErrorThrowByName(nil, fmt.Errorf("Cannot use positional argument after named argument"), "Error")
			}
			position++
		}
		if variadic >= 0 && index >= variadic {
			index = variadic
		}
		if index < len(params) && index >= 0 && bound[index] && index != variadic {
			return data.NewErrorThrowByName(nil, fmt.Errorf("Named parameter overwrites previous argument"), "Error")
		}
		var param data.GetValue
		if index >= 0 && index < len(params) {
			param = params[index]
		}
		isReference := false
		if p, ok := param.(interface{ IsReferenceParameter() bool }); ok {
			isReference = p.IsReferenceParameter()
		}
		if isReference {
			if slot == nil {
				var ctl data.Control
				slot, ctl = argumentReferenceSlot(caller, raw)
				if ctl != nil {
					return ctl
				}
			}
			value = slot.ReadValue()
		} else if value == nil {
			v, ctl := raw.GetValue(caller)
			if ctl != nil {
				return ctl
			}
			value, _ = v.(data.Value)
			if value == nil {
				value = data.NewNullValue()
			}
		}
		if param == nil {
			extra = append(extra, value)
			return nil
		}
		if index == variadic {
			var key data.Value = data.NewIntValue(rest.NextAppendIntKey())
			if hasName {
				key = data.NewStringValue(name)
			}
			if isReference {
				if !rest.BindReference(key, slot) {
					return data.NewErrorThrow(nil, fmt.Errorf("invalid variadic reference"))
				}
			} else {
				if ok, ctl := rest.AssignKey(frame, key, value); ctl != nil {
					return ctl
				} else if !ok {
					return data.NewTypeError(nil, fmt.Errorf("invalid argument key"))
				}
			}
			if !hasName {
				extra = append(extra, value)
			}
			return nil
		}
		if isReference {
			p, ok := param.(data.Variable)
			if !ok {
				return data.NewErrorThrowByName(nil, fmt.Errorf("Reference parameter cannot be bound"), "Error")
			}
			if ctl := p.SetValue(frame, data.NewZValValue(slot)); ctl != nil {
				return ctl
			}
			value = slot.ReadValue()
		} else if p, ok := param.(data.Variable); ok {
			if ctl := p.SetValue(frame, value); ctl != nil {
				return ctl
			}
			if promoted, ok := param.(*PromotedParameter); ok && receiver != nil {
				prepared, _ := frame.GetIndexValue(promoted.Index)
				if ctl := assignPromotedProperty(frame, receiver, promoted, prepared); ctl != nil {
					return ctl
				}
			}
		}
		bound[index], values[index] = true, value
		return nil
	}
	for _, argument := range args {
		if spread, ok := argument.(*SpreadArgument); ok {
			v, ctl := spread.Expr.GetValue(caller)
			if ctl != nil {
				return ctl
			}
			if array, ok := v.(*data.ArrayValue); ok {
				array = cowSeparateNestedArray(caller, spread.Expr, array).(*data.ArrayValue)
				for i, slot := range array.Range() {
					name := ""
					hasName := false
					if key, ok := slot.PHPArrayKey(i).(*data.StringValue); ok {
						name = key.Value
						hasName = true
					}
					if ctl := bind(name, hasName, data.NewZValValue(slot), slot, slot.ReadValue()); ctl != nil {
						return ctl
					}
				}
			} else {
				items, ctl := spreadToValues(caller, v)
				if ctl != nil {
					return ctl
				}
				for _, value := range items {
					if ctl := bind("", false, value, nil, value); ctl != nil {
						return ctl
					}
				}
			}
			continue
		}
		name := ""
		hasName := false
		if na, ok := argument.(*NamedArgument); ok {
			name, argument = na.Name, na.Value
			hasName = true
		}
		if ctl := bind(name, hasName, argument, nil, nil); ctl != nil {
			return ctl
		}
	}
	last := -1
	for i, param := range params {
		if i == variadic {
			if ctl := param.(data.Variable).SetValue(frame, rest); ctl != nil {
				return ctl
			}
			continue
		}
		if !bound[i] {
			v, ctl := param.GetValue(frame)
			if ctl != nil {
				return ctl
			}
			values[i], _ = v.(data.Value)
			if promoted, ok := param.(*PromotedParameter); ok && receiver != nil {
				if ctl := assignPromotedProperty(frame, receiver, promoted, values[i]); ctl != nil {
					return ctl
				}
			}
		} else {
			last = i
		}
	}
	flat := append(values[:last+1], extra...)
	if flat == nil {
		flat = []data.Value{}
	}
	frame.SetCallArgs(args)
	frame.SetFlatCallArgs(flat)
	return nil
}
