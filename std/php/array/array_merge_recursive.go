package array

import (
	"strconv"

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
	key   string
	name  string
	value data.Value
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
			if key == "" {
				key = data.IntArrayKeyName(i)
			}
			entries = append(entries, arrayEntry{key: key, name: z.Name, value: z.Value})
		}
		return entries, true
	case *data.ObjectValue:
		entries := make([]arrayEntry, 0, len(val.GetProperties()))
		for key, prop := range val.GetProperties() {
			entries = append(entries, arrayEntry{key: key, name: key, value: prop})
		}
		return entries, true
	default:
		return nil, false
	}
}

func entriesToArrayValue(entries []arrayEntry) *data.ArrayValue {
	list := make([]*data.ZVal, len(entries))
	for i, entry := range entries {
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
		merged := make(map[string]arrayEntry, len(baseEntries)+len(otherEntries))
		order := make([]string, 0, len(baseEntries)+len(otherEntries))

		addEntry := func(entry arrayEntry) {
			if existing, ok := merged[entry.key]; ok {
				merged[entry.key] = arrayEntry{
					key:   entry.key,
					name:  entry.name,
					value: mergeRecursive(existing.value, entry.value),
				}
				return
			}
			merged[entry.key] = arrayEntry{
				key:   entry.key,
				name:  entry.name,
				value: deepCopyVal(entry.value),
			}
			order = append(order, entry.key)
		}

		for _, entry := range baseEntries {
			addEntry(entry)
		}
		for _, entry := range otherEntries {
			addEntry(entry)
		}

		result := make([]arrayEntry, 0, len(order))
		for _, key := range order {
			result = append(result, merged[key])
		}
		return entriesToArrayValue(result)
	}

	if isListArray(base) && !isArrayLike(other) {
		entries, ok := arrayEntries(base)
		if ok {
			entries = append(entries, arrayEntry{
				key:   strconv.Itoa(len(entries)),
				name:  "",
				value: other,
			})
			return entriesToArrayValue(entries)
		}
	}

	return other
}

func isArrayLike(v data.Value) bool {
	switch v.(type) {
	case *data.ArrayValue, *data.ObjectValue:
		return true
	default:
		return false
	}
}

func isListArray(v data.Value) bool {
	arr, ok := v.(*data.ArrayValue)
	if !ok {
		return false
	}
	for _, z := range arr.List {
		if z != nil && z.Name != "" {
			return false
		}
	}
	return true
}

func deepCopyVal(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		out := data.NewObjectValue()
		for k, prop := range val.GetProperties() {
			out.SetProperty(k, deepCopyVal(prop))
		}
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
