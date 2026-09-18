package spl

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/std/php/core"
)

// splGetClassValue ??ctx ???? ClassValue?ClassMethodContext ??ClassValue??
func splGetClassValue(ctx data.Context) *data.ClassValue {
	if cmc, ok := ctx.(*data.ClassMethodContext); ok {
		return cmc.ClassValue
	}
	if cv, ok := ctx.(*data.ClassValue); ok {
		return cv
	}
	return nil
}

// splCallUserCallback ?? call_user_func ??????
func splCallUserCallback(ctx data.Context, cb data.GetValue, args ...data.Value) (data.GetValue, data.Control) {
	fn := core.NewCallUserFuncFunction()
	callCtx := ctx.CreateContext(fn.GetVariables())
	if v, ok := cb.(data.Value); ok {
		callCtx.SetIndexZVal(0, data.NewZVal(v))
	}
	switch len(args) {
	case 0:
	case 1:
		callCtx.SetIndexZVal(1, data.NewZVal(args[0]))
	default:
		callCtx.SetIndexZVal(1, data.NewZVal(data.NewArrayValue(args)))
	}
	return fn.Call(callCtx)
}

// splValueTruthy ?? PHP truthy
func splValueTruthy(v data.GetValue) bool {
	if v == nil {
		return false
	}
	if val, ok := v.(data.Value); ok {
		if _, isNull := val.(*data.NullValue); isNull {
			return false
		}
		if bv, ok := val.(*data.BoolValue); ok {
			return bv.Value
		}
		if ab, ok := val.(data.AsBool); ok {
			b, err := ab.AsBool()
			return err == nil && b
		}
		return true
	}
	return true
}

// splAsValue ??GetValue ?? Value???????? null??
func splAsValue(v data.GetValue) data.Value {
	if val, ok := v.(data.Value); ok {
		return val
	}
	return data.NewNullValue()
}

// splAsInt ??Value ?? int
func splAsInt(v data.GetValue) int {
	if v == nil {
		return 0
	}
	if iv, ok := v.(data.AsInt); ok {
		if n, err := iv.AsInt(); err == nil {
			return n
		}
	}
	if iv64, ok := v.(interface{ AsInt64() int64 }); ok {
		return int(iv64.AsInt64())
	}
	return 0
}
