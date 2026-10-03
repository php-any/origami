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
		collection.EditPreservingKeys(func(slots []*data.ZVal) {
			sort.SliceStable(slots, func(i, j int) bool {
				return less(slots[i].ReadValue(), slots[j].ReadValue())
			})
		})
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
	value := data.CowSeparateIndex(ctx, 0)
	ok := value != nil
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

var asortFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.Mixed{}),
	node.NewParameter(nil, "flags", 1, data.NewIntValue(0), data.Int{}),
}

func (f *AsortFunction) GetParams() []data.GetValue {
	return asortFunctionGetParams
}

var asortFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.Mixed{}),
	node.NewVariable(nil, "flags", 1, data.Int{}),
}

func (f *AsortFunction) GetVariables() []data.Variable {
	return asortFunctionGetVariables
}
