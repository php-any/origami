package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

func NewArrayReplaceFunction() data.FuncStmt {
	return &ArrayReplaceFunction{}
}

type ArrayReplaceFunction struct{}

func (fn *ArrayReplaceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
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

	result := shallowCopy(arrays[0])
	for _, arr := range arrays[1:] {
		result = replaceRecursive(result, arr)
	}
	return result, nil
}

func replaceRecursive(base, other data.Value) data.Value {
	baseObj, bOk := base.(*data.ObjectValue)
	otherObj, oOk := other.(*data.ObjectValue)
	if bOk && oOk {
		out := data.NewObjectValue()
		for k, v := range baseObj.GetProperties() {
			out.SetProperty(k, v)
		}
		for k, v := range otherObj.GetProperties() {
			out.SetProperty(k, v)
		}
		return out
	}

	baseArr, bOk := base.(*data.ArrayValue)
	otherArr, oOk := other.(*data.ArrayValue)
	if bOk && oOk {
		bVals := baseArr.ToValueList()
		oVals := otherArr.ToValueList()
		result := make([]data.Value, len(bVals))
		copy(result, bVals)
		for i, v := range oVals {
			if i < len(result) {
				result[i] = v
			} else {
				result = append(result, v)
			}
		}
		return data.NewArrayValue(result)
	}

	return other
}

func shallowCopy(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		out := data.NewObjectValue()
		for k, prop := range val.GetProperties() {
			out.SetProperty(k, prop)
		}
		return out
	case *data.ArrayValue:
		vals := val.ToValueList()
		result := make([]data.Value, len(vals))
		copy(result, vals)
		return data.NewArrayValue(result)
	}
	return v
}

func (fn *ArrayReplaceFunction) GetName() string { return "array_replace" }
func (fn *ArrayReplaceFunction) GetParams() []data.GetValue {
	return []data.GetValue{node.NewParameters(nil, "arrays", 0, nil, nil)}
}
func (fn *ArrayReplaceFunction) GetVariables() []data.Variable {
	return []data.Variable{node.NewVariable(nil, "arrays", 0, data.NewBaseType("array"))}
}
