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

	// 转换为数组
	arr, ok := arrayVal.(*data.ArrayValue)
	if !ok {
		return carry, nil
	}

	if callbackVal == nil {
		return carry, nil
	}

	// 调用回调函数
	callCallback := func(accumulator, item data.Value) (data.Value, data.Control) {
		switch cb := callbackVal.(type) {
		case *data.FuncValue:
			vars := cb.Value.GetVariables()
			fnCtx := ctx.CreateContext(vars)
			// 参数：$carry（累积值）, $item（当前元素）
			if len(vars) > 0 {
				fnCtx.SetIndexZVal(0, data.NewZVal(accumulator))
			}
			if len(vars) > 1 {
				fnCtx.SetIndexZVal(1, data.NewZVal(item))
			}
			ret, ctl := cb.Call(fnCtx)
			if ctl != nil {
				if rv, ok := ctl.(data.ReturnControl); ok {
					if v, ok := rv.ReturnValue().(data.Value); ok {
						return v, nil
					}
					return data.NewNullValue(), nil
				}
				return nil, ctl
			}
			if ret == nil {
				return data.NewNullValue(), nil
			}
			if v, ok := ret.(data.Value); ok {
				return v, nil
			}
			return data.NewNullValue(), nil
		default:
			return accumulator, nil
		}
	}

	// 遍历数组
	for _, item := range arr.List {
		var ctl data.Control
		carry, ctl = callCallback(carry, item.Value)
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
