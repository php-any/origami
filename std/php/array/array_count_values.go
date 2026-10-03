package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayCountValuesFunction 实现 array_count_values
// array_count_values(array $array): array
type ArrayCountValuesFunction struct{}

func NewArrayCountValuesFunction() data.FuncStmt {
	return &ArrayCountValuesFunction{}
}

func (f *ArrayCountValuesFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	if _, ok := arrayVal.(*data.ArrayValue); !ok {
		return nil, throwMustBeArray("array_count_values", arrayVal)
	}
	entries := toKVEntries(arrayVal)
	if len(entries) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	result := data.NewArrayValue(nil).(*data.ArrayValue)
	for _, e := range entries {
		switch e.value.(type) {
		case *data.IntValue, *data.StringValue:
		default:
			continue
		}
		k := e.value.AsString()
		count := 0
		if slot, ok := result.LookupZValByStringKey(k); ok {
			count = slot.ReadValue().(*data.IntValue).Value
		}
		result.SetStringKey(k, data.NewIntValue(count+1))
	}
	return result, nil
}

func (f *ArrayCountValuesFunction) GetName() string { return "array_count_values" }

var arrayCountValuesFunctionGetParams = []data.GetValue{node.NewParameter(nil, "array", 0, nil, nil)}

func (f *ArrayCountValuesFunction) GetParams() []data.GetValue {
	return arrayCountValuesFunctionGetParams
}

var arrayCountValuesFunctionGetVariables = []data.Variable{node.NewVariable(nil, "array", 0, data.NewBaseType("array"))}

func (f *ArrayCountValuesFunction) GetVariables() []data.Variable {
	return arrayCountValuesFunctionGetVariables
}
