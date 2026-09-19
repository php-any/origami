package data

import (
	"math"
	"strconv"
	"strings"
)

// PrepareTypedValue 按 PHP 弱类型语义把值转成声明类型（int 参数接受 "2"）。
// 转换失败时返回 ok=false，调用方应抛出类型错误。
func PrepareTypedValue(ty Types, value Value) (Value, bool) {
	if value == nil {
		value = NewNullValue()
	}
	if zv, ok := value.(*ZValValue); ok {
		if zv.ZVal == nil {
			return value, ty == nil
		}
		inner := zv.ZVal.Value
		if inner == nil {
			inner = NewNullValue()
		}
		prepared, ok := PrepareTypedValue(ty, inner)
		if !ok {
			return nil, false
		}
		zv.ZVal.Value = prepared
		return value, true
	}
	if ty == nil {
		return value, true
	}
	if _, isNull := value.(*NullValue); isNull {
		return value, true
	}
	if ty.Is(value) {
		return value, true
	}
	return coerceToType(ty, value)
}

func coerceToType(ty Types, value Value) (Value, bool) {
	switch t := ty.(type) {
	case NullableType:
		return PrepareTypedValue(t.BaseType, value)
	case UnionType:
		for _, inner := range t.Types {
			if v, ok := PrepareTypedValue(inner, value); ok {
				return v, true
			}
		}
		return nil, false
	case IntersectionType:
		if t.Is(value) {
			return value, true
		}
		return nil, false
	case Int:
		return coerceToIntValue(value)
	case Float:
		return coerceToFloatValue(value)
	case String:
		return coerceToStringValue(value)
	default:
		return nil, false
	}
}

// PHP 弱类型：string 参数接受 int/float/bool，以及实现 __toString 的对象。
func coerceToStringValue(value Value) (Value, bool) {
	switch v := value.(type) {
	case *StringValue:
		return v, true
	case *IntValue:
		return NewStringValue(strconv.Itoa(v.Value)), true
	case *FloatValue:
		return NewStringValue(strconv.FormatFloat(v.Value, 'G', -1, 64)), true
	case *BoolValue:
		if v.Value {
			return NewStringValue("1"), true
		}
		return NewStringValue(""), true
	case *ThisValue:
		return objectToStringValue(v.ClassValue)
	case *ClassValue:
		return objectToStringValue(v)
	default:
		return nil, false
	}
}

func objectToStringValue(obj *ClassValue) (Value, bool) {
	if obj == nil {
		return nil, false
	}
	toStr, ok := obj.GetMethod("__toString")
	if !ok || toStr == nil {
		return nil, false
	}
	fnCtx := obj.CreateContext(toStr.GetVariables())
	fnCtx.SetCallArgs([]GetValue{})
	val, ctl := toStr.Call(fnCtx)
	if ctl != nil || val == nil {
		return nil, false
	}
	if s, ok := val.(*StringValue); ok {
		return s, true
	}
	if v, ok := val.(Value); ok {
		return NewStringValue(v.AsString()), true
	}
	return nil, false
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
		if v.Value == math.Trunc(v.Value) {
			return NewIntValue(int(v.Value)), true
		}
		return nil, false
	case *StringValue:
		s := strings.TrimSpace(v.Value)
		if s == "" {
			return nil, false
		}
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return NewIntValue(int(i)), true
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil && !math.IsNaN(f) && !math.IsInf(f, 0) && f == math.Trunc(f) {
			return NewIntValue(int(f)), true
		}
		return nil, false
	default:
		if iv, ok := value.(AsInt); ok {
			n, err := iv.AsInt()
			if err == nil {
				return NewIntValue(n), true
			}
		}
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
		s := strings.TrimSpace(v.Value)
		if s == "" {
			return nil, false
		}
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return NewFloatValue(f), true
		}
		return nil, false
	default:
		if fv, ok := value.(AsFloat); ok {
			f, err := fv.AsFloat()
			if err == nil {
				return NewFloatValue(f), true
			}
		}
		return nil, false
	}
}
