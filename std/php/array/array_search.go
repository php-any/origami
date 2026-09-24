package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArraySearchFunction 实现 array_search 函数
// array_search(mixed $needle, array $haystack, bool $strict = false): int|string|false
//
// PHP：strict=true 时用 ===（对象/闭包比身份，标量等比类型与值）；
// strict=false 时用 == 松散比较。不可对闭包一律 AsString()=="Closure"，
// 否则 Livewire EventBus::off() 会误删第一个监听器（ExtendBlade）。
type ArraySearchFunction struct{}

func NewArraySearchFunction() data.FuncStmt {
	return &ArraySearchFunction{}
}

func (f *ArraySearchFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	needleValue, _ := ctx.GetIndexValue(0)
	haystackValue, _ := ctx.GetIndexValue(1)
	strictValue, _ := ctx.GetIndexValue(2)

	if needleValue == nil || haystackValue == nil {
		return data.NewBoolValue(false), nil
	}

	arr, ok := haystackValue.(*data.ArrayValue)
	if !ok {
		return data.NewBoolValue(false), nil
	}

	strict := false
	if strictValue != nil {
		if b, ok := strictValue.(*data.BoolValue); ok {
			strict = b.Value
		} else if b, ok := strictValue.(data.AsBool); ok {
			if bv, err := b.AsBool(); err == nil {
				strict = bv
			}
		}
	}

	for i, z := range arr.List {
		v := z.Value
		if v == nil {
			continue
		}
		match := false
		if strict {
			match = valuesIdentical(needleValue, v)
		} else {
			match = valuesLooseEqual(needleValue, v)
		}
		if match {
			if z.Name != "" {
				return data.NewStringValue(z.Name), nil
			}
			return data.NewIntValue(i), nil
		}
	}

	return data.NewBoolValue(false), nil
}

func valuesIdentical(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	switch av := a.(type) {
	case *data.FuncValue:
		bv, ok := b.(*data.FuncValue)
		return ok && av == bv
	case *data.BoundFuncValue:
		bv, ok := b.(*data.BoundFuncValue)
		return ok && av == bv
	case *data.ClassValue:
		bv, ok := b.(*data.ClassValue)
		return ok && av == bv
	case *data.ArrayValue:
		bv, ok := b.(*data.ArrayValue)
		return ok && av == bv
	case *data.StringValue:
		bv, ok := b.(*data.StringValue)
		return ok && av.Value == bv.Value
	case *data.IntValue:
		bv, ok := b.(*data.IntValue)
		return ok && av.Value == bv.Value
	case *data.FloatValue:
		bv, ok := b.(*data.FloatValue)
		return ok && av.Value == bv.Value
	case *data.BoolValue:
		bv, ok := b.(*data.BoolValue)
		return ok && av.Value == bv.Value
	case *data.NullValue:
		_, ok := b.(*data.NullValue)
		return ok
	default:
		return a == b
	}
}

func valuesLooseEqual(a, b data.Value) bool {
	if a == nil || b == nil {
		return a == b
	}
	// 对象/闭包：松散 == 在 PHP 中对对象也是身份比较
	switch a.(type) {
	case *data.FuncValue, *data.BoundFuncValue, *data.ClassValue, *data.ArrayValue:
		return valuesIdentical(a, b)
	}
	switch b.(type) {
	case *data.FuncValue, *data.BoundFuncValue, *data.ClassValue, *data.ArrayValue:
		return valuesIdentical(a, b)
	}
	return data.Compare(a, b) == 0
}

func (f *ArraySearchFunction) GetName() string {
	return "array_search"
}

var arraySearchFunctionGetParams = []data.GetValue{
	node.NewParameter(nil, "needle", 0, nil, nil),
	node.NewParameter(nil, "haystack", 1, nil, nil),
	node.NewParameter(nil, "strict", 2, data.NewBoolValue(false), nil),
}

func (f *ArraySearchFunction) GetParams() []data.GetValue {
	return arraySearchFunctionGetParams
}

var arraySearchFunctionGetVariables = []data.Variable{
	node.NewVariable(nil, "needle", 0, nil),
	node.NewVariable(nil, "haystack", 1, data.NewBaseType("array")),
	node.NewVariable(nil, "strict", 2, data.NewBaseType("bool")),
}

func (f *ArraySearchFunction) GetVariables() []data.Variable {
	return arraySearchFunctionGetVariables
}
