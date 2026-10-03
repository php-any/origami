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
	a, ok := base.(*data.ArrayValue)
	if !ok {
		return other
	}
	b, ok := other.(*data.ArrayValue)
	if !ok {
		return other
	}
	out := data.CloneArrayValue(a)
	for slots, i := b.View(), 0; i < slots.Len(); i++ {
		slot := slots.At(i)
		key := slot.PHPArrayKey(i)
		if slot.RefCount() > 0 {
			out.BindReference(key, slot)
		} else {
			out.SetKey(key, slot.ReadValue())
		}
	}
	return out
}

func setArrayNamedValue(arr *data.ArrayValue, key string, value data.Value) {
	if existing, ok := arr.LookupZValByStringKey(key); ok && existing != nil {
		existing.StoreRaw(value)
		return
	}
	arr.SetStringKey(key, value)
}

func shallowCopyPreserveKeys(v data.Value) data.Value {
	switch val := v.(type) {

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
