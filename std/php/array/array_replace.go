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

	result := shallowCopyPreserveKeys(arrays[0])
	for _, arr := range arrays[1:] {
		result = replaceByKeys(result, arr)
	}
	return result, nil
}

// replaceByKeys 按键替换（对齐 PHP array_replace），保留字符串键与稀疏整数键。
func replaceByKeys(base, other data.Value) data.Value {
	baseObj, bOk := base.(*data.ObjectValue)
	otherObj, oOk := other.(*data.ObjectValue)
	if bOk && oOk {
		out := data.NewObjectValue()
		baseObj.RangeProperties(func(k string, v data.Value) bool {
			out.SetProperty(k, v)
			return true
		})
		otherObj.RangeProperties(func(k string, v data.Value) bool {
			out.SetProperty(k, v)
			return true
		})
		return out
	}

	baseArr, bOk := base.(*data.ArrayValue)
	otherArr, oOk := other.(*data.ArrayValue)
	if !bOk || !oOk {
		return other
	}

	out := data.CloneArrayValue(baseArr)
	for i, z := range otherArr.List {
		if z == nil {
			continue
		}
		key := z.Name
		if key == "" {
			key = data.IntArrayKeyName(i)
		}
		setArrayNamedValue(out, key, z.Value)
	}
	return out
}

func setArrayNamedValue(arr *data.ArrayValue, key string, value data.Value) {
	if existing, ok := arr.LookupZValByStringKey(key); ok && existing != nil {
		existing.Value = value
		return
	}
	arr.List = append(arr.List, data.NewNamedZVal(key, value))
}

func shallowCopyPreserveKeys(v data.Value) data.Value {
	switch val := v.(type) {
	case *data.ObjectValue:
		out := data.NewObjectValue()
		val.RangeProperties(func(k string, prop data.Value) bool {
			out.SetProperty(k, prop)
			return true
		})
		return out
	case *data.ArrayValue:
		return data.CloneArrayValue(val)
	}
	return v
}

func (fn *ArrayReplaceFunction) GetName() string { return "array_replace" }
var arrayReplaceFunctionGetParams = []data.GetValue{node.NewParameters(nil, "arrays", 0, nil, nil)}

func (fn *ArrayReplaceFunction) GetParams() []data.GetValue {
	return arrayReplaceFunctionGetParams
}
var arrayReplaceFunctionGetVariables = []data.Variable{node.NewVariable(nil, "arrays", 0, data.NewBaseType("array"))}

func (fn *ArrayReplaceFunction) GetVariables() []data.Variable {
	return arrayReplaceFunctionGetVariables
}
