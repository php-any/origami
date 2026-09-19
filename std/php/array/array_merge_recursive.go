package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayMergeRecursiveFunction() data.FuncStmt {
	return &ArrayMergeRecursiveFunction{}
}

type ArrayMergeRecursiveFunction struct{}

func (fn *ArrayMergeRecursiveFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	paramsVal, _ := ctx.GetIndexValue(0)
	if paramsVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}
	paramsArr, ok := paramsVal.(*data.ArrayValue)
	if !ok {
		return data.NewArrayValue([]data.Value{}), nil
	}
	arrays := paramsArr.ToValueList()
	if len(arrays) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := deepCopyVal(arrays[0])
	for _, arr := range arrays[1:] {
		result = mergeRecursive(result, arr)
	}
	return result, nil
}

type arrayEntry struct {
	key      string
	name     string
	value    data.Value
	emptyStr bool
}

func arrayEntries(v data.Value) ([]arrayEntry, bool) {
	switch val := v.(type) {
	case *data.ArrayValue:
		entries := make([]arrayEntry, 0, len(val.List))
		for i, z := range val.List {
			if z == nil {
				continue
			}
			key := z.Name
			if z.EmptyStrKey {
				key = ""
			} else if key == "" {
				key = data.IntArrayKeyName(i)
			}
			entries = append(entries, arrayEntry{key: key, name: z.Name, value: z.Value, emptyStr: z.EmptyStrKey})
		}
		return entries, true
	case *data.ObjectValue:
		entries := make([]arrayEntry, 0)
		// 按插入顺序遍历，避免 Go map 顺序随机导致结果不稳定
		val.RangeProperties(func(key string, prop data.Value) bool {
			entries = append(entries, arrayEntry{key: key, name: key, value: prop})
			return true
		})
		return entries, true
	default:
		return nil, false
	}
}

func entriesToArrayValue(entries []arrayEntry) *data.ArrayValue {
	list := make([]*data.ZVal, len(entries))
	for i, entry := range entries {
		if entry.emptyStr {
			list[i] = data.NewEmptyStringKeyZVal(deepCopyVal(entry.value))
			continue
		}
		z := data.NewZVal(deepCopyVal(entry.value))
		z.Name = entry.name
		list[i] = z
	}
	return &data.ArrayValue{List: list}
}

func mergeRecursive(base, other data.Value) data.Value {
	baseEntries, bOk := arrayEntries(base)
	otherEntries, oOk := arrayEntries(other)
	if bOk && oOk {
		return mergeRecursiveArrays(baseEntries, otherEntries)
	}
	if bOk {
		return appendScalarToArray(baseEntries, other)
	}
	if oOk {
		return prependScalarToArray(base, otherEntries)
	}
	return entriesToArrayValue([]arrayEntry{
		{key: "0", name: "0", value: base},
		{key: "1", name: "1", value: other},
	})
}

// mergeRecursiveArrays 对齐 PHP array_merge_recursive：
// 字符串键递归合并；整数键已存在则追加新编号，绝不覆盖。
// Laravel RouteGroup 用它把外层 EncryptedCookies/StartSession 与内层 Authenticate 拼在一起。
func mergeRecursiveArrays(baseEntries, otherEntries []arrayEntry) data.Value {
	merged := make(map[string]arrayEntry, len(baseEntries)+len(otherEntries))
	order := make([]string, 0, len(baseEntries)+len(otherEntries))
	maxInt := -1

	putNew := func(entry arrayEntry) {
		merged[entry.key] = entry
		order = append(order, entry.key)
		if n, ok := data.ParseIntArrayKeyName(entry.key); ok && n > maxInt {
			maxInt = n
		}
	}

	for _, entry := range baseEntries {
		putNew(arrayEntry{
			key:      entry.key,
			name:     entryName(entry),
			value:    deepCopyVal(entry.value),
			emptyStr: entry.emptyStr,
		})
	}

	for _, entry := range otherEntries {
		val := deepCopyVal(entry.value)
		if !entry.emptyStr && !isStringArrayKey(entry.key) {
			key := entry.key
			if key == "" {
				maxInt++
				key = data.IntArrayKeyName(maxInt)
			} else if _, exists := merged[key]; exists {
				maxInt++
				key = data.IntArrayKeyName(maxInt)
			}
			putNew(arrayEntry{key: key, name: key, value: val})
			continue
		}
		if existing, ok := merged[entry.key]; ok {
			merged[entry.key] = arrayEntry{
				key:      entry.key,
				name:     entry.name,
				value:    mergeRecursive(existing.value, val),
				emptyStr: entry.emptyStr || existing.emptyStr,
			}
			continue
		}
		putNew(arrayEntry{key: entry.key, name: entry.name, value: val, emptyStr: entry.emptyStr})
	}

	result := make([]arrayEntry, 0, len(order))
	for _, key := range order {
		result = append(result, merged[key])
	}
	return entriesToArrayValue(result)
}

func entryName(entry arrayEntry) string {
	if entry.name != "" {
		return entry.name
	}
	return entry.key
}

func appendScalarToArray(entries []arrayEntry, scalar data.Value) data.Value {
	maxInt := -1
	for _, e := range entries {
		if n, ok := data.ParseIntArrayKeyName(e.key); ok && n > maxInt {
			maxInt = n
		}
	}
	maxInt++
	key := data.IntArrayKeyName(maxInt)
	entries = append(append([]arrayEntry(nil), entries...), arrayEntry{key: key, name: key, value: scalar})
	return entriesToArrayValue(entries)
}

func prependScalarToArray(scalar data.Value, entries []arrayEntry) data.Value {
	out := make([]arrayEntry, 0, len(entries)+1)
	out = append(out, arrayEntry{key: "0", name: "0", value: scalar})
	next := 1
	for _, e := range entries {
		if !isStringArrayKey(e.key) {
			key := data.IntArrayKeyName(next)
			out = append(out, arrayEntry{key: key, name: key, value: e.value})
			next++
			continue
		}
		out = append(out, e)
	}
	return entriesToArrayValue(out)
}

func deepCopyVal(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		out := data.NewObjectValue()
		// 按插入顺序遍历，避免 Go map 顺序随机
		val.RangeProperties(func(k string, prop data.Value) bool {
			out.SetProperty(k, deepCopyVal(prop))
			return true
		})
		return out
	case *data.ArrayValue:
		list := make([]*data.ZVal, len(val.List))
		for i, z := range val.List {
			if z == nil {
				continue
			}
			copied := data.NewZVal(deepCopyVal(z.Value))
			copied.Name = z.Name
			list[i] = copied
		}
		return &data.ArrayValue{List: list}
	}
	return v
}

func (fn *ArrayMergeRecursiveFunction) GetName() string { return "array_merge_recursive" }
func (fn *ArrayMergeRecursiveFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameters(nil, "arrays", 0, nil, nil)}
}
func (fn *ArrayMergeRecursiveFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "arrays", 0, data.NewBaseType("array"))}
}
