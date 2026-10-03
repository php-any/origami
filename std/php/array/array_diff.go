package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayDiffFunction() data.FuncStmt { return &ArrayDiffFunction{} }

type ArrayDiffFunction struct{}

func (fn *ArrayDiffFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	params, _ := ctx.GetIndexValue(0)
	arrays := paramsToValueList(params)
	if len(arrays) == 0 {
		return data.NewArrayValue(nil), nil
	}
	excluded := make(map[string]struct{})
	for _, array := range arrays[1:] {
		for _, entry := range toKVEntries(array) {
			excluded[entry.value.AsString()] = struct{}{}
		}
	}
	var kept []kvEntry
	for _, entry := range toKVEntries(arrays[0]) {
		if _, found := excluded[entry.value.AsString()]; !found {
			kept = append(kept, entry)
		}
	}
	return buildResultFromEntries(kept, true), nil
}

func (fn *ArrayDiffFunction) GetName() string { return "array_diff" }

var arrayDiffFunctionGetParams = []data.GetValue{
	node.NewParameters(nil, "arrays", 0, nil, nil),
}

func (fn *ArrayDiffFunction) GetParams() []data.GetValue { return arrayDiffFunctionGetParams }

var arrayDiffFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "arrays", 0, data.NewBaseType("array")),
}

func (fn *ArrayDiffFunction) GetVariables() []data.Variable { return arrayDiffFunctionGetVariables }
