package array

import (
	"github.com/php-any/origami/data"
	"github.com/php-any/origami/node"
)

// ArrayMapFunction 实现 PHP 内置函数 array_map
// array_map(callable $callback, array $array, array ...$arrays): array
// 当前实现支持：
// - $callback 为匿名函数 / Closure（*data.FuncValue）或实现 CallableValue 的对象
// - 至少一个数组参数
// - 多个数组时，按 PHP 语义以最短数组长度为准
type ArrayMapFunction struct{}

func NewArrayMapFunction() data.FuncStmt {
	return &ArrayMapFunction{}
}

func (f *ArrayMapFunction) Call(ctx data.Context) (data.GetValue, data.Control) {
	// 第一个参数：callback
	cbVal, has := ctx.GetIndexValue(0)
	if !has || cbVal == nil {
		// PHP 中 callback 省略会导致 TypeError，这里简单返回空数组
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 收集所有数组参数
	// 注意：签名里第二个参数使用了 node.NewParameters("arrays", 1, ...)，
	// 因此索引 1 上拿到的是一个 Parameters 聚合数组，里面按顺序存放所有实际传入的数组参数。
	var arrays []*data.ArrayValue
	paramsVal, _ := ctx.GetIndexValue(1)
	if paramsVal != nil {
		if paramsArr, ok := paramsVal.(*data.ArrayValue); ok {
			for _, z := range paramsArr.List {
				if z.Value == nil {
					continue
				}
				if av, ok := z.Value.(*data.ArrayValue); ok {
					arrays = append(arrays, av)
				} else {
					// 非数组参数：按照 PHP 行为，这里简化为忽略
				}
			}
		}
	}

	if len(arrays) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 计算最短数组长度
	minLen := len(arrays[0].List)
	for _, a := range arrays[1:] {
		if len(a.List) < minLen {
			minLen = len(a.List)
		}
	}
	if minLen == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	results := make([]data.Value, 0, minLen)

	for i := 0; i < minLen; i++ {
		// 收集本次调用的参数：各数组在索引 i 的元素
		args := make([]data.Value, 0, len(arrays))
		for _, arr := range arrays {
			if i < len(arr.List) {
				args = append(args, arr.List[i].Value)
			}
		}

		switch cb := cbVal.(type) {
		case *data.FuncValue:
			// 与 ArrayValueMap / forEach 一致的闭包调用方式：
			// 使用函数定义的变量创建上下文，并按顺序写入实参。
			vars := cb.Value.GetVariables()
			fnCtx := ctx.CreateContext(vars)
			for ai := 0; ai < len(vars) && ai < len(args); ai++ {
				fnCtx.SetVariableValue(data.NewVariable("", ai, nil), args[ai])
			}
			ret, ctl := cb.Value.Call(fnCtx)
			if ctl != nil {
				return nil, ctl
			}
			if v, ok := ret.(data.Value); ok {
				results = append(results, v)
			} else {
				results = append(results, data.NewNullValue())
			}
		case data.CallableValue:
			// 若实现了 CallableValue，则直接调用其 Call 接口。
			// 仅支持最多三个参数（与 ArrayValueMap/ForEach 一致的约定）。
			var arg0, arg1, arg2 data.Value = data.NewNullValue(), data.NewNullValue(), data.NewNullValue()
			if len(args) > 0 {
				arg0 = args[0]
			}
			if len(args) > 1 {
				arg1 = args[1]
			}
			if len(args) > 2 {
				arg2 = args[2]
			}
			ret, ctl := cb.Call(arg0, arg1, arg2)
			if ctl != nil {
				return nil, ctl
			}
			results = append(results, ret)
		default:
			// 非可调用，当前简单退化为 null；更严格行为可改为抛异常。
			results = append(results, data.NewNullValue())
		}
	}

	return data.NewArrayValue(results), nil
}

func (f *ArrayMapFunction) GetName() string {
	return "array_map"
}

func (f *ArrayMapFunction) GetParams() []data.GetValue {
	return []data.GetValue{
		node.NewParameter(nil, "callback", 0, nil, nil),
		node.NewParameters(nil, "arrays", 1, nil, nil),
	}
}

func (f *ArrayMapFunction) GetVariables() []data.Variable {
	return []data.Variable{
		node.NewVariable(nil, "callback", 0, data.Mixed{}),
		node.NewVariable(nil, "arrays", 1, data.Mixed{}),
	}
}
