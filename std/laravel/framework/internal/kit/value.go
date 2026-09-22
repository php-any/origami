package kit

import (
	"fmt"
	"strconv"

	"github.com/php-any/origami/data"
)

// Receiver 从 ClassMethodContext / ClassValue 取出 $this。
func Receiver(ctx data.Context) *data.ClassValue {
	if methodCtx, ok := ctx.(*data.ClassMethodContext); ok {
		return methodCtx.ClassValue
	}
	if value, ok := ctx.(*data.ClassValue); ok {
		return value
	}
	return nil
}

// Arg 取位置参数（无则 nil）。
func Arg(ctx data.Context, i int) data.Value {
	if ctx == nil {
		return nil
	}
	v, _ := ctx.GetIndexValue(i)
	return v
}

// IsNull 判断 null。
func IsNull(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

// Truthy 对齐 PHP 宽松真值（空串 / "0" 为假）。
func Truthy(v data.Value) bool {
	if v == nil || IsNull(v) {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv
	}
	s := v.AsString()
	return s != "" && s != "0"
}

// KV 数组/对象条目。
type KV struct {
	Key    data.Value
	KeyStr string
	Value  data.Value
}

// Entries 把 ArrayValue / ObjectValue 展成有序条目。
func Entries(v data.Value) []KV {
	if v == nil {
		return nil
	}
	v = Unwrap(v)
	switch arr := v.(type) {
	case *data.ArrayValue:
		entries := make([]KV, 0, len(arr.List))
		for i, z := range arr.List {
			if z == nil {
				continue
			}
			keyStr := z.Name
			var key data.Value
			if keyStr != "" {
				if n, ok := data.ParseIntArrayKeyName(keyStr); ok {
					key = data.NewIntValue(n)
				} else {
					key = data.NewStringValue(keyStr)
				}
			} else {
				key = data.NewIntValue(i)
				keyStr = data.IntArrayKeyName(i)
			}
			entries = append(entries, KV{Key: key, KeyStr: keyStr, Value: z.Value})
		}
		return entries
	case *data.ObjectValue:
		entries := make([]KV, 0)
		arr.RangeProperties(func(key string, value data.Value) bool {
			entries = append(entries, KV{Key: data.NewStringValue(key), KeyStr: key, Value: value})
			return true
		})
		return entries
	default:
		return nil
	}
}

// Unwrap 剥 ZValValue / ThisValue。
func Unwrap(v data.Value) data.Value {
	for n := 0; n < 4 && v != nil; n++ {
		switch t := v.(type) {
		case *data.ZValValue:
			if t.ZVal == nil {
				return v
			}
			v = t.ZVal.Value
		case *data.ThisValue:
			if t.ClassValue == nil {
				return v
			}
			v = t.ClassValue
		default:
			return v
		}
	}
	return v
}

// Call 调用闭包 / 绑定方法 / 函数名字符串。
func Call(ctx data.Context, cb data.Value, args ...data.Value) (data.GetValue, data.Control) {
	if cb == nil {
		return nil, nil
	}
	var fn data.FuncStmt
	switch c := cb.(type) {
	case *data.FuncValue:
		fn = c.Value
	case *data.BoundFuncValue:
		fn = c.Value
	case *data.StringValue:
		if f, ok := ctx.GetVM().GetFunc(c.AsString()); ok {
			fn = f
		}
	default:
		return nil, data.NewErrorThrow(nil, fmt.Errorf("callback is not callable"))
	}
	if fn == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("callback is not callable"))
	}
	callCtx := ctx.CreateContext(fn.GetVariables())
	data.BindDeclaredArgs(callCtx, fn, args)
	if bfv, ok := cb.(*data.BoundFuncValue); ok {
		return bfv.Call(callCtx)
	}
	return fn.Call(callCtx)
}

// KeyString 把键值转成字符串键名。
func KeyString(v data.Value) string {
	if v == nil {
		return ""
	}
	if iv, ok := v.(*data.IntValue); ok {
		return strconv.Itoa(iv.Value)
	}
	return v.AsString()
}
