package data

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// PrepareTypedValue 按 PHP 弱类型语义把值转成声明类型（int 参数接受 "2"）。
// 转换失败时返回 ok=false，调用方应抛出类型错误。
func PrepareTypedValue(ty Types, value Value) (Value, bool) {
	prepared, ok, _ := PrepareTypedValueWithControl(ty, value)
	return prepared, ok
}

// PrepareTypedValueWithControl preserves exceptions raised by user conversion
// methods. Native callback binders use this weak entry point.
func PrepareTypedValueWithControl(ty Types, value Value) (Value, bool, Control) {
	return PrepareTypedValueInContext(ty, value, nil)
}

// PrepareDeclaredValueInContext keeps compact declarations on the direct
// runtime path. Storage wrappers retain their reference identity.
func PrepareDeclaredValueInContext(ref TypeRef, value Value, ctx Context) (Value, bool, Control) {
	if value == nil {
		value = NewNullValue()
	}
	switch value.(type) {
	case *ZValValue, *ReferenceValue, *ArraySlotRef:
		return PrepareTypedValueInContext(DeclaredType(ref), value, ctx)
	}
	return declarationTypes.Prepare(ref, value, ctx)
}

// PrepareTypedValueInContext checks declared values in the active PHP unit.
// Exact matches do not read the context flag; strict mode only affects coercion.
func PrepareTypedValueInContext(ty Types, value Value, ctx Context) (Value, bool, Control) {
	if value == nil {
		value = NewNullValue()
	}
	if source, ctl := ReferenceSlot(value); source != nil || ctl != nil {
		if ctl != nil {
			return nil, false, ctl
		}
		inner := source.ReadValue()
		if inner == nil {
			inner = NewNullValue()
		}
		prepared, accepted, ctl := PrepareTypedValueInContext(ty, inner, ctx)
		if ctl != nil || !accepted {
			return nil, accepted, ctl
		}
		if source.Guard() != nil {
			checked, guardControl := source.PrepareWrite(prepared, ctx)
			ctl = guardControl
			if ctl != nil {
				return nil, false, ctl
			}
			if !samePreparedReferenceValue(prepared, checked) {
				return nil, false, NewTypeError(nil, fmt.Errorf("Reference argument type is incompatible with its property constraints"))
			}
		}
		if prepared != inner {
			CowAssign(source, prepared)
		}
		return value, true, nil
	}
	// Legacy declarations are normalized only at this compatibility boundary.
	// Exact matching and weak/strict coercion use the same compact checker.
	return declarationTypes.Prepare(DeclaredTypeRef(ty), value, ctx)
}

func coerceToBoolValue(value Value) (Value, bool, Control) {
	switch value.(type) {
	case *IntValue, *FloatValue, *StringValue:
		b, err := value.(AsBool).AsBool()
		if err == nil {
			return NewBoolValue(b), true, nil
		}
	}
	return nil, false, nil
}

// PHP 弱类型：string 参数接受 int/float/bool，以及实现 __toString 的对象。
func coerceToStringValue(value Value, ctx Context) (Value, bool, Control) {
	switch v := value.(type) {
	case *StringValue:
		return v, true, nil
	case *IntValue:
		return NewStringValue(strconv.Itoa(v.Value)), true, nil
	case *FloatValue:
		return NewStringValue(strconv.FormatFloat(v.Value, 'G', -1, 64)), true, nil
	case *BoolValue:
		if v.Value {
			return NewStringValue("1"), true, nil
		}
		return NewStringValue(""), true, nil
	case *ThisValue:
		return ObjectToStringValue(v.ClassValue, ctx)
	case *ClassValue:
		return ObjectToStringValue(v, ctx)
	default:
		return nil, false, nil
	}
}

// ObjectToStringValue invokes __toString in its declaration scope and the
// caller's request. Explicit string contexts do not apply strict_types coercion.
func ObjectToStringValue(obj *ClassValue, ctx Context) (Value, bool, Control) {
	if obj == nil {
		return nil, false, nil
	}
	toStr, ok := obj.GetMethod("__toString")
	if !ok || toStr == nil {
		return nil, false, nil
	}
	if ctx == nil {
		ctx = obj.Context
	}
	vm := ctx.GetVM()
	if obj.GetVM() != vm {
		if provider, ok := vm.(RequestScopeProvider); ok {
			obj = provider.RequestObjectScope().Object(obj)
		}
	}
	self := MethodDeclaringClass(vm, obj.Class, "__toString")
	fnCtx := WrapMethodFrame(ctx.CreateContext(toStr.GetVariables()), obj, self, obj.Class)
	defer fnCtx.ReleasePooled()
	fnCtx.SetCallArgs([]GetValue{})
	val, ctl := toStr.Call(fnCtx)
	if ctl != nil {
		return nil, false, ctl
	}
	if val == nil {
		return nil, false, nil
	}
	if s, ok := val.(*StringValue); ok {
		return s, true, nil
	}
	switch v := val.(type) {
	case *IntValue, *FloatValue, *BoolValue:
		return coerceToStringValue(v.(Value), ctx)
	}
	return nil, false, nil
}

func coerceToIntValue(value Value) (Value, bool) {
	switch v := value.(type) {
	case *IntValue:
		return v, true
	case *BoolValue:
		if v.Value {
			return NewIntValue(1), true
		}
		return NewIntValue(0), true
	case *FloatValue:
		if math.IsNaN(v.Value) || math.IsInf(v.Value, 0) {
			return nil, false
		}
		limit := math.Ldexp(1, strconv.IntSize-1)
		if v.Value < -limit || v.Value >= limit {
			return nil, false
		}
		return NewIntValue(int(v.Value)), true
	case *StringValue:
		s, numeric := numericTypeText(v.Value)
		if !numeric {
			return nil, false
		}
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return NewIntValue(int(i)), true
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return coerceToIntValue(NewFloatValue(f))
		}
		return nil, false
	default:
		return nil, false
	}
}

func coerceToFloatValue(value Value) (Value, bool) {
	switch v := value.(type) {
	case *FloatValue:
		return v, true
	case *IntValue:
		return NewFloatValue(float64(v.Value)), true
	case *BoolValue:
		if v.Value {
			return NewFloatValue(1), true
		}
		return NewFloatValue(0), true
	case *StringValue:
		s, numeric := numericTypeText(v.Value)
		if !numeric {
			return nil, false
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return NewFloatValue(f), true
		}
		return nil, false
	default:
		return nil, false
	}
}

// PHP union coercion prefers int, float, string, bool independently of source
// declaration order. Numeric strings prefer float if both numeric members exist
// and the string has decimal/exponent syntax. Exact members were checked first.
func coerceScalarKinds(kinds uint8, value Value, ctx Context) (Value, bool, Control) {
	if kinds&3 == 3 {
		if str, ok := value.(*StringValue); ok {
			if text, numeric := numericTypeText(str.Value); numeric {
				if _, err := strconv.ParseInt(text, 10, 64); err != nil {
					if prepared, ok := coerceToFloatValue(value); ok {
						return prepared, true, nil
					}
				}
			}
		}
	}
	if kinds&1 != 0 {
		if prepared, ok := coerceToIntValue(value); ok {
			return prepared, true, nil
		}
	}
	if kinds&2 != 0 {
		if prepared, ok := coerceToFloatValue(value); ok {
			return prepared, true, nil
		}
	}
	if kinds&4 != 0 {
		if prepared, ok, ctl := coerceToStringValue(value, ctx); ctl != nil {
			return nil, false, ctl
		} else if ok {
			return prepared, true, nil
		}
	}
	if kinds&8 != 0 {
		if prepared, ok, ctl := coerceToBoolValue(value); ctl != nil {
			return nil, false, ctl
		} else if ok {
			return prepared, true, nil
		}
	}
	return nil, false, nil
}

func numericTypeText(text string) (string, bool) {
	text = strings.Trim(text, " \t\r\n\v\f")
	i := 0
	if i < len(text) && (text[i] == '+' || text[i] == '-') {
		i++
	}
	start := i
	for i < len(text) && text[i] >= '0' && text[i] <= '9' {
		i++
	}
	digits := i - start
	if i < len(text) && text[i] == '.' {
		i++
		start = i
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			i++
		}
		digits += i - start
	}
	if digits == 0 {
		return text, false
	}
	if i < len(text) && (text[i] == 'e' || text[i] == 'E') {
		i++
		if i < len(text) && (text[i] == '+' || text[i] == '-') {
			i++
		}
		start = i
		for i < len(text) && text[i] >= '0' && text[i] <= '9' {
			i++
		}
		if i == start {
			return text, false
		}
	}
	return text, i == len(text)
}
