package support

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/php-any/origami/data"
)

func arrAccessible(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(isAccessible(v)), nil
}

func isAccessible(v data.Value) bool {
	if v == nil {
		return false
	}
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func arrExists(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	_, ok := arrayGet(arr, keyToString(key))
	return data.NewBoolValue(ok), nil
}

func arrHas(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keys, _ := ctx.GetIndexValue(1)
	for _, k := range keysToStrings(keys) {
		if _, ok := dataGetPath(arr, k); !ok {
			return data.NewBoolValue(false), nil
		}
	}
	return data.NewBoolValue(true), nil
}

func arrGet(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	def, _ := ctx.GetIndexValue(2)
	if key == nil || isNull(key) {
		return arr, nil
	}
	if v, ok := dataGetPath(arr, keyToString(key)); ok {
		return v, nil
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func arrSet(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	val, _ := ctx.GetIndexValue(2)
	root := asMutableArray(arr)
	if key == nil || isNull(key) {
		_ = ctx.SetVariableValue(nodeVar("array", 0), val)
		return val, nil
	}
	dataSetPath(root, keyToString(key), val)
	_ = ctx.SetVariableValue(nodeVar("array", 0), root)
	return root, nil
}

func arrAdd(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	if _, ok := dataGetPath(arr, keyToString(key)); ok {
		return arr, nil
	}
	return arrSet(ctx)
}

func arrForget(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keys, _ := ctx.GetIndexValue(1)
	root := asMutableArray(arr)
	for _, k := range keysToStrings(keys) {
		dataForgetPath(root, k)
	}
	_ = ctx.SetVariableValue(nodeVar("array", 0), root)
	return root, nil
}

func arrOnly(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keys, _ := ctx.GetIndexValue(1)
	want := map[string]bool{}
	for _, k := range keysToStrings(keys) {
		want[k] = true
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		if want[e.keyStr] {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return out, nil
}

func arrExcept(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keys, _ := ctx.GetIndexValue(1)
	deny := map[string]bool{}
	for _, k := range keysToStrings(keys) {
		deny[k] = true
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		if !deny[e.keyStr] {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return out, nil
}

func arrFirst(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	def, _ := ctx.GetIndexValue(2)
	for _, e := range toEntries(arr) {
		if cb == nil || isNull(cb) {
			return e.value, nil
		}
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ok {
			return e.value, nil
		}
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func arrLast(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	def, _ := ctx.GetIndexValue(2)
	entries := toEntries(arr)
	for i := len(entries) - 1; i >= 0; i-- {
		e := entries[i]
		if cb == nil || isNull(cb) {
			return e.value, nil
		}
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ok {
			return e.value, nil
		}
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func arrWrap(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	if v == nil || isNull(v) {
		return data.NewArrayValue(nil), nil
	}
	// PHP is_array：列表是 ArrayValue，关联数组是 ObjectValue。不可把后者再包成 [obj]。
	if av, ok := v.(*data.ArrayValue); ok {
		return av, nil
	}
	if ov, ok := v.(*data.ObjectValue); ok && ov != nil {
		return ov, nil
	}
	return data.NewArrayValue([]data.Value{v}), nil
}

func arrFrom(ctx data.Context) (data.GetValue, data.Control) {
	v, _ := ctx.GetIndexValue(0)
	return arrFromValue(ctx, v, 0)
}

func throwArrFromScalar() data.Control {
	return data.NewErrorThrowByName(nil, fmt.Errorf("Items cannot be represented by a scalar value."), "InvalidArgumentException")
}

func classIs(cv *data.ClassValue, name string) bool {
	if cv == nil {
		return false
	}
	return (data.Class{Name: name}).Is(cv)
}

func callVMFunc(ctx data.Context, name string, args ...data.Value) (data.Value, data.Control) {
	if ctx == nil || ctx.GetVM() == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined function %s()", name), "Error")
	}
	fn, ok := ctx.GetVM().GetFunc(name)
	if !ok || fn == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined function %s()", name), "Error")
	}
	callCtx := ctx.CreateContext(fn.GetVariables())
	data.BindDeclaredArgs(callCtx, fn, args)
	ret, ctl := fn.Call(callCtx)
	if ctl != nil {
		return nil, ctl
	}
	if v, ok := ret.(data.Value); ok {
		return unwrapValue(v), nil
	}
	return nil, nil
}

// arrFromValue 对齐 Illuminate\Support\Arr::from 的 match，不是 Arr::wrap。
// 标量抛 InvalidArgumentException，禁止包成 [v]。
func arrFromValue(ctx data.Context, v data.Value, depth int) (data.Value, data.Control) {
	if depth > 32 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Illuminate\\Support\\Arr::from(): maximum nesting level of 32 reached, aborting!"), "Error")
	}
	v = unwrapValue(v)
	if v == nil || isNull(v) {
		return data.NewArrayValue(nil), nil
	}
	if _, ok := v.(*data.ArrayValue); ok {
		return v, nil
	}
	if ov, ok := v.(*data.ObjectValue); ok && ov != nil {
		return ov, nil
	}
	if tv, ok := v.(*data.ThisValue); ok && tv != nil && tv.ClassValue != nil {
		v = tv.ClassValue
	}
	cv, isObj := v.(*data.ClassValue)
	if !isObj || cv == nil {
		return nil, throwArrFromScalar()
	}

	switch {
	case classIs(cv, "Illuminate\\Support\\Enumerable"):
		arr, ok, ctl := callClassNoArg(cv, "all")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::all()", cv.Class.GetName()), "Error")
		}
		if arr == cv {
			return nil, throwArrFromScalar()
		}
		return arrFromValue(ctx, arr, depth+1)
	case classIs(cv, "Illuminate\\Contracts\\Support\\Arrayable"):
		arr, ok, ctl := callClassNoArg(cv, "toArray")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::toArray()", cv.Class.GetName()), "Error")
		}
		if arr == cv {
			return nil, throwArrFromScalar()
		}
		return arrFromValue(ctx, arr, depth+1)
	case classIs(cv, "WeakMap"):
		ret, ctl := callVMFunc(ctx, "iterator_to_array", cv, data.NewBoolValue(false))
		if ctl != nil {
			return nil, ctl
		}
		return arrFromValue(ctx, ret, depth+1)
	case classIs(cv, "Traversable") || classIs(cv, "Iterator") || classIs(cv, "IteratorAggregate"):
		ret, ctl := callVMFunc(ctx, "iterator_to_array", cv)
		if ctl != nil {
			return nil, ctl
		}
		return arrFromValue(ctx, ret, depth+1)
	case classIs(cv, "Illuminate\\Contracts\\Support\\Jsonable"):
		js, ok, ctl := callClassNoArg(cv, "toJson")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::toJson()", cv.Class.GetName()), "Error")
		}
		decoded, ctl := callVMFunc(ctx, "json_decode", js, data.NewBoolValue(true))
		if ctl != nil {
			return nil, ctl
		}
		return arrFromValue(ctx, decoded, depth+1)
	case classIs(cv, "JsonSerializable"):
		arr, ok, ctl := callClassNoArg(cv, "jsonSerialize")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::jsonSerialize()", cv.Class.GetName()), "Error")
		}
		if arr == cv {
			return nil, throwArrFromScalar()
		}
		return arrFromValue(ctx, arr, depth+1)
	default:
		out := data.NewArrayValue(nil).(*data.ArrayValue)
		if cv.ObjectValue != nil {
			cv.ObjectValue.RangeProperties(func(key string, val data.Value) bool {
				setEntry(out, key, val)
				return true
			})
		}
		return out, nil
	}
}

func unwrapValue(v data.Value) data.Value {
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

func callClassNoArg(cv *data.ClassValue, name string) (data.Value, bool, data.Control) {
	if cv == nil {
		return nil, false, nil
	}
	m, ok := cv.GetMethod(name)
	if !ok || m == nil {
		return nil, false, nil
	}
	ret, ctl := m.Call(cv.CreateContext(m.GetVariables()))
	if ctl != nil {
		if rv, ok := ctl.(data.ReturnControl); ok {
			return unwrapValue(rv.ReturnValue()), true, nil
		}
		return nil, true, ctl
	}
	if ret == nil {
		return nil, true, nil
	}
	if v, ok := ret.(data.Value); ok {
		return unwrapValue(v), true, nil
	}
	return nil, true, nil
}

func foreachableArray(ctx data.Context, v data.Value) (data.Value, data.Control) {
	v = unwrapValue(v)
	if v == nil || isNull(v) {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("foreach() argument must be of type array|object, null given"), "TypeError")
	}
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return v, nil
	}
	cv, ok := v.(*data.ClassValue)
	if !ok || cv == nil {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("foreach() argument must be of type array|object, %s given", foreachValueType(v)), "TypeError")
	}
	if classIs(cv, "Illuminate\\Support\\Enumerable") {
		arr, ok, ctl := callClassNoArg(cv, "all")
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			return nil, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::all()", cv.Class.GetName()), "Error")
		}
		return foreachableArray(ctx, arr)
	}
	if classIs(cv, "Traversable") || classIs(cv, "Iterator") || classIs(cv, "IteratorAggregate") {
		return callVMFunc(ctx, "iterator_to_array", cv)
	}
	if _, has := cv.GetMethod("valid"); has {
		return callVMFunc(ctx, "iterator_to_array", cv)
	}
	if _, has := cv.GetMethod("getIterator"); has {
		return callVMFunc(ctx, "iterator_to_array", cv)
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if cv.ObjectValue != nil {
		cv.ObjectValue.RangeProperties(func(key string, val data.Value) bool {
			setEntry(out, key, val)
			return true
		})
	}
	return out, nil
}

func foreachValueType(v data.Value) string {
	switch v.(type) {
	case *data.StringValue:
		return "string"
	case *data.IntValue:
		return "int"
	case *data.FloatValue:
		return "float"
	case *data.BoolValue:
		return "bool"
	default:
		return "unknown"
	}
}

func arrPartition(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	unwrapped, ctl := foreachableArray(ctx, arr)
	if ctl != nil {
		return nil, ctl
	}
	arr = unwrapped
	passed := data.NewArrayValue(nil).(*data.ArrayValue)
	failed := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		ok := false
		if cb != nil && !isNull(cb) {
			var ctl data.Control
			ok, ctl = callBool(ctx, cb, e.value, e.key)
			if ctl != nil {
				return nil, ctl
			}
		}
		if ok {
			setEntry(passed, e.keyStr, e.value)
		} else {
			setEntry(failed, e.keyStr, e.value)
		}
	}
	return data.NewArrayValue([]data.Value{passed, failed}), nil
}

func arrCollapse(ctx data.Context) (data.GetValue, data.Control) {
	// 对齐 Illuminate\Support\Arr::collapse：
	// Collection → all() 再合并；is_array 合并一层；其它值跳过。
	// 旧实现把 ClassValue/ObjectValue 整段塞进结果，导致
	// registerConfiguredProviders 的 collapse()->toArray() 变成嵌套数组，
	// ProviderRepository 里 new $provider 拿到 ArrayValue。
	arr, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		nested, ok, ctl := collapseLayer(ctx, e.value)
		if ctl != nil {
			return nil, ctl
		}
		if !ok {
			continue
		}
		for _, ne := range nested {
			out.List = append(out.List, data.NewZVal(ne.value))
		}
	}
	return out, nil
}

func collapseLayer(ctx data.Context, v data.Value) ([]kv, bool, data.Control) {
	if v == nil || isNull(v) {
		return nil, false, nil
	}
	if cv, ok := v.(*data.ClassValue); ok && cv != nil && cv.Class != nil {
		if (data.Class{Name: "Illuminate\\Support\\Collection"}).Is(cv) {
			arr, ok, ctl := callClassNoArg(cv, "all")
			if ctl != nil {
				return nil, false, ctl
			}
			if !ok {
				return nil, false, data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::all()", cv.Class.GetName()), "Error")
			}
			return toEntries(arr), true, nil
		}
		return nil, false, nil
	}
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return toEntries(v), true, nil
	default:
		return nil, false, nil
	}
}

func arrFlatten(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	// PHP Arr::flatten($array, $depth = INF)
	depth := -1
	if d, ok := ctx.GetIndexValue(1); ok && d != nil && !isNull(d) {
		if fv, ok := d.(*data.FloatValue); ok {
			if fv.Value > 1e300 {
				depth = -1
			} else if n, err := fv.AsInt(); err == nil {
				depth = n
			}
		} else if iv, ok := d.(data.AsInt); ok {
			if n, err := iv.AsInt(); err == nil {
				depth = n
			}
		}
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if ctl := flattenInto(out, arr, depth); ctl != nil {
		return nil, ctl
	}
	return out, nil
}

func isFlattenableArray(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func flattenInto(out *data.ArrayValue, v data.Value, depth int) data.Control {
	if depth == 0 {
		out.List = append(out.List, data.NewZVal(v))
		return nil
	}
	if cv, ok := v.(*data.ClassValue); ok && cv != nil && cv.Class != nil {
		if (data.Class{Name: "Illuminate\\Support\\Collection"}).Is(cv) {
			arr, ok, ctl := callClassNoArg(cv, "all")
			if ctl != nil {
				return ctl
			}
			if !ok {
				return data.NewErrorThrowByName(nil, fmt.Errorf("Call to undefined method %s::all()", cv.Class.GetName()), "Error")
			}
			return flattenInto(out, arr, depth)
		}
	}
	if !isFlattenableArray(v) {
		out.List = append(out.List, data.NewZVal(v))
		return nil
	}
	for _, e := range toEntries(v) {
		if isFlattenableArray(e.value) || isCollectionValue(e.value) {
			next := depth - 1
			if depth < 0 {
				next = -1
			}
			if depth == 1 && isFlattenableArray(e.value) {
				for _, ne := range toEntries(e.value) {
					out.List = append(out.List, data.NewZVal(ne.value))
				}
			} else if ctl := flattenInto(out, e.value, next); ctl != nil {
				return ctl
			}
		} else {
			out.List = append(out.List, data.NewZVal(e.value))
		}
	}
	return nil
}

func isCollectionValue(v data.Value) bool {
	cv, ok := v.(*data.ClassValue)
	return ok && cv != nil && (data.Class{Name: "Illuminate\\Support\\Collection"}).Is(cv)
}

func arrDot(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	prepend := ""
	if p, ok := ctx.GetIndexValue(1); ok && p != nil && !isNull(p) {
		prepend = p.AsString()
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	dotInto(out, arr, prepend)
	return out, nil
}

func dotInto(out *data.ArrayValue, v data.Value, prefix string) {
	for _, e := range toEntries(v) {
		key := e.keyStr
		if prefix != "" {
			key = prefix + "." + e.keyStr
		}
		if isAccessible(e.value) && len(toEntries(e.value)) > 0 {
			dotInto(out, e.value, key)
		} else {
			setEntry(out, key, e.value)
		}
	}
}

func arrUndot(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		dataSetPath(out, e.keyStr, e.value)
	}
	return out, nil
}

func arrPluck(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	valueKey, _ := ctx.GetIndexValue(1)
	keyKey, _ := ctx.GetIndexValue(2)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		var extracted data.Value = e.value
		if valueKey != nil && !isNull(valueKey) {
			if v, ok := dataGetPath(e.value, keyToString(valueKey)); ok {
				extracted = v
			} else {
				extracted = data.NewNullValue()
			}
		}
		if keyKey != nil && !isNull(keyKey) {
			if k, ok := dataGetPath(e.value, keyToString(keyKey)); ok {
				setEntry(out, keyToString(k), extracted)
				continue
			}
		}
		out.List = append(out.List, data.NewZVal(extracted))
	}
	return out, nil
}

func arrMap(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		ret, ctl := callValue(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		var vv data.Value = data.NewNullValue()
		if ret != nil {
			if v, ok := ret.(data.Value); ok {
				vv = v
			}
		}
		setEntry(out, e.keyStr, vv)
	}
	return out, nil
}

func arrMapWithKeys(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		ret, ctl := callValue(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		var assoc data.Value
		if ret != nil {
			if v, ok := ret.(data.Value); ok {
				assoc = v
			}
		}
		for _, ne := range toEntries(assoc) {
			setEntry(out, ne.keyStr, ne.value)
		}
	}
	return out, nil
}

func arrIsList(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	return data.NewBoolValue(isListArray(arr)), nil
}

func arrIsAssoc(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	entries := toEntries(arr)
	return data.NewBoolValue(len(entries) > 0 && !isListArray(arr)), nil
}

func arrWhere(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	cb, _ := ctx.GetIndexValue(1)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		ok, ctl := callBool(ctx, cb, e.value, e.key)
		if ctl != nil {
			return nil, ctl
		}
		if ok {
			setEntry(out, e.keyStr, e.value)
		}
	}
	return out, nil
}

func arrPrepend(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	val, _ := ctx.GetIndexValue(1)
	key, _ := ctx.GetIndexValue(2)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	if key != nil && !isNull(key) {
		setEntry(out, keyToString(key), val)
		for _, e := range toEntries(arr) {
			if e.keyStr != keyToString(key) {
				setEntry(out, e.keyStr, e.value)
			}
		}
		return out, nil
	}
	out.List = append(out.List, data.NewZVal(val))
	for _, e := range toEntries(arr) {
		out.List = append(out.List, data.NewZVal(e.value))
	}
	return out, nil
}

func arrPull(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	key, _ := ctx.GetIndexValue(1)
	def, _ := ctx.GetIndexValue(2)
	root := asMutableArray(arr)
	ks := keyToString(key)
	v, ok := dataGetPath(root, ks)
	dataForgetPath(root, ks)
	_ = ctx.SetVariableValue(nodeVar("array", 0), root)
	if ok {
		return v, nil
	}
	if def == nil {
		return data.NewNullValue(), nil
	}
	return def, nil
}

func arrQuery(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	vals := url.Values{}
	for _, e := range toEntries(arr) {
		vals.Set(e.keyStr, e.value.AsString())
	}
	return data.NewStringValue(vals.Encode()), nil
}

func arrRandom(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	entries := toEntries(arr)
	if len(entries) == 0 {
		return data.NewNullValue(), nil
	}
	return entries[0].value, nil
}

func arrShuffle(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		out.List = append(out.List, data.NewZVal(e.value))
	}
	return out, nil
}

func arrSort(ctx data.Context) (data.GetValue, data.Control) {
	entries := toEntries(mustIndex(ctx, 0))
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].value.AsString() < entries[j].value.AsString()
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return out, nil
}

func arrSortDesc(ctx data.Context) (data.GetValue, data.Control) {
	entries := toEntries(mustIndex(ctx, 0))
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].value.AsString() > entries[j].value.AsString()
	})
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		setEntry(out, e.keyStr, e.value)
	}
	return out, nil
}

func arrSortRecursive(ctx data.Context) (data.GetValue, data.Control) {
	return arrSort(ctx)
}

func arrToCssClasses(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	parts := make([]string, 0)
	for _, e := range toEntries(arr) {
		// PHP：is_numeric($class) 时值是类名；否则键是类名、值是约束。
		if cssAssocKeyIsNumeric(e.keyStr) {
			if s := strings.TrimSpace(e.value.AsString()); s != "" {
				parts = append(parts, s)
			}
			continue
		}
		if truthy(e.value) {
			if s := strings.TrimSpace(e.keyStr); s != "" {
				parts = append(parts, s)
			}
		}
	}
	return data.NewStringValue(strings.Join(parts, " ")), nil
}

func arrToCssStyles(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	parts := make([]string, 0)
	for _, e := range toEntries(arr) {
		var s string
		if cssAssocKeyIsNumeric(e.keyStr) {
			s = strings.TrimSpace(e.value.AsString())
		} else if truthy(e.value) {
			s = strings.TrimSpace(e.keyStr)
		}
		if s == "" {
			continue
		}
		if !strings.HasSuffix(s, ";") {
			s += ";"
		}
		parts = append(parts, s)
	}
	return data.NewStringValue(strings.Join(parts, " ")), nil
}

func cssAssocKeyIsNumeric(key string) bool {
	if key == "" {
		return true
	}
	if _, ok := data.ParseIntArrayKeyName(key); ok {
		return true
	}
	_, err := strconv.Atoi(key)
	return err == nil
}

func arrJoin(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	glue := ""
	if g, ok := ctx.GetIndexValue(1); ok && g != nil {
		glue = g.AsString()
	}
	finalGlue := glue
	if fg, ok := ctx.GetIndexValue(2); ok && fg != nil && !isNull(fg) {
		finalGlue = fg.AsString()
	}
	entries := toEntries(arr)
	strs := make([]string, len(entries))
	for i, e := range entries {
		strs[i] = e.value.AsString()
	}
	if len(strs) == 0 {
		return data.NewStringValue(""), nil
	}
	if len(strs) == 1 {
		return data.NewStringValue(strs[0]), nil
	}
	if finalGlue != glue {
		return data.NewStringValue(strings.Join(strs[:len(strs)-1], glue) + finalGlue + strs[len(strs)-1]), nil
	}
	return data.NewStringValue(strings.Join(strs, glue)), nil
}

func arrKeyBy(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keyBy, _ := ctx.GetIndexValue(1)
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		k := e.keyStr
		if keyBy != nil && !isNull(keyBy) {
			if v, ok := dataGetPath(e.value, keyToString(keyBy)); ok {
				k = keyToString(v)
			}
		}
		setEntry(out, k, e.value)
	}
	return out, nil
}

func arrDivide(ctx data.Context) (data.GetValue, data.Control) {
	arr, _ := ctx.GetIndexValue(0)
	keys := data.NewArrayValue(nil).(*data.ArrayValue)
	vals := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(arr) {
		keys.List = append(keys.List, data.NewZVal(data.NewStringValue(e.keyStr)))
		vals.List = append(vals.List, data.NewZVal(e.value))
	}
	return data.NewArrayValue([]data.Value{keys, vals}), nil
}

// ---- helpers ----

type kv struct {
	key    data.Value
	keyStr string
	value  data.Value
}

func toEntries(v data.Value) []kv {
	if v == nil {
		return nil
	}
	v = unwrapValue(v)
	switch arr := v.(type) {
	case *data.ArrayValue:
		entries := make([]kv, 0, len(arr.List))
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
			entries = append(entries, kv{key: key, keyStr: keyStr, value: z.Value})
		}
		return entries
	case *data.ObjectValue:
		entries := make([]kv, 0)
		arr.RangeProperties(func(key string, value data.Value) bool {
			entries = append(entries, kv{key: data.NewStringValue(key), keyStr: key, value: value})
			return true
		})
		return entries
	default:
		return nil
	}
}

func isNull(v data.Value) bool {
	_, ok := v.(*data.NullValue)
	return ok
}

func truthy(v data.Value) bool {
	if v == nil || isNull(v) {
		return false
	}
	if b, ok := v.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv
	}
	s := v.AsString()
	return s != "" && s != "0"
}

func keyToString(v data.Value) string {
	if v == nil {
		return ""
	}
	if iv, ok := v.(*data.IntValue); ok {
		return strconv.Itoa(iv.Value)
	}
	return v.AsString()
}

func keysToStrings(v data.Value) []string {
	if v == nil {
		return nil
	}
	if _, ok := v.(*data.ArrayValue); ok {
		out := make([]string, 0)
		for _, e := range toEntries(v) {
			out = append(out, keyToString(e.value))
		}
		return out
	}
	return []string{keyToString(v)}
}

func arrayGet(arr data.Value, key string) (data.Value, bool) {
	return dataGetPath(arr, key)
}

func dataGetPath(target data.Value, path string) (data.Value, bool) {
	if path == "" {
		return target, target != nil
	}
	cur := target
	for _, seg := range strings.Split(path, ".") {
		if cur == nil {
			return nil, false
		}
		found := false
		switch t := cur.(type) {
		case *data.ArrayValue:
			if z, ok := t.LookupZValByStringKey(seg); ok && z != nil {
				cur, found = z.Value, true
			}
		case *data.ObjectValue:
			// GetProperty 对缺失键返回 NullValue，不能当成“键存在”，
			// 否则 Arr::get($aliases, 'prefix', 'prefix') 会丢掉 default。
			if t.HasProperty(seg) {
				v, ctl := t.GetProperty(seg)
				if ctl == nil {
					cur, found = v, true
				}
			}
		case *data.ClassValue:
			v, ctl := t.GetProperty(seg)
			if ctl == nil && v != nil && !isNull(v) {
				cur, found = v, true
			}
		}
		if !found {
			return nil, false
		}
	}
	return cur, true
}

func dataSetPath(target *data.ArrayValue, path string, value data.Value) {
	parts := strings.Split(path, ".")
	cur := target
	for i, seg := range parts {
		if i == len(parts)-1 {
			setEntry(cur, seg, value)
			return
		}
		nextVal, ok := dataGetPath(cur, seg)
		var next *data.ArrayValue
		if ok {
			if av, ok := nextVal.(*data.ArrayValue); ok {
				next = av
			}
		}
		if next == nil {
			next = data.NewArrayValue(nil).(*data.ArrayValue)
			setEntry(cur, seg, next)
		}
		cur = next
	}
}

func dataForgetPath(target *data.ArrayValue, path string) {
	parts := strings.Split(path, ".")
	if len(parts) == 1 {
		target.UnsetKey(data.NewStringValue(parts[0]))
		return
	}
	parentPath := strings.Join(parts[:len(parts)-1], ".")
	parent, ok := dataGetPath(target, parentPath)
	if !ok {
		return
	}
	if av, ok := parent.(*data.ArrayValue); ok {
		av.UnsetKey(data.NewStringValue(parts[len(parts)-1]))
	}
}

func asMutableArray(v data.Value) *data.ArrayValue {
	if av, ok := v.(*data.ArrayValue); ok {
		return av
	}
	out := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range toEntries(v) {
		setEntry(out, e.keyStr, e.value)
	}
	return out
}

func setEntry(arr *data.ArrayValue, key string, val data.Value) {
	if arr == nil {
		return
	}
	if n, ok := data.ParseIntArrayKeyName(key); ok {
		arr.SetIntKey(n, val)
		return
	}
	if z, ok := arr.LookupZValByStringKey(key); ok && z != nil {
		z.Value = val
		return
	}
	arr.List = append(arr.List, data.NewNamedZVal(key, val))
}

func isListArray(v data.Value) bool {
	av, ok := v.(*data.ArrayValue)
	if !ok || av == nil {
		return false
	}
	for i, z := range av.List {
		if z == nil {
			continue
		}
		if z.Name == "" {
			continue
		}
		if n, ok := data.ParseIntArrayKeyName(z.Name); !ok || n != i {
			return false
		}
	}
	return true
}

func mustIndex(ctx data.Context, i int) data.Value {
	v, _ := ctx.GetIndexValue(i)
	return v
}

func nodeVar(name string, index int) data.Variable {
	return &simpleVar{name: name, index: index}
}

type simpleVar struct {
	name  string
	index int
}

func (v *simpleVar) GetName() string     { return v.name }
func (v *simpleVar) GetIndex() int       { return v.index }
func (v *simpleVar) GetType() data.Types { return nil }
func (v *simpleVar) SetType(data.Types)  {}
func (v *simpleVar) GetValue(ctx data.Context) (data.GetValue, data.Control) {
	return ctx.GetVariableValue(v)
}
func (v *simpleVar) SetValue(ctx data.Context, value data.Value) data.Control {
	return ctx.SetVariableValue(v, value)
}

func callBool(ctx data.Context, cb data.Value, args ...data.Value) (bool, data.Control) {
	ret, ctl := callValue(ctx, cb, args...)
	if ctl != nil {
		return false, ctl
	}
	if ret == nil {
		return false, nil
	}
	if b, ok := ret.(data.AsBool); ok {
		okv, err := b.AsBool()
		return err == nil && okv, nil
	}
	if v, ok := ret.(data.Value); ok {
		return truthy(v), nil
	}
	return false, nil
}

func callValue(ctx data.Context, cb data.Value, args ...data.Value) (data.GetValue, data.Control) {
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
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Support\\Arr: callback is not callable"))
	}
	if fn == nil {
		return nil, data.NewErrorThrow(nil, fmt.Errorf("Illuminate\\Support\\Arr: callback is not callable"))
	}
	callCtx := ctx.CreateContext(fn.GetVariables())
	data.BindDeclaredArgs(callCtx, fn, args)
	if bfv, ok := cb.(*data.BoundFuncValue); ok {
		return bfv.Call(callCtx)
	}
	return fn.Call(callCtx)
}
