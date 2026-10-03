package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayColumnFunction 实现 array_column
// array_column(array $array, int|string|null $column_key, int|string|null $index_key = null): array
type ArrayColumnFunction struct{}

func NewArrayColumnFunction() data.FuncStmt {
	return &ArrayColumnFunction{}
}

func (f *ArrayColumnFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	if _, ok := arrayVal.(*data.ArrayValue); !ok {
		return nil, throwMustBeArray("array_column", arrayVal)
	}
	columnKey, _ := ctx.GetIndexValue(1)
	indexKey, hasIndexParam := ctx.GetIndexValue(2)
	useIndexKey := hasIndexParam && !isNullValue(indexKey)

	entries := toKVEntries(arrayVal)
	result := data.NewArrayValue(nil).(*data.ArrayValue)

	for _, row := range entries {
		var col data.Value
		if isNullValue(columnKey) {
			col = row.value
		} else {
			var exists bool
			col, exists = lookupElementValue(row.value, columnKey)
			if !exists {
				continue
			}
		}

		if useIndexKey {
			if idx, exists := lookupElementValue(row.value, indexKey); exists {
				result.SetKey(idx, col)
				continue
			}
		}
		result.AppendValue(col)
	}
	return result, nil
}

func (f *ArrayColumnFunction) GetName() string { return "array_column" }

var arrayColumnFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "array", 0, nil, nil),
	node.NewParameter(nil, "column_key", 1, nil, nil),
	node.NewParameter(nil, "index_key", 2, nil, nil),
}

func (f *ArrayColumnFunction) GetParams() []data.GetValue {
	return arrayColumnFunctionGetParams
}

var arrayColumnFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
	node.NewVariable(nil, "column_key", 1, data.Mixed{}),
	node.NewVariable(nil, "index_key", 2, data.Mixed{}),
}

func (f *ArrayColumnFunction) GetVariables() []data.Variable {
	return arrayColumnFunctionGetVariables
}
