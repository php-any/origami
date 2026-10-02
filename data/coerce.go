package data

import (
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

// PrepareTypedValueInContext checks declared values in the active PHP unit.
// Exact matches do not read the context flag; strict mode only affects coercion.
func PrepareTypedValueInContext(ty Types, value Value, ctx Context) (Value, bool, Control) {
	if value == nil {
		value = NewNullValue()
	}
	if zv, ok := value.(*ZValValue); ok {
		if zv.ZVal == nil {
			return value, ty == nil, nil
		}
		inner := zv.ZVal.Value
		if inner == nil {
			inner = NewNullValue()
		}
		prepared, ok, ctl := PrepareTypedValueInContext(ty, inner, ctx)
		if ctl != nil {
			return nil, false, ctl
		}
		if !ok {
			return nil, false, nil
		}
		zv.ZVal.Value = prepared
		return value, true, nil
	}
	if ty == nil {
		return value, true, nil
	}
	// Type declarations apply to the referenced value while preserving the
	// reference identity returned by a by-reference function.
	switch ref := value.(type) {
	case *ReferenceValue:
		inner, ctl := ref.Val.GetValue(ref.Ctx)
		if ctl != nil {
			return nil, false, ctl
		}
		prepared, accepted, ctl := PrepareTypedValueInContext(ty, inner.(Value), ctx)
		if ctl != nil || !accepted {
			return nil, accepted, ctl
		}
		if prepared != inner {
			if ctl := ref.Val.SetValue(ref.Ctx, prepared); ctl != nil {
				return nil, false, ctl
			}
		}
		return value, true, nil
	case *ArraySlotRef:
		if ref.Slot == nil {
			return nil, false, nil
		}
		prepared, accepted, ctl := PrepareTypedValueInContext(ty, ref.Slot.Value, ctx)
		if ctl != nil || !accepted {
			return nil, accepted, ctl
		}
		if prepared != ref.Slot.Value {
			CowAssign(ref.Slot, prepared)
		}
		return value, true, nil
	}
	if ref, ok := ty.(TypeRef); ok {
		return declarationTypes.Prepare(ref, value, ctx)
	}
	if ty.Is(value) {
		return value, true, nil
	}
	if ctx != nil && ctx.StrictTypes() {
		if integer, ok := value.(*IntValue); ok && hasFloatType(ty) {
			return NewFloatValue(float64(integer.Value)), true, nil
		}
		return nil, false, nil
	}
	return coerceToType(ty, value)
}

func hasFloatType(ty Types) bool {
	switch t := ty.(type) {
	case Float:
		return true
	case NullableType:
		return hasFloatType(t.BaseType)
	case UnionType:
		for _, member := range t.Types {
			if hasFloatType(member) {
				return true
			}
		}
	}
	return false
}

func coerceToType(ty Types, value Value) (Value, bool, Control) {
	switch t := ty.(type) {
	case NullableType:
		return PrepareTypedValueWithControl(t.BaseType, value)
	case UnionType:
		return coerceUnionValue(t, value)
	case IntersectionType:
		if t.Is(value) {
			return value, true, nil
		}
		return nil, false, nil
	case Int:
		prepared, ok := coerceToIntValue(value)
		return prepared, ok, nil
	case Float:
		prepared, ok := coerceToFloatValue(value)
		return prepared, ok, nil
	case Bool:
		switch value.(type) {
		case *IntValue, *FloatValue, *StringValue:
			b, err := value.(AsBool).AsBool()
			if err == nil {
				return NewBoolValue(b), true, nil
			}
		}
		return nil, false, nil
	case String:
		return coerceToStringValue(value)
	default:
		return nil, false, nil
	}
}

// PHP 弱类型：string 参数接受 int/float/bool，以及实现 __toString 的对象。
func coerceToStringValue(value Value) (Value, bool, Control) {
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
		return objectToStringValue(v.ClassValue)
	case *ClassValue:
		return objectToStringValue(v)
	default:
		return nil, false, nil
	}
}

func objectToStringValue(obj *ClassValue) (Value, bool, Control) {
	if obj == nil {
		return nil, false, nil
	}
	toStr, ok := obj.GetMethod("__toString")
	if !ok || toStr == nil {
		return nil, false, nil
	}
	fnCtx := obj.CreateContext(toStr.GetVariables())
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
		return coerceToStringValue(v.(Value))
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
func coerceUnionValue(ty UnionType, value Value) (Value, bool, Control) {
	var kinds uint8
	for _, member := range ty.Types {
		switch member.(type) {
		case Int:
			kinds |= 1
		case Float:
			kinds |= 2
		case String:
			kinds |= 4
		case Bool:
			kinds |= 8
		}
	}
	return coerceScalarKinds(kinds, value)
}

func coerceScalarKinds(kinds uint8, value Value) (Value, bool, Control) {
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
		if prepared, ok, ctl := coerceToStringValue(value); ctl != nil {
			return nil, false, ctl
		} else if ok {
			return prepared, true, nil
		}
	}
	if kinds&8 != 0 {
		if prepared, ok, ctl := coerceToType(Bool{}, value); ctl != nil {
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
