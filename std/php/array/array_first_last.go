package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayFirstFunction 实现 PHP 8.5 array_first(array $array): mixed
// 返回数组第一个值；空数组返回 null。
type ArrayFirstFunction struct{}

func NewArrayFirstFunction() data.FuncStmt {
	return &ArrayFirstFunction{}
}

func (f *ArrayFirstFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewNullValue(), nil
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		if len(v.List) == 0 {
			return data.NewNullValue(), nil
		}
		return v.List[0].Value, nil
	case *data.ObjectValue:
		var first data.Value
		found := false
		v.RangeProperties(func(_ string, value data.Value) bool {
			first = value
			found = true
			return false
		})
		if !found {
			return data.NewNullValue(), nil
		}
		return first, nil
	default:
		return data.NewNullValue(), nil
	}
}

func (f *ArrayFirstFunction) GetName() string { return "array_first" }

func (f *ArrayFirstFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	}
}

func (f *ArrayFirstFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Arrays{}),
	}
}

// ArrayLastFunction 实现 PHP 8.5 array_last(array $array): mixed
// 返回数组最后一个值；空数组返回 null。
type ArrayLastFunction struct{}

func NewArrayLastFunction() data.FuncStmt {
	return &ArrayLastFunction{}
}

func (f *ArrayLastFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	val, _ := ctx.GetIndexValue(0)
	if val == nil {
		return data.NewNullValue(), nil
	}

	switch v := val.(type) {
	case *data.ArrayValue:
		if len(v.List) == 0 {
			return data.NewNullValue(), nil
		}
		return v.List[len(v.List)-1].Value, nil
	case *data.ObjectValue:
		var last data.Value
		found := false
		v.RangeProperties(func(_ string, value data.Value) bool {
			last = value
			found = true
			return true
		})
		if !found {
			return data.NewNullValue(), nil
		}
		return last, nil
	default:
		return data.NewNullValue(), nil
	}
}

func (f *ArrayLastFunction) GetName() string { return "array_last" }

func (f *ArrayLastFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, data.Arrays{}),
	}
}

func (f *ArrayLastFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Arrays{}),
	}
}
