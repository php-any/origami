package array

import (
	"fmt"

	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayDiffUkeyFunction 实现 array_diff_ukey 函数
type ArrayDiffUkeyFunction struct{}

func NewArrayDiffUkeyFunction() data.FuncStmt {
	return &ArrayDiffUkeyFunction{}
}

func (f *ArrayDiffUkeyFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	raw, _ := ctx.GetIndexValue(0)
	params := paramsToValueList(raw)
	if len(params) < 3 {
		return nil, data.NewErrorThrowByName(nil, fmt.Errorf("array_diff_ukey expects at least 3 arguments"), "ArgumentCountError")
	}
	callback, ctl := node.ResolveCallback(ctx, params[len(params)-1])
	if ctl != nil {
		return nil, ctl
	}
	arrays := params[:len(params)-1]
	for _, value := range arrays {
		if _, ok := value.(*data.ArrayValue); !ok {
			return nil, throwMustBeArray("array_diff_ukey", value)
		}
	}
	var keys []data.Value
	for _, other := range arrays[1:] {
		for _, entry := range toKVEntries(other) {
			keys = append(keys, entry.key)
		}
	}
	var kept []kvEntry
	for _, entry := range toKVEntries(arrays[0]) {
		found := false
		for _, key := range keys {
			result, ctl := invokeCallback(ctx, callback, []data.Value{entry.key, key})
			if ctl != nil {
				return nil, ctl
			}
			numeric, ok := result.(data.AsInt)
			if ok {
				n, err := numeric.AsInt()
				if err == nil && n == 0 {
					found = true
					break
				}
			}
		}
		if !found {
			kept = append(kept, entry)
		}
	}
	return buildResultFromEntries(kept, true), nil
}

func (f *ArrayDiffUkeyFunction) GetName() string {
	return "array_diff_ukey"
}

var arrayDiffUkeyFunctionGetParams = []data.GetValue{
	node.NewParameters(nil, "arrays", 0, nil, data.Mixed{}),
}

func (f *ArrayDiffUkeyFunction) GetParams() []data.GetValue {
	return arrayDiffUkeyFunctionGetParams
}

var arrayDiffUkeyFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "arrays", 0, data.Mixed{}),
}

func (f *ArrayDiffUkeyFunction) GetVariables() []data.Variable {
	return arrayDiffUkeyFunctionGetVariables
}
