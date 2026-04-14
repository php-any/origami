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
	type arrayInput struct {
		values []*data.ZVal // 按顺序的值列表
		keys   []string     // ObjectValue 的键列表（若为 nil 则表示索引数组）
	}
	var inputs []arrayInput
	var firstIsObject bool

	paramsVal, _ := ctx.GetIndexValue(1)
	if paramsVal != nil {
		if paramsArr, ok := paramsVal.(*data.ArrayValue); ok {
			for idx, z := range paramsArr.List {
				if z.Value == nil {
					continue
				}
				switch v := z.Value.(type) {
				case *data.ArrayValue:
					inputs = append(inputs, arrayInput{values: v.List, keys: nil})
				case *data.ObjectValue:
					// ObjectValue 是关联数组，提取键值对并保留顺序
					var items []*data.ZVal
					var keys []string
					v.RangeProperties(func(key string, val data.Value) bool {
						items = append(items, &data.ZVal{Value: val})
						keys = append(keys, key)
						return true
					})
					inputs = append(inputs, arrayInput{values: items, keys: keys})
					if idx == 0 {
						firstIsObject = true
					}
				default:
					// 非数组参数：按照 PHP 行为，这里简化为忽略
				}
			}
		}
	}

	if len(inputs) == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	// 计算最短数组长度
	minLen := len(inputs[0].values)
	for _, inp := range inputs[1:] {
		if len(inp.values) < minLen {
			minLen = len(inp.values)
		}
	}
	if minLen == 0 {
		return data.NewArrayValue([]data.Value{}), nil
	}

	results := make([]data.Value, 0, minLen)

	for i := 0; i < minLen; i++ {
		// 收集本次调用的参数：各数组在索引 i 的元素
		args := make([]data.Value, 0, len(inputs))
		for _, inp := range inputs {
			if i < len(inp.values) {
				args = append(args, inp.values[i].Value)
			}
		}

		var result data.Value
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
				result = v
			} else {
				result = data.NewNullValue()
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
			result = ret
		default:
			// 非可调用，当前简单退化为 null；更严格行为可改为抛异常。
			result = data.NewNullValue()
		}
		results = append(results, result)
	}

	// 如果第一个输入是 ObjectValue（关联数组），则保留原始键名
	if firstIsObject && inputs[0].keys != nil {
		obj := data.NewObjectValue()
		for i, v := range results {
			if i < len(inputs[0].keys) {
				obj.SetProperty(inputs[0].keys[i], v)
			}
		}
		return obj, nil
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
