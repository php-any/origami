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
	columnKey, _ := ctx.GetIndexValue(1)
	indexKey, hasIndexParam := ctx.GetIndexValue(2)
	useIndexKey := hasIndexParam && !isNullValue(indexKey)

	if arrayVal == nil {
		return data.NewArrayValue([]data.Value{}), nil
	}

	entries := toKVEntries(arrayVal)
	result := data.NewObjectValue()
	dense := make([]data.Value, 0)

	for _, row := range entries {
		var col data.Value
		if isNullValue(columnKey) {
			col = row.value
		} else {
			col = getElementValue(row.value, columnKey)
		}

		if useIndexKey {
			idx := getElementValue(row.value, indexKey)
			result.SetProperty(idx.AsString(), col)
		} else {
			dense = append(dense, col)
		}
	}

	if useIndexKey {
		return result, nil
	}
	return data.NewArrayValue(dense), nil
}

func (f *ArrayColumnFunction) GetName() string { return "array_column" }

func (f *ArrayColumnFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "column_key", 1, nil, nil),
		node.NewParameter(nil, "index_key", 2, nil, nil),
	}
}

func (f *ArrayColumnFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
		node.NewVariable(nil, "column_key", 1, data.Mixed{}),
		node.NewVariable(nil, "index_key", 2, data.Mixed{}),
	}
}
