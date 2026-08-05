package array

import (
	"sort"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

type valueSortEntry struct {
	key   string
	value data.Value
}

func sortPreservingKeys(value data.Value, flags int, descending bool) bool {
	less := func(a, b data.Value) bool {
		if descending {
			return compareValues(b, a, flags)
		}
		return compareValues(a, b, flags)
	}

	switch collection := value.(type) {
	case *data.ArrayValue:
		sort.SliceStable(collection.List, func(i, j int) bool {
			return less(collection.List[i].Value, collection.List[j].Value)
		})
		return true
	case *data.ObjectValue:
		entries := make([]valueSortEntry, 0)
		collection.RangeProperties(func(key string, value data.Value) bool {
			entries = append(entries, valueSortEntry{key: key, value: value})
			return true
		})
		sort.SliceStable(entries, func(i, j int) bool {
			return less(entries[i].value, entries[j].value)
		})
		sorted := data.NewObjectValue()
		for _, entry := range entries {
			sorted.SetProperty(entry.key, entry.value)
		}
		*collection = *sorted
		return true
	default:
		return false
	}
}

type AsortFunction struct {
	descending bool
}

func NewAsortFunction() data.FuncStmt  { return &AsortFunction{} }
func NewArsortFunction() data.FuncStmt { return &AsortFunction{descending: true} }

func (f *AsortFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	value, ok := ctx.GetIndexValue(0)
	if !ok {
		return data.NewBoolValue(false), nil
	}
	flags := 0
	if flagValue, ok := ctx.GetIndexValue(1); ok {
		if asInt, ok := flagValue.(data.AsInt); ok {
			flags, _ = asInt.AsInt()
		}
	}
	return data.NewBoolValue(sortPreservingKeys(value, flags, f.descending)), nil
}

func (f *AsortFunction) GetName() string {
	if f.descending {
		return "arsort"
	}
	return "asort"
}
func (f *AsortFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
		node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
	}
}
func (f *AsortFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "flags", 1, data.Int{}),
	}
}
