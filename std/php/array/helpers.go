package array

import (
	"github.com/php-any/origami/data"
)

type kvEntry struct {
	key    data.Value
	keyStr string
	value  data.Value
}

func toKVEntries(v data.Value) []kvEntry {
	if v == nil {
		return nil
	}
	switch arr := v.(type) {
	case *data.ArrayValue:
		entries := make([]kvEntry, 0, len(arr.List))
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
			entries = append(entries, kvEntry{key: key, keyStr: keyStr, value: z.Value})
		}
		return entries
	case *data.ObjectValue:
		entries := make([]kvEntry, 0)
		arr.RangeProperties(func(key string, value data.Value) bool {
			entries = append(entries, kvEntry{
				key:    data.NewStringValue(key),
				keyStr: key,
				value:  value,
			})
			return true
		})
		return entries
	default:
		return nil
	}
}

func getElementValue(row data.Value, key data.Value) data.Value {
	if row == nil || key == nil {
		return data.NewNullValue()
	}
	keyStr := key.AsString()
	switch r := row.(type) {
	case *data.ArrayValue:
		if iv, ok := key.(data.AsInt); ok {
			if i, err := iv.AsInt(); err == nil {
				if z, _ := r.FindSlotByIntKey(i); z != nil {
					return z.Value
				}
			}
		}
		for _, z := range r.List {
			if z != nil && z.Name == keyStr {
				return z.Value
			}
		}
		if iv, ok := key.(data.AsInt); ok {
			if i, err := iv.AsInt(); err == nil && i >= 0 && i < len(r.List) {
				if z := r.List[i]; z != nil {
					return z.Value
				}
			}
		}
	case *data.ObjectValue:
		if val, ctl := r.GetProperty(keyStr); ctl == nil && val != nil {
			return val
		}
	case *data.ClassValue:
		if val, ctl := r.GetProperty(keyStr); ctl == nil && val != nil {
			return val
		}
	}
	return data.NewNullValue()
}

func isNestedArrayValue(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func invokeCallback(ctx data.Context, cb data.Value, args []data.Value) (data.Value, data.Control) {
	if cb == nil {
		return data.NewNullValue(), nil
	}
	switch c := cb.(type) {
	case *data.FuncValue:
		vars := c.Value.GetVariables()
		fnCtx := ctx.CreateContext(vars)
		for i := 0; i < len(vars) && i < len(args); i++ {
			fnCtx.SetVariableValue(data.NewVariable("", i, nil), args[i])
		}
		ret, ctl := c.Value.Call(fnCtx)
		if ctl != nil {
			return nil, ctl
		}
		if v, ok := ret.(data.Value); ok {
			return v, nil
		}
		return data.NewNullValue(), nil
	case *data.ArrayValue:
		if len(c.List) == 2 {
			objVal := c.List[0].Value
			methodVal := c.List[1].Value
			if obj, ok := objVal.(data.GetMethod); ok {
				methodName := methodVal.AsString()
				if method, has := obj.GetMethod(methodName); has {
					vars := method.GetVariables()
					fnCtx := ctx.CreateContext(vars)
					for i := 0; i < len(vars) && i < len(args); i++ {
						fnCtx.SetVariableValue(vars[i], args[i])
					}
					ret, ctl := method.Call(fnCtx)
					if ctl != nil {
						return nil, ctl
					}
					if v, ok := ret.(data.Value); ok {
						return v, nil
					}
				}
			}
		}
	case data.CallableValue:
		var a0, a1, a2 data.Value = data.NewNullValue(), data.NewNullValue(), data.NewNullValue()
		if len(args) > 0 {
			a0 = args[0]
		}
		if len(args) > 1 {
			a1 = args[1]
		}
		if len(args) > 2 {
			a2 = args[2]
		}
		ret, ctl := c.Call(a0, a1, a2)
		if ctl != nil {
			return nil, ctl
		}
		return ret, nil
	default:
		funcName := cb.AsString()
		if fnStmt, exists := ctx.GetVM().GetFunc(funcName); exists {
			fnValue := data.NewFuncValue(fnStmt)
			vars := fnValue.Value.GetVariables()
			fnCtx := ctx.CreateContext(vars)
			for i := 0; i < len(vars) && i < len(args); i++ {
				fnCtx.SetVariableValue(data.NewVariable("", i, nil), args[i])
			}
			ret, ctl := fnValue.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if v, ok := ret.(data.Value); ok {
				return v, nil
			}
		}
	}
	return data.NewNullValue(), nil
}

func compareCallback(ctx data.Context, cb data.Value, a, b data.Value) int {
	ret, ctl := invokeCallback(ctx, cb, []data.Value{a, b})
	if ctl != nil {
		return 0
	}
	if iv, ok := ret.(data.AsInt); ok {
		i, _ := iv.AsInt()
		return i
	}
	if fv, ok := ret.(data.AsFloat); ok {
		f, _ := fv.AsFloat()
		if f < 0 {
			return -1
		}
		if f > 0 {
			return 1
		}
		return 0
	}
	if a.AsString() < b.AsString() {
		return -1
	}
	if a.AsString() > b.AsString() {
		return 1
	}
	return 0
}

func toFloatValue(v data.Value) float64 {
	if v == nil {
		return 0
	}
	if iv, ok := v.(data.AsInt); ok {
		i, err := iv.AsInt()
		if err == nil {
			return float64(i)
		}
	}
	if fv, ok := v.(data.AsFloat); ok {
		f, err := fv.AsFloat()
		if err == nil {
			return f
		}
	}
	return 0
}

func valuesEqual(a, b data.Value) bool {
	return a.AsString() == b.AsString()
}

func buildResultFromEntries(entries []kvEntry, preserveKeys bool) data.Value {
	if len(entries) == 0 {
		return data.NewArrayValue([]data.Value{})
	}
	hasStringKey := false
	for _, e := range entries {
		if _, ok := e.key.(*data.StringValue); ok {
			if e.keyStr != "" && !isNumericKeyStr(e.keyStr) {
				hasStringKey = true
				break
			}
		}
	}
	if hasStringKey || preserveKeys {
		result := data.NewObjectValue()
		for _, e := range entries {
			result.SetProperty(e.keyStr, e.value)
		}
		return result
	}
	if preserveKeys {
		list := make([]*data.ZVal, 0, len(entries))
		for _, e := range entries {
			if _, ok := e.key.(*data.IntValue); ok && e.keyStr == data.IntArrayKeyName(e.key.(*data.IntValue).Value) {
				if n, ok := data.ParseIntArrayKeyName(e.keyStr); ok && n != len(list) {
					list = append(list, data.NewNamedZVal(e.keyStr, e.value))
					continue
				}
			}
			list = append(list, data.NewNamedZVal(e.keyStr, e.value))
		}
		return &data.ArrayValue{List: list}
	}
	vals := make([]data.Value, len(entries))
	for i, e := range entries {
		vals[i] = e.value
	}
	return data.NewArrayValue(vals)
}

func isNumericKeyStr(s string) bool {
	_, ok := data.ParseIntArrayKeyName(s)
	return ok
}

func isNullValue(v data.Value) bool {
	if v == nil {
		return true
	}
	_, ok := v.(*data.NullValue)
	return ok
}

func paramsToValueList(paramsVal data.Value) []data.Value {
	if paramsVal == nil {
		return nil
	}
	paramsArr, ok := paramsVal.(*data.ArrayValue)
	if !ok {
		return nil
	}
	return paramsArr.ToValueList()
}
