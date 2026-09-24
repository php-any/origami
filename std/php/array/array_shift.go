package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayShiftFunction 实现 array_shift 函数。
// PHP：array_shift 按引用改数组；赋值 `$b = $a` 是 copy-on-write，
// 不能把 $a 一并改掉（Symfony ArgvInput：`$this->parsed = $this->tokens` 后 shift parsed）。
type ArrayShiftFunction struct{}

func NewArrayShiftFunction() data.FuncStmt {
	return &ArrayShiftFunction{}
}

func (f *ArrayShiftFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayValue := data.CowSeparateIndex(ctx, 0)

	if arr, ok := arrayValue.(*data.ArrayValue); ok {
		if len(arr.List) == 0 {
			return data.NewNullValue(), nil
		}
		first := arr.List[0].Value
		arr.List = arr.List[1:]
		if first == nil {
			return data.NewNullValue(), nil
		}
		return first, nil
	}

	if obj, ok := arrayValue.(*data.ObjectValue); ok {
		var firstKey string
		var firstVal data.Value
		found := false
		obj.RangeProperties(func(key string, value data.Value) bool {
			firstKey = key
			firstVal = value
			found = true
			return false
		})
		if !found {
			return data.NewNullValue(), nil
		}
		obj.UnsetProperty(firstKey)
		if firstVal == nil {
			return data.NewNullValue(), nil
		}
		return firstVal, nil
	}

	return data.NewNullValue(), nil
}

func (f *ArrayShiftFunction) GetName() string {
	return "array_shift"
}

var arrayShiftFunctionGetParams = []data.GetValue{
	node.NewParameterReference(nil, "array", 0, nil, data.NewBaseType("array")),
}

func (f *ArrayShiftFunction) GetParams() []data.GetValue {
	return arrayShiftFunctionGetParams
}

var arrayShiftFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "array", 0, data.NewBaseType("array")),
}

func (f *ArrayShiftFunction) GetVariables() []data.Variable {
	return arrayShiftFunctionGetVariables
}
