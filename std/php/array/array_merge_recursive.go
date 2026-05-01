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

func mergeRecursive(base, other data.Value) data.Value {
	baseObj, bOk := base.(*data.ObjectValue)
	otherObj, oOk := other.(*data.ObjectValue)
	if bOk && oOk {
		out := data.NewObjectValue()
		for k, v := range baseObj.GetProperties() {
			out.SetProperty(k, v)
		}
		for k, v := range otherObj.GetProperties() {
			if bv, _ := baseObj.GetProperty(k); bv != nil {
				if _, isNull := bv.(*data.NullValue); !isNull {
					out.SetProperty(k, mergeRecursive(bv, v))
					continue
				}
			}
			out.SetProperty(k, v)
		}
		return out
	}

	baseArr, bOk := base.(*data.ArrayValue)
	otherArr, oOk := other.(*data.ArrayValue)
	if bOk && oOk {
		bVals := baseArr.ToValueList()
		oVals := otherArr.ToValueList()
		maxLen := len(bVals)
		if len(oVals) > maxLen {
			maxLen = len(oVals)
		}
		result := make([]data.Value, maxLen)
		for i := 0; i < maxLen; i++ {
			var bv data.Value = data.NewNullValue()
			var ov data.Value = data.NewNullValue()
			if i < len(bVals) {
				bv = bVals[i]
			}
			if i < len(oVals) {
				ov = oVals[i]
			}
			result[i] = mergeRecursive(bv, ov)
		}
		return data.NewArrayValue(result)
	}

	return other
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
		vals := val.ToValueList()
		result := make([]data.Value, len(vals))
		for i, item := range vals {
			result[i] = deepCopyVal(item)
		}
		return data.NewArrayValue(result)
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
