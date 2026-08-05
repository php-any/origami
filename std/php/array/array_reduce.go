package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayReduceFunction 实现 array_reduce 函数
// 完整签名：array_reduce(array $array, callable $callback, mixed $initial = null): mixed
// 对数组中的每个元素应用 callback 函数，并将上次的结果和当前元素作为参数传递
type ArrayReduceFunction struct{}

func NewArrayReduceFunction() data.FuncStmt {
	return &ArrayReduceFunction{}
}

func (f *ArrayReduceFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	arrayVal, _ := ctx.GetIndexValue(0)
	callbackVal, _ := ctx.GetIndexValue(1)
	initialVal, hasInitial := ctx.GetIndexValue(2)

	// 如果没有提供初始值，使用 null
	var carry data.Value
	if hasInitial && initialVal != nil {
		carry = initialVal
	} else {
		carry = data.NewNullValue()
	}

	if arrayVal == nil {
		return carry, nil
	}

	if callbackVal == nil {
		return carry, nil
	}

	// 调用回调函数（支持 FuncValue；静态闭包等）
	callCallback := func(accumulator, item data.Value) (data.Value, data.Control) {
		ret, ctl := invokeCallback(ctx, callbackVal, []data.Value{accumulator, item})
		if ctl != nil {
			return nil, ctl
		}
		if ret == nil {
			return data.NewNullValue(), nil
		}
		return ret, nil
	}

	// 遍历数组：兼容 ArrayValue 列表与 ObjectValue 关联数组
	// （iterator_to_array 在 use_keys=true 时常返回 ObjectValue）
	for _, e := range toKVEntries(arrayVal) {
		var ctl data.Control
		carry, ctl = callCallback(carry, e.value)
		if ctl != nil {
			return nil, ctl
		}
	}

	return carry, nil
}

func (f *ArrayReduceFunction) GetName() string { return "array_reduce" }

func (f *ArrayReduceFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "array", 0, nil, nil),
		node.NewParameter(nil, "callback", 1, nil, nil),
		node.NewParameter(nil, "initial", 2, node.NewNullLiteral(nil), nil),
	}
}

func (f *ArrayReduceFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "array", 0, data.Mixed{}),
		node.NewVariable(nil, "callback", 1, data.Mixed{}),
		node.NewVariable(nil, "initial", 2, data.Mixed{}),
	}
}
